package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"task-api/models"
	"time"
)

type BatchRequest struct {
	Tasks []models.TaskData `json:"tasks"`
}

type BatchResponse struct {
	Created []models.Task `json:"created"`
	Failed  []BatchError  `json:"failed"`
	Stats   BatchStats    `json:"stats"`
}

type BatchError struct {
	Index int    `json:"index"`
	Error string `json:"error"`
}

type BatchStats struct {
	Total      int   `json:"total"`
	Success    int   `json:"success"`
	Failed     int   `json:"failed"`
	DurationMS int64 `json:"duration_ms"`
}

func (app *App) BatchCreateHandler(w http.ResponseWriter, r *http.Request) {
	var req BatchRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Invalid JSON", http.StatusBadRequest, err)
		return
	}

	if len(req.Tasks) == 0 || len(req.Tasks) > 100 {
		jsonError(w, "Task count must be 1 to 100", http.StatusBadRequest, nil)
		return
	}

	start := time.Now()

	//Process Concurrently
	response := app.processBatch(r.Context(), req.Tasks)
	response.Stats.DurationMS = time.Since(start).Milliseconds()

	jsonHandler(w, http.StatusOK, response)
}

func (app *App) processBatch(ctx context.Context, tasks []models.TaskData) BatchResponse {
	var (
		wg       sync.WaitGroup
		mu       sync.Mutex
		response = BatchResponse{
			Stats: BatchStats{
				Total: len(tasks),
			},
		}
	)

	//Semaphone to limit concurrent Operations
	sem := make(chan struct{}, 10)

	for i, task := range tasks {
		wg.Add(1)

		go func(index int, t models.TaskData) {
			defer wg.Done()

			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				return
			} //Aquire

			defer func() {
				<-sem
			}()

			// Validate and normalize description to match single create behavior.
			description := strings.TrimSpace(t.Description)
			if description == "" {
				mu.Lock()
				response.Failed = append(response.Failed, BatchError{
					Index: index,
					Error: "empty description",
				})
				response.Stats.Failed++
				mu.Unlock()
				return
			}
			if len(description) < 3 {
				mu.Lock()
				response.Failed = append(response.Failed, BatchError{
					Index: index,
					Error: "description too short (min 3 chars)",
				})
				response.Stats.Failed++
				mu.Unlock()
				return
			}

			//Create Task
			created, err := app.Store.CreateTask(ctx, description)

			mu.Lock()
			if err != nil {
				response.Failed = append(response.Failed, BatchError{
					Index: index,
					Error: err.Error(),
				})
				response.Stats.Failed++
			} else {
				response.Created = append(response.Created, *created)
				response.Stats.Success++
			}
			mu.Unlock()
		}(i, task)
	}
	wg.Wait()
	return response
}
