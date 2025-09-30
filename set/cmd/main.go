package main

import (
	"fmt"
	"set/set"
)

func main() {
	setA := set.NewSet[int]()
	setB := set.NewSet[int]()

	setA.Add(1)
	setA.Add(2)
	setA.Add(3)

	setB.Add(2)
	setB.Add(3)
	setB.Add(4)

	union, err := setA.Union(setB)
	if err != nil {
		fmt.Println(err)
	}

	intersection, err := setA.Intersection(setB)
	if err != nil {
		fmt.Println(err)
	}

	difference, err := setA.Difference(setB)
	if err != nil {
		fmt.Println(err)
	}

	unionSize, _ := union.Size()
	intersectionSize, _ := intersection.Size()
	differenceSize, _ := difference.Size()

	fmt.Printf("Union size: %d\n", unionSize)
	fmt.Printf("Intersection size: %d\n", intersectionSize)
	fmt.Printf("Difference size: %d\n", differenceSize)

	fmt.Printf("SetA contains 2: %t\n", setA.Contains(2))
	fmt.Printf("SetB contains 1: %t\n", setB.Contains(1))
}
