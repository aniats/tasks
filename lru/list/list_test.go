package list

import (
	"lru/node"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want func(*List[string]) bool
	}{
		{
			name: "create empty list",
			want: func(l *List[string]) bool {
				return l != nil &&
					l.head != nil &&
					l.tail != nil &&
					l.head.Next == l.tail &&
					l.tail.Prev == l.head
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New[string]()
			assert.True(t, tt.want(l))
		})
	}
}

func TestPushFront(t *testing.T) {
	tests := []struct {
		name  string
		setup func() (*List[string], []*node.Node[string])
		check func(t *testing.T, l *List[string], nodes []*node.Node[string])
	}{
		{
			name: "push single node",
			setup: func() (*List[string], []*node.Node[string]) {
				l := New[string]()
				n := node.New("key1", "value1")
				l.PushFront(n)
				return l, []*node.Node[string]{n}
			},
			check: func(t *testing.T, l *List[string], nodes []*node.Node[string]) {
				n := nodes[0]
				assert.Equal(t, l.head, n.Prev)
				assert.Equal(t, l.tail, n.Next)
				assert.Equal(t, n, l.head.Next)
				assert.Equal(t, n, l.tail.Prev)
			},
		},
		{
			name: "push three nodes",
			setup: func() (*List[string], []*node.Node[string]) {
				l := New[string]()
				n1 := node.New("key1", "value1")
				n2 := node.New("key2", "value2")
				n3 := node.New("key3", "value3")
				l.PushFront(n1)
				l.PushFront(n2)
				l.PushFront(n3)
				return l, []*node.Node[string]{n1, n2, n3}
			},
			check: func(t *testing.T, l *List[string], nodes []*node.Node[string]) {
				n1, n2, n3 := nodes[0], nodes[1], nodes[2]
				assert.Equal(t, n3, l.head.Next)
				assert.Equal(t, n2, n3.Next)
				assert.Equal(t, n1, n2.Next)
				assert.Equal(t, l.tail, n1.Next)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, nodes := tt.setup()
			tt.check(t, l, nodes)
		})
	}
}

func TestRemove(t *testing.T) {
	tests := []struct {
		name  string
		setup func() (*List[string], []*node.Node[string], *node.Node[string])
		check func(t *testing.T, l *List[string], nodes []*node.Node[string], removed *node.Node[string])
	}{
		{
			name: "remove from two-node list",
			setup: func() (*List[string], []*node.Node[string], *node.Node[string]) {
				l := New[string]()
				n1 := node.New("key1", "value1")
				n2 := node.New("key2", "value2")
				l.PushFront(n1)
				l.PushFront(n2)
				l.Remove(n1)
				return l, []*node.Node[string]{n2}, n1
			},
			check: func(t *testing.T, l *List[string], nodes []*node.Node[string], removed *node.Node[string]) {
				n2 := nodes[0]
				assert.Equal(t, n2, l.head.Next)
				assert.Equal(t, l.tail, n2.Next)
				assert.Equal(t, n2, l.tail.Prev)
				assert.Nil(t, removed.Prev)
				assert.Nil(t, removed.Next)
			},
		},
		{
			name: "remove middle node",
			setup: func() (*List[string], []*node.Node[string], *node.Node[string]) {
				l := New[string]()
				n1 := node.New("key1", "value1")
				n2 := node.New("key2", "value2")
				n3 := node.New("key3", "value3")
				l.PushFront(n1)
				l.PushFront(n2)
				l.PushFront(n3)
				l.Remove(n2)
				return l, []*node.Node[string]{n1, n3}, n2
			},
			check: func(t *testing.T, l *List[string], nodes []*node.Node[string], removed *node.Node[string]) {
				n1, n3 := nodes[0], nodes[1]
				assert.Equal(t, n3, l.head.Next)
				assert.Equal(t, n1, n3.Next)
				assert.Equal(t, n3, n1.Prev)
				assert.Nil(t, removed.Prev)
				assert.Nil(t, removed.Next)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, nodes, removed := tt.setup()
			tt.check(t, l, nodes, removed)
		})
	}
}

func TestPopTail(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *List[string]
		want  *node.Node[string]
		check func(t *testing.T, l *List[string], result *node.Node[string])
	}{
		{
			name: "pop from empty list",
			setup: func() *List[string] {
				return New[string]()
			},
			want: nil,
			check: func(t *testing.T, l *List[string], result *node.Node[string]) {
				assert.Nil(t, result)
			},
		},
		{
			name: "pop from single node list",
			setup: func() *List[string] {
				l := New[string]()
				n1 := node.New("key1", "value1")
				l.PushFront(n1)
				return l
			},
			want: nil,
			check: func(t *testing.T, l *List[string], result *node.Node[string]) {
				assert.NotNil(t, result)
				assert.Equal(t, "key1", result.Key)
				assert.Equal(t, "value1", result.Val)
				assert.Nil(t, result.Prev)
				assert.Nil(t, result.Next)
				assert.Equal(t, l.tail, l.head.Next)
			},
		},
		{
			name: "pop from multiple node list",
			setup: func() *List[string] {
				l := New[string]()
				n1 := node.New("key1", "value1")
				n2 := node.New("key2", "value2")
				n3 := node.New("key3", "value3")
				l.PushFront(n1)
				l.PushFront(n2)
				l.PushFront(n3)
				return l
			},
			want: nil,
			check: func(t *testing.T, l *List[string], result *node.Node[string]) {
				assert.NotNil(t, result)
				assert.Equal(t, "key1", result.Key)
				assert.Equal(t, "value1", result.Val)
				assert.Nil(t, result.Prev)
				assert.Nil(t, result.Next)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := tt.setup()
			result := l.PopTail()
			tt.check(t, l, result)
		})
	}
}

func TestMoveToFront(t *testing.T) {
	tests := []struct {
		name  string
		setup func() (*List[string], []*node.Node[string], *node.Node[string])
		check func(t *testing.T, l *List[string], nodes []*node.Node[string], moved *node.Node[string])
	}{
		{
			name: "move tail to front",
			setup: func() (*List[string], []*node.Node[string], *node.Node[string]) {
				l := New[string]()
				n1 := node.New("key1", "value1")
				n2 := node.New("key2", "value2")
				n3 := node.New("key3", "value3")
				l.PushFront(n1)
				l.PushFront(n2)
				l.PushFront(n3)
				l.MoveToFront(n1)
				return l, []*node.Node[string]{n2, n3}, n1
			},
			check: func(t *testing.T, l *List[string], nodes []*node.Node[string], moved *node.Node[string]) {
				_, n3 := nodes[0], nodes[1]
				assert.Equal(t, moved, l.head.Next)
				assert.Equal(t, n3, moved.Next)
				assert.Equal(t, moved, n3.Prev)
			},
		},
		{
			name: "move middle to front",
			setup: func() (*List[string], []*node.Node[string], *node.Node[string]) {
				l := New[string]()
				n1 := node.New("key1", "value1")
				n2 := node.New("key2", "value2")
				n3 := node.New("key3", "value3")
				l.PushFront(n1)
				l.PushFront(n2)
				l.PushFront(n3)
				l.MoveToFront(n2)
				return l, []*node.Node[string]{n1, n3}, n2
			},
			check: func(t *testing.T, l *List[string], nodes []*node.Node[string], moved *node.Node[string]) {
				n1, n3 := nodes[0], nodes[1]
				assert.Equal(t, moved, l.head.Next)
				assert.Equal(t, n3, moved.Next)
				assert.Equal(t, n1, n3.Next)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, nodes, moved := tt.setup()
			tt.check(t, l, nodes, moved)
		})
	}
}
