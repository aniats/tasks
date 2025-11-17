package workerpool

import (
	"context"
	"errors"
	"log"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	if len(tasks) == 0 {
		return nil
	}

	if n <= 0 {
		n = 1
	}

	if m <= 0 {
		return runWithoutErrorLimit(tasks, n)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	taskChan := make(chan Task, len(tasks))

	var (
		wg         sync.WaitGroup
		errorCount int64
	)

	for _, task := range tasks {
		taskChan <- task
	}
	close(taskChan)

	for i := 0; i < n; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case task, ok := <-taskChan:
					if !ok {
						return
					}

					if ctx.Err() != nil {
						return
					}

					if err := task(); err != nil {
						log.Printf("Task error: %v", err)
						count := atomic.AddInt64(&errorCount, 1)
						if count >= int64(m) {
							cancel()
							return
						}
					}
				}
			}
		}()
	}

	wg.Wait()

	if atomic.LoadInt64(&errorCount) >= int64(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func runWithoutErrorLimit(tasks []Task, n int) error {
	if len(tasks) == 0 {
		return nil
	}

	taskChan := make(chan Task, len(tasks))
	var wg sync.WaitGroup

	for _, task := range tasks {
		taskChan <- task
	}
	close(taskChan)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskChan {
				if err := task(); err != nil {
					log.Printf("Task error (ignored): %v", err)
				}
			}
		}()
	}

	wg.Wait()
	return nil
}
