package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joaorufino/pomo/pkg/conf"
	"github.com/joaorufino/pomo/pkg/core/models"
	"github.com/joaorufino/pomo/pkg/runner"
	"go.uber.org/zap"
)

// Client makes requests to a listening
// pomo server to check the status of
// any currently running task session.
type RestClient struct {
	path       string
	logger     *zap.SugaredLogger
	HTTPClient *http.Client
	config     *conf.Config
}

// add requestHeaders
func addHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "ola"))
}

// makeRequest sends a message to the server
// using the protocol structure
func (c RestClient) makeRequest(req *http.Request, payload interface{}) error {
	addHeaders(req)
	res, err := c.HTTPClient.Do(req)
	maybe(err, c.logger)
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes models.Error
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errRes.Err
		}

		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}
	if payload != nil {
		err = json.NewDecoder(res.Body).Decode(payload)
	}
	return err
}

// createTask requests the creation of a task
func (c RestClient) CreateTask(task *models.Task) (int, error) {
	body, err := json.Marshal(task)
	if err != nil {
		return -1, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/tasks", c.path), bytes.NewBuffer(body))
	if err != nil {
		return -1, err
	}
	response := &models.Task{}
	err = c.makeRequest(req, response)
	maybe(err, c.logger)
	return response.ID, nil

}

// CreatePomodoro requests the server
// to append a pomodoro to a task
func (c RestClient) CreatePomodoro(taskID int, pomodoro models.Pomodoro) error {

	body, err := json.Marshal(pomodoro)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/pomodoros/%d", c.path, taskID), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	err = c.makeRequest(req, nil)
	return err
}

// DeleteTaskByID requests the server
// to delete a task
func (c RestClient) DeleteTaskByID(taskID int) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/tasks/%d", c.path, taskID), nil)
	if err != nil {
		return err
	}

	err = c.makeRequest(req, nil)
	return err
}

// GetServerStatus requests the server
// to provide the status of running tasks
func (c RestClient) GetServerStatus() (*models.Status, error) {
	c.logger.Debug("received GetServerStatus request")
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/status", c.path), nil)
	if err != nil {
		return nil, err
	}

	response := &models.Status{}
	err = c.makeRequest(req, response)
	maybe(err, c.logger)
	return response, nil
}

// GetTaskList requests the server
// to provide the list all tasks
func (c RestClient) GetTaskList() (*models.List, error) {
	c.logger.Debug("received GetTaskList request")

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/tasks", c.path), nil)
	if err != nil {
		return nil, err
	}

	response := &models.ListResults{}
	err = c.makeRequest(req, response)
	if err != nil {
		return nil, err
	}
	return &response.Results, nil
}

// GetTask requests the server
// to provide all info on specific task
func (c RestClient) GetTask(taskID int) (*models.Task, error) {
	c.logger.Debug("received GetTask request")
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/tasks/%d", c.path, taskID), nil)
	if err != nil {
		return nil, err
	}

	response := &models.Task{}
	err = c.makeRequest(req, response)
	maybe(err, c.logger)
	return response, nil
}

// StartTask starts a pomodoro
func (c RestClient) StartTask(taskID int) error {
	task, err := c.GetTask(taskID)
	maybe(err, c.logger)
	r, err := runner.NewRunner(c, task)
	maybe(err, c.logger)
	r.Start()
	r.StartUI()
	return nil
}

// UpdateStatus sends a status update to the server
func (c RestClient) UpdateStatus(status *models.Status) error {

	body, err := json.Marshal(status)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/status", c.path), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	err = c.makeRequest(req, nil)
	return err
}

func (c RestClient) Close() error {
	//
	return nil
}

func (c RestClient) Init(config *conf.Config) (*RestClient, error) {

	return &RestClient{
		HTTPClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
		logger: zap.S().With("package", "restclient"),
		path:   fmt.Sprintf("http://%s:%s/%s", config.Server.RestHost, config.Server.RestPort, config.Server.RestPath),
		config: config,
	}, nil
}

func maybe(err error, logger *zap.SugaredLogger) {
	if err != nil {
		logger.Fatalf("Error:\n%s\n", err)
		os.Exit(1)
	}
}

func (c RestClient) Config() *conf.Config {
	return c.config
}
