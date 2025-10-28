
# LeetCode Solutions - Week 1

## Problem 1: Two Sum
**Link:** https://leetcode.com/problems/two-sum/description/

```go
func twoSum(nums []int, target int) []int {
	m := make(map[int]int)
	for idx, num := range nums {
		if res, ok := m[target-num]; ok {
			return []int{idx, res}
		}
		m[num] = idx
	}
	return nil
}
```

## Problem 2: Palindrome Number
**Link:** https://leetcode.com/problems/palindrome-number/

```go
func isPalindrome(x int) bool {
	dup := x
	var pal, res int
	for dup > 0 {
		pal = dup % 10
		res = res*10 + pal
		dup /= 10
	}
	return x == res
}
```