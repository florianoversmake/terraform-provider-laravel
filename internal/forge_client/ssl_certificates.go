package forge_client

import (
	"context"
	"fmt"
	"net/http"
)

// SSLCertificate represents a certificate resource in the new API.
// In the new API, certificates are managed per domain record rather than per site.
type SSLCertificate struct {
	ID            int64  `json:"id"`
	Type          string `json:"type"`
	RequestStatus string `json:"request_status"`
	Status        string `json:"status"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type CreateCertificateRequest struct {
	Type          string `json:"type"`
	Domain        string `json:"domain,omitempty"`
	Country       string `json:"country,omitempty"`
	State         string `json:"state,omitempty"`
	City          string `json:"city,omitempty"`
	Organization  string `json:"organization,omitempty"`
	Department    string `json:"department,omitempty"`
	Key           string `json:"key,omitempty"`
	Certificate   string `json:"certificate,omitempty"`
	CertificateID *int   `json:"certificate_id,omitempty"`
}

// CreateCertificate creates a certificate on a domain record.
// In the new API, the domainID parameter identifies the domain record.
func (c *Client) CreateCertificate(ctx context.Context, serverID, siteID, domainID int, req CreateCertificateRequest) (*SSLCertificate, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/domains/%d/certificate", serverID, siteID, domainID))
	var cert SSLCertificate
	id, err := c.PostJsonApi(ctx, path, req, &cert)
	if err != nil {
		return nil, err
	}
	cert.ID = int64(id)
	return &cert, nil
}

// GetCertificate retrieves the certificate for a domain record.
// In the new API, certificates are per-domain rather than per-site.
func (c *Client) GetCertificate(ctx context.Context, serverID, siteID, domainID int) (*SSLCertificate, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/domains/%d/certificate", serverID, siteID, domainID))
	var cert SSLCertificate
	id, err := c.GetJsonApi(ctx, path, &cert)
	if err != nil {
		return nil, err
	}
	cert.ID = int64(id)
	return &cert, nil
}

type InstallCertificateRequest struct {
	Certificate      string `json:"certificate"`
	AddIntermediates bool   `json:"add_intermediates"`
}

// InstallCertificate performs a certificate action on a domain record.
func (c *Client) InstallCertificate(ctx context.Context, serverID, siteID, domainID int, req InstallCertificateRequest) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/domains/%d/certificate/actions", serverID, siteID, domainID))
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}

// DeleteCertificate deletes the certificate for a domain record.
func (c *Client) DeleteCertificate(ctx context.Context, serverID, siteID, domainID int) error {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/domains/%d/certificate", serverID, siteID, domainID))
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

type ObtainLetsencryptCertificateDNSProvider struct {
	Type                  string `json:"type"`
	CloudflareAPIToken    string `json:"cloudflare_api_token,omitempty"`
	Route53Key            string `json:"route53_key,omitempty"`
	Route53Secret         string `json:"route53_secret,omitempty"`
	DigitalOceanToken     string `json:"digitalocean_token,omitempty"`
	DNSSimpleToken        string `json:"dnssimple_token,omitempty"`
	LinodeToken           string `json:"linode_token,omitempty"`
	OVHEndpoint           string `json:"ovh_endpoint,omitempty"`
	OVHAppKey             string `json:"ovh_app_key,omitempty"`
	OVHAppSecret          string `json:"ovh_app_secret,omitempty"`
	OVHConsumerKey        string `json:"ovh_consumer_key,omitempty"`
	GoogleCredentialsFile string `json:"google_credentials_file,omitempty"`
}

type ObtainLetsencryptCertificateRequest struct {
	Domains     []string                                `json:"domains"`
	DNSProvider ObtainLetsencryptCertificateDNSProvider `json:"dns_provider"`
}

// ObtainLetsencryptCertificate creates a Let's Encrypt certificate on a domain record.
func (c *Client) ObtainLetsencryptCertificate(ctx context.Context, serverID, siteID, domainID int, req ObtainLetsencryptCertificateRequest) (*SSLCertificate, error) {
	path := c.orgPath(fmt.Sprintf("/servers/%d/sites/%d/domains/%d/certificate", serverID, siteID, domainID))
	var cert SSLCertificate
	id, err := c.PostJsonApi(ctx, path, req, &cert)
	if err != nil {
		return nil, err
	}
	cert.ID = int64(id)
	return &cert, nil
}
