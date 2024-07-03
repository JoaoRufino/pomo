package unix

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/joao.rufino/pomo/pkg/core"
	"github.com/joao.rufino/pomo/pkg/core/models"
	serverStore "github.com/joao.rufino/pomo/pkg/store"
	"github.com/knadh/koanf"
	"go.uber.org/zap"
)

func maybe(err error, logger *zap.SugaredLogger) {
	if err != nil {
		logger.Fatalf("Error:%s\n", err)
		os.Exit(1)
	}
}

// UnixServer listens on a Unix domain socket
// for Pomo status requests
type UnixServer struct {
	listener net.Listener
	running  bool
	store    core.Store
	logger   *zap.SugaredLogger
	status   models.Status
}

// listens for client requests following models.Protocol
// cid drives the response
func (s UnixServer) listen() {
	s.logger.Info("Listening")
	for s.running {
		buf := make([]byte, 1024)
		conn, err := s.listener.Accept()
		maybe(err, s.logger)
		defer conn.Close()
		n, _ := conn.Read(buf)

		//data was received
		if n != 0 {
			s.logger.Debugf("Incoming request")
			message := models.Protocol{}

			//Unmarshal to get CommandID
			err := json.Unmarshal(buf[0:n], &message)
			maybe(err, s.logger)

			switch message.Cid {
			//get server status
			case models.Cmd_GetServerStatus:
				s.logger.Debug("Incoming status request")
				_ = s.sendResponse(message.Cid, s.status, conn)

			//get all tasks
			case models.Cmd_GetList:
				s.logger.Debug("Incoming task list request")
				tasks, err := s.store.GetAllTasks(nil)
				maybe(err, s.logger)
				_ = s.sendResponse(message.Cid, tasks, conn)

			//create a task return
			case models.Cmd_CreateTask:
				s.createTask(buf[0:n], conn)
			//delete a task by id
			case models.Cmd_DeleteTask:
				s.deleteTask(buf[0:n], conn)
			//get a task by ID
			case models.Cmd_GetTask:
				s.getTask(buf[0:n], conn)

			//get a pomodoro by taskID
			case models.Cmd_CreatePomodoro:
				s.createPomodoro(buf[0:n], conn)

			//update server status
			case models.Cmd_UpdateStatus:
				s.updateStatus(buf[0:n], conn)
			}
		}
	}
}
func (s UnixServer) deleteTask(buffer []byte, conn net.Conn) {
	s.logger.Debug("Incoming delete task request")
	payload := models.Protocol{Payload: 0}
	json.Unmarshal(buffer, &payload)
	taskId, ok := payload.Payload.(float64)
	valid(ok, s.logger, payload.Payload, taskId)

	s.logger.Debugf("TaskId:%d", taskId)
	err := s.store.TaskDeleteByID(nil, int(taskId))
	maybe(err, s.logger)
	_ = s.sendResponse(payload.Cid, "", conn)
}
func (s UnixServer) getTask(buffer []byte, conn net.Conn) {
	s.logger.Debug("Incoming get task request")
	payload := models.Protocol{Payload: 0}
	json.Unmarshal(buffer, &payload)
	taskId, ok := payload.Payload.(float64)
	valid(ok, s.logger, payload.Payload, taskId)

	s.logger.Debugf("TaskId:%d", taskId)
	var task *models.Task
	task, err := s.store.TaskGetByID(nil, int(taskId))
	maybe(err, s.logger)
	err = s.store.PomodoroDeleteByTaskID(nil, int(taskId))
	maybe(err, s.logger)
	task.Pomodoros = []*models.Pomodoro{}
	_ = s.sendResponse(payload.Cid, task, conn)
}
func (s UnixServer) createPomodoro(buffer []byte, conn net.Conn) {
	s.logger.Debug("Incoming create pomodoro request")
	payload := models.Protocol{Payload: &models.PomodoroWithID{}}
	json.Unmarshal(buffer, &payload)
	pomodoro, ok := payload.Payload.(*models.PomodoroWithID)
	valid(ok, s.logger, payload.Payload, pomodoro)

	err := s.store.PomodoroSave(nil, pomodoro.TaskID, &pomodoro.Pomodoro)
	maybe(err, s.logger)
	_ = s.sendResponse(payload.Cid, err.Error(), conn)
}
func (s UnixServer) createTask(buffer []byte, conn net.Conn) {
	s.logger.Debug("Incoming create task request")
	var taskId int
	payload := models.Protocol{Payload: &models.Task{}}
	json.Unmarshal(buffer, &payload)
	task, ok := payload.Payload.(*models.Task)
	valid(ok, s.logger, payload.Payload, task)

	taskId, err := s.store.TaskSave(nil, task)
	maybe(err, s.logger)
	_ = s.sendResponse(payload.Cid, taskId, conn)
}
func (s UnixServer) updateStatus(buffer []byte, conn net.Conn) {
	s.logger.Debug("Incoming update status request")

	payload := models.Protocol{Payload: &models.Status{}}
	json.Unmarshal(buffer, &payload)

	status, ok := payload.Payload.(*models.Status)
	valid(ok, s.logger, payload.Payload, status)

	s.status = *status
	_ = s.sendResponse(payload.Cid, "", conn)
}

// makeRequest sends a message to the server
//using the protocol structure
func (s UnixServer) sendResponse(cid models.CmdID, payload interface{}, conn net.Conn) error {
	raw, err := json.Marshal(&models.Protocol{Cid: cid, Payload: payload})
	maybe(err, s.logger)
	s.logger.Debugf("writing tasks:%s", string(raw))
	conn.Write(raw)
	return nil
}

//Starts the server
func (s UnixServer) Start() {
	s.running = true
	s.listen()
}

//Stops the server
func (s UnixServer) Stop() {
	s.running = false
	s.listener.Close()
	s.store.Close()
}

//Initializes the server structure
func (s UnixServer) Init(k *koanf.Koanf, runner models.Runner) (*UnixServer, error) {
	socketPath := k.String("server.unix.socket")
	if _, err := os.Stat(socketPath); err == nil {
		_, err := net.Dial("unix", socketPath)
		//if error then sock file was saved after crash
		if err != nil {
			os.Remove(socketPath)
		} else {
			// another instance of pomo is running
			return nil, fmt.Errorf("socket %s is already in use", socketPath)
		}
	}
	store, err := serverStore.NewStore(k)
	maybe(err, s.logger)

	//open the socket
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return nil, err
	}
	server := &UnixServer{
		listener: listener,
		logger:   zap.S().With("package", "server"),
		store:    store,
		status:   models.Status{},
	}

	return server, nil

}

func valid(ok bool, logger *zap.SugaredLogger, expected interface{}, received interface{}) {
	if !ok {
		logger.Fatalf(models.ErrWrongDataType, expected, received)
	}
}
