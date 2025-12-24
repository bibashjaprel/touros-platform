package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/touros-platform/api/internal/domain"
	"github.com/touros-platform/api/internal/repository"
)

type GuideService interface {
	Create(guide *domain.Guide) error
	GetByID(id uuid.UUID) (*domain.Guide, error)
	GetByUserID(userID uuid.UUID) (*domain.Guide, error)
	Update(id uuid.UUID, updates *UpdateGuideRequest) (*domain.Guide, error)
	Delete(id uuid.UUID) error
	List(limit, offset int, status *domain.GuideStatus, agencyID *uuid.UUID) ([]domain.Guide, int64, error)
	Verify(id uuid.UUID, verifiedBy uuid.UUID) error
	Suspend(id uuid.UUID, verifiedBy uuid.UUID) error
}

type UpdateGuideRequest struct {
	PhoneNumber      *string
	EmergencyContact *string
	LicenseExpiry    *time.Time
	AgencyID         *uuid.UUID
}

type guideService struct {
	guideRepo repository.GuideRepository
	userRepo  repository.UserRepository
}

func NewGuideService(guideRepo repository.GuideRepository, userRepo repository.UserRepository) GuideService {
	return &guideService{
		guideRepo: guideRepo,
		userRepo:  userRepo,
	}
}

func (s *guideService) Create(guide *domain.Guide) error {
	existing, _ := s.guideRepo.GetByLicenseNumber(guide.LicenseNumber)
	if existing != nil {
		return errors.New("guide with this license number already exists")
	}

	existingUser, _ := s.guideRepo.GetByUserID(guide.UserID)
	if existingUser != nil {
		return errors.New("user already has a guide profile")
	}

	return s.guideRepo.Create(guide)
}

func (s *guideService) GetByID(id uuid.UUID) (*domain.Guide, error) {
	return s.guideRepo.GetByID(id)
}

func (s *guideService) GetByUserID(userID uuid.UUID) (*domain.Guide, error) {
	return s.guideRepo.GetByUserID(userID)
}

func (s *guideService) Update(id uuid.UUID, updates *UpdateGuideRequest) (*domain.Guide, error) {
	guide, err := s.guideRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if updates.PhoneNumber != nil {
		guide.PhoneNumber = *updates.PhoneNumber
	}
	if updates.EmergencyContact != nil {
		guide.EmergencyContact = *updates.EmergencyContact
	}
	if updates.LicenseExpiry != nil {
		guide.LicenseExpiry = updates.LicenseExpiry
	}
	if updates.AgencyID != nil {
		guide.AgencyID = updates.AgencyID
	}

	if err := s.guideRepo.Update(guide); err != nil {
		return nil, err
	}

	return guide, nil
}

func (s *guideService) Delete(id uuid.UUID) error {
	return s.guideRepo.Delete(id)
}

func (s *guideService) List(limit, offset int, status *domain.GuideStatus, agencyID *uuid.UUID) ([]domain.Guide, int64, error) {
	return s.guideRepo.List(limit, offset, status, agencyID)
}

func (s *guideService) Verify(id uuid.UUID, verifiedBy uuid.UUID) error {
	guide, err := s.guideRepo.GetByID(id)
	if err != nil {
		return err
	}

	now := time.Now()
	guide.Status = domain.GuideStatusVerified
	guide.VerifiedAt = &now
	guide.VerifiedBy = &verifiedBy

	return s.guideRepo.Update(guide)
}

func (s *guideService) Suspend(id uuid.UUID, verifiedBy uuid.UUID) error {
	guide, err := s.guideRepo.GetByID(id)
	if err != nil {
		return err
	}

	guide.Status = domain.GuideStatusSuspended
	guide.VerifiedBy = &verifiedBy

	return s.guideRepo.Update(guide)
}

func (s *guideService) CheckLicenseExpiry(guideID uuid.UUID) (bool, error) {
	guide, err := s.guideRepo.GetByID(guideID)
	if err != nil {
		return false, err
	}

	if guide.LicenseExpiry == nil {
		return true, nil
	}

	isExpired := guide.LicenseExpiry.Before(time.Now())
	if isExpired && guide.Status == domain.GuideStatusVerified {
		guide.Status = domain.GuideStatusSuspended
		if err := s.guideRepo.Update(guide); err != nil {
			return false, fmt.Errorf("failed to suspend guide: %w", err)
		}
	}

	return !isExpired, nil
}

