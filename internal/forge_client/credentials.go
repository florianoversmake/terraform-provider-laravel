package forge_client

import (
	"context"
)

type Credential struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Provider  string  `json:"provider"`
	InUse     bool    `json:"in_use"`
	CreatedAt *string `json:"created_at"`
	UpdatedAt *string `json:"updated_at"`
}

func (c *Client) ListCredentials(ctx context.Context) ([]Credential, error) {
	items, err := c.GetJsonApiList(ctx, c.orgPath("/server-credentials"))
	if err != nil {
		return nil, err
	}
	return unmarshalList(items, func(cr *Credential, id int64) { cr.ID = id })
}
