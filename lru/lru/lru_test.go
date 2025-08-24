package lru

import (
	"errors"
	"testing"
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
				if err == nil || !errors.Is(err, tt.wantErr) {
					t.Errorf("NewLRU() error = %v, wantErr %v", err, tt.wantErr)
				}
				if cache != nil {
					t.Errorf("NewLRU() should return nil cache on error")
				}
				return
			}

			if err != nil {
				t.Errorf("NewLRU() unexpected error = %v", err)
				return
			}

			if cache == nil {
				t.Errorf("NewLRU() returned nil cache")
				return
			}

			if cache.capacity != tt.capacity {
				t.Errorf("NewLRU() capacity = %v, want %v", cache.capacity, tt.capacity)
			}

			if cache.size != 0 {
				t.Errorf("NewLRU() size = %v, want 0", cache.size)
			}

			if cache.head.next != cache.tail {
				t.Errorf("NewLRU() head.next should point to tail")
			}

			if cache.tail.prev != cache.head {
				t.Errorf("NewLRU() tail.prev should point to head")
			}
		})
	}
}

func TestLRU_Get(t *testing.T) {
	cache, _ := NewLRU(2)

	_, err := cache.Get("")
	if !errors.Is(err, ErrorEmptyKey) {
		t.Errorf("Get(\"\") error = %v, want %v", err, ErrorEmptyKey)
	}

	_, err = cache.Get("nonexistent")
	if !errors.Is(err, ErrorKeyNotFound) {
		t.Errorf("Get(nonexistent) error = %v, want %v", err, ErrorKeyNotFound)
	}

	cache.Put("key1", "value1")
	cache.Put("key2", "value2")

	val, err := cache.Get("key1")
	if err != nil {
		t.Errorf("Get(key1) unexpected error = %v", err)
	}
	if val != "value1" {
		t.Errorf("Get(key1) = %v, want value1", val)
	}

	cache.Put("key3", "value3")

	_, err = cache.Get("key2")
	if !errors.Is(err, ErrorKeyNotFound) {
		t.Errorf("key2 should have been evicted")
	}

	_, err = cache.Get("key1")
	if err != nil {
		t.Errorf("key1 should still exist after being accessed")
	}
}

func TestLRU_Put(t *testing.T) {
	cache, _ := NewLRU(2)

	_, _, _, err := cache.Put("", "value")
	if !errors.Is(err, ErrorEmptyKey) {
		t.Errorf("Put(\"\") error = %v, want %v", err, ErrorEmptyKey)
	}

	evictedKey, evictedVal, evicted, err := cache.Put("key1", "value1")
	if err != nil {
		t.Errorf("Put(key1) unexpected error = %v", err)
	}
	if evicted {
		t.Errorf("Put(key1) should not evict on non-full cache")
	}

	evictedKey, evictedVal, evicted, err = cache.Put("key1", "updated_value1")
	if err != nil {
		t.Errorf("Put(key1 update) unexpected error = %v", err)
	}
	if evicted {
		t.Errorf("Put(key1 update) should not evict")
	}

	val, err := cache.Get("key1")
	if err != nil || val != "updated_value1" {
		t.Errorf("Key1 should be updated to 'updated_value1', got %v", val)
	}

	cache.Put("key2", "value2")

	evictedKey, evictedVal, evicted, err = cache.Put("key3", "value3")
	if err != nil {
		t.Errorf("Put(key3) unexpected error = %v", err)
	}
	if !evicted {
		t.Errorf("Put(key3) should evict oldest item")
	}
	if evictedKey != "key1" || evictedVal != "updated_value1" {
		t.Errorf("Put(key3) evicted wrong item: key=%v, val=%v, expected key1/updated_value1", evictedKey, evictedVal)
	}

	val2, err := cache.Get("key2")
	if err != nil || val2 != "value2" {
		t.Errorf("key2 should still exist with value2, got %v, err=%v", val2, err)
	}

	val3, err := cache.Get("key3")
	if err != nil || val3 != "value3" {
		t.Errorf("key3 should exist with value3, got %v, err=%v", val3, err)
	}

	_, err = cache.Get("key1")
	if !errors.Is(err, ErrorKeyNotFound) {
		t.Errorf("key1 should have been evicted")
	}
}

func TestLRU_Delete(t *testing.T) {
	cache, _ := NewLRU(3)

	err := cache.Delete("")
	if !errors.Is(err, ErrorEmptyKey) {
		t.Errorf("Delete(\"\") error = %v, want %v", err, ErrorEmptyKey)
	}

	err = cache.Delete("nonexistent")
	if !errors.Is(err, ErrorKeyNotFound) {
		t.Errorf("Delete(nonexistent) error = %v, want %v", err, ErrorKeyNotFound)
	}

	cache.Put("key1", "value1")
	cache.Put("key2", "value2")

	err = cache.Delete("key1")
	if err != nil {
		t.Errorf("Delete(key1) unexpected error = %v", err)
	}

	_, err = cache.Get("key1")
	if !errors.Is(err, ErrorKeyNotFound) {
		t.Errorf("key1 should be deleted")
	}

	size, _ := cache.Len()
	if size != 1 {
		t.Errorf("Size should be 1 after deletion, got %v", size)
	}
}

func TestLRU_Peek(t *testing.T) {
	cache, _ := NewLRU(2)

	_, err := cache.Peek("")
	if !errors.Is(err, ErrorEmptyKey) {
		t.Errorf("Peek(\"\") error = %v, want %v", err, ErrorEmptyKey)
	}

	_, err = cache.Peek("nonexistent")
	if !errors.Is(err, ErrorKeyNotFound) {
		t.Errorf("Peek(nonexistent) error = %v, want %v", err, ErrorKeyNotFound)
	}

	cache.Put("key1", "value1")
	cache.Put("key2", "value2")

	val, err := cache.Peek("key1")
	if err != nil {
		t.Errorf("Peek(key1) unexpected error = %v", err)
	}
	if val != "value1" {
		t.Errorf("Peek(key1) = %v, want value1", val)
	}

	cache.Put("key3", "value3")

	_, err = cache.Peek("key1")
	if !errors.Is(err, ErrorKeyNotFound) {
		t.Errorf("key1 should have been evicted after Peek")
	}
}

func TestLRU_Len(t *testing.T) {
	cache, _ := NewLRU(3)

	size, err := cache.Len()
	if err != nil {
		t.Errorf("Len() unexpected error = %v", err)
	}
	if size != 0 {
		t.Errorf("Len() = %v, want 0", size)
	}

	cache.Put("key1", "value1")
	cache.Put("key2", "value2")

	size, err = cache.Len()
	if err != nil {
		t.Errorf("Len() unexpected error = %v", err)
	}
	if size != 2 {
		t.Errorf("Len() = %v, want 2", size)
	}
}
