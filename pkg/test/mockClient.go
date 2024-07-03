package test

import (
	"github.com/joaorufino/pomo/pkg/core"
	"github.com/joaorufino/pomo/pkg/core/models"
	"github.com/knadh/koanf"
)

type MockClient struct {
	config  *koanf.Koanf
	options MockClientOptions
}

type MockClientOptions struct {
	List   *models.List
	status *models.Status
	taskID int
}

func NewMockClient(k *koanf.Koanf, options MockClientOptions) core.Client {
	client := MockClient{}
	client.SetServerStatus(&models.Status{})
	client.SetList(&models.List{})
	if options.List != nil {
		client.options.List = options.List

	}
	if options.status != nil {
		client.options.status = options.status

	}
	if options.taskID > 0 {
		client.options.taskID = options.taskID

	}
	client.SetConfig(k)
	return &client
}

func (c *MockClient) CreateTask(task *models.Task) (int, error) {
	return c.options.taskID, nil
}

func (c *MockClient) SetTaskID(taskID int) {
	c.options.taskID = taskID
}
func (c *MockClient) Close() error {
	return nil
}
func (c *MockClient) DeleteTaskByID(taskID int) error {
	return nil
}
func (c *MockClient) GetServerStatus() (*models.Status, error) {
	return c.options.status, nil
}

func (c *MockClient) SetServerStatus(status *models.Status) {
	c.options.status = status
}
func (c *MockClient) GetTaskList() (*models.List, error) {
	return c.options.List, nil
}

func (c *MockClient) SetList(List *models.List) {
	c.options.List = List
}

func (c *MockClient) StartTask(taskID int) error {
	return nil
}
func (c *MockClient) UpdateStatus(status *models.Status) error {
	return nil
}
func (c *MockClient) Config() *koanf.Koanf {
	return c.config
}
func (c *MockClient) SetConfig(k *koanf.Koanf) {
	c.config = k
}

func (c *MockClient) CreatePomodoro(taskID int, pomodoro models.Pomodoro) error {
	return nil
}
