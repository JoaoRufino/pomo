package rest

import (
	"net"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/joao.rufino/pomo/pkg/core"
	"github.com/joao.rufino/pomo/pkg/core/models"
	"github.com/joao.rufino/pomo/pkg/store"
	"github.com/knadh/koanf"
	"go.uber.org/zap"
)

// RestServer is the Rest web server
type RestServer struct {
	logger *zap.SugaredLogger
	router chi.Router
	conf   *koanf.Koanf
	store  core.Store
	server *http.Server
	status models.Status
}

// Setup will setup the API listener
func (s *RestServer) Setup() error {

	// Base Functions
	s.router.Post("/tasks", s.TaskSave())
	s.router.Get("/tasks/{id}", s.TaskGetByID())
	s.router.Delete("/tasks/{id}", s.TaskDeleteByID())
	s.router.Get("/tasks", s.TasksFind())

	s.router.Post("/pomodoros/{id}", s.PomodoroSave())
	s.router.Get("/pomodoros/{id}", s.PomodoroGetByID())
	s.router.Delete("/pomodoros/{id}", s.PomodoroDeleteByID())

	s.router.Get("/status", s.StatusGet())
	s.router.Post("/status", s.StatusSave())

	return nil

}

// New will setup the API listener
func New(config *koanf.Koanf) (core.Server, error) {

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	// Log Requests - Use appropriate format depending on the encoding
	if config.Bool("server.log_requests") {
		r.Use(loggerHTTPMiddlewareDefault(config.Bool("server.log_requests_body"), config.Strings("server.log_disabled_http")))
	}

	// CORS Config
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   config.Strings("server.cors.allowed_origins"),
		AllowedMethods:   config.Strings("server.cors.allowed_methods"),
		AllowedHeaders:   config.Strings("server.cors.allowed_headers"),
		AllowCredentials: config.Bool("server.cors.allowed_credentials"),
		MaxAge:           config.Int("server.cors.max_age"),
	}).Handler)

	store, _ := store.NewStore(config)

	s := &RestServer{
		conf:   config,
		logger: zap.S().With("package", "restServer"),
		router: r,
		store:  store,
	}

	// RestInterface
	if err := s.Setup(); err != nil {
		s.logger.Fatalf("Could not setup rest interface: %v", err)
	}
	return s, nil

}

// ListenAndServe will listen for requests
func (s *RestServer) Start() {

	s.server = &http.Server{
		Addr:    net.JoinHostPort(s.conf.String("server.rest.host"), s.conf.String("server.rest.port")),
		Handler: s.router,
	}

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
	s.logger.Infow("API Listening", "address", s.server.Addr, "tls", s.conf.Bool("server.tls"))

	// Enable profiler
	if s.conf.Bool("server.profiler_enabled") && s.conf.String("server.profiler_path") != "" {
		zap.S().Debugw("Profiler enabled on API", "path", s.conf.String("server.profiler_path"))
		s.router.Mount(s.conf.String("server.profiler_path"), middleware.Profiler())
	}
}

// Router returns the router
func (s *RestServer) Router() chi.Router {
	return s.router
}

func (s *RestServer) Stop() {
	s.server.Close()
	s.store.Close()
}
