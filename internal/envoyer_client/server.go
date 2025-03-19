package envoyer_client

import (
	"context"
	"fmt"
	"net/http"
)

type Server struct {
	ID                  int64   `json:"id"`
	ProjectID           int64   `json:"project_id"`
	Name                string  `json:"name"`
	ConnectAs           string  `json:"connect_as"`
	IPAddress           string  `json:"ip_address"`
	Port                string  `json:"port"`
	PHPVersion          string  `json:"php_version"`
	ReceivesCodeDeploys bool    `json:"receives_code_deployments"`
	ShouldRestartFPM    bool    `json:"should_restart_fpm"`
	DeploymentPath      string  `json:"deployment_path"`
	PHPPath             string  `json:"php_path"`
	ComposerPath        string  `json:"composer_path"`
	ConnectionStatus    string  `json:"connection_status"`
	CurrentActivity     *string `json:"current_activity"`
	PublicKey           string  `json:"public_key"`
	CreatedAt           string  `json:"created_at"`
	UpdatedAt           string  `json:"updated_at"`
}

// serversResponse is the Envoyer "List Servers" response.
type serversResponse struct {
	Servers []Server `json:"servers"`
}

// serverResponse is the Envoyer "Get/Create/Update Server" response.
type serverResponse struct {
	Server Server `json:"server"`
}

// ListServers returns the servers for a given project.
func (c *Client) ListServers(ctx context.Context, projectID int) ([]Server, error) {
	path := fmt.Sprintf("/projects/%d/servers", projectID)
	var resp serversResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return resp.Servers, nil
}

// GetServer returns a single server for the given project.
func (c *Client) GetServer(ctx context.Context, projectID, serverID int) (*Server, error) {
	path := fmt.Sprintf("/projects/%d/servers/%d", projectID, serverID)
	var resp serverResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Server, nil
}

// RefreshServerConnection refreshes the Envoyer connection for the given server.
func (c *Client) RefreshServerConnection(ctx context.Context, projectID, serverID int) error {
	path := fmt.Sprintf("/projects/%d/servers/%d/refresh", projectID, serverID)
	return c.doRequest(ctx, http.MethodPost, path, nil, nil)
}

// CreateServerRequest is the payload to create or update a server.
type CreateServerRequest struct {
	Name                string `json:"name"`
	ConnectAs           string `json:"connectAs"`
	Host                string `json:"host"`
	Port                int    `json:"port,omitempty"`
	PHPVersion          string `json:"phpVersion"`
	ReceivesCodeDeploys bool   `json:"receivesCodeDeployments"`
	DeploymentPath      string `json:"deploymentPath"`
	RestartFpm          bool   `json:"restartFpm"`
	ComposerPath        string `json:"composerPath,omitempty"`
	PHPPath             string `json:"phpPath,omitempty"`
}

// CreateServer creates a new server on a given project.
func (c *Client) CreateServer(ctx context.Context, projectID int, req CreateServerRequest) (*Server, error) {
	path := fmt.Sprintf("/projects/%d/servers", projectID)
	var resp serverResponse
	if err := c.doRequest(ctx, http.MethodPost, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp.Server, nil
}

// UpdateServer updates an existing server.
func (c *Client) UpdateServer(ctx context.Context, projectID, serverID int, req CreateServerRequest) (*Server, error) {
	path := fmt.Sprintf("/projects/%d/servers/%d", projectID, serverID)
	var resp serverResponse
	if err := c.doRequest(ctx, http.MethodPut, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp.Server, nil
}

// DeleteServer removes a server from a project.
func (c *Client) DeleteServer(ctx context.Context, projectID, serverID int) error {
	path := fmt.Sprintf("/projects/%d/servers/%d", projectID, serverID)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}
