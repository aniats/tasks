package set

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSet_Add(t *testing.T) {
	tests := []struct {
		name        string
		set         *Set[int]
		element     int
		wantErr     error
		wantContain bool
	}{
		{
			name:        "add to valid set",
			set:         NewSet[int](),
			element:     1,
			wantErr:     nil,
			wantContain: true,
		},
		{
			name:        "add to nil set",
			set:         nil,
			element:     1,
			wantErr:     ErrNilSet,
			wantContain: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.set.Add(tt.element)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantContain, tt.set.Contains(tt.element))
			}
		})
	}
}

func TestSet_Remove(t *testing.T) {
	tests := []struct {
		name        string
		setupSet    func() *Set[int]
		element     int
		wantErr     error
		wantContain bool
	}{
		{
			name: "remove existing element",
			setupSet: func() *Set[int] {
				s := NewSet[int]()
				_ = s.Add(1)
				_ = s.Add(2)
				return s
			},
			element:     1,
			wantErr:     nil,
			wantContain: false,
		},
		{
			name: "remove non-existing element",
			setupSet: func() *Set[int] {
				s := NewSet[int]()
				_ = s.Add(1)
				return s
			},
			element:     99,
			wantErr:     nil,
			wantContain: false,
		},
		{
			name:        "remove from nil set",
			setupSet:    func() *Set[int] { return nil },
			element:     1,
			wantErr:     ErrNilSet,
			wantContain: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.setupSet()
			err := s.Remove(tt.element)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantContain, s.Contains(tt.element))
			}
		})
	}
}

func TestSet_Contains(t *testing.T) {
	tests := []struct {
		name     string
		setupSet func() *Set[int]
		element  int
		want     bool
	}{
		{
			name: "contains existing element",
			setupSet: func() *Set[int] {
				s := NewSet[int]()
				_ = s.Add(1)
				_ = s.Add(2)
				return s
			},
			element: 1,
			want:    true,
		},
		{
			name: "does not contain element",
			setupSet: func() *Set[int] {
				s := NewSet[int]()
				_ = s.Add(1)
				return s
			},
			element: 99,
			want:    false,
		},
		{
			name:     "nil set contains nothing",
			setupSet: func() *Set[int] { return nil },
			element:  1,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.setupSet()
			got := s.Contains(tt.element)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSet_Size(t *testing.T) {
	tests := []struct {
		name     string
		setupSet func() *Set[int]
		wantSize int
		wantErr  error
	}{
		{
			name:     "empty set",
			setupSet: func() *Set[int] { return NewSet[int]() },
			wantSize: 0,
			wantErr:  nil,
		},
		{
			name: "set with elements",
			setupSet: func() *Set[int] {
				s := NewSet[int]()
				_ = s.Add(1)
				_ = s.Add(2)
				_ = s.Add(1)
				return s
			},
			wantSize: 2,
			wantErr:  nil,
		},
		{
			name:     "nil set",
			setupSet: func() *Set[int] { return nil },
			wantSize: 0,
			wantErr:  ErrNilSet,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.setupSet()
			size, err := s.Size()
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantSize, size)
			}
		})
	}
}

func TestSet_IsEmpty(t *testing.T) {
	tests := []struct {
		name      string
		setupSet  func() *Set[int]
		wantEmpty bool
		wantErr   error
	}{
		{
			name:      "empty set",
			setupSet:  func() *Set[int] { return NewSet[int]() },
			wantEmpty: true,
			wantErr:   nil,
		},
		{
			name: "non-empty set",
			setupSet: func() *Set[int] {
				s := NewSet[int]()
				_ = s.Add(1)
				return s
			},
			wantEmpty: false,
			wantErr:   nil,
		},
		{
			name:      "nil set",
			setupSet:  func() *Set[int] { return nil },
			wantEmpty: true,
			wantErr:   ErrNilSet,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.setupSet()
			empty, err := s.IsEmpty()
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantEmpty, empty)
			}
		})
	}
}

func TestSet_Union(t *testing.T) {
	tests := []struct {
		name         string
		setupSet1    func() *Set[int]
		setupSet2    func() *Set[int]
		wantElements []int
		wantSize     int
		wantErr      error
	}{
		{
			name: "union of two sets",
			setupSet1: func() *Set[int] {
				s := NewSet[int]()
				_ = s.Add(1)
				_ = s.Add(2)
				return s
			},
			setupSet2: func() *Set[int] {
				s := NewSet[int]()
				_ = s.Add(2)
				_ = s.Add(3)
				return s
			},
			wantElements: []int{1, 2, 3},
			wantSize:     3,
			wantErr:      nil,
		},
		{
			name: "union with nil set",
			setupSet1: func() *Set[int] {
				s := NewSet[int]()
				_ = s.Add(1)
				_ = s.Add(2)
				return s
			},
			setupSet2:    func() *Set[int] { return nil },
			wantElements: []int{1, 2},
			wantSize:     2,
			wantErr:      nil,
		},
		{
			name:         "nil set union",
			setupSet1:    func() *Set[int] { return nil },
			setupSet2:    func() *Set[int] { return NewSet[int]() },
			wantElements: nil,
			wantSize:     0,
			wantErr:      ErrNilSet,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s1 := tt.setupSet1()
			s2 := tt.setupSet2()
			result, err := s1.Union(s2)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				size, _ := result.Size()
				assert.Equal(t, tt.wantSize, size)
				for _, elem := range tt.wantElements {
					assert.True(t, result.Contains(elem))
				}
			}
		})
	}
}

func TestSet_Intersection(t *testing.T) {
	tests := []struct {
		name         string
		setupSet1    func() *Set[int]
		setupSet2    func() *Set[int]
		wantElements []int
		wantSize     int
		wantErr      error
	}{
		{
			name: "intersection of two sets",
			setupSet1: func() *Set[int] {
				s := NewSet[int]()
				_ = s.Add(1)
				_ = s.Add(2)
				_ = s.Add(3)
				return s
			},
			setupSet2: func() *Set[int] {
				s := NewSet[int]()
				_ = s.Add(2)
				_ = s.Add(3)
				_ = s.Add(4)
				return s
			},
			wantElements: []int{2, 3},
			wantSize:     2,
			wantErr:      nil,
		},
		{
			name: "intersection with nil set",
			setupSet1: func() *Set[int] {
				s := NewSet[int]()
				_ = s.Add(1)
				_ = s.Add(2)
				return s
			},
			setupSet2:    func() *Set[int] { return nil },
			wantElements: []int{},
			wantSize:     0,
			wantErr:      nil,
		},
		{
			name:         "nil set intersection",
			setupSet1:    func() *Set[int] { return nil },
			setupSet2:    func() *Set[int] { return NewSet[int]() },
			wantElements: nil,
			wantSize:     0,
			wantErr:      ErrNilSet,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s1 := tt.setupSet1()
			s2 := tt.setupSet2()
			result, err := s1.Intersection(s2)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				size, _ := result.Size()
				assert.Equal(t, tt.wantSize, size)
				for _, elem := range tt.wantElements {
					assert.True(t, result.Contains(elem))
				}
			}
		})
	}
}

func TestSet_Difference(t *testing.T) {
	tests := []struct {
		name         string
		setupSet1    func() *Set[int]
		setupSet2    func() *Set[int]
		wantElements []int
		wantSize     int
		wantErr      error
	}{
		{
			name: "difference of two sets",
			setupSet1: func() *Set[int] {
				s := NewSet[int]()
				_ = s.Add(1)
				_ = s.Add(2)
				_ = s.Add(3)
				return s
			},
			setupSet2: func() *Set[int] {
				s := NewSet[int]()
				_ = s.Add(2)
				_ = s.Add(4)
				return s
			},
			wantElements: []int{1, 3},
			wantSize:     2,
			wantErr:      nil,
		},
		{
			name: "difference with nil set",
			setupSet1: func() *Set[int] {
				s := NewSet[int]()
				_ = s.Add(1)
				_ = s.Add(2)
				_ = s.Add(3)
				return s
			},
			setupSet2:    func() *Set[int] { return nil },
			wantElements: []int{1, 2, 3},
			wantSize:     3,
			wantErr:      nil,
		},
		{
			name:         "nil set difference",
			setupSet1:    func() *Set[int] { return nil },
			setupSet2:    func() *Set[int] { return NewSet[int]() },
			wantElements: nil,
			wantSize:     0,
			wantErr:      ErrNilSet,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s1 := tt.setupSet1()
			s2 := tt.setupSet2()
			result, err := s1.Difference(s2)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				size, _ := result.Size()
				assert.Equal(t, tt.wantSize, size)
				for _, elem := range tt.wantElements {
					assert.True(t, result.Contains(elem))
				}
			}
		})
	}
}

func TestSet_Equals(t *testing.T) {
	tests := []struct {
		name      string
		setupSet1 func() *Set[int]
		setupSet2 func() *Set[int]
		wantEqual bool
		wantErr   error
	}{
		{
			name: "equal sets",
			setupSet1: func() *Set[int] {
				s := NewSet[int]()
				_ = s.Add(1)
				_ = s.Add(2)
				return s
			},
			setupSet2: func() *Set[int] {
				s := NewSet[int]()
				_ = s.Add(2)
				_ = s.Add(1)
				return s
			},
			wantEqual: true,
			wantErr:   nil,
		},
		{
			name: "unequal sets",
			setupSet1: func() *Set[int] {
				s := NewSet[int]()
				_ = s.Add(1)
				_ = s.Add(2)
				return s
			},
			setupSet2: func() *Set[int] {
				s := NewSet[int]()
				_ = s.Add(1)
				_ = s.Add(2)
				_ = s.Add(3)
				return s
			},
			wantEqual: false,
			wantErr:   nil,
		},
		{
			name:      "empty set equals nil",
			setupSet1: func() *Set[int] { return NewSet[int]() },
			setupSet2: func() *Set[int] { return nil },
			wantEqual: true,
			wantErr:   nil,
		},
		{
			name: "non-empty set not equal to nil",
			setupSet1: func() *Set[int] {
				s := NewSet[int]()
				_ = s.Add(1)
				return s
			},
			setupSet2: func() *Set[int] { return nil },
			wantEqual: false,
			wantErr:   nil,
		},
		{
			name:      "nil set error",
			setupSet1: func() *Set[int] { return nil },
			setupSet2: func() *Set[int] { return NewSet[int]() },
			wantEqual: false,
			wantErr:   ErrNilSet,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s1 := tt.setupSet1()
			s2 := tt.setupSet2()
			equal, err := s1.Equals(s2)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantEqual, equal)
			}
		})
	}
}

func TestSet_Clear(t *testing.T) {
	tests := []struct {
		name     string
		setupSet func() *Set[int]
		wantErr  error
	}{
		{
			name: "clear non-empty set",
			setupSet: func() *Set[int] {
				s := NewSet[int]()
				_ = s.Add(1)
				_ = s.Add(2)
				return s
			},
			wantErr: nil,
		},
		{
			name:     "clear empty set",
			setupSet: func() *Set[int] { return NewSet[int]() },
			wantErr:  nil,
		},
		{
			name:     "clear nil set",
			setupSet: func() *Set[int] { return nil },
			wantErr:  ErrNilSet,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.setupSet()
			err := s.Clear()
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				isEmpty, _ := s.IsEmpty()
				assert.True(t, isEmpty)
			}
		})
	}
}

func TestSet_Duplicates(t *testing.T) {
	s := NewSet[int]()
	_ = s.Add(1)
	_ = s.Add(1)
	_ = s.Add(1)

	size, _ := s.Size()
	assert.Equal(t, 1, size)
}

func TestSet_StringType(t *testing.T) {
	s := NewSet[string]()
	_ = s.Add("hello")
	_ = s.Add("world")
	_ = s.Add("hello")

	size, _ := s.Size()
	assert.Equal(t, 2, size)

	assert.True(t, s.Contains("hello"))
	assert.True(t, s.Contains("world"))
	assert.False(t, s.Contains("test"))
}
