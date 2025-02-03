package itrz_test

import (
	"testing"

	"github.com/dustin10/itrz"
	"github.com/dustin10/itrz/fn"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/constraints"
)

func Test_All(t *testing.T) {
	tests := map[string]struct {
		values   []int
		expected int
	}{
		"empty":     {values: []int{}, expected: 0},
		"nil":       {values: nil, expected: 0},
		"non-empty": {values: []int{1, 2, 3, 4}, expected: 4},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			s := itrz.All(test.values)

			count := 0
			for range s {
				count = count + 1
			}

			assert.Equal(t, test.expected, count)
		})
	}
}

func Test_Seq_AllMatch(t *testing.T) {
	tests := map[string]struct {
		values   []int
		test     fn.Predicate[int]
		expected bool
	}{
		"empty":                  {values: []int{}, test: equalOne, expected: true},
		"nil":                    {values: nil, test: equalOne, expected: true},
		"non-empty all-matching": {values: []int{1, 1, 1}, test: equalOne, expected: true},
		"non-empty non-matching": {values: []int{1, 2, 3, 4}, test: equalOne, expected: false},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			s := itrz.All(test.values)

			assert.Equal(t, test.expected, s.AllMatch(test.test))
		})
	}
}

func Test_Seq_AnyMatch(t *testing.T) {
	tests := map[string]struct {
		values   []int
		test     fn.Predicate[int]
		expected bool
	}{
		"empty":                   {values: []int{}, test: equalOne, expected: false},
		"nil":                     {values: nil, test: equalOne, expected: false},
		"non-empty all-matching":  {values: []int{1, 1, 1}, test: equalOne, expected: true},
		"non-empty one-matching":  {values: []int{1, 2, 3, 4}, test: equalOne, expected: true},
		"non-empty none-matching": {values: []int{2, 3, 4, 5}, test: equalOne, expected: false},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			s := itrz.All(test.values)

			assert.Equal(t, test.expected, s.AnyMatch(test.test))
		})
	}
}

func Test_Concat(t *testing.T) {
	s1 := itrz.Of(1, 2, 3)
	s2 := itrz.Of(4, 5, 6)

	s := itrz.Concat(s1, s2)

	elems := make([]int, 0)
	for e := range s {
		elems = append(elems, e)
	}

	assert.Equal(t, 6, len(elems))

	expected := 1
	for _, e := range elems {
		assert.Equal(t, expected, e)
		expected = expected + 1
	}
}

func Test_Seq_Count(t *testing.T) {
	tests := map[string]struct {
		values   []int
		expected int
	}{
		"empty":     {values: []int{}, expected: 0},
		"nil":       {values: nil, expected: 0},
		"non-empty": {values: []int{1, 1, 1}, expected: 3},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			s := itrz.All(test.values)

			assert.Equal(t, test.expected, s.Count())
		})
	}
}

func Test_Distinct(t *testing.T) {
	tests := map[string]struct {
		values   []int
		expected []int
	}{
		"empty":                   {values: []int{}, expected: []int{}},
		"nil":                     {values: nil, expected: []int{}},
		"non-empty no duplicates": {values: []int{1, 2, 3}, expected: []int{1, 2, 3}},
		"non-empty duplicates":    {values: []int{1, 1, 2, 2, 4}, expected: []int{1, 2, 4}},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			s := itrz.Distinct(itrz.All(test.values))

			res := make([]int, 0)
			for e := range s {
				res = append(res, e)
			}

			assert.Equal(t, len(test.expected), len(res))

			for idx := range test.expected {
				assert.Equal(t, test.expected[idx], res[idx])
			}
		})
	}
}

type TestSink []int

func (s *TestSink) Add(n int) {
	*s = append(*s, n)
}

func Test_Seq_DrainTo(t *testing.T) {
	tests := map[string]struct {
		values   []int
		expected int
	}{
		"empty":     {values: []int{}, expected: 0},
		"nil":       {values: nil, expected: 0},
		"non-empty": {values: []int{1, 2, 3}, expected: 3},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			sink := TestSink{}

			itrz.All(test.values).DrainTo(&sink)

			assert.Equal(t, test.expected, len(sink))

			for idx := range test.values {
				assert.Equal(t, test.values[idx], sink[idx])
			}
		})
	}
}

func Test_Empty(t *testing.T) {
	s := itrz.Empty[int]()

	count := 0
	for range s {
		count = count + 1
	}

	assert.Equal(t, 0, count)
}

func Test_Seq_Filter(t *testing.T) {
	tests := map[string]struct {
		values   []int
		test     fn.Predicate[int]
		expected []int
	}{
		"empty":                   {values: []int{}, test: isOdd, expected: []int{}},
		"nil":                     {values: nil, test: isOdd, expected: []int{}},
		"non-empty all-matching":  {values: []int{1, 1, 1}, test: isOdd, expected: []int{1, 1, 1}},
		"non-empty one-matching":  {values: []int{1, 2, 4}, test: isOdd, expected: []int{1}},
		"non-empty some-matching": {values: []int{1, 2, 3, 4, 5}, test: isOdd, expected: []int{1, 3, 5}},
		"non-empty none-matching": {values: []int{2, 4, 6}, test: isOdd, expected: []int{}},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			s := itrz.All(test.values).Filter(test.test)

			res := make([]int, 0)
			for n := range s {
				res = append(res, n)
			}

			assert.Equal(t, len(test.expected), len(res))

			for idx := range test.expected {
				assert.Equal(t, test.expected[idx], res[idx])
			}
		})
	}
}

func Test_Seq_FindAny(t *testing.T) {
	tests := map[string]struct {
		values   []int
		expected bool
	}{
		"empty":     {values: []int{}, expected: false},
		"nil":       {values: nil, expected: false},
		"non-empty": {values: []int{1, 1, 1}, expected: true},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res := itrz.All(test.values).FindAny()

			assert.Equal(t, test.expected, res.IsPresent())
		})
	}
}

func Test_FlatMap(t *testing.T) {

}

func Test_Seq_ForEach(t *testing.T) {

}

func Test_Generate(t *testing.T) {

}

func Test_GenerateWithLast(t *testing.T) {

}

func Test_Seq_Limit(t *testing.T) {

}

func Test_Map(t *testing.T) {

}

func Test_Seq_NoneMatch(t *testing.T) {
	tests := map[string]struct {
		values   []int
		test     fn.Predicate[int]
		expected bool
	}{
		"empty":                   {values: []int{}, test: equalOne, expected: true},
		"nil":                     {values: nil, test: equalOne, expected: true},
		"non-empty all-matching":  {values: []int{1, 1, 1}, test: equalOne, expected: false},
		"non-empty one-matching":  {values: []int{1, 2, 3, 4}, test: equalOne, expected: false},
		"non-empty none-matching": {values: []int{2, 3, 4, 5}, test: equalOne, expected: true},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			s := itrz.All(test.values)

			assert.Equal(t, test.expected, s.NoneMatch(test.test))
		})
	}
}

func Test_Of(t *testing.T) {

}

func Test_Seq_Peek(t *testing.T) {

}

func Test_Reduce(t *testing.T) {
	tests := map[string]struct {
		values   []int
		expected int
	}{
		"natural":   {values: []int{1, 2, 3, 4}, expected: 10},
		"negatives": {values: []int{-1, -2, -3, -4}, expected: -10},
		"alternate": {values: []int{-1, 1, -1, 1}, expected: 0},
		"zeros":     {values: []int{0, 0, 0, 0}, expected: 0},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			sum := itrz.Reduce(itrz.All(test.values), 0, sum)

			assert.Equal(t, test.expected, sum)
		})
	}
}

func Test_Seq_Skip(t *testing.T) {

}

func Test_Seq_ToSlice(t *testing.T) {

}

func Test_Zip(t *testing.T) {

}

func Test_ZipStrict(t *testing.T) {

}

type Number interface {
	constraints.Integer | constraints.Float
}

func sum[A Number](a A, b A) A {
	return a + b
}

func equalOne(n int) bool {
	return n == 1
}

func isOdd(n int) bool {
	return n%2 == 1
}
