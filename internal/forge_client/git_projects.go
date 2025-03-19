package forge_client

import (
	"context"
	"fmt"
	"net/http"
)

type GitProjectRequest struct {
	Provider   string `json:"provider"`
	Repository string `json:"repository"`
	Branch     string `json:"branch"`
	Composer   bool   `json:"composer,omitempty"`
}

func (c *Client) InstallGitProject(ctx context.Context, serverID, siteID int, req GitProjectRequest) error {
	path := fmt.Sprintf("/servers/%d/sites/%d/git", serverID, siteID)
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}

func (c *Client) UpdateGitProject(ctx context.Context, serverID, siteID int, req GitProjectRequest) error {
	path := fmt.Sprintf("/servers/%d/sites/%d/git", serverID, siteID)
	return c.doRequest(ctx, http.MethodPut, path, req, nil)
}

func (c *Client) RemoveGitProject(ctx context.Context, serverID, siteID int) error {
	path := fmt.Sprintf("/servers/%d/sites/%d/git", serverID, siteID)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

type DeployKeyResponse struct {
	Key string `json:"key"`
}

func (c *Client) CreateDeployKey(ctx context.Context, serverID, siteID int) (string, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/deploy-key", serverID, siteID)
	var res DeployKeyResponse
	if err := c.doRequest(ctx, http.MethodPost, path, nil, &res); err != nil {
		return "", err
	}
	return res.Key, nil
}

func (c *Client) DeleteDeployKey(ctx context.Context, serverID, siteID int) error {
	path := fmt.Sprintf("/servers/%d/sites/%d/deploy-key", serverID, siteID)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}
