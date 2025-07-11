---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "laravel_forge_scheduled_job Resource - laravel"
subcategory: ""
description: |-
  Forge scheduled job resource. This resource allows you to manage scheduled jobs on Forge servers.
---

# laravel_forge_scheduled_job (Resource)

Forge scheduled job resource. This resource allows you to manage scheduled jobs on Forge servers.

## Example Usage

```terraform
resource "laravel_forge_scheduled_job" "example" {
  server_id = 12345
  command   = "echo 'Hello World!'"
  frequency = "nightly"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `command` (String) The command to run on the server.
- `frequency` (String) The frequency in which the job should run. Valid values are `minutely`, `hourly`, `nightly`, `weekly`, `monthly`, `reboot`, and `custom`
- `server_id` (Number) The ID of the server where the scheduled job will be run.

### Optional

- `day` (String) The day at which the job should run. Required if frequency is `custom`.
- `hour` (String) The hour at which the job should run. Required if frequency is `custom`.
- `minute` (String) The minute at which the job should run. Required if frequency is `custom`.
- `month` (String) The month at which the job should run. Required if frequency is `custom`.
- `user` (String) The user under which the command will be run.
- `weekday` (String) The weekday at which the job should run. Required if frequency is `custom`.

### Read-Only

- `created_at` (String)
- `id` (Number) The ID of this resource.
- `status` (String)
