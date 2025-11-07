package main

import (
	"fmt"
)

func reverseSlice[T any](nums []T) {
	i, j := 0, len(nums)-1
	for i < j {
		nums[i], nums[j] = nums[j], nums[i]
		i++
		j--
	}
}

func main() {
	param := []string{"h", "e", "l", "l", "o"}
	reverseSlice(param)
	fmt.Println(param)
}
