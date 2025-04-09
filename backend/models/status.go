package models

import (
	"time"

	"gorm.io/gorm"
)

// StatusType represents different types of status checks
type StatusType string

const (
	StatusTypeServer   StatusType = "server"
	StatusTypeDatabase StatusType = "database"
	StatusTypeAPI      StatusType = "api"
	StatusTypeCache    StatusType = "cache"
)

// Status represents system status check results
type Status struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Type      StatusType     `json:"type" gorm:"type:varchar(20);not null"`
	Status    string         `json:"status" gorm:"type:varchar(20);not null"` // online, offline, degraded
	Message   string         `json:"message" gorm:"type:varchar(255)"`
	Latency   int64          `json:"latency" gorm:"type:bigint"` // in milliseconds
	CheckedAt time.Time      `json:"checked_at" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// StatusCheck represents a status check request
type StatusCheck struct {
	Type    StatusType `json:"type"`
	Status  string     `json:"status"`
	Message string     `json:"message"`
	Latency int64      `json:"latency"`
}

// StatusResponse represents the API response for status checks
type StatusResponse struct {
	Overall   string                 `json:"overall"`
	Checks    map[StatusType]*Status `json:"checks"`
	UpdatedAt time.Time              `json:"updated_at"`
}

func (s *Status) TableName() string {
	return "system_status"
}

func (s *Status) BeforeCreate(tx *gorm.DB) error {
	s.CheckedAt = time.Now()
	return nil
}
