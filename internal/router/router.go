package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/programmablemike/fibo/internal/fibonacci"
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

		gen := fibonacci.Generator{}
		ord, err := strconv.ParseInt(vars["ordinal"], 10, 0)
		if err != nil {
			res := GenericResponse{
				Status:  StatusError,
				Message: fmt.Sprintf("%e", err),
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(res)
			return
		}
		value := gen.Compute(int(ord))
		res := GenericResponse{
			Status:  StatusOK,
			Message: "",
			Value:   fmt.Sprintf("%d", value),
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	}).Methods("GET")

	// Step counter
	r.HandleFunc("/fibo/{ordinal}/count", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		gen := fibonacci.Generator{}
		ord, err := strconv.ParseInt(vars["ordinal"], 10, 0)
		if err != nil {
			res := GenericResponse{
				Status:  StatusError,
				Message: fmt.Sprintf("%e", err),
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(res)
			return
		}
		value := gen.Compute(int(ord))
		res := GenericResponse{
			Status:  StatusOK,
			Message: "",
			Value:   fmt.Sprintf("%d", value),
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	}).Methods("GET")

	return r
}
