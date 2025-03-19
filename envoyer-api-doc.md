# Introduction

The Envoyer API allows you to create and manage projects, servers, environments, hooks, deployments, collaborators and notifications on Envoyer through a simple REST API.

# Authentication

In order to use the API, you should authenticate your request by including your API key as a bearer token value:

`Authorization: Bearer API_KEY_HERE`

# Headers

Make sure you have the following content type headers are on every request:

```
Accept: application/json
Content-Type: application/json
```

# URI

Envoyers's API is hosted on the following base URI:

`https://envoyer.io/api`

# Errors

Envoyer uses conventional HTTP response codes to indicate the success or failure of an API request. The table below contains a summary of the typical response codes:

| Code | Description |
| --- | --- |
| 200 | Everything is ok. |
| 400 | Valid data was given but the request has failed. |
| 401 | No valid API Key was given. |
| 404 | The request resource could not be found. |
| 422 | The payload has missing required parameters or invalid data was given. |
| 429 | Too many attempts. |
| 500 | Request failed due to an internal error in Envoyer. |
| 503 | Envoyer is offline for maintenance. |

# Scopes

Scopes can limit access to your data over the API. There are a couple of points that you should be aware of:

*   All `GET` endpoints are freely accessible with any API token.
*   `*:create` tokens work for creating and updating of resources.

# Projects

## List Projects

> Response

```json
{
  "projects": [
    {
      "id": 1,
      "user_id": 1,
      "version": 1,
      "name": "Laravel",
      "provider": "github",
      "plain_repository": "laravel\/laravel",
      "repository": "git@github.com:laravel\/laravel.git",
      "type": "laravel-5",
      "branch": "main",
      "push_to_deploy": false,
      "webhook_id": null,
      "status": null,
      "should_deploy_again": 0,
      "deployment_started_at": null,
      "deployment_finished_at": "2020-10-09T10:54:38.000000Z",
      "last_deployment_status": "finished",
      "daily_deploys": 89,
      "weekly_deploys": 89,
      "last_deployment_took": 37,
      "retain_deployments": 4,
      "environment_servers": [
        4
      ],
      "folders": [],
      "monitor": "http:\/\/laravel.com",
      "new_york_status": "healthy",
      "london_status": "healthy",
      "singapore_status": "healthy",
      "token": "gJw0HbKUG0Zegum6oSNuYf0OvvHWy2z47eElFMiM",
      "created_at": "2020-09-10T13:56:28.000000Z",
      "updated_at": "2020-10-09T11:28:37.000000Z",
      "install_dev_dependencies": false,
      "install_dependencies": true,
      "quiet_composer": false,
      "servers": [
          {

          }
      ],
      "has_environment": true,
      "has_monitoring_error": false,
      "has_missing_heartbeats": true,
      "last_deployed_branch": "main",
      "last_deployment_id": 57,
      "last_deployment_author": "Dries Vints",
      "last_deployment_avatar": "https:\/\/avatars1.githubusercontent.com\/u\/594614?v=4",
      "last_deployment_hash": "c66546e",
      "last_deployment_timestamp": "4 days ago"
    }
  ]
}
```

### HTTP Request

`GET /api/projects`

## Get Project

> Response

```json
{
  "project": {
      "id": 1,
      "user_id": 1,
      "version": 1,
      "name": "Laravel",
      "provider": "github",
      "plain_repository": "laravel\/laravel",
      "repository": "git@github.com:laravel\/laravel.git",
      "type": "laravel-5",
      "branch": "main",
      "push_to_deploy": false,
      "webhook_id": null,
      "status": null,
      "should_deploy_again": 0,
      "deployment_started_at": null,
      "deployment_finished_at": "2020-10-09T10:54:38.000000Z",
      "last_deployment_status": "finished",
      "daily_deploys": 89,
      "weekly_deploys": 89,
      "last_deployment_took": 37,
      "retain_deployments": 4,
      "environment_servers": [
        4
      ],
      "folders": [],
      "monitor": "http:\/\/laravel.com",
      "new_york_status": "healthy",
      "london_status": "healthy",
      "singapore_status": "healthy",
      "token": "gJw0HbKUG0Zegum6oSNuYf0OvvHWy2z47eElFMiM",
      "created_at": "2020-09-10T13:56:28.000000Z",
      "updated_at": "2020-10-09T11:28:37.000000Z",
      "install_dev_dependencies": false,
      "install_dependencies": true,
      "quiet_composer": false,
      "servers": [
          {

          }
      ],
      "has_environment": true,
      "has_monitoring_error": false,
      "has_missing_heartbeats": true,
      "last_deployed_branch": "main",
      "last_deployment_id": 57,
      "last_deployment_author": "Dries Vints",
      "last_deployment_avatar": "https:\/\/avatars1.githubusercontent.com\/u\/594614?v=4",
      "last_deployment_hash": "c66546e",
      "last_deployment_timestamp": "4 days ago"
    }
}
```

### HTTP Request

`GET /api/projects/{projectId}`

## Create Project

> Payload

```json
{
  "name": "Laravel",
  "provider": "github",
  "repository": "laravel/laravel",
  "branch": "main",
  "type": "laravel-5",
  "retain_deployments": 5,
  "monitor": "https://my-laravel-project.com",
  "composer": true,
  "composer_dev": false,
  "composer_quiet": false
}
```

> Response

```json
{
  "project": {
    "name": "Laravel",
    "provider": "github",
    "plain_repository": "laravel/laravel",
    "repository": "git@github.com:laravel/laravel.git",
    "branch": "main",
    "type": "laravel-5",
    "version": 1,
    "token": "cegnVeJJhfdT9WMmS8ye9w0gpRegpobQ707oe6U2",
    "folders": [],
    "monitor": "https://my-laravel-project.com",
    "retain_deployments": 5,
    "install_dependencies": true,
    "install_dev_dependencies": false,
    "quiet_composer": false,
    "user_id": 1,
    "updated_at": "2020-10-14T15:12:28.000000Z",
    "created_at": "2020-10-14T15:12:28.000000Z",
    "id": 12,
    "has_environment": false,
    "has_monitoring_error": false,
    "has_missing_heartbeats": false,
    "environment_servers": [],
    "last_deployed_branch": null,
    "last_deployment_id": null,
    "last_deployment_author": null,
    "last_deployment_avatar": null,
    "last_deployment_hash": "",
    "last_deployment_timestamp": null,
    "servers": []
  }
}
```

### HTTP Request

`POST /api/projects`

### Acceptable Providers

| Key | Provider |
| --- | --- |
| `bitbucket` | Bitbucket |
| `github` | GitHub |
| `gitlab` | GitLab.com |
| `gitlab-self` | Self-hosted GitLab |

### Acceptable Project Types

| Key | Project Type |
| --- | --- |
| `laravel-5` | Laravel |
| `laravel-4` | Laravel 4 |
| `other` | Static HTML, Other PHP Projects |

### Required Scopes

This endpoint requires the `projects:create` scope.

## Update Project

> Payload

```json
{
  "name": "Laravel",
  "retain_deployments": 10,
  "monitor": "https://my-new-laravel-project.com",
  "composer": true,
  "composer_dev": false,
  "composer_quiet": false
}
```

> Response

### HTTP Request

`PUT /api/projects/{projectId}`

### Required Scopes

This endpoint requires the `projects:create` scope.

## Update Project Source

> Payload

```json
{
  "provider": "github",
  "repository": "laravel/laravel",
  "branch": "main",
  "push_to_deploy": false
}
```

### HTTP Request

`PUT /api/projects/{projectId}/source`

### Required Scopes

This endpoint requires the `projects:create` scope.

## Delete Project

### HTTP Request

`DELETE /api/projects/{projectId}`

### Required Scopes

This endpoint requires the `projects:delete` scope.

## Get Linked Folders

> Response

```json
{
  "folders": [
    {
      "from": "public\/uploads",
      "to": "storage\/app\/uploads"
    }
  ]
}
```

### HTTP Request

`GET /api/projects/{projectId}/folders`

## Create Linked Folder

> Request

```json
{
  "from": "public/uploads",
  "to": "storage/app/uploads"
}
```

> Response

```json
{
  "folders": [
    {
      "from": "public\/uploads",
      "to": "storage\/app\/uploads"
    }
  ]
}
```

### HTTP Request

`POST /api/projects/{projectId}/folders`

## Delete Linked Folder

> Request

```json
{
  "from": "public/uploads",
  "to": "storage/app/uploads"
}
```

### HTTP Request

`DELETE /api/projects/{projectId}/folders`

# Servers

## List Servers

> Response

```json
{
  "servers": [
    {
      "id": 7,
      "project_id": 1,
      "user_id": 1,
      "name": "My Test Server",
      "connect_as": "forge",
      "ip_address": "test.laravel.com",
      "port": "22",
      "php_seven": false,
      "php_version": "php74",
      "freebsd": true,
      "receives_code_deployments": true,
      "should_restart_fpm": true,
      "deployment_path": "/home/forge/test.laravel.com",
      "php_path": "php",
      "composer_path": "composer",
      "public_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQD3vYKSuh7rJf+NtWn04CFyT9+nmx+i+/sP+yMN9ueJ+Rd5Ku6d9kgscK2xwlRlkcA0sethslu0WUsG81RC1lVpF6iLrc/9O45ZhEY1CB/7dofr+7ZNwu/DJtbW6YE7oyT5G97BUW763TMq/YO9/xjMToetElTEJ4hUVWdP8q93b3MVHBazk2PEuS05wzP4p5XeQnhKq4LISetJFEgI8Y+HEpK29GiU/18fhaGZvdVwOToOxTwEwBbS3fTLNkBaUTWw9q3i7S60RRncBCHppcs2irrzw7yt7ZQOnut/BIjIGESoxx+N4ZrpTmX6P5d3/9Duk40Mfwh1ftsvze6o5AW4Xi0tki8b6bsMXmO7SapqVdiMZ5/4BWOkqHWhi926qz7I9NWoZuVFAUpSoe6fObzQBRooVp7ARw7gJ4C+Q4xc1gJJkZoQ/Wj/wHkVnbLw9M5+t5GjyWgDDOr5iyoGOyIwhuEFvATzIYH0z5B6anL1n6XQmeGh5OWKJN8wE5qVNTU= worker@envoyer.io\n",
      "connection_status": "unknown",
      "current_activity": null,
      "created_at": "2020-10-15T11:05:26.000000Z",
      "updated_at": "2020-10-15T11:05:26.000000Z"
    }
  ]
}
```

### HTTP Request

`GET /api/projects/{projectId}/servers`

## Get Server

> Response

```json
{
  "server": {
    "id": 7,
    "project_id": 1,
    "user_id": 1,
    "name": "My Test Server",
    "connect_as": "forge",
    "ip_address": "test.laravel.com",
    "port": "22",
    "php_seven": false,
    "php_version": "php74",
    "freebsd": true,
    "receives_code_deployments": true,
    "should_restart_fpm": true,
    "deployment_path": "/home/forge/test.laravel.com",
    "php_path": "php",
    "composer_path": "composer",
    "public_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQD3vYKSuh7rJf+NtWn04CFyT9+nmx+i+/sP+yMN9ueJ+Rd5Ku6d9kgscK2xwlRlkcA0sethslu0WUsG81RC1lVpF6iLrc/9O45ZhEY1CB/7dofr+7ZNwu/DJtbW6YE7oyT5G97BUW763TMq/YO9/xjMToetElTEJ4hUVWdP8q93b3MVHBazk2PEuS05wzP4p5XeQnhKq4LISetJFEgI8Y+HEpK29GiU/18fhaGZvdVwOToOxTwEwBbS3fTLNkBaUTWw9q3i7S60RRncBCHppcs2irrzw7yt7ZQOnut/BIjIGESoxx+N4ZrpTmX6P5d3/9Duk40Mfwh1ftsvze6o5AW4Xi0tki8b6bsMXmO7SapqVdiMZ5/4BWOkqHWhi926qz7I9NWoZuVFAUpSoe6fObzQBRooVp7ARw7gJ4C+Q4xc1gJJkZoQ/Wj/wHkVnbLw9M5+t5GjyWgDDOr5iyoGOyIwhuEFvATzIYH0z5B6anL1n6XQmeGh5OWKJN8wE5qVNTU= worker@envoyer.io\n",
    "connection_status": "unknown",
    "current_activity": null,
    "created_at": "2020-10-15T11:05:26.000000Z",
    "updated_at": "2020-10-15T11:05:26.000000Z"
  }
}
```

### HTTP Request

`GET /api/projects/{projectId}/servers/{serverId}`

## Refresh Server Connection

Refreshes the status of Envoyer's connection to the server.

### HTTP Request

`POST /api/projects/{projectId}/servers/{serverId}/refresh`

## Create Server

> Payload

```json
{
  "name": "My Test Server",
  "connectAs": "forge",
  "host": "test.laravel.com",
  "port": 22,
  "phpVersion": "php74",
  "receivesCodeDeployments": true,
  "deploymentPath": "/home/forge/test.laravel.com",
  "restartFpm": true,
  "composerPath": "composer"
}
```

When creating a server, you must supply at minimum:

*   `name`
*   `connectAs`
*   `host`
*   `phpVersion`

The `port` parameter will default to `22` as this is the default port that SSH listens on.

> Response

```json
{
  "server": {
    "id": 7,
    "project_id": 1,
    "user_id": 1,
    "name": "My Test Server",
    "connect_as": "forge",
    "ip_address": "test.laravel.com",
    "port": "22",
    "php_seven": false,
    "php_version": "php74",
    "freebsd": true,
    "receives_code_deployments": true,
    "should_restart_fpm": true,
    "deployment_path": "/home/forge/test.laravel.com",
    "php_path": "php",
    "composer_path": "composer",
    "public_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQD3vYKSuh7rJf+NtWn04CFyT9+nmx+i+/sP+yMN9ueJ+Rd5Ku6d9kgscK2xwlRlkcA0sethslu0WUsG81RC1lVpF6iLrc/9O45ZhEY1CB/7dofr+7ZNwu/DJtbW6YE7oyT5G97BUW763TMq/YO9/xjMToetElTEJ4hUVWdP8q93b3MVHBazk2PEuS05wzP4p5XeQnhKq4LISetJFEgI8Y+HEpK29GiU/18fhaGZvdVwOToOxTwEwBbS3fTLNkBaUTWw9q3i7S60RRncBCHppcs2irrzw7yt7ZQOnut/BIjIGESoxx+N4ZrpTmX6P5d3/9Duk40Mfwh1ftsvze6o5AW4Xi0tki8b6bsMXmO7SapqVdiMZ5/4BWOkqHWhi926qz7I9NWoZuVFAUpSoe6fObzQBRooVp7ARw7gJ4C+Q4xc1gJJkZoQ/Wj/wHkVnbLw9M5+t5GjyWgDDOr5iyoGOyIwhuEFvATzIYH0z5B6anL1n6XQmeGh5OWKJN8wE5qVNTU= worker@envoyer.io\n",
    "connection_status": "unknown",
    "current_activity": null,
    "created_at": "2020-10-15T11:05:26.000000Z",
    "updated_at": "2020-10-15T11:05:26.000000Z"
  }
}
```

### HTTP Request

`POST /api/projects/{projectId}/servers`

### Required Scopes

This endpoint requires the `servers:create` scope.

### PHP Versions

| Version | Slug |
| --- | --- |
| PHP 8.1 | `php81` |
| PHP 8.0 | `php80` |
| PHP 7.4 | `php74` |
| PHP 7.3 | `php73` |
| PHP 7.2 | `php72` |
| PHP 7.1 | `php71` |
| PHP 7.0 | `php70` |
| PHP 5.6 | `php56` |

### IP Address vs Hostname

The resulting `server` object will contain the hostname or IP address within the `ip_address` key.

## Update Server

> Payload

```json
{
  "name": "My New Server"
}
```

The payload used for updating a server is the same as creating.

> Response

```json
{
  "server": {
    "id": 7,
    "project_id": 1,
    "user_id": 1,
    "name": "My New Server",
    "connect_as": "forge",
    "ip_address": "test.laravel.com",
    "port": "22",
    "php_seven": false,
    "php_version": "php74",
    "freebsd": true,
    "receives_code_deployments": true,
    "should_restart_fpm": true,
    "deployment_path": "/home/forge/test.laravel.com",
    "php_path": "php",
    "composer_path": "composer",
    "public_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQD3vYKSuh7rJf+NtWn04CFyT9+nmx+i+/sP+yMN9ueJ+Rd5Ku6d9kgscK2xwlRlkcA0sethslu0WUsG81RC1lVpF6iLrc/9O45ZhEY1CB/7dofr+7ZNwu/DJtbW6YE7oyT5G97BUW763TMq/YO9/xjMToetElTEJ4hUVWdP8q93b3MVHBazk2PEuS05wzP4p5XeQnhKq4LISetJFEgI8Y+HEpK29GiU/18fhaGZvdVwOToOxTwEwBbS3fTLNkBaUTWw9q3i7S60RRncBCHppcs2irrzw7yt7ZQOnut/BIjIGESoxx+N4ZrpTmX6P5d3/9Duk40Mfwh1ftsvze6o5AW4Xi0tki8b6bsMXmO7SapqVdiMZ5/4BWOkqHWhi926qz7I9NWoZuVFAUpSoe6fObzQBRooVp7ARw7gJ4C+Q4xc1gJJkZoQ/Wj/wHkVnbLw9M5+t5GjyWgDDOr5iyoGOyIwhuEFvATzIYH0z5B6anL1n6XQmeGh5OWKJN8wE5qVNTU= worker@envoyer.io\n",
    "connection_status": "unknown",
    "current_activity": null,
    "created_at": "2020-10-15T11:05:26.000000Z",
    "updated_at": "2020-10-15T11:05:26.000000Z"
  }
}
```

### HTTP Request

`PUT /api/projects/{projectId}/servers/{serverId}`

### Required Scopes

This endpoint requires the `servers:create` scope.

## Delete Server

### HTTP Request

`DELETE /api/projects/{projectId}/servers/{serverId}`

### Required Scopes

This endpoint requires the `servers:delete` scope.

# Environments

## Get Environment

> Request

```json
{
  "key": "foobar"
}
```

> Response

```json
{
  "environment": "APP_NAME=Laravel\nAPP_DEBUG=true\n"
}
```

### HTTP Request

`GET /api/projects/{projectId}/environment`

## Get Environment Servers

> Response

```json
{
  "servers": [
    {
      "id": 4,
      "project_id": 1,
      "user_id": 1,
      "name": "floral-moon",
      "connect_as": "forge",
      "ip_address": "x.x.x.x",
      "port": "22",
      "php_seven": false,
      "php_version": "php74",
      "freebsd": false,
      "receives_code_deployments": true,
      "should_restart_fpm": true,
      "deployment_path": "\/home\/forge\/default",
      "php_path": "php",
      "composer_path": "composer",
      "public_key": "ssh-rsa ... worker@envoyer.io\n",
      "connection_status": "successful",
      "current_activity": null,
      "created_at": "2020-10-08T11:50:20.000000Z",
      "updated_at": "2020-10-09T11:28:41.000000Z"
    }
  ]
}
```

### HTTP Request

`GET /api/projects/{projectId}/environment/servers`

## Update Environment

> Request

```json
{
  "key": "foobar",
  "contents": "APP_NAME=Laravel\nAPP_DEBUG=false\n",
  "servers": [
    4
  ]
}
```

### HTTP Request

`PUT /api/projects/{projectId}/environment`

### Required Scopes

This endpoint requires the `environments:create` scope.

## Reset Environment Key

**Resetting your environment key will wipe the existing environment data.**

> Request

```json
{
  "key": "new-key"
}
```

### HTTP Request

`DELETE /api/projects/{projectId}/environment`

### Required Scopes

This endpoint requires the `environments:delete` scope.

# Actions

## List Actions

> Response

```json
{
  "actions": [
    {
      "id": 1,
      "version": 1,
      "name": "Clone New Release",
      "view": "scripts.deployments.CloneNewRelease",
      "sequence": 1,
      "created_at": "2020-09-11T12:01:56.000000Z",
      "updated_at": "2020-09-11T12:01:56.000000Z"
    },
    {
      "id": 2,
      "version": 1,
      "name": "Install Composer Dependencies",
      "view": "scripts.deployments.InstallComposerDependencies",
      "sequence": 2,
      "created_at": "2020-09-11T12:01:56.000000Z",
      "updated_at": "2020-09-11T12:01:56.000000Z"
    },
    {
      "id": 3,
      "version": 1,
      "name": "Activate New Release",
      "view": "scripts.deployments.ActivateNewRelease",
      "sequence": 3,
      "created_at": "2020-09-11T12:01:56.000000Z",
      "updated_at": "2020-09-11T12:01:56.000000Z"
    },
    {
      "id": 4,
      "version": 1,
      "name": "Purge Old Releases",
      "view": "scripts.deployments.PurgeOldReleases",
      "sequence": 4,
      "created_at": "2020-09-11T12:01:56.000000Z",
      "updated_at": "2020-09-11T12:01:56.000000Z"
    }
  ]
}
```

### HTTP Request

`GET /api/actions`

# Hooks

## List Hooks

> Response

```json
{
  "hooks": [
    [
        {
          "user_id": 1,
          "action_id": 2,
          "timing": "after",
          "name": "Build Frontend",
          "run_as": "forge",
          "script": "npm run prod",
          "sequence": 1,
          "project_id": 1,
          "updated_at": "2020-10-16T10:57:32.000000Z",
          "created_at": "2020-10-16T10:57:32.000000Z",
          "id": 1
        }
    ]
  ]
}
```

### HTTP Request

`GET /api/projects/{projectId}/hooks`

## Get Hook

> Response

```json
{
  "hook": {
    "user_id": 1,
    "action_id": 2,
    "timing": "after",
    "name": "Build Frontend",
    "run_as": "forge",
    "script": "npm run prod",
    "sequence": 1,
    "project_id": 1,
    "updated_at": "2020-10-16T10:57:32.000000Z",
    "created_at": "2020-10-16T10:57:32.000000Z",
    "id": 1
  }
}
```

### HTTP Request

`GET /api/projects/{projectId}/hooks/{hookId}`

## Create Hook

> Payload

```json
{
  "name": "Build Frontend",
  "script": "npm run prod",
  "runAs": "forge",
  "actionId": 2,
  "timing": "after",
  "servers": [1, 2]
}
```

> Response

```json
{
  "hook": {
    "user_id": 1,
    "action_id": 2,
    "timing": "after",
    "name": "Build Frontend",
    "run_as": "forge",
    "script": "npm run prod",
    "sequence": 1,
    "project_id": 1,
    "updated_at": "2020-10-16T10:57:32.000000Z",
    "created_at": "2020-10-16T10:57:32.000000Z",
    "id": 1
  }
}
```

### HTTP Request

`POST /api/projects/{projectId}/hooks`

### Required Scopes

This endpoint requires the `hooks:create` scope.

### Actions

You can find the ID of the action to hook into from the [Actions](#list-actions) page.

## Update Hook

> Payload

```json
{
  "servers": [2, 4, 5]
}
```

> Response

```json
{
  "hook": {
    "user_id": 1,
    "action_id": 2,
    "timing": "after",
    "name": "Build Frontend",
    "run_as": "forge",
    "script": "npm run prod",
    "sequence": 1,
    "project_id": 1,
    "updated_at": "2020-10-16T10:57:32.000000Z",
    "created_at": "2020-10-16T10:57:32.000000Z",
    "id": 1,
    "servers": [...]
  }
}
```

### HTTP Request

`PUT /api/projects/{projectId}/hooks/{hookId}`

### Required Scopes

This endpoint requires the `hooks:create` scope.

## Delete Hook

### HTTP Request

`DELETE /api/projects/{projectId}/hooks/{hookId}`

### Required Scopes

This endpoint requires the `hooks:delete` scope.

# Deployments

## List Deployments

> Response

```json
{
  "deployments": [
    {
      "id": 57,
      "project_id": 1,
      "user_id": 1,
      "commit_branch": "master",
      "commit_hash": "c66546e75fcbf208d2884b5ac7a3a858137753a3",
      "commit_author": "Dries Vints",
      "commit_avatar": "https:\/\/avatars1.githubusercontent.com\/u\/594614?v=4",
      "commit_message": "Update CHANGELOG.md",
      "status": "finished",
      "created_at": "2020-10-09T10:53:57.000000Z",
      "updated_at": "2020-10-09T10:54:38.000000Z"
    },
    {
      "id": 56,
      "project_id": 1,
      "user_id": 1,
      "commit_branch": "master",
      "commit_hash": "c66546e75fcbf208d2884b5ac7a3a858137753a3",
      "commit_author": "Dries Vints",
      "commit_avatar": "https:\/\/avatars1.githubusercontent.com\/u\/594614?v=4",
      "commit_message": "Update CHANGELOG.md",
      "status": "finished",
      "created_at": "2020-10-08T12:01:12.000000Z",
      "updated_at": "2020-10-08T12:01:48.000000Z"
    },
    {
      "id": 55,
      "project_id": 1,
      "user_id": 1,
      "commit_branch": "master",
      "commit_hash": "c66546e75fcbf208d2884b5ac7a3a858137753a3",
      "commit_author": "Dries Vints",
      "commit_avatar": "https:\/\/avatars1.githubusercontent.com\/u\/594614?v=4",
      "commit_message": "Update CHANGELOG.md",
      "status": "finished",
      "created_at": "2020-10-08T11:59:45.000000Z",
      "updated_at": "2020-10-08T12:00:30.000000Z"
    },
    {
      "id": 54,
      "project_id": 1,
      "user_id": 1,
      "commit_branch": "master",
      "commit_hash": "c66546e75fcbf208d2884b5ac7a3a858137753a3",
      "commit_author": "Dries Vints",
      "commit_avatar": "https:\/\/avatars1.githubusercontent.com\/u\/594614?v=4",
      "commit_message": "Update CHANGELOG.md",
      "status": "finished",
      "created_at": "2020-10-08T11:58:58.000000Z",
      "updated_at": "2020-10-08T11:59:54.000000Z"
    }
  ]
}
```

### HTTP Request

`GET /api/projects/{projectId}/deployments`

## Get Deployment

> Response

```json
{
  "deployment": {
    "id": 57,
    "project_id": 1,
    "user_id": 1,
    "commit_branch": "master",
    "commit_hash": "c66546e75fcbf208d2884b5ac7a3a858137753a3",
    "commit_author": "Dries Vints",
    "commit_avatar": "https:\/\/avatars1.githubusercontent.com\/u\/594614?v=4",
    "commit_message": "Update CHANGELOG.md",
    "status": "finished",
    "created_at": "2020-10-09T10:53:57.000000Z",
    "updated_at": "2020-10-09T10:54:38.000000Z",
    "processes": [
      {
        "id": 278,
        "deployment_id": 57,
        "project_id": 1,
        "server_id": 4,
        "server_name": "floral-moon",
        "sequence": 2,
        "name": "Clone New Release",
        "action_id": 1,
        "hook_id": null,
        "status": "finished",
        "started_at": 1602240838,
        "finished_at": 1602240839,
        "created_at": "2020-10-09T10:53:57.000000Z",
        "updated_at": "2020-10-09T10:53:59.000000Z",
        "server": null
      },
      {
        "id": 279,
        "deployment_id": 57,
        "project_id": 1,
        "server_id": 4,
        "server_name": "floral-moon",
        "sequence": 5,
        "name": "Install Composer Dependencies",
        "action_id": 2,
        "hook_id": null,
        "status": "finished",
        "started_at": 1602240843,
        "finished_at": 1602240869,
        "created_at": "2020-10-09T10:53:57.000000Z",
        "updated_at": "2020-10-09T10:54:29.000000Z",
        "server": null
      },
      {
        "id": 280,
        "deployment_id": 57,
        "project_id": 1,
        "server_id": 4,
        "server_name": "floral-moon",
        "sequence": 8,
        "name": "Activate New Release",
        "action_id": 3,
        "hook_id": null,
        "status": "finished",
        "started_at": 1602240870,
        "finished_at": 1602240874,
        "created_at": "2020-10-09T10:53:57.000000Z",
        "updated_at": "2020-10-09T10:54:34.000000Z",
        "server": null
      },
      {
        "id": 281,
        "deployment_id": 57,
        "project_id": 1,
        "server_id": 4,
        "server_name": "floral-moon",
        "sequence": 10,
        "name": "Purge Old Releases",
        "action_id": 4,
        "hook_id": null,
        "status": "finished",
        "started_at": 1602240874,
        "finished_at": 1602240875,
        "created_at": "2020-10-09T10:53:57.000000Z",
        "updated_at": "2020-10-09T10:54:35.000000Z",
        "server": null
      }
    ],
    "project": {
      "id": 1,
      "user_id": 1,
      "version": 1,
      "name": "Laravel",
      "provider": "github",
      "plain_repository": "laravel\/laravel",
      "repository": "git@github.com:laravel\/laravel.git",
      "type": "laravel-5",
      "branch": "master",
      "push_to_deploy": false,
      "webhook_id": null,
      "status": null,
      "should_deploy_again": 0,
      "deployment_started_at": null,
      "deployment_finished_at": "2020-10-09T10:54:38.000000Z",
      "last_deployment_status": "finished",
      "daily_deploys": 89,
      "weekly_deploys": 89,
      "last_deployment_took": 37,
      "retain_deployments": 4,
      "environment_servers": [
        4
      ],
      "folders": [],
      "monitor": "http:\/\/laravel.com",
      "new_york_status": "healthy",
      "london_status": "healthy",
      "singapore_status": "healthy",
      "token": "gJw0HbKUG0Zegum6oSNuYf0OvvHWy2z47eElFMiM",
      "created_at": "2020-09-10T13:56:28.000000Z",
      "updated_at": "2020-10-09T11:28:37.000000Z",
      "install_dev_dependencies": false,
      "install_dependencies": true,
      "quiet_composer": false,
      "has_environment": true,
      "has_monitoring_error": false,
      "has_missing_heartbeats": true,
      "last_deployed_branch": "master",
      "last_deployment_id": 57,
      "last_deployment_author": "Dries Vints",
      "last_deployment_avatar": "https:\/\/avatars1.githubusercontent.com\/u\/594614?v=4",
      "last_deployment_hash": "c66546e",
      "last_deployment_timestamp": "1 week ago",
      "servers": []
    }
  }
}
```

## Deploy Project

> Request

```json
{
    "from": "branch",
    "branch": "8.x"
}
```

### HTTP Request

`POST /api/projects/{projectId}/deployments`

### Required Scopes

This endpoint requires the `deployments:create` scope.

### Deployment Targets

By default, deployments will be made from the default configured branch. Alternatively, deployments can be made from either a `branch` or `tag`.

If you deploy from another target, you should supply the value of the `from` key as a new key with a value of the valid branch or tag:

> Deploying Via Tag

```json
{
    "from": "tag",
    "tag": "v8.4.1"
}
```

## Cancel Deployment

### HTTP Request

`DELETE /api/projects/{projectId}/deployments/{deploymentId}`

### Required Scopes

This endpoint requires the `deployments:delete` scope.

# Heartbeats

## List Heartbeats

> Response

```json
{
  "heartbeats": [
    {
      "id": 4,
      "project_id": 5,
      "name": "Test",
      "token": "jnqmHbtsNAoKaWz",
      "interval": 10,
      "status": "healthy",
      "last_checked_in_at": "2021-10-11T12:32:31.000000Z",
      "created_at": "2021-10-11T12:32:31.000000Z",
      "updated_at": "2021-10-11T12:32:31.000000Z"
    }
  ]
}
```

### HTTP Request

`GET /api/projects/{projectId}/heartbeats`

## Get Heartbeat

> Response

```json
{
  "heartbeat": {
    "id": 4,
    "project_id": 5,
    "name": "Test",
    "token": "jnqmHbtsNAoKaWz",
    "interval": 10,
    "status": "healthy",
    "last_checked_in_at": "2021-10-11T12:32:31.000000Z",
    "created_at": "2021-10-11T12:32:31.000000Z",
    "updated_at": "2021-10-11T12:32:31.000000Z"
  }
}
```

### HTTP Request

`GET /api/projects/{projectId}/heartbeats/{heartbeatId}`

## Create Heartbeat

> Payload

```json
{
  "name": "My Heartbeat",
  "interval": 10
}
```

> Response

```json
{
  "heartbeat": {
    "id": 4,
    "project_id": 5,
    "name": "Test",
    "token": "jnqmHbtsNAoKaWz",
    "interval": 10,
    "status": "healthy",
    "last_checked_in_at": "2021-10-11T12:32:31.000000Z",
    "created_at": "2021-10-11T12:32:31.000000Z",
    "updated_at": "2021-10-11T12:32:31.000000Z"
  }
}
```

### HTTP Request

`POST /api/projects/{projectId}/heartbeats`

### Required Scopes

This endpoint requires the `heartbeats:create` scope.

### Intervals

Intervals are minute based and can be one of the following:

*   10 (10 Minutes)
*   30 (30 Minutes)
*   60 (1 Hour)
*   1440 (1 Day)
*   10080 (1 Week)
*   44640 (1 Month)

## Delete Heartbeat

### HTTP Request

`DELETE /api/projects/{projectId}/heartbeats/{heartbeatId}`

### Required Scopes

This endpoint requires the `heartbeats:delete` scope.

# Collaborators

## List Collaborators

> Response

```json
{
  "collaborators": [
    {
      "id": 2,
      "name": "James Brooks",
      "email": "james@laravel.com",
      "pivot": {
        "project_id": 1,
        "collaborator_id": 2,
        "status": "invited"
      }
    }
  ]
}
```

### HTTP Request

`GET /api/projects/{projectId}/collaborators`

## Get Collaborator

> Response

```json
{
  "collaborator": {
    "id": 2,
    "name": "James Brooks",
    "email": "james@laravel.com",
    "pivot": {
      "project_id": 1,
      "collaborator_id": 2,
      "status": "invited"
    }
  }
}
```

### HTTP Request

`GET /api/projects/{projectId}/collaborators/{collaboratorId}`

## Create Collaborator

> Payload

```json
{
  "email": "mohamed@laravel.com"
}
```

> Response

```json
{
  "collaborators": [
    {
      "id": 2,
      "name": "James Brooks",
      "email": "james@laravel.com",
      "pivot": {
        "project_id": 1,
        "collaborator_id": 2,
        "status": "accepted"
      }
    }, {
      "id": 3,
      "name": "Mohamed Said",
      "email": "mohamed@laravel.com",
      "pivot": {
        "project_id": 1,
        "collaborator_id": 3,
        "status": "invited"
      }
    }
  ]
}
```

### HTTP Request

`POST /api/projects/{projectId}/collaborators`

### Required Scopes

This endpoint requires the `collaborators:create` scope.

## Delete Collaborator

### HTTP Request

> Payload

```json
{
  "email": "james@laravel.com"
}
```

`DELETE /api/projects/{projectId}/collaborators`

### Required Scopes

This endpoint requires the `collaborators:delete` scope.

# Notifications

## List Notifications

> Response

```json
{
  "notifications": [
    {
      "id": 6,
      "project_id": 1,
      "name": "Test",
      "type": "discord",
      "options": {
        "webhook": "https://discordapp.com/api/webhooks/some-discord-url"
      },
      "active": 1,
      "created_at": "2020-09-17T10:35:07.000000Z",
      "updated_at": "2020-09-17T10:35:07.000000Z"
    }
  ]
}
```

### HTTP Request

`GET /api/projects/{projectId}/notifications`

## Get Notification

> Response

```json
{
  "notification": {
    "id": 6,
    "project_id": 1,
    "name": "Test",
    "type": "discord",
    "options": {
      "webhook": "https://discordapp.com/api/webhooks/some-discord-url"
    },
    "active": 1,
    "created_at": "2020-09-17T10:35:07.000000Z",
    "updated_at": "2020-09-17T10:35:07.000000Z"
  }
}
```

### HTTP Request

`GET /api/projects/{projectId}/notifications/{notificationId}`

## Create Notification

> Payload

```json
{
  "name": "Taylor Otwell",
  "type": "email",
  "email_address": "taylor@laravel.com"
}
```

> Response

```json
{
  "notifications": [
    {
      "id": 6,
      "project_id": 1,
      "name": "Test",
      "type": "discord",
      "options": {
        "webhook": "https://discordapp.com/api/webhooks/some-discord-url"
      },
      "active": 1,
      "created_at": "2020-09-17T10:35:07.000000Z",
      "updated_at": "2020-09-17T10:35:07.000000Z"
    }, {
      "id": 7,
      "project_id": 1,
      "name": "Taylor Otwell",
      "type": "email",
      "options": {
        "email": "taylor@laravel.com"
      },
      "active": 1,
      "created_at": "2020-09-17T10:35:07.000000Z",
      "updated_at": "2020-09-17T10:35:07.000000Z"
    }
  ]
}
```

### HTTP Request

`POST /api/projects/{projectId}/notifications`

### Required Scopes

This endpoint requires the `notifications:create` scope.

### Notification Types

Envoyer supports several notification types:

*   `email`
*   `discord`
*   `slack`
*   `teams`

### Email

When creating an email notification, you **must** supply the `email_address` parameter.

### Slack

When creating a Slack notification, you **must** supply the `slack_webhook` parameter.

### Discord

When creating a Discord notification, you **must** supply the `discord_webhook` parameter.

### Microsoft Teams

When creating a Microsoft Teams notification, you **must** supply the `teams_webhook` parameter.

## Update Notification

> Payload

```json
{
  "name": "Taylor"
}
```

> Response

```json
{
 "notification": {
    "id": 7,
    "project_id": 1,
    "name": "Taylor",
    "type": "email",
    "options": {
      "email": "taylor@laravel.com"
    },
    "active": 1,
    "created_at": "2020-09-17T10:35:07.000000Z",
    "updated_at": "2020-09-17T10:35:07.000000Z"
  }
}
```

### HTTP Request

`PUT /api/projects/{projectId}/notifications/{notificationId}`

The request parameters are identical to the **Create Notification** endpoint.

### Required Scopes

This endpoint requires the `notifications:create` scope.

## Delete Notification

### HTTP Request

`DELETE /api/projects/{projectId}/notifications/{notificationId}`

### Required Scopes

This endpoint requires the `notifications:delete` scope.
