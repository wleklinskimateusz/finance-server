package optional

import "errors"

type Option[T any] interface {
	Get() (*T, error)
	IsPresent() bool
	IfPresent(consumer func(val T)) Option[T]
	OrElse(other T) T
	OrElseGet(otherFunc func() T) T
	Filter(predicate func(val T) bool) Option[T]
}

type option[T any] struct {
	val *T // pointer to the type, so that NIL can represented
}

func Of[T any](val T) Option[T] {
	return &option[T]{&val}
}

func Empty[T any]() Option[T] {
	return &option[T]{}
}

func (option *option[T]) Get() (*T, error) {
	if option.IsPresent() {
		return option.val, nil
	}
	return nil, errors.New("no value present")
}

func (option *option[T]) IsPresent() bool {
	return option.val != nil
}

func (option *option[T]) IfPresent(consumer func(val T)) Option[T] {
	val, _ := option.Get()
	if val != nil {
		consumer(*val)
	}
	return option
}

func (option *option[T]) OrElse(other T) T {
	val, _ := option.Get()
	if val != nil {
		return *val
	}
	return other
}

func (option *option[T]) OrElseGet(otherFunc func() T) T {
	val, _ := option.Get()
	if val != nil {
		return *val
	}
	return otherFunc()
}

func (option *option[T]) Filter(predicate func(val T) bool) Option[T] {
	val, _ := option.Get()
	if val != nil && predicate(*val) {
		return Of[T](*val)
	}
	return Empty[T]()
}

func Map[T any, S any](option Option[T], mapper func(val T) S) Option[S] {
	val, _ := option.Get()
	if val != nil {
		result := mapper(*val)
		return Of[S](result)
	}
	return Empty[S]()
}
