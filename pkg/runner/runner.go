package runner

import (
	"database/sql"
	"time"

	"github.com/joao.rufino/pomo/pkg/cli"
	"github.com/joao.rufino/pomo/pkg/server"
	"github.com/joao.rufino/pomo/pkg/server/models"
)

type TaskRunner struct {
	count        int
	taskID       int
	taskMessage  string
	nPomodoros   int
	origDuration time.Duration
	state        models.State
	store        *server.Store
	started      time.Time
	pause        chan bool
	toggle       chan bool
	notifier     models.Notifier
	duration     time.Duration
}

func NewTaskRunner(pomoCli cli.Cli, task *models.Task) (*TaskRunner, error) {
	store, err := server.NewStore(pomoCli.Config().String("database.path"))
	if err != nil {
		return nil, err
	}
	tr := &TaskRunner{
		taskID:       task.ID,
		taskMessage:  task.Message,
		nPomodoros:   task.NPomodoros,
		origDuration: task.Duration,
		store:        store,
		state:        models.State(0),
		pause:        make(chan bool),
		toggle:       make(chan bool),
		notifier:     models.NewXnotifier(pomoCli.Config().String("icon.path")),
		duration:     task.Duration,
	}
	return tr, nil
}

func (t *TaskRunner) Start() {
	go t.run()
}

func (t *TaskRunner) TimeRemaining() time.Duration {
	return (t.duration - time.Since(t.started)).Truncate(time.Second)
}

func (t *TaskRunner) SetState(state models.State) {
	t.state = state
}

func (t *TaskRunner) run() error {
	for t.count < t.nPomodoros {
		// Create a new pomodoro where we
		// track the start / end time of
		// of this session.
		pomodoro := &models.Pomodoro{}
		// Start this pomodoro
		pomodoro.Start = time.Now()
		// Set state to RUNNIN
		t.SetState(models.RUNNING)
		// Create a new timer
		timer := time.NewTimer(t.duration)
		// Record our started time
		t.started = pomodoro.Start
	loop:
		select {
		case <-timer.C:
			t.SetState(models.BREAKING)
			t.count++
		case <-t.toggle:
			// Catch any toggles when we
			// are not expecting them
			goto loop
		case <-t.pause:
			timer.Stop()
			// Record the remaining time of the current pomodoro
			remaining := t.TimeRemaining()
			// Change state to PAUSED
			t.SetState(models.PAUSED)
			// Wait for the user to press [p]
			<-t.pause
			// Resume the timer with previous
			// remaining time
			timer.Reset(remaining)
			// Change duration
			t.started = time.Now()
			t.duration = remaining
			// Restore state to RUNNING
			t.SetState(models.RUNNING)
			goto loop
		}
		pomodoro.End = time.Now()
		err := t.store.With(func(tx *sql.Tx) error {
			return t.store.CreatePomodoro(tx, t.taskID, *pomodoro)
		})
		if err != nil {
			return err
		}
		// All pomodoros completed
		if t.count == t.nPomodoros {
			break
		}

		t.notifier.Notify("Pomo", "It is time to take a break!")
		// Reset the duration incase it
		// was paused.
		t.duration = t.origDuration
		// User concludes the break
		<-t.toggle

	}
	t.notifier.Notify("Pomo", "Pomo session has completed!")
	t.SetState(models.COMPLETE)
	return nil
}

func (t *TaskRunner) Toggle() {
	t.toggle <- true
}

func (t *TaskRunner) Pause() {
	t.pause <- true
}

func (t *TaskRunner) Status() *models.Status {
	return &models.Status{
		State:      t.state,
		Count:      t.count,
		NPomodoros: t.nPomodoros,
		Remaining:  t.TimeRemaining(),
	}
}

func NewMockedTaskRunner(task *models.Task, store *server.Store, notifier models.Notifier) (*TaskRunner, error) {
	tr := &TaskRunner{
		taskID:       task.ID,
		taskMessage:  task.Message,
		nPomodoros:   task.NPomodoros,
		origDuration: task.Duration,
		store:        store,
		state:        models.State(0),
		pause:        make(chan bool),
		toggle:       make(chan bool),
		notifier:     notifier,
		duration:     task.Duration,
	}
	return tr, nil
}
