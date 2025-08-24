package main

import (
	"fmt"
	"lru/lru"
)

func main() {
	lru, _ := lru.NewLRU(2)

	lru.Put("a", 1)
	lru.Put("b", 2)

	if v, ok := lru.Get("a"); !ok || v.(int) != 1 {
		panic("want a=1")
	}

	// Access order now: a (MRU), b (LRU)
	lru.Put("c", 3) // should evict b

	if _, ok := lru.Get("b"); ok {
		panic("b should have been evicted")
	}
	if v, ok := lru.Get("c"); !ok || v.(int) != 3 {
		panic("want c=3")
	}
	if v, ok := lru.Get("a"); !ok || v.(int) != 1 {
		panic("want a=1")
	}

	// Update existing should NOT change size, only move to MRU
	evk, _, ev := lru.Put("a", 100)
	if ev {
		panic("updating existing must not evict")
	}
	if v, _ := lru.Get("a"); v.(int) != 100 {
		panic("want a=100")
	}
	fmt.Printf("evk=%+v\n", evk)
}
