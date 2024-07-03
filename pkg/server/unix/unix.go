package unix

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/joaorufino/pomo/pkg/conf"
	"github.com/joaorufino/pomo/pkg/core"
	"github.com/joaorufino/pomo/pkg/core/models"
	serverStore "github.com/joaorufino/pomo/pkg/store"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// UnixServer listens on a Unix domain socket for Pomo status requests
type UnixServer struct {
	listener net.Listener
	running  bool
	store    core.Store
	logger   *zap.SugaredLogger
	status   models.Status
}

// Init initializes the server structure
func (s *UnixServer) Init(config *conf.Config) (*UnixServer, error) {
	socketPath := config.Server.UnixSocket
	if _, err := os.Stat(socketPath); err == nil {
		conn, err := net.Dial("unix", socketPath)
		if err != nil {
			os.Remove(socketPath)
		} else {
			conn.Close()
			return nil, fmt.Errorf("socket %s is already in use", socketPath)
		}
	}

	logger := zap.S().With("package", "unix")
	store, err := serverStore.NewStore(config, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create store: %w", err)
	}

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on socket %s: %w", socketPath, err)
	}

	return &UnixServer{
		listener: listener,
		logger:   logger,
		store:    store,
		status:   models.Status{},
	}, nil
}

// Start starts the Unix server
func (s *UnixServer) Start() {
	s.running = true
	s.listen()
}

// Stop stops the Unix server
func (s *UnixServer) Stop() {
	s.running = false
	if err := s.listener.Close(); err != nil {
		s.logger.Errorf("Error closing listener: %v", err)
	}
	if err := s.store.Close(); err != nil {
		s.logger.Errorf("Error closing store: %v", err)
	}
}

func (s *UnixServer) listen() {
	s.logger.Info("Listening")
	for s.running {
		conn, err := s.listener.Accept()
		if err != nil {
			s.logger.Errorf("Failed to accept connection: %v", err)
			continue
		}

		go func(conn net.Conn) {
			defer conn.Close()
			s.handleConnection(conn)
		}(conn)
	}
}

func (s *UnixServer) handleConnection(conn net.Conn) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		s.logger.Errorf("Failed to read from connection: %v", err)
		return
	}

	if n == 0 {
		return
	}

	s.logger.Debugf("Incoming request")
	message := models.Protocol{}
	if err := json.Unmarshal(buf[:n], &message); err != nil {
		s.logger.Errorf("Failed to unmarshal request: %v", err)
		return
	}

	switch message.Cid {
	case models.Cmd_GetServerStatus:
		s.logger.Debug("Incoming status request")
		s.sendResponse(message.Cid, s.status, conn)
	case models.Cmd_GetList:
		s.logger.Debug("Incoming task list request")
		tasks, err := s.store.GetAllTasks(nil)
		if err != nil {
			s.logger.Errorf("Failed to get all tasks: %v", err)
			s.sendErrorResponse(message.Cid, "Failed to get tasks", conn)
			return
		}
		s.sendResponse(message.Cid, tasks, conn)
	case models.Cmd_CreateTask:
		s.createTask(buf[:n], conn)
	case models.Cmd_DeleteTask:
		s.deleteTask(buf[:n], conn)
	case models.Cmd_GetTask:
		s.getTask(buf[:n], conn)
	case models.Cmd_CreatePomodoro:
		s.createPomodoro(buf[:n], conn)
	case models.Cmd_UpdateStatus:
		s.updateStatus(buf[:n], conn)
	default:
		s.logger.Errorf("Unknown command ID: %v", message.Cid)
	}
}

func (s *UnixServer) deleteTask(buffer []byte, conn net.Conn) {
	s.logger.Debug("Incoming delete task request")
	payload := models.Protocol{}
	if err := json.Unmarshal(buffer, &payload); err != nil {
		s.logger.Errorf("Failed to unmarshal delete task request: %v", err)
		s.sendErrorResponse(payload.Cid, "Invalid request format", conn)
		return
	}

	taskId, ok := payload.Payload.(float64)
	if !ok {
		s.logger.Errorf("Invalid payload type for delete task request: %v", payload.Payload)
		s.sendErrorResponse(payload.Cid, "Invalid payload type", conn)
		return
	}

	s.logger.Debugf("TaskId: %d", int(taskId))
	err := s.store.TaskDeleteByID(nil, int(taskId))
	if err != nil {
		s.logger.Errorf("Failed to delete task: %v", err)
		s.sendErrorResponse(payload.Cid, "Failed to delete task", conn)
		return
	}
	s.sendResponse(payload.Cid, "Task deleted", conn)
}

func (s *UnixServer) getTask(buffer []byte, conn net.Conn) {
	s.logger.Debug("Incoming get task request")
	payload := models.Protocol{}
	if err := json.Unmarshal(buffer, &payload); err != nil {
		s.logger.Errorf("Failed to unmarshal get task request: %v", err)
		s.sendErrorResponse(payload.Cid, "Invalid request format", conn)
		return
	}

	taskId, ok := payload.Payload.(float64)
	if !ok {
		s.logger.Errorf("Invalid payload type for get task request: %v", payload.Payload)
		s.sendErrorResponse(payload.Cid, "Invalid payload type", conn)
		return
	}

	s.logger.Debugf("TaskId: %d", int(taskId))
	task, err := s.store.TaskGetByID(nil, int(taskId))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Errorf("Task not found: %v", err)
			s.sendErrorResponse(payload.Cid, "Task not found", conn)
		} else {
			s.logger.Errorf("Failed to get task: %v", err)
			s.sendErrorResponse(payload.Cid, "Failed to get task", conn)
		}
		return
	}

	err = s.store.PomodoroDeleteByTaskID(nil, int(taskId))
	if err != nil {
		s.logger.Errorf("Failed to delete pomodoros for task: %v", err)
		s.sendErrorResponse(payload.Cid, "Failed to delete pomodoros", conn)
		return
	}
	task.Pomodoros = []*models.Pomodoro{}
	s.sendResponse(payload.Cid, task, conn)
}

func (s *UnixServer) createPomodoro(buffer []byte, conn net.Conn) {
	s.logger.Debug("Incoming create pomodoro request")
	payload := models.Protocol{Payload: &models.PomodoroWithID{}}
	if err := json.Unmarshal(buffer, &payload); err != nil {
		s.logger.Errorf("Failed to unmarshal create pomodoro request: %v", err)
		s.sendErrorResponse(payload.Cid, "Invalid request format", conn)
		return
	}

	pomodoro, ok := payload.Payload.(*models.PomodoroWithID)
	if !ok {
		s.logger.Errorf("Invalid payload type for create pomodoro request: %v", payload.Payload)
		s.sendErrorResponse(payload.Cid, "Invalid payload type", conn)
		return
	}

	if err := s.store.PomodoroSave(nil, pomodoro.TaskID, &pomodoro.Pomodoro); err != nil {
		s.logger.Errorf("Failed to save pomodoro: %v", err)
		s.sendErrorResponse(payload.Cid, "Failed to save pomodoro", conn)
		return
	}
	s.sendResponse(payload.Cid, "Pomodoro created", conn)
}

func (s *UnixServer) createTask(buffer []byte, conn net.Conn) {
	s.logger.Debug("Incoming create task request")
	payload := models.Protocol{Payload: &models.Task{}}
	if err := json.Unmarshal(buffer, &payload); err != nil {
		s.logger.Errorf("Failed to unmarshal create task request: %v", err)
		s.sendErrorResponse(payload.Cid, "Invalid request format", conn)
		return
	}

	task, ok := payload.Payload.(*models.Task)
	if !ok {
		s.logger.Errorf("Invalid payload type for create task request: %v", payload.Payload)
		s.sendErrorResponse(payload.Cid, "Invalid payload type", conn)
		return
	}

	taskId, err := s.store.TaskSave(nil, task)
	if err != nil {
		s.logger.Errorf("Failed to save task: %v", err)
		s.sendErrorResponse(payload.Cid, "Failed to save task", conn)
		return
	}
	s.sendResponse(payload.Cid, taskId, conn)
}

func (s *UnixServer) updateStatus(buffer []byte, conn net.Conn) {
	s.logger.Debug("Incoming update status request")
	payload := models.Protocol{Payload: &models.Status{}}
	if err := json.Unmarshal(buffer, &payload); err != nil {
		s.logger.Errorf("Failed to unmarshal update status request: %v", err)
		s.sendErrorResponse(payload.Cid, "Invalid request format", conn)
		return
	}

	status, ok := payload.Payload.(*models.Status)
	if !ok {
		s.logger.Errorf("Invalid payload type for update status request: %v", payload.Payload)
		s.sendErrorResponse(payload.Cid, "Invalid payload type", conn)
		return
	}

	s.status = *status
	s.sendResponse(payload.Cid, "Status updated", conn)
}

func (s *UnixServer) sendResponse(cid models.CmdID, payload interface{}, conn net.Conn) {
	raw, err := json.Marshal(&models.Protocol{Cid: cid, Payload: payload})
	if err != nil {
		s.logger.Errorf("Failed to marshal response: %v", err)
		return
	}
	s.logger.Debugf("Writing response: %s", string(raw))
	if _, err := conn.Write(raw); err != nil {
		s.logger.Errorf("Failed to write response: %v", err)
	}
}

func (s *UnixServer) sendErrorResponse(cid models.CmdID, message string, conn net.Conn) {
	raw, err := json.Marshal(&models.Protocol{Cid: cid, Error: &message})
	if err != nil {
		s.logger.Errorf("Failed to marshal response: %v", err)
		return
	}
	s.logger.Debugf("Writing response: %s", string(raw))
	if _, err := conn.Write(raw); err != nil {
		s.logger.Errorf("Failed to write response: %v", err)
	}
}

func valid(ok bool, logger *zap.SugaredLogger, expected interface{}, received interface{}) {
	if !ok {
		logger.Fatalf("Expected %v but got %v", expected, received)
	}
}
