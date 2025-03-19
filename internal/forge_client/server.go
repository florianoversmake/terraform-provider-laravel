// Copyright (c) HashiCorp, Inc.

package forge_client

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	ID               int64   `json:"id"`
	CredentialID     int64   `json:"credential_id"`
	Name             string  `json:"name"`
	Type             string  `json:"type"`
	Provider         string  `json:"provider"`   //undocumented
	Identifier       string  `json:"identifier"` //undocumented
	Size             string  `json:"size"`
	Region           string  `json:"region"`
	UbuntuVersion    string  `json:"ubuntu_version"` //undocumented
	DBStatus         *string `json:"db_status"`      //undocumented
	RedisStatus      *string `json:"redis_status"`   //undocumented
	PHPVersion       string  `json:"php_version"`
	PHPCLIVersion    string  `json:"php_cli_version"`
	OpcacheStatus    *string `json:"opcache_status"`
	DatabaseType     string  `json:"database_type"`
	IPAddress        *string `json:"ip_address"`
	SSHPort          int     `json:"ssh_port"`
	PrivateIPAddress *string `json:"private_ip_address"`
	LocalPublicKey   string  `json:"local_public_key"` //undocumented
	BlackfireStatus  *string `json:"blackfire_status"`
	PapertrailStatus *string `json:"papertrail_status"`
	Revoked          bool    `json:"revoked"`
	CreatedAt        string  `json:"created_at"`
	IsReady          bool    `json:"is_ready"`
	Tags             []Tag   `json:"tags"` //undocumented
	Network          []int64 `json:"network"`
}

type Tag struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

type serversResponse struct {
	Servers []Server `json:"servers"`
}

type serverResponse struct {
	Server Server `json:"server"`
}

func (c *Client) ListServers(ctx context.Context) ([]Server, error) {
	var resp serversResponse
	if err := c.doRequest(ctx, http.MethodGet, "/servers", nil, &resp); err != nil {
		return nil, err
	}
	return resp.Servers, nil
}

func (c *Client) GetServer(ctx context.Context, serverID int) (*Server, error) {
	path := fmt.Sprintf("/servers/%d", serverID)
	var resp serverResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Server, nil
}

// Parameters
// Key	Description
// ubuntu_version	The version of Ubuntu to create the server with. Valid values are "20.04", "22.04", and "24.04". "24.04" is used by default if no value is defined. It is recommended to always specify a version as the default may change at any time.
// type	The type of server to create. Valid values are app, web, loadbalancer, cache, database, worker, meilisearch. app is used by default if no value is defined.
// provider	The server provider. Valid values are ocean2 for Digital Ocean, akamai (Linode), vultr2, aws, hetzner and custom.
// disk_size	The size of the disk in GB. Valid when the provider is aws. Minimum of 8GB. Example: 20.
// circle	The ID of a circle to create the server within.
// credential_id	This is only required when the provider is not custom.
// region	The name of the region where the server will be created. This value is not required you are building a Custom VPS server. Valid region identifiers.
// ip_address	The IP Address of the server. Only required when the provider is custom.
// private_ip_address	The Private IP Address of the server. Only required when the provider is custom.
// php_version	Valid values are php84, php83, php82, php81, php80, php74, php73,php72,php82, php70, and php56.
// database	The name of the database Forge should create when building the server. If omitted, forge will be used.
// database_type	Valid values are mysql8, mariadb106, mariadb1011, mariadb114, postgres, postgres13, postgres14, postgres15, postgres16 or postgres17.
// network	An array of server IDs that the server should be able to connect to.
// recipe_id	An optional ID of a recipe to run after provisioning.
// aws_vpc_id	ID of the existing VPC
// aws_subnet_id	ID of the existing subnet
// aws_vpc_name	When creating a new one
// hetzner_network_id	ID of the existing VPC
// ocean2_vpc_uuid	UUID of the existing VPC
// ocean2_vpc_name	When creating a new one
// vultr2_network_id	ID of the existing private network
// vultr2_network_name	When creating a new one

// CreateServerRequest is the payload to create a server.
type CreateServerRequest struct {
	UbuntuVersion     string  `json:"ubuntu_version"`
	Type              string  `json:"type"`
	Name              string  `json:"name"`
	Provider          string  `json:"provider"`
	Size              *string `json:"size,omitempty"`
	DiskSize          *int32  `json:"disk_size,omitempty"`
	Circle            *int64  `json:"circle,omitempty"`
	CredentialID      *int64  `json:"credential_id,omitempty"`
	Region            *string `json:"region"`
	IPAddress         *string `json:"ip_address,omitempty"`
	PrivateIPAddress  *string `json:"private_ip_address,omitempty"`
	SSHPort           *int32  `json:"ssh_port,omitempty"`
	PHPVersion        string  `json:"php_version"`
	Database          *string `json:"database,omitempty"`
	DatabaseType      *string `json:"database_type,omitempty"`
	Network           []int64 `json:"network"`
	RecipeID          *int64  `json:"recipe_id,omitempty"`
	AWSVPCID          *string `json:"aws_vpc_id,omitempty"`
	AWSSubnetID       *string `json:"aws_subnet_id,omitempty"`
	AWSVPCName        *string `json:"aws_vpc_name,omitempty"`
	HetznerNetworkID  *string `json:"hetzner_network_id,omitempty"`
	Ocean2VPCUUID     *string `json:"ocean2_vpc_uuid,omitempty"`
	Ocean2VPCName     *string `json:"ocean2_vpc_name,omitempty"`
	Vultr2NetworkID   *string `json:"vultr2_network_id,omitempty"`
	Vultr2NetworkName *string `json:"vultr2_network_name,omitempty"`
}

type CreateServerResponse struct {
	Server              Server  `json:"server"`
	SudoPassword        string  `json:"sudo_password"`
	DatabasePassword    *string `json:"database_password"`
	MeilisearchPassword *string `json:"meilisearch_password"`
	ProvisionCommand    *string `json:"provision_command"`
}

func (c *Client) CreateServer(ctx context.Context, req CreateServerRequest) (*CreateServerResponse, error) {
	var resp CreateServerResponse
	if err := c.doRequest(ctx, http.MethodPost, "/servers", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Update payload
// {
//     "name": "renamed-server",
//     "ip_address": "192.241.143.108",
//     "private_ip_address": "10.136.8.40",
//     "max_upload_size": 123,
//     "max_execution_time": 30,
//     "network": [
//         2,
//         3
//     ],
//     "timezone": "Europe/London",
//     "tags": [
//         "london-server"
//     ]
// }

type UpdateServerRequest struct {
	Name             string   `json:"name"`
	IPAddress        string   `json:"ip_address"`
	PrivateIPAddress string   `json:"private_ip_address"`
	MaxUploadSize    int      `json:"max_upload_size"`
	MaxExecutionTime int      `json:"max_execution_time"`
	Network          []int    `json:"network"`
	Timezone         string   `json:"timezone"`
	Tags             []string `json:"tags"`
}

func (c *Client) UpdateServer(ctx context.Context, serverID int, req UpdateServerRequest) (*Server, error) {
	path := fmt.Sprintf("/servers/%d", serverID)
	var resp serverResponse
	if err := c.doRequest(ctx, http.MethodPut, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp.Server, nil
}

func (c *Client) DeleteServer(ctx context.Context, serverID int) error {
	path := fmt.Sprintf("/servers/%d", serverID)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

func (c *Client) RebootServer(ctx context.Context, serverID int) error {
	path := fmt.Sprintf("/servers/%d/reboot", serverID)
	return c.doRequest(ctx, http.MethodPost, path, nil, nil)
}

func (c *Client) RevokeServer(ctx context.Context, serverID int) error {
	path := fmt.Sprintf("/servers/%d/revoke", serverID)
	return c.doRequest(ctx, http.MethodPost, path, nil, nil)
}

type ReconnectServerResponse struct {
	PublicKey string `json:"public_key"`
}

func (c *Client) ReconnectServer(ctx context.Context, serverID int) (*ReconnectServerResponse, error) {
	path := fmt.Sprintf("/servers/%d/reconnect", serverID)
	var resp ReconnectServerResponse
	if err := c.doRequest(ctx, http.MethodPost, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) ReactivateServer(ctx context.Context, serverID int) error {
	path := fmt.Sprintf("/servers/%d/reactivate", serverID)
	return c.doRequest(ctx, http.MethodPost, path, nil, nil)
}

type startServiceRequest struct {
	Service string `json:"service"`
}

func (c *Client) StartService(ctx context.Context, serverID int, service string) error {
	path := fmt.Sprintf("/servers/%d/services/start", serverID)
	req := startServiceRequest{Service: service}
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}

func (c *Client) StopService(ctx context.Context, serverID int, service string) error {
	path := fmt.Sprintf("/servers/%d/services/stop", serverID)
	req := startServiceRequest{Service: service}
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}

func (c *Client) RestartService(ctx context.Context, serverID int, service string) error {
	path := fmt.Sprintf("/servers/%d/services/restart", serverID)
	req := startServiceRequest{Service: service}
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}

func (c *Client) RebootMySQL(ctx context.Context, serverID int) error {
	path := fmt.Sprintf("/servers/%d/mysql/reboot", serverID)
	return c.doRequest(ctx, http.MethodPost, path, nil, nil)
}

func (c *Client) StopMySQL(ctx context.Context, serverID int) error {
	path := fmt.Sprintf("/servers/%d/mysql/stop", serverID)
	return c.doRequest(ctx, http.MethodPost, path, nil, nil)
}

func (c *Client) RebootNginx(ctx context.Context, serverID int) error {
	path := fmt.Sprintf("/servers/%d/nginx/reboot", serverID)
	return c.doRequest(ctx, http.MethodPost, path, nil, nil)
}

func (c *Client) StopNginx(ctx context.Context, serverID int) error {
	path := fmt.Sprintf("/servers/%d/nginx/stop", serverID)
	return c.doRequest(ctx, http.MethodPost, path, nil, nil)
}

type testNginxResponse struct {
	Result string `json:"result"`
}

func (c *Client) TestNginx(ctx context.Context, serverID int) (*testNginxResponse, error) {
	path := fmt.Sprintf("/servers/%d/nginx/test", serverID)
	var resp testNginxResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) RebootPostgres(ctx context.Context, serverID int) error {
	path := fmt.Sprintf("/servers/%d/postgres/reboot", serverID)
	return c.doRequest(ctx, http.MethodPost, path, nil, nil)
}

func (c *Client) StopPostgres(ctx context.Context, serverID int) error {
	path := fmt.Sprintf("/servers/%d/postgres/stop", serverID)
	return c.doRequest(ctx, http.MethodPost, path, nil, nil)
}

type rebootPHPRequest struct {
	Version string `json:"version"`
}

func (c *Client) RebootPHP(ctx context.Context, serverID int, version string) error {
	path := fmt.Sprintf("/servers/%d/php/reboot", serverID)
	req := rebootPHPRequest{Version: version}
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}

type installBlackfireRequest struct {
	ServerID    string `json:"server_id"`
	ServerToken string `json:"server_token"`
}

func (c *Client) InstallBlackfire(ctx context.Context, serverID int, serverToken string) error {
	path := fmt.Sprintf("/servers/%d/blackfire/install", serverID)
	req := installBlackfireRequest{ServerID: fmt.Sprint(serverID), ServerToken: serverToken}
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}

func (c *Client) RemoveBlackfire(ctx context.Context, serverID int) error {
	path := fmt.Sprintf("/servers/%d/blackfire/remove", serverID)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

type installPapertrailRequest struct {
	Host string `json:"host"`
}

func (c *Client) InstallPapertrail(ctx context.Context, serverID int, host string) error {
	path := fmt.Sprintf("/servers/%d/papertrail/install", serverID)
	req := installPapertrailRequest{Host: host}
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}

func (c *Client) RemovePapertrail(ctx context.Context, serverID int) error {
	path := fmt.Sprintf("/servers/%d/papertrail/remove", serverID)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

func (c *Client) WaitForServerToBeReady(ctx context.Context, serverID int) error {
	// GetServerer and check the IsReady field every 10 seconds until it is true or the context is cancelled.

	for {
		server, err := c.GetServer(ctx, serverID)
		if err != nil {
			return err
		}
		if server.IsReady {
			return nil
		}
		select {
		case <-time.After(10 * time.Second):
			// continue polling
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
