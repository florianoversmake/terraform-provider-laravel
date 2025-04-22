resource "laravel_forge_site" "example" {
  server_id    = 12345
  domain       = "example.com"
  wildcards    = true
  project_type = "php"
  directory    = "/current/public"
  php_version  = "php82"
}
