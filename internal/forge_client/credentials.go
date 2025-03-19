// Copyright (c) HashiCorp, Inc.

package forge_client

import (
	"context"
	"net/http"
)

type Credential struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}

type credentialsResponse struct {
	Credentials []Credential `json:"credentials"`
}

func (c *Client) ListCredentials(ctx context.Context) ([]Credential, error) {
	var res credentialsResponse
	if err := c.doRequest(ctx, http.MethodGet, "/credentials", nil, &res); err != nil {
		return nil, err
	}
	return res.Credentials, nil
}
