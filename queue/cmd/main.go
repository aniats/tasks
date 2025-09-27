package main

import (
	"fmt"
	"queue/queue"
)

func main() {
	q := queue.NewQueue()

	q.Enqueue(1)
	q.Enqueue("hello")
	q.Enqueue(3.14)
	q.Enqueue(true)

	fmt.Println("Queue after adding elements:", q)
	fmt.Println("Queue size:", q.Size())

	if front, err := q.Front(); err == nil {
		fmt.Println("First element:", front)
	}

	if back, err := q.Back(); err == nil {
		fmt.Println("Last element:", back)
	}

	fmt.Println("Extracting elements:")
	for !q.IsEmpty() {
		if item, err := q.Dequeue(); err == nil {
			fmt.Printf("Extracted: %v, remaining: %d\n", item, q.Size())
		}
	}

	if _, err := q.Dequeue(); err != nil {
		fmt.Println("Error:", err)
	}
}
