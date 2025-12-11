package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestSaveLoadTask(t *testing.T) {

	tempDir, err := os.MkdirTemp("", "task-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testFile := tempDir + "/test_task.json"

	originalTasks := []*Task{
		NewTask(1, "Test save"),
		NewTask(2, "Test Load"),
	}

	if err := SaveTasks(originalTasks, testFile); err != nil {
		t.Fatalf("Save Tasks failed: %v", err)
	}

	loadedTasks, err := LoadTasks(testFile)
	if err != nil {
		t.Fatalf("LoadTasks failed: %v", err)
	}

	if len(loadedTasks) != len(originalTasks) {
		t.Errorf("Expected %d tasks, got %d", len(originalTasks), len(loadedTasks))
	}
}

func TestLoadEmptyFile(t *testing.T) {

	const nonExistentFile = "non_existent_test_data.json"

	loadedTasks, err := LoadTasks(nonExistentFile)
	if err != nil {
		t.Fatalf("LoadTasks failed: %v", err)
	}

	if len(loadedTasks) != 0 {
		t.Errorf("Expected LoadTasks to return empty, got %d", len(loadedTasks))
	}
}

func TestNextIDCalculation(t *testing.T) {

	tempDir, _ := os.MkdirTemp("", "task-test")
	defer os.RemoveAll(tempDir)

	testFile := tempDir + "/test_task.json"

	originalTasks := []*Task{
		NewTask(1, "Test save"),
		NewTask(2, "Test Load"),
	}

	if err := SaveTasks(originalTasks, testFile); err != nil {
		t.Fatalf("Save Tasks failed: %v", err)
	}

	loadedTasks, _ := LoadTasks(testFile)
	tm := NewTaskManager()
	tm.Tasks = loadedTasks

	if len(loadedTasks) > 0 {
		tm.NextID = loadedTasks[len(loadedTasks)-1].ID + 1
	}

	if tm.NextID != 3 {
		t.Errorf("NextID expected to be 3, got %d", tm.NextID)
	}

}

func TestTaskManagerDelete(t *testing.T) {
	tm := NewTaskManager()

	tm.NextID = 3
	tm.Tasks = []*Task{
		NewTask(1, "Task A"),
		NewTask(2, "Task B"),
	}

	if err := tm.Delete(1); err != nil {
		t.Fatalf("Delete failed for existing task: %v", err)
	}
	if len(tm.Tasks) != 1 || tm.Tasks[0].ID != 2 {
		t.Errorf("Delete failed: expected 1 remaining task with ID 2, got %v", tm.Tasks)
	}

	if err := tm.Delete(99); err == nil {
		t.Error("Delete succeeded unexpectedly for non-existent ID 99")
	}
}

func TestHandleList(t *testing.T) {
	tm := NewTaskManager()
	tm.Add("First task")

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	defer func() {
		os.Stdout = oldStdout
	}()

	handleList(tm)
	w.Close()

	var capturedOutput bytes.Buffer
	io.Copy(&capturedOutput, r)

	expected := "-----Task List-----\n1. [ ] First task\n"

	if capturedOutput.String() != expected {
		t.Errorf("handleList output mismatch.\nExpected:\n%q\nGot:\n%q",
			expected, capturedOutput.String())
	}
}

func TestCompleteStatus(t *testing.T) {
	tm := NewTaskManager()

	tm.Add("Task to complete")
	err := tm.Complete(1)

	if err != nil {
		t.Fatalf("Complete failed unexpectedly for ID 1: %v", err)
	}

	completedTask := tm.Tasks[0]

	if !completedTask.Completed {
		t.Error("Task was not marked as Completed (status is false)")
	}
	if completedTask.CompletedAt == nil {
		t.Error("CompletedAt timestamp was not recorded (is nil)")
	}
}

func TestCompleteNotFound(t *testing.T) {
	tm := NewTaskManager()
	tm.Add("Task A")

	const nonExistentID = 99

	err := tm.Complete(nonExistentID)

	if err == nil {
		t.Fatal("Expected an error for non-existent ID, got nil")
	}

	if _, ok := err.(TaskNotFoundError); !ok {
		t.Errorf("Expected TaskNotFoundError, got %T", err)
	}

	if tm.Tasks[0].Completed {
		t.Error("A random task was accidentally marked completed.")
	}
}

func TestSearchCaseInsensitive(t *testing.T) {
	tm := NewTaskManager()
	tm.Add("Buy Groceries and milk")
	tm.Add("Write blog post")

	results1 := tm.Search("groceries")
	if len(results1) != 1 || results1[0].ID != 1 {
		t.Errorf("Search failed for lowercase query. Expected ID 1, got %v", results1)
	}

	results2 := tm.Search("WRITE")
	if len(results2) != 1 || results2[0].ID != 2 {
		t.Errorf("Search failed for uppercase query. Expected ID 2, got %v", results2)
	}

	results3 := tm.Search("bUy")
	if len(results3) != 1 || results3[0].ID != 1 {
		t.Errorf("Search failed for mixed-case query. Expected ID 1, got %v", results3)
	}
}
