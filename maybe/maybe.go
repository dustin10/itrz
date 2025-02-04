package maybe

import (
	"encoding/json"
	"fmt"

	"github.com/dustin10/itrz/fn"
)

// Maybe encapsulates a value that is optional. Either a value exists or it does not. This struct
// is similar to Maybe in Haskell or Option in Rust.
type Maybe[A any] struct {
	value   A
	present bool
}

// Just creates a new Maybe with a value of type A present.
func Just[A any](value A) Maybe[A] {
	return Maybe[A]{
		value:   value,
		present: true,
	}
}

// Nothing creates a Maybe with no value present.
func Nothing[A any]() Maybe[A] {
	return Maybe[A]{
		present: false,
	}
}

// FromPointer creates a new Maybe of type A. If the pointer is nil then the Maybe will be empty,
// otherwise the Maybe will contain the value the pointer is pointing to.
func FromPointer[A any](p *A) Maybe[A] {
	if p == nil {
		return Nothing[A]()
	}

	return Just(*p)
}

// FromString creates a new Maybe of type string. If the string is empty then the Maybe will be
// empty, otherwise the Maybe will contain the given string value.
func FromString(s string) Maybe[string] {
	if len(s) == 0 {
		return Nothing[string]()
	}

	return Just(s)
}

// Filter applies the given Predicate to the value contained in the Maybe. If there is no value
// present, then the function is not applied and an empty Maybe is returned.
func (m Maybe[A]) Filter(p fn.Predicate[A]) Maybe[A] {
	if m.IsEmpty() {
		return m
	}

	value := m.Get()

	if p(value) {
		return Just(value)
	}

	return Nothing[A]()
}

// Get returns the value of type A contained in the Maybe. If this function is invoked an empty
// Maybe then it will cause a panic.
func (m Maybe[A]) Get() A {
	if !m.present {
		panic("Get() called with no value is present")
	}

	return m.value
}

// IsPresent returns true if the Maybe contains a value and false otherwise.
func (m Maybe[A]) IsPresent() bool {
	return m.present
}

// IsEmpty returns true if the Maybe does not contain a value and false otherwise.
func (m Maybe[A]) IsEmpty() bool {
	return !m.IsPresent()
}

// Or returns the value contained in the Maybe if it exists, otherwise it returns
// the specified value.
func (m Maybe[A]) Or(value A) A {
	if m.present {
		return m.value
	}

	return value
}

// OrElse returns the value contained in the Maybe if it exists, otherwise it returns
// the value returned by the specified Factory function.
func (m Maybe[A]) OrElse(f fn.Factory[A]) A {
	if m.present {
		return m.value
	}

	return f()
}

// String returns a string representation of the Maybe.
func (m Maybe[A]) String() string {
	if m.present {
		return fmt.Sprintf("Just(%v)", m.value)
	} else {
		return "Nothing"
	}
}

// MarshalJSON converts the value in the Maybe, if present, to it's JSON representation.
// If not value is present then the JSON representation is null.
func (m Maybe[_]) MarshalJSON() ([]byte, error) {
	if m.present {
		return json.Marshal(m.value)
	} else {
		return []byte("null"), nil
	}
}

// UnmarshalJSON converts the JSON bytes to the value contained in the Maybe if present.
func (m *Maybe[_]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		m.present = false
		return nil
	}

	err := json.Unmarshal(data, &m.value)
	if err != nil {
		return fmt.Errorf("unmarshal Maybe value from JSON: %w", err)
	}

	m.present = true

	return nil
}

// FlatMap applies the given Function to the value in the Maybe if it exists.
func FlatMap[A, B any](m Maybe[A], f fn.Function[A, Maybe[B]]) Maybe[B] {
	if m.IsEmpty() {
		return Nothing[B]()
	}

	return f(m.Get())
}

// Map applies the given Function to the value in the Maybe if it exists.
func Map[A, B any](m Maybe[A], f fn.Function[A, B]) Maybe[B] {
	if m.IsEmpty() {
		return Nothing[B]()
	}

	return Just(f(m.Get()))
}
