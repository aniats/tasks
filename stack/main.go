package main

import (
	"errors"
	"fmt"
)

var errorStackEmpty = errors.New("стек пуст")

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
		return "Стек пуст: []"
	}

	return fmt.Sprintf("Стек: %v (сверху: %v)", s.items, s.items[len(s.items)-1])
}

func main() {
	stack := NewStack()

	stack.Push(1)
	stack.Push("hello")
	stack.Push(3.14)
	stack.Push(true)

	fmt.Println("Стек:", stack)
	fmt.Println("Размер стека:", stack.Size())

	if top, err := stack.Peek(); err == nil {
		fmt.Println("Элемент сверху:", top)
	}

	fmt.Println("Извлечение:")
	for !stack.IsEmpty() {
		if item, err := stack.Pop(); err == nil {
			fmt.Printf("Извлечен: %v, осталось: %d\n", item, stack.Size())
		}
	}

	if _, err := stack.Pop(); err != nil {
		fmt.Println("Ошибка:", err)
	}
}
