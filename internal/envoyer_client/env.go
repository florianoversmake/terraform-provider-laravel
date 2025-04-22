package envoyer_client

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Environment represents a parsed environment.
type Environment struct {
	Contents  string
	Servers   []int64
	Variables map[string]string
	Comments  []string
	ProjectID int
}

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

	// debug log the request
	tflog.Debug(ctx, "envoyer: UpdateEnvironment", map[string]any{
		"endpoint": endpoint,
		"request":  internalReq,
	})

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

// ParseEnvironment parses raw environment contents into structured format.
func ParseEnvironment(contents string) map[string]string {
	vars := make(map[string]string)

	if contents == "" {
		return vars
	}

	lines := strings.Split(contents, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Handle quoted values
		if len(value) >= 2 {
			if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
				(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
				value = value[1 : len(value)-1]
			}
		}

		vars[key] = value
	}

	return vars
}

// GetEnvironmentVariables retrieves and parses environment variables.
func (c *Client) GetEnvironmentVariables(ctx context.Context, projectID int) (map[string]string, error) {
	contents, err := c.GetEnvironment(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return ParseEnvironment(contents), nil
}

// GetEnvironmentVariable gets a specific environment variable.
func (c *Client) GetEnvironmentVariable(ctx context.Context, projectID int, key string) (string, bool, error) {
	vars, err := c.GetEnvironmentVariables(ctx, projectID)
	if err != nil {
		return "", false, err
	}

	value, exists := vars[key]
	return value, exists, nil
}

// SetEnvironmentVariable sets a specific environment variable, preserving the rest.
func (c *Client) SetEnvironmentVariable(ctx context.Context, projectID int, key, value string, servers []int64) error {
	// Get current environment
	contents, err := c.GetEnvironment(ctx, projectID)
	if err != nil {
		// If the environment doesn't exist yet, start with an empty one
		if strings.Contains(err.Error(), "404") {
			contents = ""
		} else {
			return err
		}
	}

	// Parse into lines for manual editing to preserve formatting
	vars := make(map[string]string)
	var nonVarLines []string
	var updatedLines []string

	// First pass: identify existing vars and collect non-var lines
	if contents != "" {
		lines := strings.Split(contents, "\n")
		for _, line := range lines {
			trimmedLine := strings.TrimSpace(line)
			if trimmedLine == "" || strings.HasPrefix(trimmedLine, "#") {
				nonVarLines = append(nonVarLines, line)
				continue
			}

			parts := strings.SplitN(trimmedLine, "=", 2)
			if len(parts) != 2 {
				nonVarLines = append(nonVarLines, line)
				continue
			}

			currKey := strings.TrimSpace(parts[0])
			currValue := strings.TrimSpace(parts[1])
			vars[currKey] = currValue
		}
	}

	// Second pass: reconstruct with the updated variable
	// Add non-variable lines first
	updatedLines = append(updatedLines, nonVarLines...)

	// Add or update our target variable
	vars[key] = value

	// Add all variables in sorted order for consistency
	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}

	// Sort keys alphabetically
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[i] > keys[j] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}

	// Add variables in sorted order
	for _, k := range keys {
		v := vars[k]

		// Determine if we need quotes
		needsQuotes := strings.ContainsAny(v, " \t\r\n'\"=")
		if needsQuotes {
			// Prefer double quotes unless value contains double quotes
			if strings.Contains(v, "\"") {
				v = "'" + v + "'"
			} else {
				v = "\"" + v + "\""
			}
		}

		updatedLines = append(updatedLines, fmt.Sprintf("%s=%s", k, v))
	}

	// Create the updated content
	updatedContent := strings.Join(updatedLines, "\n")

	// Update via the API
	_, err = c.UpdateEnvironment(ctx, projectID, UpdateEnvironmentRequest{
		Contents: updatedContent,
		Servers:  servers,
	})

	return err
}

// DeleteEnvironmentVariable removes a specific environment variable, preserving the rest.
func (c *Client) DeleteEnvironmentVariable(ctx context.Context, projectID int, key string, servers []int64) error {
	// Get current environment
	contents, err := c.GetEnvironment(ctx, projectID)
	if err != nil {
		return err
	}

	// Parse into lines for manual editing to preserve formatting
	vars := make(map[string]string)
	var nonVarLines []string
	var updatedLines []string
	keyFound := false

	// First pass: identify existing vars and collect non-var lines
	if contents != "" {
		lines := strings.Split(contents, "\n")
		for _, line := range lines {
			trimmedLine := strings.TrimSpace(line)
			if trimmedLine == "" || strings.HasPrefix(trimmedLine, "#") {
				nonVarLines = append(nonVarLines, line)
				continue
			}

			parts := strings.SplitN(trimmedLine, "=", 2)
			if len(parts) != 2 {
				nonVarLines = append(nonVarLines, line)
				continue
			}

			currKey := strings.TrimSpace(parts[0])
			if currKey == key {
				keyFound = true
				continue // Skip this key
			}

			currValue := strings.TrimSpace(parts[1])
			vars[currKey] = currValue
		}
	}

	// If key not found, nothing to do
	if !keyFound {
		return nil
	}

	// Second pass: reconstruct without the deleted variable
	// Add non-variable lines first
	updatedLines = append(updatedLines, nonVarLines...)

	// Add all variables in sorted order for consistency
	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}

	// Sort keys alphabetically
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[i] > keys[j] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}

	// Add variables in sorted order
	for _, k := range keys {
		v := vars[k]

		// Determine if we need quotes
		needsQuotes := strings.ContainsAny(v, " \t\r\n'\"=")
		if needsQuotes {
			// Prefer double quotes unless value contains double quotes
			if strings.Contains(v, "\"") {
				v = "'" + v + "'"
			} else {
				v = "\"" + v + "\""
			}
		}

		updatedLines = append(updatedLines, fmt.Sprintf("%s=%s", k, v))
	}

	// Create the updated content
	updatedContent := strings.Join(updatedLines, "\n")

	// Update via the API
	_, err = c.UpdateEnvironment(ctx, projectID, UpdateEnvironmentRequest{
		Contents: updatedContent,
		Servers:  servers,
	})

	return err
}

// BulkUpdateEnvironmentVariables updates multiple environment variables at once.
func (c *Client) BulkUpdateEnvironmentVariables(ctx context.Context, projectID int, vars map[string]string, servers []int64) error {
	// Get current environment
	contents, err := c.GetEnvironment(ctx, projectID)
	if err != nil {
		// If the environment doesn't exist yet, start with an empty one
		if strings.Contains(err.Error(), "404") {
			contents = ""
		} else {
			return err
		}
	}

	// Parse current environment
	currentVars := ParseEnvironment(contents)

	// Update with new values
	for k, v := range vars {
		currentVars[k] = v
	}

	// Format all variables
	var lines []string

	// Add all variables in sorted order for consistency
	keys := make([]string, 0, len(currentVars))
	for k := range currentVars {
		keys = append(keys, k)
	}

	// Sort keys alphabetically
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[i] > keys[j] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}

	// Add variables in sorted order
	for _, k := range keys {
		v := currentVars[k]

		// Determine if we need quotes
		needsQuotes := strings.ContainsAny(v, " \t\r\n'\"=")
		if needsQuotes {
			// Prefer double quotes unless value contains double quotes
			if strings.Contains(v, "\"") {
				v = "'" + v + "'"
			} else {
				v = "\"" + v + "\""
			}
		}

		lines = append(lines, fmt.Sprintf("%s=%s", k, v))
	}

	// Create the updated content
	updatedContent := strings.Join(lines, "\n")

	// Update via the API
	_, err = c.UpdateEnvironment(ctx, projectID, UpdateEnvironmentRequest{
		Contents: updatedContent,
		Servers:  servers,
	})

	return err
}
