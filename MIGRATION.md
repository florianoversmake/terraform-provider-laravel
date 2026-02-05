# Forge API Migration Progress

This document tracks the progress of migrating the Laravel Forge Terraform Provider from the deprecated v1 API to the new organization-scoped JSON:API format.

## Migration Status

### Completed

- [x] **Core Client** - Updated base URL, added organization support, JSON:API response parsing
- [x] **JSON:API Helpers** - Added `GetJsonApiList`, `GetJsonApiResource`, `PostJsonApi`, `PutJsonApi`, `unmarshalList`, `unmarshalSingle`
- [x] **Servers** - Migrated to `/orgs/{org}/servers`
- [x] **Sites** - Migrated to `/orgs/{org}/servers/{server}/sites`
- [x] **Workers** - Migrated to `/orgs/{org}/servers/{server}/background-processes` (renamed from daemons)
- [x] **Scheduled Jobs** - Migrated to `/orgs/{org}/servers/{server}/scheduled-jobs`
- [x] **SSH Keys** - Migrated to `/orgs/{org}/servers/{server}/ssh-keys`
- [x] **Recipes** - Migrated to `/orgs/{org}/recipes`
- [x] **Credentials** - Migrated to `/orgs/{org}/server-credentials`
- [x] **Certificates** - Migrated to domain-based management
- [x] **Deployment** - Migrated to new API structure
- [x] **Providers** - Split into separate `forge_providers.go` file (these endpoints do NOT require organization scope)
- [x] **Regions** - Refactored with efficient targeted lookup methods
- [x] **PHP Versions** - Migrated to new API structure
- [x] **Databases** - Migrated to `/orgs/{org}/servers/{server}/database/schemas`
- [x] **Database Users** - Migrated to `/orgs/{org}/servers/{server}/database/users`
- [x] **Backups** - Migrated to `/orgs/{org}/servers/{server}/database/backups`
- [x] **Security Rules** - Migrated to `/orgs/{org}/servers/{server}/sites/{site}/security-rules`
- [x] **Server Logs** - Migrated to `/orgs/{org}/servers/{server}/logs/{key}`
- [x] **Site Commands** - Migrated to `/orgs/{org}/servers/{server}/sites/{site}/commands`
- [x] **Webhooks** - Migrated to `/orgs/{org}/servers/{server}/sites/{site}/webhooks`
- [x] **Organizations** - Added new `organizations.go` with `ListOrganizations`, `GetOrganization`
- [x] **wordpress.go** - Added c.orgPath() wrapper (Note: endpoints may be deprecated in new API)

### New Terraform Data Sources

- [x] **forge_organizations** - List organizations the user has access to
- [x] **forge_servers** - List servers with filter support
- [x] **forge_sites** - List sites on a server with filter support
- [x] **forge_php_versions** - List PHP versions available on a server
- [x] **forge_providers** - List cloud providers (DigitalOcean, AWS, Hetzner, etc.)
- [x] **forge_regions** - List regions for a cloud provider
- [x] **forge_sizes** - List server sizes for a cloud provider

### Terraform Modules Created

- [x] **tf-test/modules/forge** - Reusable module for Forge server, site, and worker creation
- [x] **tf-test/modules/envoyer** - Reusable module for Envoyer project, server, and hook management
- [x] **tf-test/forge** - Example usage with dev overrides
- [x] **tf-test/envoyer** - Example usage with dev overrides

### Import Support Added

- [x] **EnvoyerDeploymentResource** - Added ImportState (format: project_id/deployment_id)

### Import Support Analysis

Resources WITHOUT ImportState (intentionally not added):
- **ForgeRecipeRunResource** - This is a one-time action resource that runs a recipe. It doesn't have a unique ID and cannot be imported.
- **ForgeCertificateSigningRequestInstallationResource** - This is a one-time action resource that installs a certificate. It doesn't have a unique ID and cannot be imported.

### Files Modified

- `internal/forge_client/client.go` - Added organization support, updated base URL
- `internal/forge_client/jsonapi.go` - Added JSON:API parsing helpers
- `internal/forge_client/forge_providers.go` - **NEW** - Provider API methods (no org scope required)
- `internal/forge_client/regions.go` - **REFACTORED** - Efficient targeted lookup methods instead of fetching all data
- `internal/forge_client/organizations.go` - **NEW** - Organization listing methods
- `internal/forge_client/security_rules.go` - **MIGRATED** - Updated to use JSON:API with c.orgPath()
- `internal/forge_client/server_logs.go` - **MIGRATED** - Updated to use JSON:API with key parameter
- `internal/forge_client/site_commands.go` - **MIGRATED** - Updated to use JSON:API helpers
- `internal/forge_client/webhooks.go` - **MIGRATED** - Updated to use JSON:API helpers
- `internal/forge_client/firewall.go` - **FIXED** - Port field changed from int to string to match API
- `internal/forge_client/nginx_templates.go` - **FIXED** - Corrected ListNginxTemplates endpoint path
- `internal/forge_client/wordpress.go` - **FIXED** - Added c.orgPath() wrapper
- `internal/forge_client/client_integration_test.go` - Updated tests, added 20+ new non-invasive integration tests
- `internal/provider/envoyer_deployment_resource.go` - **UPDATED** - Added ImportState support

### New Data Source Files

- `internal/provider/forge_organizations_data_source.go` - **NEW**
- `internal/provider/forge_servers_data_source.go` - **NEW**
- `internal/provider/forge_sites_data_source.go` - **NEW**
- `internal/provider/forge_php_versions_data_source.go` - **NEW**
- `internal/provider/forge_providers_data_source.go` - **NEW**
- `internal/provider/forge_regions_data_source.go` - **NEW**
- `internal/provider/forge_sizes_data_source.go` - **NEW**

### New Module Files

- `tf-test/modules/forge/main.tf` - **NEW** - Forge server module
- `tf-test/modules/forge/README.md` - **NEW** - Forge module documentation
- `tf-test/modules/envoyer/main.tf` - **NEW** - Envoyer project module
- `tf-test/modules/envoyer/README.md` - **NEW** - Envoyer module documentation
- `tf-test/forge/main.tf` - **NEW** - Forge example usage
- `tf-test/forge/terraform.tfvars.example` - **NEW** - Forge example variables
- `tf-test/envoyer/main.tf` - **NEW** - Envoyer example usage
- `tf-test/envoyer/terraform.tfvars.example` - **NEW** - Envoyer example variables

### API Endpoints That Do NOT Require Organization Scope

The following endpoints are accessed without the `/orgs/{org}` prefix:

- `GET /providers` - List all cloud providers
- `GET /providers/{provider}` - Get a specific provider
- `GET /providers/{provider}/regions` - List regions for a provider
- `GET /providers/{provider}/regions/{region}` - Get a specific region
- `GET /providers/{provider}/sizes` - List sizes for a provider
- `GET /providers/{provider}/sizes/{size}` - Get a specific size
- `GET /providers/{provider}/regions/{region}/sizes` - List sizes available in a region
- `GET /providers/{provider}/regions/{region}/sizes/{size}` - Get a specific size in a region

### Potentially Deprecated Endpoints

The following client files contain endpoints that may not exist in the new API structure (not found in openapi3.1.json):

- `phpmyadmin.go` - phpMyAdmin installation/uninstallation
- `wordpress.go` - WordPress installation/uninstallation
- `git_projects.go` - Git project management (repository is now a site attribute)

These have been updated to use `c.orgPath()` for consistency but may return 404 errors when called.

### Region/Size Lookup Refactoring

The region and size lookup methods have been refactored to use efficient targeted API calls instead of fetching all data:

**New Efficient Methods:**
- `GetProviderBySlug(slug)` - Get provider info by slug (e.g., "aws")
- `GetRegionByCode(providerSlug, regionCode)` - Get region by code (e.g., "eu-west-1")
- `GetRegionByName(providerSlug, regionName)` - Get region by name (e.g., "Ireland")
- `GetSizeByCode(providerSlug, sizeCode)` - Get size by code (e.g., "t3.small")
- `GetSizeByID(providerSlug, sizeID)` - Get size by ID

**Backward-Compatible Methods (now efficient):**
- `GetRegionIDByName` - 2 API calls instead of 100+
- `GetRegionNameByID` - 2 API calls instead of 100+
- `GetRegionSizeIDByName` - 4 API calls instead of 100+
- `GetRegionSizeNameByID` - 2 API calls instead of 100+
- `GetRegionSizeSizeByID` - 2 API calls instead of 100+

**Heavy Methods (deprecated, use with caution):**
- `GetAllRegionsWithSizes()` - Fetches ALL providers, regions, and sizes (many API calls, cached for 60s)
- `ListRegions()` - Alias for `GetAllRegionsWithSizes()` for backward compatibility

### Notes

- Provider endpoints do not require organization scope and can be called without setting `forge_organization`.
- Integration tests now run without hitting rate limits since lookup methods are efficient.

### Current Debugging

**Issue**: Server creation returns ID=0, causing WaitForServerToBeReady to fail with "No query results for model [Server] 0"

**Added debug logging** to `CreateServer` function in `internal/forge_client/server.go`:
- Logs the request body
- Logs the raw response body
- Logs the parsed Server ID

To rebuild and test:
```bash
cd /Users/flo/Code/tf-provider-laravel
go build -o terraform-provider-laravel
cd tf-test/forge
terraform apply
```

The debug output will show:
1. `[DEBUG] CreateServer request body: {...}`
2. `[DEBUG] CreateServer raw response: {...}`
3. `[DEBUG] CreateServer parsed response - Server ID: X`

This will reveal if:
- The API returns a different JSON structure than expected (`{"server": {...}}`)
- The ID is in a different field or format
- There's a JSON:API wrapper we need to handle
