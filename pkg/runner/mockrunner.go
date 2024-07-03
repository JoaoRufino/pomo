package runner

import (
	"github.com/joaorufino/pomo/pkg/core"
	"github.com/joaorufino/pomo/pkg/core/models"
)

func NewMockedTaskRunner(task *models.Task, client core.Client, notifier models.Notifier) (*TaskRunner, error) {
	tr := &TaskRunner{
		taskID:       task.ID,
		taskMessage:  task.Message,
		nPomodoros:   task.NPomodoros,
		origDuration: task.Duration,
		client:       client,
		state:        models.State(0),
		pause:        make(chan bool),
		toggle:       make(chan bool),
		notifier:     notifier,
		duration:     task.Duration,
	}
	return tr, nil
}
