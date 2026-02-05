# Retrieve all servers in the organization
data "laravel_forge_servers" "all" {}

# Filter by cloud provider
data "laravel_forge_servers" "digitalocean" {
  filter {
    name   = "provider"
    values = ["ocean2"]
  }
}

# Filter by region
data "laravel_forge_servers" "eu_servers" {
  filter {
    name   = "region"
    values = ["nyc1", "nyc3", "fra1"]
  }
}

# Filter by name
data "laravel_forge_servers" "production" {
  filter {
    name   = "name"
    values = ["production", "prod"]
  }
}

# Output
output "all_servers" {
  value = data.laravel_forge_servers.all.servers
}
