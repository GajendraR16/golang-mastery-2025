package main

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

type URLCheckResult struct {
	URL            string `json:"url"`
	Status         string `json:"status"`
	StatusCode     int    `json:"status_code"`
	ResponseTimeMs int64  `json:"response_time_ms"`
	Error          string `json:"error,omitempty"`
}

type URLCheckResponse struct {
	Results []*URLCheckResult `json:"results"`
}

type URLCheckRequest struct {
	URLs []string `json:"urls"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type PostgresStore struct {
	db *sql.DB
}

type App struct {
	Store *PostgresStore
}

func main() {

	r := mux.NewRouter()

	store, err := NewPostgresStore("postgres://postgres:postgres@localhost:5432/taskdb?sslmode=disable")

	if err != nil {
		slog.Error("Database Error", "error", err)
		os.Exit(1)
	}

	defer store.db.Close()

	app := &App{
		Store: store,
	}

	r.HandleFunc("/urls", app.CreateURLResponseHandler).Methods("POST")
	r.HandleFunc("/urls/results", app.URLResultsHandler).Methods("GET")

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

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	log.Println("Server gracefully stopped")

}

func (app *App) BatchProcessor(urls []string, r *http.Request) []*URLCheckResult {

	var (
		wg         sync.WaitGroup
		URLresults []*URLCheckResult
	)
	jobs := make(chan string)
	results := make(chan URLCheckResult)

	for range 5 {
		client := &http.Client{}
		wg.Go(func() {
			for url := range jobs {
				ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
				start := time.Now()

				req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
				if err != nil {
					cancel()
					results <- URLCheckResult{
						URL:            url,
						ResponseTimeMs: time.Since(start).Milliseconds(),
						Error:          err.Error(),
					}
					continue
				}

				resp, err := client.Do(req)
				cancel()
				if err != nil {
					results <- URLCheckResult{
						URL:            url,
						ResponseTimeMs: time.Since(start).Milliseconds(),
						Error:          err.Error(),
					}
					continue
				}

				resp.Body.Close()
				results <- URLCheckResult{
					URL:            url,
					Status:         resp.Status,
					StatusCode:     resp.StatusCode,
					ResponseTimeMs: time.Since(start).Milliseconds(),
				}
			}
		})
	}

	for _, url := range urls {
		jobs <- url
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		r := &result
		URLresults = append(URLresults, r)
	}

	return URLresults
}
