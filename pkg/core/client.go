package core

import (
	"github.com/joaorufino/pomo/pkg/core/models"
	"github.com/knadh/koanf"
)

type Client interface {
	CreateTask(task *models.Task) (int, error)
	Close() error
	DeleteTaskByID(taskID int) error
	GetServerStatus() (*models.Status, error)
	GetTaskList() (*models.List, error)
	StartTask(taskID int) error
	UpdateStatus(status *models.Status) error
	Config() *koanf.Koanf
	CreatePomodoro(taskID int, pomodoro models.Pomodoro) error
}
