package stack

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStack_Push(t *testing.T) {
	tests := []struct {
		name          string
		items         []interface{}
		expectedSize  int
		expectedEmpty bool
	}{
		{
			name:          "push single item",
			items:         []interface{}{1},
			expectedSize:  1,
			expectedEmpty: false,
		},
		{
			name:          "push multiple items",
			items:         []interface{}{1, "hello", 3.14, true},
			expectedSize:  4,
			expectedEmpty: false,
		},
		{
			name:          "push no items",
			items:         []interface{}{},
			expectedSize:  0,
			expectedEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stack := New()

			for _, item := range tt.items {
				stack.Push(item)
			}

			require.Equal(t, tt.expectedSize, stack.Size())
			require.Equal(t, tt.expectedEmpty, stack.IsEmpty())
		})
	}
}

func TestStack_Pop(t *testing.T) {
	tests := []struct {
		name          string
		stack         func() *Stack
		expectedItem  interface{}
		expectedError error
		expectedSize  int
		expectedEmpty bool
	}{
		{
			name: "pop from empty stack",
			stack: func() *Stack {
				return New()
			},
			expectedItem:  nil,
			expectedError: ErrStackEmpty,
			expectedSize:  0,
			expectedEmpty: true,
		},
		{
			name: "pop from single item stack",
			stack: func() *Stack {
				s := New()
				s.Push(42)
				return s
			},
			expectedItem:  42,
			expectedError: nil,
			expectedSize:  0,
			expectedEmpty: true,
		},
		{
			name: "pop from multiple items stack",
			stack: func() *Stack {
				s := New()
				s.Push(1)
				s.Push(2)
				s.Push(3)
				return s
			},
			expectedItem:  3,
			expectedError: nil,
			expectedSize:  2,
			expectedEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stack := tt.stack()
			item, err := stack.Pop()

			if tt.expectedError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedItem, item)
			}

			require.Equal(t, tt.expectedSize, stack.Size())
			require.Equal(t, tt.expectedEmpty, stack.IsEmpty())
		})
	}
}

func TestStack_Peek(t *testing.T) {
	tests := []struct {
		name          string
		stack         func() *Stack
		expectedItem  interface{}
		expectedError error
		expectedSize  int
	}{
		{
			name: "peek empty stack",
			stack: func() *Stack {
				return New()
			},
			expectedItem:  nil,
			expectedError: ErrStackEmpty,
			expectedSize:  0,
		},
		{
			name: "peek single item",
			stack: func() *Stack {
				s := New()
				s.Push("first")
				return s
			},
			expectedItem:  "first",
			expectedError: nil,
			expectedSize:  1,
		},
		{
			name: "peek multiple items",
			stack: func() *Stack {
				s := New()
				s.Push("first")
				s.Push("second")
				s.Push("third")
				return s
			},
			expectedItem:  "third",
			expectedError: nil,
			expectedSize:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stack := tt.stack()
			initialSize := stack.Size()

			item, err := stack.Peek()

			if tt.expectedError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedItem, item)

				item2, _ := stack.Peek()
				require.Equal(t, item, item2)
			}

			require.Equal(t, tt.expectedSize, stack.Size())
			require.Equal(t, initialSize, stack.Size(), "Peek should not change stack size")
		})
	}
}

func TestStack_IsEmpty(t *testing.T) {
	tests := []struct {
		name          string
		stack         func() *Stack
		expectedEmpty bool
	}{
		{
			name: "new stack is empty",
			stack: func() *Stack {
				return New()
			},
			expectedEmpty: true,
		},
		{
			name: "stack with items is not empty",
			stack: func() *Stack {
				s := New()
				s.Push(1)
				return s
			},
			expectedEmpty: false,
		},
		{
			name: "stack becomes empty after popping all items",
			stack: func() *Stack {
				s := New()
				s.Push(1)
				s.Pop()
				return s
			},
			expectedEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stack := tt.stack()
			require.Equal(t, tt.expectedEmpty, stack.IsEmpty())
		})
	}
}

func TestStack_Size(t *testing.T) {
	tests := []struct {
		name         string
		stack        func() *Stack
		expectedSize int
	}{
		{
			name: "empty stack size",
			stack: func() *Stack {
				return New()
			},
			expectedSize: 0,
		},
		{
			name: "single item stack",
			stack: func() *Stack {
				s := New()
				s.Push(1)
				return s
			},
			expectedSize: 1,
		},
		{
			name: "multiple items stack",
			stack: func() *Stack {
				s := New()
				s.Push(1)
				s.Push(2)
				s.Push(3)
				return s
			},
			expectedSize: 3,
		},
		{
			name: "after push and pop operations",
			stack: func() *Stack {
				s := New()
				s.Push(1)
				s.Push(2)
				s.Pop()
				s.Push(3)
				s.Pop()
				return s
			},
			expectedSize: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stack := tt.stack()
			require.Equal(t, tt.expectedSize, stack.Size())
		})
	}
}

func TestStack_Clear(t *testing.T) {
	tests := []struct {
		name  string
		stack func() *Stack
	}{
		{
			name: "clear empty stack",
			stack: func() *Stack {
				return New()
			},
		},
		{
			name: "clear single item stack",
			stack: func() *Stack {
				s := New()
				s.Push(1)
				return s
			},
		},
		{
			name: "clear multiple items stack",
			stack: func() *Stack {
				s := New()
				s.Push(1)
				s.Push(2)
				s.Push(3)
				return s
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stack := tt.stack()
			stack.Clear()

			require.True(t, stack.IsEmpty())
			require.Equal(t, 0, stack.Size())
		})
	}
}

func TestStack_String(t *testing.T) {
	tests := []struct {
		name           string
		stack          func() *Stack
		expectedString string
	}{
		{
			name: "empty stack",
			stack: func() *Stack {
				return New()
			},
			expectedString: "Stack is empty: []",
		},
		{
			name: "single item",
			stack: func() *Stack {
				s := New()
				s.Push(1)
				return s
			},
			expectedString: "Stack: [1] (top: 1)",
		},
		{
			name: "multiple items",
			stack: func() *Stack {
				s := New()
				s.Push(1)
				s.Push("hello")
				s.Push(3.14)
				return s
			},
			expectedString: "Stack: [1 hello 3.14] (top: 3.14)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stack := tt.stack()
			require.Equal(t, tt.expectedString, stack.String())
		})
	}
}
