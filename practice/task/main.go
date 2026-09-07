package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	r := mux.NewRouter()
	r.Use(LoggingMiddleware)
	r.Use(TimeoutMiddleware)

	store, err := NewPostgresStore("postgres://postgres:postgres@localhost:5432/taskdb?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	app := &App{
		taskmanager: NewTaskManager(store),
	}

	r.HandleFunc("/tasks", app.TaskHandler).Methods("GET")
	r.HandleFunc("/tasks", app.CreateTaskHandler).Methods("POST")
	r.HandleFunc("/tasks/{id:[0-9]+}", app.TaskHandlerID).Methods("GET", "PUT", "DELETE")

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		log.Println("Starting Server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err) // ✅ Only fatal on real errors
		}
	}()

	// Implement graceful shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down the server...")

	// Set a timeout for shutdown (for example, 30 seconds).
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	defer store.db.Close()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
	log.Println("Server gracefully stopped")

}
