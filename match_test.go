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
			fn := Fibonacci(tc.n)
			if fn != tc.fn {
				t.Errorf("unexpected")
			}
		})
	}
}

var fib = match.
	When(0).Return(0).
	When(1).Return(1).
	When(match.Any).Return(func(m *match.Match, n int) int { return m.Int(n-1) + m.Int(n-2) })

func Fibonacci(n int) int { return fib.Int(n) }
