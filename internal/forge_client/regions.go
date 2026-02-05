package forge_client

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
)

// RegionSize represents a size available in a region for the cached region data.
type RegionSize struct {
	ID   string `json:"id"`
	Size string `json:"size"`
	Name string `json:"name"`
}

// Region represents a cached region with its available sizes.
type Region struct {
	ID    string       `json:"id"`
	Name  string       `json:"name"`
	Sizes []RegionSize `json:"sizes"`
}

type cachedRegions struct {
	data      map[string][]Region
	timestamp time.Time
}

var (
	regionsCache      cachedRegions
	regionsCacheMutex sync.Mutex
)

// GetProviderBySlug returns a provider by its slug (e.g., "aws", "digitalocean").
func (c *Client) GetProviderBySlug(ctx context.Context, slug string) (*ProviderInfo, error) {
	providers, err := c.ListProviders(ctx)
	if err != nil {
		return nil, err
	}
	for _, p := range providers {
		if p.Slug == slug {
			return &p, nil
		}
	}
	return nil, fmt.Errorf("provider with slug %q not found", slug)
}

// GetRegionByCode returns a region by its code for a given provider slug.
func (c *Client) GetRegionByCode(ctx context.Context, providerSlug, regionCode string) (*ProviderRegionInfo, error) {
	provider, err := c.GetProviderBySlug(ctx, providerSlug)
	if err != nil {
		return nil, err
	}
	regions, err := c.ListProviderRegions(ctx, provider.ID)
	if err != nil {
		return nil, err
	}
	for _, r := range regions {
		if r.Code == regionCode {
			return &r, nil
		}
	}
	return nil, fmt.Errorf("region with code %q not found for provider %q", regionCode, providerSlug)
}

// GetRegionByName returns a region by its name for a given provider slug.
func (c *Client) GetRegionByName(ctx context.Context, providerSlug, regionName string) (*ProviderRegionInfo, error) {
	provider, err := c.GetProviderBySlug(ctx, providerSlug)
	if err != nil {
		return nil, err
	}
	regions, err := c.ListProviderRegions(ctx, provider.ID)
	if err != nil {
		return nil, err
	}
	for _, r := range regions {
		if r.Name == regionName {
			return &r, nil
		}
	}
	return nil, fmt.Errorf("region with name %q not found for provider %q", regionName, providerSlug)
}

// GetSizeByCode returns a size by its code for a given provider slug.
func (c *Client) GetSizeByCode(ctx context.Context, providerSlug, sizeCode string) (*ProviderSizeInfo, error) {
	provider, err := c.GetProviderBySlug(ctx, providerSlug)
	if err != nil {
		return nil, err
	}
	sizes, err := c.ListProviderSizes(ctx, provider.ID)
	if err != nil {
		return nil, err
	}
	for _, s := range sizes {
		if s.Code == sizeCode {
			return &s, nil
		}
	}
	return nil, fmt.Errorf("size with code %q not found for provider %q", sizeCode, providerSlug)
}

// GetSizeByID returns a size by its ID for a given provider slug.
func (c *Client) GetSizeByID(ctx context.Context, providerSlug string, sizeID int64) (*ProviderSizeInfo, error) {
	provider, err := c.GetProviderBySlug(ctx, providerSlug)
	if err != nil {
		return nil, err
	}
	return c.GetProviderSize(ctx, provider.ID, sizeID)
}

// GetRegionIDByName returns the region code for a region name and provider.
// This is a convenience method that wraps GetRegionByName.
func (c *Client) GetRegionIDByName(ctx context.Context, providerSlug, regionName string) (string, error) {
	region, err := c.GetRegionByName(ctx, providerSlug, regionName)
	if err != nil {
		// Return empty string for not found (backward compatibility)
		return "", nil
	}
	return region.Code, nil
}

// GetRegionNameByID returns the region name for a region code and provider.
// This is a convenience method that wraps GetRegionByCode.
func (c *Client) GetRegionNameByID(ctx context.Context, providerSlug, regionCode string) (string, error) {
	region, err := c.GetRegionByCode(ctx, providerSlug, regionCode)
	if err != nil {
		// Return empty string for not found (backward compatibility)
		return "", nil
	}
	return region.Name, nil
}

// GetRegionSizeIDByName returns the size ID for a size name in a region.
// This fetches sizes available in the specific region.
func (c *Client) GetRegionSizeIDByName(ctx context.Context, providerSlug, regionCode, sizeName string) (string, error) {
	provider, err := c.GetProviderBySlug(ctx, providerSlug)
	if err != nil {
		return "", nil
	}
	region, err := c.GetRegionByCode(ctx, providerSlug, regionCode)
	if err != nil {
		return "", nil
	}
	// Get all sizes for provider to get full size info
	allSizes, err := c.ListProviderSizes(ctx, provider.ID)
	if err != nil {
		return "", err
	}
	sizeByID := make(map[int64]ProviderSizeInfo)
	for _, s := range allSizes {
		sizeByID[s.ID] = s
	}
	// Get sizes available in this region
	regionSizes, err := c.ListProviderRegionSizes(ctx, provider.ID, region.ID)
	if err != nil {
		return "", err
	}
	for _, rs := range regionSizes {
		if fullSize, ok := sizeByID[rs.ID]; ok && fullSize.Name == sizeName {
			return strconv.FormatInt(rs.ID, 10), nil
		}
	}
	return "", nil
}

// GetRegionSizeNameByID returns the size name for a size ID in a region.
func (c *Client) GetRegionSizeNameByID(ctx context.Context, providerSlug, regionCode, sizeID string) (string, error) {
	provider, err := c.GetProviderBySlug(ctx, providerSlug)
	if err != nil {
		return "", nil
	}
	sizeIDInt, err := strconv.ParseInt(sizeID, 10, 64)
	if err != nil {
		return "", nil
	}
	size, err := c.GetProviderSize(ctx, provider.ID, sizeIDInt)
	if err != nil {
		return "", nil
	}
	return size.Name, nil
}

// GetRegionSizeIDBySize returns the size ID for a size code in a region.
func (c *Client) GetRegionSizeIDBySize(ctx context.Context, providerSlug, regionCode, sizeCode string) (string, error) {
	provider, err := c.GetProviderBySlug(ctx, providerSlug)
	if err != nil {
		return "", nil
	}
	sizes, err := c.ListProviderSizes(ctx, provider.ID)
	if err != nil {
		return "", err
	}
	for _, s := range sizes {
		if s.Code == sizeCode {
			return strconv.FormatInt(s.ID, 10), nil
		}
	}
	return "", nil
}

// GetRegionSizeNameBySize returns the size name for a size code in a region.
func (c *Client) GetRegionSizeNameBySize(ctx context.Context, providerSlug, regionCode, sizeCode string) (string, error) {
	size, err := c.GetSizeByCode(ctx, providerSlug, sizeCode)
	if err != nil {
		return "", nil
	}
	return size.Name, nil
}

// GetRegionSizeSizeByName returns the size code for a size name in a region.
func (c *Client) GetRegionSizeSizeByName(ctx context.Context, providerSlug, regionCode, sizeName string) (string, error) {
	provider, err := c.GetProviderBySlug(ctx, providerSlug)
	if err != nil {
		return "", nil
	}
	sizes, err := c.ListProviderSizes(ctx, provider.ID)
	if err != nil {
		return "", err
	}
	for _, s := range sizes {
		if s.Name == sizeName {
			return s.Code, nil
		}
	}
	return "", nil
}

// GetRegionSizeSizeByID returns the size code for a size ID.
func (c *Client) GetRegionSizeSizeByID(ctx context.Context, providerSlug, regionCode, sizeID string) (string, error) {
	provider, err := c.GetProviderBySlug(ctx, providerSlug)
	if err != nil {
		return "", nil
	}
	sizeIDInt, err := strconv.ParseInt(sizeID, 10, 64)
	if err != nil {
		return "", nil
	}
	size, err := c.GetProviderSize(ctx, provider.ID, sizeIDInt)
	if err != nil {
		return "", nil
	}
	return size.Code, nil
}

// GetAllRegionsWithSizes returns all regions for all providers with their available sizes.
// This method makes many API calls and caches the result for 60 seconds.
// Note: This is a heavy operation that may hit rate limits due to the number of API calls required.
// Prefer using GetRegionByCode, GetSizeByCode, etc. for targeted lookups.
func (c *Client) GetAllRegionsWithSizes(ctx context.Context) (map[string][]Region, error) {
	regionsCacheMutex.Lock()
	defer regionsCacheMutex.Unlock()

	if time.Since(regionsCache.timestamp) < 60*time.Second && regionsCache.data != nil {
		return regionsCache.data, nil
	}

	providers, err := c.ListProviders(ctx)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] GetAllRegionsWithSizes: found %d providers", len(providers))
	for i, p := range providers {
		log.Printf("[DEBUG]   provider[%d]: id=%d slug=%q name=%q", i, p.ID, p.Slug, p.Name)
	}

	apiCalls := 1 // providers call
	result := make(map[string][]Region)
	for _, provider := range providers {
		// Fetch all sizes for this provider to get size codes
		allSizes, err := c.ListProviderSizes(ctx, provider.ID)
		apiCalls++
		if err != nil {
			return nil, err
		}
		log.Printf("[DEBUG] Provider %q: %d sizes", provider.Slug, len(allSizes))
		sizeByID := make(map[int64]ProviderSizeInfo)
		for _, s := range allSizes {
			sizeByID[s.ID] = s
		}

		regions, err := c.ListProviderRegions(ctx, provider.ID)
		apiCalls++
		if err != nil {
			return nil, err
		}
		log.Printf("[DEBUG] Provider %q: %d regions", provider.Slug, len(regions))

		var regionList []Region
		for _, r := range regions {
			regionSizes, err := c.ListProviderRegionSizes(ctx, provider.ID, r.ID)
			apiCalls++
			if err != nil {
				log.Printf("[DEBUG] FAILED at api call #%d (provider %q region %q id=%d)", apiCalls, provider.Slug, r.Code, r.ID)
				return nil, err
			}

			var sizeList []RegionSize
			for _, rs := range regionSizes {
				fullSize, ok := sizeByID[rs.ID]
				if ok {
					sizeList = append(sizeList, RegionSize{
						ID:   strconv.FormatInt(rs.ID, 10),
						Size: fullSize.Code,
						Name: fullSize.Name,
					})
				}
			}

			regionList = append(regionList, Region{
				ID:    r.Code,
				Name:  r.Name,
				Sizes: sizeList,
			})
		}

		result[provider.Slug] = regionList
	}
	log.Printf("[DEBUG] GetAllRegionsWithSizes: total API calls made: %d", apiCalls)

	regionsCache = cachedRegions{
		data:      result,
		timestamp: time.Now(),
	}

	return result, nil
}

// ListRegions is an alias for GetAllRegionsWithSizes for backward compatibility.
// Deprecated: Use GetAllRegionsWithSizes or the more efficient targeted lookup methods instead.
func (c *Client) ListRegions(ctx context.Context) (map[string][]Region, error) {
	return c.GetAllRegionsWithSizes(ctx)
}
