package main

import (
	"fmt"
	"time"
)

func pipeline() {
	naturals := make(chan int, 10)
	squares := make(chan int, 10)

	go natural(naturals)
	go squared(naturals, squares)
	go printer(squares)

	time.Sleep(time.Second)
}

func natural(naturals chan<- int) {
	for x := 1; x <= 10; x++ {
		naturals <- x
	}
	close(naturals)
}

func squared(naturals <-chan int, squares chan<- int) {
	for x := range naturals {
		squares <- x * x
	}
	close(squares)
}

func printer(squares <-chan int) {
	for x := range squares {
		fmt.Println(x)
	}
}
