package forge_client

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
)

var (
	sharedOrgClient     *Client
	sharedOrgClientOnce sync.Once

	sharedNoOrgClient     *Client
	sharedNoOrgClientOnce sync.Once

	sharedLogger = log.New(os.Stderr, "[forge] ", log.Lmicroseconds)
)

func getIntegrationClient(t *testing.T) *Client {
	apiKey := os.Getenv("FORGE_API_KEY")
	org := os.Getenv("FORGE_ORG")

	if apiKey == "" {
		t.Skip("FORGE_API_KEY not set")
	}

	if org == "" {
		t.Skip("FORGE_ORG not set")
	}

	sharedOrgClientOnce.Do(func() {
		sharedOrgClient = NewClient(apiKey).WithCache(NewMemoryCache()).WithDebugLog(sharedLogger).WithOrganization(org)
	})
	t.Cleanup(func() {
		t.Logf("HTTP requests: %d, cache hits: %d", sharedOrgClient.RequestCount.Load(), sharedOrgClient.CacheHits.Load())
	})
	return sharedOrgClient
}

// getIntegrationClientWithoutOrg returns a client without organization scope.
// Use this for endpoints that don't require an organization (e.g., /providers).
func getIntegrationClientWithoutOrg(t *testing.T) *Client {
	apiKey := os.Getenv("FORGE_API_KEY")

	if apiKey == "" {
		t.Skip("FORGE_API_KEY not set")
	}

	sharedNoOrgClientOnce.Do(func() {
		sharedNoOrgClient = NewClient(apiKey).WithCache(NewMemoryCache()).WithDebugLog(sharedLogger)
	})
	t.Cleanup(func() {
		t.Logf("HTTP requests: %d, cache hits: %d", sharedNoOrgClient.RequestCount.Load(), sharedNoOrgClient.CacheHits.Load())
	})
	return sharedNoOrgClient
}

func getServerID(t *testing.T) int {
	serverIDStr := os.Getenv("FORGE_SERVER_ID")
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
	siteIDStr := os.Getenv("FORGE_SITE_ID")
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

func TestListProviders(t *testing.T) {
	// Note: Providers endpoint does not require organization scope
	client := getIntegrationClientWithoutOrg(t)
	providers, err := client.ListProviders(context.Background())
	if err != nil {
		t.Fatalf("ListProviders failed: %v", err)
	}
	if len(providers) == 0 {
		t.Error("Expected at least one provider")
	}
	t.Logf("Found %d providers", len(providers))
	for _, p := range providers {
		t.Logf("  Provider: id=%d, slug=%q, name=%q", p.ID, p.Slug, p.Name)
	}
}

func TestListProviderRegions(t *testing.T) {
	// Note: Providers endpoint does not require organization scope
	client := getIntegrationClientWithoutOrg(t)
	// AWS provider ID is typically 1
	regions, err := client.ListProviderRegions(context.Background(), 1)
	if err != nil {
		t.Fatalf("ListProviderRegions failed: %v", err)
	}
	if len(regions) == 0 {
		t.Error("Expected at least one region")
	}
	t.Logf("Found %d regions for provider 1 (AWS)", len(regions))
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
	// NOTE: This endpoint has been removed from the new Forge API
	t.Skip("deployment-failure-emails endpoint no longer exists in the new API")
}

func TestGetRegionIDByName(t *testing.T) {
	// Uses efficient targeted API calls (only 2 calls: providers + provider regions)
	client := getIntegrationClientWithoutOrg(t)
	region, err := client.GetRegionIDByName(context.Background(), "aws", "Ireland")
	if err != nil {
		t.Fatalf("GetRegionIDByName failed: %v", err)
	}
	if region != "eu-west-1" {
		t.Errorf("Expected region ID eu-west-1, got %s", region)
	}
}

func TestGetRegionNameByID(t *testing.T) {
	// Uses efficient targeted API calls (only 2 calls: providers + provider regions)
	client := getIntegrationClientWithoutOrg(t)
	region, err := client.GetRegionNameByID(context.Background(), "aws", "eu-west-1")
	if err != nil {
		t.Fatalf("GetRegionNameByID failed: %v", err)
	}
	if region != "Ireland" {
		t.Errorf("Expected region name Ireland, got %s", region)
	}
}

func TestGetRegionSizeIDByName(t *testing.T) {
	// Uses efficient targeted API calls (providers + regions + sizes + region sizes)
	client := getIntegrationClientWithoutOrg(t)

	// First, get sizes to find a valid size name
	provider, err := client.GetProviderBySlug(context.Background(), "aws")
	if err != nil {
		t.Fatalf("GetProviderBySlug failed: %v", err)
	}
	sizes, err := client.ListProviderSizes(context.Background(), provider.ID)
	if err != nil {
		t.Fatalf("ListProviderSizes failed: %v", err)
	}
	if len(sizes) == 0 {
		t.Skip("No sizes found for AWS provider")
	}

	// Find a t3.small size
	var testSize *ProviderSizeInfo
	for i := range sizes {
		if sizes[i].Code == "t3.small" {
			testSize = &sizes[i]
			break
		}
	}
	if testSize == nil {
		t.Skip("t3.small size not found for AWS provider")
	}

	// Now test the lookup by name
	sizeID, err := client.GetRegionSizeIDByName(context.Background(), "aws", "eu-west-1", testSize.Name)
	if err != nil {
		t.Fatalf("GetRegionSizeIDByName failed: %v", err)
	}
	if sizeID == "" {
		t.Errorf("Expected non-empty size ID for %q", testSize.Name)
	}
	t.Logf("Size %q has ID %s", testSize.Name, sizeID)
}

func TestGetRegionSizeNameByID(t *testing.T) {
	// Uses efficient targeted API calls (providers + size lookup)
	client := getIntegrationClientWithoutOrg(t)

	// First, get a valid size ID
	provider, err := client.GetProviderBySlug(context.Background(), "aws")
	if err != nil {
		t.Fatalf("GetProviderBySlug failed: %v", err)
	}
	sizes, err := client.ListProviderSizes(context.Background(), provider.ID)
	if err != nil {
		t.Fatalf("ListProviderSizes failed: %v", err)
	}
	if len(sizes) == 0 {
		t.Skip("No sizes found for AWS provider")
	}

	// Use the first size
	testSize := sizes[0]
	testSizeID := strconv.FormatInt(testSize.ID, 10)

	// Now test the lookup by ID
	sizeName, err := client.GetRegionSizeNameByID(context.Background(), "aws", "eu-west-1", testSizeID)
	if err != nil {
		t.Fatalf("GetRegionSizeNameByID failed: %v", err)
	}
	if sizeName != testSize.Name {
		t.Errorf("Expected size name %q, got %q", testSize.Name, sizeName)
	}
	t.Logf("Size ID %s has name %q", testSizeID, sizeName)
}

func TestGetRegionSizeIDBySize(t *testing.T) {
	// Uses efficient targeted API calls (providers + sizes)
	client := getIntegrationClientWithoutOrg(t)

	// Test the lookup by known size code directly — avoids redundant ListProviderSizes call
	// which would double the paginated API requests and trigger rate limiting timeouts.
	sizeID, err := client.GetRegionSizeIDBySize(context.Background(), "aws", "eu-west-1", "t3.small")
	if err != nil {
		t.Fatalf("GetRegionSizeIDBySize failed: %v", err)
	}
	if sizeID == "" {
		t.Fatal("Expected non-empty size ID for t3.small")
	}
	t.Logf("Size code t3.small has ID %s", sizeID)
}

// ============================================================
// Additional Non-Invasive Integration Tests (Read-only operations)
// ============================================================

func TestListSites(t *testing.T) {
	client := getIntegrationClient(t)
	serverID := getServerID(t)
	sites, err := client.ListSites(context.Background(), serverID)
	if err != nil {
		t.Fatalf("ListSites failed: %v", err)
	}
	t.Logf("Found %d sites on server %d", len(sites), serverID)
	for _, site := range sites {
		t.Logf("  Site: id=%d, name=%q, status=%q", site.ID, site.Name, site.Status)
	}
}

func TestGetSite(t *testing.T) {
	client := getIntegrationClient(t)
	serverID := getServerID(t)
	siteID := getSiteID(t)
	site, err := client.GetSite(context.Background(), serverID, siteID)
	if err != nil {
		t.Fatalf("GetSite failed: %v", err)
	}
	if site.ID != int64(siteID) {
		t.Errorf("Expected site ID %d, got %d", siteID, site.ID)
	}
	t.Logf("Site: id=%d, name=%q, status=%q, php_version=%q", site.ID, site.Name, site.Status, site.PHPVersion)
}

func TestGetServer(t *testing.T) {
	client := getIntegrationClient(t)
	serverID := getServerID(t)
	server, err := client.GetServer(context.Background(), serverID)
	if err != nil {
		t.Fatalf("GetServer failed: %v", err)
	}
	if server.ID != int64(serverID) {
		t.Errorf("Expected server ID %d, got %d", serverID, server.ID)
	}
	t.Logf("Server: id=%d, name=%q, provider=%q, region=%q", server.ID, server.Name, server.Provider, server.Region)
}

func TestListRecipes(t *testing.T) {
	client := getIntegrationClient(t)
	recipes, err := client.ListRecipes(context.Background())
	if err != nil {
		t.Fatalf("ListRecipes failed: %v", err)
	}
	t.Logf("Found %d recipes", len(recipes))
	for _, recipe := range recipes {
		t.Logf("  Recipe: id=%d, name=%q, user=%q", recipe.ID, recipe.Name, recipe.User)
	}
}

func TestListSSHKeys(t *testing.T) {
	client := getIntegrationClient(t)
	serverID := getServerID(t)
	keys, err := client.ListSSHKeys(context.Background(), serverID)
	if err != nil {
		t.Fatalf("ListSSHKeys failed: %v", err)
	}
	t.Logf("Found %d SSH keys on server %d", len(keys), serverID)
	for _, key := range keys {
		t.Logf("  SSH Key: id=%d, name=%q, user=%q, status=%q", key.ID, key.Name, key.User, key.Status)
	}
}

func TestListJobs(t *testing.T) {
	client := getIntegrationClient(t)
	serverID := getServerID(t)
	jobs, err := client.ListJobs(context.Background(), serverID)
	if err != nil {
		t.Fatalf("ListJobs failed: %v", err)
	}
	t.Logf("Found %d scheduled jobs on server %d", len(jobs), serverID)
	for _, job := range jobs {
		t.Logf("  Job: id=%d, command=%q, user=%q, frequency=%q", job.ID, job.Command, job.User, job.Frequency)
	}
}

func TestListWorkers(t *testing.T) {
	client := getIntegrationClient(t)
	serverID := getServerID(t)
	siteID := getSiteID(t)
	workers, err := client.ListWorkers(context.Background(), serverID, siteID)
	if err != nil {
		t.Fatalf("ListWorkers failed: %v", err)
	}
	t.Logf("Found %d workers (background processes) on server %d", len(workers), serverID)
	for _, worker := range workers {
		t.Logf("  Worker: id=%d, command=%q, user=%q, status=%q", worker.ID, worker.Command, worker.User, worker.Status)
	}
}

func TestListDatabases(t *testing.T) {
	client := getIntegrationClient(t)
	serverID := getServerID(t)
	databases, err := client.ListDatabases(context.Background(), serverID)
	if err != nil {
		if strings.Contains(err.Error(), "status=404") {
			t.Skip("Server does not have a database installed")
		}
		t.Fatalf("ListDatabases failed: %v", err)
	}
	t.Logf("Found %d databases on server %d", len(databases), serverID)
	for _, db := range databases {
		t.Logf("  Database: id=%d, name=%q, status=%q", db.ID, db.Name, db.Status)
	}
}

func TestListDatabaseUsers(t *testing.T) {
	client := getIntegrationClient(t)
	serverID := getServerID(t)
	users, err := client.ListDatabaseUsers(context.Background(), serverID)
	if err != nil {
		if strings.Contains(err.Error(), "status=404") {
			t.Skip("Server does not have a database installed")
		}
		t.Fatalf("ListDatabaseUsers failed: %v", err)
	}
	t.Logf("Found %d database users on server %d", len(users), serverID)
	for _, user := range users {
		t.Logf("  Database User: id=%d, name=%q, status=%q", user.ID, user.Name, user.Status)
	}
}

func TestListFirewallRules(t *testing.T) {
	client := getIntegrationClient(t)
	serverID := getServerID(t)
	rules, err := client.ListFirewallRules(context.Background(), serverID)
	if err != nil {
		t.Fatalf("ListFirewallRules failed: %v", err)
	}
	t.Logf("Found %d firewall rules on server %d", len(rules), serverID)
	for _, rule := range rules {
		t.Logf("  Firewall Rule: id=%d, name=%q, port=%s, type=%q", rule.ID, rule.Name, rule.Port, rule.Type)
	}
}

func TestListMonitors(t *testing.T) {
	client := getIntegrationClient(t)
	serverID := getServerID(t)
	monitors, err := client.ListMonitors(context.Background(), serverID)
	if err != nil {
		t.Fatalf("ListMonitors failed: %v", err)
	}
	t.Logf("Found %d monitors on server %d", len(monitors), serverID)
	for _, monitor := range monitors {
		t.Logf("  Monitor: id=%d, type=%q, status=%q, state=%q", monitor.ID, monitor.Type, monitor.Status, monitor.State)
	}
}

func TestListRedirectRules(t *testing.T) {
	client := getIntegrationClient(t)
	serverID := getServerID(t)
	siteID := getSiteID(t)
	rules, err := client.ListRedirectRules(context.Background(), serverID, siteID)
	if err != nil {
		t.Fatalf("ListRedirectRules failed: %v", err)
	}
	t.Logf("Found %d redirect rules on site %d", len(rules), siteID)
	for _, rule := range rules {
		t.Logf("  Redirect Rule: id=%d, from=%q, to=%q, type=%q", rule.ID, rule.From, rule.To, rule.Type)
	}
}

func TestListNginxTemplates(t *testing.T) {
	client := getIntegrationClient(t)
	serverID := getServerID(t)
	templates, err := client.ListNginxTemplates(context.Background(), serverID)
	if err != nil {
		if strings.Contains(err.Error(), "status=404") {
			t.Skip("No nginx templates found on server")
		}
		t.Fatalf("ListNginxTemplates failed: %v", err)
	}
	t.Logf("Found %d nginx templates on server %d", len(templates), serverID)
	for _, tmpl := range templates {
		t.Logf("  Nginx Template: id=%d, name=%q", tmpl.ID, tmpl.Name)
	}
}

func TestListDeployments(t *testing.T) {
	client := getIntegrationClient(t)
	serverID := getServerID(t)
	siteID := getSiteID(t)
	deployments, err := client.ListDeployments(context.Background(), serverID, siteID)
	if err != nil {
		t.Fatalf("ListDeployments failed: %v", err)
	}
	t.Logf("Found %d deployments on site %d", len(deployments), siteID)
	for _, d := range deployments {
		t.Logf("  Deployment: id=%d, status=%q, type=%q", d.ID, d.Status, d.Type)
	}
}

func TestListBackupConfigurations(t *testing.T) {
	client := getIntegrationClient(t)
	serverID := getServerID(t)
	backups, err := client.ListBackupConfigurations(context.Background(), serverID)
	if err != nil {
		t.Fatalf("ListBackupConfigurations failed: %v", err)
	}
	t.Logf("Found %d backup configurations on server %d", len(backups), serverID)
	for _, b := range backups {
		t.Logf("  Backup Config: id=%d, provider=%q, status=%q", b.ID, b.Provider, b.Status)
	}
}

func TestListWebhooks(t *testing.T) {
	client := getIntegrationClient(t)
	serverID := getServerID(t)
	siteID := getSiteID(t)
	webhooks, err := client.ListWebhooks(context.Background(), serverID, siteID)
	if err != nil {
		t.Fatalf("ListWebhooks failed: %v", err)
	}
	t.Logf("Found %d webhooks on site %d", len(webhooks), siteID)
	for _, w := range webhooks {
		t.Logf("  Webhook: id=%d, url=%q", w.ID, w.URL)
	}
}

func TestListSiteCommands(t *testing.T) {
	client := getIntegrationClient(t)
	serverID := getServerID(t)
	siteID := getSiteID(t)
	commands, err := client.ListSiteCommands(context.Background(), serverID, siteID)
	if err != nil {
		t.Fatalf("ListSiteCommands failed: %v", err)
	}
	t.Logf("Found %d command history entries on site %d", len(commands), siteID)
	for _, cmd := range commands {
		t.Logf("  Command: id=%d, command=%q, status=%q", cmd.ID, cmd.Command, cmd.Status)
	}
}

func TestListSecurityRules(t *testing.T) {
	client := getIntegrationClient(t)
	serverID := getServerID(t)
	siteID := getSiteID(t)
	rules, err := client.ListSecurityRules(context.Background(), serverID, siteID)
	if err != nil {
		t.Fatalf("ListSecurityRules failed: %v", err)
	}
	t.Logf("Found %d security rules on site %d", len(rules), siteID)
	for _, rule := range rules {
		t.Logf("  Security Rule: id=%d, name=%q, status=%q", rule.ID, rule.Name, rule.Status)
	}
}

func TestGetDeploymentScript(t *testing.T) {
	client := getIntegrationClient(t)
	serverID := getServerID(t)
	siteID := getSiteID(t)
	script, err := client.GetDeploymentScript(context.Background(), serverID, siteID)
	if err != nil {
		t.Fatalf("GetDeploymentScript failed: %v", err)
	}
	t.Logf("Deployment script length: %d characters", len(script))
}

func TestListOrganizations(t *testing.T) {
	client := getIntegrationClient(t)
	orgs, err := client.ListOrganizations(context.Background())
	if err != nil {
		t.Fatalf("ListOrganizations failed: %v", err)
	}
	if len(orgs) == 0 {
		t.Fatalf("Expected at least one organization")
	}
	t.Logf("Found %d organizations", len(orgs))
	for _, org := range orgs {
		t.Logf("  Organization: id=%d, name=%q, slug=%q, owner=%v", org.ID, org.Name, org.Slug, org.Owner)
	}
}

func TestGetOrganization(t *testing.T) {
	client := getIntegrationClient(t)
	// First get the list to find a valid slug
	orgs, err := client.ListOrganizations(context.Background())
	if err != nil {
		t.Fatalf("ListOrganizations failed: %v", err)
	}
	if len(orgs) == 0 {
		t.Skip("No organizations found to test")
	}
	slug := orgs[0].Slug

	org, err := client.GetOrganization(context.Background(), slug)
	if err != nil {
		t.Fatalf("GetOrganization failed: %v", err)
	}
	t.Logf("Organization: id=%d, name=%q, slug=%q, owner=%v", org.ID, org.Name, org.Slug, org.Owner)
}

func TestListProviderSizes(t *testing.T) {
	client := getIntegrationClient(t)
	providers, err := client.ListProviders(context.Background())
	if err != nil {
		t.Fatalf("ListProviders failed: %v", err)
	}
	if len(providers) == 0 {
		t.Skip("No providers found")
	}
	providerID := providers[0].ID

	sizes, err := client.ListProviderSizes(context.Background(), providerID)
	if err != nil {
		t.Fatalf("ListProviderSizes failed: %v", err)
	}
	t.Logf("Found %d sizes for provider %d (%s)", len(sizes), providerID, providers[0].Name)
	for i, size := range sizes {
		if i >= 5 {
			t.Logf("  ... and %d more", len(sizes)-5)
			break
		}
		t.Logf("  Size: id=%d, code=%q, name=%q, cpus=%d, ram=%dMB", size.ID, size.Code, size.Name, size.Cpus, size.Ram)
	}
}
