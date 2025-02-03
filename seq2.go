package itrz

import (
	"iter"

	"github.com/dustin10/itrz/fn"
)

// Seq2 is a derived type of iter.Seq2.
type Seq2[A, B any] iter.Seq2[A, B]

// FlatMap2 applies the fn.Function2, that itself returns a Seq, to each tuple element
// yielded by the Seq2 and flattens them out into one Seq.
func FlatMap2[A, B, C any](seq Seq2[A, B], f fn.Function2[A, B, Seq[C]]) Seq[C] {
	return func(yield func(C) bool) {
		for a, b := range seq {
			mapped := f(a, b)
			for c := range mapped {
				if !yield(c) {
					return
				}
			}
		}
	}
}

// Map2 returns a new Seq2 consisting of the results of applying the given fn.Function2 to
// the elements of the existing Seq2.
func Map2[A, B, C any](seq Seq2[A, B], f fn.Function2[A, B, C]) Seq[C] {
	return func(yield func(C) bool) {
		for a, b := range seq {
			if !yield(f(a, b)) {
				return
			}
		}
	}
}
