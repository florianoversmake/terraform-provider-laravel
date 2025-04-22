resource "laravel_envoyer_project" "example" {
  name = "Example Envoyer Project"

  repo_provider = "github"
  repository    = "git@github.com:org_or_user/repository.git"
  branch        = "main"

  monitor = "https://example.com/monitor"

  composer_dev = false

  delete_protection = true
}

resource "laravel_envoyer_server" "example" {
  project_id = laravel_envoyer_project.example.id

  name        = "example-server"
  connect_as  = "forge"
  ip_address  = "127.0.0.1"
  php_version = "php82"

  deployment_path = "/home/forge/example.com"
}
