package main

type ListNode struct {
	Val  int
	Next *ListNode
}

func main() {
	// Array/Map problems
	isPalindrome(121)
	twoSum([]int{2, 7, 11, 15}, 9)
	removeDuplicates([]int{1, 1, 2})
	removeElement([]int{3, 2, 2, 3}, 3)
	merge([]int{1, 2, 3, 0, 0, 0}, 3, []int{2, 5, 6}, 3)
	containsDuplicate([]int{1, 1, 1, 3, 3, 4, 3, 2, 4, 2})
	validAnagram("rat", "cat")

	// Linked list problems
	list1 := createsList([]int{1, 2, 4})
	reverseList(list1)

	list2 := createsList([]int{1, 2, 4})
	list3 := createsList([]int{1, 3, 4})
	mergeTwoLists(list2, list3)

	list4 := createsList([]int{1, 2, 3, 4, 5})
	hasCycle(list4)

	list5 := createsList([]int{1, 2, 6, 3, 4, 5, 6})
	removeElements(list5, 6)
}

// Helper function to create a linked list from a slice
func createsList(vals []int) *ListNode {
	if len(vals) == 0 {
		return nil
	}
	head := &ListNode{Val: vals[0]}
	curr := head
	for i := 1; i < len(vals); i++ {
		curr.Next = &ListNode{Val: vals[i]}
		curr = curr.Next
	}
	return head
}

func twoSum(nums []int, target int) []int {
	m := make(map[int]int, len(nums))
	/* Adding capacity imporves performance,
	   as it reduces the number of rehashing
	   and growing the map multiple times
	*/
	for idx, num := range nums {
		if res, ok := m[target-num]; ok {
			return []int{idx, res}
		}
		m[num] = idx
	}
	return nil
}

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

func removeDuplicates(nums []int) int {
	m := make(map[int]bool)
	j := 0
	for _, val := range nums {
		if !m[val] {
			m[val] = true
			nums[j] = val
			j++
		}
	}
	return j
}

func removeElement(nums []int, val int) int {
	j := 0
	for _, value := range nums {
		if value != val {
			nums[j] = value
			j++
		}
	}
	return j
}

func merge(nums1 []int, m int, nums2 []int, n int) {
	i, j, k := m-1, n-1, m+n-1
	for i >= 0 && j >= 0 && k >= 0 {
		if nums1[i] > nums2[j] {
			nums1[k] = nums1[i]
			i--
			k--
		} else {
			nums1[k] = nums2[j]
			j--
			k--
		}
	}

	for j >= 0 {
		nums1[k] = nums2[j]
		j--
		k--
	}
}

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

func reverseList(head *ListNode) *ListNode {
	var prev *ListNode
	curr := head

	for curr != nil {
		next := curr.Next
		curr.Next = prev
		prev = curr
		curr = next
	}

	return prev
}

func mergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
	dummy := &ListNode{}
	tail := dummy

	for list1 != nil && list2 != nil {
		if list1.Val < list2.Val {
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

func hasCycle(head *ListNode) bool {
	slow := head
	fast := head

	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next

		if slow == fast {
			return true
		}
	}
	return false
}

func removeElements(head *ListNode, val int) *ListNode {
	dummy := &ListNode{}
	dummy.Next = head

	prev := dummy
	curr := head

	for curr != nil {
		if curr.Val == val {
			prev.Next = curr.Next
			curr = curr.Next
		} else {
			prev = curr
			curr = curr.Next
		}
	}

	return dummy.Next
}
