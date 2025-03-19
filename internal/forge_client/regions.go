// Copyright (c) HashiCorp, Inc.

package forge_client

import (
	"context"
	"net/http"
	"sync"
	"time"
)

type RegionSize struct {
	ID   string `json:"id"`
	Size string `json:"size"`
	Name string `json:"name"`
}

type Region struct {
	ID    string       `json:"id"`
	Name  string       `json:"name"`
	Sizes []RegionSize `json:"sizes"`
}

type RegionsResponse struct {
	Regions map[string][]Region `json:"regions"`
}

type cachedRegions struct {
	data      map[string][]Region
	timestamp time.Time
}

var (
	regionsCache      cachedRegions
	regionsCacheMutex sync.Mutex
)

func (c *Client) ListRegions(ctx context.Context) (map[string][]Region, error) {
	regionsCacheMutex.Lock()
	defer regionsCacheMutex.Unlock()

	// Check if cache is valid
	if time.Since(regionsCache.timestamp) < 60*time.Second && regionsCache.data != nil {
		return regionsCache.data, nil
	}

	// Fetch regions from API
	var res RegionsResponse
	if err := c.doRequest(ctx, http.MethodGet, "/regions", nil, &res); err != nil {
		return nil, err
	}

	// Update cache
	regionsCache = cachedRegions{
		data:      res.Regions,
		timestamp: time.Now(),
	}

	return res.Regions, nil
}

// {
// "regions": {
// 	"ocean2": [
// 	{
// 		"id": "ams2",
// 		"name": "Amsterdam 2",
// 		"sizes": [
// 		{
// 			"id": "01",
// 			"size": "s-1vcpu-1gb",
// 			"name": "1GB RAM - 1 CPU Core - 25GB SSD"
// 		}
// 		]
// 	}
// 	],
// 	"linode": [],
// 	"vultr": [],
// 	"aws": []
// }
// }

func (c *Client) GetRegionIDByName(ctx context.Context, providerName string, regionName string) (string, error) {
	regions, err := c.ListRegions(ctx)
	if err != nil {
		return "", err
	}

	for _, region := range regions[providerName] {
		if region.Name == regionName {
			return region.ID, nil
		}
	}

	return "", nil
}

func (c *Client) GetRegionNameByID(ctx context.Context, providerName string, regionID string) (string, error) {
	regions, err := c.ListRegions(ctx)
	if err != nil {
		return "", err
	}

	for _, region := range regions[providerName] {
		if region.ID == regionID {
			return region.Name, nil
		}
	}

	return "", nil
}

func (c *Client) GetRegionSizeIDByName(ctx context.Context, providerName string, regionName string, sizeName string) (string, error) {
	regions, err := c.ListRegions(ctx)
	if err != nil {
		return "", err
	}

	for _, region := range regions[providerName] {
		if region.ID == regionName {
			for _, size := range region.Sizes {
				if size.Name == sizeName {
					return size.ID, nil
				}
			}
		}
	}

	return "", nil
}

func (c *Client) GetRegionSizeNameByID(ctx context.Context, providerName string, regionName string, sizeID string) (string, error) {
	regions, err := c.ListRegions(ctx)
	if err != nil {
		return "", err
	}

	for _, region := range regions[providerName] {
		if region.ID == regionName {
			for _, size := range region.Sizes {
				if size.ID == sizeID {
					return size.Name, nil
				}
			}
		}
	}

	return "", nil
}

func (c *Client) GetRegionSizeIDBySize(ctx context.Context, providerName string, regionName string, sizeName string) (string, error) {
	regions, err := c.ListRegions(ctx)
	if err != nil {
		return "", err
	}

	for _, region := range regions[providerName] {
		if region.ID == regionName {
			for _, size := range region.Sizes {
				if size.Size == sizeName {
					return size.ID, nil
				}
			}
		}
	}

	return "", nil
}

func (c *Client) GetRegionSizeNameBySize(ctx context.Context, providerName string, regionName string, sizeName string) (string, error) {
	regions, err := c.ListRegions(ctx)
	if err != nil {
		return "", err
	}

	for _, region := range regions[providerName] {
		if region.ID == regionName {
			for _, size := range region.Sizes {
				if size.Size == sizeName {
					return size.Name, nil
				}
			}
		}
	}

	return "", nil
}

func (c *Client) GetRegionSizeSizeByName(ctx context.Context, providerName string, regionName string, sizeName string) (string, error) {
	regions, err := c.ListRegions(ctx)
	if err != nil {
		return "", err
	}

	for _, region := range regions[providerName] {
		if region.ID == regionName {
			for _, size := range region.Sizes {
				if size.Name == sizeName {
					return size.Size, nil
				}
			}
		}
	}

	return "", nil
}

func (c *Client) GetRegionSizeSizeByID(ctx context.Context, providerName string, regionName string, sizeID string) (string, error) {
	regions, err := c.ListRegions(ctx)
	if err != nil {
		return "", err
	}

	for _, region := range regions[providerName] {
		if region.ID == regionName {
			for _, size := range region.Sizes {
				if size.ID == sizeID {
					return size.Size, nil
				}
			}
		}
	}

	return "", nil
}
