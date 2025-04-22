# Laravel Terraform Provider

A Terraform provider for managing Laravel Forge and Laravel Envoyer resources. This provider allows you to automate infrastructure provisioning and application deployment across these Laravel first-party services.

## Features

This provider allows you to manage resources across both Laravel Forge and Laravel Envoyer.

### Laravel Forge Resources

- Servers
- Sites
- Workers
- Recipes
- SSH Keys
- SSL Certificates & Certificate Signing Requests
- Scheduled Jobs

### Laravel Envoyer Resources

- Projects
- Deployments
- Servers
- Hooks
- Environment Variables

## Installation

Add the following to your Terraform configuration:

```hcl
terraform {
  required_providers {
    laravel = {
      source = "registry.terraform.io/florianoversmake/laravel"
      version = "~> 1.0"
    }
  }
}
```

Then run `terraform init` to download the provider.

### Provider Configuration

```hcl
provider "laravel" {
  # Forge configuration (optional if only using Envoyer)
  forge_api_token = "your-forge-api-token"
  forge_base_url  = "https://forge.laravel.com/api/v1" # Optional

  # Envoyer configuration (optional if only using Forge)
  envoyer_api_token = "your-envoyer-api-token"
  envoyer_env_key   = "your-env-key" # Optional
  envoyer_base_url  = "https://envoyer.io/api" # Optional

  # Advanced configuration (all optional)
  request_timeout = 30  # seconds
  max_retries     = 3
  retry_delay     = 5   # seconds
  enable_cache    = true
  cache_ttl       = 300 # seconds
}
```

## License

MIT
