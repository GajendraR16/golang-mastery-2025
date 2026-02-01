package models

import (
	"time"
)

type TaskData struct {
	Description string `json:"description" validate:"required,min=3"`
}

// Task Model
type Task struct {
	ID          int        `json:"id"`
	Description string     `json:"description"`
	Completed   bool       `json:"complete"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
