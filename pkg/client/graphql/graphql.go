package graphql

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
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type GraphQLClient struct {
	path       string
	logger     *zap.SugaredLogger
	HTTPClient *http.Client
	config     *conf.Config
}

func addHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "ola"))
}

func (c GraphQLClient) makeRequest(query string, variables map[string]interface{}, payload interface{}) error {
	body, err := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.path, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	addHeaders(req)
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}
	if payload != nil {
		err = json.NewDecoder(res.Body).Decode(payload)
	}
	return err
}

func (c GraphQLClient) CreateTask(task *models.Task) (int, error) {
	query := `
	mutation CreateTask($input: TaskInput!) {
		createTask(input: $input) {
			id
		}
	}`
	variables := map[string]interface{}{
		"input": task,
	}
	response := struct {
		Data struct {
			CreateTask models.Task `json:"createTask"`
		} `json:"data"`
	}{}
	err := c.makeRequest(query, variables, &response)
	return response.Data.CreateTask.ID, err
}

func (c GraphQLClient) CreatePomodoro(taskID int, pomodoro models.Pomodoro) error {
	query := `
	mutation CreatePomodoro($taskID: Int!, $input: PomodoroInput!) {
		createPomodoro(taskID: $taskID, input: $input) {
			id
		}
	}`
	variables := map[string]interface{}{
		"taskID": taskID,
		"input":  pomodoro,
	}
	return c.makeRequest(query, variables, nil)
}

func (c GraphQLClient) DeleteTaskByID(taskID int) error {
	query := `
	mutation DeleteTask($id: Int!) {
		deleteTask(id: $id)
	}`
	variables := map[string]interface{}{
		"id": taskID,
	}
	return c.makeRequest(query, variables, nil)
}

func (c GraphQLClient) GetServerStatus() (*models.Status, error) {
	query := `
	query {
		status {
			runningTasks {
				id
				title
			}
			message
		}
	}`
	response := struct {
		Data struct {
			Status models.Status `json:"status"`
		} `json:"data"`
	}{}
	err := c.makeRequest(query, nil, &response)
	return &response.Data.Status, err
}

func (c GraphQLClient) GetTaskList() (*models.List, error) {
	query := `
	query {
		tasks {
			id
			title
			description
		}
	}`
	response := struct {
		Data struct {
			Tasks models.List `json:"tasks"`
		} `json:"data"`
	}{}
	err := c.makeRequest(query, nil, &response)
	return &response.Data.Tasks, err
}

func (c GraphQLClient) GetTask(taskID int) (*models.Task, error) {
	query := `
	query Task($id: Int!) {
		task(id: $id) {
			id
			title
			description
			pomodoros {
				id
				duration
			}
		}
	}`
	variables := map[string]interface{}{
		"id": taskID,
	}
	response := struct {
		Data struct {
			Task models.Task `json:"task"`
		} `json:"data"`
	}{}
	err := c.makeRequest(query, variables, &response)
	return &response.Data.Task, err
}

func (c GraphQLClient) StartTask(taskID int) error {
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

func (c GraphQLClient) UpdateStatus(status *models.Status) error {
	query := `
	mutation UpdateStatus($input: StatusInput!) {
		updateStatus(input: $input) {
			message
		}
	}`
	variables := map[string]interface{}{
		"input": status,
	}
	return c.makeRequest(query, variables, nil)
}

func (c GraphQLClient) Close() error {
	return nil
}

func (c GraphQLClient) Init(config *conf.Config) (*GraphQLClient, error) {
	return &GraphQLClient{
		config: config,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
		logger: zap.S().With("package", "graphqlclient"),
		path:   viper.GetString("server.path"),
	}, nil
}

func maybe(err error, logger *zap.SugaredLogger) {
	if err != nil {
		logger.Fatalf("Error:\n%s\n", err)
		os.Exit(1)
	}
}

func (c GraphQLClient) Config() *conf.Config {
	return c.config
}
