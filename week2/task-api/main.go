package main

import (
	"log/slog"
	"net/http"
	"task-api/config"
	"task-api/handler"
	"task-api/middleware"
	"task-api/storage"

	"github.com/gorilla/mux"
)

func main() {
	//Database Connection/Abstraction

	cfg := config.Load()
	if cfg.DatabaseURL == "" {
		cfg.DatabaseURL = "postgres://postgres:postgres@localhost:5432/taskdb?sslmode=disable" // Local fallback
	}
	store, err := storage.NewPostgresStore(cfg.DatabaseURL)

	if err != nil {
		slog.Error("Database Error", "error", err)
	}

	app := &handler.App{
		Store: store,
	}

	router := mux.NewRouter()

	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.CorsMiddleware)
	// Specific Route First
	router.HandleFunc("/tasks", app.SearchHandler).Methods("GET").Queries("q", "{q}")

	//General Route
	router.HandleFunc("/tasks", app.TaskHandler).Methods("GET")

	router.HandleFunc("/tasks", app.CreateHandler).Methods("POST")
	router.HandleFunc("/tasks/{id:[0-9]+}", app.TaskCompleteHandler).Methods("PUT")
	router.HandleFunc("/tasks/{id:[0-9]+}", app.TaskHandlerById).Methods("GET")
	router.HandleFunc("/tasks/{id:[0-9]+}", app.DeleteHandler).Methods("DELETE")

	slog.Info("Starting server", "port", cfg.Port)
	http.ListenAndServe(":"+cfg.Port, router)
}
