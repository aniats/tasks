package lru

import (
	"errors"
	"sync"
)

var (
	ErrorZeroCapacity = errors.New("capacity must be greater than zero")
	ErrorEmptyKey     = errors.New("key cannot be empty")
)

type node[V comparable] struct {
	key  string
	val  V
	prev *node[V]
	next *node[V]
}

type list[V comparable] struct {
	head *node[V]
	tail *node[V]
}

type LRU[V comparable] struct {
	mu       sync.RWMutex
	index    map[string]*node[V]
	list     *list[V]
	capacity int64
	size     int64
}

func NewLRU[V comparable](capacity int64) (*LRU[V], error) {
	if capacity <= 0 {
		return nil, ErrorZeroCapacity
	}

	head := &node[V]{}
	tail := &node[V]{}
	head.next = tail
	tail.prev = head

	return &LRU[V]{
		capacity: capacity,
		size:     0,
		list:     &list[V]{head: head, tail: tail},
		index:    make(map[string]*node[V]),
	}, nil
}

func (l *list[V]) remove(n *node[V]) {
	p := n.prev
	q := n.next
	p.next = q
	q.prev = p
	n.prev = nil
	n.next = nil
}

func (l *list[V]) pushFront(n *node[V]) {
	n.prev = l.head
	n.next = l.head.next
	l.head.next.prev = n
	l.head.next = n
}

func (l *list[V]) popTail() *node[V] {
	lru := l.tail.prev
	if lru == l.head {
		return nil
	}

	l.remove(lru)
	return lru
}

func (l *list[V]) moveToFront(n *node[V]) {
	l.remove(n)
	l.pushFront(n)
}

func (c *LRU[V]) Get(key string) (val V, ok bool) {
	if key == "" {
		var zero V
		return zero, false
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	n, found := c.index[key]
	if !found {
		var zero V
		return zero, false
	}

	c.list.moveToFront(n)
	return n.val, true
}

func (c *LRU[V]) Put(key string, val V) error {
	if key == "" {
		return ErrorEmptyKey
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if n, ok := c.index[key]; ok {
		n.val = val
		c.list.moveToFront(n)
		return nil
	}

	n := &node[V]{key: key, val: val}
	c.list.pushFront(n)

	c.index[key] = n
	c.size++

	if c.size > c.capacity {
		lru := c.list.popTail()
		if lru != nil {
			delete(c.index, lru.key)
			c.size--
		}
	}
	return nil
}

func (c *LRU[V]) Delete(key string) bool {
	if key == "" {
		return false
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	n := c.index[key]
	if n == nil {
		return false
	}

	c.list.remove(n)
	delete(c.index, key)
	c.size--
	return true
}

func (c *LRU[V]) Peek(key string) (val V, ok bool) {
	if key == "" {
		var zero V
		return zero, false
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	n, found := c.index[key]
	if !found {
		var zero V
		return zero, false
	}
	return n.val, true
}

func (c *LRU[V]) Len() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.size
}
