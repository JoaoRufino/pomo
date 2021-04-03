package runner

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"os"

	"github.com/joao.rufino/pomo/pkg/server/models"
	"github.com/knadh/koanf"
	"go.uber.org/zap"
)

// Client makes requests to a listening
// pomo server to check the status of
// any currently running task session.
type UnixClient struct {
	path   string
	config *koanf.Koanf
	conn   net.Conn
	logger *zap.SugaredLogger
}

func (c UnixClient) makeRequest(cid string, payload interface{}) []byte {
	var err error
	c.conn, err = net.Dial("unix", c.path)
	maybe(err, c.logger)
	defer c.conn.Close()
	rxCh := make(chan []byte)
	raw, _ := json.Marshal(&models.Protocol{Cid: cid, Payload: payload})
	c.conn.Write(raw)
	go c.read(rxCh)
	return <-rxCh
}

func (c UnixClient) CreateTask(task *models.Task) (int, error) {
	cid := "create"
	message := c.makeRequest(cid, task)
	response := models.Protocol{Payload: 0}
	json.Unmarshal(message, &response)
	if response.Cid != cid {
		return -1, errors.New("receive wrong message type")
	} else {
		taskID, ok := response.Payload.(int)
		if !ok {
			log.Printf("got data of type %T but wanted int", response.Payload)
			os.Exit(1)
		}
		return taskID, nil
	}

}
func (c UnixClient) CreatePomodoro(taskID int, pomodoro models.Pomodoro) error {
	cid := "pomodoro"
	message := c.makeRequest(cid, &pomodoro)
	response := models.Protocol{Payload: ""}
	json.Unmarshal(message, &response)
	if response.Cid != cid {
		return errors.New("receive wrong message type")
	} else {
		message, ok := response.Payload.(string)
		if !ok {
			log.Printf("got data of type %T but wanted int", response.Payload)
			os.Exit(1)
		}
		if len(message) != 0 {
			return errors.New(message)
		} else {
			return nil
		}
	}

}

func (c UnixClient) DeleteTaskByID(taskID int) error {
	cid := "delete"
	message := c.makeRequest(cid, &taskID)
	response := models.Protocol{Payload: ""}
	json.Unmarshal(message, &response)
	if response.Cid != cid {
		return errors.New("receive wrong message type")
	} else {
		message, ok := response.Payload.(string)
		if !ok {
			log.Printf("got data of type %T but wanted int", response.Payload)
			os.Exit(1)
		}
		if len(message) != 0 {
			return errors.New(message)
		} else {
			return nil
		}
	}
}

func (c UnixClient) GetServerStatus() (*models.Status, error) {
	cid := "status"
	message := c.makeRequest(cid, nil)
	response := models.Protocol{Payload: &models.Status{}}
	json.Unmarshal(message, &response)
	if response.Cid != cid {
		return nil, errors.New("received wrong message type")
	} else {
		status, ok := response.Payload.(models.Status)
		if !ok {
			log.Printf("got data of type %T but wanted int", response.Payload)
			os.Exit(1)
		}
		return &status, nil
	}
}

func (c UnixClient) GetTaskList() (models.List, error) {
	cid := "list"
	message := c.makeRequest(cid, nil)
	response := models.Protocol{Payload: &models.List{}}
	json.Unmarshal(message, &response)
	c.logger.Debugf("%s", response.Payload)
	if response.Cid != cid {
		return nil, errors.New("received wrong message type")
	} else {
		list, ok := response.Payload.(*models.List)
		if !ok {
			log.Printf("got data of type %T but wanted models.List", response.Payload)
			os.Exit(1)
		}
		return *list, nil
	}
}

func (c UnixClient) GetTask(taskID int) (models.Task, error) {
	cid := "task"
	message := c.makeRequest(cid, taskID)
	response := models.Protocol{Payload: &models.Task{}}
	json.Unmarshal(message, &response)
	if response.Cid != cid {
		return models.Task{}, errors.New("received wrong message type")
	} else {
		task, ok := response.Payload.(models.Task)
		if !ok {
			log.Printf("got data of type %T but wanted Task", response.Payload)
			os.Exit(1)
		}
		return task, nil
	}
}

func (c UnixClient) StartTask(taskID int) error {
	task, err := c.GetTask(taskID)
	maybe(err, c.logger)
	runner, err := NewTaskRunner(c, &task)
	maybe(err, c.logger)
	runner.Start()
	StartUI(runner)
	return nil
}
func (c UnixClient) UpdateStatus(status *models.Status) error {
	cid := "update"
	message := c.makeRequest(cid, status)
	response := models.Protocol{Payload: ""}
	json.Unmarshal(message, &response)
	if response.Cid != cid {
		return errors.New("received wrong message type")
	} else {
		message, ok := response.Payload.(string)
		if !ok {
			log.Printf("got data of type %T but wanted int", response.Payload)
			os.Exit(1)
		}
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

func (c UnixClient) Start(path string) (*UnixClient, error) {
	c.path = path
	c.logger = zap.S().With("package", "client")
	return &c, nil
}

func (c UnixClient) read(rxCh chan []byte) {
	buf := make([]byte, 4096)
	n, err := c.conn.Read(buf)
	c.conn.Close()
	maybe(err, c.logger)
	if n > 0 {
		rxCh <- buf[0:n]
	} else {
		rxCh <- nil
	}
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
