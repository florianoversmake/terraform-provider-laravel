data "laravel_envoyer_actions" "clone" {
  filter {
    name   = "view"
    values = ["scripts.deployments.CloneNewRelease"]
  }
}
