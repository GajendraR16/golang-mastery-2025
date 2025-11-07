package main

import "fmt"

func main() {
	// Experiment 1: Watch capacity grow
	var s []int
	for i := 0; i < 10; i++ {
		s = append(s, i)
		fmt.Printf("len=%d cap=%d\n", len(s), cap(s))
	}

	// Experiment 2: Slices share underlying array
	arr := [5]int{1, 2, 3, 4, 5}
	slice1 := arr[1:4]
	slice2 := arr[1:4]
	slice1[0] = 99
	fmt.Println("slice1:", slice1)
	fmt.Println("slice2:", slice2)
	fmt.Println("arr:", arr)

	// Experiment 3: Make vs literal
	a := make([]int, 3, 5)
	b := []int{1, 2, 3}
	fmt.Printf("a: len=%d cap=%d\n", len(a), cap(a))
	fmt.Printf("b: len=%d cap=%d\n", len(b), cap(b))
}
