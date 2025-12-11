package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	tasks, _ := LoadTasks(filename)
	tm := NewTaskManager()
	tm.Tasks = tasks

	if len(tasks) > 0 {
		tm.NextID = tasks[len(tasks)-1].ID + 1
	}

	command := os.Args[1]
	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: go run . add \"description\"")
			return
		}
		handleAdd(tm, os.Args[2])

	case "complete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: go run . complete <id>")
			return
		}
		if err := handleComplete(tm, os.Args[2]); err != nil {
			fmt.Println("Error:", err)
		}

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: go run . delete <id>")
			return
		}
		if err := handleDelete(tm, os.Args[2]); err != nil {
			fmt.Println("Error:", err)
		}

	case "search":
		if len(os.Args) < 3 {
			fmt.Println("Usage: go run . search <word>")
			return
		}
		handleSearch(tm, os.Args[2])

	case "list":
		handleList(tm)
	}

	SaveTasks(tm.Tasks, filename)
}
