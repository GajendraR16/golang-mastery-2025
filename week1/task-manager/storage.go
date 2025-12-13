package main

import (
	"encoding/json"
	"os"
)

const filename = "tasks.json"

func SaveTasks(tasks []*Task, filename string) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func LoadTasks(filename string) ([]*Task, error) {
	data, err := os.ReadFile(filename)
	if os.IsNotExist(err) {
		return []*Task{}, nil // First run
	}

	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return []*Task{}, nil
	}
	
	var tasks []*Task
	err = json.Unmarshal(data, &tasks)
	return tasks, err
}
