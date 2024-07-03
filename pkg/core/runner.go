package core

import (
	"time"

	"github.com/joaorufino/pomo/pkg/core/models"
)

type Runner interface {
	TimeRemaining() time.Duration
	SetState(state models.State)
	SetStatus(status models.Status)
	Status() *models.Status
	Toggle()
	Pause()
	Start()
	StartUI()
}
