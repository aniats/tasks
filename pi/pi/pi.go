package pi

import (
	"fmt"
	"sync"
)

// CalculateLeibnizTerm calculates a single term in the Leibniz series for π/4
// π/4 = 1 - 1/3 + 1/5 - 1/7 + 1/9 - ...
// Formula: (-1)^n / (2*n + 1)
func CalculateLeibnizTerm(n int) float64 {
	denominator := 2*n + 1
	term := 1.0 / float64(denominator)

	if n%2 == 1 {
		term = -term
	}

	return term
}

type Calculator struct {
	numWorkers int
	stopChan   chan struct{}
	results    chan float64
	wg         *sync.WaitGroup
}

func New(numWorkers int) (*Calculator, error) {
	if numWorkers <= 0 {
		return nil, fmt.Errorf("number of workers must be positive, got %d", numWorkers)
	}

	return &Calculator{
		numWorkers: numWorkers,
		stopChan:   make(chan struct{}),
		results:    make(chan float64, numWorkers),
		wg:         &sync.WaitGroup{},
	}, nil
}

func (c *Calculator) worker(id int) {
	defer c.wg.Done()

	var sum float64

	fmt.Printf("Worker %d has started\n", id)

	for n := id; ; n += c.numWorkers {
		select {
		case <-c.stopChan:
			fmt.Printf("Worker %d has received stop signal\n", id)
			c.results <- sum
			return
		default:
			term := CalculateLeibnizTerm(n)
			sum += term
		}
	}
}

func (c *Calculator) Start() {
	for i := 0; i < c.numWorkers; i++ {
		c.wg.Add(1)
		go c.worker(i)
	}

	go func() {
		c.wg.Wait()
		close(c.results)
	}()
}

func (c *Calculator) Stop() {
	close(c.stopChan)
}

func (c *Calculator) Calculate() float64 {
	c.Start()

	var totalSum float64
	for result := range c.results {
		totalSum += result
	}

	return totalSum * 4
}
