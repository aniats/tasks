package main

import (
	"fmt"

	"stack/stack"
)

func main() {
	s := stack.New()

	s.Push(1)
	s.Push("hello")
	s.Push(3.14)
	s.Push(true)

	fmt.Println("Stack:", s)
	fmt.Println("Stack size:", s.Size())

	if top, err := s.Peek(); err == nil {
		fmt.Println("Top element:", top)
	}

	fmt.Println("Extracting:")
	for !s.IsEmpty() {
		if item, err := s.Pop(); err == nil {
			fmt.Printf("Extracted: %v, remaining: %d\n", item, s.Size())
		}
	}

	if _, err := s.Pop(); err != nil {
		fmt.Println("Error:", err)
	}
}