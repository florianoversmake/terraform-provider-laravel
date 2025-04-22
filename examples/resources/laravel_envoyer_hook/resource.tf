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

data "laravel_envoyer_actions" "clone" {
  filter {
    name   = "view"
    values = ["scripts.deployments.CloneNewRelease"]
  }
}

resource "laravel_envoyer_hook" "prepare-shared-folder" {
  project_id = laravel_envoyer_project.example.id
  servers    = [laravel_envoyer_server.example.id]
  name       = "Prepare Shared Folder"
  script     = <<EOF
#!/usr/bin/env bash

mkdir -p {{project}}/storage/tmp
mkdir -p {{project}}/storage/app/tmp
mkdir -p {{project}}/storage/exports
mkdir -p {{project}}/storage/framework/cache
mkdir -p {{project}}/storage/framework/cookies
mkdir -p {{project}}/storage/framework/sessions
mkdir -p {{project}}/storage/framework/views
mkdir -p {{project}}/storage/logs
mkdir -p {{project}}/storage/medialibrary/temp
EOF
  run_as     = "forge"
  action_id  = data.laravel_envoyer_actions.clone.actions[0].id
  timing     = "before"
}
