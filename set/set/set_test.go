package set

import (
	"errors"
	"testing"
)

func TestAdd(t *testing.T) {
	s, err := NewSet()
	if err != nil {
		t.Fatalf("NewSet() returned error: %v", err)
	}

	err = s.Add(1)
	if err != nil {
		t.Errorf("Add() returned error: %v", err)
	}

	if !s.Contains(1) {
		t.Error("element not added")
	}

	var nilSet *Set
	err = nilSet.Add(1)
	if !errors.Is(err, ErrNilSet) {
		t.Error("expected ErrNilSet for nil set")
	}
}

func TestRemove(t *testing.T) {
	s, err := NewSet()
	if err != nil {
		t.Fatalf("NewSet() returned error: %v", err)
	}
	s.Add(1)
	s.Add(2)

	err = s.Remove(1)
	if err != nil {
		t.Errorf("Remove() returned error: %v", err)
	}

	if s.Contains(1) {
		t.Error("element not removed")
	}

	err = s.Remove(99)
	if !errors.Is(err, ErrElementNotFound) {
		t.Error("expected ErrElementNotFound")
	}

	var nilSet *Set
	err = nilSet.Remove(1)
	if !errors.Is(err, ErrNilSet) {
		t.Error("expected ErrNilSet for nil set")
	}
}

func TestContains(t *testing.T) {
	s, err := NewSet()
	if err != nil {
		t.Fatalf("NewSet() returned error: %v", err)
	}
	s.Add(1)
	s.Add("test")

	if !s.Contains(1) {
		t.Error("Contains() returned false for existing element")
	}

	if s.Contains(99) {
		t.Error("Contains() returned true for non-existing element")
	}

	var nilSet *Set
	if nilSet.Contains(1) {
		t.Error("nil set should not contain elements")
	}
}

func TestSize(t *testing.T) {
	s, err := NewSet()
	if err != nil {
		t.Fatalf("NewSet() returned error: %v", err)
	}

	size, err := s.Size()
	if err != nil {
		t.Errorf("Size() returned error: %v", err)
	}
	if size != 0 {
		t.Errorf("expected size 0, got %d", size)
	}

	s.Add(1)
	s.Add(2)
	s.Add(1)

	size, err = s.Size()
	if err != nil {
		t.Errorf("Size() returned error: %v", err)
	}
	if size != 2 {
		t.Errorf("expected size 2, got %d", size)
	}

	var nilSet *Set
	_, err = nilSet.Size()
	if !errors.Is(err, ErrNilSet) {
		t.Error("expected ErrNilSet for nil set")
	}
}

func TestIsEmpty(t *testing.T) {
	s, err := NewSet()
	if err != nil {
		t.Fatalf("NewSet() returned error: %v", err)
	}

	isEmpty, err := s.IsEmpty()
	if err != nil {
		t.Errorf("IsEmpty() returned error: %v", err)
	}
	if !isEmpty {
		t.Error("new set should be empty")
	}

	s.Add(1)
	isEmpty, err = s.IsEmpty()
	if err != nil {
		t.Errorf("IsEmpty() returned error: %v", err)
	}
	if isEmpty {
		t.Error("set with element should not be empty")
	}

	var nilSet *Set
	_, err = nilSet.IsEmpty()
	if !errors.Is(err, ErrNilSet) {
		t.Error("expected ErrNilSet for nil set")
	}
}

func TestUnion(t *testing.T) {
	s1, err := NewSet()
	if err != nil {
		t.Fatalf("NewSet() returned error: %v", err)
	}
	s2, err := NewSet()
	if err != nil {
		t.Fatalf("NewSet() returned error: %v", err)
	}

	s1.Add(1)
	s1.Add(2)
	s2.Add(2)
	s2.Add(3)

	union, err := s1.Union(s2)
	if err != nil {
		t.Errorf("Union() returned error: %v", err)
	}

	expectedElements := []interface{}{1, 2, 3}
	for _, elem := range expectedElements {
		if !union.Contains(elem) {
			t.Errorf("union missing element %v", elem)
		}
	}

	size, _ := union.Size()
	if size != 3 {
		t.Errorf("expected union size 3, got %d", size)
	}

	var nilSet *Set
	_, err = s1.Union(nilSet)
	if !errors.Is(err, ErrNilSet) {
		t.Error("expected ErrNilSet for nil set")
	}
}

func TestIntersection(t *testing.T) {
	s1, err := NewSet()
	if err != nil {
		t.Fatalf("NewSet() returned error: %v", err)
	}
	s2, err := NewSet()
	if err != nil {
		t.Fatalf("NewSet() returned error: %v", err)
	}

	s1.Add(1)
	s1.Add(2)
	s1.Add(3)
	s2.Add(2)
	s2.Add(3)
	s2.Add(4)

	intersection, err := s1.Intersection(s2)
	if err != nil {
		t.Errorf("Intersection() returned error: %v", err)
	}

	expectedElements := []interface{}{2, 3}
	for _, elem := range expectedElements {
		if !intersection.Contains(elem) {
			t.Errorf("intersection missing element %v", elem)
		}
	}

	if intersection.Contains(1) || intersection.Contains(4) {
		t.Error("intersection contains unexpected elements")
	}

	size, _ := intersection.Size()
	if size != 2 {
		t.Errorf("expected intersection size 2, got %d", size)
	}

	var nilSet *Set
	_, err = s1.Intersection(nilSet)
	if !errors.Is(err, ErrNilSet) {
		t.Error("expected ErrNilSet for nil set")
	}
}

func TestDifference(t *testing.T) {
	s1, err := NewSet()
	if err != nil {
		t.Fatalf("NewSet() returned error: %v", err)
	}
	s2, err := NewSet()
	if err != nil {
		t.Fatalf("NewSet() returned error: %v", err)
	}

	s1.Add(1)
	s1.Add(2)
	s1.Add(3)
	s2.Add(2)
	s2.Add(4)

	difference, err := s1.Difference(s2)
	if err != nil {
		t.Errorf("Difference() returned error: %v", err)
	}

	expectedElements := []interface{}{1, 3}
	for _, elem := range expectedElements {
		if !difference.Contains(elem) {
			t.Errorf("difference missing element %v", elem)
		}
	}

	if difference.Contains(2) || difference.Contains(4) {
		t.Error("difference contains unexpected elements")
	}

	size, _ := difference.Size()
	if size != 2 {
		t.Errorf("expected difference size 2, got %d", size)
	}

	var nilSet *Set
	_, err = s1.Difference(nilSet)
	if !errors.Is(err, ErrNilSet) {
		t.Error("expected ErrNilSet for nil set")
	}
}

func TestEquals(t *testing.T) {
	s1, err := NewSet()
	if err != nil {
		t.Fatalf("NewSet() returned error: %v", err)
	}
	s2, err := NewSet()
	if err != nil {
		t.Fatalf("NewSet() returned error: %v", err)
	}

	s1.Add(1)
	s1.Add(2)
	s2.Add(2)
	s2.Add(1)

	equal, err := s1.Equals(s2)
	if err != nil {
		t.Errorf("Equals() returned error: %v", err)
	}
	if !equal {
		t.Error("sets should be equal")
	}

	s2.Add(3)
	equal, err = s1.Equals(s2)
	if err != nil {
		t.Errorf("Equals() returned error: %v", err)
	}
	if equal {
		t.Error("sets should not be equal")
	}

	var nilSet *Set
	_, err = s1.Equals(nilSet)
	if !errors.Is(err, ErrNilSet) {
		t.Error("expected ErrNilSet for nil set")
	}

	var nilSet2 *Set
	_, err = nilSet.Equals(nilSet2)
	if !errors.Is(err, ErrNilSet) {
		t.Error("expected ErrNilSet for both nil sets")
	}

	s3 := &Set{elements: nil}
	s4 := &Set{elements: nil}
	equal, err = s3.Equals(s4)
	if err != nil {
		t.Errorf("Equals() returned error: %v", err)
	}
	if !equal {
		t.Error("empty sets with nil elements should be equal")
	}

	s5, err := NewSet()
	if err != nil {
		t.Fatalf("NewSet() returned error: %v", err)
	}
	equal, err = s3.Equals(s5)
	if err != nil {
		t.Errorf("Equals() returned error: %v", err)
	}
	if !equal {
		t.Error("empty sets should be equal regardless of nil elements")
	}
}

func TestClear(t *testing.T) {
	s, err := NewSet()
	if err != nil {
		t.Fatalf("NewSet() returned error: %v", err)
	}
	s.Add(1)
	s.Add(2)

	err = s.Clear()
	if err != nil {
		t.Errorf("Clear() returned error: %v", err)
	}

	isEmpty, _ := s.IsEmpty()
	if !isEmpty {
		t.Error("set should be empty after clear")
	}

	var nilSet *Set
	err = nilSet.Clear()
	if !errors.Is(err, ErrNilSet) {
		t.Error("expected ErrNilSet for nil set")
	}
}

func TestMixedTypes(t *testing.T) {
	s, err := NewSet()
	if err != nil {
		t.Fatalf("NewSet() returned error: %v", err)
	}

	s.Add(1)
	s.Add("test")
	s.Add(3.14)
	s.Add(true)

	size, _ := s.Size()
	if size != 4 {
		t.Errorf("expected size 4, got %d", size)
	}

	if !s.Contains(1) || !s.Contains("test") || !s.Contains(3.14) || !s.Contains(true) {
		t.Error("set should contain all added elements")
	}
}

func TestDuplicates(t *testing.T) {
	s, err := NewSet()
	if err != nil {
		t.Fatalf("NewSet() returned error: %v", err)
	}
	
	s.Add(1)
	s.Add(1)
	s.Add(1)

	size, _ := s.Size()
	if size != 1 {
		t.Errorf("expected size 1 after adding duplicates, got %d", size)
	}
}
