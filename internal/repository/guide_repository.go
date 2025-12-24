package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/touros-platform/api/internal/domain"
	"gorm.io/gorm"
)

type GuideRepository interface {
	Create(guide *domain.Guide) error
	GetByID(id uuid.UUID) (*domain.Guide, error)
	GetByUserID(userID uuid.UUID) (*domain.Guide, error)
	GetByLicenseNumber(licenseNum string) (*domain.Guide, error)
	Update(guide *domain.Guide) error
	Delete(id uuid.UUID) error
	List(limit, offset int, status *domain.GuideStatus, agencyID *uuid.UUID) ([]domain.Guide, int64, error)
	UpdateLastCheckIn(guideID uuid.UUID) error
}

type guideRepository struct {
	db *gorm.DB
}

func NewGuideRepository(db *gorm.DB) GuideRepository {
	return &guideRepository{db: db}
}

func (r *guideRepository) Create(guide *domain.Guide) error {
	return r.db.Create(guide).Error
}

func (r *guideRepository) GetByID(id uuid.UUID) (*domain.Guide, error) {
	var guide domain.Guide
	err := r.db.Preload("User").Preload("Agency").Where("id = ?", id).First(&guide).Error
	if err != nil {
		return nil, err
	}
	return &guide, nil
}

func (r *guideRepository) GetByUserID(userID uuid.UUID) (*domain.Guide, error) {
	var guide domain.Guide
	err := r.db.Preload("User").Preload("Agency").Where("user_id = ?", userID).First(&guide).Error
	if err != nil {
		return nil, err
	}
	return &guide, nil
}

func (r *guideRepository) GetByLicenseNumber(licenseNum string) (*domain.Guide, error) {
	var guide domain.Guide
	err := r.db.Preload("User").Preload("Agency").Where("license_number = ?", licenseNum).First(&guide).Error
	if err != nil {
		return nil, err
	}
	return &guide, nil
}

func (r *guideRepository) Update(guide *domain.Guide) error {
	return r.db.Save(guide).Error
}

func (r *guideRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.Guide{}, id).Error
}

func (r *guideRepository) List(limit, offset int, status *domain.GuideStatus, agencyID *uuid.UUID) ([]domain.Guide, int64, error) {
	var guides []domain.Guide
	var total int64

	query := r.db.Model(&domain.Guide{}).Preload("User").Preload("Agency")
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if agencyID != nil {
		query = query.Where("agency_id = ?", *agencyID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Limit(limit).Offset(offset).Find(&guides).Error
	return guides, total, err
}

func (r *guideRepository) UpdateLastCheckIn(guideID uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&domain.Guide{}).Where("id = ?", guideID).Update("last_check_in", now).Error
}

