package main

import (
	"context"
	"fmt"
	"sync"
)

type Task struct {
	ID   int
	Data string
}

type Result struct {
	TaskID int
	Output string
}

func worker(ctx context.Context, id int, tasks <-chan Task, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case task, ok := <-tasks:
			if !ok {
				return
			}
			output := fmt.Sprintf("Worker %d processed %s", id, task.Data)
			results <- Result{TaskID: task.ID, Output: output}
		case <-ctx.Done():
			fmt.Printf("Worker %d: Received cancel signal, stopping...\n", id)
			return

		}
	}
}
