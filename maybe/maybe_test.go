package maybe_test

import (
	"testing"

	"github.com/dustin10/itrz/fn"
	"github.com/dustin10/itrz/maybe"
	"github.com/stretchr/testify/assert"
)

func Test_CreateAndPresence(t *testing.T) {
	value := "value"

	j := maybe.Just(value)

	assert.True(t, j.IsPresent(), "expected value present for Just")
	assert.False(t, j.IsEmpty(), "expected value present for Just")
	assert.Equal(t, value, j.Get(), "unexpected value for Just")

	n := maybe.Nothing[string]()

	assert.True(t, n.IsEmpty(), "expected value present for Nothing")
	assert.False(t, n.IsPresent(), "expected value present for Nothing")
	assert.Panics(t, func() { n.Get() }, "expected Get() to panic for Nothing")
}

func Test_FromPointer(t *testing.T) {

}

func Test_FromString(t *testing.T) {

}

func Test_Filter(t *testing.T) {
	tests := map[string]struct {
		value     any
		predicate fn.Predicate[any]
		expected  bool
	}{
		"nothing":              {value: nil, predicate: pass},
		"value matches":        {value: "value", predicate: pass, expected: true},
		"value does not match": {value: "value", predicate: fail},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var m maybe.Maybe[any]
			if test.value != nil {
				m = maybe.Just(test.value)
			} else {
				m = maybe.Nothing[any]()
			}

			res := m.Filter(test.predicate)

			assert.Equal(t, test.expected, res.IsPresent())
		})
	}
}

func Test_Get(t *testing.T) {
	tests := map[string]struct {
		value any
	}{
		"nil":   {},
		"value": {value: "value"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var m maybe.Maybe[any]
			if test.value != nil {
				m = maybe.Just(test.value)
			} else {
				m = maybe.Nothing[any]()
			}

			if test.value == nil {
				assert.Panics(t, func() { m.Get() })
			} else {
				assert.Equal(t, test.value, m.Get())
			}
		})
	}
}

func Test_Or(t *testing.T) {
	tests := map[string]struct {
		value  any
		dflt   any
		expect any
	}{
		"nil":   {dflt: "default", expect: "default"},
		"value": {value: "value", expect: "value"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var m maybe.Maybe[any]
			if test.value != nil {
				m = maybe.Just(test.value)
			} else {
				m = maybe.Nothing[any]()
			}

			assert.Equal(t, test.expect, m.Or(test.dflt))
		})
	}
}

func Test_OrElse(t *testing.T) {
	tests := map[string]struct {
		value  any
		f      fn.Factory[any]
		expect any
	}{
		"nil":   {f: func() any { return "default" }, expect: "default"},
		"value": {value: "value", expect: "value"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var m maybe.Maybe[any]
			if test.value != nil {
				m = maybe.Just(test.value)
			} else {
				m = maybe.Nothing[any]()
			}

			assert.Equal(t, test.expect, m.OrElse(test.f))
		})
	}
}

func Test_FlatMap(t *testing.T) {
	tests := map[string]struct {
		value  string
		f      fn.Function[string, maybe.Maybe[string]]
		expect string
	}{
		"Just when Function maps to Just":    {value: "value", f: func(string) maybe.Maybe[string] { return maybe.Just("result") }, expect: "result"},
		"Just when Function maps to Nothing": {f: func(string) maybe.Maybe[string] { return maybe.Nothing[string]() }, expect: "default"},
		"Nothing always maps to Nothing":     {expect: "default"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var m maybe.Maybe[string]
			if len(test.value) != 0 {
				m = maybe.Just(test.value)
			} else {
				m = maybe.Nothing[string]()
			}

			assert.Equal(t, test.expect, maybe.FlatMap(m, test.f).Or("default"))
		})
	}
}

func Test_Map(t *testing.T) {
	dflt := -1

	tests := map[string]struct {
		value  string
		f      fn.Function[string, int]
		expect int
	}{
		"Just should map to Just":       {value: "value", f: strlen, expect: strlen("value")},
		"Nothing should map to Nothing": {expect: dflt},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var m maybe.Maybe[string]
			if len(test.value) != 0 {
				m = maybe.Just(test.value)
			} else {
				m = maybe.Nothing[string]()
			}

			assert.Equal(t, test.expect, maybe.Map(m, test.f).Or(dflt))
		})
	}
}

func pass(any) bool {
	return true
}

func fail(any) bool {
	return false
}

func strlen(s string) int {
	return len(s)
}
