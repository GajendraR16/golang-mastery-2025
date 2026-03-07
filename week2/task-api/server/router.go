package server

import (
	"task-api/handler"
	"task-api/middleware"
	"task-api/storage"
	"time"

	"github.com/gorilla/mux"
)

func SetupRouter(store *storage.PostgresStore) *mux.Router {
	app := &handler.App{
		Store: store,
	}

	router := mux.NewRouter()

	router.Use(middleware.RequestIDMiddleware)
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.CorsMiddleware)
	router.Use(middleware.TimeoutMiddleware(time.Second * 2))
	//Health Route
	router.HandleFunc("/health", app.HealthHandler).Methods("GET")

	// Specific Route First
	router.HandleFunc("/tasks", app.SearchHandler).Methods("GET").Queries("q", "{q}")

	//General Route
	router.HandleFunc("/tasks", app.TaskHandler).Methods("GET")

	//Batch Route
	router.HandleFunc("/tasks/batch", app.BatchCreateHandler).Methods("POST")

	router.HandleFunc("/tasks", app.CreateHandler).Methods("POST")
	router.HandleFunc("/tasks/{id:[0-9]+}", app.TaskCompleteHandler).Methods("PUT")
	router.HandleFunc("/tasks/{id:[0-9]+}", app.TaskHandlerById).Methods("GET")
	router.HandleFunc("/tasks/{id:[0-9]+}", app.DeleteHandler).Methods("DELETE")

	return router
}
