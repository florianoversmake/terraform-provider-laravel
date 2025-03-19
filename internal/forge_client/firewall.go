// Copyright (c) HashiCorp, Inc.

package forge_client

import (
	"context"
	"fmt"
	"net/http"
)

type FirewallRule struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Port      int     `json:"port"`
	Type      string  `json:"type"`
	IpAddress *string `json:"ip_address"`
	Status    string  `json:"status"`
	CreatedAt string  `json:"created_at"`
}

type CreateFirewallRuleRequest struct {
	Name      string  `json:"name"`
	IpAddress *string `json:"ip_address"`
	Port      int     `json:"port"`
	Type      string  `json:"type"`
}

type FirewallRuleResponse struct {
	Rule FirewallRule `json:"rule"`
}

func (c *Client) CreateFirewallRule(ctx context.Context, serverID int, req CreateFirewallRuleRequest) (*FirewallRule, error) {
	path := fmt.Sprintf("/servers/%d/firewall-rules", serverID)
	var resp FirewallRuleResponse
	if err := c.doRequest(ctx, http.MethodPost, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp.Rule, nil
}

type firewallRulesResponse struct {
	Rules []FirewallRule `json:"rules"`
}

func (c *Client) ListFirewallRules(ctx context.Context, serverID int) ([]FirewallRule, error) {
	path := fmt.Sprintf("/servers/%d/firewall-rules", serverID)
	var resp firewallRulesResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return resp.Rules, nil
}

func (c *Client) GetFirewallRule(ctx context.Context, serverID, ruleID int) (*FirewallRule, error) {
	path := fmt.Sprintf("/servers/%d/firewall-rules/%d", serverID, ruleID)
	var resp FirewallRuleResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Rule, nil
}

func (c *Client) DeleteFirewallRule(ctx context.Context, serverID, ruleID int) error {
	path := fmt.Sprintf("/servers/%d/firewall-rules/%d", serverID, ruleID)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}
