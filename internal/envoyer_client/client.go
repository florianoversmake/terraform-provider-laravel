package envoyer_client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const DefaultBaseURL = "https://envoyer.io/api"

// ResponseFormat represents the expected format of an API response.
type ResponseFormat string

const (
	// ResponseFormatJSON indicates the response should be treated as JSON (default).
	ResponseFormatJSON ResponseFormat = "json"
	// ResponseFormatText indicates the response should be treated as plain text.
	ResponseFormatText ResponseFormat = "text"
	// ResponseFormatRaw indicates the response should be returned as raw bytes.
	ResponseFormatRaw ResponseFormat = "raw"
)

// Response represents an API response with various formats supported.
type Response struct {
	StatusCode int         // HTTP status code
	Headers    http.Header // Response headers
	Body       []byte      // Raw response body
}

// JSON unmarshals the response body as JSON into the provided value.
func (r *Response) JSON(v interface{}) error {
	if r.Body == nil || len(r.Body) == 0 {
		return fmt.Errorf("empty response body")
	}
	return json.Unmarshal(r.Body, v)
}

// Text returns the response body as a string.
func (r *Response) Text() string {
	if r.Body == nil {
		return ""
	}
	return string(r.Body)
}

// Raw returns the raw response body.
func (r *Response) Raw() []byte {
	return r.Body
}

// IsSuccess returns true if the response status code is in the 2xx range.
func (r *Response) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// RetryConfig holds configuration for request retries.
type RetryConfig struct {
	// MaxRetries is the maximum number of times to retry a failed request.
	MaxRetries int
	// BaseDelay is the initial delay before the first retry.
	BaseDelay time.Duration
	// MaxDelay is the maximum delay between retries.
	MaxDelay time.Duration
	// RetryableStatusCodes is a list of HTTP status codes that should trigger a retry.
	RetryableStatusCodes []int
	// RetryableErrors is a list of error types that should trigger a retry.
	RetryableErrors []error
}

// DefaultRetryConfig returns a default retry configuration.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:           3,
		BaseDelay:            500 * time.Millisecond,
		MaxDelay:             30 * time.Second,
		RetryableStatusCodes: []int{408, 429, 500, 502, 503, 504},
	}
}

// RequestOption is a functional option for configuring individual requests.
type RequestOption func(*requestOptions)

// requestOptions holds configuration for a single request.
type requestOptions struct {
	responseFormat ResponseFormat    // Format of the expected response
	headers        map[string]string // Additional headers for the request
	retry          *RetryConfig      // Custom retry options for this request
	debug          bool              // Enable debug for this request only
	forceNoRetry   bool              // Force no retries for this request
}

// WithResponseFormat sets the expected response format for a request.
func WithResponseFormat(format ResponseFormat) RequestOption {
	return func(opts *requestOptions) {
		opts.responseFormat = format
	}
}

// WithHeader adds a custom header to the request.
func WithHeader(key, value string) RequestOption {
	return func(opts *requestOptions) {
		if opts.headers == nil {
			opts.headers = make(map[string]string)
		}
		opts.headers[key] = value
	}
}

// WithRequestRetry sets custom retry options for a specific request.
func WithRequestRetry(config RetryConfig) RequestOption {
	return func(opts *requestOptions) {
		opts.retry = &config
	}
}

// WithDebug enables or disables debug logging for a specific request.
func WithDebug(debug bool) RequestOption {
	return func(opts *requestOptions) {
		opts.debug = debug
	}
}

// WithNoRetry disables retries for a specific request.
func WithNoRetry() RequestOption {
	return func(opts *requestOptions) {
		opts.forceNoRetry = true
	}
}

// ErrorResponse represents an error response from the API.
type ErrorResponse struct {
	Message string `json:"message,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

// ClientError represents a client-side error.
type ClientError struct {
	StatusCode  int
	Body        string
	Method      string
	URL         string
	RequestBody string
	Headers     http.Header
}

func (e *ClientError) Error() string {
	errorMsg := fmt.Sprintf("envoyer: status=%d, method=%s, url=%s", e.StatusCode, e.Method, e.URL)
	if e.Body != "" {
		errorMsg += fmt.Sprintf(", body=%s", e.Body)
	}
	if e.RequestBody != "" {
		errorMsg += fmt.Sprintf(", request=%s", e.RequestBody)
	}
	return errorMsg
}

// Client is the Envoyer API client.
type Client struct {
	httpClient      *http.Client
	baseURL         string
	EnvoyerAPIToken string
	envKey          string
	debug           bool
	retryConfig     RetryConfig
}

// NewClient creates a new Envoyer API client.
func NewClient(EnvoyerAPIToken string, envKey string) *Client {
	return &Client{
		httpClient:      http.DefaultClient,
		baseURL:         DefaultBaseURL,
		EnvoyerAPIToken: EnvoyerAPIToken,
		envKey:          envKey,
		debug:           false,
		retryConfig:     DefaultRetryConfig(),
	}
}

// WithBaseURL sets a custom base URL for the API.
func (c *Client) WithBaseURL(baseURL string) *Client {
	c.baseURL = strings.TrimSuffix(baseURL, "/")
	return c
}

// WithEnvKey sets the environment key.
func (c *Client) WithEnvKey(envKey string) *Client {
	c.envKey = envKey
	return c
}

// WithDebug enables or disables debug logging.
func (c *Client) WithDebug(debug bool) *Client {
	c.debug = debug
	return c
}

// WithHTTPClient sets a custom HTTP client.
func (c *Client) WithHTTPClient(httpClient *http.Client) *Client {
	c.httpClient = httpClient
	return c
}

// WithRetryConfig sets the retry configuration.
func (c *Client) WithRetryConfig(config RetryConfig) *Client {
	c.retryConfig = config
	return c
}

// getEffectiveRequestOptions merges default options with per-request options.
func (c *Client) getEffectiveRequestOptions(opts ...RequestOption) requestOptions {
	// Start with default options
	effective := requestOptions{
		responseFormat: ResponseFormatJSON, // Default to JSON
		debug:          c.debug,            // Use client's debug setting as default
	}

	// Apply all request-specific options
	for _, opt := range opts {
		opt(&effective)
	}

	return effective
}

// shouldRetry determines if a request should be retried based on the error and status code.
func (c *Client) shouldRetry(err error, statusCode int, retryConfig RetryConfig, attempt int) bool {
	// Check if we've exceeded max retries
	if attempt >= retryConfig.MaxRetries {
		return false
	}

	// Check if this is a retryable status code
	for _, code := range retryConfig.RetryableStatusCodes {
		if statusCode == code {
			return true
		}
	}

	// Check if this is a retryable error type
	if err != nil {
		// Network errors are generally retryable
		if strings.Contains(err.Error(), "connection refused") ||
			strings.Contains(err.Error(), "no such host") ||
			strings.Contains(err.Error(), "i/o timeout") ||
			strings.Contains(err.Error(), "connection reset") ||
			strings.Contains(err.Error(), "EOF") {
			return true
		}

		// Check custom retryable errors
		for _, retryableErr := range retryConfig.RetryableErrors {
			if strings.Contains(err.Error(), retryableErr.Error()) {
				return true
			}
		}
	}

	return false
}

// calculateBackoff calculates the backoff duration for retries.
func (c *Client) calculateBackoff(attempt int, retryConfig RetryConfig) time.Duration {
	// Exponential backoff with jitter
	delay := float64(retryConfig.BaseDelay) * math.Pow(1.5, float64(attempt))

	// Add jitter (10% randomness)
	jitter := 0.1 * delay

	// Convert the jitter to int64 for the modulo operation
	jitterInt := int64(jitter)
	if jitterInt <= 0 {
		jitterInt = 1 // Ensure we don't have a zero modulus
	}

	// Calculate random jitter value and convert back to float64
	randomJitter := float64(time.Now().UnixNano() % jitterInt)

	// Scale the random jitter to be in the range [-jitter/2, jitter/2]
	scaledJitter := (randomJitter/float64(jitterInt))*jitter - (jitter / 2)

	// Apply the jitter to the delay
	delay = delay + scaledJitter

	// Cap at max delay
	if delay > float64(retryConfig.MaxDelay) {
		delay = float64(retryConfig.MaxDelay)
	}

	return time.Duration(delay)
}

// debugLog prints debug information if debug is enabled.
func (c *Client) debugLog(reqOpts requestOptions, format string, args ...interface{}) {
	if reqOpts.debug {
		fmt.Printf("[Envoyer Debug] "+format+"\n", args...)
	}
}

// doRequestInternal performs an HTTP request with the given options and returns a Response object.
func (c *Client) doRequestInternal(ctx context.Context, method, path string, in any, opts ...RequestOption) (*Response, error) {
	reqOpts := c.getEffectiveRequestOptions(opts...)

	// Use client's retry config as default
	retryConfig := c.retryConfig
	if reqOpts.retry != nil {
		retryConfig = *reqOpts.retry
	}

	// Check if retries are disabled for this request
	if reqOpts.forceNoRetry {
		retryConfig.MaxRetries = 0
	}

	// Prepare the request body once
	var reqBody io.Reader
	var reqBodyStr string
	if in != nil {
		buf, err := json.Marshal(in)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBodyStr = string(buf)
		if reqOpts.debug {
			c.debugLog(reqOpts, "Request Body: %s", reqBodyStr)
		}
	}

	var attempt int
	for {
		// Reset the request body for each attempt
		if reqBodyStr != "" {
			reqBody = strings.NewReader(reqBodyStr)
		}

		reqURL, err := url.Parse(c.baseURL + path)
		if err != nil {
			return nil, fmt.Errorf("invalid URL '%s': %w", path, err)
		}

		req, err := http.NewRequestWithContext(ctx, method, reqURL.String(), reqBody)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Set standard headers
		req.Header.Set("Authorization", "Bearer "+c.EnvoyerAPIToken)
		req.Header.Set("Content-Type", "application/json")

		// Set appropriate accept header based on response format
		switch reqOpts.responseFormat {
		case ResponseFormatJSON:
			req.Header.Set("Accept", "application/json")
		case ResponseFormatText:
			req.Header.Set("Accept", "text/plain")
		case ResponseFormatRaw:
			req.Header.Set("Accept", "*/*")
		}

		// Add any custom headers
		for k, v := range reqOpts.headers {
			req.Header.Set(k, v)
		}

		if reqOpts.debug {
			c.debugLog(reqOpts, "Making %s request to %s", method, reqURL.String())
			c.debugLog(reqOpts, "Headers: %v", req.Header)
		}

		// Execute the request
		resp, err := c.httpClient.Do(req)
		if err != nil {
			// Check if context was canceled
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
				// Check if we should retry
				if c.shouldRetry(err, 0, retryConfig, attempt) {
					delay := c.calculateBackoff(attempt, retryConfig)
					attempt++
					c.debugLog(reqOpts, "Request failed with error: %v. Retrying in %v (attempt %d/%d)",
						err, delay, attempt, retryConfig.MaxRetries)

					select {
					case <-time.After(delay):
						continue
					case <-ctx.Done():
						return nil, ctx.Err()
					}
				}
				return nil, fmt.Errorf("request error: %w", err)
			}
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		resp.Body.Close() // Close immediately to avoid resource leak
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		if reqOpts.debug {
			c.debugLog(reqOpts, "Response Status: %d", resp.StatusCode)
			c.debugLog(reqOpts, "Response Body: %s", string(bodyBytes))
		}

		// Create response object
		response := &Response{
			StatusCode: resp.StatusCode,
			Headers:    resp.Header,
			Body:       bodyBytes,
		}

		// Check if we should retry based on status code
		if c.shouldRetry(nil, resp.StatusCode, retryConfig, attempt) {
			delay := c.calculateBackoff(attempt, retryConfig)
			attempt++
			c.debugLog(reqOpts, "Request failed with status code: %d. Retrying in %v (attempt %d/%d)",
				resp.StatusCode, delay, attempt, retryConfig.MaxRetries)

			select {
			case <-time.After(delay):
				continue
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		// Handle non-success status codes
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			clientErr := &ClientError{
				StatusCode:  resp.StatusCode,
				Body:        string(bodyBytes),
				Method:      method,
				URL:         reqURL.String(),
				RequestBody: reqBodyStr,
				Headers:     resp.Header,
			}
			return response, clientErr
		}

		return response, nil
	}
}

// doRequest is the original method for backward compatibility.
func (c *Client) doRequest(ctx context.Context, method, path string, in, out any) error {
	// Create a request option that respects the client's debug setting
	opts := []RequestOption{WithDebug(c.debug), WithResponseFormat(ResponseFormatJSON)}

	resp, err := c.doRequestInternal(ctx, method, path, in, opts...)
	if err != nil {
		return err
	}

	// If out is provided and we have response data, unmarshal it
	if out != nil && len(resp.Body) > 0 {
		if err := json.Unmarshal(resp.Body, out); err != nil {
			return fmt.Errorf("failed to unmarshal response body: %w (body: %s)", err, string(resp.Body))
		}
	}

	return nil
}

// Get performs a GET request with backward compatibility for the original API.
func (c *Client) Get(ctx context.Context, path string, out any) error {
	return c.doRequest(ctx, http.MethodGet, path, nil, out)
}

// GetWithResponse is the new enhanced version that returns a Response object.
func (c *Client) GetWithResponse(ctx context.Context, path string, opts ...RequestOption) (*Response, error) {
	return c.doRequestInternal(ctx, http.MethodGet, path, nil, opts...)
}

// GetText performs a GET request and returns the response as plain text.
func (c *Client) GetText(ctx context.Context, path string, opts ...RequestOption) (string, error) {
	opts = append(opts, WithResponseFormat(ResponseFormatText))
	resp, err := c.doRequestInternal(ctx, http.MethodGet, path, nil, opts...)
	if err != nil {
		return "", err
	}
	return resp.Text(), nil
}

// GetRaw performs a GET request and returns the raw response bytes.
func (c *Client) GetRaw(ctx context.Context, path string, opts ...RequestOption) ([]byte, error) {
	opts = append(opts, WithResponseFormat(ResponseFormatRaw))
	resp, err := c.doRequestInternal(ctx, http.MethodGet, path, nil, opts...)
	if err != nil {
		return nil, err
	}
	return resp.Raw(), nil
}

// Post performs a POST request with backward compatibility for the original API.
func (c *Client) Post(ctx context.Context, path string, in, out any) error {
	return c.doRequest(ctx, http.MethodPost, path, in, out)
}

// PostWithResponse is the new enhanced version that returns a Response object.
func (c *Client) PostWithResponse(ctx context.Context, path string, in any, opts ...RequestOption) (*Response, error) {
	return c.doRequestInternal(ctx, http.MethodPost, path, in, opts...)
}

// PostText performs a POST request and returns the response as plain text.
func (c *Client) PostText(ctx context.Context, path string, in any, opts ...RequestOption) (string, error) {
	opts = append(opts, WithResponseFormat(ResponseFormatText))
	resp, err := c.doRequestInternal(ctx, http.MethodPost, path, in, opts...)
	if err != nil {
		return "", err
	}
	return resp.Text(), nil
}

// Put performs a PUT request with backward compatibility for the original API.
func (c *Client) Put(ctx context.Context, path string, in, out any) error {
	return c.doRequest(ctx, http.MethodPut, path, in, out)
}

// PutWithResponse is the new enhanced version that returns a Response object.
func (c *Client) PutWithResponse(ctx context.Context, path string, in any, opts ...RequestOption) (*Response, error) {
	return c.doRequestInternal(ctx, http.MethodPut, path, in, opts...)
}

// PutText performs a PUT request and returns the response as plain text.
func (c *Client) PutText(ctx context.Context, path string, in any, opts ...RequestOption) (string, error) {
	opts = append(opts, WithResponseFormat(ResponseFormatText))
	resp, err := c.doRequestInternal(ctx, http.MethodPut, path, in, opts...)
	if err != nil {
		return "", err
	}
	return resp.Text(), nil
}

// Delete performs a DELETE request with backward compatibility for the original API.
func (c *Client) Delete(ctx context.Context, path string, out any) error {
	return c.doRequest(ctx, http.MethodDelete, path, nil, out)
}

// DeleteWithResponse is the new enhanced version that returns a Response object.
func (c *Client) DeleteWithResponse(ctx context.Context, path string, opts ...RequestOption) (*Response, error) {
	return c.doRequestInternal(ctx, http.MethodDelete, path, nil, opts...)
}

// DeleteText performs a DELETE request and returns the response as plain text.
func (c *Client) DeleteText(ctx context.Context, path string, opts ...RequestOption) (string, error) {
	opts = append(opts, WithResponseFormat(ResponseFormatText))
	resp, err := c.doRequestInternal(ctx, http.MethodDelete, path, nil, opts...)
	if err != nil {
		return "", err
	}
	return resp.Text(), nil
}

// GetEnvKey returns the environment key.
func (c *Client) GetEnvKey() string {
	return c.envKey
}
