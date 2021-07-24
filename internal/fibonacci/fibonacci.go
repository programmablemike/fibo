// Calculates the Fibonacci sequence
// Defined as f(n) = f(n-1) + f(n-2) for n > 1
package fibonacci

import (
	log "github.com/sirupsen/logrus"
)

type Memoizer interface {
	Write(ordinal int, value int) error
	Read(ordinal int) (int, error)
	Clear() error
}

type Generator struct {
	cache Memoizer
}

// ClearCache wipes the memoizer's Postgres DB
func (g *Generator) ClearCache() error {
	return g.cache.Clear()
}

// Compute Get the fibonacci value for the given ordinal
// Defined as f(n) = f(n-2) + f(n-1) where f(0) = 0 and f(1) = 1
// Returns -1 for invalid values
func (g Generator) Compute(n int) int {
	log.Debugf("Computing fibonacci sequence for ordinal=%d", n)

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
func (g *Generator) readCachedOrCompute(ordinal int) int {
	value, err := g.cache.Read(ordinal)
	if err != nil {
		value = g.Compute(ordinal)
		if err := g.cache.Write(ordinal, value); err != nil {
			log.Errorf("Failed to write to cache")
		}
	}
	return value
}
