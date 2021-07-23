// Calculates the Fibonacci sequence
// Defined as f(n) = f(n-1) + f(n-2) for n > 1
package fibonacci

import (
	log "github.com/sirupsen/logrus"
)

type Generator struct {
}

func NewGenerator() *Generator {
	return &Generator{}
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
		// TODO Check the cache here for a pre-computed result
		return g.Compute(n-2) + g.Compute(n-1)
	default:
		return -1
	}
}
