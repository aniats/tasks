package main

import (
	"fmt"
	"polygon/polygon"
)

func main() {
	triangle := []polygon.Point{{0, 0}, {4, 0}, {2, 3}}
	poly, err := polygon.NewPolygon(triangle)
	if err != nil {
		fmt.Printf("Error creating polygon: %v\n", err)
		return
	}

	fmt.Printf("Triangle with vertices: %v\n", poly.Vertices)
	fmt.Printf("Area: %.2f\n", poly.Area())
	fmt.Printf("Perimeter: %.2f\n", poly.Perimeter())

	fmt.Println()

	square := []polygon.Point{{0, 0}, {2, 0}, {2, 2}, {0, 2}}
	poly, _ = polygon.NewPolygon(square)

	fmt.Printf("Square with vertices: %v\n", poly.Vertices)
	fmt.Printf("Area: %.2f\n", poly.Area())
	fmt.Printf("Perimeter: %.2f\n", poly.Perimeter())

	fmt.Println()
}
