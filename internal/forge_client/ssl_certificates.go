// Copyright (c) HashiCorp, Inc.

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
	CreatedAt     string `json:"created_at"`
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
	var res struct {
		Output string `json:"output"`
	}
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &res); err != nil {
		return "", err
	}
	return res.Output, nil
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
