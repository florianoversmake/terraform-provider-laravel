package forge_client

import (
	"context"
	"fmt"
	"net/http"
)

type SSLCertificate struct {
	ID            int64  `json:"id"`
	Domain        string `json:"domain"`
	RequestStatus string `json:"request_status"`
	CreatedAt     int64  `json:"created_at"`
	Existing      bool   `json:"existing"`
	Active        bool   `json:"active"`
	Type          string `json:"type,omitempty"`
}

type certificateResponse struct {
	Certificate SSLCertificate `json:"certificate"`
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

func (c *Client) CreateCertificate(ctx context.Context, serverID, siteID int, req CreateCertificateRequest) (*SSLCertificate, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/certificates", serverID, siteID)
	var res certificateResponse
	if err := c.doRequest(ctx, http.MethodPost, path, req, &res); err != nil {
		return nil, err
	}
	return &res.Certificate, nil
}

func (c *Client) ListCertificates(ctx context.Context, serverID, siteID int) ([]SSLCertificate, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/certificates", serverID, siteID)
	var res struct {
		Certificates []SSLCertificate `json:"certificates"`
	}
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return res.Certificates, nil
}

func (c *Client) GetCertificate(ctx context.Context, serverID, siteID, certID int) (*SSLCertificate, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/certificates/%d", serverID, siteID, certID)
	var res certificateResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}
	return &res.Certificate, nil
}

func (c *Client) GetCertificateCSR(ctx context.Context, serverID, siteID, certID int) (string, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/certificates/%d/csr", serverID, siteID, certID)

	csr, err := c.GetText(ctx, path)
	if err != nil {
		return "", err
	}
	return csr, nil
}

type InstallCertificateRequest struct {
	Certificate      string `json:"certificate"`
	AddIntermediates bool   `json:"add_intermediates"`
}

func (c *Client) InstallCertificate(ctx context.Context, serverID, siteID, certID int, req InstallCertificateRequest) error {
	path := fmt.Sprintf("/servers/%d/sites/%d/certificates/%d/install", serverID, siteID, certID)
	return c.doRequest(ctx, http.MethodPost, path, req, nil)
}

func (c *Client) ActivateCertificate(ctx context.Context, serverID, siteID, certID int) error {
	path := fmt.Sprintf("/servers/%d/sites/%d/certificates/%d/activate", serverID, siteID, certID)
	return c.doRequest(ctx, http.MethodPost, path, nil, nil)
}

func (c *Client) DeleteCertificate(ctx context.Context, serverID, siteID, certID int) error {
	path := fmt.Sprintf("/servers/%d/sites/%d/certificates/%d", serverID, siteID, certID)
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

type ObtainLetsencryptCertificate struct {
	Domain        string `json:"domain"`
	Type          string `json:"type"`
	RequestStatus string `json:"request_status"`
	Status        string `json:"status"`
	CreatedAt     string `json:"created_at"`
	Id            int64  `json:"id"`
	Existing      bool   `json:"existing"`
	Active        bool   `json:"active"`
}

type ObtainLetsencryptCertificateResponse struct {
	Certificate ObtainLetsencryptCertificate `json:"certificate"`
}

func (c *Client) ObtainLetsencryptCertificate(ctx context.Context, serverID, siteID int, req ObtainLetsencryptCertificateRequest) (*ObtainLetsencryptCertificate, error) {
	path := fmt.Sprintf("/servers/%d/sites/%d/letsencrypt", serverID, siteID)
	var res ObtainLetsencryptCertificateResponse
	if err := c.doRequest(ctx, http.MethodPost, path, req, &res); err != nil {
		return nil, err
	}
	return &res.Certificate, nil
}
