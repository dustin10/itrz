package maybe

import "github.com/dustin10/itrz/fn"

type Maybe[A any] struct {
	value   A
	present bool
}

func Just[A any](value A) Maybe[A] {
	return Maybe[A]{
		value:   value,
		present: true,
	}
}

func Nothing[A any]() Maybe[A] {
	return Maybe[A]{
		present: false,
	}
}

func (m *Maybe[A]) Filter(p fn.Predicate[A]) Maybe[A] {
	if m.IsEmpty() {
		return Nothing[A]()
	}

	value := m.Get()

	if p(value) {
		return Just(value)
	}

	return Nothing[A]()
}

func (m *Maybe[A]) Get() A {
	if !m.present {
		panic("Get() called with no value is present")
	}

	return m.value
}

func (m *Maybe[A]) IsPresent() bool {
	return m.present
}

func (m *Maybe[A]) IsEmpty() bool {
	return !m.IsPresent()
}

func (m *Maybe[A]) Or(value A) A {
	if m.present {
		return m.value
	}

	return value
}

func (m *Maybe[A]) OrElse(f fn.Factory[A]) A {
	if m.present {
		return m.value
	}

	return f()
}

func FlatMap[A, B any](m Maybe[A], fn func(A) Maybe[B]) Maybe[B] {
	if m.IsEmpty() {
		return Nothing[B]()
	}

	return fn(m.Get())
}

func Map[A, B any](m Maybe[A], f fn.Function[A, B]) Maybe[B] {
	if m.IsEmpty() {
		return Nothing[B]()
	}

	return Just(f(m.Get()))
}
