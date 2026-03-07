package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"task-api/config"
	"task-api/server"
	"task-api/storage"
	"time"
)

func main() {
	//Database Connection/Abstraction

	cfg := config.Load()
	store, err := storage.NewPostgresStore(cfg.DatabaseURL)

	if err != nil {
		slog.Error("Database Error", "error", err)
		os.Exit(1)
	}
	defer store.Close()

	//Setup Router
	router := server.SetupRouter(store)

	//Configure server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server
	go func() {
		slog.Info("Server starting", slog.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server error", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	slog.Info("Server ready - Press Ctrl+C to shutdown")

	// Wait for interrupt
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server is shutting down...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown:", slog.String("error", err.Error()))
	}
	slog.Info("Server stopped gracefully")

}
