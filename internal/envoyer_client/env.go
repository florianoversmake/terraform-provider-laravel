// Copyright (c) HashiCorp, Inc.

package envoyer_client

import (
	"context"
	"fmt"
)

// UpdateEnvironmentRequest is the payload used for creating/updating an environment.
type UpdateEnvironmentRequest struct {
	Contents string  `json:"contents"`
	Servers  []int64 `json:"servers,omitempty"`
}

type UpdateEnvironmentRequestInternal struct {
	Contents string  `json:"contents"`
	Servers  []int64 `json:"servers,omitempty"`
	Key      string  `json:"key"`
}

// UpdateEnvironment creates or updates the environment file for a project.
// It returns the updated environment contents.
func (c *Client) UpdateEnvironment(ctx context.Context, projectID int, req UpdateEnvironmentRequest) (string, error) {
	var result struct {
		Environment string `json:"environment"`
	}
	endpoint := fmt.Sprintf("/projects/%d/environment", projectID)

	// Convert the request to the internal format.
	internalReq := UpdateEnvironmentRequestInternal{
		Contents: req.Contents,
		Servers:  req.Servers,
		Key:      c.envKey,
	}

	err := c.doRequest(ctx, "PUT", endpoint, internalReq, &result)
	if err != nil {
		return "", err
	}
	return result.Environment, nil
}

func (c *Client) GetEnvironment(ctx context.Context, projectID int) (string, error) {
	var result struct {
		Environment string `json:"environment"`
	}
	endpoint := fmt.Sprintf("/projects/%d/environment", projectID)
	payload := map[string]string{"key": c.envKey}
	err := c.doRequest(ctx, "GET", endpoint, payload, &result)
	if err != nil {
		return "", err
	}
	return result.Environment, nil
}

func (c *Client) GetEnvironmentServers(ctx context.Context, projectID int) ([]int64, error) {
	var result struct {
		Servers []Server `json:"servers"`
	}

	endpoint := fmt.Sprintf("/projects/%d/environment/servers", projectID)
	err := c.doRequest(ctx, "GET", endpoint, nil, &result)
	if err != nil {
		return nil, err
	}

	var servers []int64
	for _, server := range result.Servers {
		servers = append(servers, server.ID)
	}

	return servers, nil
}
