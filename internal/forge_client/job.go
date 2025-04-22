package forge_client

import (
	"context"
	"fmt"
	"net/http"
)

type Job struct {
	ID        int64  `json:"id"`
	Command   string `json:"command"`
	User      string `json:"user"`
	Frequency string `json:"frequency"`
	Cron      string `json:"cron"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

type CreateJobRequest struct {
	Command   string `json:"command"`
	Frequency string `json:"frequency"`
	User      string `json:"user"`
	Minute    string `json:"minute,omitempty"`
	Hour      string `json:"hour,omitempty"`
	Day       string `json:"day,omitempty"`
	Month     string `json:"month,omitempty"`
	Weekday   string `json:"weekday,omitempty"`
}

type jobResponse struct {
	Job Job `json:"job"`
}

type jobsResponse struct {
	Jobs []Job `json:"jobs"`
}

func (c *Client) CreateJob(ctx context.Context, serverID int, req CreateJobRequest) (*Job, error) {
	path := fmt.Sprintf("/servers/%d/jobs", serverID)
	var res jobResponse
	if err := c.doRequest(ctx, http.MethodPost, path, req, &res); err != nil {
		return nil, err
	}
	return &res.Job, nil
}

func (c *Client) ListJobs(ctx context.Context, serverID int) ([]Job, error) {
	path := fmt.Sprintf("/servers/%d/jobs", serverID)
	var res jobsResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return res.Jobs, nil
}

func (c *Client) GetJob(ctx context.Context, serverID, jobID int) (*Job, error) {
	path := fmt.Sprintf("/servers/%d/jobs/%d", serverID, jobID)
	var res jobResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return &res.Job, nil
}

func (c *Client) DeleteJob(ctx context.Context, serverID, jobID int) error {
	path := fmt.Sprintf("/servers/%d/jobs/%d", serverID, jobID)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

type jobOutputResponse struct {
	Output string `json:"output"`
}

func (c *Client) GetJobOutput(ctx context.Context, serverID, jobID int) (string, error) {
	path := fmt.Sprintf("/servers/%d/jobs/%d/output", serverID, jobID)
	var res jobOutputResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return "", err
	}
	return res.Output, nil
}
