package match_test

import (
	"fmt"
	"testing"

	"github.com/ghostec/match"
)

func TestFibonacci(t *testing.T) {
	testCases := []struct {
		n  int
		fn int
	}{
		{0, 0},
		{1, 1},
		{2, 1},
		{3, 2},
		{4, 3},
		{5, 5},
		{6, 8},
		{7, 13},
		{8, 21},
		{9, 34},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%d", tc.n), func(t *testing.T) {
			fn := fib(tc.n)
			if fn != tc.fn {
				t.Errorf("unexpected")
			}
		})
	}
}

func fib(n int) int {
	return match.
		Match(n).
		When(0).Return(0).
		When(1).Return(1).
		When(match.Any).Do(func(n int) int { return fib(n-1) + fib(n-2) }).
		Int()
}

// func TestReverseList(t *testing.T) {
// 	testCases := []struct {
// 		input  []int
// 		output []int
// 	}{
// 		{input: []int{0, 1, 2}, output: []int{2, 1, 0}},
// 	}
//
// 	for _, tc := range testCases {
// 		t.Run(fmt.Sprintf("%d", tc.n), func(t *testing.T) {
// 			output := reverse(tc.input)
// 			if !reflect.DeepEqual(output, tc.output) {
// 				t.Errorf("unexpected")
// 			}
// 		})
// 	}
// }
