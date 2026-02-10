package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/techie2000/axiom/internal/service"
)

// LEIHandler handles LEI-related HTTP requests
type LEIHandler struct {
	leiService       service.LEIService
	schedulerService service.SchedulerService
}

// NewLEIHandler creates a new LEI handler
func NewLEIHandler(leiService service.LEIService, schedulerService service.SchedulerService) *LEIHandler {
	return &LEIHandler{
		leiService:       leiService,
		schedulerService: schedulerService,
	}
}

// GetLEIByCode retrieves an LEI record by LEI code
// @Summary Get LEI record by code
// @Description Get a specific LEI record by its LEI code
// @Tags LEI
// @Accept json
// @Produce json
// @Param lei path string true "LEI code"
// @Success 200 {object} domain.LEIRecord
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/lei/{lei} [get]
func (h *LEIHandler) GetLEIByCode(c *gin.Context) {
	lei := c.Param("lei")
	
	record, err := h.leiService.GetLEIByCode(lei)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "LEI record not found"})
		return
	}
	
	c.JSON(http.StatusOK, record)
}

// GetLEIByID retrieves an LEI record by ID
// @Summary Get LEI record by ID
// @Description Get a specific LEI record by its database ID
// @Tags LEI
// @Accept json
// @Produce json
// @Param id path string true "Record ID"
// @Success 200 {object} domain.LEIRecord
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/lei/record/{id} [get]
func (h *LEIHandler) GetLEIByID(c *gin.Context) {
	id := c.Param("id")
	
	record, err := h.leiService.GetLEIByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "LEI record not found"})
		return
	}
	
	c.JSON(http.StatusOK, record)
}

// ListLEI retrieves all LEI records with pagination
// @Summary List LEI records
// @Description Get a paginated list of LEI records
// @Tags LEI
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} domain.LEIRecord
// @Failure 500 {object} map[string]string
// @Router /api/v1/lei [get]
func (h *LEIHandler) ListLEI(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	
	if limit > 100 {
		limit = 100
	}
	
	records, err := h.leiService.GetAllLEI(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve LEI records"})
		return
	}
	
	c.JSON(http.StatusOK, records)
}

// GetAuditHistory retrieves audit history for an LEI
// @Summary Get LEI audit history
// @Description Get audit trail for a specific LEI record
// @Tags LEI
// @Accept json
// @Produce json
// @Param lei path string true "LEI code"
// @Param limit query int false "Limit" default(20)
// @Success 200 {array} domain.LEIRecordAudit
// @Failure 500 {object} map[string]string
// @Router /api/v1/lei/{lei}/audit [get]
func (h *LEIHandler) GetAuditHistory(c *gin.Context) {
	lei := c.Param("lei")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	
	audits, err := h.leiService.GetAuditHistory(lei, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve audit history"})
		return
	}
	
	c.JSON(http.StatusOK, audits)
}

// TriggerFullSync manually triggers a full sync
// @Summary Trigger full LEI sync
// @Description Manually trigger a full LEI data synchronization
// @Tags LEI
// @Accept json
// @Produce json
// @Success 202 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/lei/sync/full [post]
func (h *LEIHandler) TriggerFullSync(c *gin.Context) {
	go func() {
		if err := h.schedulerService.RunDailyFullSync(); err != nil {
			// Log error but don't fail the request
		}
	}()
	
	c.JSON(http.StatusAccepted, gin.H{"message": "Full sync triggered"})
}

// TriggerDeltaSync manually triggers a delta sync
// @Summary Trigger delta LEI sync
// @Description Manually trigger a delta LEI data synchronization
// @Tags LEI
// @Accept json
// @Produce json
// @Success 202 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/lei/sync/delta [post]
func (h *LEIHandler) TriggerDeltaSync(c *gin.Context) {
	go func() {
		if err := h.schedulerService.RunDailyDeltaSync(); err != nil {
			// Log error but don't fail the request
		}
	}()
	
	c.JSON(http.StatusAccepted, gin.H{"message": "Delta sync triggered"})
}

// GetProcessingStatus retrieves processing status for a job type
// @Summary Get processing status
// @Description Get the current processing status for LEI sync jobs
// @Tags LEI
// @Accept json
// @Produce json
// @Param jobType path string true "Job type (DAILY_FULL or DAILY_DELTA)"
// @Success 200 {object} domain.FileProcessingStatus
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/lei/status/{jobType} [get]
func (h *LEIHandler) GetProcessingStatus(c *gin.Context) {
	jobType := c.Param("jobType")
	
	status, err := h.leiService.GetProcessingStatus(jobType)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Processing status not found"})
		return
	}
	
	c.JSON(http.StatusOK, status)
}

// ResumeProcessing resumes processing of a source file
// @Summary Resume file processing
// @Description Resume processing of a source file from where it left off
// @Tags LEI
// @Accept json
// @Produce json
// @Param id path string true "Source file ID"
// @Success 202 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/lei/source-file/{id}/resume [post]
func (h *LEIHandler) ResumeProcessing(c *gin.Context) {
	idStr := c.Param("id")
	
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid source file ID"})
		return
	}
	
	go func() {
		if err := h.leiService.ProcessSourceFile(id); err != nil {
			// Log error but don't fail the request
		}
	}()
	
	c.JSON(http.StatusAccepted, gin.H{"message": "Processing resumed"})
}
