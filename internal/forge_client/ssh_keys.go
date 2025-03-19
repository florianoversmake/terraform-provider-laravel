// Copyright (c) HashiCorp, Inc.

package forge_client

import (
	"context"
	"fmt"
	"net/http"
)

type SSHKey struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

type sshKeyResponse struct {
	Key SSHKey `json:"key"`
}

type sshKeysResponse struct {
	Keys []SSHKey `json:"keys"`
}

type CreateSSHKeyRequest struct {
	Name     string `json:"name"`
	Key      string `json:"key"`
	Username string `json:"username"`
}

func (c *Client) CreateSSHKey(ctx context.Context, serverID int, req CreateSSHKeyRequest) (*SSHKey, error) {
	path := fmt.Sprintf("/servers/%d/keys", serverID)
	var res sshKeyResponse
	if err := c.doRequest(ctx, http.MethodPost, path, req, &res); err != nil {
		return nil, err
	}
	return &res.Key, nil
}

func (c *Client) ListSSHKeys(ctx context.Context, serverID int) ([]SSHKey, error) {
	path := fmt.Sprintf("/servers/%d/keys", serverID)
	var res sshKeysResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return res.Keys, nil
}

func (c *Client) GetSSHKey(ctx context.Context, serverID, keyID int) (*SSHKey, error) {
	path := fmt.Sprintf("/servers/%d/keys/%d", serverID, keyID)
	var res sshKeyResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return &res.Key, nil
}

func (c *Client) DeleteSSHKey(ctx context.Context, serverID, keyID int) error {
	path := fmt.Sprintf("/servers/%d/keys/%d", serverID, keyID)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}
