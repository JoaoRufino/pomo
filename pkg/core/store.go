package core

import (
	"context"

	"github.com/joao.rufino/pomo/pkg/core/models"
)

// Store is the persistent store of tasks
type Store interface {
	TaskGetByID(ctx context.Context, id int) (*models.Task, error)
	GetAllTasks(ctx context.Context) ([]models.Task, error)
	TaskSave(ctx context.Context, task *models.Task) (int, error)
	TaskDeleteByID(ctx context.Context, id int) error

	PomodoroGetByTaskID(ctx context.Context, id int) ([]*models.Pomodoro, error)
	PomodoroSave(ctx context.Context, taskID int, pomodoro *models.Pomodoro) error
	PomodoroDeleteByTaskID(ctx context.Context, id int) error
	Close() error
	InitDB() error
}
