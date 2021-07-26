package workers

import (
	"github.com/programmablemike/fibo/internal/fibonacci"
)

type GeneratorPool struct {
	cache   fibonacci.Memoizer     // A shared cache for checking memoized values
	workers int                    // The number of workers to spawn
	input   chan uint64            // The job input queue
	output  chan *fibonacci.Number // The job output queue
	quit    chan bool              // Send the termination signal
	done    chan bool              // Result of the termination
	running bool                   // Whether or not the pool is running
}

func (gp GeneratorPool) IsRunning() bool {
	return gp.running
}

func (gp *GeneratorPool) GetInputOutput() (in chan<- uint64, out <-chan *fibonacci.Number) {
	return gp.input, gp.output
}

func (gp *GeneratorPool) Start() {
	for i := 0; i < gp.workers; i++ {
		// Spawn the goroutines for the pool, passing in the input and output queues
		go computeInGoroutine(gp.cache, gp.input, gp.output, gp.done)
	}
	gp.running = true
}

func (gp *GeneratorPool) ClearCache() error {
	wasRunning := gp.IsRunning()
	// Stop the workers if running
	if wasRunning {
		gp.Close()
	}
	gp.cache.Clear()
	// If the pool was previously running, start the workers again
	if wasRunning {
		gp.Start()
	}
	return nil
}

func (gp *GeneratorPool) Close() error {
	for i := 0; i < gp.workers; i++ {
		// Push a termination signal into the done channel for each worker to consume
		gp.quit <- true
		// For every quit signal, verify we received a done signal
		<-gp.done
	}
	gp.running = false
	return nil
}

// CreatenewGeneratorPool of size `workers`
// All workers reuse the same posgres database connection
// Returns an input channel, an output channel, and a "done" channel for signalling termination
func NewGeneratorPool(cache fibonacci.Memoizer, workers int) *GeneratorPool {
	// Input and output are oversized on purpose to allow for buffering
	input := make(chan uint64, 4*workers)
	output := make(chan *fibonacci.Number, 4*workers)
	done := make(chan bool, workers) // Must exactly match the number of workers
	gp := GeneratorPool{
		cache:   cache,
		workers: workers,
		input:   input,
		output:  output,
		done:    done,
		running: false,
	}
	return &gp
}

// computeInGoroutine wraps the generator in a Goroutine for concurrent execution
// The done channel should be sized to exactly match the number of goroutines
// For best performance, the in/out channels should be sized larger than the number of goroutines
func computeInGoroutine(c fibonacci.Memoizer, in <-chan uint64, out chan<- *fibonacci.Number, quit <-chan bool, done chan<- bool) {
	// Initialize the generator
	g := fibonacci.NewGenerator(c)
	for {
		select {
		case ord := <-in:
			// Calculate the Fibonacci number and push it in the result channel
			out <- g.Compute(ord)
		case <-quit:
			// Send confirmation of termination
			done <- true
			return
		}
	}
}
