package pi

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Result struct {
	Sum float64
}

type Calculator struct {
	numWorkers int
	stopChan   chan struct{}
	results    chan Result
	wg         *sync.WaitGroup
}

func NewCalculator(numWorkers int) (*Calculator, error) {
	if numWorkers <= 0 {
		return nil, fmt.Errorf("number of workers must be positive, got %d", numWorkers)
	}

	return &Calculator{
		numWorkers: numWorkers,
		stopChan:   make(chan struct{}),
		results:    make(chan Result, numWorkers),
		wg:         &sync.WaitGroup{},
	}, nil
}

func (c *Calculator) worker(id int) {
	defer c.wg.Done()

	var sum float64

	fmt.Printf("Worker %d has started\n", id)

	for n := int64(id); ; n += int64(c.numWorkers) {
		select {
		case <-c.stopChan:
			fmt.Printf("Worker %d has received stop signal\n", id)
			c.results <- Result{Sum: sum}
			return
		default:
			// π/4 = 1 - 1/3 + 1/5 - 1/7 + 1/9 - ...
			// Formula: (-1)^n / (2*n + 1)
			denominator := 2*n + 1
			term := 1.0 / float64(denominator)

			if n%2 == 1 {
				term = -term
			}

			sum += term
		}
	}
}

func (c *Calculator) Start() {
	for i := 0; i < c.numWorkers; i++ {
		c.wg.Add(1)
		go c.worker(i)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nStop signal received")
		close(c.stopChan)
	}()

	go func() {
		c.wg.Wait()
		close(c.results)
	}()
}

func (c *Calculator) Calculate() float64 {
	c.Start()

	var totalSum float64
	for result := range c.results {
		totalSum += result.Sum
	}

	return totalSum * 4
}
