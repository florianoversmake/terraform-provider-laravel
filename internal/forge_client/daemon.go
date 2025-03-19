// Copyright (c) HashiCorp, Inc.

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

type DaemonResponse struct {
	Daemon Daemon `json:"daemon"`
}

func (c *Client) CreateDaemon(ctx context.Context, serverID int, req CreateDaemonRequest) (*Daemon, error) {
	path := fmt.Sprintf("/servers/%d/daemons", serverID)
	var resp DaemonResponse
	if err := c.doRequest(ctx, http.MethodPost, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp.Daemon, nil
}

type daemonsResponse struct {
	Daemons []Daemon `json:"daemons"`
}

func (c *Client) ListDaemons(ctx context.Context, serverID int) ([]Daemon, error) {
	path := fmt.Sprintf("/servers/%d/daemons", serverID)
	var resp daemonsResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return resp.Daemons, nil
}

func (c *Client) GetDaemon(ctx context.Context, serverID, daemonID int) (*Daemon, error) {
	path := fmt.Sprintf("/servers/%d/daemons/%d", serverID, daemonID)
	var resp DaemonResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Daemon, nil
}

func (c *Client) DeleteDaemon(ctx context.Context, serverID, daemonID int) error {
	path := fmt.Sprintf("/servers/%d/daemons/%d", serverID, daemonID)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

func (c *Client) RestartDaemon(ctx context.Context, serverID, daemonID int) error {
	path := fmt.Sprintf("/servers/%d/daemons/%d/restart", serverID, daemonID)
	return c.doRequest(ctx, http.MethodPost, path, nil, nil)
}
