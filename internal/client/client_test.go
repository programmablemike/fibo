package client

import (
	"os"
	"testing"

	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-lib/metrics"
)

func TestMain(m *testing.M) {
	// Setup opentracing
	cfg := jaegercfg.Configuration{
		ServiceName: "client_test",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}

	jMetricsFactory := jaegerlog.NullFactory
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Metrics(jMetricsFactory),
	)
	if err != nil {
		// we could panic here instead if we want to fail fast
		log.Errorf("Failed to setup tracing")
	}
	opentracing.SetGlobalTracer(tracer)

	// Run the test
	retCode := m.Run()

	// manually close - defer doesn't play well with os.Exit()
	closer.Close()

	os.Exit(retCode)
}

func TestGet(t *testing.T) {
	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("TestGet")

}

func TestPost(t *testing.T) {

}

func TestExecute(t *testing.T) {

}
