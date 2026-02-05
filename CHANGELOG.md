## 0.2.0 (Unreleased)

BREAKING CHANGES:

* provider: New required `forge_organization` attribute. Most Forge API operations now require an organization slug.
* provider: `forge_base_url` default changed from `https://forge.laravel.com/api/v1` to `https://forge.laravel.com/api`.
* resource/laravel_forge_certificate: New required `domain_id` attribute. Certificates are now managed per domain record instead of per site. Removed `active`, `domain`, `existing` attributes. `created_at` changed from Number to String.
* resource/laravel_forge_certificate_signing_request: New required `domain_id` attribute. Removed `active`, `existing` attributes. `created_at` changed from Number to String.
* resource/laravel_forge_certificate_signing_request_installation: Replaced `certificate_signing_request_id` with `domain_id` attribute.
* resource/laravel_forge_site: Removed `database` and `nginx_template` attributes.
* resource/laravel_forge_worker: Workers are now managed as background processes. The resource now builds a `php artisan queue:work` command from the individual parameters. The API no longer returns queue-worker-specific fields (connection, timeout, sleep, etc.) separately.
* resource/laravel_forge_scheduled_job: Individual time fields (minute, hour, day, month, weekday) are now combined into a cron expression when creating jobs. `created_at` changed from String to nullable.
* resource/laravel_forge_ssh_key: `username` field renamed to `user` in the API.
* data-source/laravel_forge_credentials: `type` field renamed to `provider` in the API.

NEW DATA SOURCES:

* data-source/laravel_forge_organizations: List all organizations the authenticated user has access to, with optional filtering by name or slug.
* data-source/laravel_forge_servers: List all servers in the organization, with optional filtering by name, provider, region, or IP address.
* data-source/laravel_forge_sites: List all sites on a specific server, with optional filtering by name, status, php_version, or app_type.
* data-source/laravel_forge_php_versions: List all PHP versions available on a specific server.
* data-source/laravel_forge_providers: List all cloud providers (DigitalOcean, AWS, Hetzner, etc.), with optional filtering by name or slug.
* data-source/laravel_forge_regions: List all regions for a specific cloud provider, with optional filtering by name or code.
* data-source/laravel_forge_sizes: List all server sizes for a specific cloud provider, with optional filtering by name, code, category, series, disk_type, or architecture.

NEW TERRAFORM MODULES:

* tf-test/modules/forge: Reusable Terraform module for creating Forge servers with optional sites and workers.
* tf-test/modules/envoyer: Reusable Terraform module for creating Envoyer projects with servers and deployment hooks.

ENHANCEMENTS:

* All data sources now automatically paginate through all results using cursor-based pagination. Previously, list endpoints only returned the first page (default 30 items).
* resource/laravel_envoyer_deployment: Added import support (format: project_id/deployment_id).
* Migrated entire Forge API client from deprecated v1 API to new organization-scoped JSON:API format.
* Most API endpoints now use organization-scoped paths (`/orgs/{org}/...`).
* Provider endpoints (`/providers`, `/providers/{id}/regions`, etc.) do not require organization scope.
* Added JSON:API response parsing with generic helper functions.
* Added unit tests for JSON:API helpers, worker command builder, cron parser, and composite ID parsing.
* Split provider-related API calls into separate `forge_providers.go` file for clarity.
* Refactored region/size lookup methods to use efficient targeted API calls (2-4 calls instead of 100+).
* Added organization listing API methods in `organizations.go`.
* Added comprehensive non-invasive integration tests for all API endpoints.
* Added example Terraform configurations in tf-test/forge and tf-test/envoyer.

NOTES:

* The Laravel Forge API v1 is deprecated and will be discontinued on March 31, 2026.
* All endpoint paths have been updated to match the new API specification.
* Endpoint renames: `daemons` -> `background-processes`, `jobs` -> `scheduled-jobs`, `keys` -> `ssh-keys`, `credentials` -> `server-credentials`, `databases` -> `database/schemas`, `database-users` -> `database/users`, `backup-configs` -> `database/backups`.

## 0.1.0 (Unreleased)

FEATURES:
