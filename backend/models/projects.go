package models

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type Project struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Title       string `json:"title" gorm:"type:varchar(255);not null"`
	Slug        string `json:"slug" gorm:"uniqueIndex;not null;type:varchar(255)"`
	Description string `json:"description" gorm:"type:text"`
	ImageURL    string `json:"image_url" gorm:"type:varchar(255)"`
	LiveURL     string `json:"live_url,omitempty" gorm:"type:varchar(255)"`
	GithubURL   string `json:"github_url,omitempty" gorm:"type:varchar(255)"`

	// Project details
	Technologies []string `json:"technologies" gorm:"-"` // Stored as comma-separated in DB
	TechStack    string   `json:"-" gorm:"type:text"`    // Actual DB column
	Featured     bool     `json:"featured" gorm:"default:false"`
	Order        int      `json:"order" gorm:"default:0"` // For custom ordering

	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// BeforeSave converts technologies slice to string
func (p *Project) BeforeSave(tx *gorm.DB) error {
	p.TechStack = strings.Join(p.Technologies, ",")
	return nil
}

// AfterFind converts string to technologies slice
func (p *Project) AfterFind(tx *gorm.DB) error {
	if p.TechStack != "" {
		p.Technologies = strings.Split(p.TechStack, ",")
	}
	return nil
}
