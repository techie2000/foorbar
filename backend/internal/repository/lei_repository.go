package repository

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

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
	UpdateLEIRecord(record *domain.LEIRecord) error
	UpsertLEIRecord(record *domain.LEIRecord) (bool, error) // Returns true if updated, false if created
	DeleteLEI(id string) error

	// Source File operations
	CreateSourceFile(file *domain.SourceFile) error
	FindSourceFileByID(id string) (*domain.SourceFile, error)
	FindLatestSourceFile(fileType string) (*domain.SourceFile, error)
	UpdateSourceFile(file *domain.SourceFile) error
	FindPendingSourceFiles() ([]*domain.SourceFile, error)

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
