package forge_client

import (
	"context"
	"fmt"
	"net/http"
)

type RedirectRule struct {
	ID        int64  `json:"id"`
	From      string `json:"from"`
	To        string `json:"to"`
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
}

type CreateRedirectRuleRequest struct {
	From string `json:"from"`
	To   string `json:"to"`
	Type string `json:"type"`
}

func (c *Client) CreateRedirectRule(ctx context.Context, serverID, siteID int, req CreateRedirectRuleRequest) (*RedirectRule, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/redirect-rules", serverID, siteID))
	var rule RedirectRule
	id, err := c.PostJsonApi(ctx, path, req, &rule)
	if err != nil {
		return nil, err
	}
	rule.ID = int64(id)
	return &rule, nil
}

func (c *Client) ListRedirectRules(ctx context.Context, serverID, siteID int) ([]RedirectRule, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/redirect-rules", serverID, siteID))
	items, err := c.GetJsonApiListAll(ctx, path)
	if err != nil {
		return nil, err
	}
	return unmarshalList(items, func(r *RedirectRule, id int64) { r.ID = id })
}

func (c *Client) GetRedirectRule(ctx context.Context, serverID, siteID, ruleID int) (*RedirectRule, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/redirect-rules/%d", serverID, siteID, ruleID))
	var rule RedirectRule
	id, err := c.GetJsonApi(ctx, path, &rule)
	if err != nil {
		return nil, err
	}
	rule.ID = int64(id)
	return &rule, nil
}

func (c *Client) DeleteRedirectRule(ctx context.Context, serverID, siteID, ruleID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/redirect-rules/%d", serverID, siteID, ruleID))
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}
