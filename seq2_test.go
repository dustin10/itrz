package itrz_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dustin10/itrz"
)

func Test_FlatMap2(t *testing.T) {
	f := func(a int, b string) itrz.Seq[string] {
		return itrz.Of(fmt.Sprintf("%d: %s", a, b))
	}

	tests := map[string]struct {
		as []int
		bs []string
		fn func(int, string) itrz.Seq[string]
	}{
		"as empty": {as: []int{}, bs: []string{"1", "2", "3"}, fn: f},
		"as nil":   {as: nil, bs: []string{"1", "2", "3"}, fn: f},
		"bs empty": {as: []int{1, 2, 3}, bs: []string{}, fn: f},
		"bs nil":   {as: []int{1, 2, 3}, bs: nil, fn: f},
		"same len": {as: []int{1, 2, 3}, bs: []string{"1", "2", "3"}, fn: f},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			as := itrz.All(test.as)
			bs := itrz.All(test.bs)

			pairs := itrz.ZipToShortest(as, bs)

			s := itrz.FlatMap2(pairs, f)

			i := 0
			for c := range s {
				assert.Equal(t, fmt.Sprintf("%d: %s", test.as[i], test.bs[i]), c)
				i = i + 1
			}
		})
	}
}

func Test_Map2(t *testing.T) {
	f := func(a int, b string) string {
		return fmt.Sprintf("%d: %s", a, b)
	}

	tests := map[string]struct {
		as []int
		bs []string
		fn func(int, string) string
	}{
		"as empty": {as: []int{}, bs: []string{"1", "2", "3"}, fn: f},
		"as nil":   {as: nil, bs: []string{"1", "2", "3"}, fn: f},
		"bs empty": {as: []int{1, 2, 3}, bs: []string{}, fn: f},
		"bs nil":   {as: []int{1, 2, 3}, bs: nil, fn: f},
		"same len": {as: []int{1, 2, 3}, bs: []string{"1", "2", "3"}, fn: f},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			as := itrz.All(test.as)
			bs := itrz.All(test.bs)

			pairs := itrz.ZipToShortest(as, bs)

			s := itrz.Map2(pairs, f)

			i := 0
			for c := range s {
				assert.Equal(t, f(test.as[i], test.bs[i]), c)
				i = i + 1
			}
		})
	}

}
