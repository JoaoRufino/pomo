package rest

import (
	"net"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/joaorufino/pomo/pkg/conf"
	"github.com/joaorufino/pomo/pkg/core"
	"github.com/joaorufino/pomo/pkg/core/models"
	"github.com/joaorufino/pomo/pkg/store"
	"go.uber.org/zap"
)

// RestServer is the Rest web server
type RestServer struct {
	logger *zap.SugaredLogger
	router chi.Router
	conf   *conf.Config
	store  core.Store
	server *http.Server
	status models.Status
}

const (
	TASK_PATH        = "/tasks"
	TASK_ID_PATH     = TASK_PATH + "/{id}"
	POMODORO_PATH    = "/pomodoros"
	POMODORO_ID_PATH = POMODORO_PATH + "/{id}"
	STATUS_PATH      = "/status"
)

// Setup will setup the API listener
func (s *RestServer) Setup() error {
	// Base Functions
	s.router.Get(TASK_PATH, s.TasksFind())
	s.router.Post(TASK_PATH, s.TaskSave())
	s.router.Get(TASK_ID_PATH, s.TaskGetByID())
	s.router.Delete(TASK_ID_PATH, s.TaskDeleteByID())

	s.router.Post(POMODORO_ID_PATH, s.PomodoroSave())
	s.router.Get(POMODORO_ID_PATH, s.PomodoroGetByID())
	s.router.Delete(POMODORO_ID_PATH, s.PomodoroDeleteByID())

	s.router.Get(STATUS_PATH, s.StatusGet())
	s.router.Post(STATUS_PATH, s.StatusSave())

	return nil
}

// New will setup the API listener
func New(config *conf.Config) (core.Server, error) {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	// Log Requests - Use appropriate format depending on the encoding
	if config.Server.LogRequests {
		r.Use(loggerHTTPMiddlewareDefault(config.Server.LogRequestsBody, config.Server.LogDuration))
	}

	// CORS Config
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   config.CORS.AllowedOrigins,
		AllowedMethods:   config.CORS.AllowedMethods,
		AllowedHeaders:   config.CORS.AllowedHeaders,
		AllowCredentials: config.CORS.AllowCredentials,
		MaxAge:           config.CORS.MaxAge,
	}).Handler)

	store, _ := store.NewStore(config)

	server := &http.Server{
		Addr:    net.JoinHostPort(config.Server.RestHost, config.Server.RestPort),
		Handler: r,
	}

	s := &RestServer{
		conf:   config,
		logger: zap.S().With("package", "restServer"),
		router: r,
		store:  store,
		server: server,
	}

	// RestInterface
	if err := s.Setup(); err != nil {
		s.logger.Fatalf("Could not setup rest interface: %v", err)
	}
	return s, nil
}

// Start will listen for requests
func (s *RestServer) Start() {

	// Listen
	listener, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		s.logger.Fatalf("Could not listen on %s: %v", s.server.Addr, err)
	}

	go func() {
		if err = s.server.Serve(listener); err != nil {
			s.logger.Fatalw("API Listen error", "error", err, "address", s.server.Addr)
		}
	}()
	s.logger.Infow("API Listening", "address", s.server.Addr, "tls", s.conf.Server.TLS)
}

// Router returns the router
func (s *RestServer) Router() chi.Router {
	return s.router
}

func (s *RestServer) Stop() {
	s.server.Close()
	s.store.Close()
}
