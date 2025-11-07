package main

import "fmt"

func mergeSort(n1, n2 []int) []int {
	var i, j int
	m, n := len(n1), len(n2)
	c := make([]int, 0, m+n)
	for i < m && j < n {
		if n1[i] < n2[j] {
			c = append(c, n1[i])
			i++
		} else {
			c = append(c, n2[j])
			j++
		}
	}
	if i != m {
		c = append(c, n1[i:]...)
	}

	if j != n {
		c = append(c, n2[j:]...)
	}

	return c
}

func main() {
	a := []int{1, 3, 5}
	b := []int{2, 4, 6}
	c := mergeSort(a, b)
	fmt.Println(c)
}
