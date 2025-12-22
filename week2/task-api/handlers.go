package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func jsonError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func jsonHandler(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to Encode Json", http.StatusInternalServerError)
		return
	}
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {

	tasks, _ := LoadTasks(filename)
	tm := NewTaskManager()
	tm.Tasks = tasks
	jsonHandler(w, http.StatusOK, tasks)

}

func CreateHandler(w http.ResponseWriter, r *http.Request) {

	var task TaskData
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

	tasks, _ := LoadTasks(filename)
	tm := NewTaskManager()
	tm.Tasks = tasks

	if len(tasks) > 0 {
		tm.NextID = tasks[len(tasks)-1].ID + 1
	}

	createdTask := tm.Add(task.Description)
	SaveTasks(tm.Tasks, filename)

	jsonHandler(w, http.StatusCreated, createdTask)
}

func TaskHandlerById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"]) // Regex in router ensures this is a number

	tasks, _ := LoadTasks(filename)
	for _, t := range tasks {
		if t.ID == id {
			jsonHandler(w, http.StatusOK, t)
			return
		}
	}
	jsonError(w, "Task Not Found", http.StatusNotFound)
}

func TaskCompleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"]) // Regex in router ensures this is a number

	tasks, _ := LoadTasks(filename)
	tm := NewTaskManager()

	tm.Tasks = tasks
	task := tm.Complete(id)

	if task == nil {
		jsonError(w, "Task Not Found", http.StatusNotFound)
		return
	}

	SaveTasks(tm.Tasks, filename)
	jsonHandler(w, http.StatusOK, task)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	tasks, _ := LoadTasks(filename)
	tm := NewTaskManager()
	tm.Tasks = tasks
	if tm.Delete(id) {
		SaveTasks(tm.Tasks, filename)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	jsonError(w, "Incorrect Id", http.StatusNotFound)
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	queryParam := r.URL.Query().Get("q")
	cleanQuery := strings.Trim(queryParam, " \"")
	cleanQuery = strings.ToLower(cleanQuery)

	tasks, _ := LoadTasks(filename)
	tm := NewTaskManager()
	tm.Tasks = tasks
	results := tm.Search(cleanQuery)
	jsonHandler(w, http.StatusOK, results)

}
