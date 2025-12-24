package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/touros-platform/api/internal/domain"
	"github.com/touros-platform/api/internal/repository"
)

type SafetyService interface {
	CreateCheckIn(req *CreateCheckInRequest) (*domain.SafetyCheckIn, error)
	GetCheckInByID(id uuid.UUID) (*domain.SafetyCheckIn, error)
	ListCheckIns(guideID uuid.UUID, limit, offset int) ([]domain.SafetyCheckIn, int64, error)
	CreateIncident(req *CreateIncidentRequest) (*domain.Incident, error)
	GetIncidentByID(id uuid.UUID) (*domain.Incident, error)
	UpdateIncident(id uuid.UUID, req *UpdateIncidentRequest) (*domain.Incident, error)
	ListIncidents(limit, offset int, status *domain.IncidentStatus, guideID *uuid.UUID) ([]domain.Incident, int64, error)
	GetActiveSOS(guideID uuid.UUID) ([]domain.Incident, error)
}

type CreateCheckInRequest struct {
	GuideID   uuid.UUID
	PermitID  *uuid.UUID
	Latitude  float64
	Longitude float64
	Location  string
	Notes     string
}

type CreateIncidentRequest struct {
	IncidentType string
	GuideID      uuid.UUID
	PermitID     *uuid.UUID
	Latitude     float64
	Longitude    float64
	Location     string
	Description  string
}

type UpdateIncidentRequest struct {
	Status          *domain.IncidentStatus
	ResolutionNotes *string
	ResolvedBy      *uuid.UUID
}

type safetyService struct {
	checkInRepo  repository.SafetyCheckInRepository
	incidentRepo repository.IncidentRepository
	guideRepo    repository.GuideRepository
}

func NewSafetyService(
	checkInRepo repository.SafetyCheckInRepository,
	incidentRepo repository.IncidentRepository,
	guideRepo repository.GuideRepository,
) SafetyService {
	return &safetyService{
		checkInRepo:  checkInRepo,
		incidentRepo: incidentRepo,
		guideRepo:    guideRepo,
	}
}

func (s *safetyService) CreateCheckIn(req *CreateCheckInRequest) (*domain.SafetyCheckIn, error) {
	_, err := s.guideRepo.GetByID(req.GuideID)
	if err != nil {
		return nil, errors.New("guide not found")
	}

	checkIn := &domain.SafetyCheckIn{
		GuideID:     req.GuideID,
		PermitID:    req.PermitID,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Location:    req.Location,
		Notes:       req.Notes,
		CheckInTime: time.Now(),
	}

	if err := s.checkInRepo.Create(checkIn); err != nil {
		return nil, err
	}

	if err := s.guideRepo.UpdateLastCheckIn(req.GuideID); err != nil {
		return nil, err
	}

	return checkIn, nil
}

func (s *safetyService) GetCheckInByID(id uuid.UUID) (*domain.SafetyCheckIn, error) {
	return s.checkInRepo.GetByID(id)
}

func (s *safetyService) ListCheckIns(guideID uuid.UUID, limit, offset int) ([]domain.SafetyCheckIn, int64, error) {
	return s.checkInRepo.ListByGuideID(guideID, limit, offset)
}

func (s *safetyService) CreateIncident(req *CreateIncidentRequest) (*domain.Incident, error) {
	_, err := s.guideRepo.GetByID(req.GuideID)
	if err != nil {
		return nil, errors.New("guide not found")
	}

	incidentType := domain.IncidentType(req.IncidentType)
	if incidentType != domain.IncidentTypeCheckIn &&
		incidentType != domain.IncidentTypeSOS &&
		incidentType != domain.IncidentTypeMedical &&
		incidentType != domain.IncidentTypeWeather &&
		incidentType != domain.IncidentTypeOther {
		return nil, errors.New("invalid incident type")
	}

	incident := &domain.Incident{
		IncidentType: incidentType,
		GuideID:      req.GuideID,
		PermitID:     req.PermitID,
		Status:       domain.IncidentStatusOpen,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		Location:     req.Location,
		Description:  req.Description,
		ReportedAt:   time.Now(),
	}

	if err := s.incidentRepo.Create(incident); err != nil {
		return nil, err
	}

	return incident, nil
}

func (s *safetyService) GetIncidentByID(id uuid.UUID) (*domain.Incident, error) {
	return s.incidentRepo.GetByID(id)
}

func (s *safetyService) UpdateIncident(id uuid.UUID, req *UpdateIncidentRequest) (*domain.Incident, error) {
	incident, err := s.incidentRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Status != nil {
		incident.Status = *req.Status
		if *req.Status == domain.IncidentStatusResolved || *req.Status == domain.IncidentStatusClosed {
			now := time.Now()
			incident.ResolvedAt = &now
			if req.ResolvedBy != nil {
				incident.ResolvedBy = req.ResolvedBy
			}
		}
	}

	if req.ResolutionNotes != nil {
		incident.ResolutionNotes = *req.ResolutionNotes
	}

	if err := s.incidentRepo.Update(incident); err != nil {
		return nil, err
	}

	return incident, nil
}

func (s *safetyService) ListIncidents(limit, offset int, status *domain.IncidentStatus, guideID *uuid.UUID) ([]domain.Incident, int64, error) {
	return s.incidentRepo.List(limit, offset, status, guideID)
}

func (s *safetyService) GetActiveSOS(guideID uuid.UUID) ([]domain.Incident, error) {
	return s.incidentRepo.GetActiveSOSByGuideID(guideID)
}
