package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GuideStatus string

const (
	GuideStatusPending   GuideStatus = "pending"
	GuideStatusVerified  GuideStatus = "verified"
	GuideStatusSuspended GuideStatus = "suspended"
	GuideStatusRejected  GuideStatus = "rejected"
)

type Guide struct {
	ID               uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID           uuid.UUID   `gorm:"type:uuid;uniqueIndex;not null"`
	User             User        `gorm:"foreignKey:UserID"`
	AgencyID         *uuid.UUID  `gorm:"type:uuid;index"`
	Agency           *Agency     `gorm:"foreignKey:AgencyID"`
	LicenseNumber    string      `gorm:"column:license_number;uniqueIndex;not null"`
	PhoneNumber      string      `gorm:"column:phone_number;not null"`
	EmergencyContact string      `gorm:"column:emergency_contact;not null"`
	Status           GuideStatus `gorm:"type:varchar(20);default:'pending';index"`
	LicenseExpiry    *time.Time  `gorm:"column:license_expiry;index"`
	VerifiedAt       *time.Time  `gorm:"column:verified_at"`
	VerifiedBy       *uuid.UUID  `gorm:"type:uuid;column:verified_by"`
	LastCheckIn      *time.Time  `gorm:"column:last_check_in;index"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}

func (Guide) TableName() string {
	return "guides"
}
