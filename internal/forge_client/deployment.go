package forge_client

import (
	"context"
	"fmt"
	"net/http"
)

type Deployment struct {
	ID              int64  `json:"id"`
	ServerID        int64  `json:"server_id"`
	SiteID          int64  `json:"site_id"`
	Type            int    `json:"type"`
	CommitHash      string `json:"commit_hash"`
	CommitAuthor    string `json:"commit_author"`
	CommitMessage   string `json:"commit_message"`
	StartedAt       string `json:"started_at"`
	EndedAt         string `json:"ended_at"`
	Status          string `json:"status"`
	DisplayableType string `json:"displayable_type"`
}

func (c *Client) EnableQuickDeployment(ctx context.Context, serverID, siteID int) error {
	path := fmt.Sprintf("/servers/%d/sites/%d/deployment", serverID, siteID)
	return c.doRequest(ctx, http.MethodPost, path, nil, nil)
}

func (c *Client) DisableQuickDeployment(ctx context.Context, serverID, siteID int) error {
	path := fmt.Sprintf("/servers/%d/sites/%d/deployment", serverID, siteID)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

func (c *Client) GetDeploymentScript(ctx context.Context, serverID, siteID int) (string, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/deployment/script", serverID, siteID)
	var res struct {
		Script string `json:"script"`
	}
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return "", err
	}
	return res.Script, nil
}

type UpdateDeploymentScriptRequest struct {
	Content    string `json:"content"`
	AutoSource bool   `json:"auto_source"`
}

func (c *Client) UpdateDeploymentScript(ctx context.Context, serverID, siteID int, req UpdateDeploymentScriptRequest) error {
	path := fmt.Sprintf("/servers/%d/sites/%d/deployment/script", serverID, siteID)
	return c.doRequest(ctx, http.MethodPut, path, req, nil)
}

func (c *Client) DeployNow(ctx context.Context, serverID, siteID int) error {
	path := fmt.Sprintf("/servers/%d/sites/%d/deployment/deploy", serverID, siteID)
	return c.doRequest(ctx, http.MethodPost, path, nil, nil)
}

func (c *Client) ResetDeploymentStatus(ctx context.Context, serverID, siteID int) error {
	path := fmt.Sprintf("/servers/%d/sites/%d/deployment/reset", serverID, siteID)
	return c.doRequest(ctx, http.MethodPost, path, nil, nil)
}

func (c *Client) GetDeploymentLog(ctx context.Context, serverID, siteID int) (string, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/deployment/log", serverID, siteID)
	var res struct {
		Log string `json:"log"`
	}
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return "", err
	}
	return res.Log, nil
}

type deploymentsResponse struct {
	Deployments []Deployment `json:"deployments"`
}

type deploymentResponse struct {
	Deployment Deployment `json:"deployment"`
}

func (c *Client) ListDeployments(ctx context.Context, serverID, siteID int) ([]Deployment, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/deployment-history", serverID, siteID)
	var res deploymentsResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return res.Deployments, nil
}

func (c *Client) GetDeployment(ctx context.Context, serverID, siteID, deploymentID int) (*Deployment, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/deployment-history/%d", serverID, siteID, deploymentID)
	var res deploymentResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return &res.Deployment, nil
}

type deploymentOutputResponse struct {
	Output string `json:"output"`
}

func (c *Client) GetDeploymentOutput(ctx context.Context, serverID, siteID, deploymentID int) (string, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/deployment-history/%d/output", serverID, siteID, deploymentID)
	var res deploymentOutputResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return "", err
	}
	return res.Output, nil
}

type DeploymentFailureEmailsRequest struct {
	Emails []string `json:"emails"`
}

func (c *Client) SetDeploymentFailureEmails(ctx context.Context, serverID, siteID int, emails []string) error {
	path := fmt.Sprintf("/servers/%d/sites/%d/deployment-failure-emails", serverID, siteID)
	req := DeploymentFailureEmailsRequest{Emails: emails}
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}
