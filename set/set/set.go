package set

import "errors"

var (
	ErrNilSet = errors.New("set is nil")
)

type Set[T comparable] struct {
	elements map[T]struct{}
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		elements: make(map[T]struct{}),
	}
}

func (s *Set[T]) Add(element T) error {
	if s == nil {
		return ErrNilSet
	}

	s.elements[element] = struct{}{}
	return nil
}

func (s *Set[T]) Remove(element T) error {
	if s == nil {
		return ErrNilSet
	}

	delete(s.elements, element)
	return nil
}

func (s *Set[T]) Contains(element T) bool {
	if s == nil {
		return false
	}

	_, exists := s.elements[element]
	return exists
}

func (s *Set[T]) Size() (int, error) {
	if s == nil {
		return 0, ErrNilSet
	}

	return len(s.elements), nil
}

func (s *Set[T]) IsEmpty() (bool, error) {
	if s == nil {
		return true, ErrNilSet
	}

	return len(s.elements) == 0, nil
}

func (s *Set[T]) Union(other *Set[T]) (*Set[T], error) {
	if s == nil {
		return nil, ErrNilSet
	}

	result := NewSet[T]()

	for element := range s.elements {
		if err := result.Add(element); err != nil {
			return nil, err
		}
	}

	if other != nil {
		for element := range other.elements {
			if err := result.Add(element); err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}

func (s *Set[T]) Intersection(other *Set[T]) (*Set[T], error) {
	if s == nil {
		return nil, ErrNilSet
	}

	result := NewSet[T]()

	if other == nil {
		return result, nil
	}

	smaller, larger := s, other
	if len(other.elements) < len(s.elements) {
		smaller, larger = other, s
	}

	for element := range smaller.elements {
		if larger.Contains(element) {
			if err := result.Add(element); err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}

func (s *Set[T]) Difference(other *Set[T]) (*Set[T], error) {
	if s == nil {
		return nil, ErrNilSet
	}

	result := NewSet[T]()

	if other == nil {
		for element := range s.elements {
			if err := result.Add(element); err != nil {
				return nil, err
			}
		}
		return result, nil
	}

	for element := range s.elements {
		if !other.Contains(element) {
			if err := result.Add(element); err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}

func (s *Set[T]) Equals(other *Set[T]) (bool, error) {
	if s == nil {
		return false, ErrNilSet
	}
	if other == nil {
		return len(s.elements) == 0, nil
	}

	if len(s.elements) != len(other.elements) {
		return false, nil
	}

	for element := range s.elements {
		if !other.Contains(element) {
			return false, nil
		}
	}

	return true, nil
}

func (s *Set[T]) Clear() error {
	if s == nil {
		return ErrNilSet
	}

	s.elements = make(map[T]struct{})
	return nil
}
