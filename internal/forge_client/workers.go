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
	Data Worker `json:"data"`
}

type workersResponse struct {
	Data []Worker `json:"data"`
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
	path := fmt.Sprintf("/orgs/%s/servers/%d/sites/%d/workers", c.OrgSlug, serverID, siteID)
	var res workersResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return res.Data, nil
}

func (c *Client) CreateWorker(ctx context.Context, serverID int, siteID int, req CreateWorkerRequest) (*Worker, error) {
	path := fmt.Sprintf("/orgs/%s/servers/%d/sites/%d/workers", c.OrgSlug, serverID, siteID)
	var res workerResponse
	if err := c.doRequest(ctx, http.MethodPost, path, req, &res); err != nil {
		return nil, err
	}
	return &res.Data, nil
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
	path := fmt.Sprintf("/orgs/%s/servers/%d/sites/%d/workers/%d", c.OrgSlug, serverID, siteID, workerID)
	var res workerResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		if _, ok := err.(*ClientErrorResourceNotFound); ok {
			return nil, &ErrorWorkerNotFound{ServerID: serverID, SiteID: siteID, WorkerID: workerID}
		}
		return nil, err
	}

	return &res.Data, nil
}

func (c *Client) DeleteWorker(ctx context.Context, serverID int, siteID int, workerID int) error {
	path := fmt.Sprintf("/orgs/%s/servers/%d/sites/%d/workers/%d", c.OrgSlug, serverID, siteID, workerID)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

func (c *Client) RestartWorker(ctx context.Context, serverID int, siteID int, workerID int) error {
	path := fmt.Sprintf("/orgs/%s/servers/%d/sites/%d/workers/%d/actions", c.OrgSlug, serverID, siteID, workerID)
	req := map[string]string{"action": "restart"}
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}

type workerOutputResponse struct {
	Data struct {
		Output string `json:"output"`
	} `json:"data"`
}

func (c *Client) GetWorkerOutput(ctx context.Context, serverID int, siteID int, workerID int) (string, error) {
	path := fmt.Sprintf("/orgs/%s/servers/%d/sites/%d/workers/%d/log", c.OrgSlug, serverID, siteID, workerID)
	var res workerOutputResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return "", err
	}
	return res.Data.Output, nil
}
