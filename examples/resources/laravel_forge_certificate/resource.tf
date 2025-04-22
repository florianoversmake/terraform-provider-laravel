resource "tls_private_key" "ed25519-example" {
  algorithm = "ED25519"
}

resource "tls_self_signed_cert" "example" {
  private_key_pem       = tls_private_key.ed25519-example.private_key_pem
  validity_period_hours = 24

  subject {
    common_name         = "example.com"
    country             = "US"
    locality            = "San Francisco"
    organization        = "Example Inc."
    organizational_unit = "IT"
  }

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
  ]
}

resource "laravel_forge_certificate" "example" {
  server_id = 1234
  site_id   = 1234

  key         = tls_private_key.ed25519-example.private_key_pem
  certificate = tls_self_signed_cert.example.cert_pem
}
