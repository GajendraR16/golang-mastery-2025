package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type App struct {
	taskmanager *TaskManager
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type TaskData struct {
	Description string `json:"description"`
}

func jsonError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func jsonHandler(w http.ResponseWriter, data any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func (t *App) TaskHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	if query != "" {
		tasks, err := t.taskmanager.Search(r.Context(), query)
		if err != nil {
			jsonError(w, "Failed to search tasks", http.StatusInternalServerError)
			return
		}

		jsonHandler(w, tasks, http.StatusOK)
		return
	}

	tasks, err := t.taskmanager.List(r.Context())
	if err != nil {
		jsonError(w, "Failed to get tasks", http.StatusInternalServerError)
		return
	}

	jsonHandler(w, tasks, http.StatusOK)
}

func (a *App) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var description TaskData

	err := json.NewDecoder(r.Body).Decode(&description)
	if err != nil {
		jsonError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if description.Description == "" {
		jsonError(w, "Description cannot be empty", http.StatusBadRequest)
		return
	}

	if len(description.Description) < 3 {
		jsonError(w, "Description length must be 3 or more", http.StatusBadRequest)
		return
	}

	task, err := a.taskmanager.Add(r.Context(), description.Description)
	if err != nil {
		jsonError(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	jsonHandler(w, task, http.StatusCreated)
}

func (a *App) TaskHandlerID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		jsonError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	switch r.Method {

	case http.MethodGet:
		task, err := a.taskmanager.Get(r.Context(), id)
		if err != nil {
			jsonError(
				w,
				fmt.Sprintf("Task with id %d not found", id),
				http.StatusNotFound,
			)
			return
		}

		jsonHandler(w, task, http.StatusOK)

	case http.MethodPut:
		task, err := a.taskmanager.Complete(r.Context(), id)
		if err != nil {
			jsonError(
				w,
				fmt.Sprintf("Task with id %d not found", id),
				http.StatusNotFound,
			)
			return
		}

		jsonHandler(w, task, http.StatusOK)

	case http.MethodDelete:
		err := a.taskmanager.Delete(r.Context(), id)
		if err != nil {
			jsonError(
				w,
				fmt.Sprintf("Task with id %d not found", id),
				http.StatusNotFound,
			)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
