package set

import (
	"encoding/json"
	"fmt"
	"maps"
	"strings"

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

// FromSlice creates a new Set using the given slice as the initial data and applying
// any Options that are specified.
func FromSlice[S ~[]A, A comparable](as S, opts ...Option) Set[A] {
	s := New[A](opts...)

	for _, a := range as {
		s.Add(a)
	}

	return s
}

func create[A comparable](config Config) Set[A] {
	return Set[A]{
		config: config,
		elems:  make(map[A]struct{}, config.InitialCapacity),
	}
}

// IsEmpty returns true if the Set as zero elements and false otherwise.
func (s *Set[A]) IsEmpty() bool {
	return s.Len() == 0
}

// Len returns the number of elements in the Set.
func (s *Set[A]) Len() int {
	return len(s.elems)
}

// Add adds an element to the Set.
func (s *Set[A]) Add(a A) {
	s.elems[a] = struct{}{}
}

// Remove removes the specified element from the Set. Returns true if the value was
// removed from the Set.
func (s *Set[A]) Remove(a A) bool {
	_, exists := s.elems[a]
	if exists {
		delete(s.elems, a)
	}

	return exists
}

// Contains returns true if the Set contains the specified element or false otherwise.
func (s *Set[A]) Contains(a A) bool {
	_, exists := s.elems[a]

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

// String returns a string representation of the Maybe.
func (s Set[A]) String() string {
	f := func(a A) string {
		return fmt.Sprintf("%v", a)
	}

	as := itrz.Map(s.All(), f).ToSlice()

	return fmt.Sprintf("[%s]", strings.Join(as, ","))
}

// MarshalJSON converts the Set to it's JSON representation.
func (s Set[A]) MarshalJSON() ([]byte, error) {
	as := s.All().ToSlice()

	bytes, err := json.Marshal(&as)
	if err != nil {
		return nil, fmt.Errorf("marshal Set to JSON: %w", err)
	}

	return bytes, nil
}

// UnmarshalJSON converts the JSON bytes to the value contained in the Maybe if present.
func (s *Set[A]) UnmarshalJSON(data []byte) error {
	as := make([]A, 0)

	err := json.Unmarshal(data, &as)
	if err != nil {
		return fmt.Errorf("unmarshal JSON to Set: %w", err)
	}

	for _, a := range as {
		s.Add(a)
	}

	return nil
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
