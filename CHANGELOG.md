## 0.2.0 (Unreleased)

BREAKING CHANGES:

* provider: New required `forge_organization` attribute. All Forge API operations now require an organization slug.
* provider: `forge_base_url` default changed from `https://forge.laravel.com/api/v1` to `https://forge.laravel.com/api`.
* resource/laravel_forge_certificate: New required `domain_id` attribute. Certificates are now managed per domain record instead of per site. Removed `active`, `domain`, `existing` attributes. `created_at` changed from Number to String.
* resource/laravel_forge_certificate_signing_request: New required `domain_id` attribute. Removed `active`, `existing` attributes. `created_at` changed from Number to String.
* resource/laravel_forge_certificate_signing_request_installation: Replaced `certificate_signing_request_id` with `domain_id` attribute.
* resource/laravel_forge_site: Removed `database` and `nginx_template` attributes.
* resource/laravel_forge_worker: Workers are now managed as background processes. The resource now builds a `php artisan queue:work` command from the individual parameters. The API no longer returns queue-worker-specific fields (connection, timeout, sleep, etc.) separately.
* resource/laravel_forge_scheduled_job: Individual time fields (minute, hour, day, month, weekday) are now combined into a cron expression when creating jobs. `created_at` changed from String to nullable.
* resource/laravel_forge_ssh_key: `username` field renamed to `user` in the API.
* data-source/laravel_forge_credentials: `type` field renamed to `provider` in the API.

ENHANCEMENTS:

* Migrated entire Forge API client from deprecated v1 API to new organization-scoped JSON:API format.
* All API endpoints now use organization-scoped paths (`/orgs/{org}/...`).
* Added JSON:API response parsing with generic helper functions.
* Added unit tests for JSON:API helpers, worker command builder, cron parser, and composite ID parsing.

NOTES:

* The Laravel Forge API v1 is deprecated and will be discontinued on March 31, 2026.
* All endpoint paths have been updated to match the new API specification.
* Endpoint renames: `daemons` -> `background-processes`, `jobs` -> `scheduled-jobs`, `keys` -> `ssh-keys`, `credentials` -> `server-credentials`, `databases` -> `database/schemas`, `database-users` -> `database/users`, `backup-configs` -> `database/backups`.

## 0.1.0 (Unreleased)

FEATURES:
