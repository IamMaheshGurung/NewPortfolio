package services

import (
	"errors"

	"github.com/IamMaheshGurung/NewPortfolio/backend/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ErrSkillNotFound = errors.New("skill not found")
	ErrInvalidSkill  = errors.New("invalid skill data")
)

type SkillService struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func NewSkillService(db *gorm.DB, logger *zap.Logger) *SkillService {
	return &SkillService{
		DB:     db,
		Logger: logger,
	}
}

// CRUD Operations
func (ss *SkillService) Create(skill *models.Skill) error {
	if err := ss.DB.Create(skill).Error; err != nil {
		ss.Logger.Error("failed to create skill", zap.Error(err))
		return err
	}
	return nil
}

func (ss *SkillService) GetByID(id uint) (*models.Skill, error) {
	var skill models.Skill
	if err := ss.DB.First(&skill, id).Error; err != nil {
		ss.Logger.Error("failed to get skill", zap.Uint("id", id), zap.Error(err))
		return nil, ErrSkillNotFound
	}
	return &skill, nil
}

func (ss *SkillService) Update(skill *models.Skill) error {
	if err := ss.DB.Save(skill).Error; err != nil {
		ss.Logger.Error("failed to update skill", zap.Uint("id", skill.ID), zap.Error(err))
		return err
	}
	return nil
}

func (ss *SkillService) Delete(id uint) error {
	if err := ss.DB.Delete(&models.Skill{}, id).Error; err != nil {
		ss.Logger.Error("failed to delete skill", zap.Uint("id", id), zap.Error(err))
		return err
	}
	return nil
}

// Portfolio Specific Methods
func (ss *SkillService) ListByCategory(category string) ([]models.Skill, error) {
	var skills []models.Skill
	if err := ss.DB.Where("category = ?", category).Order("order desc").Find(&skills).Error; err != nil {
		ss.Logger.Error("failed to list skills by category", zap.String("category", category), zap.Error(err))
		return nil, err
	}
	return skills, nil
}

func (ss *SkillService) ListFeatured() ([]models.Skill, error) {
	var skills []models.Skill
	if err := ss.DB.Where("featured = ?", true).Order("proficiency desc").Find(&skills).Error; err != nil {
		ss.Logger.Error("failed to list featured skills", zap.Error(err))
		return nil, err
	}
	return skills, nil
}

func (ss *SkillService) ListAll() ([]models.Skill, error) {
	var skills []models.Skill
	if err := ss.DB.Order("order desc").Find(&skills).Error; err != nil {
		ss.Logger.Error("failed to list all skills", zap.Error(err))
		return nil, err
	}
	return skills, nil
}

func (ss *SkillService) UpdateProficiency(id uint, proficiency int) error {
	if proficiency < 0 || proficiency > 100 {
		return ErrInvalidSkill
	}
	if err := ss.DB.Model(&models.Skill{}).Where("id = ?", id).Update("proficiency", proficiency).Error; err != nil {
		ss.Logger.Error("failed to update skill proficiency", zap.Uint("id", id), zap.Int("proficiency", proficiency), zap.Error(err))
		return err
	}
	return nil
}
