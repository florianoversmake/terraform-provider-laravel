package forge_client

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Client) CreateMySQLDatabase(ctx context.Context, serverID int, req CreateDatabaseRequest) (*Database, error) {
	path := fmt.Sprintf("/servers/%d/mysql", serverID)
	var res databaseResponse
	if err := c.doRequest(ctx, http.MethodPost, path, req, &res); err != nil {
		return nil, err
	}
	return &res.Database, nil
}

func (c *Client) ListMySQLDatabases(ctx context.Context, serverID int) ([]Database, error) {
	path := fmt.Sprintf("/servers/%d/mysql", serverID)
	var res databasesResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return res.Databases, nil
}

func (c *Client) GetMySQLDatabase(ctx context.Context, serverID, databaseID int) (*Database, error) {
	path := fmt.Sprintf("/servers/%d/mysql/%d", serverID, databaseID)
	var res databaseResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return &res.Database, nil
}

func (c *Client) DeleteMySQLDatabase(ctx context.Context, serverID, databaseID int) error {
	path := fmt.Sprintf("/servers/%d/mysql/%d", serverID, databaseID)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

func (c *Client) CreateMySQLDatabaseUser(ctx context.Context, serverID int, req CreateDatabaseUserRequest) (*DatabaseUser, error) {
	path := fmt.Sprintf("/servers/%d/mysql-users", serverID)
	var res databaseUserResponse
	if err := c.doRequest(ctx, http.MethodPost, path, req, &res); err != nil {
		return nil, err
	}
	return &res.User, nil
}

func (c *Client) ListMySQLDatabaseUsers(ctx context.Context, serverID int) ([]DatabaseUser, error) {
	path := fmt.Sprintf("/servers/%d/mysql-users", serverID)
	var res databaseUsersResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return res.Users, nil
}

func (c *Client) GetMySQLDatabaseUser(ctx context.Context, serverID, userID int) (*DatabaseUser, error) {
	path := fmt.Sprintf("/servers/%d/mysql-users/%d", serverID, userID)
	var res databaseUserResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return &res.User, nil
}

func (c *Client) UpdateMySQLDatabaseUser(ctx context.Context, serverID, userID int, req UpdateDatabaseUserRequest) (*DatabaseUser, error) {
	path := fmt.Sprintf("/servers/%d/mysql-users/%d", serverID, userID)
	var res databaseUserResponse
	if err := c.doRequest(ctx, http.MethodPut, path, req, &res); err != nil {
		return nil, err
	}
	return &res.User, nil
}

func (c *Client) DeleteMySQLDatabaseUser(ctx context.Context, serverID, userID int) error {
	path := fmt.Sprintf("/servers/%d/mysql-users/%d", serverID, userID)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}
