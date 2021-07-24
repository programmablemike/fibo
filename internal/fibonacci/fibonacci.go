// Calculates the Fibonacci sequence
package fibonacci

import (
	"math/big"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// Number is a convenience wrapper around big.Int
// It allows for arbitrarily long precision - without it we would overflow very quickly
type Number = big.Int

func NewNumber(v int64) *Number {
	return big.NewInt(v)
}

// NewNumberFromDecimalString converts the decimal value string into a big.Int
// This is useful for representing values that exceed standard integer length
func NewNumberFromDecimalString(v string) (*Number, bool) {
	result := NewNumber(0)
	return result.SetString(v, 10)
}

type Memoizer interface {
	Write(ordinal uint64, value *Number) error
	Read(ordinal uint64) (*Number, error)
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
func (g Generator) Compute(n uint64) *Number {
	log.Debugf("Computing fibonacci sequence for ordinal=%s", Uint64ToString(n))

	switch {
	case n == 0:
		return NewNumber(0)
	case n == 1:
		return NewNumber(1)
	case n > 1:
		n1 := g.readCachedOrCompute(n - 1)
		n2 := g.readCachedOrCompute(n - 2)
		res := NewNumber(0) // math.big requires a target to contain the result
		return res.Add(n1, n2)
	default:
		log.Errorf("Invalid condition reached for ordinal=%s", Uint64ToString(n))
		return NewNumber(-1)
	}
}

// readCachedOrCompute will read a value from the database if it exists
// otherwise it will compute the value and store it in the cache for future use
func (g *Generator) readCachedOrCompute(ordinal uint64) *Number {
	value, err := g.cache.Read(ordinal)
	if err != nil {
		value = g.Compute(ordinal)
		if err := g.cache.Write(ordinal, value); err != nil {
			log.Errorf("Failed to write to cache")
		}
	}
	return value
}
