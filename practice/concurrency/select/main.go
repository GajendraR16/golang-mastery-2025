package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func slowJobs(ctx context.Context, job chan<- int, wg *sync.WaitGroup) {
	timer := time.NewTicker(3 * time.Second)
	defer timer.Stop()
	defer wg.Done()

	select {
	case <-timer.C:
		fmt.Println("3 Second Passed")
		job <- 2
	case <-ctx.Done():
		fmt.Println("Cancelled slowJob")
		return
	}

}

func main() {
	job := make(chan int)
	var wg sync.WaitGroup

	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go slowJobs(ctx, job, &wg)
	select {
	case val := <-job:
		fmt.Println(val)
	case <-time.After(1 * time.Second):
		fmt.Println("timeout")
		cancel()
	}
	wg.Wait()
}
