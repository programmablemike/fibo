package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/programmablemike/fibo/internal/fibonacci"
	log "github.com/sirupsen/logrus"
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

// HttpClient defines a client interface we can mock out for testing
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Defines the Fibo Client's interface
type FiboClient interface {
	Count(max *fibonacci.Number) (uint64, error)
	Calculate(ordinal uint64) (*fibonacci.Number, error)
	ClearCache() error
}

type ApiClient struct {
	client   HttpClient
	protocol string // "http" or "https"
	host     string // the IP or hostname
	port     int    // the port number
}

// NewApiClient creates a new API client for
// set port == 0 to use the protocol's default port
func NewApiClient(protocol string, host string, port int, client HttpClient) *ApiClient {
	if protocol != "http" && protocol != "https" {
		log.Fatalf("invalid protocol %s; must be http or https", protocol)
		return nil
	}
	if port < 0 && port >= 65535 {
		log.Fatalf("invalid port number %d", port)
		return nil
	}

	return &ApiClient{
		client:   client,
		protocol: protocol,
		host:     host,
		port:     port,
	}
}

func (ac ApiClient) GetApiBaseUri() string {
	if ac.port == 0 { // if port == 0, use default for protocol
		return fmt.Sprintf("%s://%s/fibo", ac.protocol, ac.host)
	} else {
		return fmt.Sprintf("%s://%s:%d/fibo", ac.protocol, ac.host, ac.port)
	}
}

// execute makes the API call and handles any error that might occur
func (ac ApiClient) execute(req *http.Request) *http.Response {
	res, err := ac.client.Do(req)
	// catastrophic error - this can't be recovered from
	if err != nil {
		log.Fatalf("unrecoverable error during API request: %s", err)
		return nil
	}
	return res
}

func (ac ApiClient) DecodeGenericResponse(res *http.Response) *GenericResponse {
	resp := &GenericResponse{}
	err := json.NewDecoder(res.Body).Decode(resp)
	if err != nil {
		log.Errorf("failed to decode response as JSON: %s", err)
		return nil
	}
	return resp
}

func (ac ApiClient) Count(max *fibonacci.Number) (string, error) {
	uri := fmt.Sprintf("%s/count/%s", ac.GetApiBaseUri(), max.String())

	req, _ := http.NewRequest("GET", uri, nil)
	res := ac.execute(req)
	defer res.Body.Close()
	resp := ac.DecodeGenericResponse(res)

	if resp.Status == StatusOK {
		return resp.Value, nil
	} else {
		return "", fmt.Errorf(resp.Message)
	}
}

func (ac ApiClient) Calculate(ordinal uint64) (string, error) {
	uri := fmt.Sprintf("%s/calculate/%s", ac.GetApiBaseUri(), fibonacci.Uint64ToString(ordinal))

	req, _ := http.NewRequest("GET", uri, nil)
	res := ac.execute(req)
	defer res.Body.Close()
	resp := ac.DecodeGenericResponse(res)
	if resp.Status == StatusOK {
		return resp.Value, nil
	} else {
		return "", fmt.Errorf(resp.Message)
	}
}

func (ac ApiClient) ClearCache() error {
	uri := fmt.Sprintf("%s/cache", ac.GetApiBaseUri())

	req, _ := http.NewRequest("DELETE", uri, nil)
	res := ac.execute(req)
	defer res.Body.Close()
	resp := ac.DecodeGenericResponse(res)

	if resp.Status == StatusOK {
		return nil
	} else {
		return fmt.Errorf(resp.Message)
	}
}
