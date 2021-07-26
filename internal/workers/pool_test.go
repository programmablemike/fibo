package workers

import (
	"fmt"
	"sync"
	"testing"

	"github.com/programmablemike/fibo/internal/fibonacci"
	"github.com/stretchr/testify/assert"
)

// MemoryCache is a naive in-memory cache implementation
// It is *not* goroutine safe
type MemoryCache struct {
	table map[uint64]*fibonacci.Number
	mut   sync.Mutex // Use a mutex to prevent goroutines from clobbering each other
}

func NewMemoryCache(values map[uint64]*fibonacci.Number) *MemoryCache {
	if values == nil {
		return &MemoryCache{table: make(map[uint64]*fibonacci.Number)}
	}
	return &MemoryCache{table: values, mut: sync.Mutex{}}
}

func (mc *MemoryCache) Write(ordinal uint64, value *fibonacci.Number) error {
	mc.mut.Lock()
	mc.table[ordinal] = value
	mc.mut.Unlock()
	return nil
}

func (mc *MemoryCache) Read(ordinal uint64) (*fibonacci.Number, error) {
	mc.mut.Lock()
	defer mc.mut.Unlock()
	if value, ok := mc.table[ordinal]; ok {
		return value, nil
	} else {
		return fibonacci.NewNumber(-1), fmt.Errorf("Value not in map")
	}
}

func (mc *MemoryCache) Clear() error {
	mc.table = make(map[uint64]*fibonacci.Number)
	return nil
}

var fibonacciTests = []struct {
	Ordinal  uint64
	Expected *fibonacci.Number
}{
	{Ordinal: 0, Expected: fibonacci.NewNumber(0)},
	{Ordinal: 1, Expected: fibonacci.NewNumber(1)},
	{Ordinal: 2, Expected: fibonacci.NewNumber(1)},
	{Ordinal: 3, Expected: fibonacci.NewNumber(2)},
	{Ordinal: 4, Expected: fibonacci.NewNumber(3)},
	{Ordinal: 5, Expected: fibonacci.NewNumber(5)},
	{Ordinal: 6, Expected: fibonacci.NewNumber(8)},
	{Ordinal: 7, Expected: fibonacci.NewNumber(13)},
	{Ordinal: 8, Expected: fibonacci.NewNumber(21)},
	{Ordinal: 9, Expected: fibonacci.NewNumber(34)},
	{Ordinal: 10, Expected: fibonacci.NewNumber(55)},
	{Ordinal: 11, Expected: fibonacci.NewNumber(89)},
	{Ordinal: 12, Expected: fibonacci.NewNumber(144)},
	{Ordinal: 13, Expected: fibonacci.NewNumber(233)},
	{Ordinal: 14, Expected: fibonacci.NewNumber(377)},
	{Ordinal: 15, Expected: fibonacci.NewNumber(610)},
	{Ordinal: 16, Expected: fibonacci.NewNumber(987)},
	{Ordinal: 17, Expected: fibonacci.NewNumber(1597)},
	{Ordinal: 18, Expected: fibonacci.NewNumber(2584)},
	{Ordinal: 19, Expected: fibonacci.NewNumber(4181)},
	{Ordinal: 20, Expected: fibonacci.NewNumber(6765)},
}

var poolSizes []int = []int{1, 2, 4, 8, 16}

func TestWorkerPool(t *testing.T) {
	for _, poolSize := range poolSizes {
		fmt.Printf("Creating a pool of size %d\n", poolSize)
		c := NewMemoryCache(nil)
		gp := NewGeneratorPool(c, poolSize)
		in, out := gp.GetInputOutput()

		fmt.Println("Starting the pool")
		assert.Equal(t, false, gp.IsRunning())
		gp.Start() // Start the worker pool
		assert.Equal(t, true, gp.IsRunning())

		for _, v := range fibonacciTests {
			fmt.Printf("Computing the ordinal %v\n", v.Ordinal)
			in <- v.Ordinal
			result := <-out
			fmt.Printf("Received value of %v\n", result)
			assert.Equal(t, v.Expected, result)
		}

		fmt.Println("Stopping the pool")
		assert.Equal(t, true, gp.IsRunning())
		gp.Close() // Stop the pool
		assert.Equal(t, false, gp.IsRunning())
	}
}
