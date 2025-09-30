package lru

import (
	"errors"
	"lru/list"
	"lru/node"
	"sync"
)

var (
	ErrorZeroCapacity = errors.New("capacity must be greater than zero")
	ErrorEmptyKey     = errors.New("key cannot be empty")
)

type LRU[V comparable] struct {
	mu       sync.RWMutex
	index    map[string]*node.Node[V]
	list     *list.List[V]
	capacity int64
	size     int64
}

func NewLRU[V comparable](capacity int64) (*LRU[V], error) {
	if capacity <= 0 {
		return nil, ErrorZeroCapacity
	}

	return &LRU[V]{
		capacity: capacity,
		size:     0,
		list:     list.New[V](),
		index:    make(map[string]*node.Node[V]),
	}, nil
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

	c.list.MoveToFront(n)
	return n.Val, true
}

func (c *LRU[V]) Put(key string, val V) error {
	if key == "" {
		return ErrorEmptyKey
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if n, ok := c.index[key]; ok {
		n.Val = val
		c.list.MoveToFront(n)
		return nil
	}

	n := node.New(key, val)
	c.list.PushFront(n)

	c.index[key] = n
	c.size++

	if c.size > c.capacity {
		lru := c.list.PopTail()
		if lru != nil {
			delete(c.index, lru.Key)
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

	c.list.Remove(n)
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
	return n.Val, true
}

func (c *LRU[V]) Len() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.size
}
