// The main HTTP router for the fibo server
package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	api "github.com/programmablemike/fibo/api"
	"github.com/programmablemike/fibo/internal/fibonacci"
	"github.com/programmablemike/fibo/internal/tracing"
	log "github.com/sirupsen/logrus"
)

// EncodeJsonResponse encodes the provided body to a JSON response
// Additionally it adds an opentracing Span to time how long encoding is taking
func EncodeJsonResponse(w http.ResponseWriter, req *http.Request, body interface{}) {
	s := tracing.StartSpanFromContext(req.Context(), "encode-json-response")
	defer s.Finish()
	json.NewEncoder(w).Encode(body)
}

// MiddlewareExtractTracer extracts any Jaeger tracers from the HTTP headers and adds
// them to the request context
func MiddlewareExtractTracer(span string, next http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		log.Debug("running MiddlewareExtractTracer")
		span := tracing.StartSpanFromRequest(span, opentracing.GlobalTracer(), req)
		defer span.Finish()
		// Save the span info to the http.Request context for retrieving later
		req = req.WithContext(tracing.SaveParentSpan(req.Context(), span))
		next(res, req)
	}
}

// MiddlewareAddHttpTraceTags adds in the default HTTP tag values to the tracer
// This *must* be used with MiddlewareExtractTracer
func MiddlewareAddHttpTraceTags(next http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		log.Debug("running MiddlewareAddHttpTraceTags")
		span := tracing.GetParentSpan(req.Context())
		span.SetTag("http.method", req.Method)
		span.SetTag("http.url", req.URL.String())
		next(res, req)
	}
}

func NewRouter(gen *fibonacci.Generator) *mux.Router {
	r := mux.NewRouter()

	// Root handler
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		res := api.GenericResponse{
			Status:  api.StatusOK,
			Message: "OK",
		}
		EncodeJsonResponse(w, r, res)
	})

	// Ordinal handler
	r.HandleFunc("/fibo/calculate/{ordinal}",
		MiddlewareExtractTracer("calculate-request",
			MiddlewareAddHttpTraceTags(
				func(w http.ResponseWriter, r *http.Request) {
					vars := mux.Vars(r)

					log.Infof("Calculating Fibonacci number for ordinal=%s...", vars["ordinal"])
					ord, err := strconv.ParseUint(vars["ordinal"], 10, 64)
					if err != nil {
						res := api.GenericResponse{
							Status:  api.StatusError,
							Message: "failed to parse ordinal value",
						}
						w.WriteHeader(http.StatusBadRequest)
						EncodeJsonResponse(w, r, res)
						return
					}
					value := gen.Compute(ord)
					res := api.GenericResponse{
						Status:  api.StatusOK,
						Message: "",
						Value:   value.String(),
					}
					w.WriteHeader(http.StatusOK)
					EncodeJsonResponse(w, r, res)
				}))).Methods("GET")

	r.HandleFunc("/fibo/cache",
		MiddlewareExtractTracer("clear-cache-request",
			MiddlewareAddHttpTraceTags(
				func(w http.ResponseWriter, r *http.Request) {
					log.Info("Clearing the memoizer cache...")

					err := gen.ClearCache()
					if err != nil {
						res := api.GenericResponse{
							Status:  api.StatusError,
							Message: fmt.Sprintf("%e", err),
						}
						w.WriteHeader(http.StatusInternalServerError)
						EncodeJsonResponse(w, r, res)
						return
					}
					res := api.GenericResponse{
						Status:  api.StatusOK,
						Message: "Cache cleared",
					}
					w.WriteHeader(http.StatusOK)
					EncodeJsonResponse(w, r, res)
				}))).Methods("DELETE")

	// Step counter
	r.HandleFunc("/fibo/count/{number}",
		MiddlewareExtractTracer("count-request",
			MiddlewareAddHttpTraceTags(
				func(w http.ResponseWriter, r *http.Request) {
					vars := mux.Vars(r)
					log.Infof("Counting ordinals between 0 and %s...", vars["number"])
					number, ok := fibonacci.NewNumberFromDecimalString(vars["number"])
					if !ok {
						res := api.GenericResponse{
							Status:  api.StatusError,
							Message: "failed to parse Fibonacci number value",
						}
						w.WriteHeader(http.StatusBadRequest)
						EncodeJsonResponse(w, r, res)
						return
					}
					value := gen.FindOrdinalsInRange(fibonacci.NewNumber(0), number)
					res := api.GenericResponse{
						Status:  api.StatusOK,
						Message: "",
						Value:   fibonacci.Uint64ToString(value),
					}
					w.WriteHeader(http.StatusOK)
					EncodeJsonResponse(w, r, res)
				}))).Methods("GET")

	return r
}
