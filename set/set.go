package set

import (
	"maps"

	"github.com/dustin10/itrz"
	"github.com/dustin10/itrz/fn"
)

const defaultCapacity = 16

type Option func(config *Config)

type Config struct {
	Capacity int
}

func WithCapacity(capacity int) Option {
	return func(config *Config) {
		config.Capacity = capacity
	}
}

type Set[A comparable] struct {
	config Config
	elems  map[A]struct{}
}

func New[A comparable](opts ...Option) Set[A] {
	config := Config{
		Capacity: defaultCapacity,
	}

	for _, opt := range opts {
		opt(&config)
	}

	return create[A](config)

}

func create[A comparable](config Config) Set[A] {
	return Set[A]{
		config: config,
		elems:  make(map[A]struct{}, config.Capacity),
	}
}

func (s *Set[A]) Add(elem A) {
	s.elems[elem] = struct{}{}
}

func (s *Set[A]) Remove(elem A) bool {
	_, exists := s.elems[elem]
	if exists {
		delete(s.elems, elem)
	}

	return exists
}

func (s *Set[A]) Contains(elem A) bool {
	_, exists := s.elems[elem]

	return exists
}

func (s *Set[A]) Clear() int {
	num := len(s.elems)

	*s = create[A](s.config)

	return num
}

func (s *Set[A]) All() itrz.Seq[A] {
	return itrz.Seq[A](maps.Keys(s.elems))
}

func Map[A, B comparable](s Set[A], f fn.Function[A, B]) Set[B] {
	result := create[B](s.config)

	s.All().ForEach(func(a A) {
		result.Add(f(a))
	})

	return result
}
