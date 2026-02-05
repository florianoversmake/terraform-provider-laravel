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

type CreateMonitorRequest struct {
	Type      string `json:"type"`
	Operator  string `json:"operator"`
	Threshold string `json:"threshold"`
	Minutes   string `json:"minutes"`
	Notify    string `json:"notify"`
}

func (c *Client) ListMonitors(ctx context.Context, serverID int) ([]Monitor, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/monitors", serverID))
	items, err := c.GetJsonApiList(ctx, path)
	if err != nil {
		return nil, err
	}
	return unmarshalList(items, func(m *Monitor, id int64) { m.ID = id })
}

func (c *Client) CreateMonitor(ctx context.Context, serverID int, req CreateMonitorRequest) (*Monitor, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/monitors", serverID))
	var monitor Monitor
	id, err := c.PostJsonApi(ctx, path, req, &monitor)
	if err != nil {
		return nil, err
	}
	monitor.ID = int64(id)
	return &monitor, nil
}

func (c *Client) GetMonitor(ctx context.Context, serverID, monitorID int) (*Monitor, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/monitors/%d", serverID, monitorID))
	var monitor Monitor
	id, err := c.GetJsonApi(ctx, path, &monitor)
	if err != nil {
		return nil, err
	}
	monitor.ID = int64(id)
	return &monitor, nil
}

func (c *Client) DeleteMonitor(ctx context.Context, serverID, monitorID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/monitors/%d", serverID, monitorID))
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}
