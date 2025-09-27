package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueue_Enqueue(t *testing.T) {
	queue := NewQueue()

	queue.Enqueue(1)
	assert.False(t, queue.IsEmpty())
	assert.Equal(t, 1, queue.Size())

	queue.Enqueue("hello")
	queue.Enqueue(3.14)
	queue.Enqueue(true)

	assert.Equal(t, 4, queue.Size())
}

func TestQueue_Dequeue(t *testing.T) {
	queue := NewQueue()

	_, err := queue.Dequeue()
	assert.Error(t, err)
	assert.ErrorIs(t, err, errorQueueEmpty)

	queue.Enqueue(1)
	queue.Enqueue(2)
	queue.Enqueue(3)

	item, err := queue.Dequeue()
	assert.NoError(t, err)
	assert.Equal(t, 1, item)
	assert.Equal(t, 2, queue.Size())

	item, _ = queue.Dequeue()
	assert.Equal(t, 2, item)

	item, _ = queue.Dequeue()
	assert.Equal(t, 3, item)

	assert.True(t, queue.IsEmpty())
}

func TestQueue_Front(t *testing.T) {
	queue := NewQueue()

	_, err := queue.Front()
	assert.Error(t, err)
	assert.ErrorIs(t, err, errorQueueEmpty)

	queue.Enqueue("first")
	queue.Enqueue("second")
	queue.Enqueue("third")

	initialSize := queue.Size()

	item, err := queue.Front()
	assert.NoError(t, err)
	assert.Equal(t, "first", item)
	assert.Equal(t, initialSize, queue.Size())

	item2, _ := queue.Front()
	assert.Equal(t, item, item2)
}

func TestQueue_Back(t *testing.T) {
	queue := NewQueue()

	_, err := queue.Back()
	assert.Error(t, err)

	queue.Enqueue("first")
	queue.Enqueue("second")
	queue.Enqueue("third")

	item, err := queue.Back()
	assert.NoError(t, err)
	assert.Equal(t, "third", item)
}

func TestQueue_IsEmpty(t *testing.T) {
	queue := NewQueue()

	assert.True(t, queue.IsEmpty())

	queue.Enqueue(1)
	assert.False(t, queue.IsEmpty())

	queue.Dequeue()
	assert.True(t, queue.IsEmpty())
}

func TestQueue_Size(t *testing.T) {
	queue := NewQueue()

	assert.Equal(t, 0, queue.Size())

	for i := 1; i <= 5; i++ {
		queue.Enqueue(i)
		assert.Equal(t, i, queue.Size())
	}

	for i := 4; i >= 0; i-- {
		queue.Dequeue()
		assert.Equal(t, i, queue.Size())
	}
}

func TestQueue_Clear(t *testing.T) {
	queue := NewQueue()

	queue.Clear()
	assert.True(t, queue.IsEmpty())

	queue.Enqueue(1)
	queue.Enqueue(2)
	queue.Enqueue(3)

	queue.Clear()

	assert.True(t, queue.IsEmpty())
	assert.Equal(t, 0, queue.Size())
}