package fibonacci

import (
	"fmt"
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
	{Ordinal: 13, Expected: 233},
	{Ordinal: 14, Expected: 377},
	{Ordinal: 15, Expected: 610},
	{Ordinal: 16, Expected: 987},
	{Ordinal: 17, Expected: 1597},
	{Ordinal: 18, Expected: 2584},
	{Ordinal: 19, Expected: 4181},
	{Ordinal: 20, Expected: 6765},
}

// Always behaves like the cache is empty
type MockEmptyCache struct {
}

func NewMockEmptyCache() *MockEmptyCache {
	return &MockEmptyCache{}
}

func (me *MockEmptyCache) Read(ordinal int) (int, error) {
	return -1, fmt.Errorf("Cache is empty")
}

func (me *MockEmptyCache) Write(ordinal int, value int) error {
	return nil
}

func (me *MockEmptyCache) Clear() error {
	return nil
}

type MockCache struct {
	table map[int]int
}

func NewMockCache(values map[int]int) *MockCache {
	if values == nil {
		return &MockCache{table: make(map[int]int)}
	}
	return &MockCache{table: values}
}

func (mc *MockCache) Write(ordinal int, value int) error {
	mc.table[ordinal] = value
	return nil
}

func (mc *MockCache) Read(ordinal int) (int, error) {
	if value, ok := mc.table[ordinal]; ok {
		return value, nil
	} else {
		return -1, fmt.Errorf("Value not in map")
	}
}

func (mc *MockCache) Clear() error {
	mc.table = make(map[int]int)
	return nil
}

func TestFibonacciNoCache(t *testing.T) {
	g := Generator{
		cache: NewMockEmptyCache(),
	}
	for _, v := range fibonacciTests {
		assert.Equal(t, v.Expected, g.Compute(v.Ordinal))
	}
}

func TestFibonacciCached(t *testing.T) {
	g := Generator{
		cache: NewMockCache(nil),
	}
	for _, v := range fibonacciTests {
		assert.Equal(t, v.Expected, g.Compute(v.Ordinal))
	}
}

func BenchmarkFibonacciNoCache(b *testing.B) {
	g := Generator{
		cache: NewMockEmptyCache(),
	}
	for i := 0; i < b.N; i++ {
		for _, v := range fibonacciTests {
			g.Compute(v.Ordinal)
		}
	}
}

func BenchmarkFibonacciCached(b *testing.B) {
	g := Generator{
		cache: NewMockCache(nil),
	}
	for i := 0; i < b.N; i++ {
		for _, v := range fibonacciTests {
			g.Compute(v.Ordinal)
		}
	}
}
