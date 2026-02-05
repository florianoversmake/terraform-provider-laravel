package forge_client

import (
	"context"
	"fmt"
)

type ServerLog struct {
	ID      int64  `json:"id"`
	Path    string `json:"path"`
	Content string `json:"content"`
}

// GetServerLog retrieves a specific server log by key.
// Valid keys include: nginx_access, nginx_error, database, php, etc.
func (c *Client) GetServerLog(ctx context.Context, serverID int, key string) (*ServerLog, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/logs/%s", serverID, key))
	var res ServerLog
	_, err := c.GetJsonApi(ctx, path, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
