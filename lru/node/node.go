package node

type Node[V comparable] struct {
	Key  string
	Val  V
	Prev *Node[V]
	Next *Node[V]
}

func New[V comparable](key string, val V) *Node[V] {
	return &Node[V]{
		Key: key,
		Val: val,
	}
}