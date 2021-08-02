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
	"github.com/spf13/viper"
)

func NewRouter(gen *fibonacci.Generator, tracer *opentracing.Tracer) *mux.Router {
	r := mux.NewRouter()

	// Root handler
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		span := tracing.StartSpanFromRequest("root-receive", opentracing.GlobalTracer(), r)
		defer span.Finish()
		w.WriteHeader(http.StatusOK)
		res := api.GenericResponse{
			Status:  api.StatusOK,
			Message: "OK",
		}
		json.NewEncoder(w).Encode(res)
	})

	// Ordinal handler
	r.HandleFunc("/fibo/calculate/{ordinal}", func(w http.ResponseWriter, r *http.Request) {
		span := tracing.StartSpanFromRequest("calculate-receive", opentracing.GlobalTracer(), r)
		defer span.Finish()

		vars := mux.Vars(r)

		log.Infof("Calculating Fibonacci number for ordinal=%s...", vars["ordinal"])
		ord, err := strconv.ParseUint(vars["ordinal"], 10, 64)
		if err != nil {
			res := api.GenericResponse{
				Status:  api.StatusError,
				Message: "failed to parse ordinal value",
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(res)
			return
		}
		value := gen.Compute(ord)
		res := api.GenericResponse{
			Status:  api.StatusOK,
			Message: "",
			Value:   value.String(),
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	}).Methods("GET")

	r.HandleFunc("/fibo/cache", func(w http.ResponseWriter, r *http.Request) {
		span := tracing.StartSpanFromRequest("clear-cache-receive", opentracing.GlobalTracer(), r)
		defer span.Finish()

		log.Info("Clearing the memoizer cache...")

		err := gen.ClearCache()
		if err != nil {
			res := api.GenericResponse{
				Status:  api.StatusError,
				Message: fmt.Sprintf("%e", err),
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(res)
			return
		}
		res := api.GenericResponse{
			Status:  api.StatusOK,
			Message: "Cache cleared",
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	}).Methods("DELETE")

	// Step counter
	r.HandleFunc("/fibo/count/{number}", func(w http.ResponseWriter, r *http.Request) {
		span := tracing.StartSpanFromRequest("count-receive", opentracing.GlobalTracer(), r)
		defer span.Finish()

		vars := mux.Vars(r)
		log.Infof("Counting ordinals between 0 and %s...", vars["number"])
		number, ok := fibonacci.NewNumberFromDecimalString(vars["number"])
		if !ok {
			res := api.GenericResponse{
				Status:  api.StatusError,
				Message: "failed to parse Fibonacci number value",
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(res)
			return
		}
		value := gen.FindOrdinalsInRange(fibonacci.NewNumber(0), number)
		res := api.GenericResponse{
			Status:  api.StatusOK,
			Message: "",
			Value:   fibonacci.Uint64ToString(value),
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	}).Methods("GET")

	return r
}

// createDsnFromConfig converts the options in the CLI flags/environment/.fiborc into a Postgres
// connection string
func createDsnFromConfig() string {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		viper.GetString("pguser"),
		viper.GetString("pgpassword"),
		viper.GetString("pghost"),
		viper.GetInt("pgport"),
		viper.GetString("pgdb"),
	)
	return dsn
}
