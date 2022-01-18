// sets provides generic convenience functions for working with sets.
package sets

import (
	"errors"

	"github.com/mcmathja/funky/constraints"
	"github.com/mcmathja/funky/pairs"
)

/* Constructors */

func FromBatch[T comparable](b func(func(T))) map[T]struct{} {
	result := make(map[T]struct{})
	b(func(ele T) {
		result[ele] = struct{}{}
	})
	return result
}

// FromChannel creates a new set containing all the distinct values received on ch.
// It only returns its results once the channel closes.
func FromChannel[T comparable](ch <-chan T) map[T]struct{} {
	result := make(map[T]struct{})
	for ele := range ch {
		result[ele] = struct{}{}
	}

	return result
}

// FromMap creates a new set containing all the key value pairs in m.
// Only maps with comparable values can be converted into sets.
func FromMap[K comparable, V comparable](m map[K]V) map[pairs.Pair[K, V]]struct{} {
	result := make(map[pairs.Pair[K, V]]struct{})
	for key, value := range m {
		result[pairs.New(key, value)] = struct{}{}
	}

	return result
}

// FromSlice creates a new set containing all the distinct values in s.
func FromSlice[T comparable](s []T) map[T]struct{} {
	result := make(map[T]struct{})
	for _, ele := range s {
		result[ele] = struct{}{}
	}

	return result
}

// New creates a new set from a sequence of elements eles.
func New[T comparable](eles ...T) map[T]struct{} {
	set := make(map[T]struct{}, len(eles))
	for _, ele := range eles {
		set[ele] = struct{}{}
	}

	return set
}

/* Operations */

// Add returns a new set with the provided elements added.
func Add[T comparable](s map[T]struct{}, eles ...T) map[T]struct{} {
	result := make(map[T]struct{}, len(s)+len(eles))
	for ele := range s {
		result[ele] = struct{}{}
	}
	for _, ele := range eles {
		result[ele] = struct{}{}
	}

	return result
}

func All[Elem comparable](set map[Elem]struct{}, fn func(Elem) bool) bool {
	for ele := range set {
		if !fn(ele) {
			return false
		}
	}

	return true
}

func Any[Elem comparable](set map[Elem]struct{}, fn func(Elem) bool) bool {
	for ele := range set {
		if fn(ele) {
			return true
		}
	}

	return false
}

func Contains[Elem comparable](set map[Elem]struct{}, ele Elem) bool {
	_, ok := set[ele]
	return ok
}

func Count[Elem comparable](set map[Elem]struct{}, fn func(Elem) bool) int {
	cnt := 0

	ForEach(set, func(ele Elem) {
		if fn(ele) {
			cnt++
		}
	})

	return cnt
}

// Difference returns the elements in s1 that are not in s2.
func Difference[T comparable](s1, s2 map[T]struct{}) map[T]struct{} {
	result := make(map[T]struct{})
	for ele := range s1 {
		if _, ok := s2[ele]; !ok {
			result[ele] = struct{}{}
		}
	}

	return result
}

// Drop returns a new set with num elements removed from the original set.
// There is no guaranteed order to the removed set.
func Drop[Elem comparable](set map[Elem]struct{}, num int) map[Elem]struct{} {
	return Take(set, len(set)-num)
}

// DropWhile returns a new set by removing values from set until fn returns false.
// This is no guaranteed order to the returned set.
func DropWhile[Elem comparable](set map[Elem]struct{}, fn func(Elem) bool) map[Elem]struct{} {
	result := make(map[Elem]struct{})
	done := false

	for ele := range set {
		if !done && !fn(ele) {
			done = true
		}
		if done {
			result[ele] = struct{}{}
		}
	}

	return result
}

// Empty checks whether or not the set contains any items.
func Empty[T comparable](set map[T]struct{}) bool {
	return len(set) == 0
}

func Equals[Elem comparable](a, b map[Elem]struct{}) bool {
	if len(a) != len(b) {
		return false
	}

	for ele := range a {
		if _, ok := b[ele]; !ok {
			return false
		}
	}

	return true
}

// Filter applies the predicate fn to each element of s
// in turn, returning a new set containing only
// the elements passing the predicate.
func Filter[T comparable](s map[T]struct{}, fn func(T) bool) map[T]struct{} {
	result := map[T]struct{}{}
	for ele := range s {
		if fn(ele) {
			result[ele] = struct{}{}
		}
	}

	return result
}

func FlatMap[T, U comparable](s map[T]struct{}, fn func(T) map[U]struct{}) map[U]struct{} {
	result := make(map[U]struct{})
	for ele := range s {
		for nestedEle := range fn(ele) {
			result[nestedEle] = struct{}{}
		}
	}

	return result
}

func Flatten[T comparable](ss []map[T]struct{}) map[T]struct{} {
	result := make(map[T]struct{})
	for _, s := range ss {
		for ele := range s {
			result[ele] = struct{}{}
		}
	}

	return result
}

// ForEach performs fn on each element in s.
// The order of operation is not guaranteed.
func ForEach[T comparable](s map[T]struct{}, fn func(T)) {
	for ele := range s {
		fn(ele)
	}
}

// GroupBy groups elements by the result of a function call.
func GroupBy[T, U comparable](s map[T]struct{}, fn func(T) U) map[U][]T {
	result := make(map[U][]T)
	for ele := range s {
		grouping := fn(ele)
		result[grouping] = append(result[grouping], ele)
	}

	return result
}

// Intersect returns the intersection of the passed in sets ss.
// If no sets are provided, it returns the empty set.
func Intersect[T comparable](ss ...map[T]struct{}) map[T]struct{} {
	if len(ss) == 0 {
		return make(map[T]struct{})
	}

	minIdx := 0
	minLen := len(ss[0])
	for idx, s := range ss {
		if len(s) < minLen {
			minIdx = idx
			minLen = len(s)
		}
	}

	result := make(map[T]struct{}, minLen)
Outer:
	for ele := range ss[minIdx] {
		for idx, s := range ss {
			if idx == minIdx {
				continue
			}

			if _, exists := s[ele]; !exists {
				continue Outer
			}
		}
		result[ele] = struct{}{}
	}

	return result
}

// Map creates a new set where every element in s
// has been mapped to a new element using fn.
func Map[T, U comparable](s map[T]struct{}, fn func(T) U) map[U]struct{} {
	result := make(map[U]struct{}, len(s))
	for ele := range s {
		result[fn(ele)] = struct{}{}
	}

	return result
}

// Max returns the highest valued element in s,
// or an error if it contains no values.
// s must consist of primitives having a total order.
func Max[T constraints.Ordered](s map[T]struct{}) (T, error) {
	if len(s) <= 0 {
		var ele T
		return ele, errors.New("no such element")
	}

	var best T
	var found bool
	for ele := range s {
		if !found {
			best = ele
			found = true
		} else if ele > best {
			best = ele
		}
	}

	return best, nil
}

// Min returns the lowest valued element in s,
// or an error if it contains no values.
// s must consist of primitives having a total order.
func Min[T constraints.Ordered](s map[T]struct{}) (T, error) {
	if len(s) <= 0 {
		var ele T
		return ele, errors.New("no such element")
	}

	var best T
	var found bool
	for ele := range s {
		if !found {
			best = ele
			found = true
		} else if ele < best {
			best = ele
		}
	}

	return best, nil
}

func Partition[Elem comparable](set map[Elem]struct{}, fn func(Elem) bool) (map[Elem]struct{}, map[Elem]struct{}) {
	a := map[Elem]struct{}{}
	b := map[Elem]struct{}{}

	ForEach(set, func(ele Elem) {
		if fn(ele) {
			a[ele] = struct{}{}
		} else {
			b[ele] = struct{}{}
		}
	})

	return a, b
}

// Product returns the product of the elements in s.
// s must consist of elements of a numeric type
// with a defined multiplication operation.
func Product[T constraints.Numeric](s map[T]struct{}) T {
	var product T
	for ele := range s {
		product *= ele
	}

	return product
}

// Reduce applies fn to each element of s in turn
// along with the value of an accumulator.
// The accumulator is initialized with init.
func Reduce[T comparable, U any](set map[T]struct{}, initial U, fn func(U, T) U) U {
	acc := &initial
	for ele := range set {
		*acc = fn(*acc, ele)
	}

	return *acc
}

// Remove returns a new set with the provided elements removed.
func Remove[T comparable](s map[T]struct{}, eles ...T) map[T]struct{} {
	result := make(map[T]struct{}, len(s))
	for ele := range s {
		result[ele] = struct{}{}
	}
	for _, ele := range eles {
		delete(result, ele)
	}

	return result
}

func Size[Elem comparable](set map[Elem]struct{}) int {
	return len(set)
}

// Sum returns the sum of the elements in s.
// s must consist of elements of a numeric type
// with a defined addition operation.
func Sum[T constraints.Numeric](s map[T]struct{}) T {
	var sum T
	for ele := range s {
		sum += ele
	}
	return sum
}

// Take returns a new set containing num elements from the original set.
// There is no guaranteed order to the selected set.
func Take[Elem comparable](set map[Elem]struct{}, num int) map[Elem]struct{} {
	if num <= 0 {
		return map[Elem]struct{}{}
	}

	result := make(map[Elem]struct{}, num)
	for ele := range set {
		num--
		if num < 0 {
			break
		}

		result[ele] = struct{}{}
	}

	return result
}

// TakeWhile returns a new set by selecting values from set until fn returns false.
// This is no guaranteed order to the selecting set.
func TakeWhile[Elem comparable](set map[Elem]struct{}, fn func(Elem) bool) map[Elem]struct{} {
	result := make(map[Elem]struct{})

	for ele := range set {
		if !fn(ele) {
			break
		}
		result[ele] = struct{}{}
	}

	return result
}

// Union returns the union of the passed in sets ss.
// If no sets are provided, it returns the empty set.
func Union[Elem comparable](ss ...map[Elem]struct{}) map[Elem]struct{} {
	result := make(map[Elem]struct{})
	for _, s := range ss {
		for ele := range s {
			result[ele] = struct{}{}
		}
	}

	return result
}
