package main

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type Task struct {
	ID          int        `json:"id"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

func NewTask(description string) *Task {
	return &Task{
		Description: description,
	}
}

func (t *Task) Complete() {
	now := time.Now()
	t.Completed = true
	t.CompletedAt = &now
}

func (t *Task) String() string {
	status := "[]"
	if t.Completed {
		status = "[✓]"
	}
	return fmt.Sprintf("%d. %s %s", t.ID, status, t.Description)
}

type TaskNotFoundError struct {
	ID int
}

func (te TaskNotFoundError) Error() string {
	return fmt.Sprintf("task with id %d not found", te.ID)
}

type Storage interface {
	Save(ctx context.Context, task *Task) error
	Get(ctx context.Context, id int) (*Task, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) ([]*Task, error)
	Complete(ctx context.Context, id int) (*Task, error)
}

type TaskManager struct {
	store  Storage
	nextID int
}

func NewTaskManager(store Storage) *TaskManager {
	return &TaskManager{
		store:  store,
		nextID: 1,
	}
}

func (tm *TaskManager) Add(ctx context.Context, description string) (*Task, error) {
	task := NewTask(description)

	err := tm.store.Save(ctx, task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (tm *TaskManager) Complete(ctx context.Context, id int) (*Task, error) {
	return tm.store.Complete(ctx, id)
}

func (tm *TaskManager) Delete(ctx context.Context, id int) error {
	return tm.store.Delete(ctx, id)
}

func (tm *TaskManager) Search(ctx context.Context, query string) ([]*Task, error) {
	tasks, err := tm.store.List(ctx)
	if err != nil {
		return nil, err
	}

	query = strings.ToLower(query)
	result := []*Task{}

	for _, task := range tasks {
		if strings.Contains(strings.ToLower(task.Description), query) {
			result = append(result, task)
		}
	}

	return result, nil
}

func (tm *TaskManager) List(ctx context.Context) ([]*Task, error) {
	return tm.store.List(ctx)
}

func (tm *TaskManager) Get(ctx context.Context, id int) (*Task, error) {
	return tm.store.Get(ctx, id)
}
