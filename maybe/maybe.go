package maybe

type Maybe[E any] struct {
	value     E
	isPresent bool
}

func Just[E any](value E) Maybe[E] {
	return Maybe[E]{
		value:     value,
		isPresent: true,
	}
}

func Nothing[E any]() Maybe[E] {
	return Maybe[E]{
		isPresent: false,
	}
}

func (m *Maybe[E]) Filter(fn func(E) bool) Maybe[E] {
	if m.IsEmpty() {
		return Nothing[E]()
	}

	value := m.Get()

	if fn(value) {
		return Just(value)
	}

	return Nothing[E]()
}

func (m *Maybe[E]) Get() E {
	if !m.isPresent {
		panic("Get() called with no value is present")
	}

	return m.value
}

func (m *Maybe[E]) IsPresent() bool {
	return m.isPresent
}

func (m *Maybe[E]) IsEmpty() bool {
	return !m.IsPresent()
}

func (m *Maybe[E]) Or(value E) E {
	if m.isPresent {
		return m.value
	}

	return value
}

func (m *Maybe[E]) OrElse(fn func() E) E {
	if m.isPresent {
		return m.value
	}

	return fn()
}

func FlatMap[I, O any](m Maybe[I], fn func(I) Maybe[O]) Maybe[O] {
	if m.IsEmpty() {
		return Nothing[O]()
	}

	return fn(m.Get())
}

func Map[I, O any](m Maybe[I], fn func(I) O) Maybe[O] {
	if m.IsEmpty() {
		return Nothing[O]()
	}

	return Just(fn(m.Get()))
}
