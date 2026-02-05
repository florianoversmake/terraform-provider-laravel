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

func (c *Client) CreateDatabase(ctx context.Context, serverID int, req CreateDatabaseRequest) (*Database, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/database/schemas", serverID))
	var db Database
	id, err := c.PostJsonApi(ctx, path, req, &db)
	if err != nil {
		return nil, err
	}
	db.ID = int64(id)
	return &db, nil
}

func (c *Client) SyncDatabase(ctx context.Context, serverID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/database/schemas/sync", serverID))
	return c.doRequest(ctx, http.MethodPost, path, nil, nil)
}

func (c *Client) ListDatabases(ctx context.Context, serverID int) ([]Database, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/database/schemas", serverID))
	items, err := c.GetJsonApiList(ctx, path)
	if err != nil {
		return nil, err
	}
	return unmarshalList(items, func(d *Database, id int64) { d.ID = id })
}

func (c *Client) GetDatabase(ctx context.Context, serverID, databaseID int) (*Database, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/database/schemas/%d", serverID, databaseID))
	var db Database
	id, err := c.GetJsonApi(ctx, path, &db)
	if err != nil {
		return nil, err
	}
	db.ID = int64(id)
	return &db, nil
}

func (c *Client) DeleteDatabase(ctx context.Context, serverID, databaseID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/database/schemas/%d", serverID, databaseID))
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}
