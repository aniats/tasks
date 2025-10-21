package stack

import (
	"errors"
	"fmt"
)

var ErrStackEmpty = errors.New("stack is empty")

type Stack struct {
	items []interface{}
}

func New() *Stack {
	return &Stack{
		items: make([]interface{}, 0),
	}
}

func (s *Stack) Push(item interface{}) {
	s.items = append(s.items, item)
}

func (s *Stack) Pop() (interface{}, error) {
	if s.IsEmpty() {
		return nil, ErrStackEmpty
	}

	index := len(s.items) - 1
	item := s.items[index]
	s.items = s.items[:index]

	return item, nil
}

func (s *Stack) Peek() (interface{}, error) {
	if s.IsEmpty() {
		return nil, ErrStackEmpty
	}

	return s.items[len(s.items)-1], nil
}

func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

func (s *Stack) Size() int {
	return len(s.items)
}

func (s *Stack) Clear() {
	s.items = s.items[:0]
}

func (s *Stack) String() string {
	if s.IsEmpty() {
		return "Stack is empty: []"
	}

	return fmt.Sprintf("Stack: %v (top: %v)", s.items, s.items[len(s.items)-1])
}