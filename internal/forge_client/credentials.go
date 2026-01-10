package forge_client

import (
	"context"
	"fmt"
	"net/http"
)

type Credential struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}

type credentialsResponse struct {
	Data []Credential `json:"data"`
}

func (c *Client) ListCredentials(ctx context.Context) ([]Credential, error) {
	path := fmt.Sprintf("/orgs/%s/server-credentials", c.OrgSlug)
	var res credentialsResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return res.Data, nil
}
