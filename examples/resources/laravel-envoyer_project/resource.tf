# Copyright (c) HashiCorp, Inc.

resource "laravel-envoyer_project" "example" {
  name          = "example"
  repo_provider = "github"
  repository    = "smakecloud/laravel-envoyer-example"
}
