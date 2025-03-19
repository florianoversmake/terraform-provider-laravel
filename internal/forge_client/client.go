package forge_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

const DefaultBaseURL = "https://forge.laravel.com/api/v1"

// CacheItem represents a cached HTTP response with metadata.
type CacheItem struct {
	Value      []byte      // The cached response body
	Expiration time.Time   // When this cache entry expires
	StatusCode int         // The HTTP status code of the response
	Headers    http.Header // Response headers
}

// Cache is the interface for the caching system.
type Cache interface {
	Get(key string) (*CacheItem, bool)
	Set(key string, item *CacheItem)
	Delete(key string)
	Clear()
	Keys() []string
}

// MemoryCache implements an in-memory cache with expiration.
type MemoryCache struct {
	items map[string]*CacheItem
	mu    sync.RWMutex
}

// NewMemoryCache creates a new in-memory cache.
func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		items: make(map[string]*CacheItem),
	}
}

// Get retrieves an item from the cache if it exists and hasn't expired.
func (c *MemoryCache) Get(key string) (*CacheItem, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	// Check if the item has expired
	if time.Now().After(item.Expiration) {
		// Item has expired, but we'll let cleanup handle deletion
		return nil, false
	}

	return item, true
}

// Set adds an item to the cache.
func (c *MemoryCache) Set(key string, item *CacheItem) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = item
}

// Delete removes an item from the cache.
func (c *MemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Clear removes all items from the cache.
func (c *MemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]*CacheItem)
}

// Keys returns all keys in the cache.
func (c *MemoryCache) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]string, 0, len(c.items))
	for k := range c.items {
		keys = append(keys, k)
	}
	return keys
}

// Cleanup removes expired items from the cache.
func (c *MemoryCache) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for k, v := range c.items {
		if now.After(v.Expiration) {
			delete(c.items, k)
		}
	}
}

// CacheConfig holds caching configuration.
type CacheConfig struct {
	Enabled             bool
	TTL                 time.Duration
	CleanupInterval     time.Duration
	MaxCacheSize        int  // Maximum number of items in cache (0 for unlimited)
	CacheErrorResponses bool // Whether to cache error responses
}

// RequestOption is a functional option for configuring individual requests.
type RequestOption func(*requestOptions)

// requestOptions holds configuration for a single request.
type requestOptions struct {
	cacheEnabled        *bool          // Whether caching is enabled for this request
	cacheTTL            *time.Duration // TTL for this specific request
	cacheErrorResponses *bool          // Whether to cache error responses for this request
	forceRefresh        bool           // Force a refresh (bypass cache)
}

// WithRequestCache enables or disables caching for a specific request.
func WithRequestCache(enabled bool) RequestOption {
	return func(opts *requestOptions) {
		opts.cacheEnabled = &enabled
	}
}

// WithRequestCacheTTL sets a specific TTL for a request.
func WithRequestCacheTTL(ttl time.Duration) RequestOption {
	return func(opts *requestOptions) {
		opts.cacheTTL = &ttl
	}
}

// WithRequestCacheErrorResponses enables or disables caching of error responses for a specific request.
func WithRequestCacheErrorResponses(enabled bool) RequestOption {
	return func(opts *requestOptions) {
		opts.cacheErrorResponses = &enabled
	}
}

// WithForceRefresh forces a refresh from the API, bypassing the cache.
func WithForceRefresh() RequestOption {
	return func(opts *requestOptions) {
		opts.forceRefresh = true
	}
}

// Client is the Forge API client.
type Client struct {
	httpClient    *http.Client
	baseURL       string
	ForgeAPIToken string

	// Configurable retry settings
	MaxRetries int           // Maximum number of retries after receiving a 429
	RetryDelay time.Duration // Default delay between retries

	// Caching configuration
	cache       Cache
	cacheConfig CacheConfig

	// Cleanup ticker for cache maintenance
	cleanupTicker *time.Ticker
	cleanupDone   chan bool
}

// NewClient creates a new Forge API client.
func NewClient(ForgeAPIToken string) *Client {
	client := &Client{
		httpClient:    http.DefaultClient,
		baseURL:       DefaultBaseURL,
		ForgeAPIToken: ForgeAPIToken,
		// Set default retry values
		MaxRetries: 6,
		RetryDelay: 10 * time.Second,
		// Default cache configuration (disabled by default)
		cacheConfig: CacheConfig{
			Enabled:             false,
			TTL:                 5 * time.Minute,
			CleanupInterval:     10 * time.Minute,
			CacheErrorResponses: false,
		},
		cleanupDone: make(chan bool),
	}

	return client
}

// WithCache enables caching with the provided cache implementation.
func (c *Client) WithCache(cache Cache) *Client {
	c.cache = cache

	// Enable caching if a cache is provided
	if cache != nil {
		c.cacheConfig.Enabled = true

		// Start the cleanup routine if needed
		if c.cleanupTicker == nil && c.cacheConfig.CleanupInterval > 0 {
			c.startCleanupRoutine()
		}
	}

	return c
}

// WithCacheConfig configures the caching behavior.
func (c *Client) WithCacheConfig(config CacheConfig) *Client {
	c.cacheConfig = config

	// Start or stop the cleanup routine based on configuration
	if c.cacheConfig.Enabled && c.cache != nil {
		if c.cleanupTicker != nil {
			c.cleanupTicker.Stop()
			c.cleanupDone <- true
		}

		if c.cacheConfig.CleanupInterval > 0 {
			c.startCleanupRoutine()
		}
	} else if c.cleanupTicker != nil {
		c.cleanupTicker.Stop()
		c.cleanupDone <- true
		c.cleanupTicker = nil
	}

	return c
}

// startCleanupRoutine starts a goroutine that periodically cleans up expired cache entries.
func (c *Client) startCleanupRoutine() {
	c.cleanupTicker = time.NewTicker(c.cacheConfig.CleanupInterval)

	go func() {
		for {
			select {
			case <-c.cleanupTicker.C:
				if memCache, ok := c.cache.(*MemoryCache); ok {
					memCache.Cleanup()
				}
			case <-c.cleanupDone:
				return
			}
		}
	}()
}

// WithBaseURL sets a custom base URL for the API.
func (c *Client) WithBaseURL(baseURL string) *Client {
	c.baseURL = strings.TrimSuffix(baseURL, "/")
	return c
}

// DisableCache disables caching.
func (c *Client) DisableCache() *Client {
	c.cacheConfig.Enabled = false
	return c
}

// EnableCache enables caching.
func (c *Client) EnableCache() *Client {
	// Create a memory cache if one doesn't exist
	if c.cache == nil {
		c.cache = NewMemoryCache()
	}
	c.cacheConfig.Enabled = true

	// Start cleanup routine if needed
	if c.cleanupTicker == nil && c.cacheConfig.CleanupInterval > 0 {
		c.startCleanupRoutine()
	}

	return c
}

// ClearCache clears all cached items.
func (c *Client) ClearCache() {
	if c.cache != nil {
		c.cache.Clear()
	}
}

// Close cleans up any resources and stops background routines.
func (c *Client) Close() {
	if c.cleanupTicker != nil {
		c.cleanupTicker.Stop()
		c.cleanupDone <- true
		c.cleanupTicker = nil
	}
}

// CacheStats returns statistics about the cache.
type CacheStats struct {
	ItemCount int
	Keys      []string
}

// GetCacheStats returns statistics about the current cache state.
func (c *Client) GetCacheStats() CacheStats {
	if c.cache == nil {
		return CacheStats{
			ItemCount: 0,
			Keys:      []string{},
		}
	}

	keys := c.cache.Keys()
	return CacheStats{
		ItemCount: len(keys),
		Keys:      keys,
	}
}

// InvalidateCacheKey removes a specific key from the cache.
func (c *Client) InvalidateCacheKey(key string) {
	if c.cache != nil {
		c.cache.Delete(key)
	}
}

// ErrorResponse represents an error response from the API.
type ErrorResponse struct {
	Message string `json:"message,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

// ClientError represents a client-side error.
type ClientError struct {
	StatusCode int
	Body       string
}

func (e *ClientError) Error() string {
	return fmt.Sprintf("forge: status=%d, body=%s", e.StatusCode, e.Body)
}

// ClientErrorResourceNotFound represents a 404 error.
type ClientErrorResourceNotFound struct {
	StatusCode int
	Body       string
}

func (e *ClientErrorResourceNotFound) Error() string {
	return fmt.Sprintf("forge: status=%d, body=%s", e.StatusCode, e.Body)
}

// generateCacheKey creates a unique key for caching based on the request details.
func (c *Client) generateCacheKey(method, path string, reqBody []byte) string {
	if reqBody != nil {
		return fmt.Sprintf("%s:%s:%x", method, path, reqBody)
	}
	return fmt.Sprintf("%s:%s", method, path)
}

// getEffectiveRequestOptions merges default options with per-request options.
func (c *Client) getEffectiveRequestOptions(opts ...RequestOption) requestOptions {
	// Start with default options
	effective := requestOptions{}

	// Apply all request-specific options
	for _, opt := range opts {
		opt(&effective)
	}

	return effective
}

// doRequestWithOptions performs an HTTP request with caching if enabled and respecting per-request options.
func (c *Client) doRequestWithOptions(ctx context.Context, method, path string, in, out any, opts ...RequestOption) error {
	// Process request-specific options
	reqOpts := c.getEffectiveRequestOptions(opts...)

	// Determine if caching is enabled for this request
	isCacheable := c.cacheConfig.Enabled && c.cache != nil && (method == http.MethodGet)
	if reqOpts.cacheEnabled != nil {
		isCacheable = *reqOpts.cacheEnabled && c.cache != nil && (method == http.MethodGet)
	}

	// Determine if error responses should be cached
	cacheErrors := c.cacheConfig.CacheErrorResponses
	if reqOpts.cacheErrorResponses != nil {
		cacheErrors = *reqOpts.cacheErrorResponses
	}

	// Get TTL for this request
	ttl := c.cacheConfig.TTL
	if reqOpts.cacheTTL != nil {
		ttl = *reqOpts.cacheTTL
	}

	// Prepare the request body once
	var reqBodyBytes []byte
	if in != nil {
		var err error
		reqBodyBytes, err = json.Marshal(in)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	// Generate cache key for this request
	cacheKey := c.generateCacheKey(method, path, reqBodyBytes)

	// Try to get the response from cache first (unless force refresh is requested)
	if isCacheable && !reqOpts.forceRefresh {
		if cachedItem, found := c.cache.Get(cacheKey); found {
			// We have a valid cached response
			if out != nil {
				if err := json.Unmarshal(cachedItem.Value, out); err != nil {
					return fmt.Errorf("failed to unmarshal cached response: %w", err)
				}
			}
			return nil
		}
	}

	var attempt int
	for {
		reqURL, err := url.Parse(c.baseURL + path)
		if err != nil {
			return fmt.Errorf("invalid URL '%s': %w", path, err)
		}

		// Use a fresh reader for each attempt
		var reqBody io.Reader
		if reqBodyBytes != nil {
			reqBody = bytes.NewReader(reqBodyBytes)
		}

		req, err := http.NewRequestWithContext(ctx, method, reqURL.String(), reqBody)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		// Set required headers
		req.Header.Set("Authorization", "Bearer "+c.ForgeAPIToken)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")

		// Execute the request
		resp, err := c.httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("request error: %w", err)
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		resp.Body.Close() // Close immediately to avoid resource leak
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		// Handle 429 (Too Many Requests) with automatic throttling
		if resp.StatusCode == http.StatusTooManyRequests {
			if attempt < c.MaxRetries {
				// Check for a Retry-After header value
				delay := c.RetryDelay
				if ra := resp.Header.Get("Retry-After"); ra != "" {
					if seconds, err := strconv.Atoi(ra); err == nil {
						delay = time.Duration(seconds) * time.Second
					}
				}
				// Wait for the delay or until the context is cancelled
				select {
				case <-time.After(delay):
					// Retry the request
					attempt++
					continue
				case <-ctx.Done():
					return ctx.Err()
				}
			} else {
				return &ClientError{
					StatusCode: resp.StatusCode,
					Body:       string(bodyBytes),
				}
			}
		}

		// Handle other non-success status codes
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			clientErr := &ClientError{
				StatusCode: resp.StatusCode,
				Body:       string(bodyBytes),
			}

			if resp.StatusCode == http.StatusNotFound {
				clientErr := &ClientErrorResourceNotFound{
					StatusCode: resp.StatusCode,
					Body:       string(bodyBytes),
				}

				// Cache error responses if configured to do so
				if isCacheable && cacheErrors {
					c.cache.Set(cacheKey, &CacheItem{
						Value:      bodyBytes,
						Expiration: time.Now().Add(ttl),
						StatusCode: resp.StatusCode,
						Headers:    resp.Header,
					})
				}

				return clientErr
			}

			// Cache error responses if configured to do so
			if isCacheable && cacheErrors {
				c.cache.Set(cacheKey, &CacheItem{
					Value:      bodyBytes,
					Expiration: time.Now().Add(ttl),
					StatusCode: resp.StatusCode,
					Headers:    resp.Header,
				})
			}

			return clientErr
		}

		// On success, store in cache if applicable
		if isCacheable {
			c.cache.Set(cacheKey, &CacheItem{
				Value:      bodyBytes,
				Expiration: time.Now().Add(ttl),
				StatusCode: resp.StatusCode,
				Headers:    resp.Header,
			})
		}

		// Decode the response into out (if provided)
		if out != nil {
			if err := json.Unmarshal(bodyBytes, out); err != nil {
				return fmt.Errorf("failed to unmarshal response body: %w", err)
			}
		}

		return nil
	}
}

// doRequest is the legacy method for backward compatibility.
func (c *Client) doRequest(ctx context.Context, method, path string, in, out any) error {
	return c.doRequestWithOptions(ctx, method, path, in, out)
}

// InvalidateCache removes all entries from the cache that match a given prefix.
func (c *Client) InvalidateByPrefix(prefix string) int {
	if c.cache == nil {
		return 0
	}

	count := 0
	for _, key := range c.cache.Keys() {
		if strings.HasPrefix(key, prefix) {
			c.cache.Delete(key)
			count++
		}
	}

	return count
}

// Get performs a GET request with caching.
func (c *Client) Get(ctx context.Context, path string, out any, opts ...RequestOption) error {
	return c.doRequestWithOptions(ctx, http.MethodGet, path, nil, out, opts...)
}

// Post performs a POST request.
func (c *Client) Post(ctx context.Context, path string, in, out any, opts ...RequestOption) error {
	return c.doRequestWithOptions(ctx, http.MethodPost, path, in, out, opts...)
}

// Put performs a PUT request.
func (c *Client) Put(ctx context.Context, path string, in, out any, opts ...RequestOption) error {
	return c.doRequestWithOptions(ctx, http.MethodPut, path, in, out, opts...)
}

// Delete performs a DELETE request.
func (c *Client) Delete(ctx context.Context, path string, out any, opts ...RequestOption) error {
	return c.doRequestWithOptions(ctx, http.MethodDelete, path, nil, out, opts...)
}

// WithHTTPClient sets a custom HTTP client.
func (c *Client) WithHTTPClient(httpClient *http.Client) *Client {
	c.httpClient = httpClient
	return c
}

// WithRetryConfig configures the retry behavior.
func (c *Client) WithRetryConfig(maxRetries int, retryDelay time.Duration) *Client {
	c.MaxRetries = maxRetries
	c.RetryDelay = retryDelay
	return c
}

func (c *Client) HTTPClient() *http.Client {
	return c.httpClient
}
