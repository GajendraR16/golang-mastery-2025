package main

import (
	"time"
)

func SlowApi(result chan<- string) {
	time.Sleep(5 * time.Second)
	result <- "Data Received"
}
