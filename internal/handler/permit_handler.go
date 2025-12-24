package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/touros-platform/api/internal/domain"
	"github.com/touros-platform/api/internal/middleware"
	"github.com/touros-platform/api/internal/service"
)

type PermitHandler struct {
	permitService service.PermitService
}

func NewPermitHandler(permitService service.PermitService) *PermitHandler {
	return &PermitHandler{
		permitService: permitService,
	}
}

type CreatePermitRequest struct {
	GuideID     uuid.UUID `json:"guide_id" binding:"required"`
	ClientID    uuid.UUID `json:"client_id" binding:"required"`
	ClientName  string    `json:"client_name" binding:"required"`
	ClientEmail string    `json:"client_email" binding:"email"`
	ClientPhone string    `json:"client_phone"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date" binding:"required"`
	Route       string    `json:"route" binding:"required"`
}

func (h *PermitHandler) Create(c *gin.Context) {
	var req CreatePermitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	issuedBy := userID.(uuid.UUID)

	serviceReq := &service.CreatePermitRequest{
		GuideID:     req.GuideID,
		ClientID:    req.ClientID,
		ClientName:  req.ClientName,
		ClientEmail: req.ClientEmail,
		ClientPhone: req.ClientPhone,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Route:       req.Route,
		IssuedBy:    issuedBy,
	}

	permit, err := h.permitService.Create(serviceReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	middleware.IncrementPermitsIssued()
	c.JSON(http.StatusCreated, permit)
}

func (h *PermitHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	permit, err := h.permitService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "permit not found"})
		return
	}

	c.JSON(http.StatusOK, permit)
}

func (h *PermitHandler) Validate(c *gin.Context) {
	permitNum := c.Param("number")

	permit, err := h.permitService.ValidatePermit(permitNum)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, permit)
}

func (h *PermitHandler) Revoke(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userID, _ := c.Get("user_id")
	revokedBy := userID.(uuid.UUID)

	if err := h.permitService.Revoke(id, revokedBy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "permit revoked"})
}

func (h *PermitHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	var guideID *uuid.UUID
	if guideIDStr := c.Query("guide_id"); guideIDStr != "" {
		if id, err := uuid.Parse(guideIDStr); err == nil {
			guideID = &id
		}
	}

	var status *string
	if statusStr := c.Query("status"); statusStr != "" {
		status = &statusStr
	}

	var permitStatus *domain.PermitStatus
	if status != nil {
		s := domain.PermitStatus(*status)
		permitStatus = &s
	}

	permits, total, err := h.permitService.List(limit, offset, guideID, permitStatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   permits,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

