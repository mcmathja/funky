package slices_test

import (
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/mcmathja/funky/batches"
	"github.com/mcmathja/funky/chans"
	"github.com/mcmathja/funky/maps"
	"github.com/mcmathja/funky/pairs"
	"github.com/mcmathja/funky/sets"
	"github.com/mcmathja/funky/slices"
)

func TestAll(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		input []int
		pred  func(int) bool
		all   bool
	}{
		"all match": {
			input: slices.New(2, 4, 6, 8),
			pred:  func(i int) bool { return i%2 == 0 },
			all:   true,
		},
		"none match": {
			input: slices.New(1, 3, 5, 7),
			pred:  func(i int) bool { return i%2 == 0 },
			all:   false,
		},
		"does not match in front": {
			input: slices.New(1, 2, 4, 6),
			pred:  func(i int) bool { return i%2 == 0 },
			all:   false,
		},
		"does not match in back": {
			input: slices.New(2, 4, 6, 7),
			pred:  func(i int) bool { return i%2 == 0 },
			all:   false,
		},
		"does not match in middle": {
			input: slices.New(2, 4, 5, 6),
			pred:  func(i int) bool { return i%2 == 0 },
			all:   false,
		},
		"does not match multiple times": {
			input: slices.New(1, 2, 4, 5, 6, 7),
			pred:  func(i int) bool { return i%2 == 0 },
			all:   false,
		},
		"empty input": {
			input: slices.New[int](),
			pred:  func(i int) bool { return i%2 == 0 },
			all:   true,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			all := slices.All(tc.input, tc.pred)

			if all != tc.all {
				t.Errorf("expected %t, but received %t", all, tc.all)
			}
		})
	}
}

func TestAny(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		input []int
		pred  func(int) bool
		any   bool
	}{
		"all match": {
			input: slices.New(2, 4, 6, 8),
			pred:  func(i int) bool { return i%2 == 0 },
			any:   true,
		},
		"none match": {
			input: slices.New(1, 3, 5, 7),
			pred:  func(i int) bool { return i%2 == 0 },
			any:   false,
		},
		"one matches in front": {
			input: slices.New(2, 3, 5, 7),
			pred:  func(i int) bool { return i%2 == 0 },
			any:   true,
		},
		"one matches in back": {
			input: slices.New(3, 5, 7, 8),
			pred:  func(i int) bool { return i%2 == 0 },
			any:   true,
		},
		"one matches in middle": {
			input: slices.New(3, 5, 6, 7),
			pred:  func(i int) bool { return i%2 == 0 },
			any:   true,
		},
		"multiple matches": {
			input: slices.New(2, 3, 5, 6, 7, 8),
			pred:  func(i int) bool { return i%2 == 0 },
			any:   true,
		},
		"empty input": {
			input: slices.New[int](),
			pred:  func(i int) bool { return i%2 == 0 },
			any:   false,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			any := slices.Any(tc.input, tc.pred)

			if any != tc.any {
				t.Errorf("expected %t, but received %t", any, tc.any)
			}
		})
	}
}

func TestAppend(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in     []int
		out    []int
		values []int
	}{
		"simple case": {
			in:     slices.New(1, 2, 3, 4),
			out:    slices.New(1, 2, 3, 4, 5),
			values: []int{5},
		},
		"multiple elements": {
			in:     slices.New(1, 2, 3, 4),
			out:    slices.New(1, 2, 3, 4, 5, 6, 7),
			values: []int{5, 6, 7},
		},
		"empty values": {
			in:     slices.New[int](),
			out:    slices.New[int](),
			values: []int{},
		},
		"nil values": {
			in:     slices.New[int](),
			out:    slices.New[int](),
			values: nil,
		},
		"empty input": {
			in:     slices.New[int](),
			out:    slices.New(3),
			values: []int{3},
		},
		"nil input": {
			in:     nil,
			out:    slices.New(3),
			values: []int{3},
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := slices.Append(tc.in, tc.values...)

			if !slices.Equal(out, tc.out) {
				t.Errorf(`expected %+v to equal %+v`, out, tc.out)
			}
		})
	}
}

func TestAtLeast(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		input   []int
		n       int
		pred    func(int) bool
		atLeast bool
	}{
		"n match": {
			input:   slices.New(1, 2, 3, 4, 5),
			n:       2,
			pred:    func(i int) bool { return i%2 == 0 },
			atLeast: true,
		},
		"more than n match": {
			input:   slices.New(1, 2, 3, 4, 5),
			n:       1,
			pred:    func(i int) bool { return i%2 == 0 },
			atLeast: true,
		},
		"fewer than n match": {
			input:   slices.New(1, 2, 3, 4, 5),
			n:       3,
			pred:    func(i int) bool { return i%2 == 0 },
			atLeast: false,
		},
		"n match when n is zero": {
			input:   slices.New(1, 2, 3, 4, 5),
			n:       0,
			pred:    func(i int) bool { return i > 6 },
			atLeast: true,
		},
		"negative n": {
			input:   slices.New(1, 2, 3, 4, 5),
			n:       -2,
			pred:    func(i int) bool { return i%2 == 0 },
			atLeast: true,
		},
		"empty input with non-zero count": {
			input:   slices.New[int](),
			pred:    func(i int) bool { return i%2 == 0 },
			n:       2,
			atLeast: false,
		},
		"nil input with non-zero n": {
			input:   slices.New[int](),
			pred:    func(i int) bool { return i%2 == 0 },
			n:       2,
			atLeast: false,
		},
		"empty input with zero n": {
			input:   slices.New[int](),
			pred:    func(i int) bool { return i%2 == 0 },
			n:       0,
			atLeast: true,
		},
		"nil input with zero n": {
			input:   nil,
			pred:    func(i int) bool { return i%2 == 0 },
			n:       0,
			atLeast: true,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			atLeast := slices.AtLeast(tc.input, tc.n, tc.pred)

			if atLeast != tc.atLeast {
				t.Errorf("expected %t, but received %t", atLeast, tc.atLeast)
			}
		})
	}
}

func TestAtMost(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		input  []int
		n      int
		pred   func(int) bool
		atMost bool
	}{
		"n match": {
			input:  slices.New(1, 2, 3, 4, 5),
			n:      2,
			pred:   func(i int) bool { return i%2 == 0 },
			atMost: true,
		},
		"more than n match": {
			input:  slices.New(1, 2, 3, 4, 5),
			n:      1,
			pred:   func(i int) bool { return i%2 == 0 },
			atMost: false,
		},
		"fewer than n match": {
			input:  slices.New(1, 2, 3, 4, 5),
			n:      3,
			pred:   func(i int) bool { return i%2 == 0 },
			atMost: true,
		},
		"n match when n is zero": {
			input:  slices.New(1, 2, 3, 4, 5),
			n:      0,
			pred:   func(i int) bool { return i > 6 },
			atMost: true,
		},
		"negative n": {
			input:  slices.New(1, 2, 3, 4, 5),
			n:      -2,
			pred:   func(i int) bool { return i%2 == 0 },
			atMost: false,
		},
		"empty input with non-zero count": {
			input:  slices.New[int](),
			pred:   func(i int) bool { return i%2 == 0 },
			n:      2,
			atMost: true,
		},
		"nil input with non-zero n": {
			input:  slices.New[int](),
			pred:   func(i int) bool { return i%2 == 0 },
			n:      2,
			atMost: true,
		},
		"empty input with zero n": {
			input:  slices.New[int](),
			pred:   func(i int) bool { return i%2 == 0 },
			n:      0,
			atMost: true,
		},
		"nil input with zero n": {
			input:  nil,
			pred:   func(i int) bool { return i%2 == 0 },
			n:      0,
			atMost: true,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			atMost := slices.AtMost(tc.input, tc.n, tc.pred)

			if atMost != tc.atMost {
				t.Errorf("expected %t, but received %t", atMost, tc.atMost)
			}
		})
	}
}

func TestCartesian(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		first  []int
		second []string
		out    []pairs.Pair[int, string]
	}{
		"simple case": {
			first:  slices.New(1, 2, 3),
			second: slices.New("one", "two", "three"),
			out: slices.New(
				pairs.New(1, "one"),
				pairs.New(1, "two"),
				pairs.New(1, "three"),
				pairs.New(2, "one"),
				pairs.New(2, "two"),
				pairs.New(2, "three"),
				pairs.New(3, "one"),
				pairs.New(3, "two"),
				pairs.New(3, "three"),
			),
		},
		"first longer": {
			first:  slices.New(1, 2, 3, 4),
			second: slices.New("one", "two", "three"),
			out: slices.New(
				pairs.New(1, "one"),
				pairs.New(1, "two"),
				pairs.New(1, "three"),
				pairs.New(2, "one"),
				pairs.New(2, "two"),
				pairs.New(2, "three"),
				pairs.New(3, "one"),
				pairs.New(3, "two"),
				pairs.New(3, "three"),
				pairs.New(4, "one"),
				pairs.New(4, "two"),
				pairs.New(4, "three"),
			),
		},
		"second longer": {
			first:  slices.New(1, 2, 3),
			second: slices.New("one", "two", "three", "four"),
			out: slices.New(
				pairs.New(1, "one"),
				pairs.New(1, "two"),
				pairs.New(1, "three"),
				pairs.New(1, "four"),
				pairs.New(2, "one"),
				pairs.New(2, "two"),
				pairs.New(2, "three"),
				pairs.New(2, "four"),
				pairs.New(3, "one"),
				pairs.New(3, "two"),
				pairs.New(3, "three"),
				pairs.New(3, "four"),
			),
		},
		"first empty": {
			first:  slices.New[int](),
			second: slices.New("one", "two", "three"),
			out:    slices.New[pairs.Pair[int, string]](),
		},
		"second empty": {
			first:  slices.New(1, 2, 3),
			second: slices.New[string](),
			out:    slices.New[pairs.Pair[int, string]](),
		},
		"both empty": {
			first:  slices.New[int](),
			second: slices.New[string](),
			out:    slices.New[pairs.Pair[int, string]](),
		},
		"first nil": {
			first:  nil,
			second: slices.New("one", "two", "three"),
			out:    slices.New[pairs.Pair[int, string]](),
		},
		"second nil": {
			first:  slices.New(1, 2, 3),
			second: nil,
			out:    slices.New[pairs.Pair[int, string]](),
		},
		"both nil": {
			first:  nil,
			second: nil,
			out:    slices.New[pairs.Pair[int, string]](),
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := slices.Cartesian(tc.first, tc.second)

			if !slices.Equal(out, tc.out) {
				t.Errorf("expected slice %v to equal %v, but did not", out, tc.out)
			}
		})
	}
}

func TestConsistsOf(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		slice      []int
		eles       []int
		consistsOf bool
	}{
		"equal slices": {
			slice:      slices.New(1, 2, 3, 4, 5),
			eles:       slices.New(1, 2, 3, 4, 5),
			consistsOf: true,
		},
		"slice rearranged": {
			slice:      slices.New(5, 3, 1, 2, 4),
			eles:       slices.New(1, 2, 3, 4, 5),
			consistsOf: true,
		},
		"eles rearranged": {
			slice:      slices.New(1, 2, 3, 4, 5),
			eles:       slices.New(5, 3, 1, 2, 4),
			consistsOf: true,
		},
		"eles missing element in slice": {
			slice:      slices.New(1, 2, 3, 4, 5),
			eles:       slices.New(1, 2, 3, 4),
			consistsOf: false,
		},
		"slice missing element in eles": {
			slice:      slices.New(1, 2, 3, 4, 5),
			eles:       slices.New(1, 2, 3, 4, 5, 6),
			consistsOf: true,
		},
		"single ele repeated": {
			slice:      slices.New(3, 3, 3, 3, 3),
			eles:       slices.New(3),
			consistsOf: true,
		},
		"slice has repeated elements": {
			slice:      slices.New(1, 2, 3, 3, 3, 4, 5),
			eles:       slices.New(1, 2, 3, 4, 5),
			consistsOf: true,
		},
		"eles has repeated elements": {
			slice:      slices.New(1, 2, 3, 4, 5),
			eles:       slices.New(1, 2, 3, 3, 3, 4, 5),
			consistsOf: true,
		},
		"slice empty": {
			slice:      slices.New[int](),
			eles:       slices.New(1, 2, 3, 4, 5),
			consistsOf: true,
		},
		"eles empty": {
			slice:      slices.New(1, 2, 3, 4, 5),
			eles:       slices.New[int](),
			consistsOf: false,
		},
		"both slice and eles empty": {
			slice:      slices.New[int](),
			eles:       slices.New[int](),
			consistsOf: true,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			consistsOf := slices.ConsistsOf(tc.slice, tc.eles...)

			if consistsOf != tc.consistsOf {
				t.Errorf(`returned %t, but expected %t`, consistsOf, tc.consistsOf)
			}
		})
	}
}

func TestComprises(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		slice     []int
		eles      []int
		comprises bool
	}{
		"equal slices": {
			slice:     slices.New(1, 2, 3, 4, 5),
			eles:      slices.New(1, 2, 3, 4, 5),
			comprises: true,
		},
		"slice rearranged": {
			slice:     slices.New(5, 3, 1, 2, 4),
			eles:      slices.New(1, 2, 3, 4, 5),
			comprises: true,
		},
		"eles rearranged": {
			slice:     slices.New(1, 2, 3, 4, 5),
			eles:      slices.New(5, 3, 1, 2, 4),
			comprises: true,
		},
		"eles missing element in slice": {
			slice:     slices.New(1, 2, 3, 4, 5),
			eles:      slices.New(1, 2, 3, 4, 6),
			comprises: false,
		},
		"slice missing element in eles": {
			slice:     slices.New(1, 2, 3, 4, 6),
			eles:      slices.New(1, 2, 3, 4, 5),
			comprises: false,
		},
		"single ele repeated": {
			slice:     slices.New(3, 3, 3, 3, 3),
			eles:      slices.New(3),
			comprises: false,
		},
		"slice has repeated elements": {
			slice:     slices.New(1, 2, 3, 3, 3, 4, 5),
			eles:      slices.New(1, 2, 3, 4, 5),
			comprises: false,
		},
		"eles has repeated elements": {
			slice:     slices.New(1, 2, 3, 4, 5),
			eles:      slices.New(1, 2, 3, 3, 3, 4, 5),
			comprises: false,
		},
		"slice empty": {
			slice:     slices.New[int](),
			eles:      slices.New(1, 2, 3, 4, 5),
			comprises: false,
		},
		"eles empty": {
			slice:     slices.New(1, 2, 3, 4, 5),
			eles:      slices.New[int](),
			comprises: false,
		},
		"both slice and eles empty": {
			slice:     slices.New[int](),
			eles:      slices.New[int](),
			comprises: true,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			comprises := slices.Comprises(tc.slice, tc.eles...)

			if comprises != tc.comprises {
				t.Errorf(`returned %t, but expected %t`, comprises, tc.comprises)
			}
		})
	}
}

func TestContains(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		slice    []int
		value    int
		contains bool
	}{
		"contains value": {
			slice:    slices.New(1, 2, 3, 4, 5),
			value:    4,
			contains: true,
		},
		"does not contain value": {
			slice:    slices.New(1, 2, 3, 4, 5),
			value:    6,
			contains: false,
		},
		"contains value at start": {
			slice:    slices.New(1, 2, 3, 4, 5),
			value:    1,
			contains: true,
		},
		"contains value at end": {
			slice:    slices.New(1, 2, 3, 4, 5),
			value:    5,
			contains: true,
		},
		"contains value multiple times": {
			slice:    slices.New(1, 2, 3, 4, 3, 2, 1),
			value:    2,
			contains: true,
		},
		"only value in slice": {
			slice:    slices.New(3),
			value:    3,
			contains: true,
		},
		"empty input": {
			slice:    slices.New[int](),
			value:    2,
			contains: false,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			contains := slices.Contains(tc.slice, tc.value)

			if contains != tc.contains {
				t.Errorf(`returned %t, but expected %t`, contains, tc.contains)
			}
		})
	}
}

func TestContainsAll(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		slice       []int
		values      []int
		containsAll bool
	}{
		"contains same values": {
			slice:       slices.New(1, 2, 3, 4, 5),
			values:      slices.New(1, 2, 3, 4, 5),
			containsAll: true,
		},
		"contains same values different order": {
			slice:       slices.New(1, 2, 3, 4, 5),
			values:      slices.New(5, 2, 1, 3, 4),
			containsAll: true,
		},
		"contains some values": {
			slice:       slices.New(1, 2, 3, 4, 5),
			values:      slices.New(2, 4, 5),
			containsAll: true,
		},
		"does not contain some values": {
			slice:       slices.New(1, 2, 3, 4, 5),
			values:      slices.New(1, 2, 8),
			containsAll: false,
		},
		"does not contain any values": {
			slice:       slices.New(1, 2, 3, 4, 5),
			values:      slices.New(9, -1, 12),
			containsAll: false,
		},
		"contains values multiple times": {
			slice:       slices.New(1, 2, 3, 4, 3, 2, 1),
			values:      slices.New(1, 2, 3),
			containsAll: true,
		},
		"contains sufficient repeated values": {
			slice:       slices.New(1, 2, 3, 1, 2, 3, 1, 2, 3),
			values:      slices.New(2, 2, 2),
			containsAll: true,
		},
		"contains more than sufficient repeated values": {
			slice:       slices.New(1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3),
			values:      slices.New(2, 2, 2),
			containsAll: true,
		},
		"contains insufficient repeated values": {
			slice:       slices.New(1, 2, 3, 1, 2, 3),
			values:      slices.New(2, 2, 2),
			containsAll: false,
		},
		"only value in slice": {
			slice:       slices.New(3),
			values:      slices.New(3),
			containsAll: true,
		},
		"empty input": {
			slice:       slices.New[int](),
			values:      slices.New(3),
			containsAll: false,
		},
		"empty values": {
			slice:       slices.New(3),
			values:      slices.New[int](),
			containsAll: true,
		},
		"both empty": {
			slice:       slices.New[int](),
			values:      slices.New[int](),
			containsAll: true,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			containsAll := slices.ContainsAll(tc.slice, tc.values...)

			if containsAll != tc.containsAll {
				t.Errorf(`returned %t, but expected %t`, containsAll, tc.containsAll)
			}
		})
	}
}

func TestContainsAny(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		slice       []int
		values      []int
		containsAny bool
	}{
		"contains same values": {
			slice:       slices.New(1, 2, 3, 4, 5),
			values:      slices.New(1, 2, 3, 4, 5),
			containsAny: true,
		},
		"contains same values different order": {
			slice:       slices.New(1, 2, 3, 4, 5),
			values:      slices.New(5, 2, 1, 3, 4),
			containsAny: true,
		},
		"contains some values": {
			slice:       slices.New(1, 2, 3, 4, 5),
			values:      slices.New(2, 4, 5),
			containsAny: true,
		},
		"does not contain some values": {
			slice:       slices.New(1, 2, 3, 4, 5),
			values:      slices.New(1, 2, 8),
			containsAny: true,
		},
		"does not contain any values": {
			slice:       slices.New(1, 2, 3, 4, 5),
			values:      slices.New(9, -1, 12),
			containsAny: false,
		},
		"contains values multiple times": {
			slice:       slices.New(1, 2, 3, 4, 3, 2, 1),
			values:      slices.New(1, 2, 3),
			containsAny: true,
		},
		"contains repeated value at least once": {
			slice:       slices.New(1, 2, 3),
			values:      slices.New(2, 2, 2),
			containsAny: true,
		},
		"does not contain repeated value": {
			slice:       slices.New(1, 2, 3),
			values:      slices.New(4, 4, 4),
			containsAny: false,
		},
		"only value in slice": {
			slice:       slices.New(3),
			values:      slices.New(3),
			containsAny: true,
		},
		"empty input": {
			slice:       slices.New[int](),
			values:      slices.New(3),
			containsAny: false,
		},
		"empty values": {
			slice:       slices.New(3),
			values:      slices.New[int](),
			containsAny: false,
		},
		"both empty": {
			slice:       slices.New[int](),
			values:      slices.New[int](),
			containsAny: false,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			containsAny := slices.ContainsAny(tc.slice, tc.values...)

			if containsAny != tc.containsAny {
				t.Errorf(`returned %t, but expected %t`, containsAny, tc.containsAny)
			}
		})
	}
}

func TestContainsSequence(t *testing.T) {
	t.Parallel()

	optCombos := map[string][]slices.ContainsSequenceOpt{
		"default":               {},
		"apostolico crochemore": {slices.ContainsSequenceApostolicoCrochemore},
		"brute force":           {slices.ContainsSequenceBruteForce},
	}

	testCases := map[string]struct {
		slice       []int
		seq         []int
		containsSeq bool
	}{
		"contains sequence": {
			slice:       slices.New(1, 2, 3, 4, 5),
			seq:         slices.New(2, 3, 4),
			containsSeq: true,
		},
		"does not contain sequence": {
			slice:       slices.New(1, 2, 3, 4, 5),
			seq:         slices.New(8, 9, 10),
			containsSeq: false,
		},
		"single element seq": {
			slice:       slices.New(1, 2, 3, 4, 5),
			seq:         slices.New(3),
			containsSeq: true,
		},
		"seq longer than slice": {
			slice:       slices.New(1, 2, 3, 4, 5),
			seq:         slices.New(1, 2, 3, 4, 5, 6),
			containsSeq: false,
		},
		"repeated element in middle": {
			slice:       slices.New(1, 3, 3, 3, 5),
			seq:         slices.New(3, 3, 3),
			containsSeq: true,
		},
		"repeated element in back": {
			slice:       slices.New(1, 2, 3, 3, 3),
			seq:         slices.New(3, 3, 3),
			containsSeq: true,
		},
		"repeated element in front": {
			slice:       slices.New(3, 3, 3, 4, 5),
			seq:         slices.New(3, 3, 3),
			containsSeq: true,
		},
		"does not contain repeated element": {
			slice:       slices.New(1, 2, 3, 4, 5),
			seq:         slices.New(3, 3, 3),
			containsSeq: false,
		},
		"partial match repeated element": {
			slice:       slices.New(1, 3, 3, 4, 5),
			seq:         slices.New(3, 3, 3),
			containsSeq: false,
		},
		"partial match in front": {
			slice:       slices.New(1, 2, 3, 4, 5),
			seq:         slices.New(2, 3, 6),
			containsSeq: false,
		},
		"partial match in back": {
			slice:       slices.New(1, 2, 3, 4, 5),
			seq:         slices.New(0, 3, 4),
			containsSeq: false,
		},
		"multiple partial matches": {
			slice:       slices.New(1, 2, 1, 2, 1),
			seq:         slices.New(1, 2, 3),
			containsSeq: false,
		},
		"partial match and full match": {
			slice:       slices.New(1, 2, 1, 2, 1, 2, 3),
			seq:         slices.New(1, 2, 3),
			containsSeq: true,
		},
		"repeated subsequence in seq": {
			slice:       slices.New(1, 2, 1, 3, 1, 2),
			seq:         slices.New(1, 2, 1, 2),
			containsSeq: false,
		},
		"partial subsequence in seq with match": {
			slice:       slices.New(1, 2, 1, 3, 1, 2, 1, 1),
			seq:         slices.New(1, 2, 1, 1),
			containsSeq: true,
		},
		"partial subsequence in seq without match": {
			slice:       slices.New(1, 2, 1, 3, 1, 2, 1, 2),
			seq:         slices.New(1, 2, 1, 1),
			containsSeq: false,
		},
		"slice nonempty seq empty": {
			slice:       slices.New(1, 2, 3, 4, 5),
			seq:         slices.New[int](),
			containsSeq: true,
		},
		"slice empty seq nonempty": {
			slice:       slices.New[int](),
			seq:         slices.New(2, 3, 4),
			containsSeq: false,
		},
		"slice empty seq empty": {
			slice:       slices.New[int](),
			seq:         slices.New[int](),
			containsSeq: true,
		},
	}

	for optsName, opts := range optCombos {
		opts := opts
		for testName, tc := range testCases {
			tc := tc
			t.Run(optsName+"/"+testName, func(t *testing.T) {
				t.Parallel()

				containsSeq := slices.ContainsSequence(tc.slice, tc.seq, opts...)

				if containsSeq != tc.containsSeq {
					t.Errorf(`returned %t, but expected %t`, containsSeq, tc.containsSeq)
				}
			})
		}
	}
}

func TestCorrespond(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		slice1     []int
		slice2     []int
		correspond bool
	}{
		"do correspond": {
			slice1:     slices.New(1, 2, 3, 4),
			slice2:     slices.New(2, 4, 6, 8),
			correspond: true,
		},
		"first slice mismatch front": {
			slice1:     slices.New(0, 2, 3, 4),
			slice2:     slices.New(2, 4, 6, 8),
			correspond: false,
		},
		"first slice mismatch back": {
			slice1:     slices.New(1, 2, 3, 5),
			slice2:     slices.New(2, 4, 6, 8),
			correspond: false,
		},
		"first slice extra elements": {
			slice1:     slices.New(1, 2, 3, 4, 5),
			slice2:     slices.New(2, 4, 6, 8),
			correspond: false,
		},
		"first slice fewer elements": {
			slice1:     slices.New(1, 2, 3),
			slice2:     slices.New(2, 4, 6, 8),
			correspond: false,
		},
		"first slice empty": {
			slice1:     slices.New[int](),
			slice2:     slices.New(2, 4, 6, 8),
			correspond: false,
		},
		"second slice mismatch front": {
			slice1:     slices.New(1, 2, 3, 4),
			slice2:     slices.New(1, 4, 6, 8),
			correspond: false,
		},
		"second slice mismatch back": {
			slice1:     slices.New(1, 2, 3, 4),
			slice2:     slices.New(2, 4, 6, 9),
			correspond: false,
		},
		"second slice extra elements": {
			slice1:     slices.New(1, 2, 3, 4),
			slice2:     slices.New(2, 4, 6, 8, 10),
			correspond: false,
		},
		"second slice fewer elements": {
			slice1:     slices.New(1, 2, 3, 4),
			slice2:     slices.New(2, 4, 6),
			correspond: false,
		},
		"second slice empty": {
			slice1:     slices.New(1, 2, 3, 4),
			slice2:     slices.New[int](),
			correspond: false,
		},
		"both slices empty": {
			slice1:     slices.New[int](),
			slice2:     slices.New[int](),
			correspond: true,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			correspond := slices.Correspond(tc.slice1, tc.slice2, func(a, b int) bool {
				return a*2 == b
			})

			if correspond != tc.correspond {
				t.Errorf(`returned %t, but expected %t`, correspond, tc.correspond)
			}
		})
	}
}

func TestCount(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		input []int
		pred  func(int) bool
		count int
	}{
		"all match": {
			input: slices.New(2, 4, 6, 8),
			pred:  func(i int) bool { return i%2 == 0 },
			count: 4,
		},
		"none match": {
			input: slices.New(1, 3, 5, 7),
			pred:  func(i int) bool { return i%2 == 0 },
			count: 0,
		},
		"one match": {
			input: slices.New(2, 3, 5, 7),
			pred:  func(i int) bool { return i%2 == 0 },
			count: 1,
		},
		"multiple matches": {
			input: slices.New(2, 3, 5, 6, 7, 8),
			pred:  func(i int) bool { return i%2 == 0 },
			count: 3,
		},
		"empty input": {
			input: slices.New[int](),
			pred:  func(i int) bool { return i%2 == 0 },
			count: 0,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			count := slices.Count(tc.input, tc.pred)

			if count != tc.count {
				t.Errorf("expected %d, but received %d", count, tc.count)
			}
		})
	}
}

func TestDistinct(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in  []int
		out []int
	}{
		"has repeating elements": {
			in:  slices.New(1, 1, 2, 2, 3, 3),
			out: slices.New(1, 2, 3),
		},
		"no repeating elements": {
			in:  slices.New(1, 2, 3, 4, 5),
			out: slices.New(1, 2, 3, 4, 5),
		},
		"preserves order": {
			in:  slices.New(4, 1, 4, 3, 1, 2, 2, 5),
			out: slices.New(4, 1, 3, 2, 5),
		},
		"empty input": {
			in:  slices.New[int](),
			out: slices.New[int](),
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := slices.Distinct(tc.in)

			if !slices.Equal(out, tc.out) {
				t.Errorf(`expected %+v to equal %+v`, out, tc.out)
			}
		})
	}
}

func TestDistinctBy(t *testing.T) {
	t.Parallel()

	type item struct {
		valueToCompare  int
		irrelevantValue int
	}

	testCases := map[string]struct {
		in  []item
		out []item
	}{
		"has repeating elements": {
			in:  slices.New([]item{{1, 0}, {1, 1}, {2, 0}, {2, 1}, {3, 0}, {3, 1}, {4, 0}, {4, 1}}...),
			out: slices.New([]item{{1, 0}, {2, 0}, {3, 0}, {4, 0}}...),
		},
		"no repeating elements": {
			in:  slices.New([]item{{1, 0}, {2, 0}, {3, 0}, {4, 0}, {5, 0}}...),
			out: slices.New([]item{{1, 0}, {2, 0}, {3, 0}, {4, 0}, {5, 0}}...),
		},
		"preserves order": {
			in:  slices.New([]item{{1, 1}, {1, 0}, {2, 0}, {2, 1}, {3, 0}, {5, 1}, {4, 1}, {5, 0}}...),
			out: slices.New([]item{{1, 1}, {2, 0}, {3, 0}, {5, 1}, {4, 1}}...),
		},
		"empty input": {
			in:  slices.New[item](),
			out: slices.New[item](),
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := slices.DistinctBy(tc.in, func(i item) int {
				return i.valueToCompare
			})

			if !slices.Equal(out, tc.out) {
				t.Errorf(`expected %+v to equal %+v`, out, tc.out)
			}
		})
	}
}

func TestDrop(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in  []int
		out []int
		n   int
	}{
		"simple case": {
			in:  slices.New(1, 2, 3, 4),
			out: slices.New(3, 4),
			n:   2,
		},
		"n is zero": {
			in:  slices.New(1, 2, 3, 4),
			out: slices.New(1, 2, 3, 4),
			n:   0,
		},
		"n equals size of slice": {
			in:  slices.New(1, 2, 3, 4),
			out: slices.New[int](),
			n:   4,
		},
		"n larger than size of slice": {
			in:  slices.New(1, 2, 3, 4),
			out: slices.New[int](),
			n:   100,
		},
		"n is negative": {
			in:  slices.New(1, 2, 3, 4),
			out: slices.New(1, 2, 3, 4),
			n:   -1,
		},
		"empty input": {
			in:  slices.New[int](),
			out: slices.New[int](),
			n:   2,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := slices.Drop(tc.in, tc.n)

			if !slices.Equal(out, tc.out) {
				t.Errorf(`expected %+v to equal %+v`, out, tc.out)
			}
		})
	}
}

func TestDropWhile(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in   []int
		out  []int
		pred func(i int) bool
	}{
		"simple case": {
			in:   slices.New(1, 2, 3, 4),
			out:  slices.New(3, 4),
			pred: func(i int) bool { return i < 3 },
		},
		"no false conditions": {
			in:   slices.New(1, 2, 3, 4),
			out:  slices.New[int](),
			pred: func(i int) bool { return i < 100 },
		},
		"no true conditions": {
			in:   slices.New(1, 2, 3, 4),
			out:  slices.New(1, 2, 3, 4),
			pred: func(i int) bool { return i < 0 },
		},
		"stops at first false": {
			in:   slices.New(1, 2, 3, 4, 5, 6, 7),
			out:  slices.New(3, 4, 5, 6, 7),
			pred: func(i int) bool { return i%3 != 0 },
		},
		"empty input": {
			in:   slices.New[int](),
			out:  slices.New[int](),
			pred: func(i int) bool { return i < 3 },
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := slices.DropWhile(tc.in, tc.pred)

			if !slices.Equal(out, tc.out) {
				t.Errorf(`expected %+v to equal %+v`, out, tc.out)
			}
		})
	}
}

func TestEmpty(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		slice []int
		empty bool
	}{
		"not empty": {
			slice: slices.New(1, 2, 3),
			empty: false,
		},
		"not empty with empty values": {
			slice: slices.New(0, 0, 0),
			empty: false,
		},
		"empty input": {
			slice: slices.New[int](),
			empty: true,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			empty := slices.Empty(tc.slice)

			if empty != tc.empty {
				t.Errorf(`expected %v to equal %v`, empty, tc.empty)
			}
		})
	}
}

func TestEndsWith(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		slice    []int
		value    int
		endsWith bool
	}{
		"ends with value": {
			slice:    slices.New(1, 2, 3, 4, 5),
			value:    5,
			endsWith: true,
		},
		"does not end with value": {
			slice:    slices.New(1, 2, 3, 4, 5),
			value:    6,
			endsWith: false,
		},
		"value not at end": {
			slice:    slices.New(1, 2, 3, 4, 5),
			value:    3,
			endsWith: false,
		},
		"only value in slice": {
			slice:    slices.New(2),
			value:    2,
			endsWith: true,
		},
		"empty input": {
			slice:    slices.New[int](),
			value:    2,
			endsWith: false,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			endsWith := slices.EndsWith(tc.slice, tc.value)

			if endsWith != tc.endsWith {
				t.Errorf(`returned %t, but expected %t`, endsWith, tc.endsWith)
			}
		})
	}
}

func TestEndsWithSequence(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		slice       []int
		seq         []int
		endsWithSeq bool
	}{
		"ends with sequence": {
			slice:       slices.New(1, 2, 3, 4, 5),
			seq:         slices.New(3, 4, 5),
			endsWithSeq: true,
		},
		"does not end with sequence": {
			slice:       slices.New(1, 2, 3, 4, 5),
			seq:         slices.New(8, 9, 10),
			endsWithSeq: false,
		},
		"partial match in front": {
			slice:       slices.New(1, 2, 3, 4, 5),
			seq:         slices.New(3, 4, 6),
			endsWithSeq: false,
		},
		"partial match in back": {
			slice:       slices.New(1, 2, 3, 4, 5),
			seq:         slices.New(0, 4, 5),
			endsWithSeq: false,
		},
		"slice nonempty seq empty": {
			slice:       slices.New(1, 2, 3, 4, 5),
			seq:         slices.New[int](),
			endsWithSeq: true,
		},
		"slice empty seq nonempty": {
			slice:       slices.New[int](),
			seq:         slices.New(3, 4, 5),
			endsWithSeq: false,
		},
		"slice empty seq empty": {
			slice:       slices.New[int](),
			seq:         slices.New[int](),
			endsWithSeq: true,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			endsWithSeq := slices.EndsWithSequence(tc.slice, tc.seq)

			if endsWithSeq != tc.endsWithSeq {
				t.Errorf(`returned %t, but expected %t`, endsWithSeq, tc.endsWithSeq)
			}
		})
	}
}

func TestEnumerate(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		slice []int
	}{
		"simple case": {
			slice: []int{1, -1, 2, -2, 3, -3},
		},
		"empty input": {
			slice: []int{},
		},
		"nil input": {
			slice: nil,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			iteration := 0
			slices.Enumerate(tc.slice, func(idx int, val int) {
				if idx != iteration {
					t.Errorf(`index %d did not match iteration %d`, idx, iteration)
				}

				if tc.slice[iteration] != val {
					t.Errorf(`value at iteration %d was %d, but expected %d`, iteration, val, tc.slice[iteration])
				}

				iteration++
			})

			if iteration != len(tc.slice) {
				t.Errorf(`length of slice was %d, but enumerated %d elements`, len(tc.slice), iteration)
			}
		})
	}
}

func TestEqual(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		slice1 []int
		slice2 []int
		equal  bool
	}{
		"do correspond": {
			slice1: slices.New(1, 2, 3, 4, 5),
			slice2: slices.New(1, 2, 3, 4, 5),
			equal:  true,
		},
		"first slice mismatch front": {
			slice1: slices.New(9, 2, 3, 4, 5),
			slice2: slices.New(1, 2, 3, 4, 5),
			equal:  false,
		},
		"first slice mismatch back": {
			slice1: slices.New(1, 2, 3, 4, 6),
			slice2: slices.New(1, 2, 3, 4, 5),
			equal:  false,
		},
		"first slice extra elements": {
			slice1: slices.New(1, 2, 3, 4, 5, 6),
			slice2: slices.New(1, 2, 3, 4, 5),
			equal:  false,
		},
		"first slice fewer elements": {
			slice1: slices.New(1, 2, 3, 4),
			slice2: slices.New(1, 2, 3, 4, 5),
			equal:  false,
		},
		"first slice empty": {
			slice1: slices.New[int](),
			slice2: slices.New(1, 2, 3, 4, 5),
			equal:  false,
		},
		"second slice mismatch front": {
			slice1: slices.New(1, 2, 3, 4, 5),
			slice2: slices.New(9, 2, 3, 4, 5),
			equal:  false,
		},
		"second slice mismatch back": {
			slice1: slices.New(1, 2, 3, 4, 5),
			slice2: slices.New(1, 2, 3, 4, 6),
			equal:  false,
		},
		"second slice extra elements": {
			slice1: slices.New(1, 2, 3, 4, 5),
			slice2: slices.New(1, 2, 3, 4, 5, 6),
			equal:  false,
		},
		"second slice fewer elements": {
			slice1: slices.New(1, 2, 3, 4, 5),
			slice2: slices.New(1, 2, 3, 4),
			equal:  false,
		},
		"second slice empty": {
			slice1: slices.New(1, 2, 3, 4, 5),
			slice2: slices.New[int](),
			equal:  false,
		},
		"both slices empty": {
			slice1: slices.New[int](),
			slice2: slices.New[int](),
			equal:  true,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			equal := slices.Equal(tc.slice1, tc.slice2)

			if equal != tc.equal {
				t.Errorf(`returned %t, but expected %t`, equal, tc.equal)
			}
		})
	}
}

func TestExactly(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		input   []int
		n       int
		pred    func(int) bool
		exactly bool
	}{
		"n match": {
			input:   slices.New(1, 2, 3, 4, 5),
			n:       2,
			pred:    func(i int) bool { return i%2 == 0 },
			exactly: true,
		},
		"more than n match": {
			input:   slices.New(1, 2, 3, 4, 5),
			n:       1,
			pred:    func(i int) bool { return i%2 == 0 },
			exactly: false,
		},
		"fewer than n match": {
			input:   slices.New(1, 2, 3, 4, 5),
			n:       3,
			pred:    func(i int) bool { return i%2 == 0 },
			exactly: false,
		},
		"n match when n is zero": {
			input:   slices.New(1, 2, 3, 4, 5),
			n:       0,
			pred:    func(i int) bool { return i > 6 },
			exactly: true,
		},
		"negative n": {
			input:   slices.New(1, 2, 3, 4, 5),
			n:       -2,
			pred:    func(i int) bool { return i%2 == 0 },
			exactly: false,
		},
		"empty input with non-zero count": {
			input:   slices.New[int](),
			pred:    func(i int) bool { return i%2 == 0 },
			n:       2,
			exactly: false,
		},
		"nil input with non-zero n": {
			input:   slices.New[int](),
			pred:    func(i int) bool { return i%2 == 0 },
			n:       2,
			exactly: false,
		},
		"empty input with zero n": {
			input:   slices.New[int](),
			pred:    func(i int) bool { return i%2 == 0 },
			n:       0,
			exactly: true,
		},
		"nil input with zero n": {
			input:   nil,
			pred:    func(i int) bool { return i%2 == 0 },
			n:       0,
			exactly: true,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			exactly := slices.Exactly(tc.input, tc.n, tc.pred)

			if exactly != tc.exactly {
				t.Errorf("expected %t, but received %t", exactly, tc.exactly)
			}
		})
	}
}

func TestFilter(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in   []int
		out  []int
		pred func(i int) bool
	}{
		"simple case": {
			in:   slices.New(1, 2, 3),
			out:  slices.New(1, 3),
			pred: func(i int) bool { return i != 2 },
		},
		"all elements pass": {
			in:   slices.New(1, 2, 3),
			out:  slices.New(1, 2, 3),
			pred: func(i int) bool { return true },
		},
		"no elements pass": {
			in:   slices.New(1, 2, 3),
			out:  slices.New[int](),
			pred: func(i int) bool { return false },
		},
		"empty input": {
			in:   slices.New[int](),
			out:  slices.New[int](),
			pred: func(i int) bool { return true },
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := slices.Filter(tc.in, tc.pred)

			if !slices.Equal(out, tc.out) {
				t.Errorf(`expected %v to equal %v`, out, tc.out)
			}
		})
	}
}

func TestFirst(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in    []int
		first int
		error bool
	}{
		"simple case": {
			in:    slices.New(1, 2, 3),
			first: 1,
		},
		"single element": {
			in:    slices.New(2),
			first: 2,
		},
		"empty input": {
			in:    slices.New[int](),
			error: true,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			first, err := slices.First(tc.in)

			if err == nil && tc.error {
				t.Errorf("should have errored, but did not")
			}

			if err != nil && !tc.error {
				t.Errorf("errored unexpectedly: %v", err)
			}

			if first != tc.first {
				t.Errorf(`expected %v to equal %v`, first, tc.first)
			}
		})
	}
}

func TestFirstIndexOf(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		slice []int
		value int
		idx   int
	}{
		"simple case": {
			slice: slices.New(1, 2, 3),
			value: 2,
			idx:   1,
		},
		"multiple matches": {
			slice: slices.New(1, 2, 3, 1, 2, 3),
			value: 2,
			idx:   1,
		},
		"no matches": {
			slice: slices.New(1, 2, 3),
			value: 0,
			idx:   -1,
		},
		"empty input": {
			slice: slices.New[int](),
			idx:   -1,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			idx := slices.FirstIndexOf(tc.slice, tc.value)

			if idx != tc.idx {
				t.Errorf(`returned %d, but expected %d`, idx, tc.idx)
			}
		})
	}
}

func TestFirstIndexWhere(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		slice []int
		idx   int
	}{
		"simple case": {
			slice: slices.New(1, 2, 3),
			idx:   1,
		},
		"multiple matches": {
			slice: slices.New(1, 2, 3, 4, 5, 6),
			idx:   1,
		},
		"no matches": {
			slice: slices.New(1, 3, 5, 7),
			idx:   -1,
		},
		"empty input": {
			slice: slices.New[int](),
			idx:   -1,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			idx := slices.FirstIndexWhere(tc.slice, func(ele int) bool {
				return ele%2 == 0
			})

			if idx != tc.idx {
				t.Errorf(`returned %d, but expected %d`, idx, tc.idx)
			}
		})
	}
}

func TestFlatMap(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in   []int
		out  []int
		pred func(i int) []int
	}{
		"simple case": {
			in:   slices.New(1, 2, 3),
			out:  slices.New(1, 2, 2, 4, 3, 6),
			pred: func(i int) []int { return slices.New(i, i*2) },
		},
		"variable size predicate result": {
			in:  slices.New(1, 2, 3),
			out: slices.New(2, 3, 4, 2, 4, 4, 5, 6),
			pred: func(i int) []int {
				if i == 2 {
					return slices.New(i, i*2)
				}
				return slices.New(i+1, i+2, i+3)
			},
		},
		"empty predicate result": {
			in:   slices.New(1, 2, 3),
			out:  slices.New[int](),
			pred: func(i int) []int { return slices.New[int]() },
		},
		"empty input": {
			in:   slices.New[int](),
			out:  slices.New[int](),
			pred: func(i int) []int { return slices.New(i, i*2) },
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := slices.FlatMap(tc.in, tc.pred)

			if !slices.Equal(out, tc.out) {
				t.Errorf(`expected %+v to equal %+v`, out, tc.out)
			}
		})
	}
}

func TestFlatten(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in  [][]int
		out []int
	}{
		"simple case": {
			in: slices.New(
				slices.New(1, 2, 3),
				slices.New(4, 5, 6),
				slices.New(7, 8, 9),
			),
			out: slices.New(1, 2, 3, 4, 5, 6, 7, 8, 9),
		},
		"some empty inner slices": {
			in: slices.New(
				slices.New(1, 2, 3),
				slices.New[int](),
				slices.New(4),
				slices.New(5, 6),
				slices.New[int](),
				slices.New(7, 8, 9),
			),
			out: slices.New(1, 2, 3, 4, 5, 6, 7, 8, 9),
		},
		"all empty inner slices": {
			in: slices.New(
				slices.New[int](),
				slices.New[int](),
				slices.New[int](),
			),
			out: slices.New[int](),
		},
		"no inner slices": {
			in:  slices.New[[]int](),
			out: slices.New[int](),
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := slices.Flatten(tc.in)

			if !slices.Equal(out, tc.out) {
				t.Errorf(`expected %+v to equal %+v`, out, tc.out)
			}
		})
	}
}

func TestForEach(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		slice []int
	}{
		"simple case": {
			slice: []int{1, -1, 2, -2, 3, -3},
		},
		"empty input": {
			slice: []int{},
		},
		"nil input": {
			slice: nil,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			iteration := 0
			slices.ForEach(tc.slice, func(val int) {
				if tc.slice[iteration] != val {
					t.Errorf(`value at iteration %d was %d, but expected %d`, iteration, val, tc.slice[iteration])
				}
				iteration++
			})

			if iteration != len(tc.slice) {
				t.Errorf(`length of slice was %d, but visited %d elements`, len(tc.slice), iteration)
			}
		})
	}
}

func TestFromBatch(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		out []int
		in  batches.Batch[int]
	}{
		"simple case": {
			out: slices.New(1, 2, 3, 4, 5),
			in:  batches.New(1, 2, 3, 4, 5),
		},
		"same element": {
			out: slices.New(0, 0, 0),
			in:  batches.New(0, 0, 0),
		},
		"empty batch": {
			out: slices.New[int](),
			in:  batches.New[int](),
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := slices.FromBatch(tc.in)

			if !slices.Equal(out, tc.out) {
				t.Errorf(`expected %v to equal %v`, out, tc.out)
			}
		})
	}
}

func TestFromChan(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		out []int
		in  <-chan int
	}{
		"simple case": {
			out: slices.New(1, 2, 3, 4, 5),
			in:  chans.New(1, 2, 3, 4, 5),
		},
		"same element": {
			out: slices.New(0, 0, 0),
			in:  chans.New(0, 0, 0),
		},
		"empty batch": {
			out: slices.New[int](),
			in:  chans.New[int](),
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := slices.FromChan(tc.in)

			if !slices.Equal(out, tc.out) {
				t.Errorf(`expected %v to equal %v`, out, tc.out)
			}
		})
	}
}

func TestFromMap(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		out []pairs.Pair[string, int]
		in  map[string]int
	}{
		"simple case": {
			out: slices.New(
				pairs.New("one", 1),
				pairs.New("two", 2),
				pairs.New("three", 3),
				pairs.New("four", 4),
				pairs.New("five", 5),
			),
			in: maps.New(
				pairs.New("one", 1),
				pairs.New("two", 2),
				pairs.New("three", 3),
				pairs.New("four", 4),
				pairs.New("five", 5),
			),
		},
		"same values": {
			out: slices.New(
				pairs.New("one", 1),
				pairs.New("two", 1),
				pairs.New("three", 1),
				pairs.New("four", 1),
				pairs.New("five", 1),
			),
			in: maps.New(
				pairs.New("one", 1),
				pairs.New("two", 1),
				pairs.New("three", 1),
				pairs.New("four", 1),
				pairs.New("five", 1),
			),
		},
		"mixed order": {
			out: slices.New(
				pairs.New("one", 1),
				pairs.New("five", 5),
				pairs.New("four", 4),
				pairs.New("two", 2),
				pairs.New("three", 3),
			),
			in: maps.New(
				pairs.New("five", 5),
				pairs.New("three", 3),
				pairs.New("two", 2),
				pairs.New("one", 1),
				pairs.New("four", 4),
			),
		},
		"no values": {
			out: slices.New[pairs.Pair[string, int]](),
			in:  map[string]int{},
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := slices.FromMap(tc.in)

			if !slices.ConsistsOf(out, tc.out...) {
				t.Errorf(`expected %v to consist of %v`, out, tc.out)
			}
		})
	}
}

func TestFromSet(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		out []int
		in  map[int]struct{}
	}{
		"simple case": {
			out: slices.New(1, 2, 3, 4, 5),
			in:  sets.New(1, 2, 3, 4, 5),
		},
		"mixed order": {
			out: slices.New(1, 5, 4, 2, 3),
			in:  sets.New(5, 3, 2, 1, 4),
		},
		"no values": {
			out: slices.New[int](),
			in:  map[int]struct{}{},
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := slices.FromSet(tc.in)

			if !slices.ConsistsOf(out, tc.out...) {
				t.Errorf(`expected %v to consist of %v`, out, tc.out)
			}
		})
	}
}

func TestGroupBy(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in  []int
		out map[string][]int
		fn  func(int) string
	}{
		"simple case": {
			in: slices.New(1, 2, 3, 4, 5),
			out: map[string][]int{
				"odd":  {1, 3, 5},
				"even": {2, 4},
			},
			fn: func(i int) string {
				if i%2 == 0 {
					return "even"
				} else {
					return "odd"
				}
			},
		},
		"single group": {
			in: slices.New(1, 2, 3, 4, 5),
			out: map[string][]int{
				"static": {1, 2, 3, 4, 5},
			},
			fn: func(i int) string {
				return "static"
			},
		},
		"no groups": {
			in: slices.New(1, 2, 3, 4, 5),
			out: map[string][]int{
				"1": {1},
				"2": {2},
				"3": {3},
				"4": {4},
				"5": {5},
			},
			fn: func(i int) string {
				return strconv.Itoa(i)
			},
		},
		"empty input": {
			in:  slices.New[int](),
			out: map[string][]int{},
			fn: func(i int) string {
				return strconv.Itoa(i)
			},
		},
		"nil input": {
			in:  nil,
			out: map[string][]int{},
			fn: func(i int) string {
				return strconv.Itoa(i)
			},
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := slices.GroupBy(tc.in, tc.fn)
			if len(tc.out) != len(out) {
				t.Errorf(`expected length of resulting map to be %d, but was %d`, len(tc.out), len(out))
			} else {
				for expectedKey, expectedValue := range tc.out {
					if actualValue, exists := out[expectedKey]; !exists {
						t.Errorf(`expected key %s to be present in output, but was missing`, expectedKey)
					} else if !slices.Equal(expectedValue, actualValue) {
						t.Errorf("expected %+v to equal %+v, but they differed", actualValue, expectedValue)
					}
				}
			}
		})
	}
}

func TestLast(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in    []int
		last  int
		error bool
	}{
		"simple case": {
			in:   slices.New(1, 2, 3),
			last: 3,
		},
		"single element": {
			in:   slices.New(2),
			last: 2,
		},
		"empty input": {
			in:    slices.New[int](),
			error: true,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			last, err := slices.Last(tc.in)

			if err == nil && tc.error {
				t.Errorf("should have errored, but did not")
			}

			if err != nil && !tc.error {
				t.Errorf("errored unexpectedly: %v", err)
			}

			if last != tc.last {
				t.Errorf(`expected %v to equal %v`, last, tc.last)
			}
		})
	}
}

func TestLastIndexOf(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		slice []int
		value int
		idx   int
	}{
		"simple case": {
			slice: slices.New(1, 2, 3),
			value: 2,
			idx:   1,
		},
		"multiple matches": {
			slice: slices.New(1, 2, 3, 1, 2, 3),
			value: 2,
			idx:   4,
		},
		"no matches": {
			slice: slices.New(1, 2, 3),
			value: 4,
			idx:   -1,
		},
		"empty input": {
			slice: slices.New[int](),
			idx:   -1,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			idx := slices.LastIndexOf(tc.slice, tc.value)

			if idx != tc.idx {
				t.Errorf(`returned %d, but expected %d`, idx, tc.idx)
			}
		})
	}
}

func TestLastIndexWhere(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		slice []int
		idx   int
	}{
		"simple case": {
			slice: slices.New(1, 2, 3),
			idx:   1,
		},
		"multiple matches": {
			slice: slices.New(1, 2, 3, 4, 5, 6),
			idx:   5,
		},
		"no matches": {
			slice: slices.New(1, 3, 5, 7),
			idx:   -1,
		},
		"empty input": {
			slice: slices.New[int](),
			idx:   -1,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			idx := slices.LastIndexWhere(tc.slice, func(ele int) bool {
				return ele%2 == 0
			})

			if idx != tc.idx {
				t.Errorf(`returned %d, but expected %d`, idx, tc.idx)
			}
		})
	}
}

func TestMap(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in   []int
		out  []int
		pred func(i int) int
	}{
		"simple case": {
			in:   slices.New(1, 2, 3),
			out:  slices.New(2, 4, 6),
			pred: func(i int) int { return i * 2 },
		},
		"identity": {
			in:   slices.New(1, 2, 3),
			out:  slices.New(1, 2, 3),
			pred: func(i int) int { return i },
		},
		"empty input": {
			in:   slices.New[int](),
			out:  slices.New[int](),
			pred: func(i int) int { return i },
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := slices.Map(tc.in, tc.pred)

			if !slices.Equal(out, tc.out) {
				t.Errorf(`expected %v to equal %v`, out, tc.out)
			}
		})
	}
}

func TestMax(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in    []int
		max   int
		error bool
	}{
		"simple case": {
			in:  slices.New(3, 1, 2, 5, 4),
			max: 5,
		},
		"max in front": {
			in:  slices.New(5, 4, 3, 2, 1),
			max: 5,
		},
		"max in back": {
			in:  slices.New(1, 2, 3, 4, 5),
			max: 5,
		},
		"negative values": {
			in:  slices.New(-3, -1, -2, -5, -4),
			max: -1,
		},
		"single element": {
			in:  slices.New(2),
			max: 2,
		},
		"empty input": {
			in:    slices.New[int](),
			error: true,
		},
		"nil input": {
			in:    nil,
			error: true,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			max, err := slices.Max(tc.in)

			if err == nil && tc.error {
				t.Errorf("should have errored, but did not")
			}

			if err != nil && !tc.error {
				t.Errorf("errored unexpectedly: %v", err)
			}

			if max != tc.max {
				t.Errorf(`expected %v to equal %v`, max, tc.max)
			}
		})
	}
}

func TestMin(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in    []int
		min   int
		error bool
	}{
		"simple case": {
			in:  slices.New(3, 1, 2, 5, 4),
			min: 1,
		},
		"min in back": {
			in:  slices.New(5, 4, 3, 2, 1),
			min: 1,
		},
		"min in front": {
			in:  slices.New(1, 2, 3, 4, 5),
			min: 1,
		},
		"negative values": {
			in:  slices.New(-3, -1, -2, -5, -4),
			min: -5,
		},
		"single element": {
			in:  slices.New(2),
			min: 2,
		},
		"empty input": {
			in:    slices.New[int](),
			error: true,
		},
		"nil input": {
			in:    nil,
			error: true,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			min, err := slices.Min(tc.in)

			if err == nil && tc.error {
				t.Errorf("should have errored, but did not")
			}

			if err != nil && !tc.error {
				t.Errorf("errored unexpectedly: %v", err)
			}

			if min != tc.min {
				t.Errorf(`expected %v to equal %v`, min, tc.min)
			}
		})
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		out    []int
		values []int
	}{
		"simple case": {
			out:    slices.New(1, 2, 3, 4, 5),
			values: []int{1, 2, 3, 4, 5},
		},
		"same element": {
			out:    slices.New(0, 0, 0),
			values: []int{0, 0, 0},
		},
		"empty slice": {
			out:    slices.New[int](),
			values: []int{},
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := slices.New(tc.values...)

			if !slices.Equal(out, tc.out) {
				t.Errorf(`expected %v to equal %v`, out, tc.out)
			}
		})
	}
}

func TestNthIndexOf(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		slice []int
		value int
		n     int
		idx   int
	}{
		"simple case": {
			slice: slices.New(1, 2, 3),
			value: 2,
			n:     1,
			idx:   1,
		},
		"first of multiple matches": {
			slice: slices.New(1, 2, 3, 1, 2, 3),
			n:     1,
			value: 2,
			idx:   1,
		},
		"second of multiple matches": {
			slice: slices.New(1, 2, 3, 1, 2, 3),
			n:     2,
			value: 2,
			idx:   4,
		},
		"n negative": {
			slice: slices.New(1, 2, 3, 1, 2, 3),
			n:     -1,
			value: 2,
			idx:   -1,
		},
		"n larger than number of matches in slice": {
			slice: slices.New(1, 2, 3, 1, 2, 3),
			n:     3,
			value: 2,
			idx:   -1,
		},
		"no matches": {
			slice: slices.New(1, 2, 3),
			n:     1,
			value: 4,
			idx:   -1,
		},
		"empty input": {
			slice: slices.New[int](),
			n:     1,
			value: 2,
			idx:   -1,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			idx := slices.NthIndexOf(tc.slice, tc.n, tc.value)

			if idx != tc.idx {
				t.Errorf(`returned %d, but expected %d`, idx, tc.idx)
			}
		})
	}
}

func TestNthIndexWhere(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		slice []int
		n     int
		idx   int
	}{
		"simple case": {
			slice: slices.New(1, 2, 3),
			n:     1,
			idx:   1,
		},
		"first of multiple matches": {
			slice: slices.New(1, 2, 3, 4, 5, 6),
			n:     1,
			idx:   1,
		},
		"second of multiple matches": {
			slice: slices.New(1, 2, 3, 4, 5, 6),
			n:     2,
			idx:   3,
		},
		"n negative": {
			slice: slices.New(1, 2, 3, 4, 5, 6),
			n:     -1,
			idx:   -1,
		},
		"n larger than number of matches in slice": {
			slice: slices.New(1, 2, 3, 4, 5, 6),
			n:     4,
			idx:   -1,
		},
		"no matches": {
			slice: slices.New(1, 3, 5, 7),
			n:     1,
			idx:   -1,
		},
		"empty input": {
			slice: slices.New[int](),
			idx:   -1,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			idx := slices.NthIndexWhere(tc.slice, tc.n, func(ele int) bool {
				return ele%2 == 0
			})

			if idx != tc.idx {
				t.Errorf(`returned %d, but expected %d`, idx, tc.idx)
			}
		})
	}
}

func TestPartition(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		input  []int
		pred   func(int) bool
		first  []int
		second []int
	}{
		"basic": {
			input:  slices.New(1, 2, 3, 4, 5),
			pred:   func(i int) bool { return i%2 == 0 },
			first:  slices.New(2, 4),
			second: slices.New(1, 3, 5),
		},
		"no elements that match": {
			input:  slices.New(1, 3, 5),
			pred:   func(i int) bool { return i%2 == 0 },
			first:  slices.New[int](),
			second: slices.New(1, 3, 5),
		},
		"no elements that don't match": {
			input:  slices.New(2, 4, 6),
			pred:   func(i int) bool { return i%2 == 0 },
			first:  slices.New(2, 4, 6),
			second: slices.New[int](),
		},
		"empty input": {
			input:  slices.New[int](),
			pred:   func(i int) bool { return i%2 == 0 },
			first:  slices.New[int](),
			second: slices.New[int](),
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			first, second := slices.Partition(tc.input, tc.pred)

			if !slices.Equal(first, tc.first) {
				t.Errorf("expected first slice %v to equal %v, but did not", first, tc.first)
			}

			if !slices.Equal(second, tc.second) {
				t.Errorf("expected second slice %v to equal %v, but did not", second, tc.second)
			}
		})
	}
}

func TestPermute(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in  []int
		out [][]int
	}{
		"simple case": {
			in: slices.New(1, 2, 3),
			out: slices.New(
				slices.New(1, 2, 3),
				slices.New(1, 3, 2),
				slices.New(2, 1, 3),
				slices.New(2, 3, 1),
				slices.New(3, 1, 2),
				slices.New(3, 2, 1),
			),
		},
		"reordered input": {
			in: slices.New(2, 3, 1),
			out: slices.New(
				slices.New(1, 2, 3),
				slices.New(1, 3, 2),
				slices.New(2, 1, 3),
				slices.New(2, 3, 1),
				slices.New(3, 1, 2),
				slices.New(3, 2, 1),
			),
		},
		"repeated elements": {
			in: slices.New(1, 3, 3),
			out: slices.New(
				slices.New(1, 3, 3),
				slices.New(1, 3, 3),
				slices.New(3, 1, 3),
				slices.New(3, 3, 1),
				slices.New(3, 1, 3),
				slices.New(3, 3, 1),
			),
		},
		"single element": {
			in:  slices.New(3),
			out: slices.New(slices.New(3)),
		},
		"empty input": {
			in:  slices.New[int](),
			out: slices.New(slices.New[int]()),
		},
		"nil input": {
			in:  nil,
			out: slices.New(slices.New[int]()),
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := slices.Permute(tc.in)
			actual := slices.Map(out, func(is []int) string {
				return strings.Join(slices.Map(is, func(i int) string {
					return strconv.Itoa(i)
				}), ",")
			})
			expected := slices.Map(tc.out, func(is []int) string {
				return strings.Join(slices.Map(is, func(i int) string {
					return strconv.Itoa(i)
				}), ",")
			})
			if !slices.Comprises(actual, expected...) {
				t.Errorf(`expected %v to be comprised of %v`, out, tc.out)
			}
		})
	}
}

func TestPrepend(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in    []int
		out   []int
		value int
	}{
		"simple case": {
			in:    slices.New(1, 2, 3, 4),
			out:   slices.New(0, 1, 2, 3, 4),
			value: 0,
		},
		"empty input": {
			in:    slices.New[int](),
			out:   slices.New(3),
			value: 3,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := slices.Prepend(tc.in, tc.value)

			if !slices.Equal(out, tc.out) {
				t.Errorf(`expected %+v to equal %+v`, out, tc.out)
			}
		})
	}
}

func TestProduct(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in      []int
		product int
	}{
		"simple case": {
			in:      slices.New(3, 1, 2, 5, 4),
			product: 120,
		},
		"commutative property": {
			in:      slices.New(3, 1, 2, 5, 4),
			product: 120,
		},
		"odd number of negative values": {
			in:      slices.New(-3, -1, -2, -5, -4),
			product: -120,
		},
		"even number of negative values": {
			in:      slices.New(-3, -1, -2, -5, -4, -1),
			product: 120,
		},
		"single element": {
			in:      slices.New(2),
			product: 2,
		},
		"empty input": {
			in:      slices.New[int](),
			product: 1,
		},
		"nil input": {
			in:      nil,
			product: 1,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			product := slices.Product(tc.in)

			if product != tc.product {
				t.Errorf(`expected %v to equal %v`, product, tc.product)
			}
		})
	}
}

func TestRange(t *testing.T) {
	t.Parallel()

	t.Run("int", func(t *testing.T) {
		testCases := map[string]struct {
			from, to, step int
			out            []int
		}{
			"simple case": {
				from: 1,
				to:   6,
				step: 1,
				out:  slices.New(1, 2, 3, 4, 5),
			},
			"larger step": {
				from: 1,
				to:   6,
				step: 2,
				out:  slices.New(1, 3, 5),
			},
			"negative step": {
				from: 6,
				to:   1,
				step: -2,
				out:  slices.New(6, 4, 2),
			},
			"zero step": {
				from: 1,
				to:   5,
				step: 0,
				out:  slices.New[int](),
			},
			"to equals from": {
				from: 6,
				to:   6,
				step: 1,
				out:  slices.New[int](),
			},
			"to less than from with positive step": {
				from: 6,
				to:   5,
				step: 1,
				out:  slices.New[int](),
			},
			"to larger than from with negative step": {
				from: 5,
				to:   6,
				step: -1,
				out:  slices.New[int](),
			},
			"positive step equal to difference": {
				from: 1,
				to:   5,
				step: 4,
				out:  slices.New(1),
			},
			"positive step larger than difference": {
				from: 1,
				to:   5,
				step: 6,
				out:  slices.New(1),
			},
			"negative step equal to difference": {
				from: 5,
				to:   1,
				step: -4,
				out:  slices.New(5),
			},
			"negative step larger than difference": {
				from: 5,
				to:   1,
				step: -6,
				out:  slices.New(5),
			},
		}

		for name, tc := range testCases {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				out := slices.Range(tc.from, tc.to, tc.step)

				if !slices.Equal(out, tc.out) {
					t.Errorf(`expected %v to equal %v`, out, tc.out)
				}
			})
		}
	})

	t.Run("float", func(t *testing.T) {
		testCases := map[string]struct {
			from, to, step float64
			out            []float64
		}{
			"simple case": {
				from: 1.3,
				to:   2.3,
				step: 0.2,
				out:  slices.New(1.3, 1.5, 1.7, 1.9, 2.1),
			},
			"larger step": {
				from: 1.3,
				to:   2.3,
				step: 0.5,
				out:  slices.New(1.3, 1.8),
			},
			"negative step": {
				from: 2.3,
				to:   1.3,
				step: -0.2,
				out:  slices.New(2.3, 2.1, 1.9, 1.7, 1.5),
			},
			"zero step": {
				from: 1.3,
				to:   2.3,
				step: 0,
				out:  slices.New[float64](),
			},
			"to equals from": {
				from: 2.3,
				to:   2.3,
				step: 0.2,
				out:  slices.New[float64](),
			},
			"to less than from with positive step": {
				from: 2.3,
				to:   2.1,
				step: 0.2,
				out:  slices.New[float64](),
			},
			"to larger than from with negative step": {
				from: 2.1,
				to:   2.3,
				step: -0.2,
				out:  slices.New[float64](),
			},
			"positive step equal to difference": {
				from: 1.3,
				to:   1.7,
				step: 0.4,
				out:  slices.New(1.3),
			},
			"positive step larger than difference": {
				from: 1.3,
				to:   1.7,
				step: 0.6,
				out:  slices.New(1.3),
			},
			"negative step equal to difference": {
				from: 1.7,
				to:   1.3,
				step: -0.4,
				out:  slices.New(1.7),
			},
			"negative step larger than difference": {
				from: 1.7,
				to:   1.3,
				step: -0.6,
				out:  slices.New(1.7),
			},
		}

		for name, tc := range testCases {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				out := slices.Range(tc.from, tc.to, tc.step)

				if !slices.Correspond(out, tc.out, func(a, b float64) bool {
					return math.Abs(a-b) < 0.0001
				}) {
					t.Errorf(`expected %v to equal %v`, out, tc.out)
				}
			})
		}
	})
}

func TestReduce(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in   []int
		init int
		pred func(acc, i int) int
		acc  int
	}{
		"simple case": {
			in:   slices.New(1, 2, 3),
			init: 0,
			pred: func(acc, i int) int { return acc + i },
			acc:  6,
		},
		"nonzero initial accumulator": {
			in:   slices.New(1, 2, 3),
			init: 10,
			pred: func(acc, i int) int { return acc + i },
			acc:  16,
		},
		"bespoke predicate logic": {
			in:   slices.New(6, 3, 5, 8, 9, 1, 2, 4, 7),
			init: 0,
			pred: func(acc, i int) int {
				if i > acc {
					return i
				}
				return acc
			},
			acc: 9,
		},
		"empty input": {
			in:   slices.New[int](),
			init: 0,
			pred: func(acc, i int) int { return acc + i },
			acc:  0,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			acc := slices.Reduce(tc.in, tc.init, tc.pred)

			if acc != tc.acc {
				t.Errorf(`expected %v to equal %v`, acc, tc.acc)
			}
		})
	}
}

func TestRepeat(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		out []int
		ele int
		num int
	}{
		"simple case": {
			out: slices.New(1, 1, 1, 1, 1),
			ele: 1,
			num: 5,
		},
		"num zero": {
			out: slices.New[int](),
			ele: 1,
			num: 0,
		},
		"num negative": {
			out: slices.New[int](),
			ele: 1,
			num: -10,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := slices.Repeat(tc.ele, tc.num)

			if !slices.Equal(out, tc.out) {
				t.Errorf(`expected %v to equal %v`, out, tc.out)
			}
		})
	}
}

func TestReversed(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in  []int
		out []int
	}{
		"simple case": {
			in:  slices.New(1, 2, 3),
			out: slices.New(3, 2, 1),
		},
		"single element": {
			in:  slices.New(3),
			out: slices.New(3),
		},
		"empty input": {
			in:  slices.New[int](),
			out: slices.New[int](),
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := slices.Reversed(tc.in)

			if !slices.Equal(out, tc.out) {
				t.Errorf(`expected %+v to equal %+v`, out, tc.out)
			}
		})
	}
}

func TestRotate(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in  []int
		out []int
		k   int
	}{
		"simple case": {
			in:  slices.New(1, 2, 3, 4, 5),
			out: slices.New(4, 5, 1, 2, 3),
			k:   2,
		},
		"k is zero": {
			in:  slices.New(1, 2, 3, 4, 5),
			out: slices.New(1, 2, 3, 4, 5),
			k:   0,
		},
		"k is n": {
			in:  slices.New(1, 2, 3, 4, 5),
			out: slices.New(1, 2, 3, 4, 5),
			k:   5,
		},
		"k is larger than n": {
			in:  slices.New(1, 2, 3, 4, 5),
			out: slices.New(5, 1, 2, 3, 4),
			k:   6,
		},
		"k is less than zero": {
			in:  slices.New(1, 2, 3, 4, 5),
			out: slices.New(3, 4, 5, 1, 2),
			k:   -2,
		},
		"k is less than -n": {
			in:  slices.New(1, 2, 3, 4, 5),
			out: slices.New(2, 3, 4, 5, 1),
			k:   -6,
		},
		"empty input": {
			in:  slices.New[int](),
			out: slices.New[int](),
			k:   2,
		},
		"nil input": {
			in:  nil,
			out: slices.New[int](),
			k:   2,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := slices.Rotate(tc.in, tc.k)

			if !slices.Equal(out, tc.out) {
				t.Errorf(`expected %v to equal %v`, out, tc.out)
			}
		})
	}
}

func TestSize(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		slice []int
		size  int
	}{
		"simple case": {
			slice: slices.New(1, 2, 3, 4, 5),
			size:  5,
		},
		"all empty values": {
			slice: slices.New(0, 0, 0, 0, 0),
			size:  5,
		},
		"empty input": {
			slice: slices.New[int](),
			size:  0,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			size := slices.Size(tc.slice)

			if size != tc.size {
				t.Errorf(`returned %d, but expected %d`, size, tc.size)
			}
		})
	}
}

func TestSort(t *testing.T) {
	t.Parallel()

	optCombos := map[string][]slices.SortOpt{
		"unstable": {},
		"stable":   {slices.SortStable},
	}

	testCases := map[string]struct {
		in  []int
		out []int
	}{
		"simple case": {
			in:  slices.New(3, 1, 2, 5, 4),
			out: slices.New(1, 2, 3, 4, 5),
		},
		"repeated elements": {
			in:  slices.New(3, 1, 3, 3, 4),
			out: slices.New(1, 3, 3, 3, 4),
		},
		"already sorted": {
			in:  slices.New(1, 2, 3, 4, 5),
			out: slices.New(1, 2, 3, 4, 5),
		},
		"empty input": {
			in:  slices.New[int](),
			out: slices.New[int](),
		},
		"nil input": {
			in:  nil,
			out: slices.New[int](),
		},
	}

	for optsName, opts := range optCombos {
		opts := opts
		for testName, tc := range testCases {
			tc := tc
			t.Run(optsName+"/"+testName, func(t *testing.T) {
				t.Parallel()

				out := slices.Sort(tc.in, opts...)

				if !slices.Equal(out, tc.out) {
					t.Errorf(`expected %+v to equal %+v`, out, tc.out)
				}
			})
		}
	}
}

func TestSortBy(t *testing.T) {
	t.Parallel()

	optCombos := map[string][]slices.SortByOpt{
		"unstable": {},
		"stable":   {slices.SortByStable},
	}

	testCases := map[string]struct {
		in  []int
		out []int
		fn  func(int, int) bool
	}{
		"simple case": {
			in:  slices.New(3, 1, 2, 5, 4),
			out: slices.New(5, 4, 3, 2, 1),
			fn: func(a, b int) bool {
				return b < a
			},
		},
		"repeated elements": {
			in:  slices.New(3, 1, 3, 3, 4),
			out: slices.New(4, 3, 3, 3, 1),
			fn: func(a, b int) bool {
				return b < a
			},
		},
		"already sorted": {
			in:  slices.New(5, 4, 3, 2, 1),
			out: slices.New(5, 4, 3, 2, 1),
			fn: func(a, b int) bool {
				return b < a
			},
		},
		"empty input": {
			in:  slices.New[int](),
			out: slices.New[int](),
			fn: func(a, b int) bool {
				return b < a
			},
		},
		"nil input": {
			in:  nil,
			out: slices.New[int](),
			fn: func(a, b int) bool {
				return b < a
			},
		},
	}

	for optsName, opts := range optCombos {
		opts := opts
		for testName, tc := range testCases {
			tc := tc
			t.Run(optsName+"/"+testName, func(t *testing.T) {
				t.Parallel()

				out := slices.SortBy(tc.in, tc.fn, opts...)

				if !slices.Equal(out, tc.out) {
					t.Errorf(`expected %+v to equal %+v`, out, tc.out)
				}
			})
		}
	}
}

func TestSplitAt(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		input  []int
		idx    int
		first  []int
		second []int
	}{
		"basic": {
			input:  slices.New(1, 2, 3, 4, 5),
			idx:    2,
			first:  slices.New(1, 2),
			second: slices.New(3, 4, 5),
		},
		"split at front": {
			input:  slices.New(1, 2, 3, 4, 5),
			idx:    0,
			first:  slices.New[int](),
			second: slices.New(1, 2, 3, 4, 5),
		},
		"split at back": {
			input:  slices.New(1, 2, 3, 4, 5),
			idx:    5,
			first:  slices.New(1, 2, 3, 4, 5),
			second: slices.New[int](),
		},
		"oob negative index": {
			input:  slices.New(1, 2, 3, 4, 5),
			idx:    -2,
			first:  slices.New[int](),
			second: slices.New(1, 2, 3, 4, 5),
		},
		"oob positive index": {
			input:  slices.New(1, 2, 3, 4, 5),
			idx:    10,
			first:  slices.New(1, 2, 3, 4, 5),
			second: slices.New[int](),
		},
		"empty slice": {
			input:  slices.New[int](),
			idx:    3,
			first:  slices.New[int](),
			second: slices.New[int](),
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			first, second := slices.SplitAt(tc.input, tc.idx)

			if !slices.Equal(first, tc.first) {
				t.Errorf("expected first slice %v to equal %v, but did not", first, tc.first)
			}

			if !slices.Equal(second, tc.second) {
				t.Errorf("expected second slice %v to equal %v, but did not", second, tc.second)
			}
		})
	}
}

func TestStartsWith(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		slice      []int
		value      int
		startsWith bool
	}{
		"starts with value": {
			slice:      slices.New(1, 2, 3, 4, 5),
			value:      1,
			startsWith: true,
		},
		"does not start with value": {
			slice:      slices.New(1, 2, 3, 4, 5),
			value:      0,
			startsWith: false,
		},
		"value not at start": {
			slice:      slices.New(1, 2, 3, 4, 5),
			value:      3,
			startsWith: false,
		},
		"only value in slice": {
			slice:      slices.New(2),
			value:      2,
			startsWith: true,
		},
		"empty input": {
			slice:      slices.New[int](),
			value:      2,
			startsWith: false,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			startsWith := slices.StartsWith(tc.slice, tc.value)

			if startsWith != tc.startsWith {
				t.Errorf(`returned %t, but expected %t`, startsWith, tc.startsWith)
			}
		})
	}
}

func TestStartsWithSequence(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		slice         []int
		seq           []int
		startsWithSeq bool
	}{
		"starts with sequence": {
			slice:         slices.New(1, 2, 3, 4, 5),
			seq:           slices.New(1, 2, 3),
			startsWithSeq: true,
		},
		"does not start with sequence": {
			slice:         slices.New(1, 2, 3, 4, 5),
			seq:           slices.New(8, 9, 10),
			startsWithSeq: false,
		},
		"partial match in front": {
			slice:         slices.New(1, 2, 3, 4, 5),
			seq:           slices.New(1, 2, 4),
			startsWithSeq: false,
		},
		"partial match in back": {
			slice:         slices.New(1, 2, 3, 4, 5),
			seq:           slices.New(0, 2, 3),
			startsWithSeq: false,
		},
		"slice nonempty seq empty": {
			slice:         slices.New(1, 2, 3, 4, 5),
			seq:           slices.New[int](),
			startsWithSeq: true,
		},
		"slice empty seq nonempty": {
			slice:         slices.New[int](),
			seq:           slices.New(1, 2, 3),
			startsWithSeq: false,
		},
		"slice empty seq empty": {
			slice:         slices.New[int](),
			seq:           slices.New[int](),
			startsWithSeq: true,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			startsWithSeq := slices.StartsWithSequence(tc.slice, tc.seq)

			if startsWithSeq != tc.startsWithSeq {
				t.Errorf(`returned %t, but expected %t`, startsWithSeq, tc.startsWithSeq)
			}
		})
	}
}

func TestSum(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in  []int
		sum int
	}{
		"simple case": {
			in:  slices.New(3, 1, 2, 5, 4),
			sum: 15,
		},
		"commutative property": {
			in:  slices.New(3, 1, 2, 5, 4),
			sum: 15,
		},
		"negative values": {
			in:  slices.New(-3, -1, -2, -5, -4),
			sum: -15,
		},
		"mix of positive and negative": {
			in:  slices.New(-3, 1, -2, 5, -4),
			sum: -3,
		},
		"single element": {
			in:  slices.New(2),
			sum: 2,
		},
		"empty input": {
			in:  slices.New[int](),
			sum: 0,
		},
		"nil input": {
			in:  nil,
			sum: 0,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			sum := slices.Sum(tc.in)

			if sum != tc.sum {
				t.Errorf(`expected %v to equal %v`, sum, tc.sum)
			}
		})
	}
}

func TestTake(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in  []int
		out []int
		n   int
	}{
		"simple case": {
			in:  slices.New(1, 2, 3, 4),
			out: slices.New(1, 2),
			n:   2,
		},
		"n is zero": {
			in:  slices.New(1, 2, 3, 4),
			out: slices.New[int](),
			n:   0,
		},
		"n equals size of slice": {
			in:  slices.New(1, 2, 3, 4),
			out: slices.New(1, 2, 3, 4),
			n:   4,
		},
		"n larger than size of slice": {
			in:  slices.New(1, 2, 3, 4),
			out: slices.New(1, 2, 3, 4),
			n:   100,
		},
		"n is negative": {
			in:  slices.New(1, 2, 3, 4),
			out: slices.New[int](),
			n:   -1,
		},
		"empty input": {
			in:  slices.New[int](),
			out: slices.New[int](),
			n:   2,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := slices.Take(tc.in, tc.n)

			if !slices.Equal(out, tc.out) {
				t.Errorf(`expected %+v to equal %+v`, out, tc.out)
			}
		})
	}
}

func TestTakeWhile(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in   []int
		out  []int
		pred func(i int) bool
	}{
		"simple case": {
			in:   slices.New(1, 2, 3, 4),
			out:  slices.New(1, 2),
			pred: func(i int) bool { return i < 3 },
		},
		"no false conditions": {
			in:   slices.New(1, 2, 3, 4),
			out:  slices.New(1, 2, 3, 4),
			pred: func(i int) bool { return i < 100 },
		},
		"no true conditions": {
			in:   slices.New(1, 2, 3, 4),
			out:  slices.New[int](),
			pred: func(i int) bool { return i < 0 },
		},
		"stops at first false": {
			in:   slices.New(1, 2, 3, 4, 5, 6, 7),
			out:  slices.New(1, 2),
			pred: func(i int) bool { return i%3 != 0 },
		},
		"empty input": {
			in:   slices.New[int](),
			out:  slices.New[int](),
			pred: func(i int) bool { return i < 3 },
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := slices.TakeWhile(tc.in, tc.pred)

			if !slices.Equal(out, tc.out) {
				t.Errorf(`expected %+v to equal %+v`, out, tc.out)
			}
		})
	}
}

func TestTally(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in  []int
		out map[int]int
	}{
		"simple case": {
			in: slices.New(1, 1, 2, 3, 3, 3),
			out: map[int]int{
				1: 2,
				2: 1,
				3: 3,
			},
		},
		"single group": {
			in: slices.New(5, 5, 5, 5),
			out: map[int]int{
				5: 4,
			},
		},
		"no groups": {
			in: slices.New(1, 2, 3, 4, 5),
			out: map[int]int{
				1: 1,
				2: 1,
				3: 1,
				4: 1,
				5: 1,
			},
		},
		"empty input": {
			in:  slices.New[int](),
			out: map[int]int{},
		},
		"nil input": {
			in:  nil,
			out: map[int]int{},
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := slices.Tally(tc.in)
			if len(tc.out) != len(out) {
				t.Errorf(`expected length of resulting map to be %d, but was %d`, len(tc.out), len(out))
			} else {
				for expectedKey, expectedValue := range tc.out {
					if actualValue, exists := out[expectedKey]; !exists {
						t.Errorf(`expected key %d to be present in output, but was missing`, expectedKey)
					} else if expectedValue != actualValue {
						t.Errorf("expected %d instances of %v, but had %d", actualValue, expectedKey, expectedValue)
					}
				}
			}
		})
	}
}

func TestTallyBy(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in  []int
		out map[string]int
		fn  func(int) string
	}{
		"simple case": {
			in: slices.New(1, 2, 3, 4, 5),
			out: map[string]int{
				"odd":  3,
				"even": 2,
			},
			fn: func(i int) string {
				if i%2 == 0 {
					return "even"
				} else {
					return "odd"
				}
			},
		},
		"single group": {
			in: slices.New(1, 2, 3, 4, 5),
			out: map[string]int{
				"static": 5,
			},
			fn: func(i int) string {
				return "static"
			},
		},
		"no groups": {
			in: slices.New(1, 2, 3, 4, 5),
			out: map[string]int{
				"1": 1,
				"2": 1,
				"3": 1,
				"4": 1,
				"5": 1,
			},
			fn: func(i int) string {
				return strconv.Itoa(i)
			},
		},
		"empty input": {
			in:  slices.New[int](),
			out: map[string]int{},
			fn: func(i int) string {
				return strconv.Itoa(i)
			},
		},
		"nil input": {
			in:  nil,
			out: map[string]int{},
			fn: func(i int) string {
				return strconv.Itoa(i)
			},
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := slices.TallyBy(tc.in, tc.fn)
			if len(tc.out) != len(out) {
				t.Errorf(`expected length of resulting map to be %d, but was %d`, len(tc.out), len(out))
			} else {
				for expectedKey, expectedValue := range tc.out {
					if actualValue, exists := out[expectedKey]; !exists {
						t.Errorf(`expected key %s to be present in output, but was missing`, expectedKey)
					} else if expectedValue != actualValue {
						t.Errorf("expected %+v to equal %+v, but they differed", actualValue, expectedValue)
					}
				}
			}
		})
	}
}

func TestTranspose(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in  [][]int
		out [][]int
		err bool
	}{
		"square matrix": {
			in: slices.New(
				slices.New(1, 2, 3),
				slices.New(4, 5, 6),
				slices.New(7, 8, 9),
			),
			out: slices.New(
				slices.New(1, 4, 7),
				slices.New(2, 5, 8),
				slices.New(3, 6, 9),
			),
		},
		"n > m": {
			in: slices.New(
				slices.New(1, 2, 3),
				slices.New(4, 5, 6),
			),
			out: slices.New(
				slices.New(1, 4),
				slices.New(2, 5),
				slices.New(3, 6),
			),
		},
		"m > n": {
			in: slices.New(
				slices.New(1, 2),
				slices.New(3, 4),
				slices.New(5, 6),
			),
			out: slices.New(
				slices.New(1, 3, 5),
				slices.New(2, 4, 6),
			),
		},
		"m == 0": {
			in:  slices.New[[]int](),
			out: slices.New[[]int](),
		},
		"n == 0": {
			in: slices.New(
				slices.New[int](),
				slices.New[int](),
				slices.New[int](),
			),
			out: slices.New[[]int](),
		},
		"m inconsistent": {
			in: slices.New(
				slices.New(1, 2, 3),
				slices.New(4, 6),
				slices.New(7, 8, 9),
			),
			err: true,
		},
		"input nil": {
			in:  nil,
			out: slices.New[[]int](),
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out, err := slices.Transpose(tc.in)

			if tc.err && err == nil {
				t.Errorf("should have errored, but did not")
			}

			if !tc.err && err != nil {
				t.Errorf("should not have errored, but got %v", err)
			}

			if !slices.Correspond(out, tc.out, slices.Equal[int]) {
				t.Errorf(`expected %v to equal %v`, out, tc.out)
			}
		})
	}
}

func TestUpdated(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in    []int
		out   []int
		idx   int
		value int
		error bool
	}{
		"simple case": {
			in:    slices.New(1, 2, 3),
			out:   slices.New(1, 100, 3),
			idx:   1,
			value: 100,
		},
		"first element": {
			in:    slices.New(1, 2, 3),
			out:   slices.New(100, 2, 3),
			idx:   0,
			value: 100,
		},
		"last element": {
			in:    slices.New(1, 2, 3),
			out:   slices.New(1, 2, 100),
			idx:   2,
			value: 100,
		},
		"only element": {
			in:    slices.New(1),
			out:   slices.New(100),
			idx:   0,
			value: 100,
		},
		"oob negative index": {
			in:    slices.New(1, 2, 3),
			idx:   -1,
			value: 100,
			error: true,
		},
		"oob positive index": {
			in:    slices.New(1, 2, 3),
			idx:   3,
			value: 100,
			error: true,
		},
		"empty input": {
			in:    slices.New[int](),
			idx:   0,
			value: 100,
			error: true,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out, err := slices.Updated(tc.in, tc.idx, tc.value)
			if tc.error && err == nil {
				t.Errorf("should have errored, but did not")
			}

			if !tc.error && err != nil {
				t.Errorf("should not have errored, but got %v", err)
			}

			if !slices.Equal(out, tc.out) {
				t.Errorf(`expected %v to equal %v`, out, tc.out)
			}
		})
	}
}

func TestZip(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		first  []int
		second []string
		out    []pairs.Pair[int, string]
	}{
		"simple case": {
			first:  slices.New(1, 2, 3),
			second: slices.New("one", "two", "three"),
			out: slices.New(
				pairs.New(1, "one"),
				pairs.New(2, "two"),
				pairs.New(3, "three"),
			),
		},
		"first longer": {
			first:  slices.New(1, 2, 3, 4),
			second: slices.New("one", "two", "three"),
			out: slices.New(
				pairs.New(1, "one"),
				pairs.New(2, "two"),
				pairs.New(3, "three"),
				pairs.New(4, ""),
			),
		},
		"second longer": {
			first:  slices.New(1, 2, 3),
			second: slices.New("one", "two", "three", "four"),
			out: slices.New(
				pairs.New(1, "one"),
				pairs.New(2, "two"),
				pairs.New(3, "three"),
				pairs.New(0, "four"),
			),
		},
		"first empty": {
			first:  slices.New[int](),
			second: slices.New("one", "two", "three"),
			out: slices.New(
				pairs.New(0, "one"),
				pairs.New(0, "two"),
				pairs.New(0, "three"),
			),
		},
		"second empty": {
			first:  slices.New(1, 2, 3),
			second: slices.New[string](),
			out: slices.New(
				pairs.New(1, ""),
				pairs.New(2, ""),
				pairs.New(3, ""),
			),
		},
		"both empty": {
			first:  slices.New[int](),
			second: slices.New[string](),
			out:    slices.New[pairs.Pair[int, string]](),
		},
		"first nil": {
			first:  nil,
			second: slices.New("one", "two", "three"),
			out: slices.New(
				pairs.New(0, "one"),
				pairs.New(0, "two"),
				pairs.New(0, "three"),
			),
		},
		"second nil": {
			first:  slices.New(1, 2, 3),
			second: nil,
			out: slices.New(
				pairs.New(1, ""),
				pairs.New(2, ""),
				pairs.New(3, ""),
			),
		},
		"both nil": {
			first:  nil,
			second: nil,
			out:    slices.New[pairs.Pair[int, string]](),
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := slices.Zip(tc.first, tc.second)

			if !slices.Equal(out, tc.out) {
				t.Errorf("expected slice %v to equal %v, but did not", out, tc.out)
			}
		})
	}
}
