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
