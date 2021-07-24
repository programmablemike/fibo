// Calculates the Fibonacci sequence
// Defined as f(n) = f(n-1) + f(n-2) for n > 1
package fibonacci

import (
	"strconv"

	log "github.com/sirupsen/logrus"
)

type Memoizer interface {
	Write(ordinal uint64, value uint64) error
	Read(ordinal uint64) (uint64, error)
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

func Uint64ToString(v uint64) string {
	return strconv.FormatUint(v, 10)
}

// Compute Get the fibonacci value for the given ordinal
// Defined as f(n) = f(n-2) + f(n-1) where f(0) = 0 and f(1) = 1
// Returns -1 for invalid values
// Note that this does not detect integer overflows which can occur quickly
func (g Generator) Compute(n uint64) uint64 {
	log.Debugf("Computing fibonacci sequence for ordinal=%s", Uint64ToString(n))

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
		log.Errorf("Invalid condition reached for ordinal=%s", Uint64ToString(n))
		return 0
	}
}

// readCachedOrCompute will read a value from the database if it exists
// otherwise it will compute the value and store it in the cache for future use
func (g *Generator) readCachedOrCompute(ordinal uint64) uint64 {
	value, err := g.cache.Read(ordinal)
	if err != nil {
		value = g.Compute(ordinal)
		if err := g.cache.Write(ordinal, value); err != nil {
			log.Errorf("Failed to write to cache")
		}
	}
	return value
}
