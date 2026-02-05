resource "laravel_forge_certificate_signing_request" "example" {
  server_id = 1234
  site_id   = 1234
  domain_id = 5678
  domain    = "example.com"

  country      = "US"
  state        = "California"
  city         = "San Francisco"
  organization = "Example Org"
  department   = "Example Department"
}

resource "cloudflare_origin_ca_certificate" "example" {
  csr                = laravel_forge_certificate_signing_request.example.certificate_signing_request
  hostnames          = ["example.com"]
  request_type       = "origin-rsa"
  requested_validity = 5475
}

resource "laravel_forge_certificate_signing_request_installation" "example" {
  domain_id         = 5678
  server_id         = 1234
  site_id           = 1234
  certificate       = cloudflare_origin_ca_certificate.example.certificate
  add_intermediates = true
}
