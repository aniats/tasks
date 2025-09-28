package main

import (
	"flag"
	"fmt"
	"log"
	"pi/pi"
)

func Calculate(numWorkers int) (float64, error) {
	calc, err := pi.NewCalculator(numWorkers)
	if err != nil {
		return 0, err
	}

	return calc.Calculate(), nil
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
