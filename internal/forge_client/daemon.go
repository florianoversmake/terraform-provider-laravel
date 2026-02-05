package forge_client

import (
	"context"
	"fmt"
	"net/http"
)

type Daemon struct {
	ID           int64  `json:"id"`
	Command      string `json:"command"`
	User         string `json:"user"`
	Directory    string `json:"directory"`
	Processes    int    `json:"processes"`
	StartSecs    int    `json:"startsecs"`
	StopWaitSecs int    `json:"stopwaitsecs"`
	StopSignal   string `json:"stopsignal"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
}

type CreateDaemonRequest struct {
	Command      string `json:"command"`
	User         string `json:"user"`
	Directory    string `json:"directory"`
	Processes    int    `json:"processes"`
	StartSecs    int    `json:"startsecs"`
	StopWaitSecs int    `json:"stopwaitsecs"`
	StopSignal   string `json:"stopsignal"`
}

func (c *Client) CreateDaemon(ctx context.Context, serverID int, req CreateDaemonRequest) (*Daemon, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/daemons", serverID))
	var daemon Daemon
	id, err := c.PostJsonApi(ctx, path, req, &daemon)
	if err != nil {
		return nil, err
	}
	daemon.ID = int64(id)
	return &daemon, nil
}

func (c *Client) ListDaemons(ctx context.Context, serverID int) ([]Daemon, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/daemons", serverID))
	items, err := c.GetJsonApiList(ctx, path)
	if err != nil {
		return nil, err
	}
	return unmarshalList(items, func(d *Daemon, id int64) { d.ID = id })
}

func (c *Client) GetDaemon(ctx context.Context, serverID, daemonID int) (*Daemon, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/daemons/%d", serverID, daemonID))
	var daemon Daemon
	id, err := c.GetJsonApi(ctx, path, &daemon)
	if err != nil {
		return nil, err
	}
	daemon.ID = int64(id)
	return &daemon, nil
}

func (c *Client) DeleteDaemon(ctx context.Context, serverID, daemonID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/daemons/%d", serverID, daemonID))
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

func (c *Client) RestartDaemon(ctx context.Context, serverID, daemonID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/daemons/%d/restart", serverID, daemonID))
	return c.doRequest(ctx, http.MethodPost, path, nil, nil)
}
