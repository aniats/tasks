package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"

	"tasks/workerpool/workerpool"
)

func main() {
	successCase()
	fmt.Println()

	errorLimitCase()
	fmt.Println()

	ignoreErrorCase()
	fmt.Println()
}

func successCase() {
	var successCount int64

	tasks := make([]workerpool.Task, 10)
	for i := range tasks {
		taskID := i + 1
		tasks[i] = func() error {
			atomic.AddInt64(&successCount, 1)
			fmt.Printf("  Task %d completed successfully\n", taskID)
			return nil
		}
	}

	start := time.Now()
	err := workerpool.Run(tasks, 3, 5)
	duration := time.Since(start)

	fmt.Printf("Result: %v\n", err)
	fmt.Printf("Successful tasks: %d\n", atomic.LoadInt64(&successCount))
	fmt.Printf("Duration: %v\n", duration)
}

func errorLimitCase() {
	var taskCount int64

	tasks := make([]workerpool.Task, 15)
	for i := range tasks {
		taskID := i + 1
		tasks[i] = func() error {
			atomic.AddInt64(&taskCount, 1)
			if taskID <= 8 {
				fmt.Printf("  Task %d failed with error\n", taskID)
				return errors.New("simulated error")
			}
			fmt.Printf("  Task %d completed successfully\n", taskID)
			return nil
		}
	}

	start := time.Now()
	err := workerpool.Run(tasks, 2, 3)
	duration := time.Since(start)

	fmt.Printf("Result: %v\n", err)
	fmt.Printf("Total tasks attempted: %d\n", atomic.LoadInt64(&taskCount))
	fmt.Printf("Duration: %v\n", duration)
}

func ignoreErrorCase() {
	var taskCount int64

	tasks := make([]workerpool.Task, 8)
	for i := range tasks {
		taskID := i + 1
		tasks[i] = func() error {
			atomic.AddInt64(&taskCount, 1)
			if rand.Float32() < 0.5 {
				fmt.Printf("  Task %d failed (but ignored)\n", taskID)
				return errors.New("ignored error")
			}
			fmt.Printf("  Task %d succeeded\n", taskID)
			return nil
		}
	}

	start := time.Now()
	err := workerpool.Run(tasks, 3, 0)
	duration := time.Since(start)

	fmt.Printf("Result: %v\n", err)
	fmt.Printf("Total tasks executed: %d\n", atomic.LoadInt64(&taskCount))
	fmt.Printf("Duration: %v\n", duration)
}
