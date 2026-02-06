package forge_client

import (
	"context"
	"fmt"
	"net/http"
)

type Worker struct {
	ID        int64   `json:"id"`
	Command   string  `json:"command"`
	User      string  `json:"user"`
	Directory *string `json:"directory"`
	Processes int     `json:"processes"`
	Status    string  `json:"status"`
	CreatedAt string  `json:"created_at"`
}

type CreateWorkerRequest struct {
	Name         string  `json:"name"`
	Command      string  `json:"command"`
	User         string  `json:"user"`
	Directory    *string `json:"directory,omitempty"`
	Processes    int     `json:"processes"`
	StartSecs    *int    `json:"startsecs,omitempty"`
	StopWaitSecs *int    `json:"stopwaitsecs,omitempty"`
	StopSignal   *string `json:"stopsignal,omitempty"`
}

func (c *Client) ListWorkers(ctx context.Context, serverID int, siteID int) ([]Worker, error) {
	items, err := c.GetJsonApiListAll(ctx, c.orgPath(fmt.Sprintf("/servers/%d/background-processes", serverID)))
	if err != nil {
		return nil, err
	}
	return unmarshalList(items, func(w *Worker, id int64) { w.ID = id })
}

func (c *Client) CreateWorker(ctx context.Context, serverID int, siteID int, req CreateWorkerRequest) (*Worker, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/background-processes", serverID))
	var worker Worker
	id, err := c.PostJsonApi(ctx, path, req, &worker)
	if err != nil {
		return nil, err
	}
	worker.ID = int64(id)
	return &worker, nil
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
	path := c.orgPath(fmt.Sprintf("/servers/%d/background-processes/%d", serverID, workerID))
	var worker Worker
	id, err := c.GetJsonApi(ctx, path, &worker)
	if err != nil {
		if _, ok := err.(*ClientErrorResourceNotFound); ok {
			return nil, &ErrorWorkerNotFound{ServerID: serverID, SiteID: siteID, WorkerID: workerID}
		}
		return nil, err
	}
	worker.ID = int64(id)
	return &worker, nil
}

func (c *Client) DeleteWorker(ctx context.Context, serverID int, siteID int, workerID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/background-processes/%d", serverID, workerID))
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

func (c *Client) RestartWorker(ctx context.Context, serverID int, siteID int, workerID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/background-processes/%d/actions", serverID, workerID))
	req := struct {
		Action string `json:"action"`
	}{Action: "restart"}
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}

func (c *Client) GetWorkerOutput(ctx context.Context, serverID int, siteID int, workerID int) (string, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/background-processes/%d/log", serverID, workerID))
	var logRes struct {
		Content string `json:"content"`
	}
	_, err := c.GetJsonApi(ctx, path, &logRes)
	if err != nil {
		return "", err
	}
	return logRes.Content, nil
}
