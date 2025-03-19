package forge_client

import (
	"context"
	"fmt"
	"net/http"
)

type IntegrationDaemon struct {
	ID           int    `json:"id"`
	Command      string `json:"command"`
	User         string `json:"user"`
	Directory    string `json:"directory"`
	Processes    int    `json:"processes"`
	StartSecs    int    `json:"startsecs"`
	StopWaitSecs int    `json:"stopwaitsecs"`
	StopSignal   string `json:"stopsignal"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
}

// Horizon.
type HorizonStatus struct {
	Enabled          bool               `json:"enabled"`
	Daemon           *IntegrationDaemon `json:"daemon"`
	HorizonInstalled bool               `json:"horizon_installed"`
}

func (c *Client) CheckHorizonStatus(ctx context.Context, serverID, siteID int) (*HorizonStatus, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/integrations/horizon", serverID, siteID)
	var res HorizonStatus
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) EnableHorizon(ctx context.Context, serverID, siteID int) (*IntegrationDaemon, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/integrations/horizon", serverID, siteID)
	var res struct {
		Daemon IntegrationDaemon `json:"daemon"`
	}
	if err := c.doRequest(ctx, http.MethodPost, path, nil, &res); err != nil {
		return nil, err
	}
	return &res.Daemon, nil
}

func (c *Client) DisableHorizon(ctx context.Context, serverID, siteID int) error {
	path := fmt.Sprintf("/servers/%d/sites/%d/integrations/horizon", serverID, siteID)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

// Octane.
type OctaneStatus struct {
	Enabled         bool               `json:"enabled"`
	Daemon          *IntegrationDaemon `json:"daemon"`
	OctaneInstalled bool               `json:"octane_installed"`
}

type EnableOctaneRequest struct {
	Port   int    `json:"port,omitempty"`
	Server string `json:"server,omitempty"`
}

func (c *Client) CheckOctaneStatus(ctx context.Context, serverID, siteID int) (*OctaneStatus, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/integrations/octane", serverID, siteID)
	var res OctaneStatus
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) EnableOctane(ctx context.Context, serverID, siteID int, port int, serverType string) (*IntegrationDaemon, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/integrations/octane", serverID, siteID)
	req := EnableOctaneRequest{Port: port, Server: serverType}
	var res struct {
		Daemon IntegrationDaemon `json:"daemon"`
	}
	if err := c.doRequest(ctx, http.MethodPost, path, req, &res); err != nil {
		return nil, err
	}
	return &res.Daemon, nil
}

func (c *Client) DisableOctane(ctx context.Context, serverID, siteID int) error {
	path := fmt.Sprintf("/servers/%d/sites/%d/integrations/octane", serverID, siteID)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

// Reverb.
type ReverbStatus struct {
	Enabled           bool               `json:"enabled"`
	Daemon            *IntegrationDaemon `json:"daemon"`
	ReverbHost        *string            `json:"reverb_host"`
	ReverbPort        *int               `json:"reverb_port"`
	ReverbConnections *int               `json:"reverb_connections"`
	ReverbInstalled   bool               `json:"reverb_installed"`
}

type EnableReverbRequest struct {
	Port        int    `json:"port,omitempty"`
	Host        string `json:"host,omitempty"`
	Connections int    `json:"connections,omitempty"`
}

func (c *Client) CheckReverbStatus(ctx context.Context, serverID, siteID int) (*ReverbStatus, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/integrations/reverb", serverID, siteID)
	var res ReverbStatus
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) EnableReverb(ctx context.Context, serverID, siteID int, port int, host string, connections int) (*IntegrationDaemon, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/integrations/reverb", serverID, siteID)
	req := EnableReverbRequest{Port: port, Host: host, Connections: connections}
	var res struct {
		Daemon            IntegrationDaemon `json:"daemon"`
		ReverbHost        string            `json:"reverb_host"`
		ReverbPort        int               `json:"reverb_port"`
		ReverbConnections int               `json:"reverb_connections"`
	}
	if err := c.doRequest(ctx, http.MethodPost, path, req, &res); err != nil {
		return nil, err
	}
	return &res.Daemon, nil
}

func (c *Client) DisableReverb(ctx context.Context, serverID, siteID int) error {
	path := fmt.Sprintf("/servers/%d/sites/%d/integrations/reverb", serverID, siteID)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

// Pulse.
type PulseStatus struct {
	Enabled        bool               `json:"enabled"`
	Daemon         *IntegrationDaemon `json:"daemon"`
	PulseInstalled bool               `json:"pulse_installed"`
}

func (c *Client) CheckPulseStatus(ctx context.Context, serverID, siteID int) (*PulseStatus, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/integrations/pulse", serverID, siteID)
	var res PulseStatus
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) EnablePulse(ctx context.Context, serverID, siteID int) (*IntegrationDaemon, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/integrations/pulse", serverID, siteID)
	var res struct {
		Daemon IntegrationDaemon `json:"daemon"`
	}
	if err := c.doRequest(ctx, http.MethodPost, path, nil, &res); err != nil {
		return nil, err
	}
	return &res.Daemon, nil
}

func (c *Client) DisablePulse(ctx context.Context, serverID, siteID int) error {
	path := fmt.Sprintf("/servers/%d/sites/%d/integrations/pulse", serverID, siteID)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

// Inertia.
type InertiaStatus struct {
	Enabled          bool               `json:"enabled"`
	Daemon           *IntegrationDaemon `json:"daemon"`
	InertiaInstalled bool               `json:"inertia_installed"`
}

type EnableInertiaRequest struct {
	DeploysRestartInertiaDaemon bool `json:"deploys_restart_inertia_daemon,omitempty"`
}

func (c *Client) CheckInertiaStatus(ctx context.Context, serverID, siteID int) (*InertiaStatus, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/integrations/inertia", serverID, siteID)
	var res InertiaStatus
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) EnableInertia(ctx context.Context, serverID, siteID int, deploysRestart bool) (*IntegrationDaemon, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/integrations/inertia", serverID, siteID)
	req := EnableInertiaRequest{DeploysRestartInertiaDaemon: deploysRestart}
	var res struct {
		Daemon IntegrationDaemon `json:"daemon"`
	}
	if err := c.doRequest(ctx, http.MethodPost, path, req, &res); err != nil {
		return nil, err
	}
	return &res.Daemon, nil
}

func (c *Client) DisableInertia(ctx context.Context, serverID, siteID int) error {
	path := fmt.Sprintf("/servers/%d/sites/%d/integrations/inertia", serverID, siteID)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

// Laravel Maintenance.
type LaravelMaintenanceStatus struct {
	Enabled          bool    `json:"enabled"`
	Status           *string `json:"status"`
	LaravelInstalled bool    `json:"laravel_installed"`
}

type EnableLaravelMaintenanceRequest struct {
	Secret string `json:"secret,omitempty"`
	Status int    `json:"status,omitempty"`
}

func (c *Client) CheckLaravelMaintenance(ctx context.Context, serverID, siteID int) (*LaravelMaintenanceStatus, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/integrations/laravel-maintenance", serverID, siteID)
	var res LaravelMaintenanceStatus
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) EnableLaravelMaintenance(ctx context.Context, serverID, siteID int, secret string, status int) error {
	path := fmt.Sprintf("/servers/%d/sites/%d/integrations/laravel-maintenance", serverID, siteID)
	req := EnableLaravelMaintenanceRequest{Secret: secret, Status: status}
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}

func (c *Client) DisableLaravelMaintenance(ctx context.Context, serverID, siteID int) error {
	path := fmt.Sprintf("/servers/%d/sites/%d/integrations/laravel-maintenance", serverID, siteID)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

// Laravel Scheduler.
type LaravelSchedulerStatus struct {
	Enabled bool `json:"enabled"`
	Job     *struct {
		ID          int    `json:"id"`
		Command     string `json:"command"`
		User        string `json:"user"`
		Frequency   string `json:"frequency"`
		Cron        string `json:"cron"`
		Status      string `json:"status"`
		CreatedAt   string `json:"created_at"`
		NextRunTime string `json:"next_run_time"`
	} `json:"job"`
	LaravelInstalled bool `json:"laravel_installed"`
}

func (c *Client) CheckLaravelScheduler(ctx context.Context, serverID, siteID int) (*LaravelSchedulerStatus, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/integrations/laravel-scheduler", serverID, siteID)
	var res LaravelSchedulerStatus
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

type LaravelSchedulerJob struct {
	Command     string `json:"command"`
	User        string `json:"user"`
	Frequency   string `json:"frequency"`
	Cron        string `json:"cron"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	ID          int    `json:"id"`
	NextRunTime string `json:"next_run_time"`
}

func (c *Client) EnableLaravelScheduler(ctx context.Context, serverID, siteID int) (*LaravelSchedulerJob, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/integrations/laravel-scheduler", serverID, siteID)
	var res struct {
		Job LaravelSchedulerJob `json:"job"`
	}
	if err := c.doRequest(ctx, http.MethodPost, path, nil, &res); err != nil {
		return nil, err
	}
	return &res.Job, nil
}

func (c *Client) DisableLaravelScheduler(ctx context.Context, serverID, siteID int) error {
	path := fmt.Sprintf("/servers/%d/sites/%d/integrations/laravel-scheduler", serverID, siteID)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

// Composer Packages Authentication.
type ComposerCredentialsResponse struct {
	Credentials map[string]map[string]string `json:"credentials"`
}

func (c *Client) GetComposerPackagesAuth(ctx context.Context, serverID, siteID int) (*ComposerCredentialsResponse, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/packages", serverID, siteID)
	var res ComposerCredentialsResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

type UpdateComposerPackagesAuthRequest struct {
	Credentials []struct {
		RepositoryURL string `json:"repository_url"`
		Username      string `json:"username"`
		Password      string `json:"password"`
	} `json:"credentials"`
}

func (c *Client) UpdateComposerPackagesAuth(ctx context.Context, serverID, siteID int, req UpdateComposerPackagesAuthRequest) error {
	path := fmt.Sprintf("/servers/%d/sites/%d/packages", serverID, siteID)
	return c.doRequest(ctx, http.MethodPut, path, req, nil)
}
