package services

import (
	"errors"

	"github.com/IamMaheshGurung/NewPortfolio/backend/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ErrProjectNotFound = errors.New("project not found")
	ErrInvalidProject  = errors.New("invalid project data")
)

type ProjectService struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func NewProjectService(db *gorm.DB, logger *zap.Logger) *ProjectService {
	return &ProjectService{
		DB:     db,
		Logger: logger,
	}
}

// CRUD Operations
func (ps *ProjectService) Create(project *models.Project) error {
	if err := ps.DB.Create(project).Error; err != nil {
		ps.Logger.Error("failed to create project", zap.Error(err))
		return err
	}
	ps.Logger.Info("project created", zap.Uint("id", project.ID))
	return nil
}

func (ps *ProjectService) GetByID(id uint) (*models.Project, error) {
	var project models.Project
	if err := ps.DB.First(&project, id).Error; err != nil {
		ps.Logger.Error("failed to get project", zap.Uint("id", id), zap.Error(err))
		return nil, ErrProjectNotFound
	}
	return &project, nil
}

func (ps *ProjectService) GetBySlug(slug string) (*models.Project, error) {
	var project models.Project
	if err := ps.DB.Where("slug = ?", slug).First(&project).Error; err != nil {
		ps.Logger.Error("failed to get project by slug", zap.String("slug", slug), zap.Error(err))
		return nil, ErrProjectNotFound
	}
	return &project, nil
}

func (ps *ProjectService) Update(project *models.Project) error {
	if err := ps.DB.Save(project).Error; err != nil {
		ps.Logger.Error("failed to update project", zap.Uint("id", project.ID), zap.Error(err))
		return err
	}
	ps.Logger.Info("project updated", zap.Uint("id", project.ID))
	return nil
}

func (ps *ProjectService) Delete(id uint) error {
	if err := ps.DB.Delete(&models.Project{}, id).Error; err != nil {
		ps.Logger.Error("failed to delete project", zap.Uint("id", id), zap.Error(err))
		return err
	}
	ps.Logger.Info("project deleted", zap.Uint("id", id))
	return nil
}

// Portfolio Specific Methods
func (ps *ProjectService) ListFeatured() ([]models.Project, error) {
	var projects []models.Project
	if err := ps.DB.Where("featured = ?", true).Order("order desc").Find(&projects).Error; err != nil {
		ps.Logger.Error("failed to list featured projects", zap.Error(err))
		return nil, err
	}
	return projects, nil
}

func (ps *ProjectService) ListAll() ([]models.Project, error) {
	var projects []models.Project
	if err := ps.DB.Order("order desc").Find(&projects).Error; err != nil {
		ps.Logger.Error("failed to list all projects", zap.Error(err))
		return nil, err
	}
	return projects, nil
}

func (ps *ProjectService) UpdateOrder(id uint, newOrder int) error {
	if err := ps.DB.Model(&models.Project{}).Where("id = ?", id).Update("order", newOrder).Error; err != nil {
		ps.Logger.Error("failed to update project order", zap.Uint("id", id), zap.Int("order", newOrder), zap.Error(err))
		return err
	}
	return nil
}

func (ps *ProjectService) ToggleFeatured(id uint) error {
	var project models.Project
	if err := ps.DB.First(&project, id).Error; err != nil {
		return ErrProjectNotFound
	}

	project.Featured = !project.Featured
	return ps.Update(&project)
}
