package polygon

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPolygonArea(t *testing.T) {
	tests := []struct {
		name     string
		points   []Point
		expected float64
	}{
		{
			name:     "unit square",
			points:   []Point{{0, 0}, {1, 0}, {1, 1}, {0, 1}},
			expected: 1.0,
		},
		{
			name:     "triangle",
			points:   []Point{{0, 0}, {2, 0}, {1, 2}},
			expected: 2.0,
		},
		{
			name:     "rectangle",
			points:   []Point{{0, 0}, {3, 0}, {3, 4}, {0, 4}},
			expected: 12.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			poly, err := New(tt.points)
			require.NoError(t, err)
			area := poly.Area()
			require.InDelta(t, tt.expected, area, 1e-9)
		})
	}
}

func TestPolygonPerimeter(t *testing.T) {
	tests := []struct {
		name     string
		points   []Point
		expected float64
	}{
		{
			name:     "unit square",
			points:   []Point{{0, 0}, {1, 0}, {1, 1}, {0, 1}},
			expected: 4.0,
		},
		{
			name:     "3-4-5 triangle",
			points:   []Point{{0, 0}, {3, 0}, {0, 4}},
			expected: 12.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			poly, err := New(tt.points)
			require.NoError(t, err)
			perimeter := poly.Perimeter()
			require.InDelta(t, tt.expected, perimeter, 1e-9)
		})
	}
}

func TestDistFunction(t *testing.T) {
	tests := []struct {
		name     string
		p1       Point
		p2       Point
		expected float64
	}{
		{
			name:     "3-4-5 triangle distance",
			p1:       Point{0, 0},
			p2:       Point{3, 4},
			expected: 5.0,
		},
		{
			name:     "same point",
			p1:       Point{1, 1},
			p2:       Point{1, 1},
			expected: 0.0,
		},
		{
			name:     "horizontal distance",
			p1:       Point{0, 0},
			p2:       Point{5, 0},
			expected: 5.0,
		},
		{
			name:     "vertical distance",
			p1:       Point{0, 0},
			p2:       Point{0, 3},
			expected: 3.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			distance := dist(tt.p1, tt.p2)
			require.InDelta(t, tt.expected, distance, 1e-9)
		})
	}
}

func TestPolygonAreaEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		polygon  Polygon
		expected float64
	}{
		{
			name:     "empty polygon",
			polygon:  Polygon{},
			expected: 0.0,
		},
		{
			name:     "single point",
			polygon:  Polygon{{0, 0}},
			expected: 0.0,
		},
		{
			name:     "two points",
			polygon:  Polygon{{0, 0}, {1, 1}},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			area := tt.polygon.Area()
			require.Equal(t, tt.expected, area)
		})
	}
}

func TestPolygonPerimeterEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		polygon  Polygon
		expected float64
	}{
		{
			name:     "empty polygon",
			polygon:  Polygon{},
			expected: 0.0,
		},
		{
			name:     "single point",
			polygon:  Polygon{{0, 0}},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			perimeter := tt.polygon.Perimeter()
			require.Equal(t, tt.expected, perimeter)
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		points    []Point
		shouldErr bool
	}{
		{
			name:      "valid triangle",
			points:    []Point{{0, 0}, {1, 0}, {0, 1}},
			shouldErr: false,
		},
		{
			name:      "valid square",
			points:    []Point{{0, 0}, {1, 0}, {1, 1}, {0, 1}},
			shouldErr: false,
		},
		{
			name:      "insufficient points - empty",
			points:    []Point{},
			shouldErr: true,
		},
		{
			name:      "insufficient points - one point",
			points:    []Point{{0, 0}},
			shouldErr: true,
		},
		{
			name:      "insufficient points - two points",
			points:    []Point{{0, 0}, {1, 0}},
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			poly, err := New(tt.points)

			if tt.shouldErr {
				require.Error(t, err)
				require.Nil(t, poly)
			} else {
				require.NoError(t, err)
				require.NotNil(t, poly)
				require.Equal(t, len(tt.points), len(poly))
			}
		})
	}
}
