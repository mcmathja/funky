// maps provides generic convenience functions for working with slices.
package maps

import (
	"errors"

	"github.com/mcmathja/funky/constraints"
	"github.com/mcmathja/funky/pairs"
)

/* Constructors */

func FromBatch[K comparable, V any](g func(func(pairs.Pair[K, V]))) map[K]V {
	result := make(map[K]V)
	g(func(pair pairs.Pair[K, V]) {
		result[pair.Left] = pair.Right
	})
	return result
}

// FromChannel creates a new map containing all the distinct values received on ch.
// It only returns its results once the channel closes.
func FromChan[K comparable, V any](ch <-chan pairs.Pair[K, V]) map[K]V {
	result := make(map[K]V)
	for kv := range ch {
		result[kv.Left] = kv.Right
	}

	return result
}

// FromSet creates a new map containing all the key value pairs in s.
// If the same key is repeated twice, a value is chosen arbitrarily.
func FromSet[K comparable, V comparable](s map[pairs.Pair[K, V]]struct{}) map[K]V {
	result := make(map[K]V)
	for kv := range s {
		result[kv.Left] = kv.Right
	}

	return result
}

// FromSlice creates a new map containing all the key value pairs in s.
// If the same key is repeated twice, the last value wins.
func FromSlice[K comparable, V any](s []pairs.Pair[K, V]) map[K]V {
	result := make(map[K]V)
	for _, kv := range s {
		result[kv.Left] = kv.Right
	}

	return result
}

func New[K comparable, V any](kvs ...pairs.Pair[K, V]) map[K]V {
	result := make(map[K]V, len(kvs))
	for _, kv := range kvs {
		result[kv.Left] = kv.Right
	}

	return result
}

/* Operations */

func Add[K comparable, V any](m map[K]V, k K, v V) map[K]V {
	result := make(map[K]V, len(m)+1)
	for key, value := range m {
		result[key] = value
	}
	result[k] = v

	return result
}

func All[K comparable, V any](m map[K]V, fn func(K, V) bool) bool {
	for k, v := range m {
		if !fn(k, v) {
			return false
		}
	}

	return true
}

func Any[K comparable, V any](m map[K]V, fn func(K, V) bool) bool {
	for k, v := range m {
		if fn(k, v) {
			return true
		}
	}

	return false
}

func ContainsKey[K comparable, V any](m map[K]V, k K) bool {
	_, ok := m[k]
	return ok
}

func ContainsValue[K comparable, V comparable](m map[K]V, v V) bool {
	for _, e := range m {
		if e == v {
			return true
		}
	}

	return false
}

func Count[K comparable, V comparable](m map[K]V, v V, fn func(K, V) bool) int {
	cnt := 0

	for k, v := range m {
		if fn(k, v) {
			cnt++
		}
	}

	return cnt
}

func Drop[K comparable, V any](m map[K]V, num int) map[K]V {
	result := make(map[K]V)

	for k, v := range m {
		if num == 0 {
			result[k] = v
		} else if num > 0 {
			num--
		}
	}

	return result
}

func DropWhile[K comparable, V any](m map[K]V, fn func(K, V) bool) map[K]V {
	result := make(map[K]V)
	done := false

	for k, v := range m {
		if !done && !fn(k, v) {
			done = true
		}
		if done {
			result[k] = v
		}
	}

	return result
}

func Empty[K comparable, V any](m map[K]V) bool {
	return len(m) == 0
}

func Equals[K comparable, V comparable](a, b map[K]V) bool {
	if len(a) != len(b) {
		return false
	}

	for k, v := range a {
		if ele, ok := b[k]; !ok || v != ele {
			return false
		}
	}

	return true
}

func Filter[K comparable, V any](m map[K]V, fn func(K, V) bool) map[K]V {
	result := make(map[K]V, len(m))
	ForEach(m, func(key K, value V) {
		if fn(key, value) {
			result[key] = value
		}
	})

	return result
}

func FlatMap[K, T comparable, V, U any](m map[K]V, fn func(K, V) map[T]U) map[T]U {
	result := make(map[T]U)
	for k, v := range m {
		for t, u := range fn(k, v) {
			result[t] = u
		}
	}

	return result
}

func Flatten[K comparable, V any](mm []map[K]V) map[K]V {
	result := make(map[K]V)
	for _, m := range mm {
		for k, v := range m {
			result[k] = v
		}
	}

	return result
}

func ForEach[K comparable, V any](m map[K]V, fn func(key K, value V)) {
	for k, v := range m {
		fn(k, v)
	}
}

func Keys[K comparable, V any](m map[K]V) []K {
	result := make([]K, len(m))
	for k := range m {
		result = append(result, k)
	}

	return result
}

func Map[K1, K2 comparable, V1, V2 any](m map[K1]V1, fn func(K1, V1) (K2, V2)) map[K2]V2 {
	result := make(map[K2]V2, len(m))
	ForEach(m, func(k1 K1, v1 V1) {
		k2, v2 := fn(k1, v1)
		result[k2] = v2
	})

	return result
}

func MaxKey[K constraints.Ordered, V any](m map[K]V) (K, error) {
	if len(m) <= 0 {
		var ele K
		return ele, errors.New("no such element")
	}

	var best K
	var found bool
	for k := range m {
		if !found {
			best = k
			found = true
		} else if k > best {
			best = k
		}
	}

	return best, nil
}

func MaxValue[K comparable, V constraints.Ordered](m map[K]V) (V, error) {
	if len(m) <= 0 {
		var ele V
		return ele, errors.New("no such element")
	}

	var best V
	var found bool
	for _, v := range m {
		if !found {
			best = v
			found = true
		} else if v > best {
			best = v
		}
	}

	return best, nil
}

func MinKey[K constraints.Ordered, V any](m map[K]V) (K, error) {
	if len(m) <= 0 {
		var ele K
		return ele, errors.New("no such element")
	}

	var best K
	var found bool
	for k := range m {
		if !found {
			best = k
			found = true
		} else if k < best {
			best = k
		}
	}

	return best, nil
}

func MinValue[K comparable, V constraints.Ordered](m map[K]V) (V, error) {
	if len(m) <= 0 {
		var ele V
		return ele, errors.New("no such element")
	}

	var best V
	var found bool
	for _, v := range m {
		if !found {
			best = v
			found = true
		} else if v < best {
			best = v
		}
	}

	return best, nil
}

func Partition[K comparable, V any](m map[K]V, fn func(K, V) bool) (map[K]V, map[K]V) {
	a := make(map[K]V)
	b := make(map[K]V)

	for k, v := range m {
		if fn(k, v) {
			a[k] = v
		} else {
			b[k] = v
		}
	}

	return a, b
}

func ProductKeys[K constraints.Numeric, V any](m map[K]V) K {
	var product K
	for k := range m {
		product *= k
	}

	return product
}

func ProductValues[K comparable, V constraints.Numeric](m map[K]V) V {
	var product V
	for _, v := range m {
		product *= v
	}

	return product
}

func Reduce[K comparable, V any, U any](m map[K]V, initial U, fn func(U, K, V) U) U {
	acc := &initial
	for k, v := range m {
		*acc = fn(*acc, k, v)
	}

	return *acc
}

func Remove[K comparable, V any](m map[K]V, k K) map[K]V {
	result := make(map[K]V, len(m))
	for key, value := range m {
		result[key] = value
	}
	delete(result, k)

	return result
}

func Size[K comparable, V any](m map[K]V) int {
	return len(m)
}

func SumKeys[K constraints.Numeric, V any](m map[K]V) K {
	var sum K
	for k := range m {
		sum += k
	}
	return sum
}

func SumValues[K comparable, V constraints.Numeric](m map[K]V) V {
	var sum V
	for _, v := range m {
		sum += v
	}
	return sum
}

func Take[K comparable, V any](m map[K]V, num int) map[K]V {
	if num <= 0 {
		return map[K]V{}
	}

	result := make(map[K]V, num)
	for k, v := range m {
		num--
		if num < 0 {
			break
		}

		result[k] = v
	}

	return result
}

func TakeWhile[K comparable, V any](m map[K]V, fn func(K, V) bool) map[K]V {
	result := make(map[K]V)

	for k, v := range m {
		if !fn(k, v) {
			break
		}
		result[k] = v
	}

	return result
}

func Values[K comparable, V any](m map[K]V) []V {
	result := make([]V, len(m))
	for _, v := range m {
		result = append(result, v)
	}

	return result
}
