package main

import (
	"errors"
	"fmt"
)

type Stack struct {
	items []int
}

// Push adds an element to the top
func (s *Stack) Push(item int) {
	s.items = append(s.items, item)

}

// Pop removes and returns the top element
func (s *Stack) Pop() (int, error) {
	if s.IsEmpty() {
		return -1, errors.New("Empty Stack")
	}
	lidx := len(s.items) - 1
	l := s.items[lidx]
	s.items = s.items[:lidx]
	return l, nil
}

// Peek returns the top element without removing it
func (s *Stack) Peek() (int, error) {
	if s.IsEmpty() {
		return -1, errors.New("Empty Stack")
	}

	return s.items[len(s.items)-1], nil
}

// IsEmpty checks if stack is empty
func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

func main() {
	var s Stack

	// Test IsEmpty
	fmt.Println("Is empty?", s.IsEmpty()) // true

	// Test Push
	s.Push(10)
	s.Push(20)
	s.Push(30)
	fmt.Println("After pushes, is empty?", s.IsEmpty()) // false

	// Test Peek
	if val, err := s.Peek(); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Peek:", val) // 30
	}

	// Test Pop
	if val, err := s.Pop(); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Popped:", val) // 30
	}

	if val, err := s.Peek(); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Peek after pop:", val) // 20
	}

	// Pop all
	s.Pop() // 20
	s.Pop() // 10

	// Test empty pop
	if val, err := s.Pop(); err != nil {
		fmt.Println("Error on empty pop:", err) // Should error
	} else {
		fmt.Println("Popped:", val)
	}
}

/*

---

## Expected Output

Is empty? true
After pushes, is empty? false
Peek: 30
Popped: 30
Peek after pop: 20
Error on empty pop: stack is empty
*/
