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
	path := c.orgPath(fmt.Sprintf("/servers/%d/php/versions", serverID))
	items, err := c.GetJsonApiListAll(ctx, path)
	if err != nil {
		return nil, err
	}
	return unmarshalList(items, func(v *PHPVersion, id int64) { v.ID = int(id) })
}

type phpVersionRequest struct {
	Version string `json:"version"`
}

func (c *Client) InstallPHPVersion(ctx context.Context, serverID int, version string) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/php/versions", serverID))
	req := phpVersionRequest{Version: version}
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}

func (c *Client) UpgradePHPPatchVersion(ctx context.Context, serverID int, version string) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/php/update", serverID))
	req := phpVersionRequest{Version: version}
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}

func (c *Client) EnableOPCache(ctx context.Context, serverID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/php/opcache", serverID))
	return c.doRequest(ctx, http.MethodPost, path, nil, nil)
}

func (c *Client) DisableOPCache(ctx context.Context, serverID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/php/opcache", serverID))
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
