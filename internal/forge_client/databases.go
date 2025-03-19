// Copyright (c) HashiCorp, Inc.

package forge_client

import (
	"context"
	"fmt"
	"net/http"
)

type Database struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

type CreateDatabaseRequest struct {
	Name     string `json:"name"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
}

type databaseResponse struct {
	Database Database `json:"database"`
}

type databasesResponse struct {
	Databases []Database `json:"databases"`
}

func (c *Client) CreateDatabase(ctx context.Context, serverID int, req CreateDatabaseRequest) (*Database, error) {
	path := fmt.Sprintf("/servers/%d/databases", serverID)
	var res databaseResponse
	if err := c.doRequest(ctx, http.MethodPost, path, req, &res); err != nil {
		return nil, err
	}
	return &res.Database, nil
}

func (c *Client) SyncDatabase(ctx context.Context, serverID int) error {
	path := fmt.Sprintf("/servers/%d/databases/sync", serverID)
	return c.doRequest(ctx, http.MethodPost, path, nil, nil)
}

func (c *Client) ListDatabases(ctx context.Context, serverID int) ([]Database, error) {
	path := fmt.Sprintf("/servers/%d/databases", serverID)
	var res databasesResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return res.Databases, nil
}

func (c *Client) GetDatabase(ctx context.Context, serverID, databaseID int) (*Database, error) {
	path := fmt.Sprintf("/servers/%d/databases/%d", serverID, databaseID)
	var res databaseResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return &res.Database, nil
}

func (c *Client) DeleteDatabase(ctx context.Context, serverID, databaseID int) error {
	path := fmt.Sprintf("/servers/%d/databases/%d", serverID, databaseID)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}
