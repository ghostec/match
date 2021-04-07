package match_test

import (
	"fmt"
	"reflect"
	"testing"

	ma "github.com/ghostec/match"
	ha "github.com/ghostec/match/handles"
	ty "github.com/ghostec/match/types"
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
				t.Errorf("unexpected output. Got: %d, Expected: %d", fn, tc.fn)
			}
		})
	}
}

var fibonacci = ma.Match{
	{ma.When(0), 0},
	{ma.When(1), 1},
	{ma.When(ha.Any(0)), func(m ma.Match, n int) int { return m.Int(n-1) + m.Int(n-2) }},
}

func Fibonacci(n int) (fn int) {
	return fibonacci.Int(n)
}

func TestJoin(t *testing.T) {
	testCases := []struct {
		input  interface{}
		by     string
		output string
	}{
		{
			input:  []int{0, 1, 2, 3, 4, 5},
			by:     "",
			output: "012345",
		},
		{
			input:  []string{"Mr", "John", "Doe"},
			by:     " ",
			output: "Mr John Doe",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.output, func(t *testing.T) {
			output := Join(ty.NewSlice(tc.input), tc.by)
			if output != tc.output {
				t.Errorf("unexpected output. Got: %s, Expected: %s", output, tc.output)
			}
		})
	}
}

var join = ma.Match{
	{
		// 0 - [], "by", "acc" -> "acc"
		ma.When(ha.Empty(), ha.String(), ha.String(0)),
		func(joined string) string { return joined },
	},
	{
		// 1 - [head|...rest], "by" -> [...], "by", "acc" // {0, 2}
		ma.When(ha.Slice(ha.Head(0, 1), ha.Slice(1)), ha.String(2)),
		func(m ma.Match, head *ty.Any, tail *ty.Slice, by string) string {
			return m.String(tail, by, fmt.Sprintf("%v", head.Get()))
		},
	},
	{
		// 2 - [head|...rest], "by", "acc" -> [...], "by", "acc" // {0, 2}
		ma.When(ha.Slice(ha.Head(0, 1), ha.Slice(1)), ha.String(2), ha.String(3)),
		func(m ma.Match, head *ty.Any, tail *ty.Slice, by, acc string) string {
			return m.String(tail, by, fmt.Sprintf("%s%s%v", acc, by, head.Get()))
		},
	},
}

func Join(slice *ty.Slice, by string) string {
	return join.String(slice, by)
}

func TestReverse(t *testing.T) {
	testCases := []struct {
		input  interface{}
		output interface{}
	}{
		{
			input:  []int{0, 1, 2, 3, 4, 5},
			output: []int{5, 4, 3, 2, 1, 0},
		},
	}

	for _, tc := range testCases {
		t.Run(Join(ty.NewSlice(tc.input), ", "), func(t *testing.T) {
			output := Reverse(ty.NewSlice(tc.input))
			if !reflect.DeepEqual(output.Get(), tc.output) {
				t.Errorf("unexpected output. Got: %s, Expected: %s", output, tc.output)
			}
		})
	}
}

var reverse = ma.Match{
	{
		ma.When(ha.Empty(), ha.Slice(0)),
		func(reversed *ty.Slice) *ty.Slice { return reversed },
	},
	{
		ma.When(ha.Slice(ha.Slice(0), ha.Tail(1, 1))),
		func(m ma.Match, head *ty.Slice, tail *ty.Any) *ty.Slice {
			return m.Slice(head, ty.NewSlice(tail.Get()))
		},
	},
	{
		ma.When(ha.Slice(ha.Slice(0), ha.Tail(1, 1)), ha.Slice(2)),
		func(m ma.Match, head *ty.Slice, tail *ty.Any, acc *ty.Slice) *ty.Slice {
			return m.Slice(head, acc.Append(tail.Get()))
		},
	},
}

func Reverse(list *ty.Slice) *ty.Slice {
	return reverse.Slice(list)
}
