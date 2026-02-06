# Retrieve all cloud providers
data "laravel_forge_providers" "all" {}

# Filter by slug
data "laravel_forge_providers" "digitalocean" {
  filter {
    name   = "slug"
    values = ["ocean2"]
  }
}

# Filter by name
data "laravel_forge_providers" "aws" {
  filter {
    name   = "name"
    values = ["Amazon"]
  }
}

# Output
output "all_providers" {
  value = data.laravel_forge_providers.all.providers
}

# Use provider ID in other data sources
output "do_provider_id" {
  value = length(data.laravel_forge_providers.digitalocean.providers) > 0 ? data.laravel_forge_providers.digitalocean.providers[0].id : null
}
