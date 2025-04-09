package services

import (
	"time"

	"github.com/IamMaheshGurung/NewPortfolio/backend/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type StatusService struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func NewStatusService(db *gorm.DB, logger *zap.Logger) *StatusService {
	return &StatusService{
		DB:     db,
		Logger: logger,
	}
}

// CheckSystem performs a complete system health check
func (s *StatusService) CheckSystem() (*models.StatusResponse, error) {
	response := &models.StatusResponse{
		Checks:    make(map[models.StatusType]*models.Status),
		UpdatedAt: time.Now(),
	}

	// Check Database
	dbStatus := s.checkDatabase()
	response.Checks[models.StatusTypeDatabase] = dbStatus

	// Check API
	apiStatus := s.checkAPI()
	response.Checks[models.StatusTypeAPI] = apiStatus

	// Determine overall status
	response.Overall = s.determineOverallStatus(response.Checks)

	return response, nil
}

// Record a new status check
func (s *StatusService) RecordStatus(check *models.StatusCheck) error {
	status := &models.Status{
		Type:      check.Type,
		Status:    check.Status,
		Message:   check.Message,
		Latency:   check.Latency,
		CheckedAt: time.Now(),
	}

	if err := s.DB.Create(status).Error; err != nil {
		s.Logger.Error("failed to record status",
			zap.String("type", string(check.Type)),
			zap.Error(err))
		return err
	}

	return nil
}

// Get latest status for each component
func (s *StatusService) GetLatestStatus() (map[models.StatusType]*models.Status, error) {
	var statuses []*models.Status
	result := make(map[models.StatusType]*models.Status)

	err := s.DB.Where("checked_at IN (?)",
		s.DB.Model(&models.Status{}).
			Select("MAX(checked_at)").
			Group("type")).
		Find(&statuses).Error

	if err != nil {
		s.Logger.Error("failed to get latest status", zap.Error(err))
		return nil, err
	}

	for _, status := range statuses {
		result[status.Type] = status
	}

	return result, nil
}

// Check database health
func (s *StatusService) checkDatabase() *models.Status {
	start := time.Now()
	status := &models.Status{
		Type:      models.StatusTypeDatabase,
		CheckedAt: time.Now(),
	}

	sqlDB, err := s.DB.DB()
	if err != nil {
		status.Status = "offline"
		status.Message = "Database connection error"
		return status
	}

	if err := sqlDB.Ping(); err != nil {
		status.Status = "offline"
		status.Message = "Database ping failed"
		return status
	}

	status.Status = "online"
	status.Message = "Database is healthy"
	status.Latency = time.Since(start).Milliseconds()

	return status
}

// Check API health
func (s *StatusService) checkAPI() *models.Status {
	status := &models.Status{
		Type:      models.StatusTypeAPI,
		Status:    "online",
		Message:   "API is responding",
		CheckedAt: time.Now(),
	}
	return status
}

// Determine overall system status
func (s *StatusService) determineOverallStatus(checks map[models.StatusType]*models.Status) string {
	for _, status := range checks {
		if status.Status == "offline" {
			return "degraded"
		}
	}
	return "healthy"
}

// Get uptime statistics
func (s *StatusService) GetUptimeStats(duration time.Duration) (map[models.StatusType]float64, error) {
	stats := make(map[models.StatusType]float64)
	since := time.Now().Add(-duration)

	var counts []struct {
		Type   models.StatusType
		Total  int64
		Online int64
	}

	err := s.DB.Raw(`
        SELECT 
            type,
            COUNT(*) as total,
            SUM(CASE WHEN status = 'online' THEN 1 ELSE 0 END) as online
        FROM system_status
        WHERE checked_at >= ?
        GROUP BY type`, since).
		Scan(&counts).Error

	if err != nil {
		return nil, err
	}

	for _, count := range counts {
		if count.Total > 0 {
			stats[count.Type] = float64(count.Online) / float64(count.Total) * 100
		}
	}

	return stats, nil
}
