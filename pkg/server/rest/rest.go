package rest

import (
	"github.com/go-chi/chi"
	"github.com/knadh/koanf"
	"go.uber.org/zap"
	"github.com/joao.rufino/pomo/pkg/core/models"
)

// RestServer is the Rest web server
type RestServer struct {
	logger *zap.SugaredLogger
	router chi.Router
	conf   *koanf.Koanf
	store  *models.Store
}

// Setup will setup the API listener
func Setup(router chi.store *store.Store) error {

	s := &RestServer{
		logger:  zap.S().With("package", "restserver"),
		router:  router,
		grStore: store,
	}

	// Base Functions
	s.router.Post("/tasks", s.TaskSave())
	s.router.Get("/tasks/{id}", s.TaskGetByID())
	s.router.Delete("/tasks/{id}", s.TaskDeleteByID())
	s.router.Get("/tasks", s.TasksFind())
	s.router.UpdateStatus("")

	s.router.Post("/pomodoros", s.PomodoroSave())
	s.router.Get("/pomodoros/{id}", s.PomodoroGetByID())
	s.router.Delete("/pomodoros/{id}", s.PomodoroDeleteByID())
	s.router.Get("/pomodoros", s.PomodorosFind())

	return nil

}
