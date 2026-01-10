package forge_client

import (
	"context"
	"fmt"
	"net/http"
)

type NginxTemplate struct {
	ID       int64  `json:"id"`
	ServerID int64  `json:"server_id"`
	Name     string `json:"name"`
	Content  string `json:"content"`
}

type nginxTemplateResponse struct {
	Data NginxTemplate `json:"data"`
}

type nginxTemplatesResponse struct {
	Data []NginxTemplate `json:"data"`
}

func (c *Client) CreateNginxTemplate(ctx context.Context, serverID int, name, content string) (*NginxTemplate, error) {
	path := fmt.Sprintf("/orgs/%s/servers/%d/nginx/templates", c.OrgSlug, serverID)
	req := map[string]string{"name": name, "content": content}
	var res nginxTemplateResponse
	if err := c.doRequest(ctx, http.MethodPost, path, req, &res); err != nil {
		return nil, err
	}
	return &res.Data, nil
}

func (c *Client) ListNginxTemplates(ctx context.Context, serverID int) ([]NginxTemplate, error) {
	path := fmt.Sprintf("/orgs/%s/servers/%d/nginx/templates/default", c.OrgSlug, serverID)
	var res nginxTemplatesResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return res.Data, nil
}

func (c *Client) GetNginxTemplate(ctx context.Context, serverID, templateID int) (*NginxTemplate, error) {
	path := fmt.Sprintf("/orgs/%s/servers/%d/nginx/templates/%d", c.OrgSlug, serverID, templateID)
	var res nginxTemplateResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return &res.Data, nil
}

func (c *Client) UpdateNginxTemplate(ctx context.Context, serverID, templateID int, name, content string) (*NginxTemplate, error) {
	path := fmt.Sprintf("/orgs/%s/servers/%d/nginx/templates/%d", c.OrgSlug, serverID, templateID)
	req := map[string]string{"name": name, "content": content}
	var res nginxTemplateResponse
	if err := c.doRequest(ctx, http.MethodPut, path, req, &res); err != nil {
		return nil, err
	}
	return &res.Data, nil
}

func (c *Client) DeleteNginxTemplate(ctx context.Context, serverID, templateID int) error {
	path := fmt.Sprintf("/orgs/%s/servers/%d/nginx/templates/%d", c.OrgSlug, serverID, templateID)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}
