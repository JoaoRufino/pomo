package runner

import (
	"fmt"
	"io/ioutil"
	"path"
	"testing"
	"time"

	"github.com/joao.rufino/pomo/pkg/core/models"
	"github.com/joao.rufino/pomo/pkg/runner/test"
)

func TestTaskRunner(t *testing.T) {
	baseDir, _ := ioutil.TempDir("/tmp", "")
	store, err := store.NewStore(path.Join(baseDir, "pomo.db"))
	if err != nil {
		t.Error(err)
	}
	err = store.InitDB(store)
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
