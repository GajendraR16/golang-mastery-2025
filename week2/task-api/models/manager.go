package models

import "strings"

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
