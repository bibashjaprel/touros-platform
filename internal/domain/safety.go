package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IncidentType string

const (
	IncidentTypeCheckIn IncidentType = "check_in"
	IncidentTypeSOS     IncidentType = "sos"
	IncidentTypeMedical IncidentType = "medical"
	IncidentTypeWeather IncidentType = "weather"
	IncidentTypeOther   IncidentType = "other"
)

type IncidentStatus string

const (
	IncidentStatusOpen       IncidentStatus = "open"
	IncidentStatusInProgress IncidentStatus = "in_progress"
	IncidentStatusResolved   IncidentStatus = "resolved"
	IncidentStatusClosed     IncidentStatus = "closed"
)

type SafetyCheckIn struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	GuideID     uuid.UUID  `gorm:"type:uuid;not null;index"`
	Guide       Guide      `gorm:"foreignKey:GuideID"`
	PermitID    *uuid.UUID `gorm:"type:uuid;index"`
	Permit      *Permit    `gorm:"foreignKey:PermitID"`
	Latitude    float64    `gorm:"type:decimal(10,8);not null"`
	Longitude   float64    `gorm:"type:decimal(11,8);not null"`
	Location    string     `gorm:"type:text"`
	Notes       string     `gorm:"type:text"`
	CheckInTime time.Time  `gorm:"column:check_in_time;default:CURRENT_TIMESTAMP;index"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (SafetyCheckIn) TableName() string {
	return "safety_check_ins"
}

type Incident struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	IncidentType    IncidentType   `gorm:"column:incident_type;type:varchar(20);not null;index"`
	GuideID         uuid.UUID      `gorm:"type:uuid;not null;index"`
	Guide           Guide          `gorm:"foreignKey:GuideID"`
	PermitID        *uuid.UUID     `gorm:"type:uuid;index"`
	Permit          *Permit        `gorm:"foreignKey:PermitID"`
	Status          IncidentStatus `gorm:"type:varchar(20);default:'open';index"`
	Latitude        float64        `gorm:"type:decimal(10,8);not null"`
	Longitude       float64        `gorm:"type:decimal(11,8);not null"`
	Location        string         `gorm:"type:text"`
	Description     string         `gorm:"type:text;not null"`
	ReportedAt      time.Time      `gorm:"column:reported_at;default:CURRENT_TIMESTAMP;index"`
	ResolvedAt      *time.Time     `gorm:"column:resolved_at"`
	ResolvedBy      *uuid.UUID     `gorm:"type:uuid;column:resolved_by"`
	ResolutionNotes string         `gorm:"column:resolution_notes;type:text"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

func (Incident) TableName() string {
	return "incidents"
}
