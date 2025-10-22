package set

type Set[T comparable] struct {
	elements map[T]struct{}
}

func New[T comparable]() *Set[T] {
	return &Set[T]{
		elements: make(map[T]struct{}),
	}
}

func (s *Set[T]) Add(element T) {
	if s == nil {
		return
	}

	s.elements[element] = struct{}{}
}

func (s *Set[T]) Remove(element T) {
	if s == nil {
		return
	}

	delete(s.elements, element)
}

func (s *Set[T]) Contains(element T) bool {
	if s == nil {
		return false
	}

	_, exists := s.elements[element]
	return exists
}

func (s *Set[T]) Size() int {
	if s == nil {
		return 0
	}

	return len(s.elements)
}

func (s *Set[T]) IsEmpty() bool {
	if s == nil {
		return true
	}

	return len(s.elements) == 0
}

func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	if s == nil {
		if other == nil {
			return New[T]()
		}

		result := New[T]()
		for element := range other.elements {
			result.Add(element)
		}
		return result
	}

	result := New[T]()

	for element := range s.elements {
		result.Add(element)
	}

	if other != nil {
		for element := range other.elements {
			result.Add(element)
		}
	}

	return result
}

func (s *Set[T]) Intersection(other *Set[T]) *Set[T] {
	result := New[T]()

	if s == nil || other == nil {
		return result
	}

	smaller, larger := s, other
	if len(other.elements) < len(s.elements) {
		smaller, larger = other, s
	}

	for element := range smaller.elements {
		if larger.Contains(element) {
			result.Add(element)
		}
	}

	return result
}

func (s *Set[T]) Difference(other *Set[T]) *Set[T] {
	result := New[T]()

	if s == nil {
		return result
	}

	if other == nil {
		for element := range s.elements {
			result.Add(element)
		}
		return result
	}

	for element := range s.elements {
		if !other.Contains(element) {
			result.Add(element)
		}
	}

	return result
}

func (s *Set[T]) Equals(other *Set[T]) bool {
	if s == nil && other == nil {
		return true
	}
	if s == nil {
		return other.Size() == 0
	}
	if other == nil {
		return s.Size() == 0
	}

	if len(s.elements) != len(other.elements) {
		return false
	}

	for element := range s.elements {
		if !other.Contains(element) {
			return false
		}
	}

	return true
}

func (s *Set[T]) Clear() {
	if s == nil {
		return
	}

	s.elements = make(map[T]struct{})
}
