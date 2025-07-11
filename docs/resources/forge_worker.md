---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "laravel_forge_worker Resource - laravel"
subcategory: ""
description: |-
  Forge worker resource. This resource allows you to manage workers in Forge.
---

# laravel_forge_worker (Resource)

Forge worker resource. This resource allows you to manage workers in Forge.

## Example Usage

```terraform
resource "laravel_forge_worker" "example" {
  server_id = 12345
  site_id   = 12345
  directory = "/home/forge/example.com/current"

  worker_connection = "sqs"
  sleep             = 10
  delay             = 0
  memory            = 256
  php_version       = "php82"
  timeout           = 60
  processes         = 2
  tries             = 3
  queue             = "default"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `server_id` (Number) The ID of the server where the worker is created.
- `site_id` (Number) The ID of the site where the worker is created.
- `worker_connection` (String) The connection string for the worker. Like `sync`, `database`, `beanstalkd`, `sqs`, `redis`...

### Optional

- `daemon` (Boolean) Whether the worker should run as a daemon. Default is true.
- `delay` (Number) The delay time for the worker in seconds. Default is 0 seconds.
- `directory` (String) The directory where the worker is located. Default is empty string (current directory).
- `force` (Boolean) To force your queue workers to process jobs even if maintenance mode is enabled, you may use force option.
- `memory` (Number) The memory limit for the worker in megabytes. Default is 128 MB.
- `php_version` (String) The PHP version to use for the worker. Default is 'php' (System default).
- `processes` (Number) The number of processes for the worker. Default is 1.
- `queue` (String) The queue name for the worker. Default is empty string (no specific queue).
- `sleep` (Number) The sleep time for the worker in seconds. Default is 3 seconds.
- `stop_wait_secs` (Number) The number of seconds to wait for the worker to stop. Default is 10 seconds. You should ensure that the value of stopwaitsecs is greater than the number of seconds consumed by your longest running job. Otherwise, Supervisor may kill the job before it is finished processing.
- `timeout` (Number) The timeout for the worker in seconds. Default is 60 seconds.
- `tries` (Number) The number of tries for the worker. Default is 0 (unlimited).

### Read-Only

- `command` (String)
- `created_at` (String)
- `id` (Number) The ID of this resource.
- `status` (String)
