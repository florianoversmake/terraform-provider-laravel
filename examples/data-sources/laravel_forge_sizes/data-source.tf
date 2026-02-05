# First get the provider ID
data "laravel_forge_providers" "digitalocean" {
  filter {
    name   = "slug"
    values = ["ocean2"]
  }
}

# Retrieve all sizes for DigitalOcean
data "laravel_forge_sizes" "do_sizes" {
  provider_id = data.laravel_forge_providers.digitalocean.providers[0].id
}

# Filter by category
data "laravel_forge_sizes" "general" {
  provider_id = data.laravel_forge_providers.digitalocean.providers[0].id
  filter {
    name   = "category"
    values = ["general"]
  }
}

# Filter by architecture
data "laravel_forge_sizes" "amd64" {
  provider_id = data.laravel_forge_providers.digitalocean.providers[0].id
  filter {
    name   = "architecture"
    values = ["amd64"]
  }
}

# Output
output "all_sizes" {
  value = data.laravel_forge_sizes.do_sizes.sizes
}

output "small_sizes" {
  value = [for s in data.laravel_forge_sizes.general.sizes : s if s.ram <= 4096]
}
