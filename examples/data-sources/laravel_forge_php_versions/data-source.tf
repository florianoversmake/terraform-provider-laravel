# Retrieve all PHP versions on a server
data "laravel_forge_php_versions" "available" {
  server_id = 12345
}

# Output
output "php_versions" {
  value = data.laravel_forge_php_versions.available.php_versions
}

output "default_php" {
  value = [for v in data.laravel_forge_php_versions.available.php_versions : v if v.used_as_default]
}
