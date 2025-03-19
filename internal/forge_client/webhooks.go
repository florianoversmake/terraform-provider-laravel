// Copyright (c) HashiCorp, Inc.

package forge_client

import (
	"context"
	"fmt"
	"net/http"
)

type Webhook struct {
	ID        int64  `json:"id"`
	URL       string `json:"url"`
	CreatedAt string `json:"created_at"`
}

type webhooksResponse struct {
	Webhooks []Webhook `json:"webhooks"`
}

type webhookResponse struct {
	Webhook Webhook `json:"webhook"`
}

func (c *Client) ListWebhooks(ctx context.Context, serverID, siteID int) ([]Webhook, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/webhooks", serverID, siteID)
	var res webhooksResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return res.Webhooks, nil
}

func (c *Client) GetWebhook(ctx context.Context, serverID, siteID, webhookID int) (*Webhook, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/webhooks/%d", serverID, siteID, webhookID)
	var res webhookResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return &res.Webhook, nil
}

func (c *Client) CreateWebhook(ctx context.Context, serverID, siteID int, urlStr string) (string, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/webhooks", serverID, siteID)
	req := map[string]string{"url": urlStr}
	var res map[string]string
	if err := c.doRequest(ctx, http.MethodPost, path, req, &res); err != nil {
		return "", err
	}
	return res["url"], nil
}

func (c *Client) DeleteWebhook(ctx context.Context, serverID, siteID, webhookID int) (string, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/webhooks/%d", serverID, siteID, webhookID)
	var res map[string]string
	if err := c.doRequest(ctx, http.MethodDelete, path, nil, &res); err != nil {
		return "", err
	}
	return res["url"], nil
}
