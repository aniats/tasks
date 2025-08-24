package main

import (
	"fmt"
	"lru/lru"
)

func main() {
	cache, err := lru.NewLRU(2)
	if err != nil {
		println(err.Error())
		return
	}

	cache.Put("a", 1)
	cache.Put("b", 2)

	if v, err := cache.Get("a"); err == nil {
		fmt.Printf("a = %v\n", v)
	}

	evictedKey, evictedVal, evicted, _ := cache.Put("c", 3)
	if evicted {
		fmt.Printf("evicted: %s = %v\n", evictedKey, evictedVal)
	}

	if v, err := cache.Get("c"); err == nil {
		fmt.Printf("c = %v\n", v)
	}

	cache.Put("a", 100)
	if v, err := cache.Get("a"); err == nil {
		fmt.Printf("a = %v (updated)\n", v)
	}

	if size, err := cache.Len(); err == nil {
		fmt.Printf("cache size: %d\n", size)
	}
}
