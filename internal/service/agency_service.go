package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/touros-platform/api/internal/domain"
	"github.com/touros-platform/api/internal/repository"
)

type AgencyService interface {
	Create(agency *domain.Agency) error
	GetByID(id uuid.UUID) (*domain.Agency, error)
	Update(id uuid.UUID, updates *UpdateAgencyRequest) (*domain.Agency, error)
	Delete(id uuid.UUID) error
	List(limit, offset int, status *domain.AgencyStatus) ([]domain.Agency, int64, error)
	Verify(id uuid.UUID, verifiedBy uuid.UUID) error
	Suspend(id uuid.UUID, verifiedBy uuid.UUID) error
}

type UpdateAgencyRequest struct {
	Name              *string
	ContactEmail      *string
	ContactPhone      *string
	Address           *string
	LicenseExpiry     *time.Time
}

type agencyService struct {
	agencyRepo repository.AgencyRepository
}

func NewAgencyService(agencyRepo repository.AgencyRepository) AgencyService {
	return &agencyService{
		agencyRepo: agencyRepo,
	}
}

func (s *agencyService) Create(agency *domain.Agency) error {
	existing, _ := s.agencyRepo.GetByRegistrationNumber(agency.RegistrationNumber)
	if existing != nil {
		return errors.New("agency with this registration number already exists")
	}

	return s.agencyRepo.Create(agency)
}

func (s *agencyService) GetByID(id uuid.UUID) (*domain.Agency, error) {
	return s.agencyRepo.GetByID(id)
}

func (s *agencyService) Update(id uuid.UUID, updates *UpdateAgencyRequest) (*domain.Agency, error) {
	agency, err := s.agencyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if updates.Name != nil {
		agency.Name = *updates.Name
	}
	if updates.ContactEmail != nil {
		agency.ContactEmail = *updates.ContactEmail
	}
	if updates.ContactPhone != nil {
		agency.ContactPhone = *updates.ContactPhone
	}
	if updates.Address != nil {
		agency.Address = *updates.Address
	}
	if updates.LicenseExpiry != nil {
		agency.LicenseExpiry = updates.LicenseExpiry
	}

	if err := s.agencyRepo.Update(agency); err != nil {
		return nil, err
	}

	return agency, nil
}

func (s *agencyService) Delete(id uuid.UUID) error {
	return s.agencyRepo.Delete(id)
}

func (s *agencyService) List(limit, offset int, status *domain.AgencyStatus) ([]domain.Agency, int64, error) {
	return s.agencyRepo.List(limit, offset, status)
}

func (s *agencyService) Verify(id uuid.UUID, verifiedBy uuid.UUID) error {
	agency, err := s.agencyRepo.GetByID(id)
	if err != nil {
		return err
	}

	now := time.Now()
	agency.Status = domain.AgencyStatusVerified
	agency.VerifiedAt = &now
	agency.VerifiedBy = &verifiedBy

	return s.agencyRepo.Update(agency)
}

func (s *agencyService) Suspend(id uuid.UUID, verifiedBy uuid.UUID) error {
	agency, err := s.agencyRepo.GetByID(id)
	if err != nil {
		return err
	}

	agency.Status = domain.AgencyStatusSuspended
	agency.VerifiedBy = &verifiedBy

	return s.agencyRepo.Update(agency)
}

