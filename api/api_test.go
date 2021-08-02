package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/programmablemike/fibo/internal/fibonacci"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockHttpClient struct {
	mock.Mock
}

func (m MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestCreateApiClient(t *testing.T) {
	assert := assert.New(t)
	client := new(MockHttpClient)
	body, _ := json.Marshal(GenericResponse{
		Status: StatusOK,
		Value:  "12345",
	})
	resp := &http.Response{
		Status:     fmt.Sprint(http.StatusOK),
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(body)),
	}
	client.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil)
	api := NewApiClient("http", "localhost", 80, client)
	v, err := api.Count(fibonacci.NewNumber(0))
	assert.NoError(err)
	assert.Equal("12345", v)
}

func TestCountApi(t *testing.T) {
	assert := assert.New(t)
	client := new(MockHttpClient)
	body, _ := json.Marshal(GenericResponse{
		Status: StatusOK,
		Value:  "12345",
	})
	resp := &http.Response{
		Status:     fmt.Sprint(http.StatusOK),
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(body)),
	}
	client.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil)
	api := NewApiClient("http", "localhost", 80, client)
	v, err := api.Count(fibonacci.NewNumber(120))
	assert.NoError(err)
	assert.Equal("12345", v)
}

func TestCalculateApi(t *testing.T) {
	assert := assert.New(t)
	client := new(MockHttpClient)
	body, _ := json.Marshal(GenericResponse{
		Status: StatusOK,
		Value:  "12345",
	})
	resp := &http.Response{
		Status:     fmt.Sprint(http.StatusOK),
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(body)),
	}
	client.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil)
	api := NewApiClient("http", "localhost", 80, client)
	v, err := api.Calculate(uint64(100))
	assert.NoError(err)
	assert.Equal("12345", v)
}

func TestClearCacheApi(t *testing.T) {
	assert := assert.New(t)
	client := new(MockHttpClient)
	body, _ := json.Marshal(GenericResponse{
		Status: StatusOK,
		Value:  "12345",
	})
	resp := &http.Response{
		Status:     fmt.Sprint(http.StatusOK),
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(body)),
	}
	client.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil)
	api := NewApiClient("http", "localhost", 80, client)
	err := api.ClearCache()
	assert.NoError(err)
}
