package unix

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/joaorufino/pomo/pkg/conf"
	"github.com/joaorufino/pomo/pkg/core/models"
	"github.com/joaorufino/pomo/pkg/runner"
	"go.uber.org/zap"
)

// Client makes requests to a listening
// pomo server to check the status of
// any currently running task session.
type UnixClient struct {
	path   string
	logger *zap.SugaredLogger
	config *conf.Config
}

// makeRequest sends a message to the server
// using the protocol structure
func (c UnixClient) makeRequest(cid models.CmdID, payload interface{}) ([]byte, error) {
	conn, err := net.Dial("unix", c.path)
	if err != nil {
		return nil, fmt.Errorf("failed to dial unix socket: %w", err)
	}
	defer conn.Close()

	rxCh := make(chan []byte)
	raw, err := json.Marshal(&models.Protocol{Cid: cid, Payload: payload})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	if _, err := conn.Write(raw); err != nil {
		return nil, fmt.Errorf("failed to write to socket: %w", err)
	}

	go c.read(rxCh, conn)
	return <-rxCh, nil
}

// read from socket
func (c UnixClient) read(rxCh chan []byte, conn net.Conn) {
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		if errors.Is(err, net.ErrClosed) {
			c.logger.Warn("connection closed by server")
		} else {
			c.logger.Error("failed to read from socket", zap.Error(err))
		}
		rxCh <- nil
		return
	}
	if n > 0 {
		rxCh <- buf[0:n]
	} else {
		rxCh <- nil
	}
}

// createTask requests the creation of a task
func (c UnixClient) CreateTask(task *models.Task) (int, error) {
	cid := models.Cmd_CreateTask
	message, err := c.makeRequest(cid, task)
	if err != nil {
		return -1, err
	}
	response := models.Protocol{Payload: 0}
	if err := json.Unmarshal(message, &response); err != nil {
		return -1, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	if response.Cid != cid {
		return -1, fmt.Errorf(models.ErrWrongMessageType, response.Cid, cid)
	}
	taskID, ok := response.Payload.(float64) // float64, for JSON numbers
	if !ok {
		return -1, fmt.Errorf(models.ErrWrongDataType, response.Payload, taskID)
	}
	return int(taskID), nil
}

// CreatePomodoro requests the server
// to append a pomodoro to a task
func (c UnixClient) CreatePomodoro(taskID int, pomodoro models.Pomodoro) error {
	c.logger.Debug("starting CreatePomodoro request")
	cid := models.Cmd_CreatePomodoro
	compose := models.PomodoroWithID{TaskID: taskID, Pomodoro: pomodoro}
	message, err := c.makeRequest(cid, &compose)
	if err != nil {
		return err
	}
	response := models.Protocol{Payload: ""}
	if err := json.Unmarshal(message, &response); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}
	if response.Cid != cid {
		return fmt.Errorf(models.ErrWrongMessageType, response.Cid, cid)
	}
	msg, ok := response.Payload.(string)
	if !ok {
		return fmt.Errorf(models.ErrWrongDataType, response.Payload, msg)
	}
	if len(msg) != 0 {
		return errors.New(msg)
	}
	return nil
}

// DeleteTaskByID requests the server
// to delete a task
func (c UnixClient) DeleteTaskByID(taskID int) error {
	c.logger.Debug("starting DeleteTaskByID request")
	cid := models.Cmd_DeleteTask
	message, err := c.makeRequest(cid, &taskID)
	if err != nil {
		return err
	}
	response := models.Protocol{Payload: ""}
	if err := json.Unmarshal(message, &response); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}
	if response.Cid != cid {
		return fmt.Errorf(models.ErrWrongMessageType, response.Cid, cid)
	}
	msg, ok := response.Payload.(string)
	if !ok {
		return fmt.Errorf(models.ErrWrongDataType, response.Payload, msg)
	}
	if len(msg) != 0 {
		return errors.New(msg)
	}
	return nil
}

// GetServerStatus requests the server
// to provide the status of running tasks
func (c UnixClient) GetServerStatus() (*models.Status, error) {
	c.logger.Debug("received GetServerStatus request")
	cid := models.Cmd_GetServerStatus
	message, err := c.makeRequest(cid, nil)
	if err != nil {
		return nil, err
	}
	response := models.Protocol{Payload: &models.Status{}}
	if err := json.Unmarshal(message, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	if response.Cid != cid {
		return nil, fmt.Errorf(models.ErrWrongMessageType, response.Cid, cid)
	}
	status, ok := response.Payload.(*models.Status)
	if !ok {
		return nil, fmt.Errorf(models.ErrWrongDataType, response.Payload, status)
	}
	return status, nil
}

// GetTaskList requests the server
// to provide the list all tasks
func (c UnixClient) GetTaskList() (*models.List, error) {
	c.logger.Debug("received GetTaskList request")
	cid := models.Cmd_GetList
	message, err := c.makeRequest(cid, nil)
	if err != nil {
		return nil, err
	}
	response := models.Protocol{Payload: &models.List{}}
	if err := json.Unmarshal(message, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	if response.Cid != cid {
		return nil, fmt.Errorf(models.ErrWrongMessageType, response.Cid, cid)
	}
	list, ok := response.Payload.(*models.List)
	if !ok {
		return nil, fmt.Errorf(models.ErrWrongDataType, response.Payload, list)
	}
	return list, nil
}

// GetTask requests the server
// to provide all info on specific task
func (c UnixClient) GetTask(taskID int) (*models.Task, error) {
	cid := models.Cmd_GetTask
	message, err := c.makeRequest(cid, taskID)
	if err != nil {
		return nil, err
	}
	response := models.Protocol{Payload: &models.Task{}}
	if err := json.Unmarshal(message, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	if response.Cid != cid {
		return nil, fmt.Errorf(models.ErrWrongMessageType, response.Cid, cid)
	}
	if response.Error != nil {
		return nil, fmt.Errorf(*response.Error, response.Cid, cid)
	}

	task, ok := response.Payload.(*models.Task)
	if !ok {
		return nil, fmt.Errorf(models.ErrWrongDataType, response.Payload, task)
	}
	return task, nil
}

// StartTask starts a pomodoro
func (c UnixClient) StartTask(taskID int) error {
	task, err := c.GetTask(taskID)
	if err != nil {
		return err
	}
	r, err := runner.NewRunner(c, task)
	if err != nil {
		return err
	}
	r.Start()
	r.StartUI()
	return nil
}

// UpdateStatus sends a status update to the server
func (c UnixClient) UpdateStatus(status *models.Status) error {
	cid := models.Cmd_UpdateStatus
	message, err := c.makeRequest(cid, status)
	if err != nil {
		return err
	}
	response := models.Protocol{Payload: ""}
	if err := json.Unmarshal(message, &response); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}
	if response.Cid != cid {
		return fmt.Errorf(models.ErrWrongMessageType, response.Cid, cid)
	}
	msg, ok := response.Payload.(string)
	if !ok {
		return fmt.Errorf(models.ErrWrongDataType, response.Payload, msg)
	}
	if len(msg) != 0 {
		return errors.New(msg)
	}
	return nil
}

func (c UnixClient) Close() error {
	return nil
}

func (c *UnixClient) Init(config *conf.Config) (*UnixClient, error) {
	c.path = config.Server.UnixSocket
	c.logger = zap.S().With("package", "client")
	c.config = config
	return c, nil
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

func (c UnixClient) Config() *conf.Config {
	return c.config
}
