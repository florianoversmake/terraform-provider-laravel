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

type CreateWebhookRequest struct {
	URL string `json:"url"`
}

func (c *Client) ListWebhooks(ctx context.Context, serverID, siteID int) ([]Webhook, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/webhooks", serverID, siteID))
	items, err := c.GetJsonApiListAll(ctx, path)
	if err != nil {
		return nil, err
	}
	return unmarshalList(items, func(w *Webhook, id int64) { w.ID = id })
}

func (c *Client) GetWebhook(ctx context.Context, serverID, siteID, webhookID int) (*Webhook, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/webhooks/%d", serverID, siteID, webhookID))
	var webhook Webhook
	id, err := c.GetJsonApi(ctx, path, &webhook)
	if err != nil {
		return nil, err
	}
	webhook.ID = int64(id)
	return &webhook, nil
}

func (c *Client) CreateWebhook(ctx context.Context, serverID, siteID int, urlStr string) (*Webhook, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/webhooks", serverID, siteID))
	req := CreateWebhookRequest{URL: urlStr}
	var webhook Webhook
	id, err := c.PostJsonApi(ctx, path, req, &webhook)
	if err != nil {
		return nil, err
	}
	webhook.ID = int64(id)
	return &webhook, nil
}

func (c *Client) DeleteWebhook(ctx context.Context, serverID, siteID, webhookID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/webhooks/%d", serverID, siteID, webhookID))
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}
