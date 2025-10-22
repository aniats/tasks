package point

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

const earthRadiusKm = 6371

type Point struct {
	Longitude float64
	Latitude  float64
}

func New(latitude, longitude float64) (Point, error) {
	if err := validateCoordinates(latitude, longitude); err != nil {
		return Point{}, err
	}
	return Point{
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

func (p Point) IsWithinRadius(center Point, radiusKm float64) bool {
	if radiusKm < 0 {
		return false
	}
	return HaversineDistance(p, center) <= radiusKm
}

func HaversineDistance(p1, p2 Point) float64 {
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

func ParsePoint(s string) (Point, error) {
	parts := strings.Split(s, ",")
	if len(parts) != 2 {
		return Point{}, fmt.Errorf("invalid point format: %s, expected lat,lng", s)
	}

	lat, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return Point{}, fmt.Errorf("invalid latitude: %s", parts[0])
	}

	lng, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return Point{}, fmt.Errorf("invalid longitude: %s", parts[1])
	}

	return New(lat, lng)
}

func ParsePoints(pointStrings []string) ([]Point, error) {
	points := make([]Point, 0, len(pointStrings))
	for _, pointStr := range pointStrings {
		point, err := ParsePoint(pointStr)
		if err != nil {
			return nil, fmt.Errorf("error parsing point %s: %w", pointStr, err)
		}
		points = append(points, point)
	}
	return points, nil
}
