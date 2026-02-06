package forge_client

import (
	"context"
	"fmt"
	"net/http"
)

type WordPressInstallRequest struct {
	Database string `json:"database"`
	User     int    `json:"user"`
}

// InstallWordPress installs WordPress on a site.
// NOTE: This endpoint may not be available in the new Forge API.
// Check the API documentation for availability.
func (c *Client) InstallWordPress(ctx context.Context, serverID, siteID int, req WordPressInstallRequest) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/wordpress", serverID, siteID))
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}

// UninstallWordPress removes WordPress from a site.
// NOTE: This endpoint may not be available in the new Forge API.
// Check the API documentation for availability.
func (c *Client) UninstallWordPress(ctx context.Context, serverID, siteID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/wordpress", serverID, siteID))
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}
