// Copyright (c) HashiCorp, Inc.

package envoyer_client

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type LinkedFolder struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type Project struct {
	ID                      int64          `json:"id"`
	UserID                  int64          `json:"user_id"`
	Version                 int64          `json:"version"`
	Name                    string         `json:"name"`
	Provider                string         `json:"provider"`
	PlainRepository         string         `json:"plain_repository"`
	Repository              string         `json:"repository"`
	Type                    string         `json:"type"`
	Branch                  string         `json:"branch"`
	PushToDeploy            bool           `json:"push_to_deploy"`
	WebhookID               *int64         `json:"webhook_id"`
	Status                  *string        `json:"status"`
	ShouldDeployAgain       int64          `json:"should_deploy_again"`
	DeploymentStartedAt     *time.Time     `json:"deployment_started_at"`
	DeploymentFinishedAt    time.Time      `json:"deployment_finished_at"`
	LastDeploymentStatus    string         `json:"last_deployment_status"`
	DailyDeploys            int64          `json:"daily_deploys"`
	WeeklyDeploys           int64          `json:"weekly_deploys"`
	LastDeploymentTook      int64          `json:"last_deployment_took"`
	RetainDeployments       int64          `json:"retain_deployments"`
	EnvironmentServers      []int64        `json:"environment_servers"`
	Folders                 []LinkedFolder `json:"folders"`
	Monitor                 string         `json:"monitor"`
	NewYorkStatus           string         `json:"new_york_status"`
	LondonStatus            string         `json:"london_status"`
	SingaporeStatus         string         `json:"singapore_status"`
	Token                   string         `json:"token"`
	CreatedAt               time.Time      `json:"created_at"`
	UpdatedAt               time.Time      `json:"updated_at"`
	InstallDevDependencies  bool           `json:"install_dev_dependencies"`
	InstallDependencies     bool           `json:"install_dependencies"`
	QuietComposer           bool           `json:"quiet_composer"`
	Servers                 []Server       `json:"servers"`
	HasEnvironment          bool           `json:"has_environment"`
	HasMonitoringError      bool           `json:"has_monitoring_error"`
	HasMissingHeartbeats    bool           `json:"has_missing_heartbeats"`
	LastDeployedBranch      string         `json:"last_deployed_branch"`
	LastDeploymentID        int64          `json:"last_deployment_id"`
	LastDeploymentAuthor    string         `json:"last_deployment_author"`
	LastDeploymentAvatar    string         `json:"last_deployment_avatar"`
	LastDeploymentHash      string         `json:"last_deployment_hash"`
	LastDeploymentTimestamp string         `json:"last_deployment_timestamp"`
}

// projectsResponse is the Envoyer "List Projects" response.
type projectsResponse struct {
	Projects []Project `json:"projects"`
}

// projectResponse is the Envoyer "Get/Update/Create Project" response.
type projectResponse struct {
	Project Project `json:"project"`
}

// ListProjects lists all projects accessible via the API token.
func (c *Client) ListProjects(ctx context.Context) ([]Project, error) {
	var resp projectsResponse
	if err := c.doRequest(ctx, http.MethodGet, "/projects", nil, &resp); err != nil {
		return nil, err
	}
	return resp.Projects, nil
}

// GetProject retrieves the details of a single project by its ID.
func (c *Client) GetProject(ctx context.Context, projectID int) (*Project, error) {
	path := fmt.Sprintf("/projects/%d", projectID)
	var resp projectResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Project, nil
}

// CreateProjectRequest is the payload required to create a new project.
type CreateProjectRequest struct {
	Name              string `json:"name"`
	Provider          string `json:"provider"`
	Repository        string `json:"repository"`
	Branch            string `json:"branch"`
	Type              string `json:"type"`
	RetainDeployments int    `json:"retain_deployments,omitempty"`
	Monitor           string `json:"monitor,omitempty"`
	Composer          bool   `json:"composer,omitempty"`
	ComposerDev       bool   `json:"composer_dev,omitempty"`
	ComposerQuiet     bool   `json:"composer_quiet,omitempty"`
}

// CreateProject creates a new Envoyer project.
func (c *Client) CreateProject(ctx context.Context, req CreateProjectRequest) (*Project, error) {
	var resp projectResponse
	if err := c.doRequest(ctx, http.MethodPost, "/projects", req, &resp); err != nil {
		return nil, err
	}
	return &resp.Project, nil
}

// UpdateProjectRequest is the payload to update a project.
type UpdateProjectRequest struct {
	Name              string `json:"name,omitempty"`
	RetainDeployments int    `json:"retain_deployments,omitempty"`
	Monitor           string `json:"monitor,omitempty"`
	Composer          bool   `json:"composer,omitempty"`
	ComposerDev       bool   `json:"composer_dev,omitempty"`
	ComposerQuiet     bool   `json:"composer_quiet,omitempty"`
}

// UpdateProject updates an existing Envoyer project.
func (c *Client) UpdateProject(ctx context.Context, projectID int, req UpdateProjectRequest) error {
	path := fmt.Sprintf("/projects/%d", projectID)
	return c.doRequest(ctx, http.MethodPut, path, req, nil)
}

// UpdateProjectSourceRequest updates the source code info of a project.
type UpdateProjectSourceRequest struct {
	Provider     string `json:"provider,omitempty"`
	Repository   string `json:"repository,omitempty"`
	Branch       string `json:"branch,omitempty"`
	PushToDeploy bool   `json:"push_to_deploy,omitempty"`
}

// UpdateProjectSource updates a project's source repository information.
func (c *Client) UpdateProjectSource(ctx context.Context, projectID int, req UpdateProjectSourceRequest) error {
	path := fmt.Sprintf("/projects/%d/source", projectID)
	return c.doRequest(ctx, http.MethodPut, path, req, nil)
}

// DeleteProject removes a project by its ID.
func (c *Client) DeleteProject(ctx context.Context, projectID int) error {
	path := fmt.Sprintf("/projects/%d", projectID)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

type foldersResponse struct {
	Folders []LinkedFolder `json:"folders"`
}

// ListLinkedFolders retrieves the project's linked folders.
func (c *Client) ListLinkedFolders(ctx context.Context, projectID int) ([]LinkedFolder, error) {
	path := fmt.Sprintf("/projects/%d/folders", projectID)
	var resp foldersResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return resp.Folders, nil
}

// CreateLinkedFolderRequest is the payload to create or delete a linked folder.
type CreateLinkedFolderRequest struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// CreateLinkedFolder creates a new linked folder on the project.
func (c *Client) CreateLinkedFolder(ctx context.Context, projectID int, folder CreateLinkedFolderRequest) ([]LinkedFolder, error) {
	path := fmt.Sprintf("/projects/%d/folders", projectID)
	var resp foldersResponse
	if err := c.doRequest(ctx, http.MethodPost, path, folder, &resp); err != nil {
		return nil, err
	}
	return resp.Folders, nil
}

// DeleteLinkedFolder removes a linked folder from the project.
func (c *Client) DeleteLinkedFolder(ctx context.Context, projectID int, folder CreateLinkedFolderRequest) error {
	path := fmt.Sprintf("/projects/%d/folders", projectID)
	return c.doRequest(ctx, http.MethodDelete, path, folder, nil)
}
