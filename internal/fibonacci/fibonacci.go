// Calculates the Fibonacci sequence
// Defined as f(n) = f(n-1) + f(n-2) for n > 1
package fibonacci

import (
	"github.com/programmablemike/fibo/internal/cache"
	log "github.com/sirupsen/logrus"
)

type Generator struct {
	cache *cache.Cache
}

func (g *Generator) Close() error {
	if g.cache != nil {
		// Close the cache
		return g.cache.Close()
	}
	return nil
}

// Compute Get the fibonacci value for the given ordinal
// Defined as f(n) = f(n-2) + f(n-1) where f(0) = 0 and f(1) = 1
// Returns -1 for invalid values
func (g Generator) Compute(n int) int {
	log.Infof("Computing fibonacci sequence for ordinal %d...", n)

	switch {
	case n == 0:
		return 0
	case n == 1:
		return 1
	case n > 1:
		n1Idx := n - 1
		n2Idx := n - 2
		n1, n1err := g.cache.ReadEntry(n1Idx)
		n2, n2err := g.cache.ReadEntry(n2Idx)
		switch {
		case n1err == nil && n2err == nil:
			log.Debugf("Found cached value for %d", n1Idx)
			log.Debugf("Found cached value for %d", n2Idx)
			return n1.Result + n2.Result
		case n1err == nil && n2err != nil:
			log.Debugf("Found cached value for %d", n1Idx)
			log.Debugf("Missed cached value for %d", n2Idx)
			n2 := cache.CacheEntry{
				Ordinal: n2Idx,
				Result:  g.Compute(n2Idx),
			}
			g.cache.WriteEntry(&n2)
			return n1.Result + n2.Result
		case n1err != nil && n2err == nil:
			log.Debugf("Missed cached value for %d", n1Idx)
			log.Debugf("Found cached value for %d", n2Idx)
			n1 := cache.CacheEntry{
				Ordinal: n1Idx,
				Result:  g.Compute(n1Idx),
			}
			g.cache.WriteEntry(&n1)
			return n1.Result + n2.Result
		default:
			log.Debugf("Missed cached value for %d", n1Idx)
			log.Debugf("Missed cached value for %d", n2Idx)
			n1 := cache.CacheEntry{
				Ordinal: n1Idx,
				Result:  g.Compute(n1Idx),
			}
			n2 := cache.CacheEntry{
				Ordinal: n2Idx,
				Result:  g.Compute(n2Idx),
			}
			g.cache.WriteEntry(&n1)
			g.cache.WriteEntry(&n2)
			return n1.Result + n2.Result
		}
	default:
		return -1
	}
}
