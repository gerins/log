package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gerins/log"
)

const (
	ProcessIDContextKey = "processID"
	ContentType         = "Content-Type"
	ApplicationJSON     = "application/json"
)

type rest struct {
	client            *http.Client
	addLogToExtraData bool
}

type Rest interface {
	Post(ctx context.Context, url string, header map[string]string, payload any) ([]byte, int, error)
	Put(ctx context.Context, url string, header map[string]string, payload any) ([]byte, int, error)
	Get(ctx context.Context, url string, header map[string]string, queryParams map[string]string) ([]byte, int, error)
	Delete(ctx context.Context, url string, header map[string]string, queryParams map[string]string) ([]byte, int, error)
}

func New(timeout time.Duration, addLogToExtraData bool) Rest {
	return &rest{
		addLogToExtraData: addLogToExtraData,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (r *rest) Post(ctx context.Context, url string, header map[string]string, payload any) ([]byte, int, error) {
	var (
		err          error
		httpResponse *http.Response
		trace        = log.NewTrace(http.MethodPost, url, header, payload, r.addLogToExtraData)
	)

	defer func() { trace.Save(ctx, httpResponse) }() // Logging the response

	// Convert payload to []byte type
	requestPayload, err := json.Marshal(payload)
	if err != nil {
		log.Context(ctx).Errorf("error marshaling payload, %v", err)
		return nil, 0, err
	}

	// Creating new request with context
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestPayload))
	if err != nil {
		log.Context(ctx).Errorf("error creating new request, %v", err)
		return nil, 0, err
	}

	if header == nil {
		header = make(map[string]string)
	}

	header[ContentType] = ApplicationJSON
	header[ProcessIDContextKey] = log.Context(ctx).ProcessID()

	// Adding header to the request
	for key, value := range header {
		req.Header.Set(key, value)
	}

	// Execute http request
	httpResponse, err = r.client.Do(req)
	if err != nil {
		log.Context(ctx).Errorf("error when making request, %v", err)
		return nil, 0, err
	}

	defer httpResponse.Body.Close()

	// Get the response body
	rawResponse, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		log.Context(ctx).Errorf("error reading response body, %v", err)
		return nil, httpResponse.StatusCode, err
	}

	trace.RawRespBody = rawResponse
	return rawResponse, httpResponse.StatusCode, nil
}

func (r *rest) Put(ctx context.Context, url string, header map[string]string, payload any) ([]byte, int, error) {
	var (
		err          error
		httpResponse *http.Response
		trace        = log.NewTrace(http.MethodPost, url, header, payload, r.addLogToExtraData)
	)

	defer func() { trace.Save(ctx, httpResponse) }() // Logging the response

	// Convert payload to []byte type
	requestPayload, err := json.Marshal(payload)
	if err != nil {
		log.Context(ctx).Errorf("error marshaling payload, %v", err)
		return nil, 0, err
	}

	// Creating new request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(requestPayload))
	if err != nil {
		log.Context(ctx).Errorf("error creating new request, %v", err)
		return nil, 0, err
	}

	if header == nil {
		header = make(map[string]string)
	}

	header[ContentType] = ApplicationJSON
	header[ProcessIDContextKey] = log.Context(ctx).ProcessID()

	// Adding header to the request
	for key, value := range header {
		req.Header.Set(key, value)
	}

	// Execute http request
	httpResponse, err = r.client.Do(req)
	if err != nil {
		log.Context(ctx).Errorf("error when making request, %v", err)
		return nil, 0, err
	}

	defer httpResponse.Body.Close()

	// Get the response body
	rawResponse, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		log.Context(ctx).Errorf("error reading response body, %v", err)
		return nil, httpResponse.StatusCode, err
	}

	trace.RawRespBody = rawResponse
	return rawResponse, httpResponse.StatusCode, nil
}

func (r *rest) Get(ctx context.Context, url string, header map[string]string, queryParams map[string]string) ([]byte, int, error) {
	var (
		err          error
		httpResponse *http.Response
		trace        = log.NewTrace(http.MethodPost, url, header, queryParams, r.addLogToExtraData)
	)

	defer func() { trace.Save(ctx, httpResponse) }() // Logging the response

	// Creating new request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Context(ctx).Errorf("error creating new request, %v", err)
		return nil, 0, err
	}

	if header == nil {
		header = make(map[string]string)
	}

	header[ProcessIDContextKey] = log.Context(ctx).ProcessID()

	// Adding header to the request
	for key, value := range header {
		req.Header.Set(key, value)
	}

	// Building query params
	query := req.URL.Query()
	for key, value := range queryParams {
		query.Add(key, value)
	}

	// Add query params to the url
	req.URL.RawQuery = query.Encode()

	// Execute http request
	httpResponse, err = r.client.Do(req)
	if err != nil {
		log.Context(ctx).Errorf("error when making request, %v", err)
		return nil, 0, err
	}

	defer httpResponse.Body.Close()

	// Get the response body
	rawResponse, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		log.Context(ctx).Errorf("error reading response body, %v", err)
		return nil, httpResponse.StatusCode, err
	}

	trace.RawRespBody = rawResponse
	return rawResponse, httpResponse.StatusCode, nil
}

func (r *rest) Delete(ctx context.Context, url string, header map[string]string, queryParams map[string]string) ([]byte, int, error) {
	var (
		err          error
		httpResponse *http.Response
		trace        = log.NewTrace(http.MethodPost, url, header, queryParams, r.addLogToExtraData)
	)

	defer func() { trace.Save(ctx, httpResponse) }() // Logging the response

	// Creating new request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		log.Context(ctx).Errorf("error creating new request, %v", err)
		return nil, 0, err
	}

	if header == nil {
		header = make(map[string]string)
	}

	header[ProcessIDContextKey] = log.Context(ctx).ProcessID()

	// Adding header to the request
	for key, value := range header {
		req.Header.Set(key, value)
	}

	// Building query params
	query := req.URL.Query()
	for key, value := range queryParams {
		query.Add(key, value)
	}

	// Add query params to the url
	req.URL.RawQuery = query.Encode()

	// Execute http request
	httpResponse, err = r.client.Do(req)
	if err != nil {
		log.Context(ctx).Errorf("error when making request, %v", err)
		return nil, 0, err
	}

	defer httpResponse.Body.Close()

	// Get the response body
	rawResponse, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		log.Context(ctx).Errorf("error reading response body, %v", err)
		return nil, httpResponse.StatusCode, err
	}

	trace.RawRespBody = rawResponse
	return rawResponse, httpResponse.StatusCode, nil
}
