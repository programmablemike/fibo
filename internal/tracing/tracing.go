package tracing

import (
	"fmt"
	"io"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-lib/metrics"
)

// SetupTracing initializes a new OpenTracing tracer
// returns an io.Closer to be deferred in main() to flush the stream
//   and/or an error if the setup failed
func Init(service string) io.Closer {
	defaultCfg := jaegercfg.Configuration{
		ServiceName: service,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LocalAgentHostPort: fmt.Sprintf("%s:%d", viper.GetString("jaegerHost"), viper.GetInt("jaegerPort")),
			LogSpans:           true,
		},
	}
	cfg, _ := defaultCfg.FromEnv()

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
	log.Info("Successfully initialized Jaeger tracing")

	return closer
}

func StartSpanFromRequest(span string, tracer opentracing.Tracer, r *http.Request) opentracing.Span {
	spanCtx, _ := Extract(tracer, r)
	return tracer.StartSpan(span, ext.RPCServerOption(spanCtx))
}

// Inject the HTTP headers for tracing into an HTTP request's headers
func Inject(span opentracing.Span, request *http.Request) error {
	return span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(request.Header),
	)
}

// Inject the HTTP headers for tracing into an HTTP response's headers
func InjectResponse(span opentracing.Span, response *http.Response) error {
	return span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(response.Header),
	)
}

// Extract the HTTP headers for tracing from an HTTP request's headers
func Extract(tracer opentracing.Tracer, r *http.Request) (opentracing.SpanContext, error) {
	return tracer.Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header),
	)
}
