package service

import (
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/touros-platform/api/internal/domain"
	"github.com/touros-platform/api/internal/repository"
)

type PermitService interface {
	Create(req *CreatePermitRequest) (*domain.Permit, error)
	GetByID(id uuid.UUID) (*domain.Permit, error)
	GetByPermitNumber(permitNum string) (*domain.Permit, error)
	ValidatePermit(permitNum string) (*domain.Permit, error)
	Revoke(id uuid.UUID, revokedBy uuid.UUID) error
	List(limit, offset int, guideID *uuid.UUID, status *domain.PermitStatus) ([]domain.Permit, int64, error)
}

type CreatePermitRequest struct {
	GuideID     uuid.UUID
	ClientID    uuid.UUID
	ClientName  string
	ClientEmail string
	ClientPhone string
	StartDate   time.Time
	EndDate     time.Time
	Route       string
	IssuedBy    uuid.UUID
}

type permitService struct {
	permitRepo repository.PermitRepository
	guideRepo  repository.GuideRepository
}

func NewPermitService(permitRepo repository.PermitRepository, guideRepo repository.GuideRepository) PermitService {
	return &permitService{
		permitRepo: permitRepo,
		guideRepo:  guideRepo,
	}
}

func (s *permitService) Create(req *CreatePermitRequest) (*domain.Permit, error) {
	guide, err := s.guideRepo.GetByID(req.GuideID)
	if err != nil {
		return nil, errors.New("guide not found")
	}

	if guide.Status != domain.GuideStatusVerified {
		return nil, errors.New("guide must be verified to issue permits")
	}

	permitNumber := s.generatePermitNumber()
	qrCode := s.generateQRCode(permitNumber)

	permit := &domain.Permit{
		PermitNumber: permitNumber,
		GuideID:      req.GuideID,
		ClientID:     req.ClientID,
		ClientName:   req.ClientName,
		ClientEmail:  req.ClientEmail,
		ClientPhone:  req.ClientPhone,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		Route:        req.Route,
		Status:       domain.PermitStatusActive,
		QRCode:       qrCode,
		IssuedBy:     req.IssuedBy,
		IssuedAt:     time.Now(),
	}

	if err := s.permitRepo.Create(permit); err != nil {
		return nil, fmt.Errorf("failed to create permit: %w", err)
	}

	return permit, nil
}

func (s *permitService) GetByID(id uuid.UUID) (*domain.Permit, error) {
	return s.permitRepo.GetByID(id)
}

func (s *permitService) GetByPermitNumber(permitNum string) (*domain.Permit, error) {
	return s.permitRepo.GetByPermitNumber(permitNum)
}

func (s *permitService) ValidatePermit(permitNum string) (*domain.Permit, error) {
	permit, err := s.permitRepo.GetByPermitNumber(permitNum)
	if err != nil {
		return nil, errors.New("permit not found")
	}

	if permit.Status != domain.PermitStatusActive {
		return nil, fmt.Errorf("permit status is %s", permit.Status)
	}

	now := time.Now()
	if now.Before(permit.StartDate) {
		return nil, errors.New("permit has not yet started")
	}

	if now.After(permit.EndDate) {
		permit.Status = domain.PermitStatusExpired
		s.permitRepo.Update(permit)
		return nil, errors.New("permit has expired")
	}

	return permit, nil
}

func (s *permitService) Revoke(id uuid.UUID, revokedBy uuid.UUID) error {
	permit, err := s.permitRepo.GetByID(id)
	if err != nil {
		return err
	}

	if permit.Status != domain.PermitStatusActive {
		return errors.New("permit is not active")
	}

	now := time.Now()
	permit.Status = domain.PermitStatusRevoked
	permit.RevokedAt = &now
	permit.RevokedBy = &revokedBy

	return s.permitRepo.Update(permit)
}

func (s *permitService) List(limit, offset int, guideID *uuid.UUID, status *domain.PermitStatus) ([]domain.Permit, int64, error) {
	return s.permitRepo.List(limit, offset, guideID, status)
}

func (s *permitService) generatePermitNumber() string {
	return fmt.Sprintf("TP-%s", uuid.New().String()[:8])
}

func (s *permitService) generateQRCode(permitNumber string) string {
	data := fmt.Sprintf("touros:permit:%s", permitNumber)
	return base64.StdEncoding.EncodeToString([]byte(data))
}

