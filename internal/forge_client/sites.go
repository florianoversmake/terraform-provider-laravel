package forge_client

import (
	"context"
	"fmt"
	"net/http"
)

type Site struct {
	ID                 int64    `json:"id"`
	Name               string   `json:"name"`
	Aliases            []string `json:"aliases"`
	Directory          string   `json:"directory"`
	Wildcards          bool     `json:"wildcards"`
	Isolated           bool     `json:"isolated"`
	Username           string   `json:"username"`
	Status             string   `json:"status"`
	Repository         *string  `json:"repository"`
	RepositoryProvider *string  `json:"repository_provider"`
	RepositoryBranch   *string  `json:"repository_branch"`
	RepositoryStatus   *string  `json:"repository_status"`
	QuickDeploy        bool     `json:"quick_deploy"`
	ProjectType        string   `json:"project_type"`
	App                *string  `json:"app"`
	PHPVersion         string   `json:"php_version"`
	AppStatus          *string  `json:"app_status"`
	SlackChannel       *string  `json:"slack_channel"`
	TelegramChatID     *string  `json:"telegram_chat_id"`
	TelegramChatTitle  *string  `json:"telegram_chat_title"`
	DeploymentURL      *string  `json:"deployment_url"`
	CreatedAt          string   `json:"created_at"`
	Tags               []string `json:"tags"`
	WebDirectory       string   `json:"web_directory"`
}

type CreateSiteRequest struct {
	Domain        string      `json:"domain"`
	ProjectType   string      `json:"project_type"`
	Aliases       []string    `json:"aliases"`
	Directory     string      `json:"directory"`
	Isolated      bool        `json:"isolated"`
	Username      string      `json:"username"`
	Database      string      `json:"database"`
	PHPVersion    string      `json:"php_version"`
	NginxTemplate interface{} `json:"nginx_template,omitempty"`
}

type siteResponse struct {
	Data Site `json:"data"`
}

type sitesResponse struct {
	Data []Site `json:"data"`
}

func (c *Client) CreateSite(ctx context.Context, serverID int, req CreateSiteRequest) (*Site, error) {
	path := fmt.Sprintf("/orgs/%s/servers/%d/sites", c.OrgSlug, serverID)
	var res siteResponse
	if err := c.doRequest(ctx, http.MethodPost, path, req, &res); err != nil {
		return nil, err
	}
	return &res.Data, nil
}

func (c *Client) ListSites(ctx context.Context, serverID int) ([]Site, error) {
	path := fmt.Sprintf("/orgs/%s/servers/%d/sites", c.OrgSlug, serverID)
	var res sitesResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return res.Data, nil
}

func (c *Client) GetSite(ctx context.Context, serverID, siteID int) (*Site, error) {
	path := fmt.Sprintf("/orgs/%s/servers/%d/sites/%d", c.OrgSlug, serverID, siteID)
	var res siteResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return &res.Data, nil
}

type UpdateSiteRequest struct {
	Directory  string   `json:"directory"`
	Name       string   `json:"name"`
	PHPVersion string   `json:"php_version"`
	Aliases    []string `json:"aliases"`
	Wildcards  bool     `json:"wildcards"`
}

func (c *Client) UpdateSite(ctx context.Context, serverID, siteID int, req UpdateSiteRequest) (*Site, error) {
	path := fmt.Sprintf("/orgs/%s/servers/%d/sites/%d", c.OrgSlug, serverID, siteID)
	var res siteResponse
	if err := c.doRequest(ctx, http.MethodPut, path, req, &res); err != nil {
		return nil, err
	}
	return &res.Data, nil
}

type ChangeSitePHPVersionRequest struct {
	Version string `json:"version"`
}

func (c *Client) ChangeSitePHPVersion(ctx context.Context, serverID, siteID int, version string) error {
	path := fmt.Sprintf("/orgs/%s/servers/%d/sites/%d/php", c.OrgSlug, serverID, siteID)
	req := ChangeSitePHPVersionRequest{Version: version}
	return c.doRequest(ctx, http.MethodPut, path, req, nil)
}

type AddSiteAliasesRequest struct {
	Aliases []string `json:"aliases"`
}

func (c *Client) AddSiteAliases(ctx context.Context, serverID, siteID int, aliases []string) (*Site, error) {
	path := fmt.Sprintf("/orgs/%s/servers/%d/sites/%d/aliases", c.OrgSlug, serverID, siteID)
	req := AddSiteAliasesRequest{Aliases: aliases}
	var res siteResponse
	if err := c.doRequest(ctx, http.MethodPut, path, req, &res); err != nil {
		return nil, err
	}
	return &res.Data, nil
}

func (c *Client) DeleteSite(ctx context.Context, serverID, siteID int) error {
	path := fmt.Sprintf("/orgs/%s/servers/%d/sites/%d", c.OrgSlug, serverID, siteID)
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

type balancingResponse struct {
	Data []balancingNode `json:"data"`
}

func (c *Client) GetSiteBalancing(ctx context.Context, serverID, siteID int) ([]balancingNode, error) {
	path := fmt.Sprintf("/orgs/%s/servers/%d/sites/%d/balancing", c.OrgSlug, serverID, siteID)
	var res balancingResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return res.Data, nil
}

func (c *Client) UpdateSiteBalancing(ctx context.Context, serverID, siteID int, req UpdateBalancingRequest) error {
	path := fmt.Sprintf("/orgs/%s/servers/%d/sites/%d/balancing", c.OrgSlug, serverID, siteID)
	return c.doRequest(ctx, http.MethodPut, path, req, nil)
}

type siteLogResponse struct {
	Data struct {
		Content string `json:"content"`
	} `json:"data"`
}

func (c *Client) GetSiteLog(ctx context.Context, serverID, siteID int) (string, error) {
	path := fmt.Sprintf("/orgs/%s/servers/%d/sites/%d/logs", c.OrgSlug, serverID, siteID)
	var res siteLogResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return "", err
	}
	return res.Data.Content, nil
}

func (c *Client) ClearSiteLog(ctx context.Context, serverID, siteID int) error {
	path := fmt.Sprintf("/orgs/%s/servers/%d/sites/%d/logs", c.OrgSlug, serverID, siteID)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}
