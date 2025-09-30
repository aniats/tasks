package main

import (
	"fmt"
	"semaphore/channel"
	"semaphore/mutex"
)

func main() {
	channelSemaphore()
	fmt.Println()
	mutexSemaphore()
	fmt.Println()
}

func channelSemaphore() {
	sem := channel.NewChannelSemaphore(3)

	if sem.TryAcquire(2) {
		fmt.Println("Acquired 2 permits successfully")
	}

	if !sem.TryAcquire(2) {
		fmt.Println("Failed to acquire 2 more permits as expected")
	}

	sem.Release(1)
	fmt.Println("Released 1 permit")

	if sem.TryAcquire(1) {
		fmt.Println("Acquired 1 permit after release")
	}
}

func mutexSemaphore() {
	sem := mutex.NewMutexSemaphore(3)

	if sem.TryAcquire(2) {
		fmt.Println("Acquired 2 permits successfully")
	}

	if !sem.TryAcquire(2) {
		fmt.Println("Failed to acquire 2 more permits as expected")
	}

	sem.Release(1)
	fmt.Println("Released 1 permit")

	if sem.TryAcquire(1) {
		fmt.Println("Acquired 1 permit after release")
	}
}
