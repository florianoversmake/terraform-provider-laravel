// Copyright (c) HashiCorp, Inc.

package forge_client

import (
	"context"
	"fmt"
	"net/http"
)

type Monitor struct {
	ID             int64  `json:"id"`
	Status         string `json:"status"`
	Type           string `json:"type"`
	Operator       string `json:"operator"`
	Threshold      int    `json:"threshold"`
	Minutes        int    `json:"minutes"`
	State          string `json:"state"`
	StateChangedAt string `json:"state_changed_at"`
}

type monitorsResponse struct {
	Monitors []Monitor `json:"monitors"`
}

type CreateMonitorRequest struct {
	Type      string `json:"type"`
	Operator  string `json:"operator"`
	Threshold string `json:"threshold"`
	Minutes   string `json:"minutes"`
	Notify    string `json:"notify"`
}

func (c *Client) ListMonitors(ctx context.Context, serverID int) ([]Monitor, error) {
	path := fmt.Sprintf("/servers/%d/monitors", serverID)
	var res monitorsResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return res.Monitors, nil
}

func (c *Client) CreateMonitor(ctx context.Context, serverID int, req CreateMonitorRequest) (*Monitor, error) {
	path := fmt.Sprintf("/servers/%d/monitors", serverID)
	var res struct {
		Monitor Monitor `json:"monitor"`
	}
	if err := c.doRequest(ctx, http.MethodPost, path, req, &res); err != nil {
		return nil, err
	}
	return &res.Monitor, nil
}

func (c *Client) GetMonitor(ctx context.Context, serverID, monitorID int) (*Monitor, error) {
	path := fmt.Sprintf("/servers/%d/monitors/%d", serverID, monitorID)
	var res struct {
		Monitor Monitor `json:"monitor"`
	}
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return &res.Monitor, nil
}

func (c *Client) DeleteMonitor(ctx context.Context, serverID, monitorID int) error {
	path := fmt.Sprintf("/servers/%d/monitors/%d", serverID, monitorID)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}
