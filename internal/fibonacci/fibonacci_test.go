package fibonacci

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var fibonacciTests = []struct {
	Ordinal  int
	Expected int
}{
	{Ordinal: -2, Expected: -1},
	{Ordinal: -1, Expected: -1},
	{Ordinal: 0, Expected: 0},
	{Ordinal: 1, Expected: 1},
	{Ordinal: 2, Expected: 1},
	{Ordinal: 3, Expected: 2},
	{Ordinal: 4, Expected: 3},
	{Ordinal: 5, Expected: 5},
	{Ordinal: 6, Expected: 8},
	{Ordinal: 7, Expected: 13},
	{Ordinal: 8, Expected: 21},
	{Ordinal: 9, Expected: 34},
	{Ordinal: 10, Expected: 55},
	{Ordinal: 11, Expected: 89},
	{Ordinal: 12, Expected: 144},
}

func TestFibonacciOk(t *testing.T) {
	g := Generator{}
	for _, v := range fibonacciTests {
		assert.Equal(t, v.Expected, g.Compute(v.Ordinal))
	}
}
