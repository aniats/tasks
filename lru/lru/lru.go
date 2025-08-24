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

type node struct {
	key  string
	val  interface{}
	prev *node
	next *node
}

type LRU struct {
	mu       sync.RWMutex
	index    map[string]*node
	head     *node
	tail     *node
	capacity int64
	size     int64
}

func NewLRU(capacity int64) (*LRU, error) {
	if capacity <= 0 {
		return nil, ErrorZeroCapacity
	}

	head := &node{}
	tail := &node{}
	head.next = tail
	tail.prev = head

	return &LRU{
		capacity: capacity,
		size:     0,
		head:     head,
		tail:     tail,
		index:    make(map[string]*node),
	}, nil
}

func (c *LRU) remove(n *node) error {
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

func (c *LRU) pushFront(n *node) error {
	if n == nil {
		return ErrorPushNilNode
	}

	n.prev = c.head
	n.next = c.head.next
	c.head.next.prev = n
	c.head.next = n
	return nil
}

func (c *LRU) popTail() (*node, error) {
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

func (c *LRU) moveToFront(n *node) error {
	if err := c.remove(n); err != nil {
		return err
	}
	return c.pushFront(n)
}

func (c *LRU) Get(key string) (val interface{}, err error) {
	if key == "" {
		return nil, ErrorEmptyKey
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	n, ok := c.index[key]
	if !ok {
		return nil, ErrorKeyNotFound
	}

	if err := c.moveToFront(n); err != nil {
		return nil, err
	}
	return n.val, nil
}

func (c *LRU) Put(key string, val interface{}) (evictedKey string, evictedVal interface{}, evicted bool, err error) {
	if key == "" {
		return "", nil, false, ErrorEmptyKey
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if n, ok := c.index[key]; ok {
		n.val = val
		if err := c.moveToFront(n); err != nil {
			return "", nil, false, err
		}
		return "", nil, false, nil
	}

	n := &node{key: key, val: val}
	if err := c.pushFront(n); err != nil {
		return "", nil, false, err
	}

	c.index[key] = n
	c.size++

	if c.size > c.capacity {
		lru, err := c.popTail()
		if err != nil {
			return "", nil, false, err
		}

		if lru != nil {
			delete(c.index, lru.key)
			c.size--
			return lru.key, lru.val, true, nil
		}
	}
	return "", nil, false, nil
}

func (c *LRU) Delete(key string) error {
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

func (c *LRU) Peek(key string) (val interface{}, err error) {
	if key == "" {
		return nil, ErrorEmptyKey
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	n, ok := c.index[key]
	if !ok {
		return nil, ErrorKeyNotFound
	}
	return n.val, nil
}

func (c *LRU) Len() (int64, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.size < 0 {
		return 0, ErrorCacheSize
	}

	return c.size, nil
}
