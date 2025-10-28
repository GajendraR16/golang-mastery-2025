package main

import "fmt"

func main() {
	x, y := 48, 18
	for y > 0 {
		x, y = y, x%y
	}
	fmt.Println(x)
}
