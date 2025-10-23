package queue

import (
	"errors"
	"fmt"
)

var errorQueueEmpty = errors.New("queue is empty")

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
	q.items = nil
}

func (q *Queue) String() string {
	if q.IsEmpty() {
		return "Queue is empty: []"
	}

	return fmt.Sprintf("Queue: %v (front: %v, back: %v)",
		q.items, q.items[0], q.items[len(q.items)-1])
}
