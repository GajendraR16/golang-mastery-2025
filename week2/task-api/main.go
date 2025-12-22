package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()

	// Specific Route First
	router.HandleFunc("/tasks", SearchHandler).Methods("GET").Queries("q", "{q}")

	//General Route
	router.HandleFunc("/tasks", TaskHandler).Methods("GET")

	router.HandleFunc("/tasks", CreateHandler).Methods("POST")
	router.HandleFunc("/tasks/{id:[0-9]+}", TaskCompleteHandler).Methods("PUT")
	router.HandleFunc("/tasks/{id:[0-9]+}", TaskHandlerById).Methods("GET")
	router.HandleFunc("/tasks/{id:[0-9]+}", DeleteHandler).Methods("DELETE")

	fmt.Println("Starting server at 8080...")
	http.ListenAndServe(":8080", router)
}
