package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {

	//1. Basic GoRoutines
	// var wg sync.WaitGroup
	// wg.Add(5)

	// for range 5 {
	// 	go numbers(&wg)
	// }
	// wg.Wait()

	//2. Channel GoRoutine example
	// pipeline()

	//3. Producer and Consumer simple worker example

	var wg sync.WaitGroup

	tasks := make(chan Task, 100)
	results := make(chan Result, 100)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	for i := range 5 {
		wg.Add(1)
		go worker(ctx, i, tasks, results, &wg)
	}

	go func() {
		for i := 1; i <= 10; i++ {
			data := fmt.Sprintf("Task %d Data", i)
			tasks <- Task{ID: i, Data: data}
			time.Sleep(1 * time.Second)
		}
		close(tasks)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		fmt.Printf("Results ID %d| %s\n", res.TaskID, res.Output)
	}

	//4.Select Statment main Code

	// result := make(chan string, 1)
	// go SlowApi(result)
	// select {
	// case msg := <-result:
	// 	fmt.Printf("API sent request after before 2 seconds %s", msg)
	// case <-time.After(4 * time.Second):
	// 	fmt.Println("API Timed Out")
	// }
}
