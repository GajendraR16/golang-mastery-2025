package main

import (
	"fmt"
	"net/http"
	"task-api/handler"
	"task-api/middleware"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()

	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.CorsMiddleware)
	// Specific Route First
	router.HandleFunc("/tasks", handler.SearchHandler).Methods("GET").Queries("q", "{q}")

	//General Route
	router.HandleFunc("/tasks", handler.TaskHandler).Methods("GET")

	router.HandleFunc("/tasks", handler.CreateHandler).Methods("POST")
	router.HandleFunc("/tasks/{id:[0-9]+}", handler.TaskCompleteHandler).Methods("PUT")
	router.HandleFunc("/tasks/{id:[0-9]+}", handler.TaskHandlerById).Methods("GET")
	router.HandleFunc("/tasks/{id:[0-9]+}", handler.DeleteHandler).Methods("DELETE")

	fmt.Println("Starting server at 8080...")
	http.ListenAndServe(":8080", router)
}
