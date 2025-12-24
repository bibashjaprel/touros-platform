package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/touros-platform/api/internal/domain"
	"github.com/touros-platform/api/internal/service"
)

type AgencyHandler struct {
	agencyService service.AgencyService
}

func NewAgencyHandler(agencyService service.AgencyService) *AgencyHandler {
	return &AgencyHandler{
		agencyService: agencyService,
	}
}

type CreateAgencyRequest struct {
	Name               string `json:"name" binding:"required"`
	RegistrationNumber string `json:"registration_number" binding:"required"`
	LicenseNumber      string `json:"license_number" binding:"required"`
	ContactEmail       string `json:"contact_email" binding:"required,email"`
	ContactPhone       string `json:"contact_phone" binding:"required"`
	Address            string `json:"address"`
}

type UpdateAgencyRequest struct {
	Name          *string `json:"name"`
	ContactEmail  *string `json:"contact_email"`
	ContactPhone  *string `json:"contact_phone"`
	Address       *string `json:"address"`
}

func (h *AgencyHandler) Create(c *gin.Context) {
	var req CreateAgencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	agency := &domain.Agency{
		Name:               req.Name,
		RegistrationNumber: req.RegistrationNumber,
		LicenseNumber:      req.LicenseNumber,
		ContactEmail:       req.ContactEmail,
		ContactPhone:       req.ContactPhone,
		Address:            req.Address,
		Status:             domain.AgencyStatusPending,
	}

	if err := h.agencyService.Create(agency); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, agency)
}

func (h *AgencyHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	agency, err := h.agencyService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "agency not found"})
		return
	}

	c.JSON(http.StatusOK, agency)
}

func (h *AgencyHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateAgencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := &service.UpdateAgencyRequest{
		Name:         req.Name,
		ContactEmail: req.ContactEmail,
		ContactPhone: req.ContactPhone,
		Address:      req.Address,
	}

	agency, err := h.agencyService.Update(id, updates)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, agency)
}

func (h *AgencyHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	var status *domain.AgencyStatus
	if statusStr := c.Query("status"); statusStr != "" {
		s := domain.AgencyStatus(statusStr)
		status = &s
	}

	agencies, total, err := h.agencyService.List(limit, offset, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   agencies,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *AgencyHandler) Verify(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userID, _ := c.Get("user_id")
	verifiedBy := userID.(uuid.UUID)

	if err := h.agencyService.Verify(id, verifiedBy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "agency verified"})
}

func (h *AgencyHandler) Suspend(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userID, _ := c.Get("user_id")
	suspendedBy := userID.(uuid.UUID)

	if err := h.agencyService.Suspend(id, suspendedBy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "agency suspended"})
}

