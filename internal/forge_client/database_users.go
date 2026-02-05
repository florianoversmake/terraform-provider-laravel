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

func (c *Client) CreateDatabaseUser(ctx context.Context, serverID int, req CreateDatabaseUserRequest) (*DatabaseUser, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/database/users", serverID))
	var user DatabaseUser
	id, err := c.PostJsonApi(ctx, path, req, &user)
	if err != nil {
		return nil, err
	}
	user.ID = int64(id)
	return &user, nil
}

func (c *Client) ListDatabaseUsers(ctx context.Context, serverID int) ([]DatabaseUser, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/database/users", serverID))
	items, err := c.GetJsonApiListAll(ctx, path)
	if err != nil {
		return nil, err
	}
	return unmarshalList(items, func(u *DatabaseUser, id int64) { u.ID = id })
}

func (c *Client) GetDatabaseUser(ctx context.Context, serverID, userID int) (*DatabaseUser, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/database/users/%d", serverID, userID))
	var user DatabaseUser
	id, err := c.GetJsonApi(ctx, path, &user)
	if err != nil {
		return nil, err
	}
	user.ID = int64(id)
	return &user, nil
}

type UpdateDatabaseUserRequest struct {
	Databases []int64 `json:"databases"`
}

func (c *Client) UpdateDatabaseUser(ctx context.Context, serverID, userID int, req UpdateDatabaseUserRequest) (*DatabaseUser, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/database/users/%d", serverID, userID))
	var user DatabaseUser
	id, err := c.PutJsonApi(ctx, path, req, &user)
	if err != nil {
		return nil, err
	}
	user.ID = int64(id)
	return &user, nil
}

func (c *Client) DeleteDatabaseUser(ctx context.Context, serverID, userID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/database/users/%d", serverID, userID))
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}
