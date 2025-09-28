package lru

import (
	"errors"
	"sync"
)

var (
	ErrorZeroCapacity         = errors.New("capacity must be greater than zero")
	ErrorEmptyKey             = errors.New("key cannot be empty")
	ErrorKeyNotFound          = errors.New("key not found")
	ErrorCacheEmpty           = errors.New("cache is empty")
	ErrorPushNilNode          = errors.New("cannot push nil node")
	ErrorCacheStructureBroken = errors.New("cache structure is corrupted")
	ErrorCacheSize            = errors.New("invalid cache size")
	ErrorRemoveNilNode        = errors.New("cannot remove nil node")
)

type node[V comparable] struct {
	key  string
	val  V
	prev *node[V]
	next *node[V]
}

type LRU[V comparable] struct {
	mu       sync.RWMutex
	index    map[string]*node[V]
	head     *node[V]
	tail     *node[V]
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
		head:     head,
		tail:     tail,
		index:    make(map[string]*node[V]),
	}, nil
}

func (c *LRU[V]) remove(n *node[V]) error {
	if n == nil {
		return ErrorRemoveNilNode
	}

	p := n.prev
	q := n.next
	p.next = q
	q.prev = p
	n.prev = nil
	n.next = nil
	return nil
}

func (c *LRU[V]) pushFront(n *node[V]) error {
	if n == nil {
		return ErrorPushNilNode
	}

	n.prev = c.head
	n.next = c.head.next
	c.head.next.prev = n
	c.head.next = n
	return nil
}

func (c *LRU[V]) popTail() (*node[V], error) {
	if c.tail == nil || c.tail.prev == nil {
		return nil, ErrorCacheStructureBroken
	}

	lru := c.tail.prev
	if lru == c.head {
		return nil, ErrorCacheEmpty
	}

	if err := c.remove(lru); err != nil {
		return nil, err
	}
	return lru, nil
}

func (c *LRU[V]) moveToFront(n *node[V]) error {
	if err := c.remove(n); err != nil {
		return err
	}
	return c.pushFront(n)
}

func (c *LRU[V]) Get(key string) (val V, err error) {
	if key == "" {
		var zero V
		return zero, ErrorEmptyKey
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	n, ok := c.index[key]
	if !ok {
		var zero V
		return zero, ErrorKeyNotFound
	}

	if err := c.moveToFront(n); err != nil {
		var zero V
		return zero, err
	}
	return n.val, nil
}

func (c *LRU[V]) Put(key string, val V) (evictedKey string, evictedVal V, evicted bool, err error) {
	if key == "" {
		var zero V
		return "", zero, false, ErrorEmptyKey
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if n, ok := c.index[key]; ok {
		n.val = val
		if err := c.moveToFront(n); err != nil {
			var zero V
			return "", zero, false, err
		}
		var zero V
		return "", zero, false, nil
	}

	n := &node[V]{key: key, val: val}
	if err := c.pushFront(n); err != nil {
		var zero V
		return "", zero, false, err
	}

	c.index[key] = n
	c.size++

	if c.size > c.capacity {
		lru, err := c.popTail()
		if err != nil {
			var zero V
			return "", zero, false, err
		}

		if lru != nil {
			delete(c.index, lru.key)
			c.size--
			return lru.key, lru.val, true, nil
		}
	}
	var zero V
	return "", zero, false, nil
}

func (c *LRU[V]) Delete(key string) error {
	if key == "" {
		return ErrorEmptyKey
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	n := c.index[key]
	if n == nil {
		return ErrorKeyNotFound
	}

	if err := c.remove(n); err != nil {
		return err
	}

	delete(c.index, key)
	c.size--
	return nil
}

func (c *LRU[V]) Peek(key string) (val V, err error) {
	if key == "" {
		var zero V
		return zero, ErrorEmptyKey
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	n, ok := c.index[key]
	if !ok {
		var zero V
		return zero, ErrorKeyNotFound
	}
	return n.val, nil
}

func (c *LRU[V]) Len() (int64, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.size < 0 {
		return 0, ErrorCacheSize
	}

	return c.size, nil
}
