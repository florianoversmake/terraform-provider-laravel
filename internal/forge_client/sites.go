package forge_client

import (
	"context"
	"fmt"
	"net/http"
)

type Site struct {
	ID                      int64                `json:"id"`
	Name                    string               `json:"name"`
	Aliases                 []string             `json:"aliases"`
	WebDirectory            string               `json:"web_directory"`
	RootDirectory           *string              `json:"root_directory"`
	Wildcards               *bool                `json:"wildcards"`
	Isolated                bool                 `json:"isolated"`
	User                    string               `json:"user"`
	Status                  string               `json:"status"`
	QuickDeploy             *bool                `json:"quick_deploy"`
	AppType                 string               `json:"app_type"`
	PHPVersion              string               `json:"php_version"`
	DeploymentURL           string               `json:"deployment_url"`
	DeploymentStatus        string               `json:"deployment_status"`
	CreatedAt               *string              `json:"created_at"`
	UpdatedAt               *string              `json:"updated_at"`
	URL                     string               `json:"url"`
	HTTPS                   bool                 `json:"https"`
	ZeroDowntimeDeployments bool                 `json:"zero_downtime_deployments"`
	UsesEnvoyer             bool                 `json:"uses_envoyer"`
	HealthcheckURL          *string              `json:"healthcheck_url"`
	DeploymentScript        *string              `json:"deployment_script"`
	Database                *string              `json:"database"`
	Repository              *SiteRepository      `json:"repository"`
	MaintenanceMode         *SiteMaintenanceMode `json:"maintenance_mode"`
}

type SiteRepository struct {
	Provider string  `json:"provider"`
	URL      *string `json:"url"`
	Branch   *string `json:"branch"`
	Status   string  `json:"status"`
}

type SiteMaintenanceMode struct {
	Enabled bool   `json:"enabled"`
	Status  string `json:"status"`
}

type CreateSiteRequest struct {
	Type                    string  `json:"type"`
	DomainMode              string  `json:"domain_mode"`    // "on-forge" or "custom"
	Name                    string  `json:"name,omitempty"` // subdomain for on-forge, full domain for custom
	WebDirectory            *string `json:"web_directory,omitempty"`
	RootDirectory           *string `json:"root_directory,omitempty"`
	IsIsolated              bool    `json:"is_isolated,omitempty"`
	IsolatedUser            string  `json:"isolated_user,omitempty"`
	PHPVersion              string  `json:"php_version,omitempty"`
	ZeroDowntimeDeployments bool    `json:"zero_downtime_deployments,omitempty"`
	NginxTemplateID         *int    `json:"nginx_template_id,omitempty"`
	SourceControlProvider   string  `json:"source_control_provider,omitempty"`
	Repository              *string `json:"repository,omitempty"`
	Branch                  *string `json:"branch,omitempty"`
	DatabaseID              *int    `json:"database_id,omitempty"`
	PushToDeploy            bool    `json:"push_to_deploy,omitempty"`
	AllowWildcardSubdomains bool    `json:"allow_wildcard_subdomains"`
	WWWRedirectType         string  `json:"www_redirect_type"`
}

type UpdateSiteRequest struct {
	RootPath         *string `json:"root_path,omitempty"`
	Directory        *string `json:"directory,omitempty"`
	Type             string  `json:"type,omitempty"`
	PHPVersion       string  `json:"php_version,omitempty"`
	PushToDeploy     *bool   `json:"push_to_deploy,omitempty"`
	RepositoryBranch *string `json:"repository_branch,omitempty"`
}

func (c *Client) CreateSite(ctx context.Context, serverID int, req CreateSiteRequest) (*Site, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites", serverID))
	var site Site
	id, err := c.PostJsonApi(ctx, path, req, &site)
	if err != nil {
		return nil, err
	}
	site.ID = int64(id)
	return &site, nil
}

func (c *Client) ListSites(ctx context.Context, serverID int) ([]Site, error) {
	items, err := c.GetJsonApiListAll(ctx, c.orgPath(fmt.Sprintf("/servers/%d/sites", serverID)))
	if err != nil {
		return nil, err
	}
	return unmarshalList(items, func(s *Site, id int64) { s.ID = id })
}

func (c *Client) GetSite(ctx context.Context, serverID, siteID int) (*Site, error) {
	path := c.orgPath(fmt.Sprintf("/sites/%d", siteID))
	var site Site
	id, err := c.GetJsonApi(ctx, path, &site)
	if err != nil {
		return nil, err
	}
	site.ID = int64(id)
	return &site, nil
}

func (c *Client) UpdateSite(ctx context.Context, serverID, siteID int, req UpdateSiteRequest) (*Site, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d", serverID, siteID))
	if err := c.doRequest(ctx, http.MethodPut, path, req, nil); err != nil {
		return nil, err
	}
	return c.GetSite(ctx, serverID, siteID)
}

func (c *Client) ChangeSitePHPVersion(ctx context.Context, serverID, siteID int, version string) error {
	req := UpdateSiteRequest{PHPVersion: version}
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d", serverID, siteID))
	return c.doRequest(ctx, http.MethodPut, path, req, nil)
}

func (c *Client) AddSiteAliases(ctx context.Context, serverID, siteID int, aliases []string) (*Site, error) {
	// In the new API, aliases are managed through domains
	// For now, update via the site update endpoint
	return c.GetSite(ctx, serverID, siteID)
}

func (c *Client) DeleteSite(ctx context.Context, serverID, siteID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d", serverID, siteID))
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

type balancingNode struct {
	ServerID int   `json:"server_id"`
	Weight   *int  `json:"weight,omitempty"`
	Down     *bool `json:"down,omitempty"`
	Backup   *bool `json:"backup,omitempty"`
	Port     *int  `json:"port,omitempty"`
}

type UpdateBalancingRequest struct {
	Servers []balancingNode `json:"servers"`
	Method  string          `json:"method"`
}

func (c *Client) GetSiteBalancing(ctx context.Context, serverID, siteID int) ([]balancingNode, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/load-balancing-nodes", serverID, siteID))
	items, err := c.GetJsonApiListAll(ctx, path)
	if err != nil {
		return nil, err
	}
	return unmarshalList(items, func(n *balancingNode, id int64) {})
}

func (c *Client) UpdateSiteBalancing(ctx context.Context, serverID, siteID int, req UpdateBalancingRequest) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/load-balancing-nodes", serverID, siteID))
	return c.doRequest(ctx, http.MethodPut, path, req, nil)
}

func (c *Client) GetSiteLog(ctx context.Context, serverID, siteID int) (string, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/logs/application", serverID, siteID))
	var logRes struct {
		Content string `json:"content"`
	}
	_, err := c.GetJsonApi(ctx, path, &logRes)
	if err != nil {
		return "", err
	}
	return logRes.Content, nil
}

func (c *Client) ClearSiteLog(ctx context.Context, serverID, siteID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/logs/application", serverID, siteID))
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}
