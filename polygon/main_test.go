package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPolygonArea(t *testing.T) {
	square := []Point{{0, 0}, {1, 0}, {1, 1}, {0, 1}}
	poly, _ := NewPolygon(square)
	area := poly.Area()
	assert.InDelta(t, 1.0, area, 1e-9)

	triangle := []Point{{0, 0}, {2, 0}, {1, 2}}
	poly, _ = NewPolygon(triangle)
	area = poly.Area()
	assert.InDelta(t, 2.0, area, 1e-9)

	rect := []Point{{0, 0}, {3, 0}, {3, 4}, {0, 4}}
	poly, _ = NewPolygon(rect)
	area = poly.Area()
	assert.InDelta(t, 12.0, area, 1e-9)
}

func TestPolygonPerimeter(t *testing.T) {
	square := []Point{{0, 0}, {1, 0}, {1, 1}, {0, 1}}
	poly, _ := NewPolygon(square)
	perimeter := poly.Perimeter()
	assert.InDelta(t, 4.0, perimeter, 1e-9)

	triangle := []Point{{0, 0}, {3, 0}, {0, 4}}
	poly, _ = NewPolygon(triangle)
	perimeter = poly.Perimeter()
	assert.InDelta(t, 12.0, perimeter, 1e-9)
}

func TestDistFunction(t *testing.T) {
	p1 := Point{0, 0}
	p2 := Point{3, 4}
	distance := dist(p1, p2)
	assert.InDelta(t, 5.0, distance, 1e-9)
}
