# First get the provider ID
data "laravel_forge_providers" "digitalocean" {
  filter {
    name   = "slug"
    values = ["ocean2"]
  }
}

# Retrieve all regions for DigitalOcean
data "laravel_forge_regions" "do_regions" {
  provider_id = data.laravel_forge_providers.digitalocean.providers[0].id
}

# Filter by code
data "laravel_forge_regions" "nyc" {
  provider_id = data.laravel_forge_providers.digitalocean.providers[0].id
  filter {
    name   = "code"
    values = ["nyc1", "nyc3"]
  }
}

# Filter by name
data "laravel_forge_regions" "amsterdam" {
  provider_id = data.laravel_forge_providers.digitalocean.providers[0].id
  filter {
    name   = "name"
    values = ["Amsterdam"]
  }
}

# Output
output "all_regions" {
  value = data.laravel_forge_regions.do_regions.regions
}
