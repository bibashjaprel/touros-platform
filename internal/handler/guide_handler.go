package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/touros-platform/api/internal/domain"
	"github.com/touros-platform/api/internal/service"
)

type GuideHandler struct {
	guideService service.GuideService
}

func NewGuideHandler(guideService service.GuideService) *GuideHandler {
	return &GuideHandler{
		guideService: guideService,
	}
}

type CreateGuideRequest struct {
	UserID          uuid.UUID `json:"user_id" binding:"required"`
	AgencyID        *uuid.UUID `json:"agency_id"`
	LicenseNumber   string    `json:"license_number" binding:"required"`
	PhoneNumber     string    `json:"phone_number" binding:"required"`
	EmergencyContact string   `json:"emergency_contact" binding:"required"`
}

type UpdateGuideRequest struct {
	PhoneNumber      *string `json:"phone_number"`
	EmergencyContact *string `json:"emergency_contact"`
	AgencyID         *uuid.UUID `json:"agency_id"`
}

func (h *GuideHandler) Create(c *gin.Context) {
	var req CreateGuideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	guide := &domain.Guide{
		UserID:          req.UserID,
		AgencyID:        req.AgencyID,
		LicenseNumber:   req.LicenseNumber,
		PhoneNumber:     req.PhoneNumber,
		EmergencyContact: req.EmergencyContact,
		Status:          domain.GuideStatusPending,
	}

	if err := h.guideService.Create(guide); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, guide)
}

func (h *GuideHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	guide, err := h.guideService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "guide not found"})
		return
	}

	c.JSON(http.StatusOK, guide)
}

func (h *GuideHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateGuideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := &service.UpdateGuideRequest{
		PhoneNumber:      req.PhoneNumber,
		EmergencyContact: req.EmergencyContact,
		AgencyID:         req.AgencyID,
	}

	guide, err := h.guideService.Update(id, updates)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, guide)
}

func (h *GuideHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	var status *domain.GuideStatus
	if statusStr := c.Query("status"); statusStr != "" {
		s := domain.GuideStatus(statusStr)
		status = &s
	}

	var agencyID *uuid.UUID
	if agencyIDStr := c.Query("agency_id"); agencyIDStr != "" {
		if id, err := uuid.Parse(agencyIDStr); err == nil {
			agencyID = &id
		}
	}

	guides, total, err := h.guideService.List(limit, offset, status, agencyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  guides,
		"total": total,
		"limit": limit,
		"offset": offset,
	})
}

func (h *GuideHandler) Verify(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userID, _ := c.Get("user_id")
	verifiedBy := userID.(uuid.UUID)

	if err := h.guideService.Verify(id, verifiedBy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "guide verified"})
}

func (h *GuideHandler) Suspend(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userID, _ := c.Get("user_id")
	suspendedBy := userID.(uuid.UUID)

	if err := h.guideService.Suspend(id, suspendedBy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "guide suspended"})
}

