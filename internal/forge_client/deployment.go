package forge_client

import (
	"context"
	"fmt"
	"net/http"
)

type Deployment struct {
	ID        int64             `json:"id"`
	Status    string            `json:"status"`
	Type      string            `json:"type"`
	StartedAt string            `json:"started_at"`
	EndedAt   string            `json:"ended_at"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
	Commit    *DeploymentCommit `json:"commit"`
}

type DeploymentCommit struct {
	Hash    *string `json:"hash"`
	Author  *string `json:"author"`
	Message *string `json:"message"`
	Branch  *string `json:"branch"`
}

type DeploymentScript struct {
	Content    *string `json:"content"`
	AutoSource bool    `json:"auto_source"`
}

func (c *Client) EnableQuickDeployment(ctx context.Context, serverID, siteID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/deployments/push-to-deploy", serverID, siteID))
	return c.doRequest(ctx, http.MethodPost, path, nil, nil)
}

func (c *Client) DisableQuickDeployment(ctx context.Context, serverID, siteID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/deployments/push-to-deploy", serverID, siteID))
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

func (c *Client) GetDeploymentScript(ctx context.Context, serverID, siteID int) (string, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/deployments/script", serverID, siteID))
	var script DeploymentScript
	_, err := c.GetJsonApi(ctx, path, &script)
	if err != nil {
		return "", err
	}
	if script.Content == nil {
		return "", nil
	}
	return *script.Content, nil
}

type UpdateDeploymentScriptRequest struct {
	Content    string `json:"content"`
	AutoSource bool   `json:"auto_source"`
}

func (c *Client) UpdateDeploymentScript(ctx context.Context, serverID, siteID int, req UpdateDeploymentScriptRequest) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/deployments/script", serverID, siteID))
	return c.doRequest(ctx, http.MethodPut, path, req, nil)
}

func (c *Client) DeployNow(ctx context.Context, serverID, siteID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/deployments", serverID, siteID))
	return c.doRequest(ctx, http.MethodPost, path, nil, nil)
}

func (c *Client) ResetDeploymentStatus(ctx context.Context, serverID, siteID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/deployments/status", serverID, siteID))
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

func (c *Client) GetDeploymentLog(ctx context.Context, serverID, siteID int) (string, error) {
	// Get the latest deployment first
	deployments, err := c.ListDeployments(ctx, serverID, siteID)
	if err != nil {
		return "", err
	}
	if len(deployments) == 0 {
		return "", nil
	}
	return c.GetDeploymentOutput(ctx, serverID, siteID, int(deployments[0].ID))
}

func (c *Client) ListDeployments(ctx context.Context, serverID, siteID int) ([]Deployment, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/deployments", serverID, siteID))
	items, err := c.GetJsonApiListAll(ctx, path)
	if err != nil {
		return nil, err
	}
	return unmarshalList(items, func(d *Deployment, id int64) { d.ID = id })
}

func (c *Client) GetDeployment(ctx context.Context, serverID, siteID, deploymentID int) (*Deployment, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/deployments/%d", serverID, siteID, deploymentID))
	var deployment Deployment
	id, err := c.GetJsonApi(ctx, path, &deployment)
	if err != nil {
		return nil, err
	}
	deployment.ID = int64(id)
	return &deployment, nil
}

func (c *Client) GetDeploymentOutput(ctx context.Context, serverID, siteID, deploymentID int) (string, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/deployments/%d/log", serverID, siteID, deploymentID))
	return c.GetText(ctx, path)
}

type DeploymentFailureEmailsRequest struct {
	Emails []string `json:"emails"`
}

// SetDeploymentFailureEmails is DEPRECATED - this endpoint no longer exists in the new Forge API.
// This function is kept for backward compatibility but will return a 404 error.
func (c *Client) SetDeploymentFailureEmails(ctx context.Context, serverID, siteID int, emails []string) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/deployment-failure-emails", serverID, siteID))
	req := DeploymentFailureEmailsRequest{Emails: emails}
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}
