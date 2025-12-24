package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/touros-platform/api/internal/domain"
	"github.com/touros-platform/api/internal/middleware"
	"github.com/touros-platform/api/internal/service"
)

type SafetyHandler struct {
	safetyService service.SafetyService
}

func NewSafetyHandler(safetyService service.SafetyService) *SafetyHandler {
	return &SafetyHandler{
		safetyService: safetyService,
	}
}

type CreateCheckInRequest struct {
	PermitID   *uuid.UUID `json:"permit_id"`
	Latitude   float64    `json:"latitude" binding:"required"`
	Longitude  float64    `json:"longitude" binding:"required"`
	Location   string     `json:"location"`
	Notes      string     `json:"notes"`
}

type CreateIncidentRequest struct {
	IncidentType string     `json:"incident_type" binding:"required"`
	PermitID     *uuid.UUID `json:"permit_id"`
	Latitude     float64    `json:"latitude" binding:"required"`
	Longitude    float64    `json:"longitude" binding:"required"`
	Location     string     `json:"location"`
	Description  string     `json:"description" binding:"required"`
}

type UpdateIncidentRequest struct {
	Status          *string `json:"status"`
	ResolutionNotes *string `json:"resolution_notes"`
}

func (h *SafetyHandler) CreateCheckIn(c *gin.Context) {
	var req CreateCheckInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	guideID := userID.(uuid.UUID)

	serviceReq := &service.CreateCheckInRequest{
		GuideID:   guideID,
		PermitID:  req.PermitID,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Location:  req.Location,
		Notes:     req.Notes,
	}

	checkIn, err := h.safetyService.CreateCheckIn(serviceReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	middleware.IncrementCheckIns()
	c.JSON(http.StatusCreated, checkIn)
}

func (h *SafetyHandler) GetCheckInByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	checkIn, err := h.safetyService.GetCheckInByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "check-in not found"})
		return
	}

	c.JSON(http.StatusOK, checkIn)
}

func (h *SafetyHandler) ListCheckIns(c *gin.Context) {
	guideID, err := uuid.Parse(c.Param("guide_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid guide_id"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	checkIns, total, err := h.safetyService.ListCheckIns(guideID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   checkIns,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *SafetyHandler) CreateIncident(c *gin.Context) {
	var req CreateIncidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	guideID := userID.(uuid.UUID)

	serviceReq := &service.CreateIncidentRequest{
		IncidentType: req.IncidentType,
		GuideID:      guideID,
		PermitID:     req.PermitID,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		Location:     req.Location,
		Description:  req.Description,
	}

	incident, err := h.safetyService.CreateIncident(serviceReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.IncidentType == string(domain.IncidentTypeSOS) {
		middleware.IncrementSOSIncidents()
	}

	c.JSON(http.StatusCreated, incident)
}

func (h *SafetyHandler) GetIncidentByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	incident, err := h.safetyService.GetIncidentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "incident not found"})
		return
	}

	c.JSON(http.StatusOK, incident)
}

func (h *SafetyHandler) UpdateIncident(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateIncidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	resolvedBy := userID.(uuid.UUID)

	var status *domain.IncidentStatus
	if req.Status != nil {
		s := domain.IncidentStatus(*req.Status)
		status = &s
	}

	serviceReq := &service.UpdateIncidentRequest{
		Status:          status,
		ResolutionNotes: req.ResolutionNotes,
		ResolvedBy:      &resolvedBy,
	}

	incident, err := h.safetyService.UpdateIncident(id, serviceReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, incident)
}

func (h *SafetyHandler) ListIncidents(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	var status *domain.IncidentStatus
	if statusStr := c.Query("status"); statusStr != "" {
		s := domain.IncidentStatus(statusStr)
		status = &s
	}

	var guideID *uuid.UUID
	if guideIDStr := c.Query("guide_id"); guideIDStr != "" {
		if id, err := uuid.Parse(guideIDStr); err == nil {
			guideID = &id
		}
	}

	incidents, total, err := h.safetyService.ListIncidents(limit, offset, status, guideID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   incidents,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *SafetyHandler) GetActiveSOS(c *gin.Context) {
	guideID, err := uuid.Parse(c.Param("guide_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid guide_id"})
		return
	}

	incidents, err := h.safetyService.GetActiveSOS(guideID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": incidents})
}

