---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "laravel Provider"
subcategory: ""
description: |-
  
---

# laravel Provider



## Example Usage

```terraform
# Copyright (c) HashiCorp, Inc.

provider "laravel" {
  envoyer_api_token = ""
  envoyer_env_key   = ""
  forge_api_token   = ""
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `envoyer_api_token` (String, Sensitive) Envoyer API token (Bearer token).
- `forge_api_token` (String, Sensitive) Forge API token (Bearer token).

### Optional

- `cache_ttl` (Number) Time-to-live for cached API responses in seconds. Default is 300 seconds (5 minutes).
- `enable_cache` (Boolean) Enable caching of API responses. Default is false.
- `envoyer_base_url` (String) Optional override of the Envoyer API base URL (defaults to `https://envoyer.io/api`).
- `envoyer_env_key` (String, Sensitive) Optional override of the Envoyer env-lock key.
- `forge_base_url` (String) Optional override of the Forge API base URL (defaults to `https://forge.laravel.com/api/v1`).
- `max_retries` (Number) Maximum number of retries for failed API requests. Default is 3.
- `request_timeout` (Number) Timeout for API requests in seconds. Default is 30 seconds.
- `retry_delay` (Number) Delay between retries in seconds. Default is 5 seconds.
