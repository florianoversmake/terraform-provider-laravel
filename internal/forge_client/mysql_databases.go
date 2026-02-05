package forge_client

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Client) CreateMySQLDatabase(ctx context.Context, serverID int, req CreateDatabaseRequest) (*Database, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/mysql", serverID))
	var db Database
	id, err := c.PostJsonApi(ctx, path, req, &db)
	if err != nil {
		return nil, err
	}
	db.ID = int64(id)
	return &db, nil
}

func (c *Client) ListMySQLDatabases(ctx context.Context, serverID int) ([]Database, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/mysql", serverID))
	items, err := c.GetJsonApiListAll(ctx, path)
	if err != nil {
		return nil, err
	}
	return unmarshalList(items, func(d *Database, id int64) { d.ID = id })
}

func (c *Client) GetMySQLDatabase(ctx context.Context, serverID, databaseID int) (*Database, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/mysql/%d", serverID, databaseID))
	var db Database
	id, err := c.GetJsonApi(ctx, path, &db)
	if err != nil {
		return nil, err
	}
	db.ID = int64(id)
	return &db, nil
}

func (c *Client) DeleteMySQLDatabase(ctx context.Context, serverID, databaseID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/mysql/%d", serverID, databaseID))
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

func (c *Client) CreateMySQLDatabaseUser(ctx context.Context, serverID int, req CreateDatabaseUserRequest) (*DatabaseUser, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/mysql-users", serverID))
	var user DatabaseUser
	id, err := c.PostJsonApi(ctx, path, req, &user)
	if err != nil {
		return nil, err
	}
	user.ID = int64(id)
	return &user, nil
}

func (c *Client) ListMySQLDatabaseUsers(ctx context.Context, serverID int) ([]DatabaseUser, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/mysql-users", serverID))
	items, err := c.GetJsonApiListAll(ctx, path)
	if err != nil {
		return nil, err
	}
	return unmarshalList(items, func(u *DatabaseUser, id int64) { u.ID = id })
}

func (c *Client) GetMySQLDatabaseUser(ctx context.Context, serverID, userID int) (*DatabaseUser, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/mysql-users/%d", serverID, userID))
	var user DatabaseUser
	id, err := c.GetJsonApi(ctx, path, &user)
	if err != nil {
		return nil, err
	}
	user.ID = int64(id)
	return &user, nil
}

func (c *Client) UpdateMySQLDatabaseUser(ctx context.Context, serverID, userID int, req UpdateDatabaseUserRequest) (*DatabaseUser, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/mysql-users/%d", serverID, userID))
	var user DatabaseUser
	id, err := c.PutJsonApi(ctx, path, req, &user)
	if err != nil {
		return nil, err
	}
	user.ID = int64(id)
	return &user, nil
}

func (c *Client) DeleteMySQLDatabaseUser(ctx context.Context, serverID, userID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/mysql-users/%d", serverID, userID))
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}
