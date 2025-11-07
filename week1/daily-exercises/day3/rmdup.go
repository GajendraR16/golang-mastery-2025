package main

import "fmt"

func main() {
	nums := []int{1, 2, 2, 3, 4, 4, 5}
	seen := make(map[int]bool)
	j := 0

	for _, v := range nums {
		if !seen[v] {
			seen[v] = true
			nums[j] = v
			j++
		}
	}
	fmt.Println(nums[:j])
}

func RemoveDupSlow(nums []int) []int {
	a := make([]int, 0, len(nums))
	for _, val := range nums {
		unique := true
		for _, value := range a {
			if val == value {
				unique = false
				break
			}
		}
		if unique {
			a = append(a, val)
		}
	}
	return a
}
