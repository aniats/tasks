package main

import (
	"errors"
	"fmt"
)

var errorQueueEmpty = errors.New("очередь пуста")

type Queue struct {
	items []interface{}
}

func NewQueue() *Queue {
	return &Queue{
		items: make([]interface{}, 0),
	}
}

func (q *Queue) Enqueue(item interface{}) {
	q.items = append(q.items, item)
}

func (q *Queue) Dequeue() (interface{}, error) {
	if q.IsEmpty() {
		return nil, errorQueueEmpty
	}

	item := q.items[0]
	q.items = q.items[1:]

	return item, nil
}

func (q *Queue) Front() (interface{}, error) {
	if q.IsEmpty() {
		return nil, errorQueueEmpty
	}

	return q.items[0], nil
}

func (q *Queue) Back() (interface{}, error) {
	if q.IsEmpty() {
		return nil, errorQueueEmpty
	}

	return q.items[len(q.items)-1], nil
}

func (q *Queue) IsEmpty() bool {
	return len(q.items) == 0
}

func (q *Queue) Size() int {
	return len(q.items)
}

func (q *Queue) Clear() {
	q.items = q.items[:0]
}

func (q *Queue) String() string {
	if q.IsEmpty() {
		return "Очередь пуста: []"
	}

	return fmt.Sprintf("Очередь: %v (спереди: %v, сзади: %v)",
		q.items, q.items[0], q.items[len(q.items)-1])
}

func main() {
	queue := NewQueue()

	queue.Enqueue(1)
	queue.Enqueue("hello")
	queue.Enqueue(3.14)
	queue.Enqueue(true)

	fmt.Println("Очередь после добавления элементов:", queue)
	fmt.Println("Размер очереди:", queue.Size())

	if front, err := queue.Front(); err == nil {
		fmt.Println("Первый элемент:", front)
	}

	if back, err := queue.Back(); err == nil {
		fmt.Println("Последний элемент:", back)
	}

	fmt.Println("Извлечение элементов:")
	for !queue.IsEmpty() {
		if item, err := queue.Dequeue(); err == nil {
			fmt.Printf("Извлечен: %v, осталось: %d\n", item, queue.Size())
		}
	}

	if _, err := queue.Dequeue(); err != nil {
		fmt.Println("Ошибка:", err)
	}
}
