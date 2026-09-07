package main

import "fmt"

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

	flag := false
	prev := dummy
	curr := l.Head

	for curr != nil {
		if curr.Value == value {
			flag = true
			prev.Next = curr.Next
		}
		prev = curr
		curr = curr.Next
	}

	l.Head = dummy.Next
	return flag
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
			fmt.Printf("->")
		}
		curr = curr.Next
	}
	fmt.Println()

}
