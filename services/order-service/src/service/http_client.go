package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// HTTPClient is a client for making HTTP requests with retry
type HTTPClient struct {
	client     *http.Client
	maxRetries int
}

// NewHTTPClient creates a new HTTP client with retry
func NewHTTPClient(timeout time.Duration, maxRetries int) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: timeout,
		},
		maxRetries: maxRetries,
	}
}

// Get makes a GET request with retry
func (c *HTTPClient) Get(url string, response interface{}) error {
	return c.doWithRetry("GET", url, nil, response)
}

// Post makes a POST request with retry
func (c *HTTPClient) Post(url string, body interface{}, response interface{}) error {
	return c.doWithRetry("POST", url, body, response)
}

// Put makes a PUT request with retry
func (c *HTTPClient) Put(url string, body interface{}, response interface{}) error {
	return c.doWithRetry("PUT", url, body, response)
}

// Delete makes a DELETE request with retry
func (c *HTTPClient) Delete(url string, response interface{}) error {
	return c.doWithRetry("DELETE", url, nil, response)
}

// doWithRetry makes an HTTP request with retry
func (c *HTTPClient) doWithRetry(method, url string, body interface{}, response interface{}) error {
	var bodyBytes []byte
	var err error

	// Marshal body if not nil
	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("error marshaling request body: %v", err)
		}
	}

	// Retry with exponential backoff
	var lastErr error
	for i := 0; i <= c.maxRetries; i++ {
		if i > 0 {
			// Exponential backoff: 1s, 2s, 4s, ...
			backoffTime := time.Duration(1<<uint(i-1)) * time.Second
			log.Printf("Retrying %s request to %s after %v...", method, url, backoffTime)
			time.Sleep(backoffTime)
		}

		// Create request
		var req *http.Request
		if body != nil {
			req, err = http.NewRequest(method, url, bytes.NewBuffer(bodyBytes))
		} else {
			req, err = http.NewRequest(method, url, nil)
		}
		if err != nil {
			lastErr = fmt.Errorf("error creating request: %v", err)
			continue
		}

		// Set headers
		req.Header.Set("Content-Type", "application/json")

		// Make request
		resp, err := c.client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("error making request: %v", err)
			continue
		}
		defer resp.Body.Close()

		// Read response body
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			lastErr = fmt.Errorf("error reading response body: %v", err)
			continue
		}

		// Check status code
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			lastErr = fmt.Errorf("error response: %s", string(respBody))
			continue
		}

		// Unmarshal response if not nil
		if response != nil {
			err = json.Unmarshal(respBody, response)
			if err != nil {
				lastErr = fmt.Errorf("error unmarshaling response: %v", err)
				continue
			}
		}

		// Success
		return nil
	}

	// All retries failed
	return fmt.Errorf("all retries failed: %v", lastErr)
}
