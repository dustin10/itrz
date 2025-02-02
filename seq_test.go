package itrz_test

import (
	"testing"

	"github.com/dustin10/itrz"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/constraints"
)

func Test_Seq_Reduce(t *testing.T) {
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

type Number interface {
	constraints.Integer | constraints.Float
}

func sum[A Number](a A, b A) A {
	return a + b
}
