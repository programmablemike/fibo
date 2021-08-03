// Calculates the Fibonacci sequence
package fibonacci

import (
	"context"
	"math/big"
	"strconv"

	"github.com/programmablemike/fibo/internal/tracing"
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
	SetContext(ctx context.Context)
}

func Uint64ToString(v uint64) string {
	return strconv.FormatUint(v, 10)
}

func StringToUint64(v string) (uint64, error) {
	value, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}

type Generator struct {
	cache   Memoizer
	context context.Context
}

// NewGenerator creates a new Fibonacci generator
// ctx is added for implementing tracing support to get more specific information around performance
func NewGenerator(ctx context.Context, cache Memoizer) *Generator {
	// TODO: Convert this to using a goroutine for concurrent execution
	return &Generator{
		context: ctx,
		cache:   cache,
	}
}

// ClearCache wipes the memoizer's Postgres DB
func (g *Generator) ClearCache() error {
	span := tracing.StartSpanFromContext(g.context, "clear-cache")
	defer span.Finish()
	return g.cache.Clear()
}

func (g *Generator) FindOrdinalsInRange(low *Number, high *Number) uint64 {
	span := tracing.StartSpanFromContext(g.context, "find-ordinals-in-range")
	defer span.Finish()

	log.Debugf("Counting ordinals in range %s to %s...", low.String(), high.String())

	// Initialize the first three fibonacci values
	f0 := NewNumber(0)
	f1 := NewNumber(1)
	f2 := NewNumber(1)
	// Contains the count of valid ordinals
	count := uint64(0)

	// Values for x.Cmp(y):
	//   -1 if x <  y
	//    0 if x == y
	//   +1 if x >  y
	for f0.Cmp(high) == -1 || f0.Cmp(high) == 0 {
		if f0.Cmp(low) == 1 || f0.Cmp(low) == 0 {
			count += 1 // Valid value, increment the count
		}
		f0 = NewNumber(0).Set(f1)
		f1 = NewNumber(0).Set(f2)
		f2 = NewNumber(0).Add(f0, f1)
	}
	return count
}

// Compute Get the fibonacci value for the given ordinal
// Defined as f(n) = f(n-2) + f(n-1) where f(0) = 0 and f(1) = 1
// Returns -1 for invalid values
func (g *Generator) Compute(n uint64) *Number {
	span := tracing.StartSpanFromContext(g.context, "compute")
	defer span.Finish()
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
	span := tracing.StartSpanFromContext(g.context, "read-cached-or-compute")
	defer span.Finish()
	g.cache.SetContext(tracing.SaveParentSpan(context.Background(), span))
	value, err := g.cache.Read(ordinal)
	if err != nil {
		value = g.Compute(ordinal)
		g.cache.SetContext(tracing.SaveParentSpan(context.Background(), span))
		if err := g.cache.Write(ordinal, value); err != nil {
			log.Errorf("Failed to write to cache")
		}
		span.SetTag("cache_hit", false)
		return value
	} else {
		span.SetTag("cache_hit", true)
		return value
	}
}
