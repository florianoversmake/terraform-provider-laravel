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

type CommandOutput struct {
	Output string `json:"output"`
}

type ExecuteCommandRequest struct {
	Command string `json:"command"`
}

func (c *Client) ExecuteSiteCommand(ctx context.Context, serverID, siteID int, cmd string) (*Command, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/commands", serverID, siteID))
	req := ExecuteCommandRequest{Command: cmd}
	var command Command
	id, err := c.PostJsonApi(ctx, path, req, &command)
	if err != nil {
		return nil, err
	}
	command.ID = int64(id)
	return &command, nil
}

func (c *Client) ListSiteCommands(ctx context.Context, serverID, siteID int) ([]Command, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/commands", serverID, siteID))
	items, err := c.GetJsonApiListAll(ctx, path)
	if err != nil {
		return nil, err
	}
	return unmarshalList(items, func(c *Command, id int64) { c.ID = id })
}

func (c *Client) GetSiteCommand(ctx context.Context, serverID, siteID, commandID int) (*Command, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/commands/%d", serverID, siteID, commandID))
	var cmd Command
	id, err := c.GetJsonApi(ctx, path, &cmd)
	if err != nil {
		return nil, err
	}
	cmd.ID = int64(id)
	return &cmd, nil
}

func (c *Client) GetSiteCommandOutput(ctx context.Context, serverID, siteID, commandID int) (string, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/commands/%d/output", serverID, siteID, commandID))
	var output CommandOutput
	_, err := c.GetJsonApi(ctx, path, &output)
	if err != nil {
		return "", err
	}
	return output.Output, nil
}

func (c *Client) DeleteSiteCommand(ctx context.Context, serverID, siteID, commandID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/commands/%d", serverID, siteID, commandID))
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}
