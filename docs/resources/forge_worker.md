---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "laravel_forge_worker Resource - laravel"
subcategory: ""
description: |-
  
---

# laravel_forge_worker (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `server_id` (Number)
- `site_id` (Number)
- `worker_connection` (String)

### Optional

- `daemon` (Boolean)
- `delay` (Number)
- `directory` (String)
- `force` (Boolean)
- `memory` (Number)
- `php_version` (String)
- `processes` (Number)
- `queue` (String)
- `sleep` (Number)
- `stop_wait_secs` (Number)
- `timeout` (Number)
- `tries` (Number)

### Read-Only

- `command` (String)
- `created_at` (String)
- `id` (Number) The ID of this resource.
- `status` (String)
