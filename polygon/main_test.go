package main

import (
	"math"
	"testing"
)

func TestPolygonArea(t *testing.T) {
	square := []Point{{0, 0}, {1, 0}, {1, 1}, {0, 1}}
	poly, _ := NewPolygon(square)
	area := poly.Area()
	expected := 1.0
	if math.Abs(area-expected) > 1e-9 {
		t.Errorf("Expected area %f, got %f", expected, area)
	}

	triangle := []Point{{0, 0}, {2, 0}, {1, 2}}
	poly, _ = NewPolygon(triangle)
	area = poly.Area()
	expected = 2.0
	if math.Abs(area-expected) > 1e-9 {
		t.Errorf("Expected area %f, got %f", expected, area)
	}

	rect := []Point{{0, 0}, {3, 0}, {3, 4}, {0, 4}}
	poly, _ = NewPolygon(rect)
	area = poly.Area()
	expected = 12.0
	if math.Abs(area-expected) > 1e-9 {
		t.Errorf("Expected area %f, got %f", expected, area)
	}
}

func TestPolygonPerimeter(t *testing.T) {
	square := []Point{{0, 0}, {1, 0}, {1, 1}, {0, 1}}
	poly, _ := NewPolygon(square)
	perimeter := poly.Perimeter()
	expected := 4.0
	if math.Abs(perimeter-expected) > 1e-9 {
		t.Errorf("Expected perimeter %f, got %f", expected, perimeter)
	}

	triangle := []Point{{0, 0}, {3, 0}, {0, 4}}
	poly, _ = NewPolygon(triangle)
	perimeter = poly.Perimeter()
	expected = 12.0
	if math.Abs(perimeter-expected) > 1e-9 {
		t.Errorf("Expected perimeter %f, got %f", expected, perimeter)
	}
}

func TestDistFunction(t *testing.T) {
	p1 := Point{0, 0}
	p2 := Point{3, 4}
	distance := dist(p1, p2)
	expected := 5.0
	if math.Abs(distance-expected) > 1e-9 {
		t.Errorf("Expected distance %f, got %f", expected, distance)
	}
}
