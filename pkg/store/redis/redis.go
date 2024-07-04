package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/joaorufino/pomo/pkg/core/models"
	"go.uber.org/zap"
)

type RedisStore struct {
	client *redis.Client
	logger *zap.SugaredLogger
}

func NewStore(addr, password string, db int, logger *zap.SugaredLogger) (*RedisStore, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Test connection
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		logger.Error("failed to connect to Redis", zap.Error(err))
		return nil, err
	}

	return &RedisStore{
		client: client,
		logger: logger,
	}, nil
}

func (s *RedisStore) InitDB() error {
	// No specific initialization required for Redis
	return nil
}

func (s *RedisStore) TaskSave(ctx context.Context, task *models.Task) (int, error) {
	key := getTaskKey(task.ID)
	data, err := json.Marshal(task)
	if err != nil {
		s.logger.Error("failed to marshal task", zap.Error(err))
		return 0, err
	}

	err = s.client.Set(ctx, key, data, 0).Err()
	if err != nil {
		s.logger.Error("failed to save task", zap.Error(err))
		return 0, err
	}
	s.logger.Info("task saved successfully", zap.Int("taskID", task.ID))
	return task.ID, nil
}

func (s *RedisStore) GetAllTasks(ctx context.Context) (models.List, error) {
	var tasks models.List

	keys, err := s.client.Keys(ctx, "task:*").Result()
	if err != nil {
		s.logger.Error("failed to get task keys", zap.Error(err))
		return nil, err
	}

	for _, key := range keys {
		data, err := s.client.Get(ctx, key).Result()
		if err != nil {
			s.logger.Error("failed to get task by key", zap.String("key", key), zap.Error(err))
			continue
		}

		var task models.Task
		if err := json.Unmarshal([]byte(data), &task); err != nil {
			s.logger.Error("failed to unmarshal task", zap.String("key", key), zap.Error(err))
			continue
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *RedisStore) TaskDeleteByID(ctx context.Context, taskID int) error {
	key := getTaskKey(taskID)
	err := s.client.Del(ctx, key).Err()
	if err != nil {
		s.logger.Error("failed to delete task by ID", zap.Int("taskID", taskID), zap.Error(err))
		return err
	}
	s.logger.Info("task deleted successfully", zap.Int("taskID", taskID))
	return nil
}

func (s *RedisStore) TaskGetByID(ctx context.Context, taskID int) (*models.Task, error) {
	key := getTaskKey(taskID)
	data, err := s.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, models.ErrNotFound
		}
		s.logger.Error("failed to get task by ID", zap.Int("taskID", taskID), zap.Error(err))
		return nil, err
	}

	var task models.Task
	if err := json.Unmarshal([]byte(data), &task); err != nil {
		s.logger.Error("failed to unmarshal task", zap.Int("taskID", taskID), zap.Error(err))
		return nil, err
	}
	return &task, nil
}

func (s *RedisStore) PomodoroSave(ctx context.Context, taskID int, pomodoro *models.Pomodoro) error {
	key := getPomodoroKey(taskID, pomodoro.ID)
	data, err := json.Marshal(pomodoro)
	if err != nil {
		s.logger.Error("failed to marshal pomodoro", zap.Error(err))
		return err
	}

	err = s.client.Set(ctx, key, data, 0).Err()
	if err != nil {
		s.logger.Error("failed to save pomodoro", zap.Error(err))
		return err
	}
	s.logger.Info("pomodoro saved successfully", zap.Int("taskID", taskID))
	return nil
}

func (s *RedisStore) PomodoroGetByTaskID(ctx context.Context, taskID int) ([]models.Pomodoro, error) {
	var pomodoros []models.Pomodoro

	keys, err := s.client.Keys(ctx, fmt.Sprintf("pomodoro:%d:*", taskID)).Result()
	if err != nil {
		s.logger.Error("failed to get pomodoro keys", zap.Int("taskID", taskID), zap.Error(err))
		return nil, err
	}

	for _, key := range keys {
		data, err := s.client.Get(ctx, key).Result()
		if err != nil {
			s.logger.Error("failed to get pomodoro by key", zap.String("key", key), zap.Error(err))
			continue
		}

		var pomodoro models.Pomodoro
		if err := json.Unmarshal([]byte(data), &pomodoro); err != nil {
			s.logger.Error("failed to unmarshal pomodoro", zap.String("key", key), zap.Error(err))
			continue
		}

		pomodoros = append(pomodoros, pomodoro)
	}

	return pomodoros, nil
}

func (s *RedisStore) PomodoroDeleteByTaskID(ctx context.Context, taskID int) error {
	keys, err := s.client.Keys(ctx, fmt.Sprintf("pomodoro:%d:*", taskID)).Result()
	if err != nil {
		s.logger.Error("failed to get pomodoro keys", zap.Int("taskID", taskID), zap.Error(err))
		return err
	}

	for _, key := range keys {
		if err := s.client.Del(ctx, key).Err(); err != nil {
			s.logger.Error("failed to delete pomodoro by key", zap.String("key", key), zap.Error(err))
			continue
		}
	}

	s.logger.Info("pomodoros deleted successfully", zap.Int("taskID", taskID))
	return nil
}

func (s *RedisStore) Close() error {
	return s.client.Close()
}

func getTaskKey(taskID int) string {
	return fmt.Sprintf("task:%d", taskID)
}

func getPomodoroKey(taskID, pomodoroID int) string {
	return fmt.Sprintf("pomodoro:%d:%d", taskID, pomodoroID)
}
