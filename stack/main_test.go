package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack_Push(t *testing.T) {
	stack := NewStack()

	stack.Push(1)
	assert.False(t, stack.IsEmpty())
	assert.Equal(t, 1, stack.Size())

	stack.Push("hello")
	stack.Push(3.14)
	stack.Push(true)

	assert.Equal(t, 4, stack.Size())
}

func TestStack_Pop(t *testing.T) {
	stack := NewStack()

	_, err := stack.Pop()
	assert.Error(t, err)
	assert.ErrorIs(t, err, errorStackEmpty)

	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	item, err := stack.Pop()
	assert.NoError(t, err)
	assert.Equal(t, 3, item)
	assert.Equal(t, 2, stack.Size())

	item, _ = stack.Pop()
	assert.Equal(t, 2, item)

	item, _ = stack.Pop()
	assert.Equal(t, 1, item)

	assert.True(t, stack.IsEmpty())
}

func TestStack_Peek(t *testing.T) {
	stack := NewStack()

	_, err := stack.Peek()
	assert.Error(t, err)
	assert.ErrorIs(t, err, errorStackEmpty)

	stack.Push("first")
	stack.Push("second")
	stack.Push("third")

	initialSize := stack.Size()

	item, err := stack.Peek()
	assert.NoError(t, err)
	assert.Equal(t, "third", item)
	assert.Equal(t, initialSize, stack.Size())

	item2, _ := stack.Peek()
	assert.Equal(t, item, item2)
}

func TestStack_IsEmpty(t *testing.T) {
	stack := NewStack()

	assert.True(t, stack.IsEmpty())

	stack.Push(1)
	assert.False(t, stack.IsEmpty())

	stack.Pop()
	assert.True(t, stack.IsEmpty())
}

func TestStack_Size(t *testing.T) {
	stack := NewStack()

	assert.Equal(t, 0, stack.Size())

	for i := 1; i <= 5; i++ {
		stack.Push(i)
		assert.Equal(t, i, stack.Size())
	}

	for i := 4; i >= 0; i-- {
		stack.Pop()
		assert.Equal(t, i, stack.Size())
	}
}

func TestStack_Clear(t *testing.T) {
	stack := NewStack()

	stack.Clear()
	assert.True(t, stack.IsEmpty())

	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	stack.Clear()

	assert.True(t, stack.IsEmpty())
	assert.Equal(t, 0, stack.Size())
}

func TestStack_Order(t *testing.T) {
	stack := NewStack()

	expected := []int{1, 2, 3, 4, 5}
	for _, v := range expected {
		stack.Push(v)
	}

	for i := len(expected) - 1; i >= 0; i-- {
		item, err := stack.Pop()
		assert.NoError(t, err)
		assert.Equal(t, expected[i], item)
	}
}