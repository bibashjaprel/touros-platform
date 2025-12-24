package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PermitStatus string

const (
	PermitStatusActive  PermitStatus = "active"
	PermitStatusExpired PermitStatus = "expired"
	PermitStatusRevoked PermitStatus = "revoked"
)

type Permit struct {
	ID           uuid.UUID    `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	PermitNumber string       `gorm:"column:permit_number;uniqueIndex;not null"`
	GuideID      uuid.UUID    `gorm:"type:uuid;not null;index"`
	Guide        Guide        `gorm:"foreignKey:GuideID"`
	ClientID     uuid.UUID    `gorm:"type:uuid;not null"`
	ClientName   string       `gorm:"column:client_name;not null"`
	ClientEmail  string       `gorm:"column:client_email"`
	ClientPhone  string       `gorm:"column:client_phone"`
	StartDate    time.Time    `gorm:"column:start_date;not null;index"`
	EndDate      time.Time    `gorm:"column:end_date;not null;index"`
	Route        string       `gorm:"type:text;not null"`
	Status       PermitStatus `gorm:"type:varchar(20);default:'active';index"`
	QRCode       string       `gorm:"column:qr_code;type:text"`
	IssuedBy     uuid.UUID    `gorm:"type:uuid;column:issued_by;not null"`
	IssuedAt     time.Time    `gorm:"column:issued_at;default:CURRENT_TIMESTAMP"`
	RevokedAt    *time.Time   `gorm:"column:revoked_at"`
	RevokedBy    *uuid.UUID   `gorm:"type:uuid;column:revoked_by"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (Permit) TableName() string {
	return "permits"
}
