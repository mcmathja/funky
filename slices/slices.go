// slices provides generic convenience functions for working with slices.
package slices

import (
	"errors"
	"sort"

	"github.com/mcmathja/funky/constraints"
	"github.com/mcmathja/funky/pairs"
)

// All returns true if all of the elements in s
// satisfy the predicate fn. Otherwise, it returns false.
func All[T any](s []T, fn func(T) bool) bool {
	for _, ele := range s {
		if !fn(ele) {
			return false
		}
	}

	return true
}

// Any returns true if any of the elements in s
// satisfy the predicate fn. Otherwise, it returns false.
func Any[T any](s []T, fn func(T) bool) bool {
	for _, ele := range s {
		if fn(ele) {
			return true
		}
	}

	return false
}

// Append returns a copy of s with ele added at the end.
func Append[T any](s []T, eles ...T) []T {
	result := make([]T, 0, len(s)+len(eles))
	result = append(result, s...)
	result = append(result, eles...)

	return result
}

// AtLeast determines whether the predicate fn
// passes for at least n elements in s.
func AtLeast[T any](s []T, n int, fn func(T) bool) bool {
	cnt := 0
	for _, ele := range s {
		if fn(ele) {
			cnt++
			if cnt >= n {
				return true
			}
		}
	}

	return cnt >= n
}

// AtMost determines whether the predicate fn
// passes for at most n elements in s.
func AtMost[T any](s []T, n int, fn func(T) bool) bool {
	cnt := 0
	for _, ele := range s {
		if fn(ele) {
			cnt++
			if cnt > n {
				return false
			}
		}
	}

	return cnt <= n
}

// Cartesian generates the cartesian product
// of all elements from s and ss.
func Cartesian[T, U any](s []T, ss []U) []pairs.Pair[T, U] {
	result := make([]pairs.Pair[T, U], 0, len(s)*len(ss))
	for _, left := range s {
		for _, right := range ss {
			result = append(result, pairs.New(left, right))
		}
	}

	return result
}

// ConsistsOf checks if s is made up of only elements
// that are also present in eles, without regard
// for arrangement or repetition.
func ConsistsOf[T comparable](s []T, eles ...T) bool {
	seen := make(map[T]struct{})
	for _, ele := range eles {
		seen[ele] = struct{}{}
	}
	for _, ele := range s {
		if _, ok := seen[ele]; !ok {
			return false
		}
	}

	return true
}

// Comprises checks if s is made up of exactly the
// elements in eles, including any repetitions,
// but ordered in any arrangement.
func Comprises[T comparable](s []T, eles ...T) bool {
	if len(eles) != len(s) {
		return false
	}

	cnts := make(map[T]int)
	for _, ele := range eles {
		cnts[ele]++
	}
	for _, ele := range s {
		if _, ok := cnts[ele]; !ok {
			return false
		} else {
			cnts[ele]--
			if cnts[ele] == 0 {
				delete(cnts, ele)
			}
		}
	}

	return len(cnts) == 0
}

// Contains checks if s contains ele.
func Contains[T comparable](s []T, ele T) bool {
	for _, e := range s {
		if e == ele {
			return true
		}
	}

	return false
}

// ContainsAll checks if s contains every element in eles.
func ContainsAll[T comparable](s []T, eles ...T) bool {
	if len(eles) == 0 {
		return true
	}

	if len(eles) > len(s) {
		return false
	}

	cnts := make(map[T]int)
	for _, ele := range eles {
		cnts[ele]++
	}
	for _, ele := range s {
		if _, ok := cnts[ele]; ok {
			cnts[ele]--
			if cnts[ele] == 0 {
				delete(cnts, ele)
			}
		}

		if len(cnts) == 0 {
			return true
		}
	}

	return false
}

// ContainsAny checks if s contains any element in eles.
func ContainsAny[T comparable](s []T, eles ...T) bool {
	if len(eles) == 0 {
		return false
	}

	seen := make(map[T]struct{})
	for _, ele := range eles {
		seen[ele] = struct{}{}
	}
	for _, ele := range s {
		if _, ok := seen[ele]; ok {
			return true
		}
	}

	return false
}

// containsSequenceSearchAlgorithm specifies the substring
// search algorithm to use when evaluating ContainsSequence.
type containsSequenceSearchAlgorithm string

const (
	containsSequenceApostolicoCrochemore containsSequenceSearchAlgorithm = "ApostolicoCrochemore"
	containsSequenceBruteForce           containsSequenceSearchAlgorithm = "BruteForce"
)

// containsSequenceArgs represent optional arguments to ContainsSequence.
type containsSequenceArgs struct {
	// algorithm indicates
	searchAlgorithm containsSequenceSearchAlgorithm
}

// ContainsSequenceOpt represent optional arguments to ContainsSequence.
type ContainsSequenceOpt func(*containsSequenceArgs)

// ContainsSequenceApostolicoCrochemore is a ContainsSequenceOpt
// that specifies we should use the Apostolico Crochemore algorithm
// to perform substring search.
func ContainsSequenceApostolicoCrochemore(args *containsSequenceArgs) {
	args.searchAlgorithm = containsSequenceApostolicoCrochemore
}

// ContainsSequenceApostolicoCrochemore is a ContainsSequenceOpt
// that specifies we should use the Brute Force algorithm
// to perform substring search.
func ContainsSequenceBruteForce(args *containsSequenceArgs) {
	args.searchAlgorithm = containsSequenceBruteForce
}

// ContainsSequence checks if s contains the provided seq.
func ContainsSequence[T comparable](s, seq []T, opts ...ContainsSequenceOpt) bool {
	if len(s) < len(seq) {
		// The slice can't possibly contain the sequence.
		return false
	}

	if len(seq) == 0 {
		// All slices contain the empty sequence.
		return true
	}

	if len(seq) == 1 {
		// Containing a single-valued sequence is equivalent
		// to containing the single element in that sequence.
		return Contains(s, seq[0])
	}

	args := containsSequenceArgs{}
	for _, opt := range opts {
		opt(&args)
	}

	switch args.searchAlgorithm {
	case containsSequenceApostolicoCrochemore:
		return apostolicoCrochemoreSearch(s, seq)
	case containsSequenceBruteForce:
		fallthrough
	default:
		return bruteForceSearch(s, seq)
	}
}

// Corresponds compares each element in s1 against its
// corresponding element in s2 using a predicate,
// returning true if the predicate returns true for every element.
// Slices of unequal lengths never correspond.
// Empty slices always correspond to each other.
func Correspond[T any](s1, s2 []T, fn func(T, T) bool) bool {
	if len(s1) != len(s2) {
		return false
	}

	for idx, ele := range s1 {
		if !fn(ele, s2[idx]) {
			return false
		}
	}

	return true
}

// Count counts the number of elements in s
// that satisfy the predicate fn.
func Count[T any](s []T, fn func(T) bool) int {
	cnt := 0

	for _, ele := range s {
		if fn(ele) {
			cnt++
		}
	}

	return cnt
}

// Distinct returns a copy of s with all duplicate elements removed.
func Distinct[T comparable](s []T) []T {
	result := make([]T, 0)
	seen := make(map[T]struct{}, 0)

	for _, ele := range s {
		if _, ok := seen[ele]; !ok {
			result = append(result, ele)
			seen[ele] = struct{}{}
		}
	}

	return result
}

// DistinctBy returns a copy of s with all duplicate elements removed,
// where duplicates are determined by the value returned by fn.
func DistinctBy[T any, Comp comparable](s []T, fn func(T) Comp) []T {
	result := make([]T, 0)
	seen := make(map[Comp]struct{}, 0)

	for _, ele := range s {
		comp := fn(ele)
		if _, ok := seen[comp]; !ok {
			result = append(result, ele)
			seen[comp] = struct{}{}
		}
	}

	return result
}

// Drop returns a new slice where the first num elements
// of s have been removed.
func Drop[T any](s []T, num int) []T {
	if num < 0 {
		num = 0
	}

	if len(s) < num {
		num = len(s)
	}

	result := make([]T, len(s)-num)
	copy(result, s[num:])

	return result
}

// DropWhile returns a new slice where the longest prefix
// of elements in s satisfying the predicate fn have been removed.
func DropWhile[T any](slice []T, fn func(T) bool) []T {
	var limit int
	for idx, ele := range slice {
		if !fn(ele) {
			break
		}
		limit = idx + 1
	}

	result := make([]T, len(slice)-limit)
	copy(result, slice[limit:])

	return result
}

// Empty checks whether s has any elements.
func Empty[T any](s []T) bool {
	return len(s) == 0
}

// EndsWith checks whether the last element of s is ele.
// If s is empty, it always returns false.
func EndsWith[T comparable](s []T, ele T) bool {
	if len(s) == 0 {
		return false
	}

	return s[len(s)-1] == ele
}

// EndsWithSequence checks if s ends with seq.
func EndsWithSequence[T comparable](s, seq []T) bool {
	lenDiff := len(s) - len(seq)
	if lenDiff < 0 {
		// It's impossible for s to contain all of seq.
		return false
	}

	for idx := len(seq) - 1; idx >= 0; idx-- {
		if s[idx+lenDiff] != seq[idx] {
			return false
		}
	}

	return true
}

// Enumerate executes fn for each index and element in s in order.
func Enumerate[T any](slice []T, fn func(int, T)) {
	for idx, ele := range slice {
		fn(idx, ele)
	}
}

// Equal compares s1 and s2 for element-wise equality.
// Slices of unequal lengths are never equal.
// Empty slices are always equal to each other.
func Equal[T comparable](s1, s2 []T) bool {
	if len(s1) != len(s2) {
		return false
	}

	for idx, ele := range s1 {
		if ele != s2[idx] {
			return false
		}
	}

	return true
}

// Exactly determines whether the predicate fn
// passes for exactly n elements in s.
func Exactly[T any](s []T, n int, fn func(T) bool) bool {
	cnt := 0
	for _, ele := range s {
		if fn(ele) {
			cnt++
			if cnt > n {
				return false
			}
		}
	}

	return cnt == n
}

// Filter applies the predicate fn to each element of s
// in turn, returning a new slice containing only
// the elements passing the predicate.
func Filter[T any](s []T, fn func(T) bool) []T {
	ss := make([]T, 0)
	for _, ele := range s {
		if fn(ele) {
			ss = append(ss, ele)
		}
	}

	return ss
}

// First returns the first item in s,
// or an error if it contains no values.
func First[T any](s []T) (T, error) {
	if len(s) <= 0 {
		var ele T
		return ele, errors.New("no such element")
	}

	return s[0], nil
}

// FirstIndexOf returns the index of the first occurrence of ele in s.
// If no matching element is found, it returns -1.
func FirstIndexOf[T comparable](s []T, ele T) int {
	for idx, e := range s {
		if e == ele {
			return idx
		}
	}

	return -1
}

// FirstIndexWhere returns the index of the first element in s
// satisfying the predicate fn. If no matching element if found,
// it returns -1.
func FirstIndexWhere[T any](s []T, fn func(T) bool) int {
	for idx, ele := range s {
		if fn(ele) {
			return idx
		}
	}

	return -1
}

// FlatMap maps each element of s to a slice of elements,
// then flattens the result into a single slice.
func FlatMap[T, U any](s []T, fn func(T) []U) []U {
	result := make([]U, 0)
	for _, ele := range s {
		result = append(result, fn(ele)...)
	}

	return result
}

// Flatten flattens the nested slice s into a single-level slice
// consisting of the elements of each subslice in order.
func Flatten[T any](s [][]T) []T {
	result := make([]T, 0)
	for _, ss := range s {
		result = append(result, ss...)
	}

	return result
}

// Enumerate executes fn for each element in s in order.
func ForEach[T any](s []T, fn func(T)) {
	for _, ele := range s {
		fn(ele)
	}
}

// FromBatch creates a new slice containing all of the values produced by b.
// It only returns its results once the batch completes.
func FromBatch[T any](b func(func(T) bool)) []T {
	result := make([]T, 0)
	b(func(ele T) bool {
		result = append(result, ele)
		return true
	})
	return result
}

// FromChannel creates a new slice containing all the values received on ch.
// It only returns its results once the channel closes.
func FromChan[T any](ch <-chan T) []T {
	result := make([]T, 0)
	for ele := range ch {
		result = append(result, ele)
	}

	return result
}

// FromMap creates a new slice containing all of the keys and values
// in m as a slice of pairs.Pair. The result order is not guaranteed.
func FromMap[K comparable, V any](m map[K]V) []pairs.Pair[K, V] {
	result := make([]pairs.Pair[K, V], 0, len(m))
	for key, value := range m {
		result = append(result, pairs.New(key, value))
	}

	return result
}

// FromSet creates a new slice containing all of the elements in s.
// The result order is not guaranteed.
func FromSet[K comparable](s map[K]struct{}) []K {
	result := make([]K, 0, len(s))
	for key := range s {
		result = append(result, key)
	}

	return result
}

// GroupBy groups elements by the result of a function call.
func GroupBy[T any, U comparable](s []T, fn func(T) U) map[U][]T {
	result := make(map[U][]T)
	for _, ele := range s {
		grouping := fn(ele)
		result[grouping] = append(result[grouping], ele)
	}

	return result
}

// Last returns the last item in s,
// or an error if it contains no values.
func Last[T any](s []T) (T, error) {
	if len(s) <= 0 {
		var ele T
		return ele, errors.New("no such element")
	}

	return s[len(s)-1], nil
}

// LastIndexOf returns the index of the last occurrence of ele in s.
// If no matching element is found, it returns -1.
func LastIndexOf[T comparable](s []T, ele T) int {
	for idx := len(s) - 1; idx >= 0; idx-- {
		if ele == s[idx] {
			return idx
		}
	}

	return -1
}

// LastIndexWhere returns the index of the last element in s
// that satisfies the predicate fn. If no matching element is found,
// it returns -1.
func LastIndexWhere[T any](s []T, fn func(T) bool) int {
	for idx := len(s) - 1; idx >= 0; idx-- {
		if fn(s[idx]) {
			return idx
		}
	}

	return -1
}

// Map creates a new slice where every element in s
// has been mapped to a new element using fn.
func Map[T, U any](s []T, fn func(T) U) []U {
	ss := make([]U, 0, len(s))
	for _, ele := range s {
		ss = append(ss, fn(ele))
	}

	return ss
}

// Max returns the highest valued element in s,
// or an error if it contains no values.
// s must consist of primitives having a total order.
func Max[T constraints.Ordered](s []T) (T, error) {
	if len(s) <= 0 {
		var ele T
		return ele, errors.New("no such element")
	}

	best := s[0]
	for idx := 1; idx < len(s); idx++ {
		if s[idx] > best {
			best = s[idx]
		}
	}

	return best, nil
}

// Min returns the lowest valued element in s,
// or an error if it contains no values.
// s must consist of primitives having a total order.
func Min[T constraints.Ordered](s []T) (T, error) {
	if len(s) <= 0 {
		var ele T
		return ele, errors.New("no such element")
	}

	best := s[0]
	for idx := 1; idx < len(s); idx++ {
		if s[idx] < best {
			best = s[idx]
		}
	}

	return best, nil
}

// New creates a new slice from eles.
func New[T any](eles ...T) []T {
	return eles
}

// NthIndexWhere returns the index of the nth occurrence of ele in s
// where n=1 is the first match, n=2 is the second, and so on.
// If an nth matching element is not found, it returns -1.
func NthIndexOf[T comparable](s []T, n int, ele T) int {
	for idx, e := range s {
		if e == ele {
			n--
			if n == 0 {
				return idx
			}
		}
	}

	return -1
}

// NthIndexWhere returns the index of the nth element in s
// satisfying the predicate fn, where n=1 is the first match,
// n=2 is the second, and so on. If an nth matching element is not found,
// it returns -1.
func NthIndexWhere[T any](s []T, n int, fn func(T) bool) int {
	for idx, ele := range s {
		if fn(ele) {
			n--
			if n == 0 {
				return idx
			}
		}
	}

	return -1
}

// Partition divides elements from s into two slices based on a predicate,
// with passing elements in the first slice and failing elements in the second.
func Partition[T any](s []T, fn func(T) bool) ([]T, []T) {
	a := make([]T, 0)
	b := make([]T, 0)

	for _, ele := range s {
		if fn(ele) {
			a = append(a, ele)
		} else {
			b = append(b, ele)
		}
	}

	return a, b
}

func Permute[T any](s []T) [][]T {
	// Set up the iteration state and the current permutation
	// as the initial arrangement of elements.
	n := len(s)
	state := make([]int, n)
	curr := make([]T, n)
	copy(curr, s)

	// Store the initial permutation.
	orig := make([]T, n)
	copy(orig, curr)
	results := [][]T{orig}

	// Generate all remaining permutations using the Heap algorithm.
	i := 1
	for i < n {
		if state[i] < i {
			if i%2 == 0 {
				curr[0], curr[i] = curr[i], curr[0]
			} else {
				curr[state[i]], curr[i] = curr[i], curr[state[i]]
			}

			// Store the resulting permutation.
			permutation := make([]T, n)
			copy(permutation, curr)
			results = append(results, permutation)

			// Restore the base case for this sub-iteration.
			state[i]++
			i = 1
		} else {
			state[i] = 0
			i++
		}
	}

	return results
}

// Prepend returns a copy of s with ele added at the beginning.
func Prepend[T any](s []T, eles ...T) []T {
	result := make([]T, 0, len(s)+len(eles))
	result = append(result, eles...)
	result = append(result, s...)

	return result
}

// Product returns the product of the elements in s.
// s must consist of elements of a numeric type
// with a defined multiplication operation.
func Product[T constraints.Numeric](s []T) T {
	var product T = 1
	for _, ele := range s {
		product *= ele
	}

	return product
}

// Range produces a new slice containing the values
// between from (inclusive) and to (exclusive) by step.
// If step is zero, an empty slice is returned.
func Range[T constraints.Real](from, to, step T) []T {
	result := make([]T, 0)
	if step == 0 {
		return result
	}

	if step > 0 {
		for num := from; num < to; num += step {
			result = append(result, num)
		}
	} else {
		for num := from; num > to; num += step {
			result = append(result, num)
		}
	}

	return result
}

// Reduce applies fn to each element of s in turn
// along with the value of an accumulator.
// The accumulator is initialized with init.
func Reduce[T, U any](s []T, init U, fn func(U, T) U) U {
	acc := &init
	for _, ele := range s {
		*acc = fn(*acc, ele)
	}

	return *acc
}

// Repeat returns a slice with ele repeated num times.
func Repeat[T any](ele T, num int) []T {
	if num < 0 {
		num = 0
	}

	result := make([]T, num)
	for idx := 0; idx < num; idx++ {
		result[idx] = ele
	}

	return result
}

// Reversed returns a copy of s with its elements reversed.
func Reversed[T any](s []T) []T {
	result := make([]T, len(s))
	for idx, ele := range s {
		result[len(s)-idx-1] = ele
	}

	return result
}

// Rotate returns a slice where each element of s has been
// moved by k positions, with elements at the end of
// the slice wrapping around to the other side.
// If k is positive, elements are moved to the right,
// otherwise they are moved to the left.
func Rotate[T any](s []T, k int) []T {
	n := len(s)
	if n != 0 {
		k = k % n
	}
	if k < 0 {
		k = n + k
	}

	result := make([]T, n)
	for idx, ele := range s {
		result[(idx+k)%n] = ele
	}
	return result
}

// Size returns the size of s.
func Size[T any](s []T) int {
	return len(s)
}

// sortArgs represent optional arguments to Sort.
type sortArgs struct {
	// stable indicates whether a stable sort should be performed.
	stable bool
}

// sortArgs represent optional arguments to Sort.
type SortOpt func(*sortArgs)

// SortStable is a SortOpt that indicates
// a stable sort should be performed.
func SortStable(o *sortArgs) {
	o.stable = true
}

// Sort returns a new slice with the elements in s in sorted order.
func Sort[T constraints.Ordered](s []T, opts ...SortOpt) []T {
	args := sortArgs{}
	for _, opt := range opts {
		opt(&args)
	}

	result := make([]T, len(s))
	copy(result, s)

	if args.stable {
		sort.SliceStable(result, func(i, j int) bool {
			return result[i] < result[j]
		})
	} else {
		sort.Slice(result, func(i, j int) bool {
			return result[i] < result[j]
		})
	}

	return result
}

// sortByArgs represent optional arguments to Sort.
type sortByArgs struct {
	// stable indicates whether a stable sort should be performed.
	stable bool
}

// sortByArgs represent optional arguments to Sort.
type SortByOpt func(*sortByArgs)

// SortByStable is a SortByOpt that indicates
// a stable sort should be performed.
func SortByStable(o *sortByArgs) {
	o.stable = true
}

// SortBy returns a new slice with the elements in s
// sorted according to the provided less function.
func SortBy[T any](s []T, less func(a, b T) bool, opts ...SortByOpt) []T {
	args := sortByArgs{}
	for _, opt := range opts {
		opt(&args)
	}

	result := make([]T, len(s))
	copy(result, s)

	if args.stable {
		sort.SliceStable(result, func(i, j int) bool {
			return less(result[i], result[j])
		})
	} else {
		sort.Slice(result, func(i, j int) bool {
			return less(result[i], result[j])
		})
	}

	return result
}

// SplitAt splits the elements of s into two slices.
// All elements in s with an index before idx
// are returned in the first slice,
// while all elements in s with an index greater than or equal to idx
// are returned in the second slice.
func SplitAt[T any](s []T, idx int) ([]T, []T) {
	if idx < 0 {
		idx = 0
	}

	if len(s) < idx {
		idx = len(s)
	}

	before := make([]T, idx)
	after := make([]T, len(s)-idx)
	copy(before, s[:idx])
	copy(after, s[idx:])

	return before, after
}

// StartsWith checks whether the first element of s is ele.
// If s is empty, it always returns false.
func StartsWith[T comparable](s []T, ele T) bool {
	if len(s) == 0 {
		return false
	}

	return s[0] == ele
}

// StartsWithSequence checks if s starts with seq.
func StartsWithSequence[T comparable](s, seq []T) bool {
	if len(s) < len(seq) {
		// It's impossible for s to contain all of seq.
		return false
	}

	for idx, ele := range seq {
		if s[idx] != ele {
			return false
		}
	}

	return true
}

// Sum returns the sum of the elements in s.
// s must consist of elements of a numeric type
// with a defined addition operation.
func Sum[T constraints.Numeric](s []T) T {
	var sum T
	for _, ele := range s {
		sum += ele
	}
	return sum
}

// Take returns a new slice containing the first num elements of s.
func Take[T any](s []T, num int) []T {
	if num < 0 {
		num = 0
	}

	if num > len(s) {
		num = len(s)
	}

	result := make([]T, num)
	copy(result, s[:num])

	return result
}

// TakeWhile returns a new slice consisting of the longest prefix
// of elements in s satisfying the predicate fn.
func TakeWhile[T any](s []T, fn func(T) bool) []T {
	var limit int
	for idx, ele := range s {
		if !fn(ele) {
			break
		}

		limit = idx + 1
	}

	result := make([]T, limit)
	copy(result, s[:limit])

	return result
}

// Tally produces a map from each distinct element in s
// to the number of occurrences of that element in s.
func Tally[T comparable](s []T) map[T]int {
	cnts := make(map[T]int)
	for _, ele := range s {
		cnts[ele]++
	}
	return cnts
}

// TallyBy produces a map from the distinct results of fn,
// applied against each element in s,
// to the number of occurrences of that result.
func TallyBy[T any, U comparable](s []T, fn func(T) U) map[U]int {
	cnts := make(map[U]int)
	for _, ele := range s {
		cnts[fn(ele)]++
	}
	return cnts
}

// Transpose returns the transposition of s:
// given s is a matrix of shape [m][n]T,
// it returns a new matrix t of shape [n][m]T,
// where t[j][i] = s[i][j] for 0 <= i < m and 0 <= j <= n.
// If n is not consistent across subslices, it returns an error.

func Transpose[T any](s [][]T) ([][]T, error) {
	m := len(s)
	if m == 0 {
		return make([][]T, 0), nil
	}

	n := len(s[0])
	for i := 1; i < len(s); i++ {
		if len(s[i]) != n {
			return nil, errors.New("all slices in s must have the same length")
		}
	}

	result := make([][]T, n)
	for i := 0; i < n; i++ {
		result[i] = make([]T, m)
		for j := 0; j < m; j++ {
			result[i][j] = s[j][i]
		}
	}

	return result, nil
}

// Updated returns a new slice with the item at index
// replaced with the provided element.
func Updated[T any](s []T, idx int, ele T) ([]T, error) {
	if idx < 0 || idx >= len(s) {
		return nil, errors.New("index out of bounds")
	}

	ss := make([]T, len(s))
	copy(ss, s)
	ss[idx] = ele
	return ss, nil
}

// Zip matches up the elements at each index in s and ss
// and returns the result as a "zipped up" slice of pairs.
// For each pair in the resulting slice, the Left value
// is the element at the corresponding index in s, and
// the Right value is the element at the corresponding
// index in ss. If the slices have unequal lengths, the
// zero value is used to fill holes left by the shorter slice.
func Zip[T, U any](s []T, ss []U) []pairs.Pair[T, U] {
	max := len(s)
	if len(ss) > max {
		max = len(ss)
	}

	result := make([]pairs.Pair[T, U], max)
	for idx := range result {
		var left T
		var right U

		if idx < len(s) {
			left = s[idx]
		}
		if idx < len(ss) {
			right = ss[idx]
		}

		result[idx].Left = left
		result[idx].Right = right
	}

	return result
}

/* Helpers */

// bruteForceSearch performs a naive brute force search
// for the subarray seq in s.
func bruteForceSearch[T comparable](s, seq []T) bool {
Outer:
	for idx := range s {
		if len(s)-idx < len(seq) {
			break
		}

		for seqIdx := range seq {
			if s[idx+seqIdx] != seq[seqIdx] {
				continue Outer
			}
		}

		return true
	}

	return false
}

// apostolicoCrochemoreSearch implements the Apostolico-Crochemore
// substring algorithm, a variant of Knuth, Morris and Pratt.
// See http://www-igm.univ-mlv.fr/~lecroq/string/node12.html.
func apostolicoCrochemoreSearch[T comparable](s, seq []T) bool {
	// Preprocess the search sequence into a shift table.
	kmpNext := make([]int, len(seq)+1)
	x := 0
	y := -1
	kmpNext[0] = -1
	for x < len(seq) {
		for y > -1 && (x == len(seq) || seq[x] != seq[y]) {
			y = kmpNext[y]
		}
		x++
		y++
		if x != len(seq) && seq[x] == seq[y] {
			kmpNext[x] = kmpNext[y]
		} else {
			kmpNext[x] = y
		}
	}

	ell := 1
	for ell < len(seq) && seq[ell-1] == seq[ell] {
		ell++
	}
	if ell == len(seq) {
		ell = 0
	}

	// Perform the search.
	i := ell
	j := 0
	k := 0
	for j <= len(s)-len(seq) {
		for i < len(seq) && seq[i] == s[i+j] {
			i++
		}
		if i >= len(seq) {
			for k < ell && seq[k] == s[j+k] {
				k++
			}
			if k >= ell {
				return true
			}
		}
		j += (i - kmpNext[i])
		if i == ell {
			if k-1 > 0 {
				k = k - 1
			} else {
				k = 0
			}
		} else {
			if kmpNext[i] <= ell {
				if kmpNext[i] > 0 {
					k = kmpNext[i]
				} else {
					k = 0
				}
				i = ell
			} else {
				k = ell
				i = kmpNext[i]
			}
		}
	}

	return false
}
