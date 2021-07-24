// Calculates the Fibonacci sequence
// Defined as f(n) = f(n-1) + f(n-2) for n > 1
package fibonacci

import (
	"strconv"

	log "github.com/sirupsen/logrus"
)

type Memoizer interface {
	Write(ordinal int64, value int64) error
	Read(ordinal int64) (int64, error)
	Clear() error
}

type Generator struct {
	cache Memoizer
}

func NewGenerator(cache Memoizer) *Generator {
	return &Generator{
		cache: cache,
	}
}

// ClearCache wipes the memoizer's Postgres DB
func (g *Generator) ClearCache() error {
	return g.cache.Clear()
}

func Int64ToString(v int64) string {
	return strconv.FormatInt(v, 10)
}

// Compute Get the fibonacci value for the given ordinal
// Defined as f(n) = f(n-2) + f(n-1) where f(0) = 0 and f(1) = 1
// Returns -1 for invalid values
func (g Generator) Compute(n int64) int64 {
	log.Debugf("Computing fibonacci sequence for ordinal=%s", Int64ToString(n))

	switch {
	case n == 0:
		return 0
	case n == 1:
		return 1
	case n > 1:
		// @TODO: Insert caching logic here
		n1 := g.readCachedOrCompute(n - 1)
		n2 := g.readCachedOrCompute(n - 2)
		return n1 + n2
	default:
		return -1
	}
}

// readCachedOrCompute will read a value from the database if it exists
// otherwise it will compute the value and store it in the cache for future use
func (g *Generator) readCachedOrCompute(ordinal int64) int64 {
	value, err := g.cache.Read(ordinal)
	if err != nil {
		value = g.Compute(ordinal)
		if err := g.cache.Write(ordinal, value); err != nil {
			log.Errorf("Failed to write to cache")
		}
	}
	return value
}
