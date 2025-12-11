package main

import (
	"fmt"
	"time"
)

type Task struct {
	ID          int        `json:"id"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
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

type TaskNotFoundError struct {
	ID int
}

func (e TaskNotFoundError) Error() string {
	return fmt.Sprintf("task with ID %d not found", e.ID)
}
