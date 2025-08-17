package main

import (
	"errors"
	"fmt"
	"math"
)

type Point struct {
	X float64
	Y float64
}

type Polygon struct {
	Vertices []Point
}

func NewPolygon(pts []Point) (Polygon, error) {
	if len(pts) < 3 {
		return Polygon{}, errors.New("polygon needs at least 3 points")
	}

	cp := make([]Point, len(pts))
	copy(cp, pts)
	return Polygon{Vertices: cp}, nil
}

func (pg *Polygon) Area() float64 {
	n := len(pg.Vertices)

	if n < 3 {
		return 0
	}

	sum := 0.0

	for i := 0; i < n; i++ {
		j := (i + 1) % n
		x1, y1 := pg.Vertices[i].X, pg.Vertices[i].Y
		x2, y2 := pg.Vertices[j].X, pg.Vertices[j].Y
		sum += x1*y2 - x2*y1
	}

	return math.Abs(sum) / 2
}

func (pg *Polygon) Perimeter() float64 {
	n := len(pg.Vertices)

	if n < 2 {
		return 0
	}

	per := 0.0

	for i := 0; i < n; i++ {
		j := (i + 1) % n
		per += dist(pg.Vertices[i], pg.Vertices[j])
	}

	return per
}

func dist(a, b Point) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return math.Hypot(dx, dy)
}

func main() {
	triangle := []Point{{0, 0}, {4, 0}, {2, 3}}
	poly, err := NewPolygon(triangle)
	if err != nil {
		fmt.Printf("Ошибка создания полигона: %v\n", err)
		return
	}

	fmt.Printf("Треугольник с вершинами: %v\n", poly.Vertices)
	fmt.Printf("Площадь: %.2f\n", poly.Area())
	fmt.Printf("Периметр: %.2f\n", poly.Perimeter())

	fmt.Println()

	square := []Point{{0, 0}, {2, 0}, {2, 2}, {0, 2}}
	poly, _ = NewPolygon(square)

	fmt.Printf("Квадрат с вершинами: %v\n", poly.Vertices)
	fmt.Printf("Площадь: %.2f\n", poly.Area())
	fmt.Printf("Периметр: %.2f\n", poly.Perimeter())

	fmt.Println()
}
