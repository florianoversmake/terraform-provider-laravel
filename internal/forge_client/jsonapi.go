package forge_client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// jsonApiResource represents a single JSON:API resource item.
type jsonApiResource struct {
	ID         string          `json:"id"`
	Type       string          `json:"type"`
	Attributes json.RawMessage `json:"attributes"`
	Meta       json.RawMessage `json:"meta,omitempty"`
}

// extractJsonApiSingle extracts the resource ID and raw attributes from a JSON:API single-resource response.
func extractJsonApiSingle(body []byte) (int, json.RawMessage, error) {
	var envelope struct {
		Data jsonApiResource `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return 0, nil, fmt.Errorf("failed to parse JSON:API response: %w", err)
	}
	id, _ := strconv.Atoi(envelope.Data.ID)
	return id, envelope.Data.Attributes, nil
}

// extractJsonApiSingleWithMeta extracts the resource ID, raw attributes, and resource-level meta
// from a JSON:API single-resource response. Some endpoints (e.g. server creation) return
// sensitive data like passwords in data.meta rather than data.attributes.
func extractJsonApiSingleWithMeta(body []byte) (int, json.RawMessage, json.RawMessage, error) {
	var envelope struct {
		Data jsonApiResource `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return 0, nil, nil, fmt.Errorf("failed to parse JSON:API response: %w", err)
	}
	id, _ := strconv.Atoi(envelope.Data.ID)
	return id, envelope.Data.Attributes, envelope.Data.Meta, nil
}

// extractJsonApiList extracts a list of resource items from a JSON:API list response.
func extractJsonApiList(body []byte) ([]jsonApiResource, error) {
	var envelope struct {
		Data []jsonApiResource `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, fmt.Errorf("failed to parse JSON:API list response: %w", err)
	}
	return envelope.Data, nil
}

// GetJsonApi performs a GET request and unwraps a single JSON:API resource.
// It returns the resource ID (from data.id) and unmarshals data.attributes into out.
func (c *Client) GetJsonApi(ctx context.Context, path string, out any) (int, error) {
	resp, err := c.doRequestInternal(ctx, http.MethodGet, path, nil)
	if err != nil {
		return 0, err
	}
	id, attrs, err := extractJsonApiSingle(resp.Body)
	if err != nil {
		return 0, err
	}
	if out != nil && len(attrs) > 0 {
		if err := json.Unmarshal(attrs, out); err != nil {
			return id, fmt.Errorf("failed to unmarshal JSON:API attributes: %w", err)
		}
	}
	return id, nil
}

// GetJsonApiWithoutCache performs a cached-bypassing GET and unwraps a single JSON:API resource.
func (c *Client) GetJsonApiWithoutCache(ctx context.Context, path string, out any) (int, error) {
	resp, err := c.doRequestInternal(ctx, http.MethodGet, path, nil, WithRequestCache(false))
	if err != nil {
		return 0, err
	}
	id, attrs, err := extractJsonApiSingle(resp.Body)
	if err != nil {
		return 0, err
	}
	if out != nil && len(attrs) > 0 {
		if err := json.Unmarshal(attrs, out); err != nil {
			return id, fmt.Errorf("failed to unmarshal JSON:API attributes: %w", err)
		}
	}
	return id, nil
}

// GetJsonApiResource performs a GET request and returns the raw JSON:API resource.
// This is useful with the unmarshalSingle generic helper.
func (c *Client) GetJsonApiResource(ctx context.Context, path string) (jsonApiResource, error) {
	resp, err := c.doRequestInternal(ctx, http.MethodGet, path, nil)
	if err != nil {
		return jsonApiResource{}, err
	}
	var envelope struct {
		Data jsonApiResource `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &envelope); err != nil {
		return jsonApiResource{}, fmt.Errorf("failed to parse JSON:API response: %w", err)
	}
	return envelope.Data, nil
}

// PostJsonApi performs a POST request and unwraps a single JSON:API resource response.
// Returns the resource ID and unmarshals data.attributes into out.
// If out is nil or the response body is empty, only the error is meaningful.
func (c *Client) PostJsonApi(ctx context.Context, path string, in, out any) (int, error) {
	resp, err := c.doRequestInternal(ctx, http.MethodPost, path, in)
	if err != nil {
		return 0, err
	}
	if out == nil || len(resp.Body) == 0 {
		return 0, nil
	}
	id, attrs, err := extractJsonApiSingle(resp.Body)
	if err != nil {
		return 0, err
	}
	if len(attrs) > 0 {
		if err := json.Unmarshal(attrs, out); err != nil {
			return id, fmt.Errorf("failed to unmarshal JSON:API attributes: %w", err)
		}
	}
	return id, nil
}

// PutJsonApi performs a PUT request and unwraps a single JSON:API resource response.
func (c *Client) PutJsonApi(ctx context.Context, path string, in, out any) (int, error) {
	resp, err := c.doRequestInternal(ctx, http.MethodPut, path, in)
	if err != nil {
		return 0, err
	}
	if out == nil || len(resp.Body) == 0 {
		return 0, nil
	}
	id, attrs, err := extractJsonApiSingle(resp.Body)
	if err != nil {
		return 0, err
	}
	if len(attrs) > 0 {
		if err := json.Unmarshal(attrs, out); err != nil {
			return id, fmt.Errorf("failed to unmarshal JSON:API attributes: %w", err)
		}
	}
	return id, nil
}

// GetJsonApiList performs a GET request and unwraps a JSON:API list response.
// Returns a slice of raw jsonApiResource items. The caller should unmarshal
// each item's Attributes field into the appropriate struct.
// NOTE: This only returns a single page. Prefer GetJsonApiListAll for complete results.
func (c *Client) GetJsonApiList(ctx context.Context, path string) ([]jsonApiResource, error) {
	resp, err := c.doRequestInternal(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	return extractJsonApiList(resp.Body)
}

// jsonApiLinks represents the links object in a JSON:API response.
// The API may return either an object {"first":"...","next":"..."} or
// an empty array [] when there are no pagination links. This custom
// type handles both cases gracefully.
type jsonApiLinks struct {
	First string `json:"first"`
	Last  string `json:"last"`
	Prev  string `json:"prev"`
	Next  string `json:"next"`
}

func (l *jsonApiLinks) UnmarshalJSON(data []byte) error {
	// If the API returns an empty array [] instead of an object, treat it as empty links.
	if len(data) > 0 && data[0] == '[' {
		return nil
	}
	// Otherwise unmarshal normally using an alias to avoid infinite recursion.
	type alias jsonApiLinks
	return json.Unmarshal(data, (*alias)(l))
}

// jsonApiMeta represents the meta object in a JSON:API paginated response.
// Same flexibility: the API may return an empty array [] when there is no
// pagination metadata.
type jsonApiMeta struct {
	Path       *string `json:"path"`
	PerPage    int     `json:"per_page"`
	NextCursor *string `json:"next_cursor"`
	PrevCursor *string `json:"prev_cursor"`
}

func (m *jsonApiMeta) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == '[' {
		return nil
	}
	type alias jsonApiMeta
	return json.Unmarshal(data, (*alias)(m))
}

// jsonApiListEnvelope represents the full JSON:API paginated list response,
// including data, links, and meta for cursor-based pagination.
type jsonApiListEnvelope struct {
	Data  []jsonApiResource `json:"data"`
	Links jsonApiLinks      `json:"links"`
	Meta  jsonApiMeta       `json:"meta"`
}

// GetJsonApiListAll performs paginated GET requests, automatically walking all
// pages using cursor-based pagination. It collects and returns all resources
// across all pages. The path should NOT include page[size] or page[cursor]
// parameters — they will be managed automatically.
func (c *Client) GetJsonApiListAll(ctx context.Context, path string) ([]jsonApiResource, error) {
	const pageSize = 100
	const maxPages = 100 // safety limit to prevent infinite loops

	// Ensure the path has the page[size] parameter set
	currentPath := appendPaginationParam(path, pageSize)

	var allItems []jsonApiResource
	var lastCursor string

	for page := 0; page < maxPages; page++ {
		resp, err := c.doRequestInternal(ctx, http.MethodGet, currentPath, nil)
		if err != nil {
			return nil, err
		}

		var envelope jsonApiListEnvelope
		if err := json.Unmarshal(resp.Body, &envelope); err != nil {
			return nil, fmt.Errorf("failed to parse JSON:API paginated response: %w", err)
		}

		allItems = append(allItems, envelope.Data...)

		// If there's no next cursor, we've reached the last page
		if envelope.Meta.NextCursor == nil || *envelope.Meta.NextCursor == "" {
			break
		}

		// Safety: if the cursor hasn't changed, the API is returning the same page — break to avoid infinite loop
		if *envelope.Meta.NextCursor == lastCursor {
			break
		}
		lastCursor = *envelope.Meta.NextCursor

		// Build the next page path using the cursor
		currentPath = appendCursorParam(path, pageSize, lastCursor)
	}

	return allItems, nil
}

// appendPaginationParam adds or replaces the page[size] query parameter on a path.
func appendPaginationParam(path string, pageSize int) string {
	sep := "?"
	if strings.Contains(path, "?") {
		sep = "&"
	}
	return fmt.Sprintf("%s%spage[size]=%d", path, sep, pageSize)
}

// appendCursorParam builds a path with both page[size] and page[cursor] query parameters.
// It strips any existing page[size]/page[cursor] params from the base path, then appends
// new ones with raw (unencoded) brackets, since the Forge API expects literal [] in param names.
func appendCursorParam(basePath string, pageSize int, cursor string) string {
	// Remove any existing page params from the query string
	clean := basePath
	for _, param := range []string{"page[size]", "page[cursor]"} {
		// Strip "param=value&" or "&param=value" or "?param=value"
		for {
			idx := strings.Index(clean, param+"=")
			if idx < 0 {
				break
			}
			// Find the end of this param=value pair
			end := strings.Index(clean[idx:], "&")
			if end < 0 {
				// Last param in the string
				if idx > 0 && clean[idx-1] == '&' {
					clean = clean[:idx-1]
				} else if idx > 0 && clean[idx-1] == '?' {
					clean = clean[:idx-1]
				} else {
					clean = clean[:idx]
				}
			} else {
				// Remove param and trailing &
				clean = clean[:idx] + clean[idx+end+1:]
			}
		}
	}

	sep := "?"
	if strings.Contains(clean, "?") {
		sep = "&"
	}
	return fmt.Sprintf("%s%spage[size]=%d&page[cursor]=%s", clean, sep, pageSize, cursor)
}

// unmarshalList is a helper that unmarshals a list of JSON:API resources into a typed slice.
// It also sets the ID on each item using the provided setter function.
func unmarshalList[T any](items []jsonApiResource, setID func(*T, int64)) ([]T, error) {
	result := make([]T, 0, len(items))
	for _, item := range items {
		var v T
		if err := json.Unmarshal(item.Attributes, &v); err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON:API item: %w", err)
		}
		id, _ := strconv.Atoi(item.ID)
		if setID != nil {
			setID(&v, int64(id))
		}
		result = append(result, v)
	}
	return result, nil
}

// unmarshalSingle is a helper that unmarshals a single JSON:API resource into a typed value.
// It also sets the ID using the provided setter function.
func unmarshalSingle[T any](item jsonApiResource, setID func(*T, int64)) (*T, error) {
	var v T
	if err := json.Unmarshal(item.Attributes, &v); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON:API item: %w", err)
	}
	id, _ := strconv.Atoi(item.ID)
	if setID != nil {
		setID(&v, int64(id))
	}
	return &v, nil
}

// orgPath creates an organization-scoped URL path.
func (c *Client) orgPath(suffix string) string {
	return "/orgs/" + c.Organization + suffix
}
