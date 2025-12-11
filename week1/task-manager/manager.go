package main

import (
	"strings"
)

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
	tm.NextID++
	return task
}

func (tm *TaskManager) List() []*Task {
	return tm.Tasks
}

func (tm *TaskManager) Complete(id int) error {
	for _, task := range tm.Tasks {
		if task.ID == id {
			task.Complete()
			return nil
		}
	}
	return TaskNotFoundError{ID: id}
}

func (tm *TaskManager) Delete(id int) error {
	for idx, task := range tm.Tasks {
		if task.ID == id {
			tm.Tasks = append(tm.Tasks[:idx], tm.Tasks[idx+1:]...)
			return nil
		}
	}
	return TaskNotFoundError{ID: id}
}

func (tm *TaskManager) Search(query string) []*Task {
	results := []*Task{}
	for _, task := range tm.Tasks {
		if strings.Contains(strings.ToLower(task.Description), strings.ToLower(query)) {
			results = append(results, task)
		}
	}
	return results
}
