package repository

import (
	"github.com/google/uuid"
	"github.com/touros-platform/api/internal/domain"
	"gorm.io/gorm"
)

type AgencyRepository interface {
	Create(agency *domain.Agency) error
	GetByID(id uuid.UUID) (*domain.Agency, error)
	GetByRegistrationNumber(regNum string) (*domain.Agency, error)
	Update(agency *domain.Agency) error
	Delete(id uuid.UUID) error
	List(limit, offset int, status *domain.AgencyStatus) ([]domain.Agency, int64, error)
}

type agencyRepository struct {
	db *gorm.DB
}

func NewAgencyRepository(db *gorm.DB) AgencyRepository {
	return &agencyRepository{db: db}
}

func (r *agencyRepository) Create(agency *domain.Agency) error {
	return r.db.Create(agency).Error
}

func (r *agencyRepository) GetByID(id uuid.UUID) (*domain.Agency, error) {
	var agency domain.Agency
	err := r.db.Where("id = ?", id).First(&agency).Error
	if err != nil {
		return nil, err
	}
	return &agency, nil
}

func (r *agencyRepository) GetByRegistrationNumber(regNum string) (*domain.Agency, error) {
	var agency domain.Agency
	err := r.db.Where("registration_number = ?", regNum).First(&agency).Error
	if err != nil {
		return nil, err
	}
	return &agency, nil
}

func (r *agencyRepository) Update(agency *domain.Agency) error {
	return r.db.Save(agency).Error
}

func (r *agencyRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.Agency{}, id).Error
}

func (r *agencyRepository) List(limit, offset int, status *domain.AgencyStatus) ([]domain.Agency, int64, error) {
	var agencies []domain.Agency
	var total int64

	query := r.db.Model(&domain.Agency{})
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Limit(limit).Offset(offset).Find(&agencies).Error
	return agencies, total, err
}
