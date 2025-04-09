package services

import (
	"errors"
	"time"

	"github.com/IamMaheshGurung/NewPortfolio/backend/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Error definitions
var (
	ErrBlogNotFound = errors.New("blog not found")
	ErrInvalidSlug  = errors.New("invalid slug")
)

// Service implementation
type BlogService struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

// Constructor
func NewBlogService(db *gorm.DB, logger *zap.Logger) BlogService {
	return BlogService{
		DB:     db,
		Logger: logger,
	}
}

// CRUD Implementations
func (bs *BlogService) GetByID(id uint) (*models.Blog, error) {
	var blog models.Blog
	if err := bs.DB.First(&blog, id).Error; err != nil {
		bs.Logger.Error("failed to get blog by id", zap.Uint("id", id), zap.Error(err))
		return nil, ErrBlogNotFound
	}
	return &blog, nil
}

func (bs *BlogService) GetBySlug(slug string) (*models.Blog, error) {
	var blog models.Blog
	if err := bs.DB.Where("slug = ?", slug).First(&blog).Error; err != nil {
		bs.Logger.Error("failed to get blog by slug", zap.String("slug", slug), zap.Error(err))
		return nil, ErrBlogNotFound
	}
	return &blog, nil
}

func (bs *BlogService) Create(blog *models.Blog) error {
	if err := bs.DB.Create(blog).Error; err != nil {
		bs.Logger.Error("failed to create blog", zap.Error(err))
		return err
	}
	bs.Logger.Info("blog created", zap.Uint("id", blog.ID))
	return nil
}

func (bs *BlogService) Update(blog *models.Blog) error {
	if err := bs.DB.Save(blog).Error; err != nil {
		bs.Logger.Error("failed to update blog", zap.Uint("id", blog.ID), zap.Error(err))
		return err
	}
	bs.Logger.Info("blog updated", zap.Uint("id", blog.ID))
	return nil
}

func (bs *BlogService) Delete(id uint) error {
	if err := bs.DB.Delete(&models.Blog{}, id).Error; err != nil {
		bs.Logger.Error("failed to delete blog", zap.Uint("id", id), zap.Error(err))
		return err
	}
	bs.Logger.Info("blog deleted", zap.Uint("id", id))
	return nil
}

// Query Implementations
func (bs *BlogService) ListPublished(page, limit int) ([]models.Blog, int64, error) {
	var blogs []models.Blog
	var total int64

	offset := (page - 1) * limit
	query := bs.DB.Model(&models.Blog{}).Where("published = ?", true)

	if err := query.Count(&total).Error; err != nil {
		bs.Logger.Error("failed to count published blogs", zap.Error(err))
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Order("published_at DESC").Find(&blogs).Error; err != nil {
		bs.Logger.Error("failed to list published blogs", zap.Error(err))
		return nil, 0, err
	}

	return blogs, total, nil
}

// Blog Management
func (bs *BlogService) Publish(id uint) error {
	result := bs.DB.Model(&models.Blog{}).Where("id = ?", id).Updates(map[string]interface{}{
		"published":    true,
		"published_at": time.Now(),
	})

	if result.Error != nil {
		bs.Logger.Error("failed to publish blog", zap.Uint("id", id), zap.Error(result.Error))
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrBlogNotFound
	}

	bs.Logger.Info("blog published", zap.Uint("id", id))
	return nil
}

func (bs *BlogService) IncrementViews(id uint) error {
	if err := bs.DB.Model(&models.Blog{}).Where("id = ?", id).
		UpdateColumn("views", gorm.Expr("views + ?", 1)).Error; err != nil {
		bs.Logger.Error("failed to increment views", zap.Uint("id", id), zap.Error(err))
		return err
	}
	return nil
}
