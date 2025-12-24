package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleAgency Role = "agency"
	RoleGuide  Role = "guide"
)

type User struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email        string     `gorm:"uniqueIndex;not null"`
	PasswordHash string     `gorm:"column:password_hash;not null"`
	Role         Role       `gorm:"type:varchar(20);not null;index"`
	FullName     string     `gorm:"column:full_name;not null"`
	IsActive     bool       `gorm:"column:is_active;default:true;index"`
	AgencyID     *uuid.UUID `gorm:"type:uuid;index"`
	Agency       *Agency    `gorm:"foreignKey:AgencyID"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (User) TableName() string {
	return "users"
}
