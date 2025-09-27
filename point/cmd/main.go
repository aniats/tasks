package main

import (
	"fmt"
	"log"

	"os"
	"point/point"

	"github.com/spf13/pflag"
)

// go run main.go --point=2,0 --point=0,0 --distance

// go run main.go --point=55.7558,37.6176 --point=59.9311,30.3609 --distance

// go run main.go \
//  --point=55.7558,37.6176 \
//  --point=59.9311,30.3609 \
//  --point=56.8431,60.6454 \
//  --distance

// go run main.go \
//   --point=55.7558,37.6176 \
//   --point=55.9,37.8 \
//   --point=50.0,40.0 \
//   --center=55.7558,37.6176 \
//   --radius=100

// go run main.go --point=90,180 --point=-90,-180 --distance

// Errors:
// go run main.go --point=55.7558 37.6176 --distance
// go run main.go --distance

func main() {
	var pointStrings []string
	var distance bool
	var radius float64
	var centerPoint string

	pflag.StringArrayVar(&pointStrings, "point", []string{},
		"Point coordinates in format lat,lng (can be specified multiple times)")
	pflag.BoolVar(&distance, "distance", false,
		"Calculate distance between points")
	pflag.Float64Var(&radius, "radius", 0,
		"Radius in km for radius check")
	pflag.StringVar(&centerPoint, "center", "",
		"Center point for radius check in format lat,lng")

	pflag.Usage = func() {
		programName := "point"
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", programName)
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s --point=55.7558,37.6176 --point=59.9311,30.3609 --distance\n", programName)
		fmt.Fprintf(os.Stderr, "  %s --point=55.7558,37.6176 --center=55.7558,37.6176 --radius=10\n", programName)
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		pflag.PrintDefaults()
	}

	pflag.Parse()

	if len(pointStrings) == 0 {
		fmt.Fprintf(os.Stderr, "Error: at least one point must be specified\n")
		pflag.Usage()
		os.Exit(1)
	}

	points, err := point.ParsePoints(pointStrings)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	if distance {
		point.PrintDistances(points)
	}

	if radius > 0 && centerPoint != "" {
		if err := point.PrintRadiusCheck(points, centerPoint, radius); err != nil {
			log.Fatalf("Error: %v", err)
		}
	}
}
