package main

import (
	"fmt"
	"strings"
	"time"
)

type TaskData struct {
	Description string `json:"description" validate:"required,min=3"`
}

// Task Model with helper methods and Constructor
type Task struct {
	ID          int        `json:"id"`
	Description string     `json:"description"`
	Completed   bool       `json:"complete"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

func NewTask(id int, description string) *Task {
	return &Task{
		ID:          id,
		Description: description,
		Completed:   false,
		CreatedAt:   time.Now(),
	}
}

func (t *Task) Complete() {
	t.Completed = true
	now := time.Now()
	t.CompletedAt = &now
}

func (t *Task) String() string {
	status := "[ ]"
	if t.Completed {
		status = "[✓]"
	}
	return fmt.Sprintf("%d. %s %s", t.ID, status, t.Description)
}

// TaskManager Model with helper methods and Constructor
type TaskManager struct {
	Tasks  []*Task
	NextID int
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		Tasks:  []*Task{},
		NextID: 1,
	}
}

func (tm *TaskManager) Add(description string) *Task {
	task := NewTask(tm.NextID, description)
	tm.Tasks = append(tm.Tasks, task)
	return task
}

func (tm *TaskManager) List() []*Task {
	return tm.Tasks
}

func (tm *TaskManager) Complete(id int) *Task {
	for _, task := range tm.Tasks {
		if task.ID == id {
			task.Complete()
			return task
		}
	}
	return nil
}

func (tm *TaskManager) Delete(id int) bool {
	for idx, task := range tm.Tasks {
		if task.ID == id {
			tm.Tasks = append(tm.Tasks[:idx], tm.Tasks[idx+1:]...)
			return true
		}
	}
	return false
}

func (tm *TaskManager) Search(query string) []*Task {
	results := []*Task{}

	for _, task := range tm.Tasks {
		if strings.Contains(strings.ToLower(task.Description), query) {
			results = append(results, task)
		}
	}
	return results
}

type ErrorResponse struct {
	Error string `json:"error"`
}
