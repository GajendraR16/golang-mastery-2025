package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"task-api/models"
	"task-api/storage"
	"testing"

	"github.com/gorilla/mux"
)

var (
	testStore *storage.PostgresStore
	testApp   *App
)

// Note: init() runs before any Test functions
func init() {
	testStore = setupTestDB()
	testApp = &App{Store: testStore}

}

func setupTestDB() *storage.PostgresStore {

	connStr := "postgres://postgres:postgres@localhost:5432/taskdb_test?sslmode=disable"
	store, err := storage.NewPostgresStore(connStr)

	if err != nil {
		log.Fatalf("Failed to connect to test DB: %v", err)
	}

	return store
}

func clearTable(store *storage.PostgresStore) {

	_, err := store.GetDB().Exec("TRUNCATE TABLE tasks RESTART IDENTITY CASCADE")
	if err != nil {
		log.Fatalf("Could not clear table: %v", err)
	}
}

func TestGetTasks(t *testing.T) {

	clearTable(testStore)

	// Create test data
	testStore.CreateTask("Test task")

	// Make request
	req := httptest.NewRequest("GET", "/tasks", nil)
	w := httptest.NewRecorder()
	testApp.TaskHandler(w, req)

	// Assert response
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
	var tasks []models.Task
	json.NewDecoder(w.Body).Decode(&tasks)
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}
}

func TestCreateTasks(t *testing.T) {

	clearTable(testStore)

    // 1. Create a buffer with JSON data
    taskData := map[string]string{"description": "New Task"}
    body, _ := json.Marshal(taskData)

    // 2. Pass the body to the request
    req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    
    w := httptest.NewRecorder()
    testApp.CreateHandler(w, req)

    // 3. Assertions
    if w.Code != http.StatusCreated && w.Code != http.StatusOK {
        t.Errorf("Expected 201 or 200, got %d", w.Code)
    }

    var task models.Task
    if err := json.NewDecoder(w.Body).Decode(&task); err != nil {
        t.Errorf("Failed to decode JSON: %v", err)
    }

    if task.ID != 1 {
        t.Errorf("Expected ID 1, got %d", task.ID)
    }

}

func TestCompleteHandler(t *testing.T) {

	clearTable(testStore)

	// Create test data
	testStore.CreateTask("Test Data Complete")

	// Make request
	req := httptest.NewRequest("PUT", "/tasks/1", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": "1",
	})
	w := httptest.NewRecorder()
	testApp.TaskCompleteHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var task models.Task
	if err := json.NewDecoder(w.Body).Decode(&task); err != nil {
		t.Errorf("Failed to decode JSON: %v", err)
	}

	// 1. Ensure it isn't nil
	if task.CompletedAt == nil {
		t.Error("Expected CompletedAt to be set, but it was nil")
	} else {
		// 2. Ensure it isn't the "Zero Time" (0001-01-01)
		if task.CompletedAt.IsZero() {
			t.Error("Expected a real timestamp, but got Zero time")
		}
	}
}

func TestDeleteHandler(t *testing.T) {

	clearTable(testStore)

	// Create test data
	testStore.CreateTask("Test Data Delete")

	// Make request
	req := httptest.NewRequest("DELETE", "/tasks/1", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": "1",
	})
	w := httptest.NewRecorder()
	testApp.DeleteHandler(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected 204, got %d", w.Code)
	}

	// Make request
	req = httptest.NewRequest("GET", "/tasks/1", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": "1",
	})
	w = httptest.NewRecorder()
	testApp.TaskHandlerById(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected 404, got %d", w.Code)
	}

}

func TestSearchHandler(t *testing.T) {

	clearTable(testStore)

	// Create test data
	testStore.CreateTask("Test Data 1")
	testStore.CreateTask("Test Data 2")

	// Make request
	req := httptest.NewRequest("GET", "/tasks/?q=Test", nil)
	w := httptest.NewRecorder()
	testApp.SearchHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	// Make request
	var tasks []models.Task
	if err := json.NewDecoder(w.Body).Decode(&tasks); err != nil {
		t.Errorf("Failed to decode JSON: %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("Expect 2 instead got %d", len(tasks))
	}

}
