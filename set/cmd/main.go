package main

import (
	"fmt"
	"set/set"
)

func main() {
	setA := set.New[int]()
	setB := set.New[int]()

	setA.Add(1)
	setA.Add(2)
	setA.Add(3)

	setB.Add(2)
	setB.Add(3)
	setB.Add(4)

	union := setA.Union(setB)
	intersection := setA.Intersection(setB)
	difference := setA.Difference(setB)

	unionSize := union.Size()
	intersectionSize := intersection.Size()
	differenceSize := difference.Size()

	fmt.Printf("Union size: %d\n", unionSize)
	fmt.Printf("Intersection size: %d\n", intersectionSize)
	fmt.Printf("Difference size: %d\n", differenceSize)

	fmt.Printf("SetA contains 2: %t\n", setA.Contains(2))
	fmt.Printf("SetB contains 1: %t\n", setB.Contains(1))
}
