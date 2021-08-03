package tracing

import (
	"testing"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
)

func TestSetupTracer(t *testing.T) {
	closer := Init("TestSetupTracer")
	// defer closing the tracer
	defer func() {
		err := closer.Close()
		assert.NoError(t, err)
	}()
	time.Sleep(1 * time.Second)
}

func TestLogSpans(t *testing.T) {
	closer := Init("TestLogSpans")
	// defer closing the tracer while checking for any errors
	defer func() {
		err := closer.Close()
		assert.NoError(t, err)
	}()
	// setup the tracer
	tracer := opentracing.GlobalTracer()

	span := tracer.StartSpan("test-span")
	defer span.Finish()
	// This should happen in the context of the parent span
	time.Sleep(1 * time.Second)
	childSpan := tracer.StartSpan("test-child-span", opentracing.ChildOf(span.Context()))
	defer childSpan.Finish()
	// This should happen in the context of the child and parent span
	time.Sleep(1 * time.Second)
}
