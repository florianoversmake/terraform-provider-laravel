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
	Template NginxTemplate `json:"template"`
}

type nginxTemplatesResponse struct {
	Templates []NginxTemplate `json:"templates"`
}

func (c *Client) CreateNginxTemplate(ctx context.Context, serverID int, name, content string) (*NginxTemplate, error) {
	path := fmt.Sprintf("/servers/%d/nginx/templates", serverID)
	req := map[string]string{"name": name, "content": content}
	var res nginxTemplateResponse
	if err := c.doRequest(ctx, http.MethodPost, path, req, &res); err != nil {
		return nil, err
	}
	return &res.Template, nil
}

func (c *Client) ListNginxTemplates(ctx context.Context, serverID int) ([]NginxTemplate, error) {
	path := fmt.Sprintf("/servers/%d/nginx/templates/default", serverID)
	var res nginxTemplatesResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return res.Templates, nil
}

func (c *Client) GetNginxTemplate(ctx context.Context, serverID, templateID int) (*NginxTemplate, error) {
	path := fmt.Sprintf("/servers/%d/nginx/templates/%d", serverID, templateID)
	var res nginxTemplateResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return &res.Template, nil
}

func (c *Client) UpdateNginxTemplate(ctx context.Context, serverID, templateID int, name, content string) (*NginxTemplate, error) {
	path := fmt.Sprintf("/servers/%d/nginx/templates/%d", serverID, templateID)
	req := map[string]string{"name": name, "content": content}
	var res nginxTemplateResponse
	if err := c.doRequest(ctx, http.MethodPut, path, req, &res); err != nil {
		return nil, err
	}
	return &res.Template, nil
}

func (c *Client) DeleteNginxTemplate(ctx context.Context, serverID, templateID int) error {
	path := fmt.Sprintf("/servers/%d/nginx/templates/%d", serverID, templateID)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}
