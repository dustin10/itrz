package itrz

import (
	"iter"

	"github.com/dustin10/itrz/fn"
	"github.com/dustin10/itrz/maybe"
)

type Seq[A any] iter.Seq[A]

func All[S ~[]A, A any](as S) Seq[A] {
	return func(yield func(a A) bool) {
		for _, a := range as {
			if !yield(a) {
				return
			}
		}
	}
}

func (s Seq[A]) AllMatch(p fn.Predicate[A]) bool {
	for a := range s {
		if !p(a) {
			return false
		}
	}

	return true
}

func (s Seq[A]) AnyMatch(p fn.Predicate[A]) bool {
	for a := range s {
		if p(a) {
			return true
		}
	}

	return false
}

func Concat[A any](seqs ...Seq[A]) Seq[A] {
	return func(yield func(A) bool) {
		for _, seq := range seqs {
			for a := range seq {
				if !yield(a) {
					return
				}
			}
		}
	}
}

func (s Seq[A]) Count() int {
	count := 0
	for range s {
		count = count + 1
	}

	return count
}

func Distinct[A comparable](seq Seq[A]) Seq[A] {
	set := make(map[A]struct{}, 0)

	return func(yield func(A) bool) {
		for a := range seq {
			_, exists := set[a]
			if exists {
				continue
			}

			set[a] = struct{}{}

			if !yield(a) {
				return
			}
		}
	}
}

func Empty[A any]() Seq[A] {
	return func(func(A) bool) {
		return
	}
}

func (s Seq[A]) Filter(p fn.Predicate[A]) Seq[A] {
	return func(yield func(a A) bool) {
		for a := range s {
			if p(a) && !yield(a) {
				return
			}
		}
	}
}

func (s Seq[A]) FindAny() maybe.Maybe[A] {
	next, stop := iter.Pull(iter.Seq[A](s))
	defer stop()

	if val, exists := next(); exists {
		return maybe.Just(val)
	}

	return maybe.Nothing[A]()
}

func FlatMap[A, B any](seq Seq[A], f func(A) Seq[B]) Seq[B] {
	return func(yield func(B) bool) {
		for a := range seq {
			mapped := f(a)
			for b := range mapped {
				if !yield(b) {
					return
				}
			}
		}
	}
}

func (s Seq[A]) ForEach(c fn.Consumer[A]) {
	for a := range s {
		c(a)
	}
}

func Generate[A any](f fn.Factory[A]) Seq[A] {
	return func(yield func(A) bool) {
		for {
			if !yield(f()) {
				return
			}
		}
	}
}

func (s Seq[A]) Limit(limit int) Seq[A] {
	return func(yield func(A) bool) {
		count := 0
		for a := range s {
			if count == limit || !yield(a) {
				return
			}

			count = count + 1
		}
	}
}

func Map[A, B any](seq Seq[A], f fn.Function[A, B]) Seq[B] {
	return func(yield func(B) bool) {
		for a := range seq {
			if !yield(f(a)) {
				return
			}
		}
	}
}

func (s Seq[A]) NoneMatch(p fn.Predicate[A]) bool {
	for a := range s {
		if p(a) {
			return false
		}
	}

	return true
}

func Of[A any](as ...A) Seq[A] {
	return func(yield func(A) bool) {
		for _, a := range as {
			if !yield(a) {
				return
			}
		}
	}
}

func (s Seq[A]) Peek(c fn.Consumer[A]) Seq[A] {
	return func(yield func(A) bool) {
		for a := range s {
			c(a)

			if !yield(a) {
				return
			}
		}
	}
}

func Reduce[A, B any](seq Seq[A], identity B, f func(A, B) B) B {
	result := identity

	for a := range seq {
		result = f(a, result)
	}

	return result
}

func (s Seq[A]) Skip(n int) Seq[A] {
	skipped := 0
	return func(yield func(A) bool) {
		for a := range s {
			if skipped < n {
				skipped = skipped + 1
				continue
			}

			if !yield(a) {
				return
			}
		}
	}
}

func (s Seq[A]) ToSlice() []A {
	as := make([]A, 0)
	for a := range s {
		as = append(as, a)
	}

	return as
}

func Zip[A, B any](as Seq[A], bs Seq[B]) Seq2[A, B] {
	return func(yield func(a A, b B) bool) {
		nextA, stopA := iter.Pull(iter.Seq[A](as))
		defer stopA()

		nextB, stopB := iter.Pull(iter.Seq[B](bs))
		defer stopB()

		for {
			a, existsA := nextA()
			b, existsB := nextB()

			if !existsA && !existsB {
				return
			}

			if !yield(a, b) {
				return
			}
		}
	}
}

func ZipStrict[A, B any](as Seq[A], bs Seq[B]) Seq2[A, B] {
	return func(yield func(a A, b B) bool) {
		nextA, stopA := iter.Pull(iter.Seq[A](as))
		defer stopA()

		nextB, stopB := iter.Pull(iter.Seq[B](bs))
		defer stopB()

		for {
			a, existsA := nextA()
			b, existsB := nextB()

			if !existsA && !existsB {
				return
			}

			if !existsA || !existsB {
				panic("sequences are not the same length")
			}

			if !yield(a, b) {
				return
			}
		}
	}
}
