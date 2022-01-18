package batches

import (
	"github.com/mcmathja/funky/pairs"
)

// Batch
type Batch[T any] func(next func(T) bool)

/* Constructors */

func FromChan[T any](ch <-chan T) Batch[T] {
	return func(next func(T) bool) {
		for ele := range ch {
			if !next(ele) {
				break
			}
		}
	}
}

func FromMap[K comparable, V any](m map[K]V) Batch[pairs.Pair[K, V]] {
	return func(next func(pairs.Pair[K, V]) bool) {
		for k, v := range m {
			if !next(pairs.New(k, v)) {
				break
			}
		}
	}
}

func FromSet[T comparable](s map[T]struct{}) Batch[T] {
	return func(next func(T) bool) {
		for ele := range s {
			if !next(ele) {
				break
			}
		}
	}
}

func FromSlice[T any](s []T) Batch[T] {
	return func(next func(T) bool) {
		for _, ele := range s {
			if !next(ele) {
				break
			}
		}
	}
}

func New[T any](eles ...T) Batch[T] {
	return func(next func(T) bool) {
		for _, ele := range eles {
			if !next(ele) {
				break
			}
		}
	}
}

/* Operations */

func Append[T any](b Batch[T], ele T) Batch[T] {
	return func(next func(T) bool) {
		b(func(in T) bool {
			return next(in)
		})
		next(ele)
	}
}

func Distinct[T comparable](b Batch[T]) Batch[T] {
	return func(next func(T) bool) {
		seen := make(map[T]struct{}, 0)
		b(func(ele T) bool {
			if _, ok := seen[ele]; !ok {
				seen[ele] = struct{}{}
				return next(ele)
			}
			return true
		})
	}
}

func DistinctBy[T any, U comparable](b Batch[T], fn func(T) U) Batch[T] {
	return func(next func(T) bool) {
		seen := make(map[U]struct{}, 0)
		b(func(ele T) bool {
			comp := fn(ele)
			if _, ok := seen[comp]; !ok {
				seen[comp] = struct{}{}
				return next(ele)
			}
			return true
		})
	}
}

func Drop[T any](b Batch[T], num int) Batch[T] {
	return func(next func(T) bool) {
		b(func(in T) bool {
			if num > 0 {
				num--
			}
			if num <= 0 {
				return next(in)
			}
			return true
		})
	}
}

func DropWhile[T any](b Batch[T], fn func(T) bool) Batch[T] {
	return func(next func(T) bool) {
		done := false
		b(func(in T) bool {
			if !done && fn(in) {
				done = true
			}
			if done {
				return next(in)
			}
			return true
		})
	}
}

func Filter[T any](b Batch[T], fn func(T) bool) Batch[T] {
	return func(next func(T) bool) {
		b(func(in T) bool {
			if fn(in) {
				return next(in)
			}
			return true
		})
	}
}

func FlatMap[T, U any](b Batch[T], fn func(T) []U) Batch[U] {
	return func(next func(U) bool) {
		b(func(in T) bool {
			for _, ele := range fn(in) {
				if !next(ele) {
					return false
				}
			}
			return true
		})
	}
}

func Flatten[T any](b Batch[[]T]) Batch[T] {
	return func(next func(T) bool) {
		b(func(eles []T) bool {
			for _, ele := range eles {
				if !next(ele) {
					return false
				}
			}
			return true
		})
	}
}

func Map[T, U any](b Batch[T], fn func(T) U) Batch[U] {
	return func(next func(U) bool) {
		b(func(in T) bool {
			return next(fn(in))
		})
	}
}

func Prepend[T any](b Batch[T], ele T) Batch[T] {
	return func(next func(T) bool) {
		if next(ele) {
			b(func(in T) bool {
				return next(in)
			})
		}
	}
}

func Take[T any](b Batch[T], num int) Batch[T] {
	return func(next func(T) bool) {
		b(func(in T) bool {
			if num > 0 {
				num--
				return next(in)
			}
			return false
		})
	}
}

func TakeWhile[T any](b Batch[T], fn func(T) bool) Batch[T] {
	return func(next func(T) bool) {
		done := false
		b(func(in T) bool {
			if !done && fn(in) {
				done = true
			}
			if !done {
				return next(in)
			}
			return false
		})
	}
}
