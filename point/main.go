package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
)

const earthRadiusKm = 6371

type Point struct {
	Longitude float64
	Latitude  float64
}

func NewPoint(latitude, longitude float64) (*Point, error) {
	if err := validateCoordinates(latitude, longitude); err != nil {
		return nil, err
	}
	return &Point{
		Latitude:  latitude,
		Longitude: longitude,
	}, nil
}

func validateCoordinates(lat, lng float64) error {
	if lat < -90 || lat > 90 {
		return fmt.Errorf("latitude must be between -90 and 90, got: %f", lat)
	}
	if lng < -180 || lng > 180 {
		return fmt.Errorf("longitude must be between -180 and 180, got: %f", lng)
	}
	return nil
}

func (p Point) String() string {
	return fmt.Sprintf("Point(lat=%.6f, lng=%.6f)", p.Latitude, p.Longitude)
}

func (p Point) DistanceTo(other Point) float64 {
	return haversineDistance(p, other)
}

func (p Point) IsWithinRadius(center Point, radiusKm float64) bool {
	if radiusKm < 0 {
		return false
	}
	return p.DistanceTo(center) <= radiusKm
}

func haversineDistance(p1, p2 Point) float64 {
	lat1Rad := toRadians(p1.Latitude)
	lat2Rad := toRadians(p2.Latitude)
	deltaLatRad := toRadians(p2.Latitude - p1.Latitude)
	deltaLngRad := toRadians(p2.Longitude - p1.Longitude)

	a := math.Sin(deltaLatRad/2)*math.Sin(deltaLatRad/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLngRad/2)*math.Sin(deltaLngRad/2)

	return 2 * earthRadiusKm * math.Asin(math.Sqrt(a))
}

func toRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func parsePoint(s string) (*Point, error) {
	parts := strings.Split(s, ",")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid point format: %s, expected lat,lng", s)
	}

	lat, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return nil, fmt.Errorf("invalid latitude: %s", parts[0])
	}

	lng, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return nil, fmt.Errorf("invalid longitude: %s", parts[1])
	}

	return NewPoint(lat, lng)
}

func parsePoints(pointStrings []string) ([]*Point, error) {
	var points []*Point
	for _, pointStr := range pointStrings {
		point, err := parsePoint(pointStr)
		if err != nil {
			return nil, fmt.Errorf("error parsing point %s: %w", pointStr, err)
		}
		points = append(points, point)
	}
	return points, nil
}

func printDistances(points []*Point) {
	fmt.Printf("\nDistances between points:\n")

	if len(points) < 2 {
		fmt.Printf("Need at least 2 points to calculate distance\n")
		return
	}

	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			dist := points[i].DistanceTo(*points[j])
			fmt.Printf("  %s to %s: %.2f km\n",
				points[i], points[j], dist)
		}
	}
}

func printRadiusCheck(points []*Point, centerPointStr string, radius float64) error {
	center, err := parsePoint(centerPointStr)
	if err != nil {
		return fmt.Errorf("error parsing center point %s: %w", centerPointStr, err)
	}

	fmt.Printf("\nRadius check (center: %s, radius: %.2f km):\n",
		center, radius)

	for i, point := range points {
		within := point.IsWithinRadius(*center, radius)
		distance := point.DistanceTo(*center)
		fmt.Printf("  Point %d: %s - Distance: %.2f km, Within radius: %v\n",
			i+1, point, distance, within)
	}

	return nil
}

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

	points, err := parsePoints(pointStrings)
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
