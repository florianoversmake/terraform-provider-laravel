resource "laravel_envoyer_project" "example" {
  name = "Example Envoyer Project"

  repo_provider = "github"
  repository    = "git@github.com:org_or_user/repository.git"
  branch        = "main"

  monitor = "https://example.com/monitor"

  composer_dev = false

  delete_protection = true
}
