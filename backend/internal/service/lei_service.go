package service

import (
	"archive/zip"
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

// GLEIF API endpoints and data directory configuration
const (
	// Base URL for GLEIF golden copy API
	GLEIFBaseURL = "https://goldencopy.gleif.org"

	// Discovery endpoint to get latest file URLs
	// Returns bulk file metadata - format differs from single LEI API queries
	GLEIFLatestPublishesURL = "https://goldencopy.gleif.org/api/v2/golden-copies/publishes/latest"

	// Data directory for downloaded files (relative to working directory)
	DefaultDataDirectory = "./data/lei"
)

// GLEIFPublishesResponse represents the response from the GLEIF latest publishes endpoint
type GLEIFPublishesResponse struct {
	Data GLEIFPublishesData `json:"data"`
}

type GLEIFPublishesData struct {
	LEI2 GLEIFFileFormats `json:"lei2"`
}

type GLEIFFileFormats struct {
	Type        string               `json:"type"`
	PublishDate string               `json:"publish_date"`
	FullFile    GLEIFFileGroup       `json:"full_file"`
	DeltaFiles  GLEIFDeltaFileGroups `json:"delta_files"`
}

type GLEIFFileGroup struct {
	JSON GLEIFJSONFileInfo `json:"json"`
}

type GLEIFDeltaFileGroups struct {
	IntraDay  GLEIFFileGroup `json:"IntraDay"`
	LastDay   GLEIFFileGroup `json:"LastDay"`
	LastWeek  GLEIFFileGroup `json:"LastWeek"`
	LastMonth GLEIFFileGroup `json:"LastMonth"`
}

type GLEIFJSONFileInfo struct {
	URL         string `json:"url"`
	Size        int64  `json:"size"`
	RecordCount int    `json:"record_count"`
	PublishedAt string `json:"published_at"`
	DeltaType   string `json:"delta_type"`
}

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
	CountLEIRecords() (int64, error)
	UpdateLEIRecord(record *domain.LEIRecord) error

	// Audit and history
	GetAuditHistory(lei string, limit int) ([]*domain.LEIRecordAudit, error)

	// Processing status
	GetProcessingStatus(jobType string) (*domain.FileProcessingStatus, error)
	UpdateProcessingStatus(status *domain.FileProcessingStatus) error
}

type leiService struct {
	repo    repository.LEIRepository
	dataDir string // Directory to store downloaded files
}

// NewLEIService creates a new LEI service
func NewLEIService(repo repository.LEIRepository, dataDir string) LEIService {
	return &leiService{
		repo:    repo,
		dataDir: dataDir,
	}
}

// getLatestFileURLs fetches the latest file URLs from GLEIF API
func (s *leiService) getLatestFileURLs() (*GLEIFPublishesResponse, error) {
	log.Info().Str("url", GLEIFLatestPublishesURL).Msg("Fetching latest file URLs from GLEIF")

	resp, err := http.Get(GLEIFLatestPublishesURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest publishes: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch latest publishes: HTTP %d", resp.StatusCode)
	}

	// Read the response body for debugging
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var publishesResp GLEIFPublishesResponse
	if err := json.Unmarshal(body, &publishesResp); err != nil {
		log.Error().Err(err).Str("body_preview", string(body[:500])).Msg("Failed to parse GLEIF API response")
		return nil, fmt.Errorf("failed to decode publishes response: %w", err)
	}

	// Debug: Log the parsed structure to verify unmarshaling
	fullURL := publishesResp.Data.LEI2.FullFile.JSON.URL
	deltaURL := publishesResp.Data.LEI2.DeltaFiles.LastWeek.JSON.URL

	log.Info().
		Str("full_url", fullURL).
		Int64("full_size", publishesResp.Data.LEI2.FullFile.JSON.Size).
		Int("full_records", publishesResp.Data.LEI2.FullFile.JSON.RecordCount).
		Str("delta_url", deltaURL).
		Int64("delta_size", publishesResp.Data.LEI2.DeltaFiles.LastWeek.JSON.Size).
		Int("delta_records", publishesResp.Data.LEI2.DeltaFiles.LastWeek.JSON.RecordCount).
		Msgf("Retrieved latest file information (full empty: %v, delta empty: %v)", fullURL == "", deltaURL == "")

	return &publishesResp, nil
}

// DownloadFullFile downloads the full LEI data file from GLEIF
func (s *leiService) DownloadFullFile() (*domain.SourceFile, error) {
	publishes, err := s.getLatestFileURLs()
	if err != nil {
		return nil, fmt.Errorf("failed to get latest file URLs: %w", err)
	}

	url := publishes.Data.LEI2.FullFile.JSON.URL
	publishedAt := publishes.Data.LEI2.PublishDate
	return s.downloadFile(url, "FULL", publishedAt)
}

// DownloadDeltaFile downloads the delta LEI data file from GLEIF
func (s *leiService) DownloadDeltaFile() (*domain.SourceFile, error) {
	publishes, err := s.getLatestFileURLs()
	if err != nil {
		return nil, fmt.Errorf("failed to get latest file URLs: %w", err)
	}

	url := publishes.Data.LEI2.DeltaFiles.LastWeek.JSON.URL
	publishedAt := publishes.Data.LEI2.PublishDate
	return s.downloadFile(url, "DELTA", publishedAt)
}

// downloadFile downloads a file from GLEIF and creates a SourceFile record
func (s *leiService) downloadFile(url, fileType, publishedAt string) (*domain.SourceFile, error) {
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

	// Parse publication date
	var publicationDate time.Time
	if publishedAt != "" {
		if t, err := time.Parse(time.RFC3339, publishedAt); err == nil {
			publicationDate = t
		} else {
			publicationDate = time.Now()
		}
	} else {
		publicationDate = time.Now()
	}

	// Create SourceFile record
	sourceFile := &domain.SourceFile{
		FileName:         fileName,
		FileType:         fileType,
		FileURL:          url,
		FileSize:         fileSize,
		FileHash:         fileHash,
		DownloadedAt:     time.Now(),
		PublicationDate:  publicationDate,
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
// GLEIF JSON format: {"records": [ {...}, {...}, ... ]}
func (s *leiService) processJSONFile(jsonPath string, sourceFile *domain.SourceFile, resumeFromLEI string) error {
	file, err := os.Open(jsonPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a JSON decoder
	decoder := json.NewDecoder(file)

	// Read the opening brace
	token, err := decoder.Token()
	if err != nil {
		return fmt.Errorf("failed to read opening brace: %w", err)
	}
	if delim, ok := token.(json.Delim); !ok || delim != '{' {
		return fmt.Errorf("expected '{', got %v", token)
	}

	// Read until we find the "records" key
	for decoder.More() {
		token, err := decoder.Token()
		if err != nil {
			return fmt.Errorf("failed to read token: %w", err)
		}

		if key, ok := token.(string); ok && key == "records" {
			// Found the records array, start processing
			return s.processRecordsArray(decoder, sourceFile, resumeFromLEI)
		}

		// Skip the value for non-records keys
		var skipValue interface{}
		if err := decoder.Decode(&skipValue); err != nil {
			return fmt.Errorf("failed to skip value: %w", err)
		}
	}

	return fmt.Errorf("records array not found in JSON file")
}

// processRecordsArray processes the records array from the JSON decoder
func (s *leiService) processRecordsArray(decoder *json.Decoder, sourceFile *domain.SourceFile, resumeFromLEI string) error {
	// Read the opening bracket of the records array
	token, err := decoder.Token()
	if err != nil {
		return fmt.Errorf("failed to read array opening: %w", err)
	}
	if delim, ok := token.(json.Delim); !ok || delim != '[' {
		return fmt.Errorf("expected '[', got %v", token)
	}

	var totalRecords int
	var processedRecords int
	var failedRecords int
	var shouldProcess bool = (resumeFromLEI == "")

	// Process each record in the array
	for decoder.More() {
		var jsonRecord LEIJSONRecord
		if err := decoder.Decode(&jsonRecord); err != nil {
			log.Error().Err(err).Msg("Failed to decode LEI JSON record")
			failedRecords++
			continue
		}

		totalRecords++

		// Check if we should start processing (resume logic)
		if !shouldProcess {
			lei := s.extractLEI(&jsonRecord)
			if lei == resumeFromLEI {
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

	// Final update
	sourceFile.TotalRecords = totalRecords
	sourceFile.ProcessedRecords = processedRecords
	sourceFile.FailedRecords = failedRecords

	return nil
}

// extractLEI extracts the LEI string from a JSON record (handles nested $ structure)
func (s *leiService) extractLEI(jsonRecord *LEIJSONRecord) string {
	return jsonRecord.LEI.Value
}

// LEIJSONRecord represents the JSON structure from GLEIF bulk files
// NOTE: This is the BULK FILE FORMAT. The single LEI API query format is different.
// Bulk files use nested objects with $ properties for values.
// Single LEI queries return a different structure - implement separately if needed.
type LEIJSONRecord struct {
	LEI          LEIValueField   `json:"LEI"`
	Entity       LEIEntity       `json:"Entity"`
	Registration LEIRegistration `json:"Registration"`
}

// LEIValueField represents a simple value field with $ property
type LEIValueField struct {
	Value string `json:"$"`
}

type LEIEntity struct {
	LegalName                      LEILegalName             `json:"LegalName"`
	OtherEntityNames               LEIOtherEntityNames      `json:"OtherEntityNames"`
	TransliteratedOtherEntityNames LEIOtherEntityNames      `json:"TransliteratedOtherEntityNames"`
	LegalAddress                   LEIAddress               `json:"LegalAddress"`
	HeadquartersAddress            LEIAddress               `json:"HeadquartersAddress"`
	RegistrationAuthority          LEIRegistrationAuthority `json:"RegistrationAuthority"`
	LegalJurisdiction              LEIValueField            `json:"LegalJurisdiction"`
	EntityCategory                 LEIValueField            `json:"EntityCategory"`
	LegalForm                      LEILegalForm             `json:"LegalForm"`
	EntityStatus                   LEIValueField            `json:"EntityStatus"`
}

type LEILegalName struct {
	Value    string `json:"$"`
	Language string `json:"@xml:lang"`
}

type LEIOtherEntityNames struct {
	OtherEntityName []LEIOtherName `json:"OtherEntityName"`
}

type LEIOtherName struct {
	Value    string `json:"$"`
	Type     string `json:"@type"`
	Language string `json:"@xml:lang"`
}

type LEIAddress struct {
	FirstAddressLine      LEIValueField   `json:"FirstAddressLine"`
	AdditionalAddressLine []LEIValueField `json:"AdditionalAddressLine"`
	City                  LEIValueField   `json:"City"`
	Region                LEIValueField   `json:"Region"`
	Country               LEIValueField   `json:"Country"`
	PostalCode            LEIValueField   `json:"PostalCode"`
	Language              string          `json:"@xml:lang"`
}

type LEIRegistrationAuthority struct {
	RegistrationAuthorityID       LEIValueField `json:"RegistrationAuthorityID"`
	RegistrationAuthorityEntityID LEIValueField `json:"RegistrationAuthorityEntityID"`
}

type LEILegalForm struct {
	EntityLegalFormCode LEIValueField `json:"EntityLegalFormCode"`
	OtherLegalForm      LEIValueField `json:"OtherLegalForm"`
}

type LEIRegistration struct {
	InitialRegistrationDate LEIValueField          `json:"InitialRegistrationDate"`
	LastUpdateDate          LEIValueField          `json:"LastUpdateDate"`
	RegistrationStatus      LEIValueField          `json:"RegistrationStatus"`
	NextRenewalDate         LEIValueField          `json:"NextRenewalDate"`
	ManagingLOU             LEIValueField          `json:"ManagingLOU"`
	ValidationSources       LEIValueField          `json:"ValidationSources"`
	ValidationAuthority     LEIValidationAuthority `json:"ValidationAuthority"`
}

type LEIValidationAuthority struct {
	ValidationAuthorityID       LEIValueField `json:"ValidationAuthorityID"`
	ValidationAuthorityEntityID LEIValueField `json:"ValidationAuthorityEntityID"`
}

// jsonToDomainRecord converts a JSON record to a domain.LEIRecord
func (s *leiService) jsonToDomainRecord(jsonRecord *LEIJSONRecord, sourceFileID uuid.UUID) *domain.LEIRecord {
	record := &domain.LEIRecord{
		LEI:                    jsonRecord.LEI.Value,
		LegalName:              jsonRecord.Entity.LegalName.Value,
		LegalAddressLine1:      jsonRecord.Entity.LegalAddress.FirstAddressLine.Value,
		LegalAddressCity:       jsonRecord.Entity.LegalAddress.City.Value,
		LegalAddressRegion:     jsonRecord.Entity.LegalAddress.Region.Value,
		LegalAddressCountry:    jsonRecord.Entity.LegalAddress.Country.Value,
		LegalAddressPostalCode: jsonRecord.Entity.LegalAddress.PostalCode.Value,
		RegistrationAuthority:  jsonRecord.Entity.RegistrationAuthority.RegistrationAuthorityID.Value,
		RegistrationNumber:     jsonRecord.Entity.RegistrationAuthority.RegistrationAuthorityEntityID.Value,
		EntityCategory:         jsonRecord.Entity.EntityCategory.Value,
		EntityLegalForm:        jsonRecord.Entity.LegalForm.EntityLegalFormCode.Value,
		EntityStatus:           jsonRecord.Entity.EntityStatus.Value,
		ManagingLOU:            jsonRecord.Registration.ManagingLOU.Value,
		SourceFileID:           &sourceFileID,
		// Initialize JSONB fields with valid JSON
		OtherNames:        "[]",
		ValidationSources: "{}",
		ChangedFields:     "{}",
	}

	// Handle additional address lines
	if len(jsonRecord.Entity.LegalAddress.AdditionalAddressLine) > 0 {
		record.LegalAddressLine2 = jsonRecord.Entity.LegalAddress.AdditionalAddressLine[0].Value
	}
	if len(jsonRecord.Entity.LegalAddress.AdditionalAddressLine) > 1 {
		record.LegalAddressLine3 = jsonRecord.Entity.LegalAddress.AdditionalAddressLine[1].Value
	}
	if len(jsonRecord.Entity.LegalAddress.AdditionalAddressLine) > 2 {
		record.LegalAddressLine4 = jsonRecord.Entity.LegalAddress.AdditionalAddressLine[2].Value
	}

	// Handle headquarters address
	if jsonRecord.Entity.HeadquartersAddress.FirstAddressLine.Value != "" {
		record.HQAddressLine1 = jsonRecord.Entity.HeadquartersAddress.FirstAddressLine.Value
		record.HQAddressCity = jsonRecord.Entity.HeadquartersAddress.City.Value
		record.HQAddressRegion = jsonRecord.Entity.HeadquartersAddress.Region.Value
		record.HQAddressCountry = jsonRecord.Entity.HeadquartersAddress.Country.Value
		record.HQAddressPostalCode = jsonRecord.Entity.HeadquartersAddress.PostalCode.Value

		if len(jsonRecord.Entity.HeadquartersAddress.AdditionalAddressLine) > 0 {
			record.HQAddressLine2 = jsonRecord.Entity.HeadquartersAddress.AdditionalAddressLine[0].Value
		}
		if len(jsonRecord.Entity.HeadquartersAddress.AdditionalAddressLine) > 1 {
			record.HQAddressLine3 = jsonRecord.Entity.HeadquartersAddress.AdditionalAddressLine[1].Value
		}
		if len(jsonRecord.Entity.HeadquartersAddress.AdditionalAddressLine) > 2 {
			record.HQAddressLine4 = jsonRecord.Entity.HeadquartersAddress.AdditionalAddressLine[2].Value
		}
	}

	// Parse dates (ISO 8601 format)
	if jsonRecord.Registration.InitialRegistrationDate.Value != "" {
		if t, err := time.Parse("2006-01-02T15:04:05Z", jsonRecord.Registration.InitialRegistrationDate.Value); err == nil {
			record.InitialRegistrationDate = t
		} else if t, err := time.Parse("2006-01-02", jsonRecord.Registration.InitialRegistrationDate.Value); err == nil {
			record.InitialRegistrationDate = t
		}
	}
	if jsonRecord.Registration.LastUpdateDate.Value != "" {
		if t, err := time.Parse("2006-01-02T15:04:05Z", jsonRecord.Registration.LastUpdateDate.Value); err == nil {
			record.LastUpdateDate = t
		} else if t, err := time.Parse("2006-01-02", jsonRecord.Registration.LastUpdateDate.Value); err == nil {
			record.LastUpdateDate = t
		}
	}
	if jsonRecord.Registration.NextRenewalDate.Value != "" {
		if t, err := time.Parse("2006-01-02T15:04:05Z", jsonRecord.Registration.NextRenewalDate.Value); err == nil {
			record.NextRenewalDate = t
		} else if t, err := time.Parse("2006-01-02T17:00:00Z", jsonRecord.Registration.NextRenewalDate.Value); err == nil {
			// Handle the specific format with 17:00:00 timezone
			record.NextRenewalDate = t
		} else if t, err := time.Parse("2006-01-02", jsonRecord.Registration.NextRenewalDate.Value); err == nil {
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

// CountLEIRecords returns the total count of LEI records
func (s *leiService) CountLEIRecords() (int64, error) {
	return s.repo.CountLEIRecords()
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
