package models

import (
	"fmt"
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
