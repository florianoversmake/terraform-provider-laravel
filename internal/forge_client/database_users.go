// Copyright (c) HashiCorp, Inc.

package forge_client

import (
	"context"
	"fmt"
	"net/http"
)

type DatabaseUser struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Status    string  `json:"status"`
	CreatedAt string  `json:"created_at"`
	Databases []int64 `json:"databases"`
}

type CreateDatabaseUserRequest struct {
	Name      string  `json:"name"`
	Password  string  `json:"password"`
	Databases []int64 `json:"databases"`
}

type databaseUserResponse struct {
	User DatabaseUser `json:"user"`
}

type databaseUsersResponse struct {
	Users []DatabaseUser `json:"users"`
}

func (c *Client) CreateDatabaseUser(ctx context.Context, serverID int, req CreateDatabaseUserRequest) (*DatabaseUser, error) {
	path := fmt.Sprintf("/servers/%d/database-users", serverID)
	var res databaseUserResponse
	if err := c.doRequest(ctx, http.MethodPost, path, req, &res); err != nil {
		return nil, err
	}
	return &res.User, nil
}

func (c *Client) ListDatabaseUsers(ctx context.Context, serverID int) ([]DatabaseUser, error) {
	path := fmt.Sprintf("/servers/%d/database-users", serverID)
	var res databaseUsersResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return res.Users, nil
}

func (c *Client) GetDatabaseUser(ctx context.Context, serverID, userID int) (*DatabaseUser, error) {
	path := fmt.Sprintf("/servers/%d/database-users/%d", serverID, userID)
	var res databaseUserResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return &res.User, nil
}

type UpdateDatabaseUserRequest struct {
	Databases []int64 `json:"databases"`
}

func (c *Client) UpdateDatabaseUser(ctx context.Context, serverID, userID int, req UpdateDatabaseUserRequest) (*DatabaseUser, error) {
	path := fmt.Sprintf("/servers/%d/database-users/%d", serverID, userID)
	var res databaseUserResponse
	if err := c.doRequest(ctx, http.MethodPut, path, req, &res); err != nil {
		return nil, err
	}
	return &res.User, nil
}

func (c *Client) DeleteDatabaseUser(ctx context.Context, serverID, userID int) error {
	path := fmt.Sprintf("/servers/%d/database-users/%d", serverID, userID)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}
