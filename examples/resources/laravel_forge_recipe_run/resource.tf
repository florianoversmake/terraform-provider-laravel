resource "laravel_forge_recipe" "example" {
  name   = "Bootstrap"
  user   = "root"
  script = <<EOVF
#!/bin/bash
export DEBIAN_FRONTEND=noninteractive
apt-get update
apt-get install -y htop libnotify-bin tmpreaper
sed -i -e 's/SHOWWARNING/#&/g' /etc/tmpreaper.conf
EOVF
}

resource "laravel_forge_recipe_run" "example" {
  servers = [
    12345,
    12346,
  ]
  recipe_id = laravel_forge_recipe.example.id
  notify    = false
}
