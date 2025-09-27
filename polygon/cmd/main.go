package main

import (
	"fmt"
	"polygon/polygon"
)

func main() {
	triangle := []polygon.Point{{0, 0}, {4, 0}, {2, 3}}
	poly, err := polygon.NewPolygon(triangle)
	if err != nil {
		fmt.Printf("Ошибка создания полигона: %v\n", err)
		return
	}

	fmt.Printf("Треугольник с вершинами: %v\n", poly.Vertices)
	fmt.Printf("Площадь: %.2f\n", poly.Area())
	fmt.Printf("Периметр: %.2f\n", poly.Perimeter())

	fmt.Println()

	square := []polygon.Point{{0, 0}, {2, 0}, {2, 2}, {0, 2}}
	poly, _ = polygon.NewPolygon(square)

	fmt.Printf("Квадрат с вершинами: %v\n", poly.Vertices)
	fmt.Printf("Площадь: %.2f\n", poly.Area())
	fmt.Printf("Периметр: %.2f\n", poly.Perimeter())

	fmt.Println()
}
