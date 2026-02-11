package repository

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/techie2000/axiom/internal/domain"
	"gorm.io/gorm"
)

// LEIRepository interface
type LEIRepository interface {
	// LEI Record operations
	CreateLEIRecord(record *domain.LEIRecord) error
	FindLEIByLEI(lei string) (*domain.LEIRecord, error)
	FindLEIByID(id string) (*domain.LEIRecord, error)
	FindAllLEI(limit, offset int) ([]*domain.LEIRecord, error)
	CountLEIRecords() (int64, error)
	UpdateLEIRecord(record *domain.LEIRecord) error
	UpsertLEIRecord(record *domain.LEIRecord) (bool, error) // Returns true if updated, false if created
	BatchUpsertLEIRecords(records []*domain.LEIRecord) (int, int, error) // Returns (created, updated, error)
	DeleteLEI(id string) error

	// Source File operations
	CreateSourceFile(file *domain.SourceFile) error
	FindSourceFileByID(id string) (*domain.SourceFile, error)
	FindLatestSourceFile(fileType string) (*domain.SourceFile, error)
	UpdateSourceFile(file *domain.SourceFile) error
	FindPendingSourceFiles() ([]*domain.SourceFile, error)
	FindRetryableFailedFiles() ([]*domain.SourceFile, error)
	ResetFailedFileForRetry(fileID uuid.UUID) error

	// File Processing Status operations
	FindProcessingStatus(jobType string) (*domain.FileProcessingStatus, error)
	UpdateProcessingStatus(status *domain.FileProcessingStatus) error

	// Audit operations
	CreateAuditRecord(audit *domain.LEIRecordAudit) error
	FindAuditHistoryByLEI(lei string, limit int) ([]*domain.LEIRecordAudit, error)
}

type leiRepository struct {
	db *gorm.DB
}

// NewLEIRepository creates a new LEI repository instance
func NewLEIRepository(db *gorm.DB) LEIRepository {
	return &leiRepository{db: db}
}

// CreateLEIRecord creates a new LEI record
func (r *leiRepository) CreateLEIRecord(record *domain.LEIRecord) error {
	return r.db.Create(record).Error
}

// FindLEIByLEI finds an LEI record by LEI code
func (r *leiRepository) FindLEIByLEI(lei string) (*domain.LEIRecord, error) {
	var record domain.LEIRecord
	if err := r.db.Where("lei = ?", lei).Preload("SourceFile").First(&record).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

// FindLEIByID finds an LEI record by ID
func (r *leiRepository) FindLEIByID(id string) (*domain.LEIRecord, error) {
	var record domain.LEIRecord
	if err := r.db.Preload("SourceFile").First(&record, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

// FindAllLEI retrieves all LEI records with pagination
func (r *leiRepository) FindAllLEI(limit, offset int) ([]*domain.LEIRecord, error) {
	var records []*domain.LEIRecord
	if err := r.db.Limit(limit).Offset(offset).Preload("SourceFile").Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

// CountLEIRecords returns the total count of LEI records
func (r *leiRepository) CountLEIRecords() (int64, error) {
	var count int64
	if err := r.db.Model(&domain.LEIRecord{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// UpdateLEIRecord updates an existing LEI record
func (r *leiRepository) UpdateLEIRecord(record *domain.LEIRecord) error {
	return r.db.Save(record).Error
}

// UpsertLEIRecord creates or updates an LEI record with change detection
// Returns true if updated, false if created
func (r *leiRepository) UpsertLEIRecord(record *domain.LEIRecord) (bool, error) {
	existing, err := r.FindLEIByLEI(record.LEI)

	// If not found, create new record
	if err == gorm.ErrRecordNotFound {
		record.CreatedBy = "system"
		record.UpdatedBy = "system"
		if err := r.CreateLEIRecord(record); err != nil {
			return false, err
		}

		// Create audit record for creation
		auditRecord := &domain.LEIRecordAudit{
			LEIRecordID:    record.ID,
			LEI:            record.LEI,
			Action:         "CREATE",
			RecordSnapshot: r.recordToJSON(record),
			ChangedFields:  "{}",
			SourceFileID:   record.SourceFileID,
			ChangedBy:      "system",
		}
		if err := r.CreateAuditRecord(auditRecord); err != nil {
			return false, fmt.Errorf("failed to create audit record: %w", err)
		}

		return false, nil
	}

	if err != nil {
		return false, err
	}

	// Detect changes
	changes := r.detectChanges(existing, record)

	// If no changes detected, don't update
	if len(changes) == 0 {
		return false, nil
	}

	// Convert changes to JSON
	changesJSON, err := json.Marshal(changes)
	if err != nil {
		return false, fmt.Errorf("failed to marshal changes: %w", err)
	}

	// Update the record
	record.ID = existing.ID
	record.CreatedAt = existing.CreatedAt
	record.CreatedBy = existing.CreatedBy
	record.UpdatedBy = "system"
	record.ChangedFields = string(changesJSON)

	if err := r.UpdateLEIRecord(record); err != nil {
		return false, err
	}

	// Create audit record for update
	auditRecord := &domain.LEIRecordAudit{
		LEIRecordID:    record.ID,
		LEI:            record.LEI,
		Action:         "UPDATE",
		RecordSnapshot: r.recordToJSON(record),
		ChangedFields:  string(changesJSON),
		SourceFileID:   record.SourceFileID,
		ChangedBy:      "system",
	}
	if err := r.CreateAuditRecord(auditRecord); err != nil {
		return false, fmt.Errorf("failed to create audit record: %w", err)
	}

	return true, nil
}

// BatchUpsertLEIRecords performs batch upsert of LEI records for high-performance bulk imports
// Returns (created_count, updated_count, error)
func (r *leiRepository) BatchUpsertLEIRecords(records []*domain.LEIRecord) (int, int, error) {
	if len(records) == 0 {
		return 0, 0, nil
	}

	// Set created_by and updated_by for all records
	now := time.Now()
	leiCodes := make([]string, len(records))
	for i, record := range records {
		if record.CreatedAt.IsZero() {
			record.CreatedAt = now
		}
		if record.UpdatedAt.IsZero() {
			record.UpdatedAt = now
		}
		if record.CreatedBy == "" {
			record.CreatedBy = "system"
		}
		if record.UpdatedBy == "" {
			record.UpdatedBy = "system"
		}
		leiCodes[i] = record.LEI
	}

	// Count existing records before batch insert
	var existingCount int64
	if err := r.db.Model(&domain.LEIRecord{}).
		Where("lei IN ?", leiCodes).
		Count(&existingCount).Error; err != nil {
		return 0, 0, fmt.Errorf("failed to count existing records: %w", err)
	}

	// Use GORM's CreateInBatches with ON CONFLICT for reliable upsert
	// Note: GORM doesn't have built-in upsert, so we use raw SQL with ON CONFLICT
	// Build VALUES clause for PostgreSQL INSERT ... ON CONFLICT
	if len(records) > 0 {
		// Use transaction for atomicity
		tx := r.db.Begin()
		if tx.Error != nil {
			return 0, 0, tx.Error
		}
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		// Insert in batches of 100 for optimal performance
		batchSize := 100
		for i := 0; i < len(records); i += batchSize {
			end := i + batchSize
			if end > len(records) {
				end = len(records)
			}
			batch := records[i:end]

			// Build SQL with placeholders
			valueStrings := make([]string, 0, len(batch))
			valueArgs := make([]interface{}, 0, len(batch)*20) // Approx 20 fields per record

			for _, record := range batch {
				valueStrings = append(valueStrings, "(gen_random_uuid(), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW(), 'system', 'system', '{}')")
				valueArgs = append(valueArgs,
					record.LEI,
					record.LegalName,
					record.EntityStatus,
					record.LegalAddressLine1,
					record.LegalAddressCity,
					record.LegalAddressRegion,
					record.LegalAddressCountry,
					record.LegalAddressPostalCode,
					record.HQAddressLine1,
					record.HQAddressCity,
					record.HQAddressRegion,
					record.HQAddressCountry,
					record.HQAddressPostalCode,
					record.InitialRegistrationDate,
					record.LastUpdateDate,
					record.NextRenewalDate,
					record.ManagingLOU,
					record.ValidationSources,
					record.EntityCategory,
					record.EntitySubCategory,
					record.SourceFileID,
				)
			}

			stmt := fmt.Sprintf(`
				INSERT INTO lei_raw.lei_records (
					id, lei, legal_name, entity_status,
					legal_address_line_1, legal_address_city, legal_address_region,
					legal_address_country, legal_address_postal_code,
					hq_address_line_1, hq_address_city, hq_address_region,
					hq_address_country, hq_address_postal_code,
					initial_registration_date, last_update_date, next_renewal_date,
					managing_lou, validation_sources, entity_category, entity_sub_category,
					source_file_id,
					created_at, updated_at, created_by, updated_by, changed_fields
				) VALUES %s
				ON CONFLICT (lei) DO UPDATE SET
					legal_name = EXCLUDED.legal_name,
					entity_status = EXCLUDED.entity_status,
					legal_address_line_1 = EXCLUDED.legal_address_line_1,
					legal_address_city = EXCLUDED.legal_address_city,
					legal_address_region = EXCLUDED.legal_address_region,
					legal_address_country = EXCLUDED.legal_address_country,
					legal_address_postal_code = EXCLUDED.legal_address_postal_code,
					hq_address_line_1 = EXCLUDED.hq_address_line_1,
					hq_address_city = EXCLUDED.hq_address_city,
					hq_address_region = EXCLUDED.hq_address_region,
					hq_address_country = EXCLUDED.hq_address_country,
					hq_address_postal_code = EXCLUDED.hq_address_postal_code,
					initial_registration_date = EXCLUDED.initial_registration_date,
					last_update_date = EXCLUDED.last_update_date,
					next_renewal_date = EXCLUDED.next_renewal_date,
					managing_lou = EXCLUDED.managing_lou,
					validation_sources = EXCLUDED.validation_sources,
					entity_category = EXCLUDED.entity_category,
					entity_sub_category = EXCLUDED.entity_sub_category,
					source_file_id = EXCLUDED.source_file_id,
					updated_at = NOW(),
					updated_by = 'system'
			`, strings.Join(valueStrings, ","))

			if err := tx.Exec(stmt, valueArgs...).Error; err != nil {
				tx.Rollback()
				// Log detailed error information
				log.Error().
					Err(err).
					Int("batch_start", i).
					Int("batch_end", end).
					Int("batch_size", len(batch)).
					Str("first_lei", batch[0].LEI).
					Str("last_lei", batch[len(batch)-1].LEI).
					Int("value_args_count", len(valueArgs)).
					Msg("CRITICAL: Batch upsert SQL execution failed")
				return 0, 0, fmt.Errorf("failed to batch upsert records %d-%d: %w", i, end, err)
			}

			log.Debug().
				Int("batch_start", i).
				Int("batch_end", end).
				Int("batch_size", len(batch)).
				Msg("Batch SQL executed successfully")
		}

		if err := tx.Commit().Error; err != nil {
			return 0, 0, fmt.Errorf("failed to commit transaction: %w", err)
		}
	}

	totalCount := int64(len(records))
	createdCount := totalCount - existingCount
	updatedCount := existingCount

	return int(createdCount), int(updatedCount), nil
}

// DeleteLEI soft deletes an LEI record
func (r *leiRepository) DeleteLEI(id string) error {
	// Get the record before deleting for audit
	record, err := r.FindLEIByID(id)
	if err != nil {
		return err
	}

	// Soft delete
	if err := r.db.Delete(&domain.LEIRecord{}, "id = ?", id).Error; err != nil {
		return err
	}

	// Create audit record for deletion
	auditRecord := &domain.LEIRecordAudit{
		LEIRecordID:    record.ID,
		LEI:            record.LEI,
		Action:         "DELETE",
		RecordSnapshot: r.recordToJSON(record),
		ChangedFields:  "{}",
		ChangedBy:      "system",
	}
	return r.CreateAuditRecord(auditRecord)
}

// CreateSourceFile creates a new source file record
func (r *leiRepository) CreateSourceFile(file *domain.SourceFile) error {
	return r.db.Create(file).Error
}

// FindSourceFileByID finds a source file by ID
func (r *leiRepository) FindSourceFileByID(id string) (*domain.SourceFile, error) {
	var file domain.SourceFile
	if err := r.db.First(&file, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

// FindLatestSourceFile finds the latest source file of a given type
func (r *leiRepository) FindLatestSourceFile(fileType string) (*domain.SourceFile, error) {
	var file domain.SourceFile
	if err := r.db.Where("file_type = ?", fileType).Order("publication_date DESC").First(&file).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

// UpdateSourceFile updates a source file record
func (r *leiRepository) UpdateSourceFile(file *domain.SourceFile) error {
	return r.db.Save(file).Error
}

// FindPendingSourceFiles finds all source files pending processing
func (r *leiRepository) FindPendingSourceFiles() ([]*domain.SourceFile, error) {
	var files []*domain.SourceFile
	if err := r.db.Where("processing_status IN ?", []string{"PENDING", "IN_PROGRESS"}).
		Order("publication_date ASC").
		Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}

// FindRetryableFailedFiles finds FAILED files that are eligible for retry
func (r *leiRepository) FindRetryableFailedFiles() ([]*domain.SourceFile, error) {
	var files []*domain.SourceFile
	if err := r.db.Where("processing_status = ? AND retry_count < max_retries", "FAILED").
		Where("failure_category IN ? OR failure_category IS NULL", []string{"SCHEMA_ERROR", "NETWORK_ERROR", "UNKNOWN"}).
		Order("publication_date ASC").
		Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}

// ResetFailedFileForRetry resets a failed file to PENDING for retry
func (r *leiRepository) ResetFailedFileForRetry(fileID uuid.UUID) error {
	return r.db.Model(&domain.SourceFile{}).
		Where("id = ?", fileID).
		Updates(map[string]interface{}{
			"processing_status": "PENDING",
			"retry_count":        gorm.Expr("retry_count + 1"),
			"processing_error":   "",
		}).Error
}

// FindProcessingStatus finds the processing status for a job type
func (r *leiRepository) FindProcessingStatus(jobType string) (*domain.FileProcessingStatus, error) {
	var status domain.FileProcessingStatus
	if err := r.db.Where("job_type = ?", jobType).Preload("CurrentSourceFile").First(&status).Error; err != nil {
		return nil, err
	}
	return &status, nil
}

// UpdateProcessingStatus updates the processing status
func (r *leiRepository) UpdateProcessingStatus(status *domain.FileProcessingStatus) error {
	return r.db.Save(status).Error
}

// CreateAuditRecord creates a new audit record
func (r *leiRepository) CreateAuditRecord(audit *domain.LEIRecordAudit) error {
	return r.db.Create(audit).Error
}

// FindAuditHistoryByLEI retrieves audit history for an LEI
func (r *leiRepository) FindAuditHistoryByLEI(lei string, limit int) ([]*domain.LEIRecordAudit, error) {
	var audits []*domain.LEIRecordAudit
	query := r.db.Where("lei = ?", lei).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&audits).Error; err != nil {
		return nil, err
	}
	return audits, nil
}

// detectChanges compares two LEI records and returns a map of changed fields
func (r *leiRepository) detectChanges(old, new *domain.LEIRecord) map[string]domain.LEIChangeDetection {
	changes := make(map[string]domain.LEIChangeDetection)

	oldVal := reflect.ValueOf(*old)
	newVal := reflect.ValueOf(*new)
	oldType := oldVal.Type()

	for i := 0; i < oldVal.NumField(); i++ {
		field := oldType.Field(i)
		fieldName := field.Name

		// Skip internal fields and timestamps
		if fieldName == "ID" || fieldName == "CreatedAt" || fieldName == "UpdatedAt" ||
			fieldName == "DeletedAt" || fieldName == "CreatedBy" || fieldName == "UpdatedBy" ||
			fieldName == "ChangedFields" || fieldName == "SourceFile" || fieldName == "SourceFileID" {
			continue
		}

		oldFieldVal := oldVal.Field(i).Interface()
		newFieldVal := newVal.Field(i).Interface()

		// Compare values
		if !reflect.DeepEqual(oldFieldVal, newFieldVal) {
			// Special handling for time.Time zero values
			if field.Type == reflect.TypeOf(time.Time{}) {
				oldTime := oldFieldVal.(time.Time)
				newTime := newFieldVal.(time.Time)
				if oldTime.IsZero() && newTime.IsZero() {
					continue
				}
			}

			changes[fieldName] = domain.LEIChangeDetection{
				FieldName: fieldName,
				OldValue:  oldFieldVal,
				NewValue:  newFieldVal,
			}
		}
	}

	return changes
}

// recordToJSON converts an LEI record to JSON string
func (r *leiRepository) recordToJSON(record *domain.LEIRecord) string {
	jsonBytes, err := json.Marshal(record)
	if err != nil {
		return "{}"
	}
	return string(jsonBytes)
}
