package handler

import (
	"encoding/json"
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

func jsonError(w http.ResponseWriter, message string, code int) {
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

	tasks, err := app.Store.GetAllTasks()
	if err != nil {
		jsonError(w, err.Error(), 500)
		return
	}

	jsonHandler(w, 200, tasks)

}

func (app *App) CreateHandler(w http.ResponseWriter, r *http.Request) {

	var task models.TaskData
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		jsonError(w, "Invalid Json", http.StatusBadRequest)
		return
	}

	task.Description = strings.TrimSpace(task.Description)

	if task.Description == "" {
		jsonError(w, "Description cannot be empty", http.StatusBadRequest)
		return
	}

	if len(task.Description) < 3 {
		jsonError(w, "Description is too short (min 3 chars)", http.StatusBadRequest)
		return
	}

	createdTask, err := app.Store.CreateTask(task.Description)
	if err != nil {
		jsonError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	jsonHandler(w, http.StatusCreated, createdTask)
}

func (app *App) TaskHandlerById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"]) // Regex in router ensures this is a number

	task, err := app.Store.GetTaskById(id)
	if err != nil {
		jsonError(w, err.Error(), http.StatusNotFound)
		return
	}

	jsonHandler(w, http.StatusFound, task)
}

func (app *App) TaskCompleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"]) // Regex in router ensures this is a number

	task, err := app.Store.CompletedTaskById(id)

	if err != nil {
		jsonError(w, "Task Not Found", http.StatusNotFound)
		return
	}

	jsonHandler(w, http.StatusOK, task)
}

func (app *App) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	err := app.Store.DeleteTaskById(id)
	if err != nil {
		jsonError(w, "Incorrect Id", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *App) SearchHandler(w http.ResponseWriter, r *http.Request) {
	queryParam := r.URL.Query().Get("q")
	cleanQuery := strings.Trim(queryParam, " \"")
	cleanQuery = strings.ToLower(cleanQuery)

	tasks, _ := app.Store.GetAllTasks()
	tm := models.NewTaskManager()
	tm.Tasks = tasks
	results := tm.Search(cleanQuery)
	jsonHandler(w, http.StatusOK, results)

}
