package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"task-api/models"
	"task-api/storage"

	"github.com/gorilla/mux"
)

type App struct {
	Store *storage.PostgresStore
}

func jsonError(w http.ResponseWriter, message string, code int, err error) {
	if err != nil && code == 500 {
		slog.Error(message, "error", err.Error())
	} else {
		slog.Warn(message, "status", code)
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(models.ErrorResponse{Error: message})
}

func jsonHandler(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to Encode Json", http.StatusInternalServerError)
		return
	}
}

func (app *App) TaskHandler(w http.ResponseWriter, r *http.Request) {

	tasks, err := app.Store.GetAllTasks(r.Context())
	if err != nil {
		if err == context.DeadlineExceeded {
			jsonError(w, "Database query timeout", http.StatusRequestTimeout, err)
			return
		}
		jsonError(w, "internal error", 500, err)
		return
	}

	jsonHandler(w, 200, tasks)

}

func (app *App) CreateHandler(w http.ResponseWriter, r *http.Request) {

	var task models.TaskData
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		jsonError(w, "Invalid Json", http.StatusBadRequest, err)
		return
	}

	task.Description = strings.TrimSpace(task.Description)

	if task.Description == "" {
		jsonError(w, "Description cannot be empty", http.StatusBadRequest, nil)
		return
	}

	if len(task.Description) < 3 {
		jsonError(w, "Description is too short (min 3 chars)", http.StatusBadRequest, nil)
		return
	}

	createdTask, err := app.Store.CreateTask(r.Context(), task.Description)
	if err != nil {
		if err == context.DeadlineExceeded {
			jsonError(w, "Database query timeout", http.StatusRequestTimeout, err)
			return
		}
		jsonError(w, "Internal Server Error", http.StatusInternalServerError, err)
		return
	}

	jsonHandler(w, http.StatusCreated, createdTask)
}

func (app *App) TaskHandlerById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"]) // Regex in router ensures this is a number

	task, err := app.Store.GetTaskById(r.Context(), id)
	if err != nil {
		if err == context.DeadlineExceeded {
			jsonError(w, "Database query timeout", http.StatusRequestTimeout, err)
			return
		}
		jsonError(w, err.Error(), http.StatusNotFound, err)
		return
	}

	jsonHandler(w, http.StatusFound, task)
}

func (app *App) TaskCompleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"]) // Regex in router ensures this is a number

	task, err := app.Store.CompletedTaskById(r.Context(), id)

	if err != nil {
		if err == context.DeadlineExceeded {
			jsonError(w, "Database query timeout", http.StatusRequestTimeout, err)
			return
		}
		jsonError(w, "Task Not Found", http.StatusNotFound, err)
		return
	}

	jsonHandler(w, http.StatusOK, task)
}

func (app *App) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	err := app.Store.DeleteTaskById(r.Context(), id)
	if err != nil {
		if err == context.DeadlineExceeded {
			jsonError(w, "Database query timeout", http.StatusRequestTimeout, err)
			return
		}
		jsonError(w, "Incorrect Id", http.StatusNotFound, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *App) SearchHandler(w http.ResponseWriter, r *http.Request) {
	queryParam := r.URL.Query().Get("q")
	cleanQuery := strings.Trim(queryParam, " \"")

	tasks, err := app.Store.SearchTasks(r.Context(), cleanQuery)
	if err != nil {
		if err == context.DeadlineExceeded {
			jsonError(w, "Database query timeout", http.StatusRequestTimeout, err)
			return
		}
		jsonError(w, "Failed to fetch tasks", http.StatusInternalServerError, err)
		return
	}
	jsonHandler(w, http.StatusOK, tasks)

}
