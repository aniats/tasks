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
			assert.Equal(t, cache.list.tail, cache.list.head.next)
			assert.Equal(t, cache.list.head, cache.list.tail.prev)
		})
	}
}

func TestLRU_Get(t *testing.T) {
	cache, _ := NewLRU[string](2)

	_, ok := cache.Get("")
	assert.False(t, ok)

	_, ok = cache.Get("nonexistent")
	assert.False(t, ok)

	cache.Put("key1", "value1")
	cache.Put("key2", "value2")

	val, ok := cache.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, "value1", val)

	cache.Put("key3", "value3")

	_, ok = cache.Get("key2")
	assert.False(t, ok)

	_, ok = cache.Get("key1")
	assert.True(t, ok)
}

func TestLRU_Put(t *testing.T) {
	cache, _ := NewLRU[string](2)

	err := cache.Put("", "value")
	assert.ErrorIs(t, err, ErrorEmptyKey)

	err = cache.Put("key1", "value1")
	assert.NoError(t, err)

	err = cache.Put("key1", "updated_value1")
	assert.NoError(t, err)

	val, ok := cache.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, "updated_value1", val)

	cache.Put("key2", "value2")

	err = cache.Put("key3", "value3")
	assert.NoError(t, err)

	val2, ok := cache.Get("key2")
	assert.True(t, ok)
	assert.Equal(t, "value2", val2)

	val3, ok := cache.Get("key3")
	assert.True(t, ok)
	assert.Equal(t, "value3", val3)

	_, ok = cache.Get("key1")
	assert.False(t, ok)
}

func TestLRU_Delete(t *testing.T) {
	cache, _ := NewLRU[string](3)

	deleted := cache.Delete("")
	assert.False(t, deleted)

	deleted = cache.Delete("nonexistent")
	assert.False(t, deleted)

	cache.Put("key1", "value1")
	cache.Put("key2", "value2")

	deleted = cache.Delete("key1")
	assert.True(t, deleted)

	_, ok := cache.Get("key1")
	assert.False(t, ok)

	size := cache.Len()
	assert.Equal(t, int64(1), size)
}

func TestLRU_Peek(t *testing.T) {
	cache, _ := NewLRU[string](2)

	_, ok := cache.Peek("")
	assert.False(t, ok)

	_, ok = cache.Peek("nonexistent")
	assert.False(t, ok)

	cache.Put("key1", "value1")
	cache.Put("key2", "value2")

	val, ok := cache.Peek("key1")
	assert.True(t, ok)
	assert.Equal(t, "value1", val)

	cache.Put("key3", "value3")

	_, ok = cache.Peek("key1")
	assert.False(t, ok)
}

func TestLRU_Len(t *testing.T) {
	cache, _ := NewLRU[string](3)

	size := cache.Len()
	assert.Equal(t, int64(0), size)

	cache.Put("key1", "value1")
	cache.Put("key2", "value2")

	size = cache.Len()
	assert.Equal(t, int64(2), size)
}
