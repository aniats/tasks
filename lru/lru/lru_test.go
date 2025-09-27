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
			cache, err := NewLRU(tt.capacity)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, cache)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, cache)
			assert.Equal(t, tt.capacity, cache.capacity)
			assert.Equal(t, int64(0), cache.size)
			assert.Equal(t, cache.tail, cache.head.next)
			assert.Equal(t, cache.head, cache.tail.prev)
		})
	}
}

func TestLRU_Get(t *testing.T) {
	cache, _ := NewLRU(2)

	_, err := cache.Get("")
	assert.ErrorIs(t, err, ErrorEmptyKey)

	_, err = cache.Get("nonexistent")
	assert.ErrorIs(t, err, ErrorKeyNotFound)

	cache.Put("key1", "value1")
	cache.Put("key2", "value2")

	val, err := cache.Get("key1")
	assert.NoError(t, err)
	assert.Equal(t, "value1", val)

	cache.Put("key3", "value3")

	_, err = cache.Get("key2")
	assert.ErrorIs(t, err, ErrorKeyNotFound)

	_, err = cache.Get("key1")
	assert.NoError(t, err)
}

func TestLRU_Put(t *testing.T) {
	cache, _ := NewLRU(2)

	_, _, _, err := cache.Put("", "value")
	assert.ErrorIs(t, err, ErrorEmptyKey)

	evictedKey, evictedVal, evicted, err := cache.Put("key1", "value1")
	assert.NoError(t, err)
	assert.False(t, evicted)

	evictedKey, evictedVal, evicted, err = cache.Put("key1", "updated_value1")
	assert.NoError(t, err)
	assert.False(t, evicted)

	val, err := cache.Get("key1")
	assert.NoError(t, err)
	assert.Equal(t, "updated_value1", val)

	cache.Put("key2", "value2")

	evictedKey, evictedVal, evicted, err = cache.Put("key3", "value3")
	assert.NoError(t, err)
	assert.True(t, evicted)
	assert.Equal(t, "key1", evictedKey)
	assert.Equal(t, "updated_value1", evictedVal)

	val2, err := cache.Get("key2")
	assert.NoError(t, err)
	assert.Equal(t, "value2", val2)

	val3, err := cache.Get("key3")
	assert.NoError(t, err)
	assert.Equal(t, "value3", val3)

	_, err = cache.Get("key1")
	assert.ErrorIs(t, err, ErrorKeyNotFound)
}

func TestLRU_Delete(t *testing.T) {
	cache, _ := NewLRU(3)

	err := cache.Delete("")
	assert.ErrorIs(t, err, ErrorEmptyKey)

	err = cache.Delete("nonexistent")
	assert.ErrorIs(t, err, ErrorKeyNotFound)

	cache.Put("key1", "value1")
	cache.Put("key2", "value2")

	err = cache.Delete("key1")
	assert.NoError(t, err)

	_, err = cache.Get("key1")
	assert.ErrorIs(t, err, ErrorKeyNotFound)

	size, _ := cache.Len()
	assert.Equal(t, int64(1), size)
}

func TestLRU_Peek(t *testing.T) {
	cache, _ := NewLRU(2)

	_, err := cache.Peek("")
	assert.ErrorIs(t, err, ErrorEmptyKey)

	_, err = cache.Peek("nonexistent")
	assert.ErrorIs(t, err, ErrorKeyNotFound)

	cache.Put("key1", "value1")
	cache.Put("key2", "value2")

	val, err := cache.Peek("key1")
	assert.NoError(t, err)
	assert.Equal(t, "value1", val)

	cache.Put("key3", "value3")

	_, err = cache.Peek("key1")
	assert.ErrorIs(t, err, ErrorKeyNotFound)
}

func TestLRU_Len(t *testing.T) {
	cache, _ := NewLRU(3)

	size, err := cache.Len()
	assert.NoError(t, err)
	assert.Equal(t, int64(0), size)

	cache.Put("key1", "value1")
	cache.Put("key2", "value2")

	size, err = cache.Len()
	assert.NoError(t, err)
	assert.Equal(t, int64(2), size)
}
