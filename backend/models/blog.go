package models

import (
	"database/sql"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Blog represents a blog post in the system
type Blog struct {
	ID            uint   `json:"id" gorm:"primaryKey"`
	UserID        uint   `json:"user_id" gorm:"not null;index"`
	Title         string `json:"title" gorm:"not null;type:varchar(255)"`
	Slug          string `json:"slug" gorm:"uniqueIndex;not null;type:varchar(255)"`
	Excerpt       string `json:"excerpt" gorm:"type:varchar(500)"`
	Content       string `json:"content" gorm:"type:text;not null"`
	Tags          string `json:"tags" gorm:"type:varchar(255)"`
	Category      string `json:"category" gorm:"type:varchar(100);index"`
	FeaturedImage string `json:"featured_image" gorm:"type:varchar(255)"`

	// Status fields
	Published bool `json:"published" gorm:"default:false;index"`
	Featured  bool `json:"featured" gorm:"default:false;index"`
	Views     int  `json:"views" gorm:"default:0"`

	// Time fields using proper PostgreSQL timestamp types
	CreatedAt   time.Time      `json:"created_at" gorm:"index;not null;default:now()"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"not null;default:now()"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	PublishedAt sql.NullTime   `json:"published_at,omitempty" gorm:"index;type:timestamp with time zone"`

	// SEO fields
	MetaTitle       string `json:"meta_title" gorm:"type:varchar(255)"`
	MetaDescription string `json:"meta_description" gorm:"type:varchar(320)"`
	CanonicalURL    string `json:"canonical_url" gorm:"type:varchar(255)"`
}

// TableName specifies the table name for the Blog model
func (Blog) TableName() string {
	return "blogs"
}

// BeforeCreate hook runs before creating a new blog post
func (b *Blog) BeforeCreate(tx *gorm.DB) error {
	// If the blog is published and PublishedAt is not set, set it to current time
	if b.Published && !b.PublishedAt.Valid {
		b.PublishedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
	}

	// If no meta title is provided, use the regular title
	if b.MetaTitle == "" {
		b.MetaTitle = b.Title
	}

	// If no excerpt is provided but content exists, create one
	if b.Excerpt == "" && len(b.Content) > 0 {
		// Take first 160 chars of content (standard meta description length)
		contentLen := len(b.Content)
		if contentLen > 160 {
			contentLen = 160
		}
		b.Excerpt = b.Content[:contentLen] + "..."
	}

	return nil
}

// BeforeUpdate hook runs before updating a blog post
func (b *Blog) BeforeUpdate(tx *gorm.DB) error {
	// If the blog is being published now (wasn't before), set PublishedAt
	if b.Published {
		var oldBlog Blog
		if err := tx.First(&oldBlog, b.ID).Error; err == nil {
			if !oldBlog.Published && !b.PublishedAt.Valid {
				b.PublishedAt = sql.NullTime{
					Time:  time.Now(),
					Valid: true,
				}
			}
		}
	}

	return nil
}

// PostgreSQL-specific indexes and constraints (run this in a migration)
func (Blog) Migrations() []string {
	return []string{
		// Full text search index for PostgreSQL
		`CREATE INDEX IF NOT EXISTS blogs_content_search_idx ON blogs USING gin(to_tsvector('english', title || ' ' || content))`,

		// Index for finding latest posts quickly
		`CREATE INDEX IF NOT EXISTS blogs_published_date_idx ON blogs(published_at DESC) WHERE published = true`,

		// Add constraint to ensure published blogs have a published_at date
		`ALTER TABLE blogs ADD CONSTRAINT chk_published_date CHECK (NOT published OR published_at IS NOT NULL)`,
	}
}

// Structure for returning blog responses to clients
type BlogResponse struct {
	ID            uint      `json:"id"`
	Title         string    `json:"title"`
	Slug          string    `json:"slug"`
	Excerpt       string    `json:"excerpt"`
	FeaturedImage string    `json:"featured_image"`
	Category      string    `json:"category"`
	Tags          []string  `json:"tags,omitempty"`
	Published     bool      `json:"published"`
	Featured      bool      `json:"featured"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	PublishedAt   time.Time `json:"published_at,omitempty"`
	Content       string    `json:"content,omitempty"` // Only included in detailed views
	Views         int       `json:"views"`
	Author        struct {
		ID       uint   `json:"id"`
		Name     string `json:"name"`
		Username string `json:"username"`
		Avatar   string `json:"avatar,omitempty"`
	} `json:"author,omitempty"`
}

// ToResponse converts a Blog model to a BlogResponse
func (b *Blog) ToResponse(includeContent bool, includeAuthor bool) BlogResponse {
	response := BlogResponse{
		ID:            b.ID,
		Title:         b.Title,
		Slug:          b.Slug,
		Excerpt:       b.Excerpt,
		FeaturedImage: b.FeaturedImage,
		Category:      b.Category,
		Published:     b.Published,
		Featured:      b.Featured,
		CreatedAt:     b.CreatedAt,
		UpdatedAt:     b.UpdatedAt,
		Views:         b.Views,
	}

	// Convert tags string to array
	if b.Tags != "" {
		response.Tags = strings.Split(b.Tags, ",")
	}

	// Include published date if available
	if b.PublishedAt.Valid {
		response.PublishedAt = b.PublishedAt.Time
	}

	// Include content only when needed
	if includeContent {
		response.Content = b.Content
	}

	return response
}
