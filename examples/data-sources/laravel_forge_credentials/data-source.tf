data "laravel_forge_credentials" "aws_production" {
  filter {
    name   = "name"
    values = ["Amazon (Production)"]
  }
}
