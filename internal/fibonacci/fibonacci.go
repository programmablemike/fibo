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

func Uint64ToString(v uint64) string {
	return strconv.FormatUint(v, 10)
}

type Generator struct {
	cache Memoizer
}

func NewGenerator(cache Memoizer) *Generator {
	// TODO: Convert this to using a goroutine for concurrent execution
	return &Generator{
		cache: cache,
	}
}

// ClearCache wipes the memoizer's Postgres DB
func (g *Generator) ClearCache() error {
	return g.cache.Clear()
}

// FindOrdinalInRange calculates the number of Fibonacci numbers in a given range
/*
func (g *Generator) FindOrdinalsInRange(low int, high int) int {
	f0 := 0
	f1 := 1
	f2 := 1
	count := 0
	for f0 <= high {
		if f0 >= low {
			count += 1
		}
		f0 = f1
		f1 = f2
		f2 = f0 + f1
	}
	return count
}
*/

func (g *Generator) FindOrdinalsInRange(low *Number, high *Number) uint64 {
	// Initialize the first three fibonacci values
	f0 := NewNumber(0)
	f1 := NewNumber(1)
	f2 := NewNumber(1)
	// Contains the count of valid ordinals
	count := uint64(0)

	log.Infof("Counting ordinals in range %s to %s...", low.String(), high.String())
	// Values for Number.Cmp:
	//   -1 if x <  y
	//    0 if x == y
	//   +1 if x >  y
	for f0.Cmp(high) == -1 || f0.Cmp(high) == 0 {
		if f0.Cmp(low) == 1 || f0.Cmp(low) == 0 {
			log.Infof("Found %s", f0.String())
			count += 1 // Valid value, increment the count
		}
		log.Infof("old f0: %s", f0.String())
		log.Infof("old f1: %s", f1.String())
		log.Infof("old f2: %s", f2.String())
		f0 = NewNumber(0).Set(f1)
		f1 = NewNumber(0).Set(f2)
		f2 = NewNumber(0).Add(f0, f1)
		log.Infof("new f0: %s", f0.String())
		log.Infof("new f1: %s", f1.String())
		log.Infof("new f2: %s", f2.String())
	}
	return count
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
