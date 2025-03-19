package forge_client

import (
	"context"
	"os"
	"strconv"
	"testing"
)

var (
	apiKey      = os.Getenv("FORGE_API_KEY")
	serverIDStr = os.Getenv("FORGE_SERVER_ID")
	siteIDStr   = os.Getenv("FORGE_SITE_ID")
)

func getIntegrationClient(t *testing.T) *Client {
	if apiKey == "" {
		t.Skip("FORGE_API_KEY not set")
	}
	return NewClient(apiKey)
}

func getServerID(t *testing.T) int {
	if serverIDStr == "" {
		t.Skip("FORGE_SERVER_ID not set")
	}
	id, err := strconv.Atoi(serverIDStr)
	if err != nil {
		t.Fatalf("Invalid FORGE_SERVER_ID: %v", err)
	}
	return id
}

func getSiteID(t *testing.T) int {
	if siteIDStr == "" {
		t.Skip("FORGE_SITE_ID not set")
	}
	id, err := strconv.Atoi(siteIDStr)
	if err != nil {
		t.Fatalf("Invalid FORGE_SITE_ID: %v", err)
	}
	return id
}

func TestGetUser(t *testing.T) {
	client := getIntegrationClient(t)
	user, err := client.GetUser(context.Background())
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}
	if user.Email == "" {
		t.Error("Expected non-empty email in user response")
	}
}

func TestListServers(t *testing.T) {
	client := getIntegrationClient(t)
	servers, err := client.ListServers(context.Background())
	if err != nil {
		t.Fatalf("ListServers failed: %v", err)
	}
	t.Logf("Found %d servers", len(servers))
}

func TestListPHPVersions(t *testing.T) {
	client := getIntegrationClient(t)
	serverID := getServerID(t)
	versions, err := client.ListPHPVersions(context.Background(), serverID)
	if err != nil {
		t.Fatalf("ListPHPVersions failed: %v", err)
	}
	if len(versions) == 0 {
		t.Error("Expected at least one PHP version")
	}
}

func TestListRegions(t *testing.T) {
	client := getIntegrationClient(t)
	regions, err := client.ListRegions(context.Background())
	if err != nil {
		t.Fatalf("ListRegions failed: %v", err)
	}
	if len(regions) == 0 {
		t.Error("Expected non-empty regions map")
	}
}

func TestListCredentials(t *testing.T) {
	client := getIntegrationClient(t)
	creds, err := client.ListCredentials(context.Background())
	if err != nil {
		t.Fatalf("ListCredentials failed: %v", err)
	}
	t.Logf("Found %d credentials", len(creds))
}

func TestSetDeploymentFailureEmails(t *testing.T) {
	client := getIntegrationClient(t)
	serverID := getServerID(t)
	siteID := getSiteID(t)
	emails := []string{"test@example.com"}
	if err := client.SetDeploymentFailureEmails(context.Background(), serverID, siteID, emails); err != nil {
		t.Fatalf("SetDeploymentFailureEmails failed: %v", err)
	}
}

func TestGetRegionIDByName(t *testing.T) {
	client := getIntegrationClient(t)
	region, err := client.GetRegionIDByName(context.Background(), "aws", "Ireland")
	if err != nil {
		t.Fatalf("ListRegions failed: %v", err)
	}
	if region != "eu-west-1" {
		t.Errorf("Expected region ID eu-west-1, got %s", region)
	}
}

func TestGetRegionNameByID(t *testing.T) {
	client := getIntegrationClient(t)
	region, err := client.GetRegionNameByID(context.Background(), "aws", "eu-west-1")
	if err != nil {
		t.Fatalf("ListRegions failed: %v", err)
	}
	if region != "Ireland" {
		t.Errorf("Expected region name Ireland, got %s", region)
	}
}

func TestGetRegionSizeIDByName(t *testing.T) {
	client := getIntegrationClient(t)
	size, err := client.GetRegionSizeIDByName(context.Background(), "aws", "eu-west-1", "2GB RAM (t3.small) (64-bit x86) - 2 vCPUs")
	if err != nil {
		t.Fatalf("ListRegions failed: %v", err)
	}
	if size != "10" {
		t.Errorf("Expected size ID 10, got %s", size)
	}
}

func TestGetRegionSizeNameByID(t *testing.T) {
	client := getIntegrationClient(t)
	size, err := client.GetRegionSizeNameByID(context.Background(), "aws", "eu-west-1", "10")
	if err != nil {
		t.Fatalf("ListRegions failed: %v", err)
	}
	if size != "2GB RAM (t3.small) (64-bit x86) - 2 vCPUs" {
		t.Errorf("Expected size name 2GB RAM (t3.small) (64-bit x86) - 2 vCPUs, got %s", size)
	}
}

func TestGetRegionSizeIDBySize(t *testing.T) {
	client := getIntegrationClient(t)
	size, err := client.GetRegionSizeIDBySize(context.Background(), "aws", "eu-west-1", "t3.small")
	if err != nil {
		t.Fatalf("ListRegions failed: %v", err)
	}
	if size != "10" {
		t.Errorf("Expected size ID 10, got %s", size)
	}
}
