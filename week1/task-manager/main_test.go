package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// --- Task Manager Core Logic Tests (Manager Methods) ---

func TestManager_AddList(t *testing.T) {
	tm := NewTaskManager()

	// Test Add
	task1 := tm.Add("First")
	if task1.ID != 1 || len(tm.Tasks) != 1 || tm.NextID != 2 {
		t.Errorf("Add failed. Expected ID 1, Count 1, NextID 2. Got %d, %d, %d", task1.ID, len(tm.Tasks), tm.NextID)
	}

	// Test List
	tasks := tm.List()
	if len(tasks) != 1 || tasks[0].Description != "First" {
		t.Errorf("List failed. Expected 1 task, got %d", len(tasks))
	}
}

func TestManager_Complete(t *testing.T) {
	tm := NewTaskManager()
	tm.Add("Task to complete")

	// Case 1: Complete existing task (ID 1)
	err := tm.Complete(1)
	if err != nil {
		t.Fatalf("Complete failed unexpectedly for ID 1: %v", err)
	}

	completedTask := tm.Tasks[0]
	if !completedTask.Completed || completedTask.CompletedAt == nil {
		t.Error("Complete failed to mark status or set timestamp.")
	}

	// Case 2: Complete non-existent task (ID 99)
	err = tm.Complete(99)
	if err == nil {
		t.Fatal("Expected an error for non-existent ID 99, got nil")
	}
	if _, ok := err.(TaskNotFoundError); !ok {
		t.Errorf("Expected TaskNotFoundError, got %T", err)
	}
}

func TestManager_Delete(t *testing.T) {
	tm := NewTaskManager()
	tm.Add("Task A")
	tm.Add("Task B")
	initialCount := len(tm.Tasks)

	// Case 1: Delete existing task (ID 1)
	if err := tm.Delete(1); err != nil {
		t.Fatalf("Delete failed for existing task: %v", err)
	}
	if len(tm.Tasks) != initialCount-1 || tm.Tasks[0].ID != 2 {
		t.Error("Delete failed: task was not removed or wrong task remained.")
	}

	// Case 2: Delete non-existent task (ID 99)
	if err := tm.Delete(99); err == nil {
		t.Error("Delete succeeded unexpectedly for non-existent ID 99")
	}
}

func TestManager_Search(t *testing.T) {
	tm := NewTaskManager()
	tm.Add("Buy Groceries and milk")
	tm.Add("Write blog post")
	tm.Add("Sell old books")

	// Case 1: Case-Insensitive Search
	results := tm.Search("WRITE")
	if len(results) != 1 || results[0].ID != 2 {
		t.Errorf("Case-insensitive search failed. Expected ID 2, got %v", results)
	}

	// Case 2: Partial Match Search
	results = tm.Search("buy")
	if len(results) != 1 {
		t.Errorf("Partial search failed. Expected 1 result, got %d", len(results))
	}

	// Case 3: Empty Query Search (should return all)
	results = tm.Search("")
	if len(results) != 3 {
		t.Errorf("Empty query search failed. Expected all 3 tasks, got %d", len(results))
	}

	// Case 4: No Results
	results = tm.Search("xyz")
	if len(results) != 0 {
		t.Error("Search failed. Expected 0 results.")
	}
}

// --- Task Struct Methods Tests ---

func TestTask_InitializationAndState(t *testing.T) {
	task := NewTask(10, "Sample task")

	if task.ID != 10 || task.Description != "Sample task" || task.Completed {
		t.Error("New task initialized incorrectly.")
	}
	if task.CompletedAt != nil || task.CreatedAt.IsZero() {
		t.Error("Time fields initialized incorrectly.")
	}

	// Test Complete method
	task.Complete()
	if !task.Completed || task.CompletedAt == nil || task.CompletedAt.IsZero() {
		t.Error("Complete method failed to set status and timestamp.")
	}
}

func TestTask_String(t *testing.T) {
	// Case 1: Incomplete task string
	task := NewTask(5, "Test description")
	if str := task.String(); str != "5. [ ] Test description" {
		t.Errorf("Incomplete string format error. Got %q", str)
	}

	// Case 2: Complete task string
	task.Complete()
	if str := task.String(); !strings.Contains(str, "[✓]") {
		t.Errorf("Complete string format error. Got %q", str)
	}
}

func TestTaskNotFoundError_Error(t *testing.T) {
	err := TaskNotFoundError{ID: 99}
	expected := "task with ID 99 not found"

	if err.Error() != expected {
		t.Errorf("TaskNotFoundError message error. Expected %q, got %q", expected, err.Error())
	}
}

// --- Persistence Tests ---

func TestPersistence_SaveLoadNextID(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "task-test")
	defer os.RemoveAll(tempDir)
	testFile := tempDir + "/test_task.json"

	// 1. Arrange & Act: Save tasks (IDs 1, 2)
	originalTasks := []*Task{NewTask(1, "A"), NewTask(2, "B")}
	if err := SaveTasks(originalTasks, testFile); err != nil {
		t.Fatalf("SaveTasks failed: %v", err)
	}

	// 2. Act: Load tasks
	loadedTasks, err := LoadTasks(testFile)
	if err != nil {
		t.Fatalf("LoadTasks failed: %v", err)
	}

	// 3. Assert Load & NextID calculation
	if len(loadedTasks) != 2 {
		t.Errorf("Expected 2 loaded tasks, got %d", len(loadedTasks))
	}

	tm := NewTaskManager()
	tm.Tasks = loadedTasks
	if len(loadedTasks) > 0 {
		tm.NextID = loadedTasks[len(loadedTasks)-1].ID + 1
	}
	if tm.NextID != 3 {
		t.Errorf("NextID calculation failed. Expected 3, got %d", tm.NextID)
	}
}

func TestPersistence_LoadEmptyAndErrors(t *testing.T) {
	// Case 1: Load non-existent file (First run scenario)
	const nonExistentFile = "non_existent_test_data.json"
	loadedTasks, err := LoadTasks(nonExistentFile)
	if err != nil {
		t.Fatalf("LoadTasks failed on non-existent file: %v", err)
	}
	if len(loadedTasks) != 0 {
		t.Error("LoadTasks failed to return empty slice for non-existent file.")
	}

	// Case 2: Load corrupted JSON file
	tempDir, _ := os.MkdirTemp("", "task-test-corrupt")
	defer os.RemoveAll(tempDir)
	testFile := tempDir + "/corrupted.json"
	os.WriteFile(testFile, []byte("invalid json: {"), 0644)

	if _, err := LoadTasks(testFile); err == nil {
		t.Error("Expected error when loading corrupted JSON, got nil.")
	}

	// Case 3: SaveTasks error (Permission/Invalid path)
	tasks := []*Task{NewTask(1, "Test")}
	if err := SaveTasks(tasks, "/invalid/path/tasks.json"); err == nil {
		t.Error("Expected error when saving to invalid path, got nil.")
	}
}

// --- CLI Handler Tests (Check Output and Error Handling) ---

// Helper function to redirect stdout and capture output
func captureOutput(_ *testing.T, fn func()) string {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Ensure restoration and cleanup
	defer func() {
		os.Stdout = oldStdout
	}()

	fn() // Run the function that prints
	w.Close()

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestHandleAdd(t *testing.T) {
	tm := NewTaskManager()
	handleAdd(tm, "New task")

	if len(tm.Tasks) != 1 || tm.Tasks[0].Description != "New task" {
		t.Error("HandleAdd failed to create task correctly.")
	}
}

func TestHandleList(t *testing.T) {
	tm := NewTaskManager()
	tm.Add("List task 1")

	// Case 1: List with tasks
	output := captureOutput(t, func() { handleList(tm) })
	expected := "-----Task List-----\n1. [ ] List task 1\n"
	if output != expected {
		t.Errorf("handleList output mismatch.\nExpected:\n%q\nGot:\n%q", expected, output)
	}

	// Case 2: List with no tasks
	tmEmpty := NewTaskManager()
	output = captureOutput(t, func() { handleList(tmEmpty) })
	if !strings.Contains(output, "No tasks found") {
		t.Error("HandleList failed to display 'No tasks found'.")
	}
}

func TestHandleCompleteAndDeleteErrors(t *testing.T) {
	tm := NewTaskManager()
	tm.Add("Task to test errors")

	// Test non-numeric ID in complete
	err := handleComplete(tm, "abc")
	if err == nil {
		t.Error("handleComplete expected error on non-numeric ID")
	}

	// Test non-numeric ID in delete
	err = handleDelete(tm, "xyz")
	if err == nil {
		t.Error("handleDelete expected error on non-numeric ID")
	}
}

func TestHandleSearch(t *testing.T) {
	tm := NewTaskManager()
	tm.Add("Buy groceries")

	output := captureOutput(t, func() { handleSearch(tm, "buy") })

	if !strings.Contains(output, "Buy groceries") {
		t.Errorf("HandleSearch failed. Output: %q", output)
	}
}
