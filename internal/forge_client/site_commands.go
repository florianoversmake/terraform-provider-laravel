// Copyright (c) HashiCorp, Inc.

package forge_client

import (
	"context"
	"fmt"
	"net/http"
)

type Command struct {
	ID              int64  `json:"id"`
	ServerID        int64  `json:"server_id"`
	SiteID          int64  `json:"site_id"`
	UserID          int64  `json:"user_id"`
	EventID         int64  `json:"event_id"`
	Command         string `json:"command"`
	Status          string `json:"status"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	ProfilePhotoURL string `json:"profile_photo_url"`
	UserName        string `json:"user_name"`
}

type commandResponse struct {
	Command Command `json:"command"`
	Output  string  `json:"output,omitempty"`
}

type commandsResponse struct {
	Commands []Command `json:"commands"`
}

type ExecuteCommandRequest struct {
	Command string `json:"command"`
}

func (c *Client) ExecuteSiteCommand(ctx context.Context, serverID, siteID int, cmd string) error {
	path := fmt.Sprintf("/servers/%d/sites/%d/commands", serverID, siteID)
	req := ExecuteCommandRequest{Command: cmd}
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}

func (c *Client) ListSiteCommands(ctx context.Context, serverID, siteID int) ([]Command, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/commands", serverID, siteID)
	var res commandsResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return res.Commands, nil
}

func (c *Client) GetSiteCommand(ctx context.Context, serverID, siteID, commandID int) (*Command, string, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/commands/%d", serverID, siteID, commandID)
	var res commandResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, "", err
	}
	return &res.Command, res.Output, nil
}
