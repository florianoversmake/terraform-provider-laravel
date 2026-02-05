package forge_client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// jsonApiResource represents a single JSON:API resource item.
type jsonApiResource struct {
	ID         string          `json:"id"`
	Type       string          `json:"type"`
	Attributes json.RawMessage `json:"attributes"`
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
func (c *Client) GetJsonApiList(ctx context.Context, path string) ([]jsonApiResource, error) {
	resp, err := c.doRequestInternal(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	return extractJsonApiList(resp.Body)
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

// orgPath creates an organization-scoped URL path.
func (c *Client) orgPath(suffix string) string {
	return "/orgs/" + c.Organization + suffix
}
