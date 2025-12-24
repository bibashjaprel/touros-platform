package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/touros-platform/api/internal/domain"
	"gorm.io/gorm"
)

type SafetyCheckInRepository interface {
	Create(checkIn *domain.SafetyCheckIn) error
	GetByID(id uuid.UUID) (*domain.SafetyCheckIn, error)
	ListByGuideID(guideID uuid.UUID, limit, offset int) ([]domain.SafetyCheckIn, int64, error)
	ListRecentByGuideID(guideID uuid.UUID, since time.Time) ([]domain.SafetyCheckIn, error)
}

type safetyCheckInRepository struct {
	db *gorm.DB
}

func NewSafetyCheckInRepository(db *gorm.DB) SafetyCheckInRepository {
	return &safetyCheckInRepository{db: db}
}

func (r *safetyCheckInRepository) Create(checkIn *domain.SafetyCheckIn) error {
	return r.db.Create(checkIn).Error
}

func (r *safetyCheckInRepository) GetByID(id uuid.UUID) (*domain.SafetyCheckIn, error) {
	var checkIn domain.SafetyCheckIn
	err := r.db.Preload("Guide.User").Preload("Permit").Where("id = ?", id).First(&checkIn).Error
	if err != nil {
		return nil, err
	}
	return &checkIn, nil
}

func (r *safetyCheckInRepository) ListByGuideID(guideID uuid.UUID, limit, offset int) ([]domain.SafetyCheckIn, int64, error) {
	var checkIns []domain.SafetyCheckIn
	var total int64

	if err := r.db.Model(&domain.SafetyCheckIn{}).Where("guide_id = ?", guideID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Preload("Guide.User").Preload("Permit").
		Where("guide_id = ?", guideID).
		Limit(limit).Offset(offset).
		Order("check_in_time DESC").
		Find(&checkIns).Error
	return checkIns, total, err
}

func (r *safetyCheckInRepository) ListRecentByGuideID(guideID uuid.UUID, since time.Time) ([]domain.SafetyCheckIn, error) {
	var checkIns []domain.SafetyCheckIn
	err := r.db.Preload("Guide.User").Preload("Permit").
		Where("guide_id = ? AND check_in_time >= ?", guideID, since).
		Order("check_in_time DESC").
		Find(&checkIns).Error
	return checkIns, err
}

type IncidentRepository interface {
	Create(incident *domain.Incident) error
	GetByID(id uuid.UUID) (*domain.Incident, error)
	Update(incident *domain.Incident) error
	List(limit, offset int, status *domain.IncidentStatus, guideID *uuid.UUID) ([]domain.Incident, int64, error)
	GetActiveSOSByGuideID(guideID uuid.UUID) ([]domain.Incident, error)
}

type incidentRepository struct {
	db *gorm.DB
}

func NewIncidentRepository(db *gorm.DB) IncidentRepository {
	return &incidentRepository{db: db}
}

func (r *incidentRepository) Create(incident *domain.Incident) error {
	return r.db.Create(incident).Error
}

func (r *incidentRepository) GetByID(id uuid.UUID) (*domain.Incident, error) {
	var incident domain.Incident
	err := r.db.Preload("Guide.User").Preload("Guide.Agency").Preload("Permit").Where("id = ?", id).First(&incident).Error
	if err != nil {
		return nil, err
	}
	return &incident, nil
}

func (r *incidentRepository) Update(incident *domain.Incident) error {
	return r.db.Save(incident).Error
}

func (r *incidentRepository) List(limit, offset int, status *domain.IncidentStatus, guideID *uuid.UUID) ([]domain.Incident, int64, error) {
	var incidents []domain.Incident
	var total int64

	query := r.db.Model(&domain.Incident{}).Preload("Guide.User").Preload("Guide.Agency").Preload("Permit")
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if guideID != nil {
		query = query.Where("guide_id = ?", *guideID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Limit(limit).Offset(offset).Order("reported_at DESC").Find(&incidents).Error
	return incidents, total, err
}

func (r *incidentRepository) GetActiveSOSByGuideID(guideID uuid.UUID) ([]domain.Incident, error) {
	var incidents []domain.Incident
	err := r.db.Preload("Guide.User").Preload("Guide.Agency").Preload("Permit").
		Where("guide_id = ? AND incident_type = ? AND status IN ?", 
			guideID, domain.IncidentTypeSOS, []domain.IncidentStatus{domain.IncidentStatusOpen, domain.IncidentStatusInProgress}).
		Find(&incidents).Error
	return incidents, err
}

