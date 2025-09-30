package list

import "lru/node"

type List[V comparable] struct {
	head *node.Node[V]
	tail *node.Node[V]
}

func New[V comparable]() *List[V] {
	head := &node.Node[V]{}
	tail := &node.Node[V]{}
	head.Next = tail
	tail.Prev = head

	return &List[V]{
		head: head,
		tail: tail,
	}
}

func (l *List[V]) Remove(n *node.Node[V]) {
	p := n.Prev
	q := n.Next
	p.Next = q
	q.Prev = p
	n.Prev = nil
	n.Next = nil
}

func (l *List[V]) PushFront(n *node.Node[V]) {
	n.Prev = l.head
	n.Next = l.head.Next
	l.head.Next.Prev = n
	l.head.Next = n
}

func (l *List[V]) PopTail() *node.Node[V] {
	lru := l.tail.Prev
	if lru == l.head {
		return nil
	}

	l.Remove(lru)
	return lru
}

func (l *List[V]) MoveToFront(n *node.Node[V]) {
	l.Remove(n)
	l.PushFront(n)
}