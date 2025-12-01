
# LeetCode Solutions - Week 1

## Problem 1: Two Sum
**Link:** https://leetcode.com/problems/two-sum/description/

```go
func twoSum(nums []int, target int) []int {
	m := make(map[int]int, len(nums))
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

## Problem 6: Contains Duplicate	
**Link:** https://leetcode.com/problems/contains-duplicate/
```go
func containsDuplicate(nums []int) bool {
	m := make(map[int]bool, len(nums)-1)
	for _, num := range nums {
		if m[num] {
			return true
		}
		m[num] = true
	}
	return false
}
```

## Problem 7: Valid Anagram
**Link:** https://leetcode.com/problems/valid-anagram/
```go
func validAnagram(s, t string) bool {
	if len(s) != len(t) {
		return false
	}

	count := make(map[rune]int, len(s))

	for _, str := range s {
		count[str]++
	}

	for _, v := range t {
		count[v]--
		if count[v] < 0 {
			return false
		}
	}
	return true
}
```

## Problem 8: Reverse Linked List
**Link:** https://leetcode.com/problems/reverse-linked-list/
```go
func reverseList(head *ListNode) *ListNode {
	var prev *ListNode
	curr := head

	for curr != nil{
		next := curr.Next
		curr.Next = prev
		prev = curr
		curr = next
	}

	return prev
}
```

## Problem 9: Merge Two Sorted Lists
**Link:** https://leetcode.com/problems/merge-two-sorted-lists/
```go
func mergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode  {
	dummy := &ListNode{}
	tail := dummy

	for list1 != nil && list2 != nil{
		if list1.Val < list2.Val{
			tail.Next = list1
			list1 = list1.Next
		} else {
			tail.Next = list2
			list2 = list2.Next
		}
		tail = tail.Next
	}

	for list1 != nil {
		tail.Next = list1
		list1 = list1.Next
		tail = tail.Next
	}

	for list2 != nil {
		tail.Next = list2
		list2 = list2.Next
		tail = tail.Next
	}

	return dummy.Next
}
```