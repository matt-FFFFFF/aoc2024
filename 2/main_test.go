package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsSafe(t *testing.T) {
	testCases := []struct {
		desc     string
		in       []int8
		errCount int
		result   bool
	}{
		{
			desc:     "simple",
			in:       []int8{1, 2, 3, 4},
			errCount: 0,
			result:   true,
		},
		{
			desc:     "simple-3",
			in:       []int8{1, 4, 7, 10},
			errCount: 0,
			result:   true,
		},
		{
			desc:     "simplerev-3",
			in:       []int8{10, 7, 4, 1},
			errCount: 0,
			result:   true,
		},
		{
			desc:     "reverse-error-diff",
			in:       []int8{10, 7, 3, 1},
			errCount: 2,
			result:   false,
		},
		{
			desc:     "reverse-dupe-single-error-ok",
			in:       []int8{10, 7, 7, 6},
			errCount: 1,
			result:   true,
		},
		{
			desc:     "reverse-increase-single-error-ok",
			in:       []int8{10, 7, 8, 6},
			errCount: 1,
			result:   true,
		},
		{
			desc:     "increase-flip-then-dupe-not-ok",
			in:       []int8{5, 4, 5, 6},
			errCount: 2,
			result:   false,
		},
		{
			desc:     "increase-flip-single-error-ok",
			in:       []int8{5, 6, 7, 5, 8},
			errCount: 1,
			result:   true,
		},
		{
			desc:     "increase-flip-single-error-ok",
			in:       []int8{8, 6, 9, 10, 11},
			errCount: 1,
			result:   true,
		},
	}
	setupLogger()
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			res, count := isSafe(tC.in)
			assert.Equalf(t, tC.errCount, count, "Error count incorrect: got %d want %d", count, tC.errCount)
			assert.Equalf(t, tC.result, res, "Result incorrect: got %t want %t", res, tC.result)
		})
	}
}
