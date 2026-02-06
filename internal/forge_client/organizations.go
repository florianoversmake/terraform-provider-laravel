package forge_client

import (
	"context"
	"fmt"
)

// Organization represents a Forge organization.
type Organization struct {
	ID        int64   `json:"id"`
	Slug      string  `json:"slug"`
	Name      string  `json:"name"`
	Owner     bool    `json:"owner"`
	CreatedAt *string `json:"created_at"`
	UpdatedAt *string `json:"updated_at"`
}

// ListOrganizations returns all organizations the user has access to.
func (c *Client) ListOrganizations(ctx context.Context) ([]Organization, error) {
	items, err := c.GetJsonApiListAll(ctx, "/orgs")
	if err != nil {
		return nil, err
	}
	return unmarshalList(items, func(o *Organization, id int64) { o.ID = id })
}

// GetOrganization retrieves a specific organization by slug.
func (c *Client) GetOrganization(ctx context.Context, slug string) (*Organization, error) {
	path := fmt.Sprintf("/orgs/%s", slug)
	var org Organization
	id, err := c.GetJsonApi(ctx, path, &org)
	if err != nil {
		return nil, err
	}
	org.ID = int64(id)
	return &org, nil
}
