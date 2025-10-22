package lru

import (
	"errors"
	"lru/list"
	"lru/node"
	"sync"
)

var (
	ErrorZeroCapacity = errors.New("capacity must be greater than zero")
)

type CacheItem[V comparable] struct {
	Key  string
	Node *node.Node[V]
}

type LRU[V comparable] struct {
	mu       *sync.RWMutex
	index    map[string]*CacheItem[V]
	list     *list.List[V]
	capacity int64
	size     int64
}

func NewLRU[V comparable](capacity int64) (*LRU[V], error) {
	if capacity <= 0 {
		return nil, ErrorZeroCapacity
	}

	return &LRU[V]{
		mu:       &sync.RWMutex{},
		capacity: capacity,
		size:     0,
		list:     list.New[V](),
		index:    make(map[string]*CacheItem[V]),
	}, nil
}

func (c *LRU[V]) Get(key string) (val V, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, found := c.index[key]
	if !found {
		var zero V
		return zero, false
	}

	c.list.MoveToFront(item.Node)
	return item.Node.Val, true
}

func (c *LRU[V]) Put(key string, val V) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, ok := c.index[key]; ok {
		item.Node.Val = val
		c.list.MoveToFront(item.Node)
		return nil
	}

	n := node.New(key, val)
	c.list.PushFront(n)

	item := &CacheItem[V]{
		Key:  key,
		Node: n,
	}
	c.index[key] = item
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
	c.mu.Lock()
	defer c.mu.Unlock()

	item := c.index[key]
	if item == nil {
		return false
	}

	c.list.Remove(item.Node)
	delete(c.index, key)
	c.size--
	return true
}

func (c *LRU[V]) Peek(key string) (val V, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.index[key]
	if !found {
		var zero V
		return zero, false
	}
	return item.Node.Val, true
}

func (c *LRU[V]) Len() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.size
}
