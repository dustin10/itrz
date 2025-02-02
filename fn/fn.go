package fn

// Predicate is a derived type that represents a function that takes a value of
// type A and returns a boolean.
type Predicate[A any] func(A) bool

// Consumer is a derived type that represents a function that takes a value of
// type A and has no return value.
type Consumer[A any] func(A)

// Factory is a derived type that represents a function that is capable of creating
// a value of type A without any arguments.
type Factory[A any] func() A

// Function is a derived type that represents a function that takes a value of type
// A as input and produces a value of type B as output.
type Function[A, B any] func(A) B
