package queue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueue_Enqueue(t *testing.T) {
	q := NewQueue()

	q.Enqueue(1)
	assert.False(t, q.IsEmpty())
	assert.Equal(t, 1, q.Size())

	q.Enqueue("hello")
	q.Enqueue(3.14)
	q.Enqueue(true)

	assert.Equal(t, 4, q.Size())
}

func TestQueue_Dequeue(t *testing.T) {
	q := NewQueue()

	_, err := q.Dequeue()
	assert.Error(t, err)
	assert.ErrorIs(t, err, errorQueueEmpty)

	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	item, err := q.Dequeue()
	assert.NoError(t, err)
	assert.Equal(t, 1, item)
	assert.Equal(t, 2, q.Size())

	item, _ = q.Dequeue()
	assert.Equal(t, 2, item)

	item, _ = q.Dequeue()
	assert.Equal(t, 3, item)

	assert.True(t, q.IsEmpty())
}

func TestQueue_Front(t *testing.T) {
	q := NewQueue()

	_, err := q.Front()
	assert.Error(t, err)
	assert.ErrorIs(t, err, errorQueueEmpty)

	q.Enqueue("first")
	q.Enqueue("second")
	q.Enqueue("third")

	initialSize := q.Size()

	item, err := q.Front()
	assert.NoError(t, err)
	assert.Equal(t, "first", item)
	assert.Equal(t, initialSize, q.Size())

	item2, _ := q.Front()
	assert.Equal(t, item, item2)
}

func TestQueue_Back(t *testing.T) {
	q := NewQueue()

	_, err := q.Back()
	assert.Error(t, err)

	q.Enqueue("first")
	q.Enqueue("second")
	q.Enqueue("third")

	item, err := q.Back()
	assert.NoError(t, err)
	assert.Equal(t, "third", item)
}

func TestQueue_IsEmpty(t *testing.T) {
	q := NewQueue()

	assert.True(t, q.IsEmpty())

	q.Enqueue(1)
	assert.False(t, q.IsEmpty())

	q.Dequeue()
	assert.True(t, q.IsEmpty())
}

func TestQueue_Size(t *testing.T) {
	q := NewQueue()

	assert.Equal(t, 0, q.Size())

	for i := 1; i <= 5; i++ {
		q.Enqueue(i)
		assert.Equal(t, i, q.Size())
	}

	for i := 4; i >= 0; i-- {
		q.Dequeue()
		assert.Equal(t, i, q.Size())
	}
}

func TestQueue_Clear(t *testing.T) {
	q := NewQueue()

	q.Clear()
	assert.True(t, q.IsEmpty())

	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	q.Clear()

	assert.True(t, q.IsEmpty())
	assert.Equal(t, 0, q.Size())
}
