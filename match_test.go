package match_test

import (
	"fmt"
	"reflect"
	"testing"

	ma "github.com/ghostec/match"
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
	{ma.When(ma.Any(0)), func(m ma.Match, n int) int {
		return m.Int(n-1) + m.Int(n-2)
	}},
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
			output := Join(tc.input, tc.by)
			if output != tc.output {
				t.Errorf("unexpected output. Got: %s, Expected: %s", output, tc.output)
			}
		})
	}
}

var join = ma.Match{
	{
		ma.When(ma.Empty, ma.String, ma.String(0)),
		func(joined string) string {
			return joined
		},
	},
	{
		ma.When(ma.Slice(0), ma.String(1)),
		func(m ma.Match, list interface{}, by string) string {
			return m.String(list, by, "")
		},
	},
	{
		ma.When(
			ma.Slice(ma.Head(0, 1), ma.Slice(1)),
			ma.String(2),
			ma.Empty,
		),
		func(m ma.Match, head, tail interface{}, by string) string {
			return m.String(tail, by, fmt.Sprintf("%v", head))
		},
	},
	{
		ma.When(
			ma.Slice(ma.Head(0, 1), ma.Slice(1)),
			ma.String(2),
			ma.String(3),
		),
		func(m ma.Match, head, tail interface{}, by, acc string) string {
			return m.String(tail, by, fmt.Sprintf("%s%s%v", acc, by, head))
		},
	},
}

func Join(slice interface{}, by string) string {
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
		t.Run(Join(tc.input, ", "), func(t *testing.T) {
			output := Reverse(tc.input)
			if reflect.DeepEqual(output, tc.output) {
				t.Errorf("unexpected output. Got: %s, Expected: %s", output, tc.output)
			}
		})
	}
}

var reverse = ma.Match{
	{
		ma.When(ma.Empty, ma.Slice(0)),
		func(reversed *ma.SliceType) *ma.SliceType {
			return reversed
		},
	},
	{
		ma.When(
			ma.Slice(ma.Slice(0), ma.Tail(1, 1)),
		),
		func(m ma.Match, head *ma.SliceType, tail interface{}) *ma.SliceType {
			return m.Slice(head, tail)
		},
	},
	{
		ma.When(
			ma.Slice(ma.Slice(0), ma.Tail(1, 1)),
			ma.Slice(2),
		),
		func(m ma.Match, head *ma.SliceType, tail interface{}, acc *ma.SliceType) *ma.SliceType {
			return m.Slice(head, acc.Append(tail))
		},
	},
}

func Reverse(list interface{}) interface{} {
	return reverse.Slice(list).Get()
}
