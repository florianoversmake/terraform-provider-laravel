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

func (c *Client) InstallWordPress(ctx context.Context, serverID, siteID int, req WordPressInstallRequest) error {
	path := fmt.Sprintf("/servers/%d/sites/%d/wordpress", serverID, siteID)
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}

func (c *Client) UninstallWordPress(ctx context.Context, serverID, siteID int) error {
	path := fmt.Sprintf("/servers/%d/sites/%d/wordpress", serverID, siteID)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}
