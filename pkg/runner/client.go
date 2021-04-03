package runner

import (
	"errors"

	"github.com/joao.rufino/pomo/pkg/server/models"
	"github.com/knadh/koanf"
)

type Client interface {
	CreateTask(task *models.Task) (int, error)
	Close() error
	DeleteTaskByID(taskID int) error
	GetServerStatus() (*models.Status, error)
	GetTaskList() (models.List, error)
	StartTask(taskID int) error
	UpdateStatus(status *models.Status) error
	Config() *koanf.Koanf
	CreatePomodoro(taskID int, pomodoro models.Pomodoro) error
}

func NewClient(k *koanf.Koanf) (Client, error) {
	switch k.String("server.type") {
	case "unix":
		unix := &UnixClient{config: k}
		return unix.Start(k.String("server.socket"))
	}
	return nil, errors.New("error creating client")
}
