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

func (c *Client) CreateNginxTemplate(ctx context.Context, serverID int, name, content string) (*NginxTemplate, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/nginx/templates", serverID))
	req := map[string]string{"name": name, "content": content}
	var tmpl NginxTemplate
	id, err := c.PostJsonApi(ctx, path, req, &tmpl)
	if err != nil {
		return nil, err
	}
	tmpl.ID = int64(id)
	return &tmpl, nil
}

func (c *Client) ListNginxTemplates(ctx context.Context, serverID int) ([]NginxTemplate, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/nginx/templates/default", serverID))
	items, err := c.GetJsonApiList(ctx, path)
	if err != nil {
		return nil, err
	}
	return unmarshalList(items, func(t *NginxTemplate, id int64) { t.ID = id })
}

func (c *Client) GetNginxTemplate(ctx context.Context, serverID, templateID int) (*NginxTemplate, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/nginx/templates/%d", serverID, templateID))
	var tmpl NginxTemplate
	id, err := c.GetJsonApi(ctx, path, &tmpl)
	if err != nil {
		return nil, err
	}
	tmpl.ID = int64(id)
	return &tmpl, nil
}

func (c *Client) UpdateNginxTemplate(ctx context.Context, serverID, templateID int, name, content string) (*NginxTemplate, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/nginx/templates/%d", serverID, templateID))
	req := map[string]string{"name": name, "content": content}
	var tmpl NginxTemplate
	id, err := c.PutJsonApi(ctx, path, req, &tmpl)
	if err != nil {
		return nil, err
	}
	tmpl.ID = int64(id)
	return &tmpl, nil
}

func (c *Client) DeleteNginxTemplate(ctx context.Context, serverID, templateID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/nginx/templates/%d", serverID, templateID))
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}
