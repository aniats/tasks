package node

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		val      string
		wantKey  string
		wantVal  string
		wantPrev *Node[string]
		wantNext *Node[string]
	}{
		{
			name:     "create string node",
			key:      "key1",
			val:      "value1",
			wantKey:  "key1",
			wantVal:  "value1",
			wantPrev: nil,
			wantNext: nil,
		},
		{
			name:     "create empty string node",
			key:      "",
			val:      "",
			wantKey:  "",
			wantVal:  "",
			wantPrev: nil,
			wantNext: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := New(tt.key, tt.val)

			assert.NotNil(t, n)
			assert.Equal(t, tt.wantKey, n.Key)
			assert.Equal(t, tt.wantVal, n.Val)
			assert.Equal(t, tt.wantPrev, n.Prev)
			assert.Equal(t, tt.wantNext, n.Next)
		})
	}
}

func TestNodeLinking(t *testing.T) {
	tests := []struct {
		name  string
		setup func() (*Node[string], *Node[string])
		check func(t *testing.T, n1, n2 *Node[string])
	}{
		{
			name: "link two nodes",
			setup: func() (*Node[string], *Node[string]) {
				n1 := New("key1", "value1")
				n2 := New("key2", "value2")
				n1.Next = n2
				return n1, n2
			},
			check: func(t *testing.T, n1, n2 *Node[string]) {
				assert.Equal(t, n2, n1.Next)
				assert.Nil(t, n2.Next)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n1, n2 := tt.setup()
			tt.check(t, n1, n2)
		})
	}
}
