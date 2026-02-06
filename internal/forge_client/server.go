package forge_client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
	Revoked          *bool   `json:"revoked"`
	CreatedAt        string  `json:"created_at"`
	IsReady          bool    `json:"is_ready"`
	Tags             []Tag   `json:"tags"` //undocumented
	Network          []int64 `json:"network"`

	// Additional fields returned only on creation
	SudoPassword        string  `json:"sudo_password,omitempty"`
	DatabasePassword    *string `json:"database_password,omitempty"`
	MeilisearchPassword *string `json:"meilisearch_password,omitempty"`
	ProvisionCommand    *string `json:"provision_command,omitempty"`
}

type Tag struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

func (c *Client) ListServers(ctx context.Context) ([]Server, error) {
	items, err := c.GetJsonApiListAll(ctx, c.orgPath("/servers"))
	if err != nil {
		return nil, err
	}
	return unmarshalList(items, func(s *Server, id int64) { s.ID = id })
}

func (c *Client) GetServer(ctx context.Context, serverID int) (*Server, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d", serverID))
	var server Server
	_, err := c.GetJsonApi(ctx, path, &server)
	if err != nil {
		return nil, err
	}
	return &server, nil
}

func (c *Client) GetServerWithoutCache(ctx context.Context, serverID int) (*Server, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d", serverID))
	var server Server
	_, err := c.GetJsonApiWithoutCache(ctx, path, &server)
	if err != nil {
		return nil, err
	}
	return &server, nil
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
	Name          string   `json:"name"`
	Provider      string   `json:"provider"`
	Type          string   `json:"type"`
	UbuntuVersion string   `json:"ubuntu_version"`
	CredentialID  *int64   `json:"credential_id,omitempty"`
	PHPVersion    string   `json:"php_version,omitempty"`
	DatabaseType  *string  `json:"database_type,omitempty"`
	Database      *string  `json:"database,omitempty"`
	Circle        *int64   `json:"circle,omitempty"`
	Network       []int64  `json:"network,omitempty"`
	RecipeID      *int64   `json:"recipe_id,omitempty"`
	Tags          []string `json:"tags,omitempty"`

	// Provider-specific nested objects (new API format)
	AWS     *AWSServerConfig     `json:"aws,omitempty"`
	Ocean2  *Ocean2ServerConfig  `json:"ocean2,omitempty"`
	Hetzner *HetznerServerConfig `json:"hetzner,omitempty"`
	Vultr   *VultrServerConfig   `json:"vultr,omitempty"`
	Akamai  *AkamaiServerConfig  `json:"akamai,omitempty"`
	Laravel *LaravelServerConfig `json:"laravel,omitempty"`
	Custom  *CustomServerConfig  `json:"custom,omitempty"`
}

// AWSServerConfig holds AWS-specific server configuration.
type AWSServerConfig struct {
	RegionID   string `json:"region_id"`
	SizeID     string `json:"size_id"`
	VPCUUID    string `json:"vpc_uuid,omitempty"`
	SubnetUUID string `json:"subnet_uuid,omitempty"`
	DiskSize   string `json:"disk_size"`
}

// Ocean2ServerConfig holds DigitalOcean-specific server configuration.
type Ocean2ServerConfig struct {
	RegionID            string `json:"region_id"`
	SizeID              string `json:"size_id"`
	VPCUUID             string `json:"vpc_uuid,omitempty"`
	EnableWeeklyBackups string `json:"enable_weekly_backups,omitempty"`
}

// HetznerServerConfig holds Hetzner-specific server configuration.
type HetznerServerConfig struct {
	RegionID           string `json:"region_id"`
	SizeID             string `json:"size_id"`
	NetworkID          string `json:"network_id,omitempty"`
	EnableDailyBackups string `json:"enable_daily_backups,omitempty"`
}

// VultrServerConfig holds Vultr-specific server configuration.
type VultrServerConfig struct {
	RegionID  string `json:"region_id"`
	SizeID    string `json:"size_id"`
	NetworkID string `json:"network_id,omitempty"`
}

// AkamaiServerConfig holds Akamai/Linode-specific server configuration.
type AkamaiServerConfig struct {
	RegionID string `json:"region_id"`
	SizeID   string `json:"size_id"`
}

// LaravelServerConfig holds Laravel VPS-specific server configuration.
type LaravelServerConfig struct {
	RegionID string `json:"region_id"`
	SizeID   string `json:"size_id"`
}

// CustomServerConfig holds custom VPS server configuration.
type CustomServerConfig struct {
	IPAddress        string `json:"ip_address"`
	PrivateIPAddress string `json:"private_ip_address,omitempty"`
	SSHPort          string `json:"ssh_port,omitempty"`
	BehindNAT        string `json:"behind_nat,omitempty"`
	NATSSHPort       string `json:"nat_ssh_port,omitempty"`
}

type CreateServerResponse struct {
	Server              Server  `json:"-"` // Populated from JSON:API data.attributes
	SudoPassword        string  `json:"sudo_password"`
	DatabasePassword    *string `json:"database_password"`
	MeilisearchPassword *string `json:"meilisearch_password"`
	ProvisionCommand    *string `json:"provision_command"`
}

// createServerMeta represents the data.meta object returned by the create server endpoint.
// Passwords and provision commands are returned here, not in data.attributes.
type createServerMeta struct {
	SudoPassword        string  `json:"sudo_password"`
	DatabasePassword    *string `json:"database_password"`
	MeilisearchPassword *string `json:"meilisearch_password"`
	ProvisionCommand    *string `json:"provision_command"`
}

// CreateServer creates a new server. The API returns JSON:API format with
// server details in data.attributes and passwords in data.meta.
func (c *Client) CreateServer(ctx context.Context, req CreateServerRequest) (*CreateServerResponse, error) {
	path := c.orgPath("/servers")

	// Debug: log the request
	reqBytes, _ := json.Marshal(req)
	log.Printf("[DEBUG] CreateServer request body: %s", string(reqBytes))

	// Make the POST request and get the raw response
	rawResp, err := c.doRequestInternal(ctx, http.MethodPost, path, req)
	if err != nil {
		return nil, err
	}

	// Extract attributes and meta from data
	serverID, attrs, meta, err := extractJsonApiSingleWithMeta(rawResp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal server from data.attributes
	var server Server
	if len(attrs) > 0 {
		if err := json.Unmarshal(attrs, &server); err != nil {
			return nil, fmt.Errorf("failed to unmarshal server attributes: %w", err)
		}
	}
	server.ID = int64(serverID)
	log.Printf("[DEBUG] CreateServer - Server ID from JSON:API: %d", serverID)

	// Unmarshal passwords from data.meta
	var serverMeta createServerMeta
	if len(meta) > 0 {
		if err := json.Unmarshal(meta, &serverMeta); err != nil {
			log.Printf("[WARN] CreateServer - failed to unmarshal data.meta: %v", err)
		}
	}

	resp := &CreateServerResponse{
		Server:              server,
		SudoPassword:        serverMeta.SudoPassword,
		DatabasePassword:    serverMeta.DatabasePassword,
		MeilisearchPassword: serverMeta.MeilisearchPassword,
		ProvisionCommand:    serverMeta.ProvisionCommand,
	}

	return resp, nil
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
	path := c.orgPath(fmt.Sprintf("/servers/%d", serverID))
	var server Server
	_, err := c.PutJsonApi(ctx, path, req, &server)
	if err != nil {
		return nil, err
	}
	return &server, nil
}

func (c *Client) DeleteServer(ctx context.Context, serverID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d", serverID))
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

// serverActionRequest is the payload for server-level actions (reboot, power-cycle).
type serverActionRequest struct {
	Action string `json:"action"`
}

func (c *Client) RebootServer(ctx context.Context, serverID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/actions", serverID))
	req := serverActionRequest{Action: "reboot"}
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}

func (c *Client) RevokeServer(ctx context.Context, serverID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/revoke", serverID))
	return c.doRequest(ctx, http.MethodPost, path, nil, nil)
}

type ReconnectServerResponse struct {
	PublicKey string `json:"public_key"`
}

func (c *Client) ReconnectServer(ctx context.Context, serverID int) (*ReconnectServerResponse, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/reconnect", serverID))
	var resp ReconnectServerResponse
	if err := c.doRequest(ctx, http.MethodPost, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) ReactivateServer(ctx context.Context, serverID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/reactivate", serverID))
	return c.doRequest(ctx, http.MethodPost, path, nil, nil)
}

// serviceActionRequest is the payload for service-level actions (reboot, stop, start).
type serviceActionRequest struct {
	Action string `json:"action"`
}

func (c *Client) StartService(ctx context.Context, serverID int, service string) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/services/%s/actions", serverID, service))
	req := serviceActionRequest{Action: "start"}
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}

func (c *Client) StopService(ctx context.Context, serverID int, service string) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/services/%s/actions", serverID, service))
	req := serviceActionRequest{Action: "stop"}
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}

func (c *Client) RestartService(ctx context.Context, serverID int, service string) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/services/%s/actions", serverID, service))
	req := serviceActionRequest{Action: "reboot"}
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}

func (c *Client) RebootMySQL(ctx context.Context, serverID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/services/mysql/actions", serverID))
	req := serviceActionRequest{Action: "reboot"}
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}

func (c *Client) StopMySQL(ctx context.Context, serverID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/services/mysql/actions", serverID))
	req := serviceActionRequest{Action: "stop"}
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}

func (c *Client) RebootNginx(ctx context.Context, serverID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/services/nginx/actions", serverID))
	req := serviceActionRequest{Action: "reboot"}
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}

func (c *Client) StopNginx(ctx context.Context, serverID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/services/nginx/actions", serverID))
	req := serviceActionRequest{Action: "stop"}
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}

type testNginxResponse struct {
	Result string `json:"result"`
}

func (c *Client) TestNginx(ctx context.Context, serverID int) (*testNginxResponse, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/nginx/test", serverID))
	var resp testNginxResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) RebootPostgres(ctx context.Context, serverID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/services/postgres/actions", serverID))
	req := serviceActionRequest{Action: "reboot"}
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}

func (c *Client) StopPostgres(ctx context.Context, serverID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/services/postgres/actions", serverID))
	req := serviceActionRequest{Action: "stop"}
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}

func (c *Client) RebootPHP(ctx context.Context, serverID int, version string) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/services/php/actions", serverID))
	req := serviceActionRequest{Action: "reboot"}
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}

type installBlackfireRequest struct {
	ServerID    string `json:"server_id"`
	ServerToken string `json:"server_token"`
}

func (c *Client) InstallBlackfire(ctx context.Context, serverID int, serverToken string) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/blackfire/install", serverID))
	req := installBlackfireRequest{ServerID: fmt.Sprint(serverID), ServerToken: serverToken}
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}

func (c *Client) RemoveBlackfire(ctx context.Context, serverID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/blackfire/remove", serverID))
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

type installPapertrailRequest struct {
	Host string `json:"host"`
}

func (c *Client) InstallPapertrail(ctx context.Context, serverID int, host string) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/papertrail/install", serverID))
	req := installPapertrailRequest{Host: host}
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}

func (c *Client) RemovePapertrail(ctx context.Context, serverID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/papertrail/remove", serverID))
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

func (c *Client) WaitForServerToBeReady(ctx context.Context, serverID int) error {
	for {
		server, err := c.GetServerWithoutCache(ctx, serverID)
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
