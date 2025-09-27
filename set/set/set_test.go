package set

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	s, err := NewSet()
	assert.NoError(t, err)

	err = s.Add(1)
	assert.NoError(t, err)
	assert.True(t, s.Contains(1))

	var nilSet *Set
	err = nilSet.Add(1)
	assert.ErrorIs(t, err, ErrNilSet)
}

func TestRemove(t *testing.T) {
	s, err := NewSet()
	assert.NoError(t, err)
	s.Add(1)
	s.Add(2)

	err = s.Remove(1)
	assert.NoError(t, err)
	assert.False(t, s.Contains(1))

	err = s.Remove(99)
	assert.ErrorIs(t, err, ErrElementNotFound)

	var nilSet *Set
	err = nilSet.Remove(1)
	assert.ErrorIs(t, err, ErrNilSet)
}

func TestContains(t *testing.T) {
	s, err := NewSet()
	assert.NoError(t, err)
	s.Add(1)
	s.Add("test")

	assert.True(t, s.Contains(1))
	assert.False(t, s.Contains(99))

	var nilSet *Set
	assert.False(t, nilSet.Contains(1))
}

func TestSize(t *testing.T) {
	s, err := NewSet()
	assert.NoError(t, err)

	size, err := s.Size()
	assert.NoError(t, err)
	assert.Equal(t, 0, size)

	s.Add(1)
	s.Add(2)
	s.Add(1)

	size, err = s.Size()
	assert.NoError(t, err)
	assert.Equal(t, 2, size)

	var nilSet *Set
	_, err = nilSet.Size()
	assert.ErrorIs(t, err, ErrNilSet)
}

func TestIsEmpty(t *testing.T) {
	s, err := NewSet()
	assert.NoError(t, err)

	isEmpty, err := s.IsEmpty()
	assert.NoError(t, err)
	assert.True(t, isEmpty)

	s.Add(1)
	isEmpty, err = s.IsEmpty()
	assert.NoError(t, err)
	assert.False(t, isEmpty)

	var nilSet *Set
	_, err = nilSet.IsEmpty()
	assert.ErrorIs(t, err, ErrNilSet)
}

func TestUnion(t *testing.T) {
	s1, err := NewSet()
	assert.NoError(t, err)
	s2, err := NewSet()
	assert.NoError(t, err)

	s1.Add(1)
	s1.Add(2)
	s2.Add(2)
	s2.Add(3)

	union, err := s1.Union(s2)
	assert.NoError(t, err)

	expectedElements := []interface{}{1, 2, 3}
	for _, elem := range expectedElements {
		assert.True(t, union.Contains(elem))
	}

	size, _ := union.Size()
	assert.Equal(t, 3, size)

	var nilSet *Set
	_, err = s1.Union(nilSet)
	assert.ErrorIs(t, err, ErrNilSet)
}

func TestIntersection(t *testing.T) {
	s1, err := NewSet()
	assert.NoError(t, err)
	s2, err := NewSet()
	assert.NoError(t, err)

	s1.Add(1)
	s1.Add(2)
	s1.Add(3)
	s2.Add(2)
	s2.Add(3)
	s2.Add(4)

	intersection, err := s1.Intersection(s2)
	assert.NoError(t, err)

	expectedElements := []interface{}{2, 3}
	for _, elem := range expectedElements {
		assert.True(t, intersection.Contains(elem))
	}

	assert.False(t, intersection.Contains(1))
	assert.False(t, intersection.Contains(4))

	size, _ := intersection.Size()
	assert.Equal(t, 2, size)

	var nilSet *Set
	_, err = s1.Intersection(nilSet)
	assert.ErrorIs(t, err, ErrNilSet)
}

func TestDifference(t *testing.T) {
	s1, err := NewSet()
	assert.NoError(t, err)
	s2, err := NewSet()
	assert.NoError(t, err)

	s1.Add(1)
	s1.Add(2)
	s1.Add(3)
	s2.Add(2)
	s2.Add(4)

	difference, err := s1.Difference(s2)
	assert.NoError(t, err)

	expectedElements := []interface{}{1, 3}
	for _, elem := range expectedElements {
		assert.True(t, difference.Contains(elem))
	}

	assert.False(t, difference.Contains(2))
	assert.False(t, difference.Contains(4))

	size, _ := difference.Size()
	assert.Equal(t, 2, size)

	var nilSet *Set
	_, err = s1.Difference(nilSet)
	assert.ErrorIs(t, err, ErrNilSet)
}

func TestEquals(t *testing.T) {
	s1, err := NewSet()
	assert.NoError(t, err)
	s2, err := NewSet()
	assert.NoError(t, err)

	s1.Add(1)
	s1.Add(2)
	s2.Add(2)
	s2.Add(1)

	equal, err := s1.Equals(s2)
	assert.NoError(t, err)
	assert.True(t, equal)

	s2.Add(3)
	equal, err = s1.Equals(s2)
	assert.NoError(t, err)
	assert.False(t, equal)

	var nilSet *Set
	_, err = s1.Equals(nilSet)
	assert.ErrorIs(t, err, ErrNilSet)

	var nilSet2 *Set
	_, err = nilSet.Equals(nilSet2)
	assert.ErrorIs(t, err, ErrNilSet)

	s3 := &Set{elements: nil}
	s4 := &Set{elements: nil}
	equal, err = s3.Equals(s4)
	assert.NoError(t, err)
	assert.True(t, equal)

	s5, err := NewSet()
	assert.NoError(t, err)
	equal, err = s3.Equals(s5)
	assert.NoError(t, err)
	assert.True(t, equal)
}

func TestClear(t *testing.T) {
	s, err := NewSet()
	assert.NoError(t, err)
	s.Add(1)
	s.Add(2)

	err = s.Clear()
	assert.NoError(t, err)

	isEmpty, _ := s.IsEmpty()
	assert.True(t, isEmpty)

	var nilSet *Set
	err = nilSet.Clear()
	assert.ErrorIs(t, err, ErrNilSet)
}

func TestMixedTypes(t *testing.T) {
	s, err := NewSet()
	assert.NoError(t, err)

	s.Add(1)
	s.Add("test")
	s.Add(3.14)
	s.Add(true)

	size, _ := s.Size()
	assert.Equal(t, 4, size)

	assert.True(t, s.Contains(1))
	assert.True(t, s.Contains("test"))
	assert.True(t, s.Contains(3.14))
	assert.True(t, s.Contains(true))
}

func TestDuplicates(t *testing.T) {
	s, err := NewSet()
	assert.NoError(t, err)

	s.Add(1)
	s.Add(1)
	s.Add(1)

	size, _ := s.Size()
	assert.Equal(t, 1, size)
}
