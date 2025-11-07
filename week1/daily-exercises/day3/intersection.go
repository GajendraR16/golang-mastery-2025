package main

import (
	"fmt"
)

func intersectionBF(n1, n2 []int) []int {
	c := make([]int, 0, len(n1))
	seen := make(map[int]bool)
	for _, val := range n1 {
		same := false
		for _, value := range n2 {
			if val == value {
				same = true
				break
			}
		}
		if same && !seen[val] {
			c = append(c, val)
			seen[val] = true
		}
	}
	return c
}

func intersectionMP(n1, n2 []int) []int {
	c := make([]int, 0, len(n1))
	m := make(map[int]bool)
	for _, val := range n1 {
		if _, ok := m[val]; !ok {
			m[val] = true
		}
	}
	for _, value := range n2 {
		if m[value] {
			c = append(c, value)
			delete(m, value)
		}
	}
	return c
}

func main() {
	n1 := []int{1, 2, 3}
	n2 := []int{2, 2, 3}

	c := intersectionBF(n1, n2)
	fmt.Println(c)
	// expected -> [3, 4]

}
