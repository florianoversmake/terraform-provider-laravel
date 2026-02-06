package forge_client

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
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

type installPHPVersionRequest struct {
	Version     string `json:"version"`
	CLIDefault  bool   `json:"cli_default"`
	SiteDefault bool   `json:"site_default"`
}

type updatePHPDefaultVersionRequest struct {
	PHPVersion string `json:"php_version"`
}

func (c *Client) InstallPHPVersion(ctx context.Context, serverID int, version string) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/php/versions", serverID))
	req := installPHPVersionRequest{Version: version, CLIDefault: false, SiteDefault: false}
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}

// InstallPHPVersionAsDefault installs a PHP version and sets it as both CLI and site default.
func (c *Client) InstallPHPVersionAsDefault(ctx context.Context, serverID int, version string) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/php/versions", serverID))
	req := installPHPVersionRequest{Version: version, CLIDefault: true, SiteDefault: true}
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}

// UpdatePHPCLIVersion updates the default PHP CLI version for a server.
// The version should be in dotted format (e.g., "8.2", "8.4").
func (c *Client) UpdatePHPCLIVersion(ctx context.Context, serverID int, version string) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/php/cli-version", serverID))
	req := updatePHPDefaultVersionRequest{PHPVersion: version}
	return c.doRequest(ctx, http.MethodPut, path, req, nil)
}

// UpdatePHPSiteVersion updates the default PHP site version for a server.
// The version should be in dotted format (e.g., "8.2", "8.4").
func (c *Client) UpdatePHPSiteVersion(ctx context.Context, serverID int, version string) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/php/site-version", serverID))
	req := updatePHPDefaultVersionRequest{PHPVersion: version}
	return c.doRequest(ctx, http.MethodPut, path, req, nil)
}

// WaitForPHPVersionInstalled polls until the given PHP version is installed on the server.
// The version parameter should be in Forge format (e.g., "php82", "php84").
func (c *Client) WaitForPHPVersionInstalled(ctx context.Context, serverID int, version string) error {
	dotted := PHPVersionToDotted(version)
	for {
		versions, err := c.ListPHPVersions(ctx, serverID)
		if err != nil {
			return err
		}
		for _, v := range versions {
			if v.Version == dotted && v.Status == "installed" {
				return nil
			}
		}
		select {
		case <-time.After(10 * time.Second):
			// continue polling
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// IsPHPVersionInstalled checks if a PHP version is already installed on the server.
// The version parameter should be in Forge format (e.g., "php82", "php84").
func (c *Client) IsPHPVersionInstalled(ctx context.Context, serverID int, version string) (bool, error) {
	dotted := PHPVersionToDotted(version)
	versions, err := c.ListPHPVersions(ctx, serverID)
	if err != nil {
		return false, err
	}
	for _, v := range versions {
		if v.Version == dotted && v.Status == "installed" {
			return true, nil
		}
	}
	return false, nil
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

// PHPVersionToDotted converts a Forge php_version string (e.g., "php82") to dotted format (e.g., "8.2").
// If the input is already in dotted format or unrecognized, it is returned as-is.
func PHPVersionToDotted(version string) string {
	v := strings.TrimPrefix(version, "php")
	if v == version {
		// No "php" prefix — might already be dotted or a different format.
		return version
	}
	switch v {
	case "5", "56-old":
		return "5.6"
	}
	if len(v) == 2 {
		return string(v[0]) + "." + string(v[1])
	}
	// Fallback for longer strings like "56".
	if len(v) >= 2 {
		return string(v[0]) + "." + v[1:]
	}
	return version
}

// PHPDottedToForgeVersion converts a dotted PHP version (e.g., "8.2") to Forge format (e.g., "php82").
func PHPDottedToForgeVersion(dotted string) string {
	return "php" + strings.ReplaceAll(dotted, ".", "")
}
