package fn

type Predicate[A any] func(A) bool

type Consumer[A any] func(A)

type Factory[A any] func() A

type Function[A, B any] func(A) B
