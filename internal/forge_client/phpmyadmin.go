// Copyright (c) HashiCorp, Inc.

package forge_client

import (
	"context"
	"fmt"
	"net/http"
)

type PhpMyAdminInstallRequest struct {
	Database string `json:"database"`
	User     int    `json:"user"`
}

func (c *Client) InstallPhpMyAdmin(ctx context.Context, serverID, siteID int, req PhpMyAdminInstallRequest) error {
	path := fmt.Sprintf("/servers/%d/sites/%d/phpmyadmin", serverID, siteID)
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}

func (c *Client) UninstallPhpMyAdmin(ctx context.Context, serverID, siteID int) error {
	path := fmt.Sprintf("/servers/%d/sites/%d/phpmyadmin", serverID, siteID)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}
