package forge_client

import (
	"context"
	"fmt"
	"net/http"
)

type ServerLog struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

func (c *Client) GetServerLog(ctx context.Context, serverID int) (*ServerLog, error) {
	path := fmt.Sprintf("/servers/%d/logs", serverID)
	var res ServerLog
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
