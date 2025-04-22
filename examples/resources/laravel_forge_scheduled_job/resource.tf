resource "laravel_forge_scheduled_job" "example" {
  server_id = 12345
  command   = "echo 'Hello World!'"
  frequency = "nightly"
}
