package client

import (
	"fmt"
	"io"
	"net/http"

	"github.com/opentracing/opentracing-go"
)

type Method int

const (
	Get Method = iota
	Head
	Post
	Put
	Delete
	Options
	Trace
	Patch
)

func (m Method) String() string {
	switch m {
	case Get:
		return "GET"
	case Head:
		return "HEAD"
	case Post:
		return "POST"
	case Put:
		return "PUT"
	case Delete:
		return "DELETE"
	case Options:
		return "OPTIONS"
	case Trace:
		return "TRACE"
	case Patch:
		return "PATCH"
	default:
		return "unknown"
	}
}

// HttpClient wraps an http.Client for request tracing
type HttpClient struct {
	http.Client
}

// Get is a convenience method for sending GET requests
func (hc HttpClient) Get(uri string, body io.Reader, parentSpan opentracing.Span) (*http.Response, error) {
	return hc.Execute(Get, uri, body, parentSpan)
}

// Delete is a convenience method for sending DELETE requests
func (hc HttpClient) Delete(uri string, body io.Reader, parentSpan opentracing.Span) (*http.Response, error) {
	return hc.Execute(Delete, uri, body, parentSpan)
}

// Execute wraps the request in a new tracing span and then calls the http.Client to send it
func (hc HttpClient) Execute(method Method, uri string, body io.Reader, parentSpan opentracing.Span) (res *http.Response, err error) {
	tracer := opentracing.GlobalTracer()

	// derive a new child span only for the request
	spanName := fmt.Sprintf("%s-%s", method, uri)
	span := opentracing.StartSpan(spanName, opentracing.ChildOf(parentSpan.Context()))
	defer span.Finish()

	client := &http.Client{}
	req, _ := http.NewRequest(method.String(), uri, body)

	// inject tracing headers to match up client requests with server responses
	tracer.Inject(span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header))

	res, err = client.Do(req)
	return res, err
}
