# Upgrade Guide: v0.1.x to v0.2.0

This guide covers the breaking changes introduced in v0.2.0 when migrating from the deprecated Laravel Forge API v1 to the new organization-scoped JSON:API.

## Why This Migration?

Laravel Forge API v1 is deprecated and will be **discontinued on March 31, 2026**. The new API uses:

- **Organization-scoped endpoints** (`/orgs/{organization}/...`)
- **JSON:API response format** (`application/vnd.api+json`)
- **Restructured endpoints** with several path and model changes

## Provider Configuration

### Required: Add `forge_organization`

All Forge API operations now require an organization slug. Add this to your provider configuration:

```hcl
provider "laravel" {
  forge_api_token    = "your-forge-token"
  forge_organization = "your-organization-slug"
}
```

You can find your organization slug in the Forge dashboard URL: `https://forge.laravel.com/orgs/{slug}/servers`.

### Base URL Change

The default `forge_base_url` changed from `https://forge.laravel.com/api/v1` to `https://forge.laravel.com/api`. If you had explicitly set the base URL, update it accordingly.

## Resource Changes

### `laravel_forge_certificate`

**New required attribute:** `domain_id`

Certificates are now managed per domain record instead of per site. You need to provide the domain ID.

```hcl
# Before (v0.1.x)
resource "laravel_forge_certificate" "example" {
  server_id   = 1234
  site_id     = 5678
  key         = "..."
  certificate = "..."
  active      = true
}

# After (v0.2.0)
resource "laravel_forge_certificate" "example" {
  server_id   = 1234
  site_id     = 5678
  domain_id   = 9012  # NEW: required
  key         = "..."
  certificate = "..."
  # Removed: active, domain, existing
}
```

**Removed attributes:** `active`, `domain`, `existing`
**Type change:** `created_at` changed from Number to String

### `laravel_forge_certificate_signing_request`

**New required attribute:** `domain_id`

```hcl
# Before (v0.1.x)
resource "laravel_forge_certificate_signing_request" "example" {
  server_id = 1234
  site_id   = 5678
  domain    = "example.com"
  # ...
}

# After (v0.2.0)
resource "laravel_forge_certificate_signing_request" "example" {
  server_id = 1234
  site_id   = 5678
  domain_id = 9012  # NEW: required
  domain    = "example.com"
  # ...
}
```

**Removed attributes:** `active`, `existing`
**Type change:** `created_at` changed from Number to String

### `laravel_forge_certificate_signing_request_installation`

**Replaced:** `certificate_signing_request_id` with `domain_id`

```hcl
# Before (v0.1.x)
resource "laravel_forge_certificate_signing_request_installation" "example" {
  certificate_signing_request_id = laravel_forge_certificate_signing_request.example.id
  server_id                      = 1234
  site_id                        = 5678
  certificate                    = "..."
}

# After (v0.2.0)
resource "laravel_forge_certificate_signing_request_installation" "example" {
  domain_id   = 9012  # Replaces certificate_signing_request_id
  server_id   = 1234
  site_id     = 5678
  certificate = "..."
}
```

### `laravel_forge_site`

**Removed attributes:** `database`, `nginx_template`

These attributes are no longer part of the site creation/update API. If you were using them, remove them from your configuration.

### `laravel_forge_worker`

Workers are now managed as generic background processes in the new API. The resource continues to accept the same queue-worker-specific attributes (connection, timeout, sleep, delay, etc.) but internally constructs a `php artisan queue:work` command from them.

**Behavioral change:** The API no longer returns individual queue-worker fields. During `terraform plan`, you may see changes on the first run after upgrading as the state is reconciled.

### `laravel_forge_scheduled_job`

No configuration changes required. Internally, the individual time fields (minute, hour, day, month, weekday) are now combined into a cron expression when creating jobs via the API.

**Type change:** `created_at` is now nullable.

### `laravel_forge_ssh_key`

No configuration changes required. The API field `username` was renamed to `user` internally.

### `laravel_forge_credentials` (data source)

No configuration changes required. The API field `type` was renamed to `provider` internally.

## State Migration

After upgrading, you may need to update your Terraform state for resources that have new required attributes (especially certificate resources that now require `domain_id`). Options:

1. **Import with new IDs:** Remove the resource from state and re-import with the new composite ID format:
   ```bash
   terraform state rm laravel_forge_certificate.example
   terraform import laravel_forge_certificate.example "server_id:site_id:domain_id"
   ```

2. **Recreate resources:** For certificate resources, it may be simpler to destroy and recreate them with the new `domain_id` attribute.

## API Endpoint Changes Reference

| Old Path (v1) | New Path |
|---|---|
| `/servers/{id}/daemons` | `/orgs/{org}/servers/{id}/background-processes` |
| `/servers/{id}/jobs` | `/orgs/{org}/servers/{id}/scheduled-jobs` |
| `/servers/{id}/keys` | `/orgs/{org}/servers/{id}/ssh-keys` |
| `/credentials` | `/orgs/{org}/server-credentials` |
| `/servers/{id}/databases` | `/orgs/{org}/servers/{id}/database/schemas` |
| `/servers/{id}/database-users` | `/orgs/{org}/servers/{id}/database/users` |
| `/servers/{id}/backup-configs` | `/orgs/{org}/servers/{id}/database/backups` |
| `/servers/{id}/sites/{id}/certificates` | `/orgs/{org}/servers/{id}/sites/{id}/domains/{id}/certificate` |
