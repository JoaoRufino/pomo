package unix

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/joao.rufino/pomo/pkg/core/models"
	"github.com/joao.rufino/pomo/pkg/runner"
	"github.com/knadh/koanf"
	"go.uber.org/zap"
)

// Client makes requests to a listening
// pomo server to check the status of
// any currently running task session.
type UnixClient struct {
	path   string
	config *koanf.Koanf
	logger *zap.SugaredLogger
}

// makeRequest sends a message to the server
//using the protocol structure
func (c UnixClient) makeRequest(cid models.CmdID, payload interface{}) []byte {
	conn, err := net.Dial("unix", c.path)
	maybe(err, c.logger)
	defer conn.Close()
	rxCh := make(chan []byte)
	raw, _ := json.Marshal(&models.Protocol{Cid: cid, Payload: payload})
	conn.Write(raw)
	go c.read(rxCh, conn)
	return <-rxCh
}

// read from socket
func (c UnixClient) read(rxCh chan []byte, conn net.Conn) {
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	conn.Close()
	maybe(err, c.logger)
	if n > 0 {
		rxCh <- buf[0:n]
	} else {
		rxCh <- nil
	}
}

// createTask requests the creation of a task
func (c UnixClient) CreateTask(task *models.Task) (int, error) {
	cid := models.Cmd_CreateTask
	message := c.makeRequest(cid, task)
	response := models.Protocol{Payload: 0}
	json.Unmarshal(message, &response)
	if response.Cid != cid {
		return -1, fmt.Errorf(models.ErrWrongMessageType, response.Cid, cid)
	} else {
		taskID, ok := response.Payload.(float64) //float64, for JSON numbers
		valid(ok, c.logger, response.Payload, message)
		return int(taskID), nil
	}

}

// CreatePomodoro requests the server
// to append a pomodoro to a task
func (c UnixClient) CreatePomodoro(taskID int, pomodoro models.Pomodoro) error {
	c.logger.Debug("starting CreatePomodoro request")
	cid := models.Cmd_CreatePomodoro
	compose := models.PomodoroWithID{TaskID: taskID, Pomodoro: pomodoro}
	message := c.makeRequest(cid, &compose)
	response := models.Protocol{Payload: ""}
	json.Unmarshal(message, &response)
	if response.Cid != cid {
		return fmt.Errorf(models.ErrWrongMessageType, response.Cid, cid)
	} else {
		message, ok := response.Payload.(string)
		valid(ok, c.logger, response.Payload, message)
		if len(message) != 0 {
			return errors.New(message)
		} else {
			return nil
		}
	}

}

// DeleteTaskByID requests the server
// to delete a task
func (c UnixClient) DeleteTaskByID(taskID int) error {
	c.logger.Debug("starting DeleteTaskByID request")
	cid := models.Cmd_DeleteTask
	message := c.makeRequest(cid, &taskID)
	response := models.Protocol{Payload: ""}
	json.Unmarshal(message, &response)
	if response.Cid != cid {
		return fmt.Errorf(models.ErrWrongMessageType, response.Cid, cid)
	} else {
		message, ok := response.Payload.(string)
		valid(ok, c.logger, response.Payload, message)
		if len(message) != 0 {
			return errors.New(message)
		} else {
			return nil
		}
	}
}

// GetServerStatus requests the server
// to provide the status of running tasks
func (c UnixClient) GetServerStatus() (*models.Status, error) {
	c.logger.Debug("received GetServerStatus request")
	cid := models.Cmd_GetServerStatus
	message := c.makeRequest(cid, nil)
	response := models.Protocol{Payload: &models.Status{}}
	json.Unmarshal(message, &response)
	if response.Cid != cid {
		return nil, fmt.Errorf(models.ErrWrongMessageType, response.Cid, cid)
	} else {
		status, ok := response.Payload.(*models.Status)
		valid(ok, c.logger, response.Payload, message)
		return status, nil
	}
}

// GetTaskList requests the server
// to provide the list all tasks
func (c UnixClient) GetTaskList() (*models.List, error) {
	c.logger.Debug("received GetTaskList request")
	cid := models.Cmd_GetList
	message := c.makeRequest(cid, nil)
	response := models.Protocol{Payload: &models.List{}}
	json.Unmarshal(message, &response)
	if response.Cid != cid {
		return nil, fmt.Errorf(models.ErrWrongMessageType, response.Cid, cid)
	} else {
		list, ok := response.Payload.(*models.List)
		valid(ok, c.logger, response.Payload, message)
		return list, nil
	}
}

// GetTask requests the server
// to provide all info on specific task
func (c UnixClient) GetTask(taskID int) (*models.Task, error) {
	cid := models.Cmd_GetTask
	message := c.makeRequest(cid, taskID)
	response := models.Protocol{Payload: &models.Task{}}
	json.Unmarshal(message, &response)
	if response.Cid != cid {
		return nil, fmt.Errorf(models.ErrWrongMessageType, response.Cid, cid)
	} else {
		task, ok := response.Payload.(*models.Task)
		valid(ok, c.logger, response.Payload, message)
		return task, nil
	}
}

// StartTask starts a pomodoro
func (c UnixClient) StartTask(taskID int) error {
	task, err := c.GetTask(taskID)
	maybe(err, c.logger)
	r, err := runner.NewRunner(c, task)
	maybe(err, c.logger)
	r.Start()
	r.StartUI()
	return nil
}

// UpdateStatus sends a status update to the server
func (c UnixClient) UpdateStatus(status *models.Status) error {
	cid := models.Cmd_UpdateStatus
	message := c.makeRequest(cid, status)
	response := models.Protocol{Payload: ""}
	json.Unmarshal(message, &response)
	if response.Cid != cid {
		return fmt.Errorf(models.ErrWrongMessageType, response.Cid, cid)
	} else {
		message, ok := response.Payload.(string)
		valid(ok, c.logger, response.Payload, message)
		if len(message) != 0 {
			return errors.New(message)
		} else {
			return nil
		}
	}
}

func (c UnixClient) Close() error {
	return nil
}

func (c UnixClient) Init(k *koanf.Koanf) (*UnixClient, error) {
	c.path = k.String("server.unix.socket")
	c.config = k
	c.logger = zap.S().With("package", "client")
	return &c, nil
}

func (c UnixClient) Config() *koanf.Koanf {
	return c.config
}

func maybe(err error, logger *zap.SugaredLogger) {
	if err != nil {
		logger.Fatalf("Error:\n%s\n", err)
		os.Exit(1)
	}
}

func valid(ok bool, logger *zap.SugaredLogger, expected interface{}, received interface{}) {
	if !ok {
		logger.Fatalf(models.ErrWrongDataType, expected, received)
	}
}
