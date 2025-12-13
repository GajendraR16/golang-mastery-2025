package main

import (
	"fmt"
)

// Implement these methods:
// - Insert(value int)
// - Delete(value int) bool
// - Search(value int) *Node
// - Reverse()
// - Display()
// - Length() int

type Node struct {
	Value int
	Next  *Node
}

type LinkedList struct {
	Head *Node
}

func (l *LinkedList) Insert(value int) {
	if l.Head == nil {
		l.Head = &Node{Value: value}
		return
	}

	curr := l.Head
	for curr.Next != nil {
		curr = curr.Next
	}

	curr.Next = &Node{Value: value}
}

func (l *LinkedList) Delete(value int) bool {
	dummy := &Node{}
	dummy.Next = l.Head

	prev := dummy
	curr := l.Head
	flag := false

	for curr != nil {
		if curr.Value == value {
			prev.Next = curr.Next
			curr = curr.Next
			flag = true
		} else {
			prev = curr
			curr = curr.Next
		}
	}
	l.Head = dummy.Next
	return flag
}

func (l *LinkedList) Search(value int) *Node {
	if l.Head == nil {
		return nil
	}

	curr := l.Head
	for curr != nil {
		if curr.Value == value {
			return curr
		}
		curr = curr.Next
	}
	return nil
}

func (l *LinkedList) Reverse() {
	var prev *Node
	curr := l.Head

	for curr != nil {
		next := curr.Next
		curr.Next = prev
		prev = curr
		curr = next
	}

	l.Head = prev
}

func (l *LinkedList) Display() {
	curr := l.Head

	for curr != nil {
		fmt.Print(curr.Value)
		if curr.Next != nil {
			fmt.Print("->")
		}
		curr = curr.Next
	}
	fmt.Println()
}

func (l *LinkedList) Length() int {
	c := 0

	curr := l.Head

	for curr != nil {
		c++
		curr = curr.Next
	}
	return c
}

func main() {
	l := &LinkedList{}

	fmt.Println("Inserting values: 10, 20, 30, 40")
	l.Insert(10)
	l.Insert(20)
	l.Insert(30)
	l.Insert(40)
	l.Display()

	fmt.Printf("Length: %d\n", l.Length())

	fmt.Println("\nSearching for 20:")
	node := l.Search(20)
	if node != nil {
		fmt.Println("Found:", node.Value)
	} else {
		fmt.Println("Not found")
	}

	fmt.Println("\nDeleting 30:")
	if l.Delete(30) {
		fmt.Println("Deleted 30")
	} else {
		fmt.Println("30 not found")
	}
	l.Display()

	fmt.Println("\nDeleting 10 (head):")
	if l.Delete(10) {
		fmt.Println("Deleted 10")
	}
	l.Display()

	fmt.Printf("Length: %d\n", l.Length())

	fmt.Println("\nReversing list:")
	l.Reverse()
	l.Display()
}
