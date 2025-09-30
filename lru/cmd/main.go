package main

import (
	"fmt"
	"lru/lru"
)

func main() {
	cache, err := lru.NewLRU[int](2)
	if err != nil {
		println(err.Error())
		return
	}

	cache.Put("a", 1)
	cache.Put("b", 2)

	if v, ok := cache.Get("a"); ok {
		fmt.Printf("a = %v\n", v)
	}

	cache.Put("c", 3)

	if v, ok := cache.Get("c"); ok {
		fmt.Printf("c = %v\n", v)
	}

	cache.Put("a", 100)
	if v, ok := cache.Get("a"); ok {
		fmt.Printf("a = %v (updated)\n", v)
	}

	size := cache.Len()
	fmt.Printf("cache size: %d\n", size)
}
