package forge_client

import (
	"context"
	"net/http"
)

type Organization struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Type        string `json:"type"`
	Avatar      string `json:"avatar"`
	TotalSeats  int    `json:"total_seats"`
	UsedSeats   int    `json:"used_seats"`
	CanManage   bool   `json:"can_manage"`
	BillingLink string `json:"billing_link"`
	CreatedAt   string `json:"created_at"`
}

type organizationsResponse struct {
	Data []Organization `json:"data"`
}

func (c *Client) ListOrganizations(ctx context.Context) ([]Organization, error) {
	var res organizationsResponse
	if err := c.doRequest(ctx, http.MethodGet, "/orgs", nil, &res); err != nil {
		return nil, err
	}
	return res.Data, nil
}
