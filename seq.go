package itrz

import (
	"iter"

	"github.com/dustin10/itrz/fn"
	"github.com/dustin10/itrz/maybe"
)

// Seq is a type derived from iter.Seq.
type Seq[A any] iter.Seq[A]

// All returns a Seq from the specified slice.
func All[S ~[]A, A any](as S) Seq[A] {
	return func(yield func(a A) bool) {
		for _, a := range as {
			if !yield(a) {
				return
			}
		}
	}
}

// AllMatch returns true if all elements in the Seq match the specified fn.Predicate.
func (s Seq[A]) AllMatch(p fn.Predicate[A]) bool {
	for a := range s {
		if !p(a) {
			return false
		}
	}

	return true
}

// AnyMatch returns true if any of the elements in the Seq match the specified fn.Predicate.
func (s Seq[A]) AnyMatch(p fn.Predicate[A]) bool {
	for a := range s {
		if p(a) {
			return true
		}
	}

	return false
}

// Concat concatenates the specified sequences together into one Seq.
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

// Count returns the number of elements yielded by the Seq.
func (s Seq[A]) Count() int {
	count := 0
	for range s {
		count = count + 1
	}

	return count
}

// Distinct takes the elements from the specified Seq and returns a new Seq that will only
// yield the distince values from the original one.
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

// Sink is the interface that the accepting data structure must implement in order to have
// a Seq be drained into it.
type Sink[A any] interface {
	// Add adds an element of type A.
	Add(A)
}

// DrainTo iterates all elements of the Seq and calls the Sink.Add function for each one.
func (s Seq[A]) DrainTo(sink Sink[A]) {
	s.ForEach(sink.Add)
}

// Empty returns a Seq which yields zero elements.
func Empty[A any]() Seq[A] {
	return func(func(A) bool) {
		return
	}
}

// Filter returns a Seq that only yields elements from the original Seq that match the
// specified fn.Predicate.
func (s Seq[A]) Filter(p fn.Predicate[A]) Seq[A] {
	return func(yield func(a A) bool) {
		for a := range s {
			if p(a) && !yield(a) {
				return
			}
		}
	}
}

// FindAny returns a maybe.Maybe that contains some element of the Seq, or is empty if the
// Seq has no elements.
func (s Seq[A]) FindAny() maybe.Maybe[A] {
	next, stop := iter.Pull(iter.Seq[A](s))
	defer stop()

	if val, exists := next(); exists {
		return maybe.Just(val)
	}

	return maybe.Nothing[A]()
}

// FlatMap applies the function, that itself returns a Seq, to each element yielded by the
// Seq and flattens them out into one Seq.
func FlatMap[A, B any](seq Seq[A], f fn.Function[A, Seq[B]]) Seq[B] {
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

// ForEach applies the given fn.Consumer to each element yielded by the Seq.
func (s Seq[A]) ForEach(c fn.Consumer[A]) {
	for a := range s {
		c(a)
	}
}

// Generate returns a Seq that will yield elements produced by the specified fn.Factory
// function.
func Generate[A any](f fn.Factory[A]) Seq[A] {
	return func(yield func(A) bool) {
		for {
			if !yield(f()) {
				return
			}
		}
	}
}

// Limit returns a new Seq that will only yield limit number of elements.
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

// Map returns a new Seq consisting of the results of applying the given fn.Function to
// the elements of the existing Seq.
func Map[A, B any](seq Seq[A], f fn.Function[A, B]) Seq[B] {
	return func(yield func(B) bool) {
		for a := range seq {
			if !yield(f(a)) {
				return
			}
		}
	}
}

// NoneMatch returns true if no element yielded by the Seq matches the fn.Predicate.
func (s Seq[A]) NoneMatch(p fn.Predicate[A]) bool {
	for a := range s {
		if p(a) {
			return false
		}
	}

	return true
}

// Of returns a Seq that yields the specified elements.
func Of[A any](as ...A) Seq[A] {
	return func(yield func(A) bool) {
		for _, a := range as {
			if !yield(a) {
				return
			}
		}
	}
}

// Peek returns a Seq consisting of the elements of this Seq, additionally performing the
// provided action on each element as elements are consumed from the resulting Seq.
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

// Reduce performs a reduction on the elements of the Seq, using the provided identity value
// and an associative accumulation function, and returns the reduced value.
func Reduce[A, B any](seq Seq[A], identity B, f fn.Function2[A, B, B]) B {
	result := identity

	for a := range seq {
		result = f(a, result)
	}

	return result
}

// Skip returns a Seq consisting of the remaining elements of the existing Seq after discarding
// the first n elements.
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

// ToSlice returns a slice containing the elements of the Seq.
func (s Seq[A]) ToSlice() []A {
	as := make([]A, 0)
	for a := range s {
		as = append(as, a)
	}

	return as
}

// Zip combines the elements of the two specified Seq instances into a Seq2 that contains
// the pair-wise tuples. If the two sequences are not of the same length then the Seq2 will
// stop yielding tuple elements once the first Seq is exhausted.
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

// ZipStrict combines the elements of the two specified Seq instances into a Seq2 that contains
// the pair-wise tuples. If the two sequences are not of the same length then the Seq2 will
// panic while yielding the tuple elements.
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
