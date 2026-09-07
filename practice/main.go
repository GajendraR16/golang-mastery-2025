package main

import "fmt"

func FizzBuzz() {
	for val := 1; val < 20; val++ {
		if val%3 == 0 && val%5 == 0 {
			fmt.Println("FizzBuzz")
		} else if val%5 == 0 {
			fmt.Println("Buzz")
		} else if val%3 == 0 {
			fmt.Println("Fizz")
		} else {
			fmt.Println(val)
		}
	}
}

func RevereString(s string) string {
	b := []rune(s)
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}

func twoSum(nums []int, target int) []int {
	m := make(map[int]int)
	for idx, num := range nums {
		if val, ok := m[target-num]; ok {
			return []int{idx, val}
		}
		m[num] = idx
	}
	return nil
}

func isValid(s string) bool {
	m := map[rune]rune{
		'}': '{',
		')': '(',
		']': '[',
	}

	var stack []rune
	for _, brac := range s {
		if opening, ok := m[brac]; ok {
			if len(stack) == 0 || stack[len(stack)-1] != opening {
				return false
			}
			stack = stack[:len(stack)-1]
		} else {
			stack = append(stack, brac)
		}
	}
	return len(stack) == 0
}

// func reverseList(head *ListNode) *ListNode {
//     var prev *ListNode
//     curr := head

//     for curr != nil {
//         next := curr.Next
//         curr.Next = prev
//         prev = curr
//         curr = next
//     }

//     return prev
// }

func isAnagram(s string, t string) bool {
	if len(s) != len(t) {
		return false
	}

	m := make(map[rune]int)
	for _, val := range s {
		m[val]++
	}

	for _, val := range t {
		if num, ok := m[val]; !ok || num <= 0 {
			return false
		} else {
			m[val]--
		}
	}
	return true
}

func merge(nums1 []int, m int, nums2 []int, n int) {
	i, j, k := m-1, n-1, m+n-1
	for i >= 0 && j >= 0 {
		if nums1[i] < nums2[j] {
			nums1[k] = nums2[j]
			j--
		}
		nums1[k] = nums1[i]
		i--
		k--
	}

	for j >= 0 {
		nums1[k] = nums2[j]
		j--
		k--
	}
}

func removeElement(nums []int, val int) int {
	i, j := 0, len(nums)-1
	for i < j {
		if nums[i] == val {
			if nums[j] == val{
				j--
			} else if nums[j] != val {
				nums[i], nums[j] = nums[j], nums[i]
				i++
			}
		} else {
			i++
		}
			
	}
	return j
}

func main() {
	nums := []int{3, 2, 2, 3}
	val := 3
	removeElement(nums, val)
}
