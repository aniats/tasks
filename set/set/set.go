package set

import "errors"

var (
	ErrElementNotFound = errors.New("element not found in set")
	ErrNilSet          = errors.New("set is nil")
	ErrEmptySet        = errors.New("set is empty")
)

type Set struct {
	elements map[interface{}]struct{}
}

func NewSet() (*Set, error) {
	return &Set{
		elements: make(map[interface{}]struct{}),
	}, nil
}

func (s *Set) Add(element interface{}) error {
	if s == nil {
		return ErrNilSet
	}

	s.elements[element] = struct{}{}
	return nil
}

func (s *Set) Remove(element interface{}) error {
	if s == nil {
		return ErrNilSet
	}
	if s.elements == nil {
		return ErrEmptySet
	}

	if !s.Contains(element) {
		return ErrElementNotFound
	}

	delete(s.elements, element)
	return nil
}

func (s *Set) Contains(element interface{}) bool {
	if s == nil || s.elements == nil {
		return false
	}

	_, exists := s.elements[element]
	return exists
}

func (s *Set) Size() (int, error) {
	if s == nil {
		return 0, ErrNilSet
	}
	if s.elements == nil {
		return 0, nil
	}

	return len(s.elements), nil
}

func (s *Set) IsEmpty() (bool, error) {
	if s == nil {
		return true, ErrNilSet
	}
	if s.elements == nil {
		return true, nil
	}

	return len(s.elements) == 0, nil
}

func (s *Set) Union(other *Set) (*Set, error) {
	if s == nil || other == nil {
		return nil, ErrNilSet
	}

	result, err := NewSet()
	if err != nil {
		return nil, err
	}

	if s.elements != nil {
		for element := range s.elements {
			if err := result.Add(element); err != nil {
				return nil, err
			}
		}
	}

	if other.elements != nil {
		for element := range other.elements {
			if err := result.Add(element); err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}

func (s *Set) Intersection(other *Set) (*Set, error) {
	if s == nil || other == nil {
		return nil, ErrNilSet
	}

	result, err := NewSet()
	if err != nil {
		return nil, err
	}

	if s.elements == nil || other.elements == nil {
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

func (s *Set) Difference(other *Set) (*Set, error) {
	if s == nil {
		return nil, ErrNilSet
	}
	if other == nil {
		return nil, ErrNilSet
	}

	result, err := NewSet()
	if err != nil {
		return nil, err
	}

	if s.elements == nil {
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

func (s *Set) Equals(other *Set) (bool, error) {
	if s == nil || other == nil {
		return false, ErrNilSet
	}

	sSize, err := s.Size()
	if err != nil {
		return false, err
	}

	otherSize, err := other.Size()
	if err != nil {
		return false, err
	}

	if sSize != otherSize {
		return false, nil
	}

	if s.elements != nil {
		for element := range s.elements {
			if !other.Contains(element) {
				return false, nil
			}
		}
	}

	return true, nil
}

func (s *Set) Clear() error {
	if s == nil {
		return ErrNilSet
	}

	s.elements = make(map[interface{}]struct{})
	return nil
}
