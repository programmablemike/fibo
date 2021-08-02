package tracing

import (
	"io"

	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-lib/metrics"
)

// SetupTracing initializes a new OpenTracing tracer
// returns an io.Closer to be deferred in main() to flush the stream
//   and/or an error if the setup failed
func SetupTracing(serviceName string) io.Closer {
	cfg := jaegercfg.Configuration{
		ServiceName: serviceName,
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
		return nil
	}
	opentracing.SetGlobalTracer(tracer)

	return closer
}
