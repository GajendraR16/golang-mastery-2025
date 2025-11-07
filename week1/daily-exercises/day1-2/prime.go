package main

import (
	"fmt"
)

func main() {
	for i := 2; i < 10; i++ {
		if isPrime(i) {
			fmt.Println(i)
		}
	}
}

func isPrime(num int) bool {
	for i := 2; i <= num>>1; i++ {
		if num%i == 0 {
			return false

		}
	}
	return true
}
