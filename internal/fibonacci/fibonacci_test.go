package fibonacci

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var fibonacciTests = []struct {
	Ordinal  uint64
	Expected *Number
}{
	{Ordinal: 0, Expected: NewNumber(0)},
	{Ordinal: 1, Expected: NewNumber(1)},
	{Ordinal: 2, Expected: NewNumber(1)},
	{Ordinal: 3, Expected: NewNumber(2)},
	{Ordinal: 4, Expected: NewNumber(3)},
	{Ordinal: 5, Expected: NewNumber(5)},
	{Ordinal: 6, Expected: NewNumber(8)},
	{Ordinal: 7, Expected: NewNumber(13)},
	{Ordinal: 8, Expected: NewNumber(21)},
	{Ordinal: 9, Expected: NewNumber(34)},
	{Ordinal: 10, Expected: NewNumber(55)},
	{Ordinal: 11, Expected: NewNumber(89)},
	{Ordinal: 12, Expected: NewNumber(144)},
	{Ordinal: 13, Expected: NewNumber(233)},
	{Ordinal: 14, Expected: NewNumber(377)},
	{Ordinal: 15, Expected: NewNumber(610)},
	{Ordinal: 16, Expected: NewNumber(987)},
	{Ordinal: 17, Expected: NewNumber(1597)},
	{Ordinal: 18, Expected: NewNumber(2584)},
	{Ordinal: 19, Expected: NewNumber(4181)},
	{Ordinal: 20, Expected: NewNumber(6765)},
}

// Always behaves like the cache is empty
type MockEmptyCache struct {
}

func NewMockEmptyCache() *MockEmptyCache {
	return &MockEmptyCache{}
}

func (me *MockEmptyCache) Read(ordinal uint64) (*Number, error) {
	return NewNumber(-1), fmt.Errorf("Cache is empty")
}

func (me *MockEmptyCache) Write(ordinal uint64, value *Number) error {
	return nil
}

func (me *MockEmptyCache) Clear() error {
	return nil
}

type MemoryCache struct {
	table map[uint64]*Number
}

func NewMemoryCache(values map[uint64]*Number) *MemoryCache {
	if values == nil {
		return &MemoryCache{table: make(map[uint64]*Number)}
	}
	return &MemoryCache{table: values}
}

func (mc *MemoryCache) Write(ordinal uint64, value *Number) error {
	mc.table[ordinal] = value
	return nil
}

func (mc *MemoryCache) Read(ordinal uint64) (*Number, error) {
	if value, ok := mc.table[ordinal]; ok {
		return value, nil
	} else {
		return NewNumber(-1), fmt.Errorf("Value not in map")
	}
}

func (mc *MemoryCache) Clear() error {
	mc.table = make(map[uint64]*Number)
	return nil
}

func TestFibonacciNoCache(t *testing.T) {
	g := NewGenerator(NewMockEmptyCache())
	for _, v := range fibonacciTests {
		assert.Equal(t, v.Expected, g.Compute(v.Ordinal))
	}
}

func TestFibonacciCached(t *testing.T) {
	g := NewGenerator(NewMemoryCache(nil))
	for _, v := range fibonacciTests {
		assert.Equal(t, v.Expected, g.Compute(v.Ordinal))
	}
}

func TestFibonacciLargeValue(t *testing.T) {
	g := NewGenerator(NewMemoryCache(nil))
	ord := uint64(100)
	v, _ := NewNumberFromDecimalString("354224848179261915075")

	assert.Equal(t, v, g.Compute(ord))
}

func TestFibonacciVeryLargeValue(t *testing.T) {
	g := NewGenerator(NewMemoryCache(nil))
	ord := uint64(1000)
	v, _ := NewNumberFromDecimalString("43466557686937456435688527675040625802564660517371780402481729089536555417949051890403879840079255169295922593080322634775209689623239873322471161642996440906533187938298969649928516003704476137795166849228875")

	assert.Equal(t, v, g.Compute(ord))
}

func BenchmarkFibonacciNoCache(b *testing.B) {
	g := NewGenerator(NewMockEmptyCache())
	for i := 0; i < b.N; i++ {
		for _, v := range fibonacciTests {
			g.Compute(v.Ordinal)
		}
	}
}

func BenchmarkFibonacciCached(b *testing.B) {
	g := NewGenerator(NewMemoryCache(nil))
	for i := 0; i < b.N; i++ {
		for _, v := range fibonacciTests {
			g.Compute(v.Ordinal)
		}
	}
}
