package service

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/techie2000/axiom/internal/domain"
)

// SchedulerService handles scheduled jobs for LEI data acquisition
type SchedulerService interface {
	Start() error
	Stop()
	RunDailyFullSync() error
	RunDailyDeltaSync() error
}

type schedulerService struct {
	leiService LEIService
	stopChan   chan struct{}
	running    bool
}

// NewSchedulerService creates a new scheduler service
func NewSchedulerService(leiService LEIService) SchedulerService {
	return &schedulerService{
		leiService: leiService,
		stopChan:   make(chan struct{}),
		running:    false,
	}
}

// Start begins the scheduler
func (s *schedulerService) Start() error {
	if s.running {
		log.Warn().Msg("Scheduler already running")
		return nil
	}
	
	s.running = true
	log.Info().Msg("Starting LEI scheduler service")
	
	// Start goroutine for daily delta sync (runs every hour to check for updates)
	go s.dailyDeltaSyncLoop()
	
	// Start goroutine for weekly full sync (runs every Sunday at 2 AM)
	go s.weeklyFullSyncLoop()
	
	return nil
}

// Stop stops the scheduler
func (s *schedulerService) Stop() {
	if !s.running {
		return
	}
	
	log.Info().Msg("Stopping LEI scheduler service")
	s.running = false
	close(s.stopChan)
}

// dailyDeltaSyncLoop runs delta sync every hour
func (s *schedulerService) dailyDeltaSyncLoop() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	
	// Run immediately on start
	if err := s.RunDailyDeltaSync(); err != nil {
		log.Error().Err(err).Msg("Failed to run initial delta sync")
	}
	
	for {
		select {
		case <-ticker.C:
			if err := s.RunDailyDeltaSync(); err != nil {
				log.Error().Err(err).Msg("Failed to run scheduled delta sync")
			}
		case <-s.stopChan:
			log.Info().Msg("Stopping delta sync loop")
			return
		}
	}
}

// weeklyFullSyncLoop runs full sync every Sunday at 2 AM
func (s *schedulerService) weeklyFullSyncLoop() {
	for {
		// Calculate next Sunday 2 AM
		now := time.Now()
		nextRun := time.Date(now.Year(), now.Month(), now.Day(), 2, 0, 0, 0, now.Location())
		
		// Add days until Sunday
		daysUntilSunday := (7 - int(now.Weekday())) % 7
		if daysUntilSunday == 0 && now.Hour() >= 2 {
			daysUntilSunday = 7 // Next week if we've already passed 2 AM today
		}
		nextRun = nextRun.AddDate(0, 0, daysUntilSunday)
		
		// If the next run is in the past, add a week
		if nextRun.Before(now) {
			nextRun = nextRun.AddDate(0, 0, 7)
		}
		
		duration := nextRun.Sub(now)
		log.Info().
			Time("next_run", nextRun).
			Dur("wait_duration", duration).
			Msg("Scheduled next full sync")
		
		select {
		case <-time.After(duration):
			if err := s.RunDailyFullSync(); err != nil {
				log.Error().Err(err).Msg("Failed to run scheduled full sync")
			}
		case <-s.stopChan:
			log.Info().Msg("Stopping full sync loop")
			return
		}
	}
}

// RunDailyDeltaSync downloads and processes delta file
func (s *schedulerService) RunDailyDeltaSync() error {
	log.Info().Msg("Starting daily delta sync")
	
	// Update processing status
	status, err := s.leiService.GetProcessingStatus("DAILY_DELTA")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get processing status")
		// Create new status if not found
		status = &domain.FileProcessingStatus{
			JobType: "DAILY_DELTA",
			Status:  "IDLE",
		}
	}
	
	// Check if already running
	if status.Status == "RUNNING" {
		log.Warn().Msg("Delta sync already running, skipping")
		return nil
	}
	
	// Update status
	status.Status = "RUNNING"
	now := time.Now()
	status.LastRunAt = &now
	if err := s.leiService.UpdateProcessingStatus(status); err != nil {
		log.Error().Err(err).Msg("Failed to update processing status")
	}
	
	// Download delta file
	sourceFile, err := s.leiService.DownloadDeltaFile()
	if err != nil {
		status.Status = "FAILED"
		status.ErrorMessage = err.Error()
		s.leiService.UpdateProcessingStatus(status)
		return err
	}
	
	// Update status with current file
	status.CurrentSourceFileID = &sourceFile.ID
	s.leiService.UpdateProcessingStatus(status)
	
	// Process file
	if err := s.leiService.ProcessSourceFile(sourceFile.ID); err != nil {
		status.Status = "FAILED"
		status.ErrorMessage = err.Error()
		s.leiService.UpdateProcessingStatus(status)
		return err
	}
	
	// Update status
	status.Status = "COMPLETED"
	status.LastSuccessAt = &now
	status.NextRunAt = calculateNextRun(1 * time.Hour)
	status.ErrorMessage = ""
	if err := s.leiService.UpdateProcessingStatus(status); err != nil {
		log.Error().Err(err).Msg("Failed to update processing status")
	}
	
	log.Info().Msg("Daily delta sync completed successfully")
	return nil
}

// RunDailyFullSync downloads and processes full file
func (s *schedulerService) RunDailyFullSync() error {
	log.Info().Msg("Starting daily full sync")
	
	// Update processing status
	status, err := s.leiService.GetProcessingStatus("DAILY_FULL")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get processing status")
		// Create new status if not found
		status = &domain.FileProcessingStatus{
			JobType: "DAILY_FULL",
			Status:  "IDLE",
		}
	}
	
	// Check if already running
	if status.Status == "RUNNING" {
		log.Warn().Msg("Full sync already running, skipping")
		return nil
	}
	
	// Update status
	status.Status = "RUNNING"
	now := time.Now()
	status.LastRunAt = &now
	if err := s.leiService.UpdateProcessingStatus(status); err != nil {
		log.Error().Err(err).Msg("Failed to update processing status")
	}
	
	// Download full file
	sourceFile, err := s.leiService.DownloadFullFile()
	if err != nil {
		status.Status = "FAILED"
		status.ErrorMessage = err.Error()
		s.leiService.UpdateProcessingStatus(status)
		return err
	}
	
	// Update status with current file
	status.CurrentSourceFileID = &sourceFile.ID
	s.leiService.UpdateProcessingStatus(status)
	
	// Process file (can resume if interrupted)
	var resumeLEI string
	if sourceFile.LastProcessedLEI != "" {
		resumeLEI = sourceFile.LastProcessedLEI
		log.Info().Str("resume_from", resumeLEI).Msg("Resuming file processing")
	}
	
	if err := s.leiService.ProcessSourceFileWithResume(sourceFile.ID, resumeLEI); err != nil {
		status.Status = "FAILED"
		status.ErrorMessage = err.Error()
		s.leiService.UpdateProcessingStatus(status)
		return err
	}
	
	// Update status
	status.Status = "COMPLETED"
	status.LastSuccessAt = &now
	status.NextRunAt = calculateNextWeeklyRun()
	status.ErrorMessage = ""
	if err := s.leiService.UpdateProcessingStatus(status); err != nil {
		log.Error().Err(err).Msg("Failed to update processing status")
	}
	
	log.Info().Msg("Daily full sync completed successfully")
	return nil
}

// calculateNextRun calculates the next run time based on interval
func calculateNextRun(interval time.Duration) *time.Time {
	next := time.Now().Add(interval)
	return &next
}

// calculateNextWeeklyRun calculates next Sunday at 2 AM
func calculateNextWeeklyRun() *time.Time {
	now := time.Now()
	nextRun := time.Date(now.Year(), now.Month(), now.Day(), 2, 0, 0, 0, now.Location())
	
	daysUntilSunday := (7 - int(now.Weekday())) % 7
	if daysUntilSunday == 0 && now.Hour() >= 2 {
		daysUntilSunday = 7
	}
	nextRun = nextRun.AddDate(0, 0, daysUntilSunday)
	
	if nextRun.Before(now) {
		nextRun = nextRun.AddDate(0, 0, 7)
	}
	
	return &nextRun
}
