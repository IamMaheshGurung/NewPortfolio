package models

import (
	"errors"
	"html"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User represents a user in the system (admin or regular visitor with account)
type User struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	FirstName string `json:"first_name" gorm:"type:varchar(100);not null"`
	LastName  string `json:"last_name" gorm:"type:varchar(100);not null"`
	Email     string `json:"email" gorm:"uniqueIndex;type:varchar(255);not null"`
	Username  string `json:"username" gorm:"uniqueIndex;type:varchar(50);not null"`
	Password  string `json:"-" gorm:"type:varchar(255);not null"` // Never expose in JSON responses
	Bio       string `json:"bio,omitempty" gorm:"type:text"`
	Avatar    string `json:"avatar,omitempty" gorm:"type:varchar(255)"`
	Role      string `json:"role" gorm:"type:varchar(20);default:'user';not null"` // admin, user

	// Social profiles
	GithubURL   string `json:"github_url,omitempty" gorm:"type:varchar(255)"`
	LinkedinURL string `json:"linkedin_url,omitempty" gorm:"type:varchar(255)"`
	TwitterURL  string `json:"twitter_url,omitempty" gorm:"type:varchar(255)"`
	WebsiteURL  string `json:"website_url,omitempty" gorm:"type:varchar(255)"`

	// Security fields
	EmailVerified  bool       `json:"email_verified" gorm:"default:false"`
	LastLogin      time.Time  `json:"last_login,omitempty"`
	FailedAttempts int        `json:"-" gorm:"default:0"`
	LockoutEnd     *time.Time `json:"-"` // Account lockout for security
	ResetToken     string     `json:"-" gorm:"type:varchar(255)"`
	ResetTokenExp  *time.Time `json:"-"`

	// Timestamps
	CreatedAt time.Time      `json:"created_at" gorm:"index;not null"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"not null"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships (include related blogs by this user)
	Blogs []Blog `json:"-" gorm:"foreignKey:UserID"`
}

// TableName specifies the database table name
func (User) TableName() string {
	return "users"
}

// BeforeSave is a GORM hook that runs before saving the user
func (u *User) BeforeSave(tx *gorm.DB) error {
	// Only hash the password if it's been modified (not already hashed)
	if u.Password != "" && len(u.Password) < 60 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}

	// Sanitize user inputs to prevent XSS
	u.FirstName = html.EscapeString(strings.TrimSpace(u.FirstName))
	u.LastName = html.EscapeString(strings.TrimSpace(u.LastName))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	u.Bio = html.EscapeString(strings.TrimSpace(u.Bio))

	// Sanitize URLs
	u.GithubURL = sanitizeURL(u.GithubURL)
	u.LinkedinURL = sanitizeURL(u.LinkedinURL)
	u.TwitterURL = sanitizeURL(u.TwitterURL)
	u.WebsiteURL = sanitizeURL(u.WebsiteURL)

	// Ensure the role is valid
	if u.Role != "admin" && u.Role != "user" {
		u.Role = "user" // Default to user if invalid role
	}

	return nil
}

// Validate performs validation on user fields
func (u *User) Validate() error {
	if u.FirstName == "" {
		return errors.New("first name is required")
	}

	if u.LastName == "" {
		return errors.New("last name is required")
	}

	if u.Email == "" {
		return errors.New("email is required")
	}

	if u.Username == "" {
		return errors.New("username is required")
	}

	// Could add more validation like email format, password strength, etc.
	return nil
}

// CheckPassword verifies the provided password against the stored hash
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// sanitizeURL cleans a URL input
func sanitizeURL(url string) string {
	url = html.EscapeString(strings.TrimSpace(url))

	// Only allow http/https URLs
	if url != "" && !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	return url
}

// FullName returns the user's full name
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// UpdatePassword securely updates a user's password
func (u *User) UpdatePassword(newPassword string) error {
	// Validate password strength (example: at least 8 chars)
	if len(newPassword) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	return nil
}

// UserResponse is a safe representation of User for API responses
type UserResponse struct {
	ID        uint      `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	FullName  string    `json:"full_name"`
	Username  string    `json:"username"`
	Email     string    `json:"email,omitempty"` // Only include for the user themselves or admins
	Bio       string    `json:"bio,omitempty"`
	Avatar    string    `json:"avatar,omitempty"`
	Role      string    `json:"role,omitempty"` // Only include for admins
	CreatedAt time.Time `json:"created_at"`

	// Social profiles
	GithubURL   string `json:"github_url,omitempty"`
	LinkedinURL string `json:"linkedin_url,omitempty"`
	TwitterURL  string `json:"twitter_url,omitempty"`
	WebsiteURL  string `json:"website_url,omitempty"`
}

// ToResponse converts a User to a safe UserResponse
func (u *User) ToResponse(includePrivate bool) UserResponse {
	response := UserResponse{
		ID:          u.ID,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		FullName:    u.FullName(),
		Username:    u.Username,
		Bio:         u.Bio,
		Avatar:      u.Avatar,
		CreatedAt:   u.CreatedAt,
		GithubURL:   u.GithubURL,
		LinkedinURL: u.LinkedinURL,
		TwitterURL:  u.TwitterURL,
		WebsiteURL:  u.WebsiteURL,
	}

	// Only include private data if authorized
	if includePrivate {
		response.Email = u.Email
		response.Role = u.Role
	}

	return response
}

// RecordLoginAttempt updates login security metrics
func (u *User) RecordLoginAttempt(successful bool) {
	if successful {
		u.LastLogin = time.Now()
		u.FailedAttempts = 0
		u.LockoutEnd = nil
	} else {
		u.FailedAttempts++

		// Implement account lockout after 5 failed attempts
		if u.FailedAttempts >= 5 {
			lockoutEnd := time.Now().Add(15 * time.Minute)
			u.LockoutEnd = &lockoutEnd
		}
	}
}

// IsLockedOut checks if the account is currently locked out
func (u *User) IsLockedOut() bool {
	if u.LockoutEnd == nil {
		return false
	}
	return time.Now().Before(*u.LockoutEnd)
}

// GeneratePasswordResetToken creates a secure token for password reset
func (u *User) GeneratePasswordResetToken() (string, error) {
	// Generate a random token (in a real app, use a secure random generator)
	token := "reset-" + strconv.FormatInt(time.Now().UnixNano(), 10)

	// Set expiration to 24 hours from now
	expiration := time.Now().Add(24 * time.Hour)

	u.ResetToken = token
	u.ResetTokenExp = &expiration

	return token, nil
}

// ValidateResetToken checks if a reset token is valid
func (u *User) ValidateResetToken(token string) bool {
	return u.ResetToken == token &&
		u.ResetTokenExp != nil &&
		time.Now().Before(*u.ResetTokenExp)
}
