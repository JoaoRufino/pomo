package runner

import (
	"time"

	"github.com/joaorufino/pomo/pkg/core"
	"github.com/joaorufino/pomo/pkg/core/models"
)

func NewRunner(client core.Client, task *models.Task) (core.Runner, error) {
	return NewTaskRunner(client, task)
}

type TaskRunner struct {
	count        int
	taskID       int
	taskMessage  string
	nPomodoros   int
	origDuration time.Duration
	state        models.State
	client       core.Client
	started      time.Time
	pause        chan bool
	toggle       chan bool
	notifier     models.Notifier
	duration     time.Duration
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
		// new ticker for periodic status updates
		ticker := time.NewTicker(5 * time.Second)
		// Record our started time
		t.started = pomodoro.Start
		t.client.UpdateStatus(t.Status())
	loop:
		select {
		case <-timer.C:
			t.SetState(models.BREAKING)
			t.count++
			t.client.UpdateStatus(t.Status())
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
			t.client.UpdateStatus(t.Status())
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
			t.client.UpdateStatus(t.Status())
			goto loop
		case <-ticker.C:
			t.client.UpdateStatus(t.Status())
			goto loop
		}
		pomodoro.End = time.Now()
		err := t.client.CreatePomodoro(t.taskID, *pomodoro)
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
	t.notifier.Notify("Pomo", "Pomo session has been completed!")
	t.SetState(models.COMPLETE)
	t.client.UpdateStatus(t.Status())
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
func (t *TaskRunner) SetStatus(status models.Status) {
	t.state = status.State
	t.count = status.Count
	t.nPomodoros = status.NPomodoros
}

func NewTaskRunner(client core.Client, task *models.Task) (*TaskRunner, error) {
	config := client.Config()
	tr := &TaskRunner{
		taskID:       task.ID,
		taskMessage:  task.Message,
		nPomodoros:   task.NPomodoros,
		origDuration: task.Duration,
		client:       client,
		state:        models.State(0),
		pause:        make(chan bool),
		toggle:       make(chan bool),
		notifier:     models.NewXnotifier(config.Runner.IconPath),
		duration:     task.Duration,
	}
	return tr, nil
}
