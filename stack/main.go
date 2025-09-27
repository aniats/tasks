package main

import (
	"errors"
	"fmt"
)

var errorStackEmpty = errors.New("stack is empty")

type Stack struct {
	items []interface{}
}

func NewStack() *Stack {
	return &Stack{
		items: make([]interface{}, 0),
	}
}

func (s *Stack) Push(item interface{}) {
	s.items = append(s.items, item)
}

func (s *Stack) Pop() (interface{}, error) {
	if s.IsEmpty() {
		return nil, errorStackEmpty
	}

	index := len(s.items) - 1
	item := s.items[index]
	s.items = s.items[:index]

	return item, nil
}

func (s *Stack) Peek() (interface{}, error) {
	if s.IsEmpty() {
		return nil, errorStackEmpty
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

func main() {
	stack := NewStack()

	stack.Push(1)
	stack.Push("hello")
	stack.Push(3.14)
	stack.Push(true)

	fmt.Println("Stack:", stack)
	fmt.Println("Stack size:", stack.Size())

	if top, err := stack.Peek(); err == nil {
		fmt.Println("Top element:", top)
	}

	fmt.Println("Extracting:")
	for !stack.IsEmpty() {
		if item, err := stack.Pop(); err == nil {
			fmt.Printf("Extracted: %v, remaining: %d\n", item, stack.Size())
		}
	}

	if _, err := stack.Pop(); err != nil {
		fmt.Println("Error:", err)
	}
}
