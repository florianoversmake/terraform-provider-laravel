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

type redirectRuleResponse struct {
	RedirectRule RedirectRule `json:"redirect_rule"`
}

type redirectRulesResponse struct {
	RedirectRules []RedirectRule `json:"redirect_rules"`
}

type CreateRedirectRuleRequest struct {
	From string `json:"from"`
	To   string `json:"to"`
	Type string `json:"type"`
}

func (c *Client) CreateRedirectRule(ctx context.Context, serverID, siteID int, req CreateRedirectRuleRequest) (*RedirectRule, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/redirect-rules", serverID, siteID)
	var res redirectRuleResponse
	if err := c.doRequest(ctx, http.MethodPost, path, req, &res); err != nil {
		return nil, err
	}
	return &res.RedirectRule, nil
}

func (c *Client) ListRedirectRules(ctx context.Context, serverID, siteID int) ([]RedirectRule, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/redirect-rules", serverID, siteID)
	var res redirectRulesResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return res.RedirectRules, nil
}

func (c *Client) GetRedirectRule(ctx context.Context, serverID, siteID, ruleID int) (*RedirectRule, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/redirect-rules/%d", serverID, siteID, ruleID)
	var res redirectRuleResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return &res.RedirectRule, nil
}

func (c *Client) DeleteRedirectRule(ctx context.Context, serverID, siteID, ruleID int) error {
	path := fmt.Sprintf("/servers/%d/sites/%d/redirect-rules/%d", serverID, siteID, ruleID)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}
