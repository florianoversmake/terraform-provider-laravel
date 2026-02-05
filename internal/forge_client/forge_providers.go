package forge_client

import (
	"context"
	"fmt"
)

// ProviderInfo represents a cloud provider from the JSON:API /providers endpoint.
type ProviderInfo struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	Slug              string `json:"slug"`
	SimpleName        string `json:"simple_name"`
	Currency          string `json:"currency"`
	CurrencySymbol    string `json:"currency_symbol"`
	DefaultSizeCode   string `json:"default_size_code"`
	DefaultRegionCode string `json:"default_region_code"`
}

// ProviderRegionInfo represents a region from the JSON:API /providers/{id}/regions endpoint.
type ProviderRegionInfo struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Code          string `json:"code"`
	AlternateCode string `json:"alternate_code"`
}

// ProviderSizeInfo represents a size from the JSON:API /providers/{id}/sizes endpoint.
type ProviderSizeInfo struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Code         string `json:"code"`
	Series       string `json:"series"`
	Category     string `json:"category"`
	Cpus         int    `json:"cpus"`
	DiskType     string `json:"disk_type"`
	Architecture string `json:"architecture"`
	Ram          int    `json:"ram"`
	Disk         int    `json:"disk"`
}

// ProviderRegionSizeInfo represents a size available in a region from the JSON:API
// /providers/{id}/regions/{id}/sizes endpoint.
type ProviderRegionSizeInfo struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// ListProviders returns all cloud providers.
// Note: This endpoint does not require an organization scope.
func (c *Client) ListProviders(ctx context.Context) ([]ProviderInfo, error) {
	items, err := c.GetJsonApiList(ctx, "/providers?page[size]=100")
	if err != nil {
		return nil, fmt.Errorf("failed to list providers: %w", err)
	}
	return unmarshalList(items, func(p *ProviderInfo, id int64) { p.ID = id })
}

// GetProvider returns a single cloud provider by ID.
// Note: This endpoint does not require an organization scope.
func (c *Client) GetProvider(ctx context.Context, providerID int64) (*ProviderInfo, error) {
	path := fmt.Sprintf("/providers/%d", providerID)
	resp, err := c.GetJsonApiResource(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}
	return unmarshalSingle(resp, func(p *ProviderInfo, id int64) { p.ID = id })
}

// ListProviderRegions returns all regions for a specific provider.
// Note: This endpoint does not require an organization scope.
func (c *Client) ListProviderRegions(ctx context.Context, providerID int64) ([]ProviderRegionInfo, error) {
	path := fmt.Sprintf("/providers/%d/regions?page[size]=100", providerID)
	items, err := c.GetJsonApiList(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to list provider regions: %w", err)
	}
	return unmarshalList(items, func(r *ProviderRegionInfo, id int64) { r.ID = id })
}

// GetProviderRegion returns a specific region for a provider.
// Note: This endpoint does not require an organization scope.
func (c *Client) GetProviderRegion(ctx context.Context, providerID, regionID int64) (*ProviderRegionInfo, error) {
	path := fmt.Sprintf("/providers/%d/regions/%d", providerID, regionID)
	resp, err := c.GetJsonApiResource(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider region: %w", err)
	}
	return unmarshalSingle(resp, func(r *ProviderRegionInfo, id int64) { r.ID = id })
}

// ListProviderSizes returns all sizes for a specific provider.
// Note: This endpoint does not require an organization scope.
func (c *Client) ListProviderSizes(ctx context.Context, providerID int64) ([]ProviderSizeInfo, error) {
	path := fmt.Sprintf("/providers/%d/sizes?page[size]=100", providerID)
	items, err := c.GetJsonApiList(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to list provider sizes: %w", err)
	}
	return unmarshalList(items, func(s *ProviderSizeInfo, id int64) { s.ID = id })
}

// GetProviderSize returns a specific size for a provider.
// Note: This endpoint does not require an organization scope.
func (c *Client) GetProviderSize(ctx context.Context, providerID, sizeID int64) (*ProviderSizeInfo, error) {
	path := fmt.Sprintf("/providers/%d/sizes/%d", providerID, sizeID)
	resp, err := c.GetJsonApiResource(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider size: %w", err)
	}
	return unmarshalSingle(resp, func(s *ProviderSizeInfo, id int64) { s.ID = id })
}

// ListProviderRegionSizes returns all sizes available in a specific region.
// Note: This endpoint does not require an organization scope.
func (c *Client) ListProviderRegionSizes(ctx context.Context, providerID, regionID int64) ([]ProviderRegionSizeInfo, error) {
	path := fmt.Sprintf("/providers/%d/regions/%d/sizes?page[size]=100", providerID, regionID)
	items, err := c.GetJsonApiList(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to list provider region sizes: %w", err)
	}
	return unmarshalList(items, func(s *ProviderRegionSizeInfo, id int64) { s.ID = id })
}

// GetProviderRegionSize returns a specific size in a region.
// Note: This endpoint does not require an organization scope.
func (c *Client) GetProviderRegionSize(ctx context.Context, providerID, regionID, sizeID int64) (*ProviderRegionSizeInfo, error) {
	path := fmt.Sprintf("/providers/%d/regions/%d/sizes/%d", providerID, regionID, sizeID)
	resp, err := c.GetJsonApiResource(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider region size: %w", err)
	}
	return unmarshalSingle(resp, func(s *ProviderRegionSizeInfo, id int64) { s.ID = id })
}
