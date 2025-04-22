# Copyright (c) HashiCorp, Inc.

# Envoyer Only Example
provider "laravel" {
  envoyer_api_token = "your-envoyer-token"
  envoyer_env_key   = "your-envoyer-env-key"
}

# Forge Only Example
provider "laravel" {
  forge_api_token = "your-forge-token"
}

# Envoyer and Forge Example
provider "laravel" {
  envoyer_api_token = "your-envoyer-token"
  envoyer_env_key   = "your-envoyer-env-key"

  forge_api_token = "your-forge-token"
}
