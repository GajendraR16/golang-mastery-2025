package main

import (
	"fmt"
	"math"
)

func main() {
	for i := 2; i < 10; i++ {
		if isPrime(i) {
			fmt.Println(i)
		}
	}
}

func isPrime(num int) bool {
	for i := 2; i <= int(math.Sqrt(float64(num))); i++ {
		if num%i == 0 {
			return false

		}
	}
	return true
}
