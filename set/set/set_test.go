package set

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSet_Add(t *testing.T) {
	tests := []struct {
		name     string
		set      *Set[int]
		element  int
		expected []int
	}{
		{
			name:     "add to valid set",
			set:      New[int](),
			element:  1,
			expected: []int{1},
		},
		{
			name:     "add to nil set",
			set:      nil,
			element:  1,
			expected: nil,
		},
		{
			name: "add duplicate element",
			set: func() *Set[int] {
				s := New[int]()
				s.Add(1)
				return s
			}(),
			element:  1,
			expected: []int{1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.set.Add(tt.element)

			if tt.set == nil {
				return
			}

			for _, expected := range tt.expected {
				require.True(t, tt.set.Contains(expected))
			}
			require.Equal(t, len(tt.expected), tt.set.Size())
		})
	}
}

func TestSet_Remove(t *testing.T) {
	tests := []struct {
		name             string
		initialElements  []int
		removeElement    int
		isNil            bool
		expectedSize     int
		shouldContain    []int
		shouldNotContain []int
	}{
		{
			name:             "remove existing element",
			initialElements:  []int{1, 2, 3},
			removeElement:    2,
			expectedSize:     2,
			shouldContain:    []int{1, 3},
			shouldNotContain: []int{2},
		},
		{
			name:             "remove non-existing element",
			initialElements:  []int{1, 2, 3},
			removeElement:    99,
			expectedSize:     3,
			shouldContain:    []int{1, 2, 3},
			shouldNotContain: []int{99},
		},
		{
			name:             "remove from empty set",
			initialElements:  []int{},
			removeElement:    1,
			expectedSize:     0,
			shouldContain:    []int{},
			shouldNotContain: []int{1},
		},
		{
			name:          "remove from nil set",
			isNil:         true,
			removeElement: 1,
			expectedSize:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s *Set[int]
			if tt.isNil {
				s = nil
			} else {
				s = New[int]()
				for _, elem := range tt.initialElements {
					s.Add(elem)
				}
			}

			s.Remove(tt.removeElement)

			if tt.isNil {
				return
			}

			require.Equal(t, tt.expectedSize, s.Size())
			for _, elem := range tt.shouldContain {
				require.True(t, s.Contains(elem))
			}
			for _, elem := range tt.shouldNotContain {
				require.False(t, s.Contains(elem))
			}
		})
	}
}

func TestSet_Contains(t *testing.T) {
	tests := []struct {
		name     string
		elements []int
		search   int
		isNil    bool
		expected bool
	}{
		{
			name:     "contains existing element",
			elements: []int{1, 2, 3},
			search:   2,
			expected: true,
		},
		{
			name:     "does not contain element",
			elements: []int{1, 2, 3},
			search:   99,
			expected: false,
		},
		{
			name:     "empty set contains nothing",
			elements: []int{},
			search:   1,
			expected: false,
		},
		{
			name:     "nil set contains nothing",
			isNil:    true,
			search:   1,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s *Set[int]
			if tt.isNil {
				s = nil
			} else {
				s = New[int]()
				for _, elem := range tt.elements {
					s.Add(elem)
				}
			}

			result := s.Contains(tt.search)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestSet_Size(t *testing.T) {
	tests := []struct {
		name     string
		elements []int
		isNil    bool
		expected int
	}{
		{
			name:     "empty set",
			elements: []int{},
			expected: 0,
		},
		{
			name:     "single element",
			elements: []int{1},
			expected: 1,
		},
		{
			name:     "multiple elements",
			elements: []int{1, 2, 3, 4, 5},
			expected: 5,
		},
		{
			name:     "duplicate elements",
			elements: []int{1, 1, 2, 2, 3},
			expected: 3,
		},
		{
			name:     "nil set size",
			isNil:    true,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s *Set[int]
			if tt.isNil {
				s = nil
			} else {
				s = New[int]()
				for _, elem := range tt.elements {
					s.Add(elem)
				}
			}

			require.Equal(t, tt.expected, s.Size())
		})
	}
}

func TestSet_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		elements []int
		isNil    bool
		expected bool
	}{
		{
			name:     "empty set",
			elements: []int{},
			expected: true,
		},
		{
			name:     "non-empty set",
			elements: []int{1},
			expected: false,
		},
		{
			name:     "multiple elements",
			elements: []int{1, 2, 3},
			expected: false,
		},
		{
			name:     "nil set is empty",
			isNil:    true,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s *Set[int]
			if tt.isNil {
				s = nil
			} else {
				s = New[int]()
				for _, elem := range tt.elements {
					s.Add(elem)
				}
			}

			require.Equal(t, tt.expected, s.IsEmpty())
		})
	}
}

func TestSet_Union(t *testing.T) {
	tests := []struct {
		name     string
		set1     []int
		set2     []int
		set1Nil  bool
		set2Nil  bool
		expected []int
	}{
		{
			name:     "union of two non-empty sets",
			set1:     []int{1, 2},
			set2:     []int{2, 3, 4},
			expected: []int{1, 2, 3, 4},
		},
		{
			name:     "union with empty set",
			set1:     []int{1, 2},
			set2:     []int{},
			expected: []int{1, 2},
		},
		{
			name:     "union of empty sets",
			set1:     []int{},
			set2:     []int{},
			expected: []int{},
		},
		{
			name:     "union with identical sets",
			set1:     []int{1, 2, 3},
			set2:     []int{1, 2, 3},
			expected: []int{1, 2, 3},
		},
		{
			name:     "union nil with nil",
			set1Nil:  true,
			set2Nil:  true,
			expected: []int{},
		},
		{
			name:     "union nil with non-empty",
			set1Nil:  true,
			set2:     []int{1, 2},
			expected: []int{1, 2},
		},
		{
			name:     "union non-empty with nil",
			set1:     []int{1, 2},
			set2Nil:  true,
			expected: []int{1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s1, s2 *Set[int]

			if tt.set1Nil {
				s1 = nil
			} else {
				s1 = New[int]()
				for _, elem := range tt.set1 {
					s1.Add(elem)
				}
			}

			if tt.set2Nil {
				s2 = nil
			} else {
				s2 = New[int]()
				for _, elem := range tt.set2 {
					s2.Add(elem)
				}
			}

			result := s1.Union(s2)

			require.NotNil(t, result)
			require.Equal(t, len(tt.expected), result.Size())
			for _, elem := range tt.expected {
				require.True(t, result.Contains(elem))
			}
		})
	}
}

func TestSet_Intersection(t *testing.T) {
	tests := []struct {
		name     string
		set1     []int
		set2     []int
		set1Nil  bool
		set2Nil  bool
		expected []int
	}{
		{
			name:     "intersection of overlapping sets",
			set1:     []int{1, 2, 3, 4},
			set2:     []int{3, 4, 5, 6},
			expected: []int{3, 4},
		},
		{
			name:     "intersection with no overlap",
			set1:     []int{1, 2},
			set2:     []int{3, 4},
			expected: []int{},
		},
		{
			name:     "intersection with identical sets",
			set1:     []int{1, 2, 3},
			set2:     []int{1, 2, 3},
			expected: []int{1, 2, 3},
		},
		{
			name:     "intersection with empty set",
			set1:     []int{1, 2, 3},
			set2:     []int{},
			expected: []int{},
		},
		{
			name:     "intersection nil with nil",
			set1Nil:  true,
			set2Nil:  true,
			expected: []int{},
		},
		{
			name:     "intersection nil with non-empty",
			set1Nil:  true,
			set2:     []int{1, 2},
			expected: []int{},
		},
		{
			name:     "intersection non-empty with nil",
			set1:     []int{1, 2},
			set2Nil:  true,
			expected: []int{},
		},
		{
			name:     "intersection single elements",
			set1:     []int{5},
			set2:     []int{5},
			expected: []int{5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s1, s2 *Set[int]

			if tt.set1Nil {
				s1 = nil
			} else {
				s1 = New[int]()
				for _, elem := range tt.set1 {
					s1.Add(elem)
				}
			}

			if tt.set2Nil {
				s2 = nil
			} else {
				s2 = New[int]()
				for _, elem := range tt.set2 {
					s2.Add(elem)
				}
			}

			result := s1.Intersection(s2)

			require.NotNil(t, result)
			require.Equal(t, len(tt.expected), result.Size())
			for _, elem := range tt.expected {
				require.True(t, result.Contains(elem))
			}
		})
	}
}

func TestSet_Difference(t *testing.T) {
	tests := []struct {
		name     string
		set1     []int
		set2     []int
		set1Nil  bool
		set2Nil  bool
		expected []int
	}{
		{
			name:     "difference of overlapping sets",
			set1:     []int{1, 2, 3, 4},
			set2:     []int{3, 4, 5, 6},
			expected: []int{1, 2},
		},
		{
			name:     "difference with no overlap",
			set1:     []int{1, 2},
			set2:     []int{3, 4},
			expected: []int{1, 2},
		},
		{
			name:     "difference with identical sets",
			set1:     []int{1, 2, 3},
			set2:     []int{1, 2, 3},
			expected: []int{},
		},
		{
			name:     "difference with empty set",
			set1:     []int{1, 2, 3},
			set2:     []int{},
			expected: []int{1, 2, 3},
		},
		{
			name:     "difference nil with nil",
			set1Nil:  true,
			set2Nil:  true,
			expected: []int{},
		},
		{
			name:     "difference nil with non-empty",
			set1Nil:  true,
			set2:     []int{1, 2},
			expected: []int{},
		},
		{
			name:     "difference non-empty with nil",
			set1:     []int{1, 2, 3},
			set2Nil:  true,
			expected: []int{1, 2, 3},
		},
		{
			name:     "difference empty with non-empty",
			set1:     []int{},
			set2:     []int{1, 2, 3},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s1, s2 *Set[int]

			if tt.set1Nil {
				s1 = nil
			} else {
				s1 = New[int]()
				for _, elem := range tt.set1 {
					s1.Add(elem)
				}
			}

			if tt.set2Nil {
				s2 = nil
			} else {
				s2 = New[int]()
				for _, elem := range tt.set2 {
					s2.Add(elem)
				}
			}

			result := s1.Difference(s2)

			require.NotNil(t, result)
			require.Equal(t, len(tt.expected), result.Size())
			for _, elem := range tt.expected {
				require.True(t, result.Contains(elem))
			}
		})
	}
}

func TestSet_Equals(t *testing.T) {
	tests := []struct {
		name     string
		set1     []int
		set2     []int
		set1Nil  bool
		set2Nil  bool
		expected bool
	}{
		{
			name:     "equal sets same order",
			set1:     []int{1, 2, 3},
			set2:     []int{1, 2, 3},
			expected: true,
		},
		{
			name:     "equal sets different order",
			set1:     []int{1, 2, 3},
			set2:     []int{3, 1, 2},
			expected: true,
		},
		{
			name:     "different size sets",
			set1:     []int{1, 2, 3},
			set2:     []int{1, 2},
			expected: false,
		},
		{
			name:     "different elements",
			set1:     []int{1, 2, 3},
			set2:     []int{1, 2, 4},
			expected: false,
		},
		{
			name:     "both empty sets",
			set1:     []int{},
			set2:     []int{},
			expected: true,
		},
		{
			name:     "both nil sets",
			set1Nil:  true,
			set2Nil:  true,
			expected: true,
		},
		{
			name:     "nil equals empty",
			set1Nil:  true,
			set2:     []int{},
			expected: true,
		},
		{
			name:     "empty equals nil",
			set1:     []int{},
			set2Nil:  true,
			expected: true,
		},
		{
			name:     "nil not equal to non-empty",
			set1Nil:  true,
			set2:     []int{1},
			expected: false,
		},
		{
			name:     "non-empty not equal to nil",
			set1:     []int{1},
			set2Nil:  true,
			expected: false,
		},
		{
			name:     "single element sets equal",
			set1:     []int{42},
			set2:     []int{42},
			expected: true,
		},
		{
			name:     "single element sets not equal",
			set1:     []int{42},
			set2:     []int{43},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s1, s2 *Set[int]

			if tt.set1Nil {
				s1 = nil
			} else {
				s1 = New[int]()
				for _, elem := range tt.set1 {
					s1.Add(elem)
				}
			}

			if tt.set2Nil {
				s2 = nil
			} else {
				s2 = New[int]()
				for _, elem := range tt.set2 {
					s2.Add(elem)
				}
			}

			result := s1.Equals(s2)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestSet_Clear(t *testing.T) {
	tests := []struct {
		name     string
		elements []int
		isNil    bool
	}{
		{
			name:     "clear non-empty set",
			elements: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "clear empty set",
			elements: []int{},
		},
		{
			name:     "clear single element set",
			elements: []int{1},
		},
		{
			name:     "clear large set",
			elements: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
		},
		{
			name:  "clear nil set",
			isNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.isNil {
				var s *Set[int]
				s.Clear()
				return
			}

			s := New[int]()
			for _, elem := range tt.elements {
				s.Add(elem)
			}

			s.Clear()

			require.True(t, s.IsEmpty())
			require.Equal(t, 0, s.Size())
		})
	}
}

func TestSet_StringType(t *testing.T) {
	tests := []struct {
		name          string
		elements      []string
		testElement   string
		shouldContain bool
		expectedSize  int
	}{
		{
			name:          "string set basic operations",
			elements:      []string{"hello", "world", "hello"},
			testElement:   "hello",
			shouldContain: true,
			expectedSize:  2,
		},
		{
			name:          "string set missing element",
			elements:      []string{"hello", "world"},
			testElement:   "test",
			shouldContain: false,
			expectedSize:  2,
		},
		{
			name:          "empty string in set",
			elements:      []string{"", "hello", "world"},
			testElement:   "",
			shouldContain: true,
			expectedSize:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New[string]()
			for _, elem := range tt.elements {
				s.Add(elem)
			}

			require.Equal(t, tt.expectedSize, s.Size())
			require.Equal(t, tt.shouldContain, s.Contains(tt.testElement))
		})
	}
}
