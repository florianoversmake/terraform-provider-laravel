package envoyer_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const DefaultBaseURL = "https://envoyer.io/api"

type Client struct {
	httpClient      *http.Client
	baseURL         string
	EnvoyerAPIToken string
	envKey          string
}

func NewClient(EnvoyerAPIToken string, envKey string) *Client {
	return &Client{
		httpClient:      http.DefaultClient,
		baseURL:         DefaultBaseURL,
		EnvoyerAPIToken: EnvoyerAPIToken,
		envKey:          envKey,
	}
}

func (c *Client) WithBaseURL(baseURL string) *Client {
	c.baseURL = strings.TrimSuffix(baseURL, "/")
	return c
}

func (c *Client) WithEnvKey(envKey string) *Client {
	c.envKey = envKey
	return c
}

type ErrorResponse struct {
	Message string `json:"message,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

type ClientError struct {
	StatusCode int
	Body       string
}

func (e *ClientError) Error() string {
	return fmt.Sprintf("envoyer: status=%d, body=%s", e.StatusCode, e.Body)
}

func (c *Client) doRequest(ctx context.Context, method, path string, in, out any) error {
	// Build the URL.
	reqURL, err := url.Parse(c.baseURL + path)
	if err != nil {
		return fmt.Errorf("invalid URL '%s': %w", path, err)
	}

	var reqBody io.Reader
	if in != nil {
		buf, err := json.Marshal(in)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(buf)
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL.String(), reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set required headers.
	req.Header.Set("Authorization", "Bearer "+c.EnvoyerAPIToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Execute the request.
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	// Read response body.
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Check HTTP status code for error.
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return &ClientError{
			StatusCode: resp.StatusCode,
			Body:       string(bodyBytes),
		}
	}

	// If an output struct is provided, decode into it.
	if out != nil {
		if err := json.Unmarshal(bodyBytes, out); err != nil {
			return fmt.Errorf("failed to unmarshal response body: %w", err)
		}
	}

	return nil
}

func (c *Client) GetEnvKey() string {
	return c.envKey
}
