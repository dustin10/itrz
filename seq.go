package itrz

import (
	"iter"

	"github.com/dustin10/itrz/maybe"
)

type Predicate[E any] func(E) bool

func All[S ~[]E, E any](es S) iter.Seq[E] {
	return func(yield func(e E) bool) {
		for _, e := range es {
			if !yield(e) {
				return
			}
		}
	}
}

func AllMatch[E any](seq iter.Seq[E], p Predicate[E]) bool {
	for e := range seq {
		if !p(e) {
			return false
		}
	}

	return true
}

func AnyMatch[E any](seq iter.Seq[E], p Predicate[E]) bool {
	for e := range seq {
		if p(e) {
			return true
		}
	}

	return false
}

func Concat[E any](seqs ...iter.Seq[E]) iter.Seq[E] {
	return func(yield func(E) bool) {
		for _, seq := range seqs {
			for elem := range seq {
				if !yield(elem) {
					return
				}
			}
		}
	}
}

func Count[E any](seq iter.Seq[E]) int {
	count := 0
	for range seq {
		count = count + 1
	}

	return count
}

func Distinct[E comparable](seq iter.Seq[E]) iter.Seq[E] {
	set := make(map[E]struct{}, 0)

	return func(yield func(E) bool) {
		for elem := range seq {
			_, exists := set[elem]
			if exists {
				continue
			}

			set[elem] = struct{}{}

			if !yield(elem) {
				return
			}
		}
	}
}

func Empty[E any]() iter.Seq[E] {
	return func(func(E) bool) {
		return
	}
}

func Filter[E any](seq iter.Seq[E], p Predicate[E]) iter.Seq[E] {
	return func(yield func(e E) bool) {
		for e := range seq {
			if p(e) && !yield(e) {
				return
			}
		}
	}
}

func FindAny[E any](seq iter.Seq[E]) maybe.Maybe[E] {
	next, stop := iter.Pull(seq)
	defer stop()

	if val, exists := next(); exists {
		return maybe.Just(val)
	}

	return maybe.Nothing[E]()
}

func FlatMap[I, O any](seq iter.Seq[I], fn func(I) iter.Seq[O]) iter.Seq[O] {
	return func(yield func(O) bool) {
		for elem := range seq {
			innerSeq := fn(elem)
			for innerElem := range innerSeq {
				if !yield(innerElem) {
					return
				}
			}
		}
	}
}

func ForEach[E any](seq iter.Seq[E], fn func(E)) {
	for elem := range seq {
		fn(elem)
	}
}

func Generate[E any](fn func() E) iter.Seq[E] {
	return func(yield func(E) bool) {
		for {
			if !yield(fn()) {
				return
			}
		}
	}
}

func Limit[E any](seq iter.Seq[E], limit int) iter.Seq[E] {
	return func(yield func(E) bool) {
		count := 0
		for elem := range seq {
			if count == limit || !yield(elem) {
				return
			}

			count = count + 1
		}
	}
}

func Map[I, O any](seq iter.Seq[I], fn func(I) O) iter.Seq[O] {
	return func(yield func(O) bool) {
		for elem := range seq {
			if !yield(fn(elem)) {
				return
			}
		}
	}
}

func NoneMatch[E any](seq iter.Seq[E], p Predicate[E]) bool {
	for e := range seq {
		if p(e) {
			return false
		}
	}

	return true
}

func Of[E any](es ...E) iter.Seq[E] {
	return func(yield func(E) bool) {
		for _, elem := range es {
			if !yield(elem) {
				return
			}
		}
	}
}

func Peek[E any](seq iter.Seq[E], fn func(E)) iter.Seq[E] {
	return func(yield func(E) bool) {
		for elem := range seq {
			fn(elem)

			if !yield(elem) {
				return
			}
		}
	}
}

func Skip[E any](seq iter.Seq[E], n int) iter.Seq[E] {
	skipped := 0
	return func(yield func(E) bool) {
		for elem := range seq {
			if skipped < n {
				skipped = skipped + 1
				continue
			}

			if !yield(elem) {
				return
			}
		}
	}
}

func ToSlice[E any](seq iter.Seq[E]) []E {
	es := make([]E, 0)
	for e := range seq {
		es = append(es, e)
	}

	return es
}

func Zip[E, F any](es iter.Seq[E], fs iter.Seq[F]) iter.Seq2[E, F] {
	return func(yield func(e E, f F) bool) {
		nextE, stopE := iter.Pull(es)
		defer stopE()

		nextF, stopF := iter.Pull(fs)
		defer stopF()

		for {
			e, existsE := nextE()
			f, existsF := nextF()

			if !existsE && !existsF {
				return
			}

			if !yield(e, f) {
				return
			}
		}
	}
}

func ZipStrict[E, F any](es iter.Seq[E], fs iter.Seq[F]) iter.Seq2[E, F] {
	return func(yield func(e E, f F) bool) {
		nextE, stopE := iter.Pull(es)
		defer stopE()

		nextF, stopF := iter.Pull(fs)
		defer stopF()

		for {
			e, existsE := nextE()
			f, existsF := nextF()

			if !existsE && !existsF {
				return
			}

			if !existsE || !existsF {
				panic("sequences are not the same length")
			}

			if !yield(e, f) {
				return
			}
		}
	}
}
