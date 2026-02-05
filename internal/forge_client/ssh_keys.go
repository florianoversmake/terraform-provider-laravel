package forge_client

import (
	"context"
	"fmt"
	"net/http"
)

type SSHKey struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	User      string `json:"user"`
	Status    string `json:"status"`
	CreatedBy *int64 `json:"created_by"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CreateSSHKeyRequest struct {
	Name string  `json:"name"`
	Key  string  `json:"key"`
	User *string `json:"user,omitempty"`
}

func (c *Client) CreateSSHKey(ctx context.Context, serverID int, req CreateSSHKeyRequest) (*SSHKey, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/ssh-keys", serverID))
	var key SSHKey
	id, err := c.PostJsonApi(ctx, path, req, &key)
	if err != nil {
		return nil, err
	}
	key.ID = int64(id)
	return &key, nil
}

func (c *Client) ListSSHKeys(ctx context.Context, serverID int) ([]SSHKey, error) {
	items, err := c.GetJsonApiList(ctx, c.orgPath(fmt.Sprintf("/servers/%d/ssh-keys", serverID)))
	if err != nil {
		return nil, err
	}
	return unmarshalList(items, func(k *SSHKey, id int64) { k.ID = id })
}

func (c *Client) GetSSHKey(ctx context.Context, serverID, keyID int) (*SSHKey, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/ssh-keys/%d", serverID, keyID))
	var key SSHKey
	id, err := c.GetJsonApi(ctx, path, &key)
	if err != nil {
		return nil, err
	}
	key.ID = int64(id)
	return &key, nil
}

func (c *Client) DeleteSSHKey(ctx context.Context, serverID, keyID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/ssh-keys/%d", serverID, keyID))
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}
