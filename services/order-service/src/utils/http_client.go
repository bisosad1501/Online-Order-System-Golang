package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

// HTTPClient is a wrapper around http.Client with retry and timeout
type HTTPClient struct {
	client     *http.Client
	maxRetries int
	retryDelay time.Duration
	timeout    time.Duration
}

// NewHTTPClient creates a new HTTPClient with default settings
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		client:     &http.Client{},
		maxRetries: 2,
		retryDelay: 100 * time.Millisecond,
		timeout:    3 * time.Second,
	}
}

// NewHTTPClientWithOptions creates a new HTTPClient with custom settings
func NewHTTPClientWithOptions(maxRetries int, retryDelay time.Duration, timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		client:     &http.Client{},
		maxRetries: maxRetries,
		retryDelay: retryDelay,
		timeout:    timeout,
	}
}

// Get performs a GET request with retry and timeout
func (c *HTTPClient) Get(url string, result interface{}) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	return c.doRequest(req, result)
}

// Post performs a POST request with retry and timeout
func (c *HTTPClient) Post(url string, body interface{}, result interface{}) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	return c.doRequest(req, result)
}

// isTemporaryError kiểm tra xem lỗi có phải là lỗi tạm thời không
func isTemporaryError(err error, statusCode int) bool {
	// Lỗi mạng hoặc timeout
	if netErr, ok := err.(net.Error); ok {
		return netErr.Temporary() || netErr.Timeout()
	}

	// Lỗi timeout
	if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "deadline exceeded") {
		return true
	}

	// Lỗi connection refused/reset có thể là tạm thời
	if strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "connection reset") {
		return true
	}

	// Lỗi server (5xx)
	if statusCode >= 500 && statusCode < 600 {
		return true
	}

	return false
}

// doRequest performs the HTTP request with retry and timeout
func (c *HTTPClient) doRequest(req *http.Request, result interface{}) error {
	var lastErr error

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			// Wait before retrying with exponential backoff
			backoffTime := c.retryDelay * time.Duration(1<<attempt) // 100ms, 200ms, 400ms, ...
			log.Printf("Retrying request to %s after %v (attempt %d/%d)",
				req.URL.String(), backoffTime, attempt+1, c.maxRetries+1)
			time.Sleep(backoffTime)
			// Create a fresh request for retry to avoid reusing a closed body
			newReq := cloneRequest(req)
			req = newReq
		}

		// Create a context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
		reqWithCtx := req.WithContext(ctx)

		// Perform the request
		resp, err := c.client.Do(reqWithCtx)

		if err != nil {
			cancel() // Cancel the context to release resources
			lastErr = err

			// Kiểm tra xem có phải lỗi tạm thời không
			if isTemporaryError(err, 0) {
				log.Printf("Temporary error when calling %s: %v", req.URL.String(), err)
				continue // Retry on temporary errors
			} else {
				log.Printf("Permanent error when calling %s: %v", req.URL.String(), err)
				return fmt.Errorf("permanent error: %v", err)
			}
		}

		// Check if the status code indicates a server error (5xx)
		if resp.StatusCode >= 500 {
			resp.Body.Close()
			cancel()
			lastErr = fmt.Errorf("server error: %d", resp.StatusCode)
			log.Printf("Server error %d when calling %s", resp.StatusCode, req.URL.String())
			continue // Retry on server errors
		}

		// Check if the status code indicates a client error (4xx)
		if resp.StatusCode >= 400 {
			resp.Body.Close()
			cancel()
			log.Printf("Client error %d when calling %s", resp.StatusCode, req.URL.String())
			return fmt.Errorf("client error: %d", resp.StatusCode)
		}

		// If result is nil, we don't need to parse the response
		if result == nil {
			resp.Body.Close()
			cancel()
			return nil
		}

		// Parse the response
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		cancel() // Cancel the context to release resources

		if err != nil {
			lastErr = err
			log.Printf("Error reading response body from %s: %v", req.URL.String(), err)
			if isTemporaryError(err, 0) {
				continue // Retry on temporary errors
			} else {
				return fmt.Errorf("error reading response: %v", err)
			}
		}

		err = json.Unmarshal(body, result)
		if err != nil {
			lastErr = err
			log.Printf("Error unmarshaling JSON from %s: %v", req.URL.String(), err)
			return fmt.Errorf("error parsing response: %v", err)
		}

		return nil
	}

	return fmt.Errorf("max retries exceeded: %v", lastErr)
}

// cloneRequest creates a clone of the provided request
// This is needed because the original request body may be closed after first attempt
func cloneRequest(req *http.Request) *http.Request {
	newReq := &http.Request{
		Method:        req.Method,
		URL:           req.URL,
		Proto:         req.Proto,
		ProtoMajor:    req.ProtoMajor,
		ProtoMinor:    req.ProtoMinor,
		Header:        make(http.Header),
		Host:          req.Host,
		ContentLength: req.ContentLength,
	}

	// Copy headers
	for key, values := range req.Header {
		for _, value := range values {
			newReq.Header.Add(key, value)
		}
	}

	// Copy body if it exists
	if req.Body != nil && req.Body != http.NoBody {
		// For our tests, we know the body is a JSON payload
		// In a real implementation, you might want to handle this differently
		// or store the original body somewhere
		if req.GetBody != nil {
			newBody, err := req.GetBody()
			if err == nil {
				newReq.Body = newBody
				newReq.GetBody = req.GetBody
			}
		}
	}

	return newReq
}
