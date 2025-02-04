package set_test

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dustin10/itrz/set"
)

func Test_New(t *testing.T) {
	s := set.New[int]()

	assert.True(t, s.IsEmpty())
	assert.Equal(t, 0, s.Len())
}

func Test_FromSlice(t *testing.T) {
	tests := map[string]struct {
		values   []int
		expected []int
	}{
		"empty":      {values: []int{}, expected: []int{}},
		"nil":        {values: nil, expected: []int{}},
		"one":        {values: []int{1}, expected: []int{1}},
		"many":       {values: []int{1, 2, 3, 4}, expected: []int{1, 2, 3, 4}},
		"duplicates": {values: []int{1, 2, 3, 1, 4, 3}, expected: []int{1, 2, 3, 4}},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			s := set.FromSlice(test.values)

			assert.Equal(t, len(test.expected), s.Len())

			for _, e := range test.expected {
				assert.True(t, s.Contains(e))
			}
		})
	}
}

func Test_WithIntialCapacity(t *testing.T) {
	cfg := set.Config{}

	set.WithInitialCapacity(100)(&cfg)

	assert.Equal(t, 100, cfg.InitialCapacity)
}

func Test_Set_IsEmpty(t *testing.T) {
	tests := map[string]struct {
		values   []int
		expected bool
	}{
		"empty":      {values: []int{}, expected: true},
		"nil":        {values: nil, expected: true},
		"one":        {values: []int{1}},
		"many":       {values: []int{1, 2, 3, 4}},
		"duplicates": {values: []int{1, 2, 3, 1, 4, 3}},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			s := set.FromSlice(test.values)

			assert.Equal(t, test.expected, s.IsEmpty())
		})
	}
}

func Test_Set_Len(t *testing.T) {
	tests := map[string]struct {
		values   []int
		expected int
	}{
		"empty":      {values: []int{}},
		"nil":        {values: nil},
		"one":        {values: []int{1}, expected: 1},
		"many":       {values: []int{1, 2, 3, 4}, expected: 4},
		"duplicates": {values: []int{1, 2, 3, 1, 4, 3}, expected: 4},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			s := set.FromSlice(test.values)

			assert.Equal(t, test.expected, s.Len())
		})
	}
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
			s := set.FromSlice(test.values)

			assert.Equal(t, len(test.expected), s.Len())

			for _, e := range test.expected {
				assert.True(t, s.Contains(e))
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
			s := set.FromSlice(test.values)

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
			s := set.FromSlice(test.values)

			assert.Equal(t, test.expect, s.Contains(test.needle))
		})
	}
}

func Test_Set_Clear(t *testing.T) {
	s := set.New[int]()
	s.Add(1)
	s.Add(2)
	s.Add(3)

	s.Clear()

	assert.Equal(t, 0, s.Len())
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
			s := set.FromSlice(test.values)

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
	f := func(n int) set.Set[int] {
		s := set.New[int]()
		s.Add(n)
		return s
	}

	tests := map[string]struct {
		values []int
		fn     func(n int) set.Set[int]
	}{
		"empty":     {values: []int{}, fn: f},
		"nil":       {values: nil, fn: f},
		"non-empty": {values: []int{1, 2, 3}, fn: f},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			s := set.FromSlice(test.values)

			res := set.FlatMap(s, test.fn)

			for _, v := range test.values {
				assert.True(t, res.Contains(v))
			}
		})
	}
}

func Test_Map(t *testing.T) {
	f := func(n int) int {
		return 2 * n
	}

	tests := map[string]struct {
		values []int
		fn     func(n int) int
	}{
		"empty":     {values: []int{}, fn: f},
		"nil":       {values: nil, fn: f},
		"non-empty": {values: []int{1, 2, 3}, fn: f},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			s := set.FromSlice(test.values)

			res := set.Map(s, test.fn)

			for _, v := range test.values {
				assert.True(t, res.Contains(test.fn(v)))
			}
		})
	}
}
