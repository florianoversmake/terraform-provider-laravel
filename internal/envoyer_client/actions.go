package envoyer_client

import (
	"context"
	"fmt"
)

// -------------------------------
// New Types for Hooks, Actions, and Environment
// -------------------------------

// Action represents an action in Envoyer.
type Action struct {
	ID        int64  `json:"id"`
	Version   int64  `json:"version"`
	Name      string `json:"name"`
	View      string `json:"view"`
	Sequence  int64  `json:"sequence"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// Hook represents a deployment hook in Envoyer.
type Hook struct {
	ID        int64  `json:"id"`
	ProjectID int64  `json:"project_id"`
	ActionID  int64  `json:"action_id"`
	Timing    string `json:"timing"`
	Name      string `json:"name"`
	RunAs     string `json:"run_as"`
	Script    string `json:"script"`
	Sequence  int64  `json:"sequence"`
	// Servers is a slice of server IDs (optional)
	Servers   []int64 `json:"servers,omitempty"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

// CreateHookRequest is the payload used when creating a new hook.
type CreateHookRequest struct {
	// ProjectID is not sent as part of the JSON payload;
	// it is used to build the endpoint URL.
	ProjectID int64   `json:"-"`
	ActionID  int64   `json:"actionId"`
	Timing    string  `json:"timing"`
	Name      string  `json:"name"`
	RunAs     string  `json:"runAs"`
	Script    string  `json:"script"`
	Servers   []int64 `json:"servers,omitempty"`
}

// UpdateHookRequest is the payload for updating an existing hook.
type UpdateHookRequest struct {
	Servers []int64 `json:"servers"`
}

// -------------------------------
// New Client Methods
// -------------------------------

// ListActions retrieves the list of available actions from Envoyer.
func (c *Client) ListActions(ctx context.Context) ([]Action, error) {
	var result struct {
		Actions []Action `json:"actions"`
	}
	err := c.doRequest(ctx, "GET", "/actions", nil, &result)
	if err != nil {
		return nil, err
	}
	return result.Actions, nil
}

// CreateHook creates a new hook on the given project.
func (c *Client) CreateHook(ctx context.Context, req CreateHookRequest) (Hook, error) {
	var result struct {
		Hook Hook `json:"hook"`
	}
	endpoint := fmt.Sprintf("/projects/%d/hooks", req.ProjectID)
	err := c.doRequest(ctx, "POST", endpoint, req, &result)
	if err != nil {
		return Hook{}, err
	}
	return result.Hook, nil
}

// GetHook retrieves the details of a hook by project ID and hook ID.
func (c *Client) GetHook(ctx context.Context, projectID, hookID int) (Hook, error) {
	var result struct {
		Hook Hook `json:"hook"`
	}
	endpoint := fmt.Sprintf("/projects/%d/hooks/%d", projectID, hookID)
	err := c.doRequest(ctx, "GET", endpoint, nil, &result)
	if err != nil {
		return Hook{}, err
	}
	return result.Hook, nil
}

// UpdateHook updates an existing hook (e.g. to change the server list).
func (c *Client) UpdateHook(ctx context.Context, projectID, hookID int, req UpdateHookRequest) error {
	endpoint := fmt.Sprintf("/projects/%d/hooks/%d", projectID, hookID)
	return c.doRequest(ctx, "PUT", endpoint, req, nil)
}

// DeleteHook deletes a hook from the project.
func (c *Client) DeleteHook(ctx context.Context, projectID, hookID int) error {
	endpoint := fmt.Sprintf("/projects/%d/hooks/%d", projectID, hookID)
	return c.doRequest(ctx, "DELETE", endpoint, nil, nil)
}
