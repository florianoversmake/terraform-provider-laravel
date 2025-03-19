// Copyright (c) HashiCorp, Inc.

package forge_client

import (
	"context"
	"fmt"
	"net/http"
)

type Worker struct {
	ID                    int64   `json:"id"`
	Connection            string  `json:"connection"`
	Command               string  `json:"command"`
	Queue                 *string `json:"queue"`
	Timeout               int     `json:"timeout"`
	Delay                 int     `json:"delay"`
	Sleep                 int     `json:"sleep"`
	Tries                 *int    `json:"tries"`
	Processes             int     `json:"processes"`
	StopWaitSecs          *int    `json:"stopwaitsecs"`
	Environment           *string `json:"environment"`
	PHPVersion            string  `json:"php_version"`
	Daemon                bool    `json:"daemon"`
	Force                 bool    `json:"force"`
	Status                string  `json:"status"`
	CreatedAt             string  `json:"created_at"`
	DisplayablePHPVersion string  `json:"displayable_php_version"`
}

type workerResponse struct {
	Worker Worker `json:"worker"`
}

type workersResponse struct {
	Workers []Worker `json:"workers"`
}

type CreateWorkerRequest struct {
	Connection   string  `json:"connection"`
	TimeOut      int     `json:"timeout"`
	Delay        int     `json:"delay"`
	Sleep        int     `json:"sleep"`
	Tries        *int    `json:"tries"`
	Processes    int     `json:"processes"`
	StopWaitSecs *int    `json:"stopwaitsecs,omitempty"`
	Daemon       bool    `json:"daemon"`
	Force        bool    `json:"force"`
	PHPVersion   string  `json:"php_version"`
	Queue        *string `json:"queue"`
	Memory       int     `json:"memory"`
	Directory    string  `json:"directory"`
}

func (c *Client) ListWorkers(ctx context.Context, serverID int, siteID int) ([]Worker, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/workers", serverID, siteID)
	var res workersResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return res.Workers, nil
}

func (c *Client) CreateWorker(ctx context.Context, serverID int, siteID int, req CreateWorkerRequest) (*Worker, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/workers", serverID, siteID)
	var res workerResponse
	if err := c.doRequest(ctx, http.MethodPost, path, req, &res); err != nil {
		return nil, err
	}
	return &res.Worker, nil
}

type ErrorWorkerNotFound struct {
	ServerID int
	SiteID   int
	WorkerID int
}

func (e *ErrorWorkerNotFound) Error() string {
	return fmt.Sprintf("worker not found: server=%d, site=%d, worker=%d", e.ServerID, e.SiteID, e.WorkerID)
}

func (c *Client) GetWorker(ctx context.Context, serverID int, siteID int, workerID int) (*Worker, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/workers/%d", serverID, siteID, workerID)
	var res workerResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		if _, ok := err.(*ClientErrorResourceNotFound); ok {
			return nil, &ErrorWorkerNotFound{ServerID: serverID, SiteID: siteID, WorkerID: workerID}
		}
		return nil, err
	}

	return &res.Worker, nil
}

func (c *Client) DeleteWorker(ctx context.Context, serverID int, siteID int, workerID int) error {
	path := fmt.Sprintf("/servers/%d/sites/%d/workers/%d", serverID, siteID, workerID)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

func (c *Client) RestartWorker(ctx context.Context, serverID int, siteID int, workerID int) error {
	path := fmt.Sprintf("/servers/%d/sites/%d/workers/%d/restart", serverID, siteID, workerID)
	return c.doRequest(ctx, http.MethodPost, path, nil, nil)
}

type workerOutputResponse struct {
	Output string `json:"output"`
}

func (c *Client) GetWorkerOutput(ctx context.Context, serverID int, siteID int, workerID int) (string, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/workers/%d/output", serverID, siteID, workerID)
	var res workerOutputResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return "", err
	}
	return res.Output, nil
}
