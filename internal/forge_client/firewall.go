package forge_client

import (
	"context"
	"fmt"
	"net/http"
)

type FirewallRule struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Port      string  `json:"port"`
	Type      string  `json:"type"`
	IpAddress *string `json:"ip_address"`
	Status    string  `json:"status"`
	CreatedAt string  `json:"created_at"`
}

type CreateFirewallRuleRequest struct {
	Name      string  `json:"name"`
	IpAddress *string `json:"ip_address"`
	Port      string  `json:"port"`
	Type      string  `json:"type"`
}

func (c *Client) CreateFirewallRule(ctx context.Context, serverID int, req CreateFirewallRuleRequest) (*FirewallRule, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/firewall-rules", serverID))
	var rule FirewallRule
	id, err := c.PostJsonApi(ctx, path, req, &rule)
	if err != nil {
		return nil, err
	}
	rule.ID = int64(id)
	return &rule, nil
}

func (c *Client) ListFirewallRules(ctx context.Context, serverID int) ([]FirewallRule, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/firewall-rules", serverID))
	items, err := c.GetJsonApiList(ctx, path)
	if err != nil {
		return nil, err
	}
	return unmarshalList(items, func(r *FirewallRule, id int64) { r.ID = id })
}

func (c *Client) GetFirewallRule(ctx context.Context, serverID, ruleID int) (*FirewallRule, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/firewall-rules/%d", serverID, ruleID))
	var rule FirewallRule
	id, err := c.GetJsonApi(ctx, path, &rule)
	if err != nil {
		return nil, err
	}
	rule.ID = int64(id)
	return &rule, nil
}

func (c *Client) DeleteFirewallRule(ctx context.Context, serverID, ruleID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/firewall-rules/%d", serverID, ruleID))
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}
