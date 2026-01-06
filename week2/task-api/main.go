package main

import (
	"fmt"
	"log"
	"net/http"
	"task-api/handler"
	"task-api/middleware"
	"task-api/storage"

	"github.com/gorilla/mux"
)

func main() {
	//Database Connection/Abstraction

	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=taskdb sslmode=disable"
	store, err := storage.NewPostgresStore(connStr)

	if err != nil {
		log.Fatal(err)
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

	fmt.Println("Starting server at 8080...")
	http.ListenAndServe(":8080", router)
}
