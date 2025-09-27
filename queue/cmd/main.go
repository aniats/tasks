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

	fmt.Println("Очередь после добавления элементов:", q)
	fmt.Println("Размер очереди:", q.Size())

	if front, err := q.Front(); err == nil {
		fmt.Println("Первый элемент:", front)
	}

	if back, err := q.Back(); err == nil {
		fmt.Println("Последний элемент:", back)
	}

	fmt.Println("Извлечение элементов:")
	for !q.IsEmpty() {
		if item, err := q.Dequeue(); err == nil {
			fmt.Printf("Извлечен: %v, осталось: %d\n", item, q.Size())
		}
	}

	if _, err := q.Dequeue(); err != nil {
		fmt.Println("Ошибка:", err)
	}
}
