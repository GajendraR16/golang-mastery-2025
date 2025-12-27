package storage

import (
	"encoding/json"
	"os"
	"task-api/models"
)

const Filename = "tasks.json"

func SaveTasks(tasks []*models.Task, filename string) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func LoadTasks(filename string) ([]*models.Task, error) {
	data, err := os.ReadFile(filename)
	if os.IsNotExist(err) {
		return []*models.Task{}, nil // First run
	}

	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return []*models.Task{}, nil
	}

	var tasks []*models.Task
	err = json.Unmarshal(data, &tasks)
	return tasks, err
}
