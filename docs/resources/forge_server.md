---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "laravel_forge_server Resource - laravel"
subcategory: ""
description: |-
  Forge server resource. This resource allows you to manage servers in Forge.
---

# laravel_forge_server (Resource)

Forge server resource. This resource allows you to manage servers in Forge.

## Example Usage

```terraform
data "laravel_forge_credentials" "aws_production" {
  filter {
    name   = "name"
    values = ["Amazon (Production)"]
  }
}

data "aws_vpc" "example" {
  id = "vpc-123abc"
}

data "aws_subnets" "example" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.example.id]
  }
}

resource "laravel_forge_server" "example" {
  depends_on      = [data.laravel_forge_credentials.aws_production]
  name            = "example-server"
  server_provider = "aws"
  type            = "app"
  region          = "eu-west-1"
  size            = "c6a.2xlarge"
  ubuntu_version  = "24.04"
  credential_id   = data.laravel_forge_credentials.aws_production.credentials[0].id

  aws_vpc_id    = data.aws_vpc.example.id
  aws_subnet_id = data.aws_subnets.example.ids[0]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String)
- `server_provider` (String)

### Optional

- `aws_subnet_id` (String) The ID of the subnet to launch the server in.
- `aws_vpc_id` (String) The ID of the VPC to launch the server in.
- `aws_vpc_name` (String) When creating a new one
- `circle` (Number)
- `credential_id` (Number)
- `database` (String) The name of the database Forge should create when building the server. If omitted, forge will be used.
- `database_type` (String) Valid values are mysql8, mariadb106, mariadb1011, mariadb114, postgres, postgres13, postgres14, postgres15, postgres16 or postgres17.
- `delete_protection` (Boolean) This is a virtual attribute and not in the API. It is used to prevent accidental deletion of the server.
- `disk_size` (Number)
- `ip_address` (String)
- `network` (List of Number) An array of server IDs that the server should be able to connect to.
- `php_version` (String)
- `private_ip_address` (String)
- `recipe_id` (Number) An optional ID of a recipe to run after provisioning.
- `region` (String)
- `revoked` (Boolean)
- `size` (String)
- `ssh_port` (Number)
- `type` (String)
- `ubuntu_version` (String)

### Read-Only

- `database_password` (String)
- `id` (Number) The ID of this resource.
- `identifier` (String)
- `is_ready` (Boolean)
- `local_public_key` (String)
- `meilisearch_password` (String)
- `provision_command` (String)
- `sudo_password` (String)
