package runner

import (
	"fmt"
	"io/ioutil"
	"path"
	"testing"
	"time"

	"github.com/joao.rufino/pomo/pkg/runner/test"
	"github.com/joao.rufino/pomo/pkg/server"
	"github.com/joao.rufino/pomo/pkg/server/models"
)

func TestTaskRunner(t *testing.T) {
	baseDir, _ := ioutil.TempDir("/tmp", "")
	store, err := server.NewStore(path.Join(baseDir, "pomo.db"))
	if err != nil {
		t.Error(err)
	}
	err = server.InitDB(store)
	if err != nil {
		t.Error(err)
	}
	runner, err := test.NewMockedTaskRunner(&models.Task{
		Duration:   time.Second * 2,
		NPomodoros: 2,
		Message:    fmt.Sprint("Test Task"),
	}, store, models.NoopNotifier{})
	if err != nil {
		t.Error(err)
	}

	runner.Start()

	runner.Toggle()
	runner.Toggle()

	runner.Toggle()
	runner.Toggle()
}
