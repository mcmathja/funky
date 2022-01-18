package sets_test

import (
	"testing"

	"github.com/mcmathja/funky/sets"
)

func TestAdd(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		in     map[int]struct{}
		out    map[int]struct{}
		values []int
	}{
		"simple case": {
			in:     sets.New(1, 2, 3, 4),
			out:    sets.New(1, 2, 3, 4, 5),
			values: []int{5},
		},
		"multiple elements": {
			in:     sets.New(1, 2, 3, 4),
			out:    sets.New(1, 2, 3, 4, 5, 6, 7),
			values: []int{5, 6, 7},
		},
		"duplicate elements": {
			in:     sets.New(1, 2, 3, 4),
			out:    sets.New(1, 2, 3, 4, 5),
			values: []int{5, 5, 5},
		},
		"empty values": {
			in:     sets.New[int](),
			out:    sets.New[int](),
			values: []int{},
		},
		"nil values": {
			in:     sets.New[int](),
			out:    sets.New[int](),
			values: nil,
		},
		"empty input": {
			in:     sets.New[int](),
			out:    sets.New(3),
			values: []int{3},
		},
		"nil input": {
			in:     nil,
			out:    sets.New(3),
			values: []int{3},
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := sets.Add(tc.in, tc.values...)

			if !sets.Equals(out, tc.out) {
				t.Errorf(`expected %+v to equal %+v`, out, tc.out)
			}
		})
	}
}

func TestAll(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		input map[int]struct{}
		pred  func(int) bool
		all   bool
	}{
		"all match": {
			input: sets.New(2, 4, 6, 8),
			pred:  func(i int) bool { return i%2 == 0 },
			all:   true,
		},
		"none match": {
			input: sets.New(1, 3, 5, 7),
			pred:  func(i int) bool { return i%2 == 0 },
			all:   false,
		},
		"some match": {
			input: sets.New(1, 2, 4, 5, 6, 7),
			pred:  func(i int) bool { return i%2 == 0 },
			all:   false,
		},
		"empty input": {
			input: sets.New[int](),
			pred:  func(i int) bool { return i%2 == 0 },
			all:   true,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			all := sets.All(tc.input, tc.pred)

			if all != tc.all {
				t.Errorf("expected %t, but received %t", all, tc.all)
			}
		})
	}
}

func TestAny(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		input map[int]struct{}
		pred  func(int) bool
		any   bool
	}{
		"all match": {
			input: sets.New(2, 4, 6, 8),
			pred:  func(i int) bool { return i%2 == 0 },
			any:   true,
		},
		"none match": {
			input: sets.New(1, 3, 5, 7),
			pred:  func(i int) bool { return i%2 == 0 },
			any:   false,
		},
		"some match": {
			input: sets.New(1, 2, 4, 5, 6, 7),
			pred:  func(i int) bool { return i%2 == 0 },
			any:   true,
		},
		"empty input": {
			input: sets.New[int](),
			pred:  func(i int) bool { return i%2 == 0 },
			any:   false,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			any := sets.Any(tc.input, tc.pred)

			if any != tc.any {
				t.Errorf("expected %t, but received %t", any, tc.any)
			}
		})
	}
}
