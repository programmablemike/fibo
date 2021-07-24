// The main HTTP router for the fibo server
package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/programmablemike/fibo/internal/cache"
	"github.com/programmablemike/fibo/internal/fibonacci"
	"github.com/spf13/viper"
)

type GenericResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Value   string `json:"value"`
}

const (
	StatusOK    string = "OK"
	StatusError string = "ERROR"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	// Root handler
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		res := GenericResponse{
			Status:  StatusOK,
			Message: "OK",
		}
		json.NewEncoder(w).Encode(res)
	})

	// Ordinal handler
	r.HandleFunc("/fibo/{ordinal}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		dsn := createDsnFromConfig()
		gen := fibonacci.NewGenerator(cache.NewCache(dsn))
		ord, err := strconv.ParseUint(vars["ordinal"], 10, 64)
		if err != nil {
			res := GenericResponse{
				Status:  StatusError,
				Message: "failed to parse ordinal value",
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(res)
			return
		}
		value := gen.Compute(uint64(ord))
		res := GenericResponse{
			Status:  StatusOK,
			Message: "",
			Value:   value.String(),
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	}).Methods("GET")

	r.HandleFunc("/fibo/cache", func(w http.ResponseWriter, r *http.Request) {
		dsn := createDsnFromConfig()
		gen := fibonacci.NewGenerator(cache.NewCache(dsn))
		err := gen.ClearCache()
		if err != nil {
			res := GenericResponse{
				Status:  StatusError,
				Message: fmt.Sprintf("%e", err),
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(res)
			return
		}
		res := GenericResponse{
			Status:  StatusOK,
			Message: "Cache cleared",
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	}).Methods("DELETE")

	// Step counter
	r.HandleFunc("/fibo/{ordinal}/count", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		dsn := createDsnFromConfig()
		gen := fibonacci.NewGenerator(cache.NewCache(dsn))
		ord, err := strconv.ParseUint(vars["ordinal"], 10, 64)
		if err != nil {
			res := GenericResponse{
				Status:  StatusError,
				Message: "failed to parse ordinal value",
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(res)
			return
		}
		value := gen.Compute(uint64(ord))
		res := GenericResponse{
			Status:  StatusOK,
			Message: "",
			Value:   value.String(),
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
