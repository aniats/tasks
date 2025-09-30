package lru

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLRU(t *testing.T) {
	tests := []struct {
		name     string
		capacity int64
		wantErr  error
	}{
		{
			name:     "valid capacity",
			capacity: 10,
			wantErr:  nil,
		},
		{
			name:     "zero capacity",
			capacity: 0,
			wantErr:  ErrorZeroCapacity,
		},
		{
			name:     "negative capacity",
			capacity: -5,
			wantErr:  ErrorZeroCapacity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache, err := NewLRU[string](tt.capacity)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, cache)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, cache)
			assert.Equal(t, tt.capacity, cache.capacity)
			assert.Equal(t, int64(0), cache.size)
		})
	}
}

func TestLRU_Get(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *LRU[string]
		key      string
		wantVal  string
		wantOk   bool
	}{
		{
			name: "get from empty cache",
			setup: func() *LRU[string] {
				cache, _ := NewLRU[string](2)
				return cache
			},
			key:     "nonexistent",
			wantVal: "",
			wantOk:  false,
		},
		{
			name: "get with empty key",
			setup: func() *LRU[string] {
				cache, _ := NewLRU[string](2)
				return cache
			},
			key:     "",
			wantVal: "",
			wantOk:  false,
		},
		{
			name: "get existing key",
			setup: func() *LRU[string] {
				cache, _ := NewLRU[string](2)
				cache.Put("key1", "value1")
				cache.Put("key2", "value2")
				return cache
			},
			key:     "key1",
			wantVal: "value1",
			wantOk:  true,
		},
		{
			name: "get evicted key",
			setup: func() *LRU[string] {
				cache, _ := NewLRU[string](2)
				cache.Put("key1", "value1")
				cache.Put("key2", "value2")
				cache.Put("key3", "value3") // should evict key1
				return cache
			},
			key:     "key1",
			wantVal: "",
			wantOk:  false,
		},
		{
			name: "get after LRU update",
			setup: func() *LRU[string] {
				cache, _ := NewLRU[string](2)
				cache.Put("key1", "value1")
				cache.Put("key2", "value2")
				cache.Get("key1") // moves key1 to front
				cache.Put("key3", "value3") // should evict key2
				return cache
			},
			key:     "key1",
			wantVal: "value1",
			wantOk:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := tt.setup()
			val, ok := cache.Get(tt.key)
			assert.Equal(t, tt.wantOk, ok)
			assert.Equal(t, tt.wantVal, val)
		})
	}
}

func TestLRU_Put(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *LRU[string]
		key     string
		val     string
		wantErr error
		check   func(t *testing.T, cache *LRU[string])
	}{
		{
			name: "put with empty key",
			setup: func() *LRU[string] {
				cache, _ := NewLRU[string](2)
				return cache
			},
			key:     "",
			val:     "value",
			wantErr: ErrorEmptyKey,
			check: func(t *testing.T, cache *LRU[string]) {
				assert.Equal(t, int64(0), cache.Len())
			},
		},
		{
			name: "put new key",
			setup: func() *LRU[string] {
				cache, _ := NewLRU[string](2)
				return cache
			},
			key:     "key1",
			val:     "value1",
			wantErr: nil,
			check: func(t *testing.T, cache *LRU[string]) {
				val, ok := cache.Get("key1")
				assert.True(t, ok)
				assert.Equal(t, "value1", val)
				assert.Equal(t, int64(1), cache.Len())
			},
		},
		{
			name: "update existing key",
			setup: func() *LRU[string] {
				cache, _ := NewLRU[string](2)
				cache.Put("key1", "value1")
				return cache
			},
			key:     "key1",
			val:     "updated_value1",
			wantErr: nil,
			check: func(t *testing.T, cache *LRU[string]) {
				val, ok := cache.Get("key1")
				assert.True(t, ok)
				assert.Equal(t, "updated_value1", val)
				assert.Equal(t, int64(1), cache.Len())
			},
		},
		{
			name: "put beyond capacity",
			setup: func() *LRU[string] {
				cache, _ := NewLRU[string](2)
				cache.Put("key1", "value1")
				cache.Put("key2", "value2")
				return cache
			},
			key:     "key3",
			val:     "value3",
			wantErr: nil,
			check: func(t *testing.T, cache *LRU[string]) {
				// key1 should be evicted
				_, ok := cache.Get("key1")
				assert.False(t, ok)

				val2, ok := cache.Get("key2")
				assert.True(t, ok)
				assert.Equal(t, "value2", val2)

				val3, ok := cache.Get("key3")
				assert.True(t, ok)
				assert.Equal(t, "value3", val3)

				assert.Equal(t, int64(2), cache.Len())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := tt.setup()
			err := cache.Put(tt.key, tt.val)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			tt.check(t, cache)
		})
	}
}

func TestLRU_Delete(t *testing.T) {
	tests := []struct {
		name   string
		setup  func() *LRU[string]
		key    string
		wantOk bool
		check  func(t *testing.T, cache *LRU[string])
	}{
		{
			name: "delete with empty key",
			setup: func() *LRU[string] {
				cache, _ := NewLRU[string](3)
				return cache
			},
			key:    "",
			wantOk: false,
			check: func(t *testing.T, cache *LRU[string]) {
				assert.Equal(t, int64(0), cache.Len())
			},
		},
		{
			name: "delete nonexistent key",
			setup: func() *LRU[string] {
				cache, _ := NewLRU[string](3)
				return cache
			},
			key:    "nonexistent",
			wantOk: false,
			check: func(t *testing.T, cache *LRU[string]) {
				assert.Equal(t, int64(0), cache.Len())
			},
		},
		{
			name: "delete existing key",
			setup: func() *LRU[string] {
				cache, _ := NewLRU[string](3)
				cache.Put("key1", "value1")
				cache.Put("key2", "value2")
				return cache
			},
			key:    "key1",
			wantOk: true,
			check: func(t *testing.T, cache *LRU[string]) {
				_, ok := cache.Get("key1")
				assert.False(t, ok)
				assert.Equal(t, int64(1), cache.Len())
			},
		},
		{
			name: "delete only key",
			setup: func() *LRU[string] {
				cache, _ := NewLRU[string](3)
				cache.Put("key1", "value1")
				return cache
			},
			key:    "key1",
			wantOk: true,
			check: func(t *testing.T, cache *LRU[string]) {
				assert.Equal(t, int64(0), cache.Len())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := tt.setup()
			deleted := cache.Delete(tt.key)
			assert.Equal(t, tt.wantOk, deleted)
			tt.check(t, cache)
		})
	}
}

func TestLRU_Peek(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *LRU[string]
		key     string
		wantVal string
		wantOk  bool
		check   func(t *testing.T, cache *LRU[string])
	}{
		{
			name: "peek with empty key",
			setup: func() *LRU[string] {
				cache, _ := NewLRU[string](2)
				return cache
			},
			key:     "",
			wantVal: "",
			wantOk:  false,
			check:   func(t *testing.T, cache *LRU[string]) {},
		},
		{
			name: "peek nonexistent key",
			setup: func() *LRU[string] {
				cache, _ := NewLRU[string](2)
				return cache
			},
			key:     "nonexistent",
			wantVal: "",
			wantOk:  false,
			check:   func(t *testing.T, cache *LRU[string]) {},
		},
		{
			name: "peek existing key",
			setup: func() *LRU[string] {
				cache, _ := NewLRU[string](2)
				cache.Put("key1", "value1")
				cache.Put("key2", "value2")
				return cache
			},
			key:     "key1",
			wantVal: "value1",
			wantOk:  true,
			check:   func(t *testing.T, cache *LRU[string]) {},
		},
		{
			name: "peek doesn't affect LRU order",
			setup: func() *LRU[string] {
				cache, _ := NewLRU[string](2)
				cache.Put("key1", "value1")
				cache.Put("key2", "value2")
				cache.Peek("key1") // should not move key1 to front
				cache.Put("key3", "value3") // should still evict key1
				return cache
			},
			key:     "key1",
			wantVal: "",
			wantOk:  false,
			check: func(t *testing.T, cache *LRU[string]) {
				// Verify key1 was evicted and key2, key3 remain
				_, ok := cache.Get("key2")
				assert.True(t, ok)
				_, ok = cache.Get("key3")
				assert.True(t, ok)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := tt.setup()
			val, ok := cache.Peek(tt.key)
			assert.Equal(t, tt.wantOk, ok)
			assert.Equal(t, tt.wantVal, val)
			tt.check(t, cache)
		})
	}
}

func TestLRU_Len(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *LRU[string]
		wantSize int64
	}{
		{
			name: "empty cache",
			setup: func() *LRU[string] {
				cache, _ := NewLRU[string](3)
				return cache
			},
			wantSize: 0,
		},
		{
			name: "cache with items",
			setup: func() *LRU[string] {
				cache, _ := NewLRU[string](3)
				cache.Put("key1", "value1")
				cache.Put("key2", "value2")
				return cache
			},
			wantSize: 2,
		},
		{
			name: "cache at capacity",
			setup: func() *LRU[string] {
				cache, _ := NewLRU[string](2)
				cache.Put("key1", "value1")
				cache.Put("key2", "value2")
				return cache
			},
			wantSize: 2,
		},
		{
			name: "cache with eviction",
			setup: func() *LRU[string] {
				cache, _ := NewLRU[string](2)
				cache.Put("key1", "value1")
				cache.Put("key2", "value2")
				cache.Put("key3", "value3") // evicts key1
				return cache
			},
			wantSize: 2,
		},
		{
			name: "cache after deletion",
			setup: func() *LRU[string] {
				cache, _ := NewLRU[string](3)
				cache.Put("key1", "value1")
				cache.Put("key2", "value2")
				cache.Delete("key1")
				return cache
			},
			wantSize: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := tt.setup()
			size := cache.Len()
			assert.Equal(t, tt.wantSize, size)
		})
	}
}
