resource "laravel_forge_certificate_signing_request" "example" {
  server_id = 1234
  site_id   = 1234
  domain    = "example.com"

  country      = "US"
  state        = "California"
  city         = "San Francisco"
  organization = "Example Org"
  department   = "Example Department"
}
