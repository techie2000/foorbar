package service

import (
	"archive/zip"
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/techie2000/axiom/internal/domain"
	"github.com/techie2000/axiom/internal/repository"
)

// GLEIF URLs for downloading LEI data
const (
	// Level 1 data (who is who) - full file published daily (JSON format)
	GLEIFLevel1FullJSONURL = "https://goldencopy.gleif.org/api/v2/golden-copies/publishes/lei2-json/latest/download"
	GLEIFLevel1DeltaJSONURL = "https://goldencopy.gleif.org/api/v2/golden-copies/publishes/lei2-delta-json/latest/download"
)

// LEIService interface
type LEIService interface {
	// File download and management
	DownloadFullFile() (*domain.SourceFile, error)
	DownloadDeltaFile() (*domain.SourceFile, error)
	
	// File processing
	ProcessSourceFile(sourceFileID uuid.UUID) error
	ProcessSourceFileWithResume(sourceFileID uuid.UUID, resumeFromLEI string) error
	
	// Record management
	CreateLEIRecord(record *domain.LEIRecord) error
	GetLEIByCode(lei string) (*domain.LEIRecord, error)
	GetLEIByID(id string) (*domain.LEIRecord, error)
	GetAllLEI(limit, offset int) ([]*domain.LEIRecord, error)
	UpdateLEIRecord(record *domain.LEIRecord) error
	
	// Audit and history
	GetAuditHistory(lei string, limit int) ([]*domain.LEIRecordAudit, error)
	
	// Processing status
	GetProcessingStatus(jobType string) (*domain.FileProcessingStatus, error)
	UpdateProcessingStatus(status *domain.FileProcessingStatus) error
}

type leiService struct {
	repo repository.LEIRepository
	dataDir string // Directory to store downloaded files
}

// NewLEIService creates a new LEI service
func NewLEIService(repo repository.LEIRepository, dataDir string) LEIService {
	return &leiService{
		repo: repo,
		dataDir: dataDir,
	}
}

// DownloadFullFile downloads the full LEI data file from GLEIF
func (s *leiService) DownloadFullFile() (*domain.SourceFile, error) {
	return s.downloadFile(GLEIFLevel1FullJSONURL, "FULL")
}

// DownloadDeltaFile downloads the delta LEI data file from GLEIF
func (s *leiService) DownloadDeltaFile() (*domain.SourceFile, error) {
	return s.downloadFile(GLEIFLevel1DeltaJSONURL, "DELTA")
}

// downloadFile downloads a file from GLEIF and creates a SourceFile record
func (s *leiService) downloadFile(url, fileType string) (*domain.SourceFile, error) {
	log.Info().Str("url", url).Str("type", fileType).Msg("Starting file download from GLEIF")
	
	// Create data directory if it doesn't exist
	if err := os.MkdirAll(s.dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}
	
	// Download file
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download file: HTTP %d", resp.StatusCode)
	}
	
	// Generate filename with timestamp
	timestamp := time.Now().Format("20060102-150405")
	fileName := fmt.Sprintf("lei-%s-%s.json.zip", fileType, timestamp)
	filePath := filepath.Join(s.dataDir, fileName)
	
	// Create file
	out, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()
	
	// Calculate hash while downloading
	hash := sha256.New()
	multiWriter := io.MultiWriter(out, hash)
	
	// Copy data
	fileSize, err := io.Copy(multiWriter, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}
	
	fileHash := hex.EncodeToString(hash.Sum(nil))
	
	log.Info().
		Str("file", fileName).
		Int64("size", fileSize).
		Str("hash", fileHash).
		Msg("File downloaded successfully")
	
	// Create SourceFile record
	sourceFile := &domain.SourceFile{
		FileName:        fileName,
		FileType:        fileType,
		FileURL:         url,
		FileSize:        fileSize,
		FileHash:        fileHash,
		DownloadedAt:    time.Now(),
		PublicationDate: time.Now(), // TODO: Parse from GLEIF metadata if available
		ProcessingStatus: "PENDING",
	}
	
	if err := s.repo.CreateSourceFile(sourceFile); err != nil {
		return nil, fmt.Errorf("failed to create source file record: %w", err)
	}
	
	return sourceFile, nil
}

// ProcessSourceFile processes a downloaded source file
func (s *leiService) ProcessSourceFile(sourceFileID uuid.UUID) error {
	return s.ProcessSourceFileWithResume(sourceFileID, "")
}

// ProcessSourceFileWithResume processes a source file, optionally resuming from a specific LEI
func (s *leiService) ProcessSourceFileWithResume(sourceFileID uuid.UUID, resumeFromLEI string) error {
	log.Info().Str("source_file_id", sourceFileID.String()).Str("resume_from", resumeFromLEI).Msg("Starting file processing")
	
	// Get source file
	sourceFile, err := s.repo.FindSourceFileByID(sourceFileID.String())
	if err != nil {
		return fmt.Errorf("failed to find source file: %w", err)
	}
	
	// Update status to IN_PROGRESS
	sourceFile.ProcessingStatus = "IN_PROGRESS"
	startTime := time.Now()
	sourceFile.ProcessingStartedAt = &startTime
	if err := s.repo.UpdateSourceFile(sourceFile); err != nil {
		return fmt.Errorf("failed to update source file status: %w", err)
	}
	
	// Extract and process file
	filePath := filepath.Join(s.dataDir, sourceFile.FileName)
	
	// Unzip file
	jsonPath, err := s.extractZipFile(filePath)
	if err != nil {
		sourceFile.ProcessingStatus = "FAILED"
		sourceFile.ProcessingError = err.Error()
		s.repo.UpdateSourceFile(sourceFile)
		return fmt.Errorf("failed to extract file: %w", err)
	}
	defer os.Remove(jsonPath) // Clean up extracted JSON
	
	// Parse and process JSON
	if err := s.processJSONFile(jsonPath, sourceFile, resumeFromLEI); err != nil {
		sourceFile.ProcessingStatus = "FAILED"
		sourceFile.ProcessingError = err.Error()
		s.repo.UpdateSourceFile(sourceFile)
		return fmt.Errorf("failed to process JSON file: %w", err)
	}
	
	// Update status to COMPLETED
	sourceFile.ProcessingStatus = "COMPLETED"
	completedTime := time.Now()
	sourceFile.ProcessingCompletedAt = &completedTime
	if err := s.repo.UpdateSourceFile(sourceFile); err != nil {
		return fmt.Errorf("failed to update source file status: %w", err)
	}
	
	log.Info().
		Str("source_file_id", sourceFileID.String()).
		Int("total", sourceFile.TotalRecords).
		Int("processed", sourceFile.ProcessedRecords).
		Int("failed", sourceFile.FailedRecords).
		Msg("File processing completed")
	
	return nil
}

// extractZipFile extracts the JSON file from a ZIP archive
func (s *leiService) extractZipFile(zipPath string) (string, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", err
	}
	defer r.Close()
	
	// Find the JSON file in the ZIP
	for _, f := range r.File {
		if filepath.Ext(f.Name) == ".json" || filepath.Ext(f.Name) == ".jsonl" {
			rc, err := f.Open()
			if err != nil {
				return "", err
			}
			defer rc.Close()
			
			// Create output file
			jsonPath := zipPath + ".extracted.json"
			outFile, err := os.Create(jsonPath)
			if err != nil {
				return "", err
			}
			defer outFile.Close()
			
			// Copy content
			_, err = io.Copy(outFile, rc)
			if err != nil {
				return "", err
			}
			
			return jsonPath, nil
		}
	}
	
	return "", fmt.Errorf("no JSON file found in ZIP archive")
}

// processJSONFile parses and processes the LEI JSON file
// GLEIF JSON format is typically JSON Lines (one JSON object per line)
func (s *leiService) processJSONFile(jsonPath string, sourceFile *domain.SourceFile, resumeFromLEI string) error {
	file, err := os.Open(jsonPath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	// Increase buffer size for large JSON lines
	buf := make([]byte, 0, 1024*1024) // 1MB buffer
	scanner.Buffer(buf, 10*1024*1024) // 10MB max token size
	
	var totalRecords int
	var processedRecords int
	var failedRecords int
	var shouldProcess bool = (resumeFromLEI == "")
	
	// Process JSON lines (JSON Lines format - one record per line)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		
		totalRecords++
		
		var jsonRecord LEIJSONRecord
		if err := json.Unmarshal(line, &jsonRecord); err != nil {
			log.Error().Err(err).Msg("Failed to decode LEI JSON record")
			failedRecords++
			continue
		}
		
		// Check if we should start processing (resume logic)
		if !shouldProcess {
			if jsonRecord.LEI == resumeFromLEI {
				shouldProcess = true
			} else {
				continue
			}
		}
		
		// Convert JSON record to domain model
		record := s.jsonToDomainRecord(&jsonRecord, sourceFile.ID)
		
		// Upsert record (handles change detection)
		if _, err := s.repo.UpsertLEIRecord(record); err != nil {
			log.Error().Err(err).Str("lei", record.LEI).Msg("Failed to upsert LEI record")
			failedRecords++
		} else {
			processedRecords++
		}
		
		// Update progress every 1000 records
		if processedRecords%1000 == 0 {
			sourceFile.TotalRecords = totalRecords
			sourceFile.ProcessedRecords = processedRecords
			sourceFile.FailedRecords = failedRecords
			sourceFile.LastProcessedLEI = record.LEI
			s.repo.UpdateSourceFile(sourceFile)
			
			log.Info().
				Int("total", totalRecords).
				Int("processed", processedRecords).
				Int("failed", failedRecords).
				Str("last_lei", record.LEI).
				Msg("Processing progress")
		}
	}
	
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading JSON file: %w", err)
	}
	
	// Final update
	sourceFile.TotalRecords = totalRecords
	sourceFile.ProcessedRecords = processedRecords
	sourceFile.FailedRecords = failedRecords
	
	return nil
}

// LEIJSONRecord represents the JSON structure from GLEIF
// Based on GLEIF Level 1 JSON format
type LEIJSONRecord struct {
	LEI                     string                  `json:"LEI"`
	Entity                  LEIEntity               `json:"Entity"`
	Registration            LEIRegistration         `json:"Registration"`
}

type LEIEntity struct {
	LegalName               LEILegalName            `json:"LegalName"`
	OtherNames              []LEIOtherName          `json:"OtherEntityNames"`
	TransliteratedOtherNames []LEIOtherName         `json:"TransliteratedOtherEntityNames"`
	LegalAddress            LEIAddress              `json:"LegalAddress"`
	HeadquartersAddress     LEIAddress              `json:"HeadquartersAddress"`
	RegistrationAuthority   LEIRegistrationAuthority `json:"RegistrationAuthority"`
	LegalJurisdiction       string                  `json:"LegalJurisdiction"`
	EntityCategory          string                  `json:"EntityCategory"`
	LegalForm               LEILegalForm            `json:"LegalForm"`
	EntityStatus            string                  `json:"EntityStatus"`
}

type LEILegalName struct {
	Value                   string                  `json:"$"`
	Language                string                  `json:"@xml:lang"`
}

type LEIOtherName struct {
	Value                   string                  `json:"$"`
	Type                    string                  `json:"@type"`
}

type LEIAddress struct {
	FirstAddressLine        string                  `json:"FirstAddressLine"`
	AdditionalAddressLine   []string                `json:"AdditionalAddressLine"`
	City                    string                  `json:"City"`
	Region                  string                  `json:"Region"`
	Country                 string                  `json:"Country"`
	PostalCode              string                  `json:"PostalCode"`
}

type LEIRegistrationAuthority struct {
	RegistrationAuthorityID string                  `json:"RegistrationAuthorityID"`
	RegistrationAuthorityEntityID string            `json:"RegistrationAuthorityEntityID"`
}

type LEILegalForm struct {
	EntityLegalFormCode     string                  `json:"EntityLegalFormCode"`
}

type LEIRegistration struct {
	InitialRegistrationDate string                  `json:"InitialRegistrationDate"`
	LastUpdateDate          string                  `json:"LastUpdateDate"`
	RegistrationStatus      string                  `json:"RegistrationStatus"`
	NextRenewalDate         string                  `json:"NextRenewalDate"`
	ManagingLOU             string                  `json:"ManagingLOU"`
	ValidationSources       string                  `json:"ValidationSources"`
}

// jsonToDomainRecord converts a JSON record to a domain.LEIRecord
func (s *leiService) jsonToDomainRecord(jsonRecord *LEIJSONRecord, sourceFileID uuid.UUID) *domain.LEIRecord {
	record := &domain.LEIRecord{
		LEI:                   jsonRecord.LEI,
		LegalName:             jsonRecord.Entity.LegalName.Value,
		LegalAddressLine1:     jsonRecord.Entity.LegalAddress.FirstAddressLine,
		LegalAddressCity:      jsonRecord.Entity.LegalAddress.City,
		LegalAddressRegion:    jsonRecord.Entity.LegalAddress.Region,
		LegalAddressCountry:   jsonRecord.Entity.LegalAddress.Country,
		LegalAddressPostalCode: jsonRecord.Entity.LegalAddress.PostalCode,
		RegistrationAuthority: jsonRecord.Entity.RegistrationAuthority.RegistrationAuthorityID,
		RegistrationNumber:    jsonRecord.Entity.RegistrationAuthority.RegistrationAuthorityEntityID,
		EntityCategory:        jsonRecord.Entity.EntityCategory,
		EntityLegalForm:       jsonRecord.Entity.LegalForm.EntityLegalFormCode,
		EntityStatus:          jsonRecord.Registration.RegistrationStatus,
		ManagingLOU:           jsonRecord.Registration.ManagingLOU,
		SourceFileID:          &sourceFileID,
	}
	
	// Handle additional address lines
	if len(jsonRecord.Entity.LegalAddress.AdditionalAddressLine) > 0 {
		record.LegalAddressLine2 = jsonRecord.Entity.LegalAddress.AdditionalAddressLine[0]
	}
	if len(jsonRecord.Entity.LegalAddress.AdditionalAddressLine) > 1 {
		record.LegalAddressLine3 = jsonRecord.Entity.LegalAddress.AdditionalAddressLine[1]
	}
	if len(jsonRecord.Entity.LegalAddress.AdditionalAddressLine) > 2 {
		record.LegalAddressLine4 = jsonRecord.Entity.LegalAddress.AdditionalAddressLine[2]
	}
	
	// Handle headquarters address
	if jsonRecord.Entity.HeadquartersAddress.FirstAddressLine != "" {
		record.HQAddressLine1 = jsonRecord.Entity.HeadquartersAddress.FirstAddressLine
		record.HQAddressCity = jsonRecord.Entity.HeadquartersAddress.City
		record.HQAddressRegion = jsonRecord.Entity.HeadquartersAddress.Region
		record.HQAddressCountry = jsonRecord.Entity.HeadquartersAddress.Country
		record.HQAddressPostalCode = jsonRecord.Entity.HeadquartersAddress.PostalCode
		
		if len(jsonRecord.Entity.HeadquartersAddress.AdditionalAddressLine) > 0 {
			record.HQAddressLine2 = jsonRecord.Entity.HeadquartersAddress.AdditionalAddressLine[0]
		}
		if len(jsonRecord.Entity.HeadquartersAddress.AdditionalAddressLine) > 1 {
			record.HQAddressLine3 = jsonRecord.Entity.HeadquartersAddress.AdditionalAddressLine[1]
		}
		if len(jsonRecord.Entity.HeadquartersAddress.AdditionalAddressLine) > 2 {
			record.HQAddressLine4 = jsonRecord.Entity.HeadquartersAddress.AdditionalAddressLine[2]
		}
	}
	
	// Parse dates (ISO 8601 format)
	if jsonRecord.Registration.InitialRegistrationDate != "" {
		if t, err := time.Parse("2006-01-02T15:04:05Z", jsonRecord.Registration.InitialRegistrationDate); err == nil {
			record.InitialRegistrationDate = t
		} else if t, err := time.Parse("2006-01-02", jsonRecord.Registration.InitialRegistrationDate); err == nil {
			record.InitialRegistrationDate = t
		}
	}
	if jsonRecord.Registration.LastUpdateDate != "" {
		if t, err := time.Parse("2006-01-02T15:04:05Z", jsonRecord.Registration.LastUpdateDate); err == nil {
			record.LastUpdateDate = t
		} else if t, err := time.Parse("2006-01-02", jsonRecord.Registration.LastUpdateDate); err == nil {
			record.LastUpdateDate = t
		}
	}
	if jsonRecord.Registration.NextRenewalDate != "" {
		if t, err := time.Parse("2006-01-02T15:04:05Z", jsonRecord.Registration.NextRenewalDate); err == nil {
			record.NextRenewalDate = t
		} else if t, err := time.Parse("2006-01-02", jsonRecord.Registration.NextRenewalDate); err == nil {
			record.NextRenewalDate = t
		}
	}
	
	return record
}

// CreateLEIRecord creates a new LEI record
func (s *leiService) CreateLEIRecord(record *domain.LEIRecord) error {
	return s.repo.CreateLEIRecord(record)
}

// GetLEIByCode retrieves an LEI record by LEI code
func (s *leiService) GetLEIByCode(lei string) (*domain.LEIRecord, error) {
	return s.repo.FindLEIByLEI(lei)
}

// GetLEIByID retrieves an LEI record by ID
func (s *leiService) GetLEIByID(id string) (*domain.LEIRecord, error) {
	return s.repo.FindLEIByID(id)
}

// GetAllLEI retrieves all LEI records with pagination
func (s *leiService) GetAllLEI(limit, offset int) ([]*domain.LEIRecord, error) {
	return s.repo.FindAllLEI(limit, offset)
}

// UpdateLEIRecord updates an LEI record
func (s *leiService) UpdateLEIRecord(record *domain.LEIRecord) error {
	return s.repo.UpdateLEIRecord(record)
}

// GetAuditHistory retrieves audit history for an LEI
func (s *leiService) GetAuditHistory(lei string, limit int) ([]*domain.LEIRecordAudit, error) {
	return s.repo.FindAuditHistoryByLEI(lei, limit)
}

// GetProcessingStatus retrieves processing status for a job type
func (s *leiService) GetProcessingStatus(jobType string) (*domain.FileProcessingStatus, error) {
	return s.repo.FindProcessingStatus(jobType)
}

// UpdateProcessingStatus updates processing status
func (s *leiService) UpdateProcessingStatus(status *domain.FileProcessingStatus) error {
	return s.repo.UpdateProcessingStatus(status)
}
