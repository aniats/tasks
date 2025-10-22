package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"pi/pi"
	"syscall"
)

func Calculate(numWorkers int) (float64, error) {
	calc, err := pi.New(numWorkers)
	if err != nil {
		return 0, err
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	resultChan := make(chan float64)
	go func() {
		result := calc.Calculate()
		resultChan <- result
	}()
	
	go func() {
		<-sigChan
		fmt.Println("\nStop signal received")
		calc.Stop()
	}()

	return <-resultChan, nil
}

func main() {
	numWorkers := flag.Int("n", 4, "Number of goroutines")
	flag.Parse()

	piValue, err := Calculate(*numWorkers)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Approximately pi is:", piValue)
}
