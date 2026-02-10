package service

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
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
	// Level 1 data (who is who) - full file published daily
	GLEIFLevel1FullXMLURL = "https://goldencopy.gleif.org/api/v2/golden-copies/publishes/lei2/latest/download"
	GLEIFLevel1DeltaXMLURL = "https://goldencopy.gleif.org/api/v2/golden-copies/publishes/lei2-delta/latest/download"
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
	return s.downloadFile(GLEIFLevel1FullXMLURL, "FULL")
}

// DownloadDeltaFile downloads the delta LEI data file from GLEIF
func (s *leiService) DownloadDeltaFile() (*domain.SourceFile, error) {
	return s.downloadFile(GLEIFLevel1DeltaXMLURL, "DELTA")
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
	fileName := fmt.Sprintf("lei-%s-%s.xml.zip", fileType, timestamp)
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
	xmlPath, err := s.extractZipFile(filePath)
	if err != nil {
		sourceFile.ProcessingStatus = "FAILED"
		sourceFile.ProcessingError = err.Error()
		s.repo.UpdateSourceFile(sourceFile)
		return fmt.Errorf("failed to extract file: %w", err)
	}
	defer os.Remove(xmlPath) // Clean up extracted XML
	
	// Parse and process XML
	if err := s.processXMLFile(xmlPath, sourceFile, resumeFromLEI); err != nil {
		sourceFile.ProcessingStatus = "FAILED"
		sourceFile.ProcessingError = err.Error()
		s.repo.UpdateSourceFile(sourceFile)
		return fmt.Errorf("failed to process XML file: %w", err)
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

// extractZipFile extracts the XML file from a ZIP archive
func (s *leiService) extractZipFile(zipPath string) (string, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", err
	}
	defer r.Close()
	
	// Find the XML file in the ZIP
	for _, f := range r.File {
		if filepath.Ext(f.Name) == ".xml" {
			rc, err := f.Open()
			if err != nil {
				return "", err
			}
			defer rc.Close()
			
			// Create output file
			xmlPath := zipPath + ".extracted.xml"
			outFile, err := os.Create(xmlPath)
			if err != nil {
				return "", err
			}
			defer outFile.Close()
			
			// Copy content
			_, err = io.Copy(outFile, rc)
			if err != nil {
				return "", err
			}
			
			return xmlPath, nil
		}
	}
	
	return "", fmt.Errorf("no XML file found in ZIP archive")
}

// processXMLFile parses and processes the LEI XML file
func (s *leiService) processXMLFile(xmlPath string, sourceFile *domain.SourceFile, resumeFromLEI string) error {
	file, err := os.Open(xmlPath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	decoder := xml.NewDecoder(file)
	
	var totalRecords int
	var processedRecords int
	var failedRecords int
	var shouldProcess bool = (resumeFromLEI == "")
	
	// Process XML elements
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("XML parsing error: %w", err)
		}
		
		switch se := token.(type) {
		case xml.StartElement:
			if se.Name.Local == "LEIRecord" {
				totalRecords++
				
				var xmlRecord LEIXMLRecord
				if err := decoder.DecodeElement(&xmlRecord, &se); err != nil {
					log.Error().Err(err).Msg("Failed to decode LEI record")
					failedRecords++
					continue
				}
				
				// Check if we should start processing (resume logic)
				if !shouldProcess {
					if xmlRecord.LEI == resumeFromLEI {
						shouldProcess = true
					} else {
						continue
					}
				}
				
				// Convert XML record to domain model
				record := s.xmlToDomainRecord(&xmlRecord, sourceFile.ID)
				
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
		}
	}
	
	// Final update
	sourceFile.TotalRecords = totalRecords
	sourceFile.ProcessedRecords = processedRecords
	sourceFile.FailedRecords = failedRecords
	
	return nil
}

// LEIXMLRecord represents the XML structure from GLEIF
// This is a simplified version - actual GLEIF XML is more complex
type LEIXMLRecord struct {
	LEI                     string    `xml:"LEI"`
	LegalName               string    `xml:"Entity>LegalName"`
	LegalAddressLine1       string    `xml:"Entity>LegalAddress>Line1"`
	LegalAddressLine2       string    `xml:"Entity>LegalAddress>Line2"`
	LegalAddressCity        string    `xml:"Entity>LegalAddress>City"`
	LegalAddressRegion      string    `xml:"Entity>LegalAddress>Region"`
	LegalAddressCountry     string    `xml:"Entity>LegalAddress>Country"`
	LegalAddressPostalCode  string    `xml:"Entity>LegalAddress>PostalCode"`
	RegistrationAuthority   string    `xml:"Entity>RegistrationAuthority"`
	RegistrationNumber      string    `xml:"Entity>RegistrationNumber"`
	EntityCategory          string    `xml:"Entity>EntityCategory"`
	EntityStatus            string    `xml:"Registration>RegistrationStatus"`
	InitialRegistrationDate string    `xml:"Registration>InitialRegistrationDate"`
	LastUpdateDate          string    `xml:"Registration>LastUpdateDate"`
	NextRenewalDate         string    `xml:"Registration>NextRenewalDate"`
}

// xmlToDomainRecord converts an XML record to a domain.LEIRecord
func (s *leiService) xmlToDomainRecord(xmlRecord *LEIXMLRecord, sourceFileID uuid.UUID) *domain.LEIRecord {
	record := &domain.LEIRecord{
		LEI:                   xmlRecord.LEI,
		LegalName:             xmlRecord.LegalName,
		LegalAddressLine1:     xmlRecord.LegalAddressLine1,
		LegalAddressLine2:     xmlRecord.LegalAddressLine2,
		LegalAddressCity:      xmlRecord.LegalAddressCity,
		LegalAddressRegion:    xmlRecord.LegalAddressRegion,
		LegalAddressCountry:   xmlRecord.LegalAddressCountry,
		LegalAddressPostalCode: xmlRecord.LegalAddressPostalCode,
		RegistrationAuthority: xmlRecord.RegistrationAuthority,
		RegistrationNumber:    xmlRecord.RegistrationNumber,
		EntityCategory:        xmlRecord.EntityCategory,
		EntityStatus:          xmlRecord.EntityStatus,
		SourceFileID:          &sourceFileID,
	}
	
	// Parse dates
	if xmlRecord.InitialRegistrationDate != "" {
		if t, err := time.Parse("2006-01-02", xmlRecord.InitialRegistrationDate); err == nil {
			record.InitialRegistrationDate = t
		}
	}
	if xmlRecord.LastUpdateDate != "" {
		if t, err := time.Parse("2006-01-02", xmlRecord.LastUpdateDate); err == nil {
			record.LastUpdateDate = t
		}
	}
	if xmlRecord.NextRenewalDate != "" {
		if t, err := time.Parse("2006-01-02", xmlRecord.NextRenewalDate); err == nil {
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
