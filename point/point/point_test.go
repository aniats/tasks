package point

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPoint(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		lat       float64
		lng       float64
		shouldErr bool
	}{
		{"Valid coordinates", 55.7558, 37.6176, false},
		{"Zero coordinates", 0, 0, false},
		{"Boundary coordinates", 90, -180, false},
		{"Invalid latitude too high", 91, 0, true},
		{"Invalid latitude too low", -91, 0, true},
		{"Invalid longitude too high", 0, 181, true},
		{"Invalid longitude too low", 0, -181, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			point, err := NewPoint(tt.lat, tt.lng)

			if tt.shouldErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.lat, point.Latitude)
			assert.Equal(t, tt.lng, point.Longitude)
		})
	}
}

func TestValidateCoordinates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		lat       float64
		lng       float64
		shouldErr bool
	}{
		{"Valid coordinates", 55.7558, 37.6176, false},
		{"Boundary valid", 90, 180, false},
		{"Boundary valid negative", -90, -180, false},
		{"Invalid lat high", 90.1, 0, true},
		{"Invalid lat low", -90.1, 0, true},
		{"Invalid lng high", 0, 180.1, true},
		{"Invalid lng low", 0, -180.1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validateCoordinates(tt.lat, tt.lng)

			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPointString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		point    Point
		expected string
	}{
		{
			"Moscow coordinates",
			Point{Latitude: 55.7558, Longitude: 37.6176},
			"Point(lat=55.755800, lng=37.617600)",
		},
		{
			"Zero coordinates",
			Point{Latitude: 0, Longitude: 0},
			"Point(lat=0.000000, lng=0.000000)",
		},
		{
			"Negative coordinates",
			Point{Latitude: -33.8688, Longitude: -151.2093},
			"Point(lat=-33.868800, lng=-151.209300)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := tt.point.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDistanceTo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		p1        Point
		p2        Point
		expected  float64
		tolerance float64
	}{
		{
			"Moscow to Saint Petersburg",
			Point{Latitude: 55.7558, Longitude: 37.6176},
			Point{Latitude: 59.9311, Longitude: 30.3609},
			635.0,
			10.0,
		},
		{
			"Same point",
			Point{Latitude: 55.7558, Longitude: 37.6176},
			Point{Latitude: 55.7558, Longitude: 37.6176},
			0.0,
			0.001,
		},
		{
			"Equator points",
			Point{Latitude: 0, Longitude: 0},
			Point{Latitude: 0, Longitude: 1},
			111.32,
			1.0,
		},
		{
			"Example from task: (2,0) to (0,0)",
			Point{Latitude: 2, Longitude: 0},
			Point{Latitude: 0, Longitude: 0},
			222.64,
			5.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			distance := tt.p1.DistanceTo(tt.p2)
			assert.InDelta(t, tt.expected, distance, tt.tolerance)

			reverseDistance := tt.p2.DistanceTo(tt.p1)
			assert.InDelta(t, distance, reverseDistance, 0.001)
		})
	}
}

func TestIsWithinRadius(t *testing.T) {
	t.Parallel()

	moscow := Point{Latitude: 55.7558, Longitude: 37.6176}
	stPetersburg := Point{Latitude: 59.9311, Longitude: 30.3609}

	tests := []struct {
		name     string
		point    Point
		center   Point
		radius   float64
		expected bool
	}{
		{"Point within radius", moscow, moscow, 10.0, true},
		{"Point outside radius", stPetersburg, moscow, 100.0, false},
		{"Point exactly on radius boundary", stPetersburg, moscow, 635.0, true},
		{"Negative radius", moscow, moscow, -10.0, false},
		{"Zero radius same point", moscow, moscow, 0.0, true},
		{"Zero radius different points", stPetersburg, moscow, 0.0, false},
		{
			"Example case: point (2,0) within radius 300km of (0,0)",
			Point{Latitude: 2, Longitude: 0},
			Point{Latitude: 0, Longitude: 0},
			300.0,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := tt.point.IsWithinRadius(tt.center, tt.radius)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHaversineDistance(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		p1        Point
		p2        Point
		expected  float64
		tolerance float64
	}{
		{
			"Known distance calculation",
			Point{Latitude: 52.5200, Longitude: 13.4050},
			Point{Latitude: 48.8566, Longitude: 2.3522},
			878.0,
			20.0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			distance := haversineDistance(tt.p1, tt.p2)
			assert.InDelta(t, tt.expected, distance, tt.tolerance)
		})
	}
}

func TestToRadians(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		degrees  float64
		expected float64
	}{
		{"Zero degrees", 0, 0},
		{"90 degrees", 90, math.Pi / 2},
		{"180 degrees", 180, math.Pi},
		{"360 degrees", 360, 2 * math.Pi},
		{"Negative degrees", -90, -math.Pi / 2},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := toRadians(tt.degrees)
			assert.InDelta(t, tt.expected, result, 0.000001)
		})
	}
}

func TestParsePoint(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     string
		expected  *Point
		shouldErr bool
	}{
		{"Valid point", "55.7558,37.6176", &Point{Latitude: 55.7558, Longitude: 37.6176}, false},
		{"Valid with spaces", " 55.7558 , 37.6176 ", &Point{Latitude: 55.7558, Longitude: 37.6176}, false},
		{"Invalid format missing comma", "55.7558 37.6176", nil, true},
		{"Invalid format too many parts", "55.7558,37.6176,10", nil, true},
		{"Invalid latitude", "invalid,37.6176", nil, true},
		{"Invalid longitude", "55.7558,invalid", nil, true},
		{"Out of range latitude", "91.0,37.6176", nil, true},
		{"Out of range longitude", "55.7558,181.0", nil, true},
		{"Example coordinates", "2,0", &Point{Latitude: 2, Longitude: 0}, false},
		{"Example coordinates 2", "0,0", &Point{Latitude: 0, Longitude: 0}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := parsePoint(tt.input)

			if tt.shouldErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected.Latitude, result.Latitude)
			assert.Equal(t, tt.expected.Longitude, result.Longitude)
		})
	}
}

func TestParsePoints(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     []string
		expected  int
		shouldErr bool
	}{
		{
			"Valid multiple points",
			[]string{"55.7558,37.6176", "59.9311,30.3609"},
			2,
			false,
		},
		{
			"Single point",
			[]string{"0,0"},
			1,
			false,
		},
		{
			"Empty slice",
			[]string{},
			0,
			false,
		},
		{
			"Invalid point in slice",
			[]string{"55.7558,37.6176", "invalid"},
			0,
			true,
		},
		{
			"Example points from task",
			[]string{"2,0", "0,0"},
			2,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := parsePoints(tt.input)

			if tt.shouldErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, result, tt.expected)
		})
	}
}
