package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joaorufino/pomo/pkg/client/graphql"
	"github.com/joaorufino/pomo/pkg/client/rest"
	"github.com/joaorufino/pomo/pkg/client/unix"
	"github.com/joaorufino/pomo/pkg/conf"
	"github.com/joaorufino/pomo/pkg/core"
	"go.uber.org/zap"
)

func NewClient(conf *conf.Config) (core.Client, error) {
	switch conf.Client.Type {
	case "unix":
		unixClient := &unix.UnixClient{}
		return unixClient.Init(conf)
	case "rest":
		restClient := &rest.RestClient{}
		return restClient.Init(conf)
	case "graphql":
		graphqlClient := &graphql.GraphQLClient{}
		return graphqlClient.Init(conf)
	}

	return nil, errors.New("error creating client")
}

type Task struct {
	ID         int      `json:"id"`
	Message    string   `json:"message"`
	NPomodoros int      `json:"n_pomodoros"`
	Tags       []string `json:"tags"`
}

type TaskResponse struct {
	Count   int    `json:"count"`
	Results []Task `json:"results"`
}

type PageData struct {
	Title  string
	ApiUrl string
	Tasks  []Task
}

func fetchTasks(apiUrl string) ([]Task, error) {
	resp, err := http.Get(apiUrl + "/tasks")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var taskResponse TaskResponse
	err = json.Unmarshal(body, &taskResponse)
	if err != nil {
		return nil, err
	}

	return taskResponse.Results, nil
}

func handler(apiUrl string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tasks, err := fetchTasks(apiUrl)
		if err != nil {
			http.Error(w, "Failed to fetch tasks: "+err.Error(), http.StatusInternalServerError)
			return
		}

		data := &PageData{
			Title:  "Pomodoro App",
			ApiUrl: apiUrl,
			Tasks:  tasks,
		}

		tmpl, err := template.ParseFiles("ui/templates/layout/base.html", "ui/templates/task/list.html")
		if err != nil {
			http.Error(w, "Failed to load template: "+err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.ExecuteTemplate(w, "base", data)
	}
}

func newTaskHandler(apiUrl string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := &PageData{
			ApiUrl: apiUrl,
		}

		tmpl, err := template.ParseFiles("ui/templates/task/form.html")
		if err != nil {
			http.Error(w, "Failed to load template: "+err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, data)
	}
}

func createTaskHandler(apiUrl string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		body, err := json.Marshal(task)
		if err != nil {
			http.Error(w, "Failed to create task", http.StatusInternalServerError)
			return
		}

		req, err := http.NewRequest("POST", apiUrl+"/tasks", bytes.NewBuffer(body))
		if err != nil {
			http.Error(w, "Failed to create request", http.StatusInternalServerError)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != http.StatusCreated {
			http.Error(w, "Failed to create task", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func StartWebServer(conf *conf.Config, log *zap.SugaredLogger) {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// CORS setup
	cors := cors.New(cors.Options{
		AllowedOrigins:   conf.CORS.AllowedOrigins,
		AllowedMethods:   conf.CORS.AllowedMethods,
		AllowedHeaders:   conf.CORS.AllowedHeaders,
		AllowCredentials: conf.CORS.AllowCredentials,
		MaxAge:           conf.CORS.MaxAge,
	})
	r.Use(cors.Handler)

	switch conf.Client.HostType {
	case "nextjs":
		// Serve Next.js static files
		staticPath := "ui/pomo-client/out"
		fs := http.FileServer(http.Dir(staticPath))
		r.Handle("/*", fs)
	case "gotpl":
		apiUrl := fmt.Sprintf("http://%s:%s/%s", conf.Server.RestHost, conf.Server.RestPort, conf.Server.RestPath)
		r.Get("/", handler(apiUrl))
		r.Get("/tasks/new", newTaskHandler(apiUrl))
		r.Post("/tasks", createTaskHandler(apiUrl))
	}

	http.ListenAndServe(":8081", r)
}
