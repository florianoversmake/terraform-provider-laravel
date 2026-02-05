package forge_client

import (
	"context"
	"fmt"
	"net/http"
)

type UpdateConfigurationRequest struct {
	Content string `json:"content"`
}

func (c *Client) GetNginxConfiguration(ctx context.Context, serverID, siteID int) (string, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/nginx", serverID, siteID))
	var res struct {
		Content string `json:"content"`
	}
	if _, err := c.GetJsonApi(ctx, path, &res); err != nil {
		return "", err
	}
	return res.Content, nil
}

func (c *Client) UpdateNginxConfiguration(ctx context.Context, serverID, siteID int, content string) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/nginx", serverID, siteID))
	req := UpdateConfigurationRequest{Content: content}
	return c.doRequest(ctx, http.MethodPut, path, req, nil)
}

func (c *Client) GetEnvFile(ctx context.Context, serverID, siteID int) (string, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/env", serverID, siteID))
	var res struct {
		Content string `json:"content"`
	}
	if _, err := c.GetJsonApi(ctx, path, &res); err != nil {
		return "", err
	}
	return res.Content, nil
}

func (c *Client) UpdateEnvFile(ctx context.Context, serverID, siteID int, content string) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/env", serverID, siteID))
	req := UpdateConfigurationRequest{Content: content}
	return c.doRequest(ctx, http.MethodPut, path, req, nil)
}
