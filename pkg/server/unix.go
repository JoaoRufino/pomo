package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/user"
	"path"

	"github.com/joao.rufino/pomo/pkg/runner"
	"github.com/joao.rufino/pomo/pkg/server/models"
	"go.uber.org/zap"
)

// UnixServer listens on a Unix domain socket
// for Pomo status requests
type UnixServer struct {
	listener net.Listener
	running  bool
	runner   runner.TaskRunner
	store    *Store
	logger   *zap.SugaredLogger
}

const TASK_STATUS_REQUEST = "status"
const TASK_LIST_REQUEST = "list"
const TASK_STATUS_UPDATE = "update"

func (s UnixServer) listen() {
	s.logger.Info("Listening")
	for s.running {
		buf := make([]byte, 1024)
		conn, err := s.listener.Accept()
		maybe(err, s.logger)
		defer conn.Close()
		n, _ := conn.Read(buf)

		if n != 0 {
			s.logger.Debugf("Incoming request")
			message := models.Protocol{}
			err := json.Unmarshal(buf[0:n], &message)
			maybe(err, s.logger)
			s.logger.Debug(string(buf[0:n]))
			switch string(message.Cid) {

			case TASK_STATUS_REQUEST:
				s.logger.Debug("Incoming status request")
				if &s.runner == nil {
					raw, _ := json.Marshal(s.runner.Status())
					conn.Write(raw)
				} else {
					conn.Write([]byte{})
				}

			case TASK_LIST_REQUEST:
				s.logger.Debug("Incoming task list request")
				u, err := user.Current()
				db, err := NewStore(path.Join(u.HomeDir, ".pomo/pomo.db"))
				maybe(err, s.logger)
				defer db.Close()
				tasks := models.List{}
				err = db.With(func(tx *sql.Tx) error {
					tasks, err = db.ReadTasks(tx)
					maybe(err, s.logger)
					return nil
				})
				raw, err := json.Marshal(models.Protocol{Cid: "list", Payload: tasks})
				maybe(err, s.logger)
				conn.Write(raw)
				log.Printf("writing tasks:%s", string(raw))
				maybe(err, s.logger)
			}
		}
	}
}

func (s UnixServer) Start() {
	s.running = true
	s.listen()
}

func (s UnixServer) Stop() {
	s.running = false
	s.listener.Close()
}

func (s UnixServer) Init(socketPath string, runner models.Runner) (*UnixServer, error) {
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
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return nil, err
	}
	server := &UnixServer{
		listener: listener,
		logger:   zap.S().With("package", "server"),
	}

	return server, nil

}
