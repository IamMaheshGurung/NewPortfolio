package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/IamMaheshGurung/NewPortfolio/backend/models"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	// Auth methods
	Login(email, password string) (*models.User, error)
	ValidateCredentials(user *models.User, password string) error

	// Profile methods
	GetUserByID(id uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateProfile(user *models.User) error
	UpdateAvatar(userID uint, avatarURL string) error

	// Social profile methods
	UpdateSocialLinks(userID uint, github, linkedin, twitter, website string) error

	// Password management
	GenerateResetToken(email string) (string, error)
	ValidateResetToken(token string) (*models.User, error)
	ResetPassword(userID uint, newPassword string) error

	// Security methods

	IsAccountLocked(user *models.User) bool
	RecordLoginAttempt(user *models.User, successful bool) error
}

type userService struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func NewUserService(db *gorm.DB, logger *zap.Logger) UserService {
	return &userService{
		DB:     db,
		Logger: logger}
}

// Auth methods
func (us *userService) Login(email, password string) (*models.User, error) {
	var user models.User
	if err := us.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	if us.IsAccountLocked(&user) {
		return nil, errors.New("account is locked")
	}

	if err := us.ValidateCredentials(&user, password); err != nil {
		us.RecordLoginAttempt(&user, false)
		return nil, err
	}

	us.RecordLoginAttempt(&user, true)
	return &user, nil
}

func (us *userService) ValidateCredentials(user *models.User, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

// Profile methods
func (us *userService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := us.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (us *userService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := us.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (us *userService) UpdateProfile(user *models.User) error {
	return us.DB.Save(user).Error
}

func (us *userService) UpdateAvatar(userID uint, avatarURL string) error {
	return us.DB.Model(&models.User{}).Where("id = ?", userID).
		Update("avatar", avatarURL).Error
}

// Social profile methods
func (us *userService) UpdateSocialLinks(userID uint, github, linkedin, twitter, website string) error {
	updates := map[string]interface{}{
		"github_url":   github,
		"linkedin_url": linkedin,
		"twitter_url":  twitter,
		"website_url":  website,
	}
	return us.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error
}

// Password management methods
func (us *userService) GenerateResetToken(email string) (string, error) {
	user, err := us.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("user not found")
	}

	// Generate secure token
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := base64.URLEncoding.EncodeToString(b)

	// Set token expiration
	expiration := time.Now().Add(24 * time.Hour)

	// Update user with reset token
	user.ResetToken = token
	user.ResetTokenExp = &expiration
	if err := us.UpdateProfile(user); err != nil {
		return "", err
	}

	return token, nil
}

func (us *userService) ValidateResetToken(token string) (*models.User, error) {
	var user models.User
	if err := us.DB.Where("reset_token = ?", token).First(&user).Error; err != nil {
		return nil, errors.New("invalid token")
	}

	if user.ResetTokenExp == nil || time.Now().After(*user.ResetTokenExp) {
		return nil, errors.New("token expired")
	}

	return &user, nil
}

func (us *userService) ResetPassword(userID uint, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	updates := map[string]interface{}{
		"password":        string(hashedPassword),
		"reset_token":     nil,
		"reset_token_exp": nil,
	}

	return us.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error
}

// Security methods
func (us *userService) IsAccountLocked(user *models.User) bool {
	return user.LockoutEnd != nil && time.Now().Before(*user.LockoutEnd)
}

func (us *userService) RecordLoginAttempt(user *models.User, successful bool) error {
	if successful {
		updates := map[string]interface{}{
			"last_login":      time.Now(),
			"failed_attempts": 0,
			"lockout_end":     nil,
		}
		return us.DB.Model(user).Updates(updates).Error
	}

	user.FailedAttempts++
	if user.FailedAttempts >= 5 {
		lockoutEnd := time.Now().Add(15 * time.Minute)
		user.LockoutEnd = &lockoutEnd
	}

	return us.DB.Save(user).Error
}
