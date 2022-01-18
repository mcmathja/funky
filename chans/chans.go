// chans provides generic convenience functions for working with channels.
package chans

import (
	"sync"

	"github.com/mcmathja/funky/pairs"
)

/* Constructors */

func FromBatch[T any](b func(func(T))) <-chan T {
	result := make(chan T)
	go func() {
		defer close(result)
		b(func(ele T) {
			result <- ele
		})
	}()

	return result
}

func FromFunc[Elem any](fn func() Elem) <-chan Elem {
	result := make(chan Elem)
	go func() {
		defer close(result)
		result <- fn()
	}()

	return result
}

func FromMap[K comparable, V any](m map[K]V) <-chan pairs.Pair[K, V] {
	result := make(chan pairs.Pair[K, V])

	go func() {
		defer close(result)
		for key, value := range m {
			result <- pairs.New(key, value)
		}
	}()

	return result
}

func FromSet[T comparable](m map[T]struct{}) <-chan T {
	result := make(chan T)

	go func() {
		defer close(result)
		for ele := range m {
			result <- ele
		}
	}()

	return result
}

func FromSlice[T comparable](m []T) <-chan T {
	result := make(chan T)

	go func() {
		defer close(result)
		for _, ele := range m {
			result <- ele
		}
	}()

	return result
}

// New creates a new channel from a sequence of elements.
func New[Elem any](eles ...Elem) <-chan Elem {
	result := make(chan Elem)
	go func() {
		defer close(result)
		for _, ele := range eles {
			result <- ele
		}
	}()

	return result
}

/* Operations */

func All[Elem any](ch <-chan Elem, fn func(Elem) bool) bool {
	for ele := range ch {
		if !fn(ele) {
			return false
		}
	}

	return true
}

func Any[Elem any](ch <-chan Elem, fn func(Elem) bool) bool {
	for ele := range ch {
		if fn(ele) {
			return true
		}
	}

	return false
}

func Append[Elem any](ch <-chan Elem, ele Elem) <-chan Elem {
	result := make(chan Elem)
	go func() {
		defer close(result)

		for e := range ch {
			result <- e
		}
		result <- ele
	}()

	return result
}

func Broadcast[Elem any](ch <-chan Elem, cnt int) []<-chan Elem {
	if cnt <= 0 {
		return []<-chan Elem{}
	}

	rwResults := make([]chan Elem, cnt)
	roResults := make([]<-chan Elem, cnt)
	for idx := 0; idx < cnt; idx++ {
		result := make(chan Elem)
		rwResults[idx] = result
		roResults[idx] = result
	}

	go func() {
		for _, result := range rwResults {
			result := result
			defer close(result)
		}

		for ele := range ch {
			for _, result := range rwResults {
				result <- ele
			}
		}
	}()

	return roResults
}

func Buffer[Elem any](ch <-chan Elem, size int) <-chan []Elem {
	if size <= 0 {
		return Map(ch, func(ele Elem) []Elem {
			return []Elem{}
		})
	}

	result := make(chan []Elem)
	go func() {
		defer close(result)

		buffer := []Elem{}
		for ele := range ch {
			buffer = append(buffer, ele)
			if len(buffer) >= size {
				result <- buffer
				buffer = []Elem{}
			}
		}

		if len(buffer) > 0 {
			result <- buffer
		}
	}()

	return result
}

func Concat[Elem any](chs ...<-chan Elem) <-chan Elem {
	result := make(chan Elem)
	go func() {
		defer close(result)
		for _, ch := range chs {
			for ele := range ch {
				result <- ele
			}
		}
	}()

	return result
}

func Contains[Elem comparable](ch <-chan Elem, ele Elem) bool {
	for e := range ch {
		if e == ele {
			return true
		}
	}

	return false
}

func ContainsSequence[Elem comparable](ch <-chan Elem, subseq []Elem) bool {
	if len(subseq) <= 0 {
		return true
	}

	window := []Elem{}
Outer:
	for ele := range ch {
		window = append(window, ele)
		if len(window) > len(subseq) {
			window = window[1:]
		}

		if len(window) < len(subseq) {
			continue
		}

		for idx, ele := range subseq {
			if window[idx] != ele {
				continue Outer
			}
		}

		return true
	}

	return false
}

func Corresponds[Elem any](a, b <-chan Elem, fn func(Elem, Elem) bool) bool {
	for {
		aVal, aDone := <-a
		bVal, bDone := <-b

		if aDone != bDone {
			return false
		}

		if aDone && bDone {
			return true
		}

		if !fn(aVal, bVal) {
			return false
		}
	}
}

func Count[Elem any](ch <-chan Elem, fn func(Elem) bool) int {
	n := 0
	for ele := range ch {
		if fn(ele) {
			n++
		}
	}

	return n
}

func Distinct[Elem comparable](ch <-chan Elem) <-chan Elem {
	return DistinctBy(ch, func(ele Elem) Elem {
		return ele
	})
}

func DistinctBy[Elem any, Comp comparable](ch <-chan Elem, fn func(Elem) Comp) <-chan Elem {
	result := make(chan Elem)

	go func() {
		seen := map[Comp]struct{}{}
		for ele := range ch {
			comp := fn(ele)
			if _, ok := seen[comp]; !ok {
				result <- ele
				seen[comp] = struct{}{}
			}
		}
	}()

	return result
}

func Distribute[Elem any](ch <-chan Elem, cnt int) []<-chan Elem {
	if cnt <= 0 {
		return []<-chan Elem{}
	}

	rwResults := make([]chan Elem, cnt)
	roResults := make([]<-chan Elem, cnt)
	for idx := 0; idx < cnt; idx++ {
		result := make(chan Elem)
		rwResults[idx] = result
		roResults[idx] = result
	}

	for _, rwResult := range rwResults {
		go func(rwResult chan Elem) {
			defer close(rwResult)
			for ele := range ch {
				rwResult <- ele
			}
		}(rwResult)
	}

	return roResults
}

func Drop[Elem any](ch <-chan Elem, num int) <-chan Elem {
	result := make(chan Elem)
	go func() {
		defer close(result)
		for ele := range ch {
			if num > 0 {
				num--
				continue
			}

			result <- ele
		}
	}()

	return result
}

func DropWhile[Elem any](ch <-chan Elem, fn func(Elem) bool) <-chan Elem {
	result := make(chan Elem)
	go func() {
		defer close(result)
		done := false
		for ele := range ch {
			if !done && fn(ele) {
				done = true
			}
			if done {
				result <- ele
			}
		}
	}()

	return result
}

func EndsWith[Elem comparable](ch <-chan Elem, ele Elem) bool {
	matched := false
	for e := range ch {
		matched = e == ele
	}

	return matched
}

func EndsWithSequence[Elem comparable](ch <-chan Elem, subseq []Elem) bool {
	if len(subseq) <= 0 {
		return true
	}

	window := []Elem{}
	for ele := range ch {
		window = append(window, ele)
		if len(window) > len(subseq) {
			window = window[1:]
		}
	}

	if len(subseq) != len(window) {
		return false
	}

	for idx, ele := range subseq {
		if window[idx] != ele {
			return false
		}
	}

	return true
}

func Equals[Elem comparable](a, b <-chan Elem) bool {
	return Corresponds(a, b, func(i, j Elem) bool {
		return i == j
	})
}

func Filter[Elem any](ch <-chan Elem, fn func(Elem) bool) <-chan Elem {
	result := make(chan Elem)
	go func() {
		defer close(result)
		for ele := range ch {
			if fn(ele) {
				result <- ele
			}
		}
	}()

	return result
}

func First[Elem any](ch <-chan Elem) Elem {
	return <-ch
}

func FirstWhere[Elem any](ch <-chan Elem, fn func(Elem) bool) <-chan Elem {
	return NthWhere(ch, 1, fn)
}

func FlatMap[From, To any](ch <-chan From, fn func(From) <-chan To) <-chan To {
	return Flatten(Map(ch, fn))
}

func Flatten[Elem any](ch <-chan <-chan Elem) <-chan Elem {
	result := make(chan Elem)
	go func() {
		defer close(result)
		for subch := range ch {
			for ele := range subch {
				result <- ele
			}
		}
	}()

	return result
}

func Last[Elem any](ch <-chan Elem, fn func(Elem) bool) <-chan Elem {
	result := make(chan Elem)
	go func() {
		defer close(result)

		var found bool
		var target Elem
		for ele := range ch {
			found = true
			target = ele
		}

		if found {
			result <- target
		}
	}()

	return result
}

func LastWhere[Elem any](ch <-chan Elem, fn func(Elem) bool) <-chan Elem {
	result := make(chan Elem)
	go func() {
		defer close(result)

		var found bool
		var target Elem
		for ele := range ch {
			if fn(ele) {
				found = true
				target = ele
			}
		}

		if found {
			result <- target
		}
	}()

	return result
}

func Map[From, To any](ch <-chan From, fn func(From) To) <-chan To {
	result := make(chan To)
	go func() {
		defer close(result)
		for ele := range ch {
			result <- fn(ele)
		}
	}()

	return result
}

func Merge[Elem any](chs ...<-chan Elem) <-chan Elem {
	result := make(chan Elem)

	var wg sync.WaitGroup
	for _, ch := range chs {
		wg.Add(1)
		go func(ch <-chan Elem) {
			defer wg.Done()
			for ele := range ch {
				result <- ele
			}
		}(ch)
	}
	go func() {
		defer close(result)
		wg.Wait()
	}()

	return result
}

func NthWhere[Elem any](ch <-chan Elem, n int, fn func(Elem) bool) <-chan Elem {
	result := make(chan Elem)
	go func() {
		defer close(result)
		if n <= 0 {
			return
		}

		for ele := range ch {
			if fn(ele) {
				n--
				if n == 0 {
					result <- ele
					break
				}
			}
		}
	}()

	return result
}

func Partition[Elem any](ch <-chan Elem, fn func(Elem) bool) (<-chan Elem, <-chan Elem) {
	left := make(chan Elem)
	right := make(chan Elem)
	go func() {
		defer close(left)
		defer close(right)
		for ele := range ch {
			if fn(ele) {
				left <- ele
			} else {
				right <- ele
			}
		}
	}()

	return left, right
}

func Prepend[Elem any](ch <-chan Elem, ele Elem) <-chan Elem {
	result := make(chan Elem)
	go func() {
		defer close(result)

		result <- ele
		for e := range ch {
			result <- e
		}
	}()

	return result
}

func Reduce[Elem any, Acc any](ch <-chan Elem, initial Acc, fn func(Acc, Elem) Acc) <-chan Acc {
	result := make(chan Acc)
	go func() {
		defer close(result)
		acc := &initial
		for ele := range ch {
			*acc = fn(*acc, ele)
			result <- *acc
		}
	}()

	return result
}

func Size[Elem any](ch <-chan Elem, fn func(Elem) bool) int {
	return Count(ch, func(ele Elem) bool {
		return true
	})
}

func SplitAt[Elem any](ch <-chan Elem, n int) (<-chan Elem, <-chan Elem) {
	if n < 0 {
		n = 0
	}

	return Partition(ch, func(ele Elem) bool {
		n--
		return n < 0
	})
}

func StartsWith[Elem comparable](ch <-chan Elem, ele Elem) bool {
	return ele == <-ch
}

func StartsWithSequence[Elem comparable](ch <-chan Elem, subseq []Elem) bool {
	if len(subseq) <= 0 {
		return true
	}

	idx := 0
	for ele := range ch {
		if subseq[idx] != ele {
			return false
		}
		idx++
		if idx >= len(subseq) {
			break
		}
	}

	return true
}

func Take[Elem any](ch <-chan Elem, num int) <-chan Elem {
	result := make(chan Elem)
	go func() {
		defer close(result)
		for ele := range ch {
			num--
			if num < 0 {
				break
			}

			result <- ele
		}
	}()

	return result
}

func TakeWhile[Elem any](ch <-chan Elem, fn func(Elem) bool) <-chan Elem {
	result := make(chan Elem)
	go func() {
		defer close(result)
		for ele := range ch {
			if !fn(ele) {
				break
			}

			result <- ele
		}
	}()

	return result
}

func Window[Elem any](ch <-chan Elem, size int) <-chan []Elem {
	if size <= 0 {
		return Map(ch, func(ele Elem) []Elem {
			return []Elem{}
		})
	}

	result := make(chan []Elem)
	go func() {
		defer close(result)

		window := []Elem{}
		for ele := range ch {
			window = append(window, ele)
			if len(window) > size {
				window = window[1:]
			}
			result <- window
		}
	}()

	return result
}
