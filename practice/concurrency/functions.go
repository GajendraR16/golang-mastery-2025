package main

import (
	"fmt"
	"sync"
)

// Exercise 1: Basic goroutines (15 mins)
func Basic() {
	var wg sync.WaitGroup
	wg.Go(func() {
		fmt.Println("Spun Up a go Routine")
	})
	wg.Done()
}

// Exercise 2: Pipeline (20 mins)
func counter(out chan<- int) {
	for i := 1; i <= 10; i++ {
		out <- i
	}
	close(out)
}

func squarer(out chan<- int, in <-chan int) {
	for val := range in {
		out <- val * val
	}
	close(out)
}

func printer(out <-chan int) {
	for val := range out {
		fmt.Println(val)
	}
}

// Exercise 3: Worker pool (30 mins)
// 3 workers, 15 jobs, buffered channel
// Each job prints "Worker X processed job Y"

func worker(id int, jobs <-chan int) {
	for job := range jobs {
		fmt.Printf("Worker %d processed job %d\n", id, job)
	}
}

func jobs() {
	var wg sync.WaitGroup
	jobs := make(chan int)

	for i := 1; i <= 3; i++ {
		workerID := i
		wg.Go(func() {
			worker(workerID, jobs)
		})
	}

	for i := 1; i <= 15; i++ {
		jobs <- i
	}

	close(jobs)
	wg.Wait()
}
