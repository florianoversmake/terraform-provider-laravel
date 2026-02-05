# Upgrade Guide: v0.1.x to v0.2.0

This guide covers the breaking changes introduced in v0.2.0 when migrating from the deprecated Laravel Forge API v1 to the new organization-scoped JSON:API.

## Why This Migration?

Laravel Forge API v1 is deprecated and will be **discontinued on March 31, 2026**. The new API uses:

- **Organization-scoped endpoints** (`/orgs/{organization}/...`)
- **JSON:API response format** (`application/vnd.api+json`)
- **Restructured endpoints** with several path and model changes

## Provider Configuration

### Required: Add `forge_organization`

Most Forge API operations now require an organization slug. Add this to your provider configuration:

```hcl
provider "laravel" {
  forge_api_token    = "your-forge-token"
  forge_organization = "your-organization-slug"
}
```

You can find your organization slug in the Forge dashboard URL: `https://forge.laravel.com/orgs/{slug}/servers`.

> **Note:** Some endpoints like `/providers` (cloud provider information) do not require organization scope and can be accessed without setting `forge_organization`. However, most resource operations (servers, sites, workers, etc.) require it.

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

## New Data Sources

v0.2.0 introduces several new data sources that can help you query Forge resources:

### `laravel_forge_organizations`

List all organizations the authenticated user has access to:

```hcl
data "laravel_forge_organizations" "all" {}

# With filter
data "laravel_forge_organizations" "personal" {
  filter {
    name   = "slug"
    values = ["personal"]
  }
}
```

### `laravel_forge_servers`

List all servers in the organization:

```hcl
data "laravel_forge_servers" "all" {}

# With filter
data "laravel_forge_servers" "production" {
  filter {
    name   = "provider"
    values = ["digitalocean"]
  }
}
```

### `laravel_forge_sites`

List all sites on a specific server:

```hcl
data "laravel_forge_sites" "all" {
  server_id = 1234
}

# With filter
data "laravel_forge_sites" "active" {
  server_id = 1234
  filter {
    name   = "status"
    values = ["installed"]
  }
}
```

### `laravel_forge_php_versions`

List all PHP versions available on a server:

```hcl
data "laravel_forge_php_versions" "available" {
  server_id = 1234
}
```

### `laravel_forge_providers`

List all cloud providers supported by Forge:

```hcl
data "laravel_forge_providers" "all" {}

# With filter
data "laravel_forge_providers" "aws" {
  filter {
    name   = "slug"
    values = ["aws"]
  }
}
```

### `laravel_forge_regions`

List all regions for a specific cloud provider:

```hcl
data "laravel_forge_regions" "do_regions" {
  provider_id = data.laravel_forge_providers.do.providers[0].id
}

# With filter
data "laravel_forge_regions" "eu_regions" {
  provider_id = 1
  filter {
    name   = "code"
    values = ["eu-west-1", "eu-central-1"]
  }
}
```

### `laravel_forge_sizes`

List all server sizes for a specific cloud provider:

```hcl
data "laravel_forge_sizes" "do_sizes" {
  provider_id = 1
}

# With filter
data "laravel_forge_sizes" "small_instances" {
  provider_id = 1
  filter {
    name   = "category"
    values = ["general"]
  }
}
```

## Terraform Modules

v0.2.0 includes reusable Terraform modules in `tf-test/modules/`:

### Forge Module

The Forge module creates a complete server setup with optional site and worker:

```hcl
module "production_server" {
  source = "./tf-test/modules/forge"

  server_provider = "ocean2"  # DigitalOcean
  credentials_id  = data.laravel_forge_credentials.all.credentials[0].id
  server_name     = "production-server"
  region          = "nyc3"
  size            = "s-1vcpu-1gb"
  php_version     = "php83"
  database_type   = "mysql8"

  create_site       = true
  site_domain       = "example.com"
  site_project_type = "php"
  site_directory    = "/public"

  create_worker     = true
  worker_connection = "redis"
  worker_queue      = "default"
}

output "server_ip" {
  value = module.production_server.server_ip
}
```

### Envoyer Module

The Envoyer module creates a project with servers and deployment hooks:

```hcl
module "my_project" {
  source = "./tf-test/modules/envoyer"

  project_name  = "my-laravel-app"
  repo_provider = "github"
  repository    = "my-org/my-repo"
  branch        = "main"

  servers = [
    {
      name            = "web-1"
      connect_as      = "forge"
      ip_address      = "192.168.1.100"
      php_version     = "php83"
      deployment_path = "/home/forge/my-app"
    }
  ]

  hooks = [
    {
      action_id = 1  # Clone New Release
      timing    = "after"
      name      = "Install Dependencies"
      run_as    = "forge"
      script    = "cd {{release}} && composer install --no-interaction"
    }
  ]
}
```

## Import Support

Most resources support importing existing infrastructure into Terraform state.

### Import Formats

| Resource | Import Format |
|---|---|
| `laravel_forge_server` | `server_id` |
| `laravel_forge_site` | `server_id:site_id` |
| `laravel_forge_worker` | `server_id:site_id:worker_id` |
| `laravel_forge_scheduled_job` | `server_id:job_id` |
| `laravel_forge_ssh_key` | `server_id:key_id` |
| `laravel_forge_certificate` | `server_id:site_id:domain_id` |
| `laravel_forge_certificate_signing_request` | `server_id:site_id:domain_id` |
| `laravel_forge_recipe` | `recipe_id` |
| `laravel_forge_deployment_settings` | `server_id:site_id` |
| `laravel_envoyer_project` | `project_id` |
| `laravel_envoyer_server` | `project_id:server_id` |
| `laravel_envoyer_hook` | `project_id:hook_id` |
| `laravel_envoyer_environment` | `project_id` |
| `laravel_envoyer_deployment` | `project_id/deployment_id` |

### Example Import

```bash
# Import an existing Forge server
terraform import laravel_forge_server.example 12345

# Import an existing site
terraform import laravel_forge_site.example "12345:67890"

# Import an existing Envoyer deployment
terraform import laravel_envoyer_deployment.example "100/200"
```

### Resources Without Import Support

The following resources are one-time actions and cannot be imported:

- `laravel_forge_recipe_run` - Runs a recipe once, no persistent state
- `laravel_forge_certificate_signing_request_installation` - Installs a certificate once, no persistent state
