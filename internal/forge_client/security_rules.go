package forge_client

import (
	"context"
	"fmt"
	"net/http"
)

type SecurityCredential struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
}

type SecurityRule struct {
	ID          int64                `json:"id"`
	Name        string               `json:"name"`
	Path        *string              `json:"path"`
	CreatedAt   string               `json:"created_at"`
	Credentials []SecurityCredential `json:"credentials"`
}

type securityRuleResponse struct {
	SecurityRule SecurityRule `json:"security_rule"`
}

type securityRulesResponse struct {
	SecurityRules []SecurityRule `json:"security_rules"`
}

type CreateSecurityRuleRequest struct {
	Name        string  `json:"name"`
	Path        *string `json:"path"`
	Credentials []struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"credentials"`
}

func (c *Client) CreateSecurityRule(ctx context.Context, serverID, siteID int, req CreateSecurityRuleRequest) (*SecurityRule, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/security-rules", serverID, siteID)
	var res securityRuleResponse
	if err := c.doRequest(ctx, http.MethodPost, path, req, &res); err != nil {
		return nil, err
	}
	return &res.SecurityRule, nil
}

func (c *Client) ListSecurityRules(ctx context.Context, serverID, siteID int) ([]SecurityRule, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/security-rules", serverID, siteID)
	var res securityRulesResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return res.SecurityRules, nil
}

func (c *Client) GetSecurityRule(ctx context.Context, serverID, siteID, ruleID int) (*SecurityRule, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/security-rules/%d", serverID, siteID, ruleID)
	var res securityRuleResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return &res.SecurityRule, nil
}

func (c *Client) DeleteSecurityRule(ctx context.Context, serverID, siteID, ruleID int) error {
	path := fmt.Sprintf("/servers/%d/sites/%d/security-rules/%d", serverID, siteID, ruleID)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}
