package util

import "fmt"

type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(data T) {
	s.items = append(s.items, data)
}

func (s *Stack[T]) Pop() T {
	if s.IsEmpty() {
		return *new(T)
	}
	t := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return t
}

func (s *Stack[T]) Top() (T, error) {
	if s.IsEmpty() {
		return *new(T), fmt.Errorf("stack is empty")
	}
	return s.items[len(s.items)-1], nil
}

func (s *Stack[T]) IsEmpty() bool {
	if len(s.items) == 0 {
		return true
	}
	return false
}

func (s *Stack[T]) Len() int {
	return len(s.items)
}

func (s *Stack[T]) Print() {
	for _, item := range s.items {
		fmt.Print(item, " ")
	}
	fmt.Println()
}
