package itrz_test

import (
	"math"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/constraints"

	"github.com/dustin10/itrz"
	"github.com/dustin10/itrz/fn"
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

			assert.Equal(t, test.expected, countElems(s))
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

			assert.True(t, slices.Equal(test.expected, res))
		})
	}
}

type TestSink []int

func (s *TestSink) Add(n int) {
	*s = append(*s, n)
}

func Test_Seq_DrainTo(t *testing.T) {
	tests := map[string]struct {
		values []int
	}{
		"empty":     {values: []int{}},
		"nil":       {values: nil},
		"non-empty": {values: []int{1, 2, 3}},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			sink := TestSink{}

			itrz.All(test.values).DrainTo(&sink)

			assert.True(t, slices.Equal(test.values, sink))
		})
	}
}

func Test_Empty(t *testing.T) {
	s := itrz.Empty[int]()

	assert.Equal(t, 0, countElems(s))
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

			assert.True(t, slices.Equal(test.expected, res))
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
	f := func(n int) itrz.Seq[int] {
		return itrz.Of(n)
	}

	tests := map[string]struct {
		values []int
		fn     func(n int) itrz.Seq[int]
	}{
		"empty":     {values: []int{}, fn: f},
		"nil":       {values: nil, fn: f},
		"non-empty": {values: []int{1, 2, 3}, fn: f},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			s := itrz.FlatMap(itrz.All(test.values), test.fn)

			res := make([]int, 0)
			for e := range s {
				res = append(res, e)
			}

			assert.True(t, slices.Equal(test.values, res))
		})
	}
}

func Test_Seq_ForEach(t *testing.T) {
	var sum int
	f := func(n int) { sum = sum + n }

	tests := map[string]struct {
		values   []int
		fn       func(n int)
		expected int
	}{
		"empty":     {values: []int{}, fn: f},
		"nil":       {values: nil, fn: f},
		"non-empty": {values: []int{1, 2, 3}, fn: f, expected: 6},
	}

	for name, test := range tests {
		sum = 0

		t.Run(name, func(t *testing.T) {
			itrz.All(test.values).ForEach(test.fn)

			assert.Equal(t, test.expected, sum)
		})
	}
}

func Test_Generate(t *testing.T) {
	last := 0

	s := itrz.Generate(func() int {
		last = last + 1
		return last
	})

	count := 0
	for e := range s {
		count = count + 1

		if e == 10 {
			break
		}
	}

	assert.Equal(t, 10, count)
}

func Test_GenerateWithLast(t *testing.T) {
	s := itrz.GenerateWithLast(0, func(n int) int {
		return n + 1
	})

	count := 0
	for e := range s {
		count = count + 1

		if e == 10 {
			break
		}
	}

	assert.Equal(t, 10, count)
}

func Test_Seq_Limit(t *testing.T) {
	tests := map[string]struct {
		values []int
		limit  int
	}{
		"empty":     {values: []int{}},
		"nil":       {values: nil},
		"non-empty": {values: []int{1, 2, 3, 4}, limit: 2},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			s := itrz.All(test.values).Limit(test.limit)

			assert.Equal(t, test.limit, countElems(s))
		})
	}
}

func Test_Map(t *testing.T) {
	f := func(n int) int {
		return 2 * n
	}

	tests := map[string]struct {
		values   []int
		fn       func(n int) int
		expected []int
	}{
		"empty":     {values: []int{}, fn: f},
		"nil":       {values: nil, fn: f},
		"non-empty": {values: []int{1, 2, 3}, fn: f, expected: []int{2, 4, 6}},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			s := itrz.Map(itrz.All(test.values), test.fn)

			res := make([]int, 0)
			for e := range s {
				res = append(res, e)
			}

			assert.True(t, slices.Equal(test.expected, res))
		})
	}
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
	tests := map[string]struct {
		values   []int
		expected int
	}{
		"empty":  {values: []int{}},
		"nil":    {values: nil},
		"single": {values: []int{1}, expected: 1},
		"multi":  {values: []int{1, 2, 3, 4}, expected: 4},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			s := itrz.Of(test.values...)

			assert.Equal(t, test.expected, countElems(s))
		})
	}
}

func Test_Seq_Peek(t *testing.T) {
	var res []int
	f := func(n int) { res = append(res, n) }

	tests := map[string]struct {
		values []int
		fn     func(n int)
	}{
		"empty":     {values: []int{}, fn: f},
		"nil":       {values: nil, fn: f},
		"non-empty": {values: []int{1, 2, 3}, fn: f},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res = make([]int, 0)

			s := itrz.All(test.values).Peek(test.fn)

			for range s {
			}

			assert.True(t, slices.Equal(test.values, res))
		})
	}
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
	tests := map[string]struct {
		values   []int
		n        int
		expected []int
	}{
		"empty":            {values: []int{}},
		"nil":              {values: nil},
		"non-empty skip 0": {values: []int{1, 2, 3}, expected: []int{1, 2, 3}},
		"non-empty skip":   {values: []int{1, 2, 3, 4}, n: 2, expected: []int{3, 4}},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			s := itrz.All(test.values).Skip(test.n)

			res := make([]int, 0)
			for e := range s {
				res = append(res, e)
			}

			assert.True(t, slices.Equal(test.expected, res))
		})
	}
}

func Test_Seq_ToSlice(t *testing.T) {
	tests := map[string]struct {
		values []int
	}{
		"empty":     {values: []int{}},
		"nil":       {values: nil},
		"non-empty": {values: []int{1, 2, 3, 4}},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res := itrz.All(test.values).ToSlice()

			assert.True(t, slices.Equal(test.values, res))
		})
	}
}

func Test_Zip(t *testing.T) {
	tests := map[string]struct {
		as []int
		bs []string
	}{
		"as empty": {as: []int{}, bs: []string{"1", "2", "3"}},
		"as nil":   {as: nil, bs: []string{"1", "2", "3"}},
		"bs empty": {as: []int{1, 2, 3}, bs: []string{}},
		"bs nil":   {as: []int{1, 2, 3}, bs: nil},
		"same len": {as: []int{1, 2, 3}, bs: []string{"1", "2", "3"}},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			as := itrz.All(test.as)
			bs := itrz.All(test.bs)

			max := int(math.Max(float64(len(test.as)), float64(len(test.bs))))

			resAs := make([]int, 0)
			resBs := make([]string, 0)

			s := itrz.Zip(as, bs)
			for a, b := range s {
				resAs = append(resAs, a)
				resBs = append(resBs, b)
			}

			assert.Equal(t, max, len(resAs))
			assert.Equal(t, max, len(resBs))
		})
	}
}

func Test_ZipToShortest(t *testing.T) {
	tests := map[string]struct {
		as []int
		bs []string
	}{
		"as empty": {as: []int{}, bs: []string{"1", "2", "3"}},
		"as nil":   {as: nil, bs: []string{"1", "2", "3"}},
		"bs empty": {as: []int{1, 2, 3}, bs: []string{}},
		"bs nil":   {as: []int{1, 2, 3}, bs: nil},
		"same len": {as: []int{1, 2, 3}, bs: []string{"1", "2", "3"}},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			as := itrz.All(test.as)
			bs := itrz.All(test.bs)

			min := int(math.Min(float64(len(test.as)), float64(len(test.bs))))

			resAs := make([]int, 0)
			resBs := make([]string, 0)

			s := itrz.ZipToShortest(as, bs)
			for a, b := range s {
				resAs = append(resAs, a)
				resBs = append(resBs, b)
			}

			assert.Equal(t, min, len(resAs))
			assert.Equal(t, min, len(resBs))

			for idx := range min {
				assert.Equal(t, test.as[idx], resAs[idx])
				assert.Equal(t, test.bs[idx], resBs[idx])
			}
		})
	}
}

func Test_ZipStrict(t *testing.T) {
	tests := map[string]struct {
		as []int
		bs []string
	}{
		"as empty": {as: []int{}, bs: []string{"1", "2", "3"}},
		"as nil":   {as: nil, bs: []string{"1", "2", "3"}},
		"bs empty": {as: []int{1, 2, 3}, bs: []string{}},
		"bs nil":   {as: []int{1, 2, 3}, bs: nil},
		"same len": {as: []int{1, 2, 3}, bs: []string{"1", "2", "3"}},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			as := itrz.All(test.as)
			bs := itrz.All(test.bs)

			resAs := make([]int, 0)
			resBs := make([]string, 0)

			panics := len(test.as) != len(test.bs)

			if panics {
				assert.Panics(t, func() {
					for range itrz.ZipStrict(as, bs) {
					}
				})
			} else {
				s := itrz.ZipStrict(as, bs)
				for a, b := range s {
					resAs = append(resAs, a)
					resBs = append(resBs, b)
				}

				for idx := range test.as {
					assert.Equal(t, test.as[idx], resAs[idx])
					assert.Equal(t, test.bs[idx], resBs[idx])
				}
			}
		})
	}
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

func countElems[A any](s itrz.Seq[A]) int {
	n := 0
	for range s {
		n = n + 1
	}

	return n
}
