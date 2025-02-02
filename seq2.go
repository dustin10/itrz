package itrz

import "iter"

type Seq2[A, B any] iter.Seq2[A, B]

func FlatMap2[A, B, C any](seq Seq2[A, B], f func(A, B) Seq[C]) Seq[C] {
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

func Map2[A, B, C any](seq Seq2[A, B], f func(A, B) C) Seq[C] {
	return func(yield func(C) bool) {
		for a, b := range seq {
			if !yield(f(a, b)) {
				return
			}
		}
	}
}
