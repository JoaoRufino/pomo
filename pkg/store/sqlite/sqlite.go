package sqlite

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/joao.rufino/pomo/pkg/core"
	"github.com/joao.rufino/pomo/pkg/core/models"
	_ "github.com/mattn/go-sqlite3"
)

// 2018-01-16 19:05:21.752851759+08:00
const datetimeFmt = "2006-01-02 15:04:05.999999999-07:00"

type SqliteStoreFunc func(tx *sql.Tx) error

type SqliteStore struct {
	db *sql.DB
}

func NewStore(path string) (core.Store, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	return &SqliteStore{db: db}, nil
}

// With applies all of the given functions with
// a single transaction, rolling back on failure
// and commiting on success.
func (s SqliteStore) With(fns ...func(tx *sql.Tx) error) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	for _, fn := range fns {
		err = fn(tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (s SqliteStore) TaskSave(context context.Context, task *models.Task) (int, error) {
	var taskID int

	err := s.With(func(tx *sql.Tx) error {
		_, err := tx.Exec(
			"INSERT INTO task (message,pomodoros,duration,tags) VALUES ($1,$2,$3,$4)",
			task.Message,
			task.NPomodoros,
			task.Duration.String(),
			strings.Join(task.Tags, ","))
		if err != nil {
			return err
		}
		err = tx.QueryRow("SELECT last_insert_rowid() FROM task").Scan(&taskID)
		if err != nil {
			return err
		}
		err = tx.QueryRow("SELECT last_insert_rowid() FROM task").Scan(&taskID)
		if err != nil {
			return err
		}
		return nil
	})
	return taskID, err
}

func (s SqliteStore) GetAllTasks(context context.Context) ([]models.Task, error) {
	tasks := []models.Task{}

	err := s.With(func(tx *sql.Tx) error {
		rows, err := tx.Query(`SELECT rowid,message,pomodoros,duration,tags FROM task`)
		if err != nil {
			return err
		}
		for rows.Next() {
			var (
				tags        string
				strDuration string
			)
			task := &models.Task{Pomodoros: []*models.Pomodoro{}}
			err = rows.Scan(&task.ID, &task.Message, &task.NPomodoros, &strDuration, &tags)
			if err != nil {
				return err
			}
			duration, _ := time.ParseDuration(strDuration)
			task.Duration = duration
			if tags != "" {
				task.Tags = strings.Split(tags, ",")
			}
			pomodoros, err := s.PomodoroGetByTaskID(context, task.ID)
			if err != nil {
				return err
			}
			task.Pomodoros = append(task.Pomodoros, pomodoros...)
			tasks = append(tasks, *task)
		}
		return nil
	})
	return tasks, err
}

func (s SqliteStore) TaskDeleteByID(context context.Context, taskID int) error {

	err := s.With(func(tx *sql.Tx) error {
		_, err := tx.Exec("DELETE FROM task WHERE rowid = $1", &taskID)
		if err != nil {
			return err
		}
		_, err = tx.Exec("DELETE FROM pomodoro WHERE task_id = $1", &taskID)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s SqliteStore) TaskGetByID(context context.Context, taskID int) (*models.Task, error) {
	task := &models.Task{}

	err := s.With(func(tx *sql.Tx) error {
		var (
			tags        string
			strDuration string
		)
		err := tx.QueryRow(`SELECT rowid,message,pomodoros,duration,tags FROM task WHERE rowid = $1`, &taskID).
			Scan(&task.ID, &task.Message, &task.NPomodoros, &strDuration, &tags)
		if err != nil {
			return nil
		}
		duration, _ := time.ParseDuration(strDuration)
		task.Duration = duration
		if tags != "" {
			task.Tags = strings.Split(tags, ",")
		}
		return nil
	})
	return task, err
}

func (s SqliteStore) PomodoroSave(context context.Context, taskID int, pomodoro *models.Pomodoro) error {
	err := s.With(func(tx *sql.Tx) error {
		_, err := tx.Exec(
			`INSERT INTO pomodoro (task_id, start, end) VALUES ($1, $2, $3)`,
			taskID,
			pomodoro.Start,
			pomodoro.End,
		)
		return err
	})
	return err
}

func (s SqliteStore) PomodoroGetByTaskID(context context.Context, taskID int) ([]*models.Pomodoro, error) {
	pomodoros := []*models.Pomodoro{}
	err := s.With(func(tx *sql.Tx) error {
		rows, err := tx.Query(`SELECT start,end FROM pomodoro WHERE task_id = $1`, &taskID)
		if err != nil {
			return nil
		}
		for rows.Next() {
			var (
				startStr string
				endStr   string
			)
			pomodoro := &models.Pomodoro{}
			err = rows.Scan(&startStr, &endStr)
			if err != nil {
				return err
			}
			start, _ := time.Parse(datetimeFmt, startStr)
			end, _ := time.Parse(datetimeFmt, endStr)
			pomodoro.Start = start
			pomodoro.End = end
			pomodoros = append(pomodoros, pomodoro)
		}
		return nil
	})
	return pomodoros, err
}

func (s SqliteStore) PomodoroDeleteByTaskID(context context.Context, taskID int) error {
	err := s.With(func(tx *sql.Tx) error {
		_, err := tx.Exec("DELETE FROM pomodoro WHERE task_id = $1", &taskID)
		return err
	})
	return err
}

func (s SqliteStore) Close() error { return s.db.Close() }

func (s SqliteStore) InitDB() error {
	stmt := `
    CREATE TABLE task (
	message TEXT,
	pomodoros INTEGER,
	duration TEXT,
	tags TEXT
    );
    CREATE TABLE pomodoro (
	task_id INTEGER,
	start DATETTIME,
	end DATETTIME
    );
    `
	_, err := s.db.Exec(stmt)
	return err
}
