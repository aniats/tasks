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
		printDistances(points)
	}

	if radius > 0 && centerPoint != "" {
		if err := printRadiusCheck(points, centerPoint, radius); err != nil {
			log.Fatalf("Error: %v", err)
		}
	}
}

func printDistances(points []point.Point) {
	fmt.Printf("\nDistances between points:\n")

	if len(points) < 2 {
		fmt.Printf("Need at least 2 points to calculate distance\n")
		return
	}

	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			dist := point.HaversineDistance(points[i], points[j])
			fmt.Printf("  %s to %s: %.2f km\n",
				points[i], points[j], dist)
		}
	}
}

func printRadiusCheck(points []point.Point, centerPointStr string, radius float64) error {
	center, err := point.ParsePoint(centerPointStr)
	if err != nil {
		return fmt.Errorf("error parsing center point %s: %w", centerPointStr, err)
	}

	fmt.Printf("\nRadius check (center: %s, radius: %.2f km):\n",
		center, radius)

	for i, pt := range points {
		within := pt.IsWithinRadius(center, radius)
		distance := point.HaversineDistance(pt, center)
		fmt.Printf("  Point %d: %s - Distance: %.2f km, Within radius: %v\n",
			i+1, pt, distance, within)
	}

	return nil
}
