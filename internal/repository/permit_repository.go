package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/touros-platform/api/internal/domain"
	"gorm.io/gorm"
)

type PermitRepository interface {
	Create(permit *domain.Permit) error
	GetByID(id uuid.UUID) (*domain.Permit, error)
	GetByPermitNumber(permitNum string) (*domain.Permit, error)
	Update(permit *domain.Permit) error
	Delete(id uuid.UUID) error
	List(limit, offset int, guideID *uuid.UUID, status *domain.PermitStatus) ([]domain.Permit, int64, error)
	GetActiveByGuideID(guideID uuid.UUID) ([]domain.Permit, error)
}

type permitRepository struct {
	db *gorm.DB
}

func NewPermitRepository(db *gorm.DB) PermitRepository {
	return &permitRepository{db: db}
}

func (r *permitRepository) Create(permit *domain.Permit) error {
	return r.db.Create(permit).Error
}

func (r *permitRepository) GetByID(id uuid.UUID) (*domain.Permit, error) {
	var permit domain.Permit
	err := r.db.Preload("Guide.User").Preload("Guide.Agency").Where("id = ?", id).First(&permit).Error
	if err != nil {
		return nil, err
	}
	return &permit, nil
}

func (r *permitRepository) GetByPermitNumber(permitNum string) (*domain.Permit, error) {
	var permit domain.Permit
	err := r.db.Preload("Guide.User").Preload("Guide.Agency").Where("permit_number = ?", permitNum).First(&permit).Error
	if err != nil {
		return nil, err
	}
	return &permit, nil
}

func (r *permitRepository) Update(permit *domain.Permit) error {
	return r.db.Save(permit).Error
}

func (r *permitRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.Permit{}, id).Error
}

func (r *permitRepository) List(limit, offset int, guideID *uuid.UUID, status *domain.PermitStatus) ([]domain.Permit, int64, error) {
	var permits []domain.Permit
	var total int64

	query := r.db.Model(&domain.Permit{}).Preload("Guide.User").Preload("Guide.Agency")
	if guideID != nil {
		query = query.Where("guide_id = ?", *guideID)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&permits).Error
	return permits, total, err
}

func (r *permitRepository) GetActiveByGuideID(guideID uuid.UUID) ([]domain.Permit, error) {
	var permits []domain.Permit
	now := time.Now()
	err := r.db.Preload("Guide.User").Preload("Guide.Agency").
		Where("guide_id = ? AND status = ? AND start_date <= ? AND end_date >= ?", 
			guideID, domain.PermitStatusActive, now, now).
		Find(&permits).Error
	return permits, err
}

