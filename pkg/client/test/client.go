// Code generated by MockGen. DO NOT EDIT.
// Source: ../../pkg/core/client.go
//
// Generated by this command:
//
//	mockgen -source ../../pkg/core/client.go
//
// Package mock_core is a generated GoMock package.
package test

import (
	reflect "reflect"

	models "github.com/joaorufino/pomo/pkg/core/models"
	koanf "github.com/knadh/koanf"
	gomock "go.uber.org/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockClient) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockClientMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockClient)(nil).Close))
}

// Config mocks base method.
func (m *MockClient) Config() *koanf.Koanf {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Config")
	ret0, _ := ret[0].(*koanf.Koanf)
	return ret0
}

// Config indicates an expected call of Config.
func (mr *MockClientMockRecorder) Config() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Config", reflect.TypeOf((*MockClient)(nil).Config))
}

// CreatePomodoro mocks base method.
func (m *MockClient) CreatePomodoro(taskID int, pomodoro models.Pomodoro) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePomodoro", taskID, pomodoro)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreatePomodoro indicates an expected call of CreatePomodoro.
func (mr *MockClientMockRecorder) CreatePomodoro(taskID, pomodoro any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePomodoro", reflect.TypeOf((*MockClient)(nil).CreatePomodoro), taskID, pomodoro)
}

// CreateTask mocks base method.
func (m *MockClient) CreateTask(task *models.Task) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTask", task)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTask indicates an expected call of CreateTask.
func (mr *MockClientMockRecorder) CreateTask(task any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTask", reflect.TypeOf((*MockClient)(nil).CreateTask), task)
}

// DeleteTaskByID mocks base method.
func (m *MockClient) DeleteTaskByID(taskID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTaskByID", taskID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTaskByID indicates an expected call of DeleteTaskByID.
func (mr *MockClientMockRecorder) DeleteTaskByID(taskID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTaskByID", reflect.TypeOf((*MockClient)(nil).DeleteTaskByID), taskID)
}

// GetServerStatus mocks base method.
func (m *MockClient) GetServerStatus() (*models.Status, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetServerStatus")
	ret0, _ := ret[0].(*models.Status)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetServerStatus indicates an expected call of GetServerStatus.
func (mr *MockClientMockRecorder) GetServerStatus() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetServerStatus", reflect.TypeOf((*MockClient)(nil).GetServerStatus))
}

// GetTaskList mocks base method.
func (m *MockClient) GetTaskList() (*models.List, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTaskList")
	ret0, _ := ret[0].(*models.List)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTaskList indicates an expected call of GetTaskList.
func (mr *MockClientMockRecorder) GetTaskList() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTaskList", reflect.TypeOf((*MockClient)(nil).GetTaskList))
}

// StartTask mocks base method.
func (m *MockClient) StartTask(taskID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartTask", taskID)
	ret0, _ := ret[0].(error)
	return ret0
}

// StartTask indicates an expected call of StartTask.
func (mr *MockClientMockRecorder) StartTask(taskID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartTask", reflect.TypeOf((*MockClient)(nil).StartTask), taskID)
}

// UpdateStatus mocks base method.
func (m *MockClient) UpdateStatus(status *models.Status) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStatus", status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateStatus indicates an expected call of UpdateStatus.
func (mr *MockClientMockRecorder) UpdateStatus(status any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStatus", reflect.TypeOf((*MockClient)(nil).UpdateStatus), status)
}
