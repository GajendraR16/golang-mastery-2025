package main

import (
	"fmt"
	"strconv"
)

func printUsage() {
	fmt.Print(`Usage:
go run . add 'Buy Groceries'
go run . list
go run . complete 1
go run . delete 2
go run . search 'buy'
`)
}

func handleAdd(tm *TaskManager, description string) {
	tm.Add(description)
}

func handleList(tm *TaskManager) {
	fmt.Println("-----Task List-----")
	tasks := tm.List()
	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}
	for _, task := range tasks {
		fmt.Println(task)
	}
}

func handleComplete(tm *TaskManager, args string) error {
	id, err := strconv.Atoi(args)
	if err != nil {
		return fmt.Errorf("invalid ID: %v", err)
	}

	if err := tm.Complete(id); err != nil {
		return err
	}
	fmt.Println("Task completed.")
	return nil
}

func handleDelete(tm *TaskManager, args string) error {
	id, err := strconv.Atoi(args)
	if err != nil {
		return fmt.Errorf("invalid ID: %v", err)
	}

	if err := tm.Delete(id); err != nil {
		return err
	}
	fmt.Println("Task deleted.")
	return nil
}

func handleSearch(tm *TaskManager, args string) {
	results := tm.Search(args)
	for _, task := range results {
		fmt.Println(task)
	}
}
