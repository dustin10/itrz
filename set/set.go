package set

import (
	"maps"

	"github.com/dustin10/itrz"
	"github.com/dustin10/itrz/fn"
)

// defaultCapacity defines the default initial capacity for a Set.
const defaultInitialCapacity = 16

// Option defines a function that can be used to custmize the configuration used to
// create a Set.
type Option func(config *Config)

// Config contains the supported configuration of a Set.
type Config struct {
	// Capacity defines the initial size of the Set.
	InitialCapacity int
}

// WithInitialCapacity is an Option that can be used to configure the initial capacity
// of a Set.
func WithInitialCapacity(capacity int) Option {
	return func(config *Config) {
		config.InitialCapacity = capacity
	}
}

// Set is a collection that contains no duplicate elements.
type Set[A comparable] struct {
	config Config
	elems  map[A]struct{}
}

// New creates a new Set applying any Options that are specified.
func New[A comparable](opts ...Option) Set[A] {
	config := Config{
		InitialCapacity: defaultInitialCapacity,
	}

	for _, opt := range opts {
		opt(&config)
	}

	return create[A](config)

}

func create[A comparable](config Config) Set[A] {
	return Set[A]{
		config: config,
		elems:  make(map[A]struct{}, config.InitialCapacity),
	}
}

// Add adds an element to the Set.
func (s *Set[A]) Add(elem A) {
	s.elems[elem] = struct{}{}
}

// Remove removes the specified element from the Set. Returns true if the value was
// removed from the Set.
func (s *Set[A]) Remove(elem A) bool {
	_, exists := s.elems[elem]
	if exists {
		delete(s.elems, elem)
	}

	return exists
}

// Contains returns true if the Set contains the specified element or false otherwise.
func (s *Set[A]) Contains(elem A) bool {
	_, exists := s.elems[elem]

	return exists
}

// Clear removes all values in the Set.
func (s *Set[A]) Clear() int {
	num := len(s.elems)

	*s = create[A](s.config)

	return num
}

// All returns an itrz.Seq that can be used to range over the Set.
func (s *Set[A]) All() itrz.Seq[A] {
	return itrz.Seq[A](maps.Keys(s.elems))
}

// FlatMap applies a given function, that itself returns a Set, to each value in the Set and
// returns a new Set with the flattened values.
func FlatMap[A, B comparable](s Set[A], f fn.Function[A, Set[B]]) Set[B] {
	result := create[B](Config{
		InitialCapacity: len(s.elems),
	})

	mapper := func(a A) itrz.Seq[B] {
		s := f(a)
		return s.All()
	}

	itrz.FlatMap(s.All(), mapper).DrainTo(&result)

	return result
}

// Map applies a given function to each value in the Set returning a new Set with the
// mapped values.
func Map[A, B comparable](s Set[A], f fn.Function[A, B]) Set[B] {
	result := create[B](Config{
		InitialCapacity: len(s.elems),
	})

	itrz.Map(s.All(), f).DrainTo(&result)

	return result
}
