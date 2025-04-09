package models

import (
	"time"

	"gorm.io/gorm"
)

type Skill struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"type:varchar(100);not null;unique"`
	Category    string         `json:"category" gorm:"type:varchar(50);not null"` // e.g., Frontend, Backend, DevOps
	IconURL     string         `json:"icon_url" gorm:"type:varchar(255)"`
	Proficiency int            `json:"proficiency" gorm:"type:int;default:0"` // 1-100
	Featured    bool           `json:"featured" gorm:"default:false"`
	Order       int            `json:"order" gorm:"default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
