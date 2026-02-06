package forge_client

import (
	"context"
	"fmt"
	"net/http"
)

type Job struct {
	ID          int64   `json:"id"`
	Name        *string `json:"name"`
	Command     string  `json:"command"`
	User        string  `json:"user"`
	Frequency   string  `json:"frequency"`
	Cron        string  `json:"cron"`
	NextRunTime string  `json:"next_run_time"`
	Status      string  `json:"status"`
	CreatedAt   *string `json:"created_at"`
	UpdatedAt   *string `json:"updated_at"`
}

type CreateJobRequest struct {
	Name      *string `json:"name,omitempty"`
	Command   string  `json:"command"`
	User      string  `json:"user"`
	Frequency string  `json:"frequency"`
	Cron      *string `json:"cron,omitempty"`
}

func (c *Client) CreateJob(ctx context.Context, serverID int, req CreateJobRequest) (*Job, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/scheduled-jobs", serverID))
	var job Job
	id, err := c.PostJsonApi(ctx, path, req, &job)
	if err != nil {
		return nil, err
	}
	job.ID = int64(id)
	return &job, nil
}

func (c *Client) ListJobs(ctx context.Context, serverID int) ([]Job, error) {
	items, err := c.GetJsonApiListAll(ctx, c.orgPath(fmt.Sprintf("/servers/%d/scheduled-jobs", serverID)))
	if err != nil {
		return nil, err
	}
	return unmarshalList(items, func(j *Job, id int64) { j.ID = id })
}

func (c *Client) GetJob(ctx context.Context, serverID, jobID int) (*Job, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/scheduled-jobs/%d", serverID, jobID))
	var job Job
	id, err := c.GetJsonApi(ctx, path, &job)
	if err != nil {
		return nil, err
	}
	job.ID = int64(id)
	return &job, nil
}

func (c *Client) DeleteJob(ctx context.Context, serverID, jobID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/scheduled-jobs/%d", serverID, jobID))
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

func (c *Client) GetJobOutput(ctx context.Context, serverID, jobID int) (string, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/scheduled-jobs/%d/output", serverID, jobID))
	var logRes struct {
		Content string `json:"content"`
	}
	_, err := c.GetJsonApi(ctx, path, &logRes)
	if err != nil {
		return "", err
	}
	return logRes.Content, nil
}
