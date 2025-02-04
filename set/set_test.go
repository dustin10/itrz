package set

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	s := New[int]()

	assert.Equal(t, defaultInitialCapacity, s.config.InitialCapacity)
	assert.Equal(t, 0, len(s.elems))
}

func Test_WithIntialCapacity(t *testing.T) {
	cfg := Config{}

	WithInitialCapacity(100)(&cfg)

	assert.Equal(t, 100, cfg.InitialCapacity)
}

func Test_Set_Add(t *testing.T) {
	tests := map[string]struct {
		values   []int
		expected []int
	}{
		"none":       {values: []int{}, expected: []int{}},
		"one":        {values: []int{1}, expected: []int{1}},
		"many":       {values: []int{1, 2, 3, 4}, expected: []int{1, 2, 3, 4}},
		"duplicates": {values: []int{1, 2, 3, 1, 4, 3}, expected: []int{1, 2, 3, 4}},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			s := New[int]()
			for _, e := range test.values {
				s.Add(e)
			}

			assert.Equal(t, len(test.expected), len(s.elems))

			for _, e := range test.expected {
				_, exists := s.elems[e]
				assert.True(t, exists)
			}
		})
	}
}

func Test_Set_Remove(t *testing.T) {
	tests := map[string]struct {
		values []int
		remove int
		expect bool
	}{
		"remove empty":          {values: []int{}, remove: 1},
		"remove does not exist": {values: []int{2, 3}, remove: 1},
		"remove exists":         {values: []int{1, 2, 3}, remove: 1, expect: true},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			s := New[int]()
			for _, e := range test.values {
				s.Add(e)
			}

			assert.Equal(t, test.expect, s.Remove(test.remove))
		})
	}
}

func Test_Set_Contains(t *testing.T) {
	tests := map[string]struct {
		values []int
		needle int
		expect bool
	}{
		"empty":          {values: []int{}, needle: 1},
		"does not exist": {values: []int{2, 3}, needle: 1},
		"exists":         {values: []int{1, 2, 3}, needle: 1, expect: true},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			s := New[int]()
			for _, e := range test.values {
				s.Add(e)
			}

			assert.Equal(t, test.expect, s.Contains(test.needle))
		})
	}
}

func Test_Set_Clear(t *testing.T) {
	s := New[int]()
	s.Add(1)
	s.Add(2)
	s.Add(3)

	s.Clear()

	assert.Equal(t, 0, len(s.elems))
}

func Test_Set_All(t *testing.T) {
	tests := map[string]struct {
		values []int
	}{
		"empty": {values: []int{}},
		"one":   {values: []int{1}},
		"many":  {values: []int{1, 2, 3}},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			s := New[int]()
			for _, e := range test.values {
				s.Add(e)
			}

			res := make([]int, 0)
			for e := range s.All() {
				res = append(res, e)
			}

			assert.Equal(t, len(test.values), len(res))

			for _, e := range test.values {
				assert.True(t, slices.Contains(res, e))
			}
		})
	}
}

func Test_FlatMap(t *testing.T) {

}

func Test_Map(t *testing.T) {

}
