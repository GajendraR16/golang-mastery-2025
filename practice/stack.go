package main

type Stack struct {
	list []int
}

func (s *Stack) Push(value int) {
	s.list = append(s.list, value)
}

func (s *Stack) Pop() (int, bool) {
	if len(s.list) == 0 {
		return 0, false
	}
	val := s.list[len(s.list)-1]
	s.list = s.list[:len(s.list)-1]
	return val, true
}

func (s *Stack) Peek() (int, bool) {
	if len(s.list) == 0 {
		return 0, false
	}
	return s.list[len(s.list)-1], true
}
