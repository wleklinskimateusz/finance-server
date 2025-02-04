// option/option.go
package option

import "fmt"

// Option represents an optional value.
// If valid is true, then value is present (Some); otherwise, it is None.
type Option[T any] struct {
	value T
	valid bool
}

// Some creates an Option containing a value.
func Some[T any](value T) Option[T] {
	return Option[T]{value: value, valid: true}
}

// None creates an Option with no value.
func None[T any]() Option[T] {
	var zero T
	return Option[T]{value: zero, valid: false}
}

// IsSome returns true if the Option contains a value.
func (o Option[T]) IsSome() bool {
	return o.valid
}

// IsNone returns true if the Option does not contain a value.
func (o Option[T]) IsNone() bool {
	return !o.valid
}

// Unwrap returns the contained value if present; it panics if the Option is None.
func (o Option[T]) Unwrap() T {
	if !o.valid {
		panic("called Unwrap on a None Option")
	}
	return o.value
}

// UnwrapOr returns the contained value if present; otherwise, it returns defaultValue.
func (o Option[T]) UnwrapOr(defaultValue T) T {
	if o.valid {
		return o.value
	}
	return defaultValue
}

// Option represents an optional value.

// Map applies a function to the Option’s value, if present,
// returning a new Option holding the result.
func Map[T, U any](o Option[T], f func(T) U) Option[U] {
	if o.valid {
		return Some(f(o.value))
	}
	return None[U]()
}

// FlatMap applies a function that returns an Option to the contained value.
func FlatMap[T, U any](o Option[T], f func(T) Option[U]) Option[U] {
	if o.valid {
		return f(o.value)
	}
	return None[U]()
}

// Filter returns the Option itself if the contained value satisfies the predicate,
// otherwise it returns None.
func (o Option[T]) Filter(predicate func(T) bool) Option[T] {
	if o.valid && predicate(o.value) {
		return o
	}
	return None[T]()
}

// String implements the Stringer interface so that Option values can be printed nicely.
func (o Option[T]) String() string {
	if o.valid {
		return fmt.Sprintf("Some(%v)", o.value)
	}
	return "None"
}
