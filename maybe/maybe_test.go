package maybe_test

import (
	"testing"

	"github.com/dustin10/itrz/maybe"
	"github.com/stretchr/testify/assert"
)

func Test_Create(t *testing.T) {
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
