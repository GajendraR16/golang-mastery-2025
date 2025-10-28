package main

import "fmt"

func recurfact(n int) int {
	if n == 1 {
		return 1
	}
	return n * recurfact(n-1)
}

func main() {
	num := 1
	for i := 1; i <= 5; i++ {
		num *= i
	}
	fmt.Printf("Non Recursive Factorial %d\n", num)
	fmt.Printf("Recursive Factorial %d", recurfact(5))
}
