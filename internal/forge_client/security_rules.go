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
	Status      string               `json:"status"`
	CreatedAt   string               `json:"created_at"`
	Credentials []SecurityCredential `json:"credentials"`
}

type SecurityRuleCredentialRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateSecurityRuleRequest struct {
	Name        string                          `json:"name"`
	Path        *string                         `json:"path,omitempty"`
	Credentials []SecurityRuleCredentialRequest `json:"credentials"`
}

func (c *Client) CreateSecurityRule(ctx context.Context, serverID, siteID int, req CreateSecurityRuleRequest) (*SecurityRule, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/security-rules", serverID, siteID))
	var rule SecurityRule
	id, err := c.PostJsonApi(ctx, path, req, &rule)
	if err != nil {
		return nil, err
	}
	rule.ID = int64(id)
	return &rule, nil
}

func (c *Client) ListSecurityRules(ctx context.Context, serverID, siteID int) ([]SecurityRule, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/security-rules", serverID, siteID))
	items, err := c.GetJsonApiListAll(ctx, path)
	if err != nil {
		return nil, err
	}
	return unmarshalList(items, func(r *SecurityRule, id int64) { r.ID = id })
}

func (c *Client) GetSecurityRule(ctx context.Context, serverID, siteID, ruleID int) (*SecurityRule, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/security-rules/%d", serverID, siteID, ruleID))
	var rule SecurityRule
	id, err := c.GetJsonApi(ctx, path, &rule)
	if err != nil {
		return nil, err
	}
	rule.ID = int64(id)
	return &rule, nil
}

func (c *Client) DeleteSecurityRule(ctx context.Context, serverID, siteID, ruleID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/security-rules/%d", serverID, siteID, ruleID))
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}
