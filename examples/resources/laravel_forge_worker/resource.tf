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
