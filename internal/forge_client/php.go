package forge_client

import (
	"context"
	"fmt"
	"net/http"
)

type PHPVersion struct {
	ID                 int    `json:"id"`
	Version            string `json:"version"`
	Status             string `json:"status"`
	DisplayableVersion string `json:"displayable_version"`
	BinaryName         string `json:"binary_name"`
	UsedAsDefault      bool   `json:"used_as_default"`
	UsedOnCLI          bool   `json:"used_on_cli"`
}

func (c *Client) ListPHPVersions(ctx context.Context, serverID int) ([]PHPVersion, error) {
	// Fetch PHP versions from API
	path := fmt.Sprintf("/orgs/%s/servers/%d/php", c.OrgSlug, serverID)
	var res struct {
		Data []PHPVersion `json:"data"`
	}
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}

	return res.Data, nil
}

type phpVersionRequest struct {
	Version string `json:"version"`
}

func (c *Client) InstallPHPVersion(ctx context.Context, serverID int, version string) error {
	path := fmt.Sprintf("/orgs/%s/servers/%d/php", c.OrgSlug, serverID)
	req := phpVersionRequest{Version: version}
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}

func (c *Client) UpgradePHPPatchVersion(ctx context.Context, serverID int, version string) error {
	path := fmt.Sprintf("/orgs/%s/servers/%d/php/update", c.OrgSlug, serverID)
	req := phpVersionRequest{Version: version}
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}

func (c *Client) EnableOPCache(ctx context.Context, serverID int) error {
	path := fmt.Sprintf("/orgs/%s/servers/%d/php/opcache", c.OrgSlug, serverID)
	return c.doRequest(ctx, http.MethodPost, path, nil, nil)
}

func (c *Client) DisableOPCache(ctx context.Context, serverID int) error {
	path := fmt.Sprintf("/orgs/%s/servers/%d/php/opcache", c.OrgSlug, serverID)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

func (c *Client) GetPHPVersionFromDisplayableVersion(ctx context.Context, serverID int, displayableVersion string) (*PHPVersion, error) {
	versions, err := c.ListPHPVersions(ctx, serverID)
	if err != nil {
		return nil, err
	}

	for _, version := range versions {
		if version.DisplayableVersion == displayableVersion {
			return &version, nil
		}
	}

	return nil, fmt.Errorf("php version not found: %s", displayableVersion)
}
