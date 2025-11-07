
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
## Problem 3: Remove Duplicates from a sorted array
**Link:** https://leetcode.com/problems/remove-duplicates-from-sorted-array/
```go
func removeDuplicates(nums []int) int {
    m := make(map[int]bool)
    j := 0
    for _, val := range nums{
        if !m[val]{
            m[val] = true
            nums[j] = val
            j++
        }
    }
    return j
}
```

## Problem 4: Remove Element
**Link:** https://leetcode.com/problems/remove-element/
```go
func removeElement(nums []int, val int) int {
    j := 0
    for _, value := range nums{
        if value != val {
            nums[j] = value
            j++
        } 
    }
    return j
}
```

## Problem 5: Merge Sorted Array	
**Link:** https://leetcode.com/problems/merge-sorted-array/
```go
func merge(nums1 []int, m int, nums2 []int, n int)  {
    i,j,k := m-1, n-1, m+n-1
    for i>=0 && j>=0 && k>=0{
        if nums1[i] > nums2[j]{
            nums1[k] = nums1[i]
            i--
            k--
        } else {
            nums1[k] = nums2[j]
            j--
            k--
        }
    }

    for j >=0 {
        nums1[k] = nums2[j]
        j--
        k--
    }
}
```