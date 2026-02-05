# Retrieve all sites on a server
data "laravel_forge_sites" "all" {
  server_id = 12345
}

# Filter by status
data "laravel_forge_sites" "installed" {
  server_id = 12345
  filter {
    name   = "status"
    values = ["installed"]
  }
}

# Filter by PHP version
data "laravel_forge_sites" "php82" {
  server_id = 12345
  filter {
    name   = "php_version"
    values = ["php82", "php83"]
  }
}

# Filter by application type
data "laravel_forge_sites" "laravel" {
  server_id = 12345
  filter {
    name   = "app_type"
    values = ["php"]
  }
}

# Output
output "all_sites" {
  value = data.laravel_forge_sites.all.sites
}
