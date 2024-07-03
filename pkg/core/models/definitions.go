package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"time"

	"github.com/0xAX/notificator"
)

type Runner interface {
	Status() *Status
	SetStatus(status Status)
}
type State int

func (s State) String() string {
	switch s {
	case RUNNING:
		return "RUNNING"
	case BREAKING:
		return "BREAKING"
	case COMPLETE:
		return "COMPLETE"
	case PAUSED:
		return "PAUSED"
	}
	return ""
}

const (
	RUNNING State = iota + 1
	BREAKING
	COMPLETE
	PAUSED
)

// Wheel keeps track of an ASCII spinner
type Wheel int

func (w *Wheel) String() string {
	switch int(*w) {
	case 0:
		*w++
		return "|"
	case 1:
		*w++
		return "/"
	case 2:
		*w++
		return "-"
	case 3:
		*w = 0
		return "\\"
	}
	return ""
}

// Task describes some activity
type Task struct {
	ID      int    `json:"id" gorm:"primaryKey"`
	Message string `json:"message" "gorm:"not null"`
	// Array of completed pomodoros
	Pomodoros []*Pomodoro `json:"pomodoros" "gorm:"foreignKey:TaskID"`
	// Free-form tags associated with this task
	Tags Tags `json:"tags"`
	// Number of pomodoros for this task
	NPomodoros int `json:"n_pomodoros"`
	// Duration of each pomodoro
	Duration time.Duration `json:"duration"`
}

// Tags is a custom type to handle a slice of strings in the database
type Tags []string

// Value implements the driver.Valuer interface for database serialization
func (t Tags) Value() (driver.Value, error) {
	return json.Marshal(t)
}

// Scan implements the sql.Scanner interface for database deserialization
func (t *Tags) Scan(value interface{}) error {
	if value == nil {
		*t = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to convert database value to []byte")
	}
	return json.Unmarshal(bytes, t)
}

type ListResults struct {
	Count   int64 `json:"count"`
	Results List  `json:"results"`
}
type List []Task

// ByID is a sortable array of tasks
type ByID []*Task

func (b ByID) Len() int           { return len(b) }
func (b ByID) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByID) Less(i, j int) bool { return b[i].ID < b[j].ID }

// After returns tasks that were started after the
// provided start time.
func After(start time.Time, tasks []Task) []Task {
	filtered := []Task{}
	for _, task := range tasks {
		if len(task.Pomodoros) > 0 {
			if start.Before(task.Pomodoros[0].Start) {
				filtered = append(filtered, task)
			}
		}
	}
	return filtered
}

// Pomodoro is a unit of time to spend working
// on a single task.
type Pomodoro struct {
	ID     int       `json:"-" gorm:"primaryKey"`
	TaskID int       `json:"-" gorm:"not null"`
	Start  time.Time `json:"start" gorm:"not null"`
	End    time.Time `json:"end" gorm:"not null"`
}

// PomodoroWithID is a unit for requesting
// an update to a pomodoro to the server
type PomodoroWithID struct {
	TaskID   int
	Pomodoro Pomodoro
}

// Duration returns the runtime of the pomodoro
func (p Pomodoro) Duration() time.Duration {
	return (p.End.Sub(p.Start))
}

// Status is used to communicate the state
// of a running Pomodoro session
type Status struct {
	State      State         `json:"state"`
	Remaining  time.Duration `json:"remaining"`
	Count      int           `json:"count"`
	NPomodoros int           `json:"n_pomodoros"`
}

// Notifier sends a system notification
type Notifier interface {
	Notify(string, string) error
}

// NoopNotifier does nothing
type NoopNotifier struct{}

func (n NoopNotifier) Notify(string, string) error { return nil }

// Xnotifier can push notifications to mac, linux and windows.
type Xnotifier struct {
	*notificator.Notificator
	iconPath string
}

func NewXnotifier(iconPath string) Notifier {
	// Write the built-in tomato icon if it
	// doesn't already exist.
	_, err := os.Stat(iconPath)
	if os.IsNotExist(err) {
		raw := MustAsset("tomato-icon.png")
		_ = ioutil.WriteFile(iconPath, raw, 0644)
	}
	return Xnotifier{
		Notificator: notificator.New(notificator.Options{}),
		iconPath:    iconPath,
	}
}

// Notify sends a notification to the OS.
func (n Xnotifier) Notify(title, body string) error {
	return n.Push(title, body, n.iconPath, notificator.UR_NORMAL)
}
