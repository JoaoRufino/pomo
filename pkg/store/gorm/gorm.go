package gormStore

import (
	"context"

	"github.com/joaorufino/pomo/pkg/core/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type GormStore struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func NewStore(db *gorm.DB, logger *zap.SugaredLogger) *GormStore {
	return &GormStore{
		db:     db,
		logger: logger,
	}
}

func (s *GormStore) TaskSave(ctx context.Context, task *models.Task) (int, error) {
	if err := s.db.WithContext(ctx).Create(task).Error; err != nil {
		s.logger.Error("failed to save task", zap.Error(err))
		return 0, err
	}
	s.logger.Info("task saved successfully", zap.Int("taskID", task.ID))
	return task.ID, nil
}

func (s *GormStore) GetAllTasks(ctx context.Context) (models.List, error) {
	var tasks models.List
	if err := s.db.WithContext(ctx).Preload("Pomodoros").Find(&tasks).Error; err != nil {
		s.logger.Error("failed to get all tasks", zap.Error(err))
		return nil, err
	}

	return tasks, nil
}

func (s *GormStore) TaskDeleteByID(ctx context.Context, taskID int) error {
	if err := s.db.WithContext(ctx).Delete(&models.Task{}, taskID).Error; err != nil {
		s.logger.Error("failed to delete task by ID", zap.Int("taskID", taskID), zap.Error(err))
		return err
	}
	s.logger.Info("task deleted successfully", zap.Int("taskID", taskID))
	return nil
}

func (s *GormStore) TaskGetByID(ctx context.Context, taskID int) (*models.Task, error) {
	var task models.Task
	if err := s.db.WithContext(ctx).Preload("Pomodoros").First(&task, taskID).Error; err != nil {
		s.logger.Error("failed to get task by ID", zap.Int("taskID", taskID), zap.Error(err))
		return nil, err
	}
	return &task, nil
}

func (s *GormStore) PomodoroSave(ctx context.Context, taskID int, pomodoro *models.Pomodoro) error {
	pomodoro.TaskID = taskID
	if err := s.db.WithContext(ctx).Create(pomodoro).Error; err != nil {
		s.logger.Error("failed to save pomodoro", zap.Int("taskID", taskID), zap.Error(err))
		return err
	}
	s.logger.Info("pomodoro saved successfully", zap.Int("taskID", taskID))
	return nil
}

func (s *GormStore) PomodoroGetByTaskID(ctx context.Context, taskID int) ([]models.Pomodoro, error) {
	var pomodoros []models.Pomodoro
	if err := s.db.WithContext(ctx).Where("task_id = ?", taskID).Find(&pomodoros).Error; err != nil {
		s.logger.Error("failed to get pomodoros by task ID", zap.Int("taskID", taskID), zap.Error(err))
		return nil, err
	}
	return pomodoros, nil
}

func (s *GormStore) PomodoroDeleteByTaskID(ctx context.Context, taskID int) error {
	if err := s.db.WithContext(ctx).Where("task_id = ?", taskID).Delete(&models.Pomodoro{}).Error; err != nil {
		s.logger.Error("failed to delete pomodoros by task ID", zap.Int("taskID", taskID), zap.Error(err))
		return err
	}
	s.logger.Info("pomodoros deleted successfully", zap.Int("taskID", taskID))
	return nil
}

func (s *GormStore) Close() error {
	if s.db == nil {
		return nil
	}
	sqlDB, err := s.db.DB()
	if err != nil {
		s.logger.Error("failed to get sql.DB from gorm.DB", zap.Error(err))
		return err
	}
	return sqlDB.Close()

}

func (s *GormStore) Init() error {
	// Auto migrate schema
	if err := s.db.AutoMigrate(&models.Task{}, &models.Pomodoro{}); err != nil {
		s.logger.Error("failed to migrate database schema", zap.Error(err))
		return err
	}
	return nil
}
