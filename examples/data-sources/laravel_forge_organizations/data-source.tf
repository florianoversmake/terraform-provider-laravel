# Retrieve all organizations
data "laravel_forge_organizations" "all" {}

# Filter by slug
data "laravel_forge_organizations" "personal" {
  filter {
    name   = "slug"
    values = ["personal"]
  }
}

# Filter by name (partial match)
data "laravel_forge_organizations" "team" {
  filter {
    name   = "name"
    values = ["Team"]
  }
}

# Output
output "all_organizations" {
  value = data.laravel_forge_organizations.all.organizations
}
