package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AgencyStatus string

const (
	AgencyStatusPending   AgencyStatus = "pending"
	AgencyStatusVerified  AgencyStatus = "verified"
	AgencyStatusSuspended AgencyStatus = "suspended"
	AgencyStatusRejected  AgencyStatus = "rejected"
)

type Agency struct {
	ID                 uuid.UUID    `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name               string       `gorm:"not null"`
	RegistrationNumber string       `gorm:"column:registration_number;uniqueIndex;not null"`
	LicenseNumber      string       `gorm:"column:license_number;uniqueIndex;not null"`
	ContactEmail       string       `gorm:"column:contact_email;not null"`
	ContactPhone       string       `gorm:"column:contact_phone;not null"`
	Address            string       `gorm:"type:text"`
	Status             AgencyStatus `gorm:"type:varchar(20);default:'pending';index"`
	LicenseExpiry      *time.Time   `gorm:"column:license_expiry;index"`
	VerifiedAt         *time.Time   `gorm:"column:verified_at"`
	VerifiedBy         *uuid.UUID   `gorm:"type:uuid;column:verified_by"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          gorm.DeletedAt `gorm:"index"`
}

func (Agency) TableName() string {
	return "agencies"
}
