package set

import (
	"iter"

	"github.com/dustin10/itrz"
)

const defaultCapacity = 16

type Set[E comparable] struct {
	elems map[E]struct{}
}

func New[E comparable]() Set[E] {
	return NewWithCapacity[E](defaultCapacity)
}

func NewWithCapacity[E comparable](capacity int) Set[E] {
	return Set[E]{
		elems: make(map[E]struct{}, capacity),
	}
}

func (s *Set[E]) Add(elem E) {
	s.elems[elem] = struct{}{}
}

func (s *Set[E]) Remove(elem E) bool {
	_, exists := s.elems[elem]
	if exists {
		delete(s.elems, elem)
	}

	return exists
}

func (s *Set[E]) Contains(elem E) bool {
	return itrz.AnyMatch[E](s.All(), func(e E) bool {
		return e == elem
	})
}

func (s *Set[E]) Clear() int {
	num := len(s.elems)
	s.elems = make(map[E]struct{}, defaultCapacity)

	return num
}

func (s *Set[E]) All() iter.Seq[E] {
	return func(yield func(E) bool) {
		for e := range s.elems {
			if !yield(e) {
				return
			}
		}
	}
}

func (s *Set[E]) Filter(p itrz.Predicate[E]) iter.Seq[E] {
	seq := func(yield func(e E) bool) {
		for e := range s.elems {
			if !yield(e) {
				return
			}
		}
	}

	return itrz.Filter(seq, p)
}
