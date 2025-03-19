# Introduction

The Forge API allows you to create and interact with servers and sites on Laravel Forge through a simple REST API.

# Authentication

In order to use the API, you should authenticate your request by including your API key as a bearer token value:

`Authorization: Bearer API_KEY_HERE`

You may generate API keys in the [API Dashboard](https://forge.laravel.com/user-profile/api).

# Headers

Make sure you have the following content type headers are set on every request:

`Accept: application/json` `Content-Type: application/json`

# URI

Forge API is hosted on the following base URI:

`https://forge.laravel.com/api/v1`

# Errors

Forge uses conventional HTTP response codes to indicate the success or failure of an API request. The table below contains a summary of the typical response codes:

| Code | Description |
| --- | --- |
| 200 | Everything is ok. |
| 400 | Valid data was given but the request has failed. |
| 401 | No valid API Key was given. |
| 404 | The request resource could not be found. |
| 422 | The payload has missing required parameters or invalid data was given. |
| 429 | Too many attempts. |
| 500 | Request failed due to an internal error in Forge. |
| 503 | Forge is offline for maintenance. |

# User

## show

> Response

```json
{
    "user": {
        "id": 1,
        "name": "Mohamed Said",
        "email": "mail@gmail.com",
        "card_last_four": "1881",
        "connected_to_github": true,
        "connected_to_gitlab": true,
        "connected_to_bitbucket_two": true,
        "connected_to_digitalocean": true,
        "connected_to_linode": true,
        "connected_to_vultr": true,
        "connected_to_aws": true,
        "connected_to_hetzner": true,
        "ready_for_billing": true,
        "stripe_is_active": 1,
        "stripe_price": "yearly-basic-199-trial",
        "subscribed": 1,
        "can_create_servers": true,
        "2fa_enabled": false
    }
}
```

### HTTP Request

`GET /api/v1/user`

# Servers

## Create Server

> Payload

```json
{
    "provider": "ocean2",
    "ubuntu_version": "22.04",
    "type": "web",
    "credential_id": 1,
    "name": "test-via-api",
    "size": "01",
    "database": "test123",
    "php_version": "php82",
    "region": "ams2",
    "recipe_id": null
}
```

> Response

```json
{
    "server": {
        "id": 16,
        "credential_id": 1,
        "name": "test-via-api",
        "type": "web",
        "size": "01",
        "region": "ams2",
        "php_version": "php82",
        "php_cli_version": "php82",
        "opcache_status": "enabled",
        "database_type": "mysql8",
        "ip_address": null,
        "private_ip_address": null,
        "blackfire_status": null,
        "papertrail_status": null,
        "revoked": false,
        "created_at": "2016-12-15 15:04:05",
        "is_ready": false,
        "network": []
    },
    "sudo_password": "baracoda",
    "database_password": "spotted_eagle_ray"
}
```

### HTTP Request

`POST /api/v1/servers`

### Parameters

| Key | Description |
| --- | --- |
| ubuntu\_version | The version of Ubuntu to create the server with. Valid values are `"20.04"`, `"22.04"`, and `"24.04"`. `"24.04"` is used by default if no value is defined. It is recommended to always specify a version as the default may change at any time. |
| type | The type of server to create. Valid values are `app`, `web`, `loadbalancer`, `cache`, `database`, `worker`, `meilisearch`. `app` is used by default if no value is defined. |
| provider | The server provider. Valid values are `ocean2` for Digital Ocean, `akamai` (Linode), `vultr2`, `aws`, `hetzner` and `custom`. |
| disk\_size | The size of the disk in GB. Valid when the provider is `aws`. Minimum of 8GB. Example: `20`. |
| circle | The ID of a circle to create the server within. |
| credential\_id | This is only required when the provider is not `custom`. |
| region | The name of the region where the server will be created. This value is not required you are building a Custom VPS server. [Valid region identifiers](/api-documentation#regions). |
| ip\_address | The IP Address of the server. Only required when the provider is `custom`. |
| private\_ip\_address | The Private IP Address of the server. Only required when the provider is `custom`. |
| php\_version | Valid values are `php84`, `php83`, `php82`, `php81`, `php80`, `php74`, `php73`,`php72`,`php82`, `php70`, and `php56`. |
| database | The name of the database Forge should create when building the server. If omitted, `forge` will be used. |
| database\_type | Valid values are `mysql8`, `mariadb106`, `mariadb1011`, `mariadb114`, `postgres`, `postgres13`, `postgres14`, `postgres15`, `postgres16` or `postgres17`. |
| network | An array of server IDs that the server should be able to connect to. |
| recipe\_id | An optional ID of a recipe to run after provisioning. |
| aws\_vpc\_id | ID of the existing VPC |
| aws\_subnet\_id | ID of the existing subnet |
| aws\_vpc\_name | When creating a new one |
| hetzner\_network\_id | ID of the existing VPC |
| ocean2\_vpc\_uuid | UUID of the existing VPC |
| ocean2\_vpc\_name | When creating a new one |
| vultr2\_network\_id | ID of the existing private network |
| vultr2\_network\_name | When creating a new one |

### Server Status

Servers take about 10 minutes to provision. Once the server is ready to be used, the `is_ready` parameter on the server will be `true`. **You should not repeatedly ping the Forge API asking if the server is ready. Instead, consider pinging the endpoint once every 2 minutes.**

### Valid Sizes

Check the [Regions endpoint](https://forge.laravel.com/api-documentation#regions) for the available sizes and regions IDs.

### Custom VPS

While creating a custom VPS, the response of this endpoint will contain a `provision_command` attribute:

```json
{
    "provision_command": "wget -O forge.sh https://..."
}
```

## List Servers

> Response

```json
{
    "servers": [
       {
            "id": 1,
            "credential_id": 1,
            "name": "test-via-api",
            "size": "s-1vcpu-1gb",
            "region": "Amsterdam 2",
            "php_version": "php82",
            "php_cli_version": "php82",
            "opcache_status": "enabled",
            "database_type": "mysql8",
            "ip_address": "37.139.3.148",
            "private_ip_address": "10.129.3.252",
            "blackfire_status": null,
            "papertrail_status": null,
            "revoked": false,
            "created_at": "2016-12-15 18:38:18",
            "is_ready": true,
            "network": []
        }
    ]
}
```

### HTTP Request

`GET /api/v1/servers`

## Get Server

> Response

```json
{
    "server": {
        "id": 1,
        "credential_id": 1,
        "name": "test-via-api",
        "size": "s-1vcpu-1gb",
        "region": "Amsterdam 2",
        "php_version": "php82",
        "php_cli_version": "php82",
        "opcache_status": "enabled",
        "database_type": "mysql8",
        "ip_address": "37.139.3.148",
        "private_ip_address": "10.129.3.252",
        "blackfire_status": null,
        "papertrail_status": null,
        "revoked": false,
        "created_at": "2016-12-15 18:38:18",
        "is_ready": true,
        "network": []
    }
}
```

### HTTP Request

`GET /api/v1/servers/{id}`

## Update Server

> Payload

```json
{
    "name": "renamed-server",
    "ip_address": "192.241.143.108",
    "private_ip_address": "10.136.8.40",
    "max_upload_size": 123,
    "max_execution_time": 30,
    "network": [
        2,
        3
    ],
    "timezone": "Europe/London",
    "tags": [
        "london-server"
    ]
}
```

> Response

```json
{
    "server": {
        "id": 16,
        "credential_id": 1,
        "name": "test-via-api",
        "size": "s-1vcpu-1gb",
        "region": "Amsterdam 2",
        "php_version": "php82",
        "php_cli_version": "php82",
        "opcache_status": "enabled",
        "database_type": "mysql8",
        "ip_address": null,
        "private_ip_address": null,
        "blackfire_status": null,
        "papertrail_status": null,
        "revoked": false,
        "created_at": "2016-12-15 15:04:05",
        "is_ready": false,
        "network": [
            2,
            3
        ],
        "tags": [
            "london-server"
        ]
    }
}
```

### HTTP Request

`PUT /api/v1/servers/{id}`

## Update Database Password

> Payload

```json
{
    "password": "maeve"
}
```

### HTTP Request

`PUT /api/v1/servers/{serverId}/database-password`

This endpoint will update Forge's copy of the primary database password which should be used to authenticate the creation of new databases and database users. This is typically only needed if you are working with a Forge server that was built before database administration was added to Forge.

## Delete Server

### HTTP Request

`DELETE /api/v1/servers/{id}`

## Reboot Server

### HTTP Request

`POST /api/v1/servers/{id}/reboot`

## Revoke Forge access to server

### HTTP Request

`POST /api/v1/servers/{id}/revoke`

## Reconnect revoked server

> Response

```json
{
    "public_key": "CONTENT_OF_THE_PUBLIC_KEY"
}
```

### HTTP Request

`POST /api/v1/servers/{id}/reconnect`

This endpoint will return an SSH key which you will need to add to the server. Once the key has been added to the server, you may "reactivate" it.

## Reactivate revoked server

### HTTP Request

`POST /api/v1/servers/{id}/reactivate`

## Get Recent Events

> Response

```json
[
   {
        "server_id": 18,
        "ran_as": "forge",
        "server_name": "billowing-cliff",
        "description": "Deploying PHP Info Page.",
        "created_at": "2017-04-28 18:08:44"
    }
]
```

### HTTP Request

`GET /api/v1/servers/events`

### Parameters

| Key | Description |
| --- | --- |
| server\_id | Optionally specify a server\_id to get recent events of only this server. |

## Get Server Events

> Response

```json
{
    "events": [
       {
            "server_id": 18,
            "ran_as": "forge",
            "server_name": "billowing-cliff",
            "description": "Deploying PHP Info Page.",
            "created_at": "2017-04-28 18:08:44"
        }
    ]
}
```

### HTTP Request

`GET /api/v1/servers/{id}/events`

## Get Server Event Output

> Response

```json
{
    "output": "Some command output."
}
```

### HTTP Request

`GET /api/v1/servers/{id}/events/{event}`

# Services

## Start Service

### HTTP Request

`POST /api/v1/servers/{id}/services/start`

> Payload

```json
{
    "service": "service-name"
}
```

## Stop Service

### HTTP Request

`POST /api/v1/servers/{id}/services/stop`

> Payload

```json
{
    "service": "service-name"
}
```

## Restart Service

### HTTP Request

`POST /api/v1/servers/{id}/services/restart`

> Payload

```json
{
    "service": "service-name"
}
```

## Reboot MySQL

### HTTP Request

`POST /api/v1/servers/{id}/mysql/reboot`

## Stop MySQL

### HTTP Request

`POST /api/v1/servers/{id}/mysql/stop`

## Reboot Nginx

### HTTP Request

`POST /api/v1/servers/{id}/nginx/reboot`

## Stop Nginx

### HTTP Request

`POST /api/v1/servers/{id}/nginx/stop`

## Test Nginx

> Response

```json
{
    "result": ⊕"nginx: [emerg] a duplicate lis ..."⊖"nginx: [emerg] a duplicate listen 0.0.0.0:80 in /etc/nginx/sites-enabled/default:6\nnginx: configuration file /etc/nginx/nginx.conf test failed\n"
}
```

### HTTP Request

`GET /api/v1/servers/{id}/nginx/test`

## Reboot Postgres

### HTTP Request

`POST /api/v1/servers/{id}/postgres/reboot`

## Stop Postgres

### HTTP Request

`POST /api/v1/servers/{id}/postgres/stop`

## Reboot PHP

> Payload

```json
{
    "version": "php74"
}
```

### HTTP Request

`POST /api/v1/servers/{id}/php/reboot`

## Install Blackfire

> Payload

```json
{
    "server_id": "...",
    "server_token": "..."
}
```

### HTTP Request

`POST /api/v1/servers/{id}/blackfire/install`

## Remove Blackfire

### HTTP Request

`DELETE /api/v1/servers/{id}/blackfire/remove`

## Install Papertrail

> Payload

```json
{
    "host": "192.241.143.108"
}
```

### HTTP Request

`POST /api/v1/servers/{id}/papertrail/install`

## Remove Papertrail

### HTTP Request

`DELETE /api/v1/servers/{id}/papertrail/remove`

# Daemons

## Create Daemon

> Payload

```json
{
    "command": "COMMAND",
    "user": "root",
    "directory": "/home/forge/foo.com"
}
```

> Response

```json
{
    "daemon": {
        "id": 2,
        "command": "php artisan queue:work",
        "user": "forge",
        "directory": "/home/forge/foo.com",
        "processes": 1,
        "startsecs": 1,
        "stopwaitsecs": 10,
        "stopsignal": "SIGTERM",
        "status": "installing",
        "created_at": "2022-02-14 09:29:18"
    }
}
```

### HTTP Request

`POST /api/v1/servers/{serverId}/daemons`

## List Daemons

> Response

```json
{
    "daemons": [
       {
            "id": 2,
            "command": "php artisan queue:work",
            "user": "forge",
            "directory": "/home/forge/foo.com",
            "processes": 1,
            "startsecs": 1,
            "stopwaitsecs": 10,
            "stopsignal": "SIGTERM",
            "status": "installing",
            "created_at": "2022-02-14 09:29:18"
        }
    ]
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/daemons`

## Get Daemon

> Response

```json
{
    "daemon": {
        "id": 2,
        "command": "php artisan queue:work",
        "user": "forge",
        "directory": "/home/forge/foo.com",
        "processes": 1,
        "startsecs": 1,
        "stopwaitsecs": 10,
        "stopsignal": "SIGTERM",
        "status": "installing",
        "created_at": "2022-02-14 09:29:18"
    }
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/daemons/{daemonId}`

## Delete Daemon

### HTTP Request

`DELETE /api/v1/servers/{serverId}/daemons/{daemonId}`

## Restart Daemon

### HTTP Request

`POST /api/v1/servers/{serverId}/daemons/{daemonId}/restart`

# Firewall Rules

## Create Rule

> Payload

```json
{
    "name": "rule name",
    "ip_address": "192.168.1.1",
    "port": 88,
    "type": "allow"
}
```

> Response

```json
{
    "rule": {
        "id": 4,
        "name": "rule",
        "port": 123,
        "type": "allow",
        "ip_address": null,
        "status": "installing",
        "created_at": "2016-12-16 15:50:17"
    }
}
```

### HTTP Request

`POST /api/v1/servers/{serverId}/firewall-rules`

### Available Rule Types

You may specify `allow` or `deny` as the rule type.

## List Rules

> Response

```json
{
    "rules": [
       {
            "id": 4,
            "name": "rule",
            "port": 123,
            "type": "allow",
            "ip_address": null,
            "status": "installing",
            "created_at": "2016-12-16 15:50:17"
        }
    ]
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/firewall-rules`

## Get Rule

> Response

```json
{
    "rule": {
        "id": 4,
        "name": "rule",
        "port": 123,
        "type": "allow",
        "ip_address": null,
        "status": "installing",
        "created_at": "2016-12-16 15:50:17"
    }
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/firewall-rules/{ruleId}`

## Delete Rule

### HTTP Request

`DELETE /api/v1/servers/{serverId}/firewall-rules/{ruleId}`

# Scheduled Jobs

## Create Job

> Payload

```json
{
    "command": "COMMAND_THE_JOB_RUNS",
    "frequency": "custom",
    "user": "root",
    "minute": "*",
    "hour": "*",
    "day": "*",
    "month": "*",
    "weekday": "*"
}
```

> Response

```json
{
    "job": {
        "id": 2,
        "command": "COMMAND_THE_JOB_RUNS",
        "user": "root",
        "frequency": "Nightly",
        "cron": "0 0 * * *",
        "status": "installing",
        "created_at": "2016-12-16 15:56:59"
    }
}
```

### HTTP Request

`POST /api/v1/servers/{serverId}/jobs`

### Parameters

| Key | Description |
| --- | --- |
| frequency | The frequency in which the job should run. Valid values are `minutely`, `hourly`, `nightly`, `weekly`, `monthly`, `reboot`, and `custom` |
| minute | Required if the frequency is `custom`. |
| hour | Required if the frequency is `custom`. |
| day | Required if the frequency is `custom`. |
| month | Required if the frequency is `custom`. |
| weekday | Required if the frequency is `custom`. |

## List Jobs

> Response

```json
{
    "jobs": [
       {
            "id": 2,
            "command": "COMMAND_THE_JOB_RUNS",
            "user": "root",
            "frequency": "nightly",
            "cron": "0 0 * * *",
            "status": "installing",
            "created_at": "2016-12-16 15:56:59"
        }
    ]
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/jobs`

## Get Job

> Response

```json
{
    "job": {
        "id": 2,
        "command": "COMMAND_THE_JOB_RUNS",
        "user": "root",
        "frequency": "Nightly",
        "cron": "0 0 * * *",
        "status": "installing",
        "created_at": "2016-12-16 15:56:59"
    }
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/jobs/{jobId}`

## Delete Job

### HTTP Request

`DELETE /api/v1/servers/{serverId}/jobs/{jobId}`

## Get Job Output

> Response

```json
{
    "output": ⊕"The output of the job will be  ..."⊖"The output of the job will be returned here."
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/jobs/{jobId}/output`

# PHP

## List PHP Versions

> Response

```json
[
      {
            "id": 29,
            "version": "php74",
            "status": "installed",
            "displayable_version": "PHP 7.4",
            "binary_name": "php7.4",
            "used_as_default": false,
            "used_on_cli": false
      },
      {
            "id": 30,
            "version": "php73",
            "status": "installed",
            "displayable_version": "PHP 7.3",
            "binary_name": "php7.3",
            "used_as_default": true,
            "used_on_cli": true
      },
      {
            "id": 31,
            "version": "php72",
            "status": "installed",
            "displayable_version": "PHP 7.2",
            "binary_name": "php7.2",
            "used_as_default": false,
            "used_on_cli": false
      },
      {
            "id": 32,
            "version": "php71",
            "status": "installed",
            "displayable_version": "PHP 7.1",
            "binary_name": "php7.1",
            "used_as_default": false,
            "used_on_cli": false
      },
      {
            "id": 33,
            "version": "php56",
            "status": "installed",
            "displayable_version": "PHP 5.6",
            "binary_name": "php5.6",
            "used_as_default": false,
            "used_on_cli": false
      }
]
```

### HTTP Request

`GET /api/v1/servers/{serverId}/php`

## Install PHP Version

### HTTP Request

`POST /api/v1/servers/{serverId}/php`

> Payload

```json
{
    "version": "php74"
}
```

### Available Versions

| Key | Description |
| --- | --- |
| `php84` | PHP 8.4 |
| `php83` | PHP 8.3 |
| `php82` | PHP 8.2 |
| `php81` | PHP 8.1 |
| `php80` | PHP 8.0 |
| `php74` | PHP 7.4 |
| `php73` | PHP 7.3 |
| `php72` | PHP 7.2 |
| `php71` | PHP 7.1 |
| `php70` | PHP 7.0 |
| `php56` | PHP 5.6 |

## Upgrade PHP Patch Version

You must supply the version to be patched.

### HTTP Request

`POST /api/v1/servers/{serverId}/php/update`

> Payload

```json
{
    "version": "php74"
}
```

## Enable OPCache

### HTTP Request

`POST /api/v1/servers/{serverId}/php/opcache`

## Disable OPCache

### HTTP Request

`DELETE /api/v1/servers/{serverId}/php/opcache`

# Databases

## Create Database

> Payload

```json
{
    "name": "forge",
    "user": "forge",
    "password": "dolores"
}
```

> Response

```json
{
    "database": {
        "id": 1,
        "name": "forge",
        "status": "installing",
        "created_at": "2016-12-16 16:12:22"
    }
}
```

### HTTP Request

`POST /api/v1/servers/{serverId}/databases`

### Parameters

| Key | Description |
| --- | --- |
| user | This field is optional. If passed, it will be used to create a new Database User with access to the newly created database. |
| password | This field is only required when a `user` value is given. |

## Sync Database

### HTTP Request

`POST /api/v1/servers/{serverId}/databases/sync`

## List Databases

> Response

```json
{
    "databases": [
        {
            "id": 1,
            "name": "forge",
            "status": "installing",
            "created_at": "2016-12-16 16:12:22"
        }
    ]
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/databases`

## Get Database

> Response

```json
{
    "database": {
        "id": 1,
        "name": "forge",
        "status": "installing",
        "created_at": "2016-12-16 16:12:22"
    }
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/databases/{databaseId}`

## Delete Database

### HTTP Request

`DELETE /api/v1/servers/{serverId}/databases/{databaseId}`

# Database Users

## Create User

> Payload

```json
{
    "name": "forge",
    "password": "dolores",
    "databases": [1]
}
```

> Response

```json
{
    "user": {
        "id": 2,
        "name": "forge",
        "status": "installing",
        "created_at": "2016-12-16 16:19:01",
        "databases": [
            1
        ]
    }
}
```

### HTTP Request

`POST /api/v1/servers/{serverId}/database-users`

### Parameters

| Key | Description |
| --- | --- |
| databases | An array of database IDs referencing the databases the user has access to. |
| password | The password to assign the user. |

## List Users

> Response

```json
{
    "users": [
        {
            "id": 2,
            "name": "forge",
            "status": "installing",
            "created_at": "2016-12-16 16:19:01",
            "databases": [
                1
            ]
        }
    ]
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/database-users`

## Get User

> Response

```json
{
    "user": {
        "id": 2,
        "name": "forge",
        "status": "installing",
        "created_at": "2016-12-16 16:19:01",
        "databases": [
            1
        ]
    }
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/database-users/{userId}`

## Update User

> Payload

```json
{
    "databases": [2]
}
```

> Response

```json
{
    "user": {
        "id": 2,
        "name": "forge",
        "status": "installing",
        "created_at": "2016-12-16 16:19:01",
        "databases": [
            1
        ]
    }
}
```

### HTTP Request

`PUT /api/v1/servers/{serverId}/database-users/{userId}`

This endpoint may be used to update the databases the Database User has access to.

## Delete User

### HTTP Request

`DELETE /api/v1/servers/{serverId}/database-users/{userId}`

# MySQL Databases

The `/mysql` endpoints is now deprecated in favour of the new [`/databases`](#databases) endpoint.

## Create Database

> Payload

```json
{
    "name": "forge",
    "user": "forge",
    "password": "dolores"
}
```

> Response

```json
{
    "database": {
        "id": 1,
        "name": "forge",
        "status": "installing",
        "created_at": "2016-12-16 16:12:22"
    }
}
```

### HTTP Request

`POST /api/v1/servers/{serverId}/mysql`

### Parameters

| Key | Description |
| --- | --- |
| user | This field is optional. If passed, it will be used to create a new MySQL user with access to the newly created database. |
| password | This field is only required when a `user` value is given. |

## List Databases

> Response

```json
{
    "databases": [
        {
            "id": 1,
            "name": "forge",
            "status": "installing",
            "created_at": "2016-12-16 16:12:22"
        }
    ]
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/mysql`

## Get Database

> Response

```json
{
    "database": {
        "id": 1,
        "name": "forge",
        "status": "installing",
        "created_at": "2016-12-16 16:12:22"
    }
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/mysql/{databaseId}`

## Delete Database

### HTTP Request

`DELETE /api/v1/servers/{serverId}/mysql/{databaseId}`

# MySQL Database Users

The `/mysql-users` endpoint is now deprecated in favour of the new [`/database-users`](#database-users) endpoint.

## Create User

> Payload

```json
{
    "name": "forge",
    "password": "dolores",
    "databases": [1]
}
```

> Response

```json
{
    "user": {
        "id": 2,
        "name": "forge",
        "status": "installing",
        "created_at": "2016-12-16 16:19:01",
        "databases": [
            1
        ]
    }
}
```

### HTTP Request

`POST /api/v1/servers/{serverId}/mysql-users`

### Parameters

| Key | Description |
| --- | --- |
| databases | An array of database IDs referencing the databases the user has access to. |

## List Users

> Response

```json
{
    "users": [
        {
            "id": 2,
            "name": "forge",
            "status": "installing",
            "created_at": "2016-12-16 16:19:01",
            "databases": [
                1
            ]
        }
    ]
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/mysql-users`

## Get User

> Response

```json
{
    "user": {
        "id": 2,
        "name": "forge",
        "status": "installing",
        "created_at": "2016-12-16 16:19:01",
        "databases": [
            1
        ]
    }
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/mysql-users/{userId}`

## Update User

> Payload

```json
{
    "databases": [2]
}
```

> Response

```json
{
    "user": {
        "id": 2,
        "name": "forge",
        "status": "installing",
        "created_at": "2016-12-16 16:19:01",
        "databases": [
            1
        ]
    }
}
```

### HTTP Request

`PUT /api/v1/servers/{serverId}/mysql-users/{userId}`

This endpoint may be used to update the databases the MySQL user has access to.

## Delete User

### HTTP Request

`DELETE /api/v1/servers/{serverId}/mysql-users/{userId}`

# Nginx Templates

## Create Template

> Payload

```json
{
    "name": "My Nginx Template",
    "content": "server { listen {{ PORT }}; location = / { ... } }"
}
```

> Response

```json
{
    "template": {
        "id": 1,
        "server_id": 50,
        "name": "My Nginx Template",
        "content": "server { listen {{ PORT }}; location = / { ... } }"
    }
}
```

### HTTP Request

`POST /api/v1/servers/{serverId}/nginx/templates`

### Variables

Nginx templates support multiple variables that will be replaced with real data when the site is being created. For more information on variables, see the [Nginx templates documentation](https://forge.laravel.com/docs/servers/nginx-templates.html#template-variables).

## List Nginx Templates

> Response

```json
{
    "templates": [
        {
            "id": 1,
            "server_id": 50,
            "name": "My Nginx Template",
            "content": "server { listen {{ PORT }}; location = / { ... } }"
        }
    ]
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/nginx/templates/default`

## Get Default Nginx Template

When a site is created without a custom Nginx template selected, this is the Nginx configuration that Forge will use.

> Response

```json
{
    "template": {
        "server_id": 50,
        "name": "Forge Default",
        "content": "..."
    }
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/nginx/templates/{templateId}`

## Get Nginx Template

> Response

```json
{
    "template": {
        "id": 1,
        "server_id": 50,
        "name": "My Nginx Template",
        "content": "server { listen {{ PORT }}; location = / { ... } }"
    }
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/nginx/templates/{templateId}`

## Update Nginx Template

> Request

```json
{
    "name": "My New Name",
    "content": "My new content"
}
```

> Response

```json
{
    "template": {
        "id": 1,
        "server_id": 50,
        "name": "My New Name",
        "content": "My new content"
    }
}
```

### HTTP Request

`PUT /api/v1/servers/{serverId}/nginx/templates/{templateId}`

## Delete Nginx Template

### HTTP Request

`DELETE /api/v1/servers/{serverId}/nginx/templates/{templateId}`

# Sites

## Create Site

> Payload

```json
{
    "domain": "site.com",
    "project_type": "php",
    "aliases": ["alias1.com", "alias2.com"],
    "directory": "/test",
    "isolated": true,
    "username": "laravel",
    "database": "site-com-db",
    "php_version": "php81",
    "nginx_template": 1
}
```

> Response

```json
{
    "site": {
        "id": 2,
        "name": "site.com",
        "aliases": ["alias1.com", "alias2.com"],
        "directory": "/test",
        "wildcards": false,
        "isolated": true,
        "username": "forge",
        "status": "installing",
        "repository": null,
        "repository_provider": null,
        "repository_branch": null,
        "repository_status": null,
        "quick_deploy": false,
        "project_type": "php",
        "app": null,
        "php_version": "php81",
        "app_status": null,
        "slack_channel": null,
        "telegram_chat_id": null,
        "telegram_chat_title": null,
        "created_at": "2016-12-16 16:38:08",
        "deployment_url": "...",
        "tags": [],
        "web_directory": "/home/forge/torphy.com/public"
    }
}
```

### HTTP Request

`POST /api/v1/servers/{serverId}/sites`

### Available site types

| Key | Description |
| --- | --- |
| `php` | PHP / Laravel / Symfony |
| `html` | Static HTML / Nuxt.js / Next.js |

### Nginx Templates

You may leave the `nginx_template` key off to use the `default` template, or supply `nginx_template: "default"`.

## List Sites

> Response

```json
{
    "sites": [
        {
            "id": 2,
            "name": "site.com",
            "username": "laravel",
            "directory": "/test",
            "wildcards": false,
            "status": "installing",
            "repository": null,
            "repository_provider": null,
            "repository_branch": null,
            "repository_status": null,
            "quick_deploy": false,
            "project_type": "php",
            "app": null,
            "php_version": "php81",
            "app_status": null,
            "slack_channel": null,
            "telegram_chat_id": null,
            "telegram_chat_title": null,
            "deployment_url": "...",
            "created_at": "2016-12-16 16:38:08",
            "tags": [],
            "web_directory": "/home/forge/torphy.com/public"
        }
    ]
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/sites`

## Get Site

> Response

```json
{
    "site": {
        "id": 2,
        "name": "site.com",
        "aliases": ["alias1.com"],
        "username": "laravel",
        "directory": "/test",
        "wildcards": false,
        "status": "installing",
        "repository": null,
        "repository_provider": null,
        "repository_branch": null,
        "repository_status": null,
        "quick_deploy": false,
        "project_type": "php",
        "app": null,
        "php_version": "php81",
        "app_status": null,
        "slack_channel": null,
        "telegram_chat_id": null,
        "telegram_chat_title": null,
        "deployment_url": "...",
        "created_at": "2016-12-16 16:38:08",
        "tags": [],
        "web_directory": "/home/forge/torphy.com/public"
    }
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/sites/{siteId}`

## Update Site

> Payload

```json
{
    "directory": "/some/path",
    "name": "site-new-name.com",
    "php_version": "php81",
    "aliases": ["alias1.com", "alias2.com"],
    "wildcards": true
}
```

> Response

```json
{
    "site": {
        "id": 2,
        "name": "site-new-name.com",
        "aliases": ["alias1.com", "alias2.com"],
        "username": "laravel",
        "directory": "/some/path",
        "wildcards": false,
        "status": "installing",
        "repository": null,
        "repository_provider": null,
        "repository_branch": null,
        "repository_status": null,
        "quick_deploy": false,
        "project_type": "php",
        "app": null,
        "app_status": null,
        "slack_channel": null,
        "telegram_chat_id": null,
        "telegram_chat_title": null,
        "deployment_url": "...",
        "created_at": "2016-12-16 16:38:08",
        "tags": [],
        "web_directory": "/home/forge/torphy.com/public"
    }
}
```

### HTTP Request

`PUT /api/v1/servers/{serverId}/sites/{siteId}`

This endpoint is used to update the "web directory", primary name, aliases or whether to use wildcard sub-domains for a given site.

## Change Site PHP Version

> Payload

```json
{
    "version": "php74"
}
```

### HTTP Request

`PUT /api/v1/servers/{serverId}/sites/{siteId}/php`

## Add Site Aliases

> Payload

```json
{
    "aliases": ["alias1.com", "alias2.com"]
}
```

> Response

```json
{
    "site": {
        "id": 2,
        "name": "site.com",
        "aliases": ["alias1.com", "alias2.com"],
        "username": "laravel",
        "directory": "/",
        "wildcards": false,
        "status": "installing",
        "repository": null,
        "repository_provider": null,
        "repository_branch": null,
        "repository_status": null,
        "quick_deploy": false,
        "project_type": "php",
        "app": null,
        "app_status": null,
        "slack_channel": null,
        "telegram_chat_id": null,
        "telegram_chat_title": null,
        "deployment_url": "...",
        "created_at": "2016-12-16 16:38:08",
        "tags": [],
        "web_directory": "/home/forge/torphy.com/public"
    }
}
```

### HTTP Request

`PUT /api/v1/servers/{serverId}/sites/{siteId}/aliases`

Use this endpoint to add additional site aliases and keep the existing ones.

## Delete Site

### HTTP Request

`DELETE /api/v1/servers/{serverId}/sites/{siteId}`

## Load Balancing

> Response

```json
{
  "nodes": [
    {
      "server_id": 2,
      "weight": 5,
      "down": false,
      "backup": false,
      "port": 80
    }, {
      "server_id": 3,
      "weight": 1,
      "down": false,
      "backup": true,
      "port": 80
    }, {
      "server_id": 4,
      "weight": 1,
      "down": true,
      "backup": false,
      "port": 80
    }
  ]
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/sites/{siteId}/balancing`

## Update Load Balancing

> Payload

```json
{
    "servers": [{
        "id": 2,
        "weight": 5
    }, {
        "id": 3,
        "backup": true
    }, {
        "id": 4,
        "down": true
    }],
    "method": "least_conn"
}
```

### HTTP Request

`PUT /api/v1/servers/{serverId}/sites/{siteId}/balancing`

If the server is a load balancer, this endpoint may be used to specify the servers the load balancer should send traffic to.

#### Load Balancing Methods

| Key | Description |
| --- | --- |
| `round_robin` | Requests are evenly distributed across servers |
| `least_conn` | Requests are sent to the server with the least number of active connections |
| `ip_hash` | The server to which a request is sent is determined from the client IP address |

## Site Log

### HTTP Request

`GET /api/v1/servers/{serverId}/sites/{siteId}/logs`

> Response

```json
{
    "content": "[2020-08-18 10:32:56] local.INFO: Test  \n"
}
```

## Clear Site Log

### HTTP Request

`DELETE /api/v1/servers/{serverId}/sites/{siteId}/logs`

> Response (204)

# SSL Certificates

## Create Certificate

> Payload

```json
{
    "type": "new",
    "domain": "domain.com",
    "country": "US",
    "state": "NY",
    "city": "New York",
    "organization": "Company Name",
    "department": "IT"
}
```

> Response

```json
{
    "certificate": {
        "domain": "domain.com",
        "request_status": "creating",
        "created_at": "2016-12-17 07:02:35",
        "id": 3,
        "existing": false,
        "active": false
    }
}
```

### HTTP Request

`POST /api/v1/servers/{serverId}/sites/{siteId}/certificates`

## Installing An Existing Certificate

> Payload

```json
{
    "type": "existing",
    "key": "PRIVATE_KEY_HERE",
    "certificate": "CERTIFICATE_HERE"
}
```

> Response

```json
{
    "certificate": {
        "domain": "domain.com",
        "request_status": "creating",
        "created_at": "2016-12-17 07:02:35",
        "id": 3,
        "existing": false,
        "active": false
    }
}
```

### HTTP Request

`POST /api/v1/servers/{serverId}/sites/{siteId}/certificates`

## Cloning An Existing Certificate

> Payload

```json
{
    "type": "clone",
    "certificate_id": 1
}
```

> Response

```json
{
    "certificate": {
        "domain": "domain.com",
        "request_status": "creating",
        "created_at": "2016-12-17 07:02:35",
        "id": 3,
        "existing": false,
        "active": false
    }
}
```

### HTTP Request

`POST /api/v1/servers/{serverId}/sites/{siteId}/certificates`

## Obtain A LetsEncrypt Certificate

The `dns_provider` object is only required for wildcard sub-domains.

> Payload

```json
{
    "domains": ["www.site.com"],
    "dns_provider": {
        "type": "xxx",
        "cloudflare_api_token": "xxx",
        "route53_key": "xxx",
        "route53_secret": "xxx",
        "digitalocean_token": "xxx",
        "dnssimple_token": "xxx",
        "linode_token": "xxx",
        "ovh_endpoint": "xxx",
        "ovh_app_key": "xxx",
        "ovh_app_secret": "xxx",
        "ovh_consumer_key": "xxx",
        "google_credentials_file": "xxx",
    }
}
```

> Response

```json
{
    "certificate": {
        "domain": "www.test.com",
        "type": "letsencrypt",
        "request_status": "created",
        "status": "installing",
        "created_at": "2017-02-09 17:14:34",
        "id": 1,
        "existing": true,
        "active": false
    }
}
```

### HTTP Request

`POST /api/v1/servers/{serverId}/sites/{siteId}/certificates/letsencrypt`

### DNS Provider Types

| Type | Extra Fields |
| --- | --- |
| `cloudflare` | `cloudflare_api_token`. |
| `route53` | `route53_key` and `route53_secret` |
| `digitalocean` | `digitalocean_token` |
| `dnssimple` | `dnssimple_token` |
| `linode` | `linode_token` |
| `ovh` | `ovh_endpoint`, `ovh_app_key`, `ovh_app_secret` and `ovh_consumer_key` |
| `google` | `google_credentials_file` |

## List Certificates

> Response

```json
{
    "certificates": [
        {
            "domain": "domain.com",
            "request_status": "creating",
            "created_at": "2016-12-17 07:02:35",
            "id": 3,
            "existing": false,
            "active": false
        }
    ]
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/sites/{siteId}/certificates`

## Get Certificate

> Response

```json
{
    "certificate": {
        "domain": "domain.com",
        "request_status": "creating",
        "created_at": "2016-12-17 07:02:35",
        "id": 3,
        "existing": false,
        "active": false
    }
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/sites/{siteId}/certificates/{id}`

## Get Signing Request

### HTTP Request

`GET /api/v1/servers/{serverId}/sites/{siteId}/certificates/{id}/csr`

This endpoint may be used to get the full certificate signing request content.

## Install Certificate

> Payload

```json
{
    "certificate": "certificate content",
    "add_intermediates": false
}
```

### HTTP Request

`POST /api/v1/servers/{serverId}/sites/{siteId}/certificates/{id}/install`

## Activate Certificate

### HTTP Request

`POST /api/v1/servers/{serverId}/sites/{siteId}/certificates/{id}/activate`

## Delete Certificate

### HTTP Request

`DELETE /api/v1/servers/{serverId}/sites/{siteId}/certificates/{id}`

# SSH Keys

## Create Key

> Payload

```json
{
    "name": "test-key",
    "key": "KEY_CONTENT_HERE",
    "username": "forge"
}
```

> Response

```json
{
    "key": {
        "id": 9,
        "name": "test-key",
        "username": "forge",
        "status": "installing",
        "created_at": "2016-12-16 16:31:16"
    }
}
```

### HTTP Request

`POST /api/v1/servers/{serverId}/keys`

## List Keys

> Response

```json
{
    "keys": [
        {
            "id": 9,
            "name": "test-key",
            "username": "forge",
            "status": "installing",
            "created_at": "2016-12-16 16:31:16"
        }
    ]
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/keys`

## Get Key

> Response

```json
{
    "key": {
        "id": 9,
        "name": "test-key",
        "username": "forge",
        "status": "installing",
        "created_at": "2016-12-16 16:31:16"
    }
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/keys/{keyId}`

## Delete Key

### HTTP Request

`DELETE /api/v1/servers/{serverId}/keys/{keyId}`

# Workers

## Create Worker

You may pass `php` as the `php_version` value to use the server's default PHP CLI version.

> Payload

```json
{
    "connection": "sqs",
    "timeout": 90,
    "sleep": 60,
    "tries": null,
    "processes": 1,
    "stopwaitsecs": 600,
    "daemon": true,
    "force": false,
    "php_version": "php72",
    "queue": "hotfix"
}
```

> Response

```json
{
    "worker": {
        "id": 1,
        "connection": "rule",
        "command": "php7.2 /home/forge/default/artisan queue:work rule --sleep=60 --daemon --quiet --timeout=90",
        "queue": null,
        "timeout": 90,
        "sleep": 60,
        "tries": null,
        "processes": 1,
        "stopwaitsecs": 600,
        "environment": null,
        "php_version": "php72",
        "daemon": 1,
        "force": 0,
        "status": "installing",
        "created_at": "2016-12-17 07:15:03"
    }
}
```

### Queue Options

The `queue` key can be left blank for the default queue.

### HTTP Request

`POST /api/v1/servers/{serverId}/sites/{siteId}/workers`

## List Workers

> Response

```json
{
    "workers": [
        {
            "id": 1,
            "connection": "rule",
            "command": "php7.2 /home/forge/default/artisan queue:work rule --sleep=60 --daemon --quiet --timeout=90",
            "queue": null,
            "timeout": 90,
            "sleep": 60,
            "tries": null,
            "processes": 1,
            "stopwaitsecs": null,
            "environment": null,
            "php_version": "php72",
            "daemon": 1,
            "force": 0,
            "status": "installing",
            "created_at": "2016-12-17 07:15:03"
        }
    ]
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/sites/{siteId}/workers`

## Get Worker

> Response

```json
{
    "worker": {
        "id": 1,
        "connection": "rule",
        "command": "php7.2 /home/forge/default/artisan queue:work rule --sleep=60 --daemon --quiet --timeout=90",
        "queue": null,
        "timeout": 90,
        "sleep": 60,
        "tries": null,
        "processes": 1,
        "stopwaitsecs": null,
        "environment": null,
        "php_version": "php72",
        "daemon": 1,
        "force": 0,
        "status": "installing",
        "created_at": "2016-12-17 07:15:03"
    }
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/sites/{siteId}/workers/{id}`

## Delete Worker

### HTTP Request

`DELETE /api/v1/servers/{serverId}/sites/{siteId}/workers/{id}`

## Restart Worker

### HTTP Request

`POST /api/v1/servers/{serverId}/sites/{siteId}/workers/{id}/restart`

## Get Worker Output

> Response

```json
{
    "output": "The output of the worker will be returned here."
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/sites/{siteId/}/workers/{workerId}/output`

# Redirect Rules

## Create Rule

> Payload

```json
{
    "from": "/docs",
    "to": "/docs/1.1",
    "type": "redirect"
}
```

> Response

```json
{
    "redirect_rule": {
        "id": 15,
        "from": "/docs",
        "to": "/docs/1.1",
        "type": "redirect",
        "created_at": "2018-03-07 16:33:20"
    }
}
```

### HTTP Request

`POST /api/v1/servers/{serverId}/sites/{siteId}/redirect-rules`

| Type | Description |
| --- | --- |
| redirect | Creates a temporary 302 redirect |
| permanent | Create a permanent 301 redirect |

## List Redirect Rules

> Response

```json
{
    "redirect_rules": [
        {
            "id": 15,
            "from": "/docs",
            "to": "/docs/1.1",
            "type": "redirect",
            "created_at": "2018-03-07 16:33:20"
        }
    ]
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/sites/{siteId}/redirect-rules`

## Get Rule

> Response

```json
{
    "redirect_rule": {
        "id": 15,
        "from": "/docs",
        "to": "/docs/1.1",
        "type": "redirect",
        "created_at": "2018-03-07 16:33:20"
    }
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/sites/{siteId}/redirect-rules/{id}`

## Delete Rule

### HTTP Request

`DELETE /api/v1/servers/{serverId}/sites/{siteId}/redirect-rules/{id}`

# Security Rules

## Create Security Rule

> Payload

```json
{
    "name": "Access Restricted",
    "path": null,
    "credentials": [
      {
        "username": "taylor.otwell",
        "password": "password123"
      }, {
        "username": "james.brooks",
        "password": "secret123"
      }
    ]
}
```

> Response

```json
{
    "security_rule": {
        "id": 15,
        "name": "Access Restricted",
        "path": null,
        "created_at": "2020-07-30 10:11:10",
        "credentials": [
          {
            "id": 20,
            "username": "taylor.otwell",
            "created_at": "2020-07-30 10:11:10"
          }, {
             "id": 21,
             "username": "james.brooks",
             "created_at": "2020-07-30 10:11:10"
          }
        ]
    }
}
```

### HTTP Request

`POST /api/v1/servers/{serverId}/sites/{siteId}/security-rules`

You may leave `path` empty to protect all routes within your site.

## List Security Rules

> Response

```json
{
    "security_rules": [
        {
             "id": 15,
             "name": "Access Restricted",
             "path": null,
             "created_at": "2020-07-30 10:11:10",
             "credentials": [
               {
                 "id": 20,
                 "username": "taylor.otwell",
                 "created_at": "2020-07-30 10:11:10"
               }, {
                  "id": 21,
                  "username": "james.brooks",
                  "created_at": "2020-07-30 10:11:10"
               }
             ]
         }
    ]
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/sites/{siteId}/security-rules`

## Get Security Rule

> Response

```json
{
    "security_rule": {
        "id": 15,
        "name": "Access Restricted",
        "path": null,
        "created_at": "2020-07-30 10:11:10",
        "credentials": [
          {
            "id": 20,
            "username": "taylor.otwell",
            "created_at": "2020-07-30 10:11:10"
          }, {
             "id": 21,
             "username": "james.brooks",
             "created_at": "2020-07-30 10:11:10"
          }
        ]
    }
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/sites/{siteId}/security-rules/{id}`

## Delete Security Rule

### HTTP Request

`DELETE /api/v1/servers/{serverId}/sites/{siteId}/security-rules/{id}`

# Deployment

## Enable Quick Deployment

### HTTP Request

`POST /api/v1/servers/{serverId}/sites/{siteId}/deployment`

## Disable Quick Deployment

### HTTP Request

`DELETE /api/v1/servers/{serverId}/sites/{siteId}/deployment`

## Get Deployment Script

The response is a string for this request.

### HTTP Request

`GET /api/v1/servers/{serverId}/sites/{siteId}/deployment/script`

## Update Deployment Script

> Payload

```json
{
    "content": "CONTENT_OF_THE_SCRIPT",
    "auto_source": false
}
```

### HTTP Request

`PUT /api/v1/servers/{serverId}/sites/{siteId}/deployment/script`

### Parameters

| Key | Description |
| --- | --- |
| `content` | The contents of the deployment script. This field is required. |
| `auto_source` | Whether to automatically source environment variables into the deployment script. |

## Deploy Now

### HTTP Request

`POST /api/v1/servers/{serverId}/sites/{siteId}/deployment/deploy`

## Reset Deployment Status

### HTTP Request

`POST /api/v1/servers/{serverId}/sites/{siteId}/deployment/reset`

## Get Deployment Log

The response is a string for this request.

### HTTP Request

`GET /api/v1/servers/{serverId}/sites/{siteId}/deployment/log`

# Deployment History

## List Deployments

> Response

```json
{
  "deployments": [
    {
      "id": 71,
      "server_id": 196,
      "site_id": 110,
      "type": 4,
      "commit_hash": "1aa50f0e4c49fed3a2335e866b03d4178ab93c4e",
      "commit_author": "Dries Vints",
      "commit_message": "Merge branch '8.x'\n\n# Conflicts:\n#\tCHANGELOG.md",
      "started_at": "2020-11-05 12:56:05",
      "ended_at": "2020-11-05 12:56:11",
      "status": "failed",
      "displayable_type": "Deployment API"
    }
  ]
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/sites/{siteId}/deployment-history`

## Get Deployment

> Response

```json
{
  "deployment": {
      "id": 71,
      "server_id": 196,
      "site_id": 110,
      "type": 4,
      "commit_hash": "1aa50f0e4c49fed3a2335e866b03d4178ab93c4e",
      "commit_author": "Dries Vints",
      "commit_message": "Merge branch '8.x'\n\n# Conflicts:\n#\tCHANGELOG.md",
      "started_at": "2020-11-05 12:56:05",
      "ended_at": "2020-11-05 12:56:11",
      "status": "failed",
      "displayable_type": "Deployment API"
    }
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/sites/{siteId}/deployment-history/{deploymentId}`

## Get Deployment Output

> Response

```json
{
  "output": "Thu 05 Nov 2020 12:42:30 PM UTC\nFrom github.com:laravel\/laravel\n * branch              master     -> FETCH_HEAD\nAlready up to date.\nInstalling dependencies from lock file (including require-dev)\nVerifying lock file contents can be installed on current platform.\nNothing to install, update or remove\nGenerating optimized autoload files\n> Illuminate\\Foundation\\ComposerScripts::postAutoloadDump\n> @php artisan package:discover --ansi\nDiscovered Package: [32mfideloper\/proxy[39m\nDiscovered Package: [32mfruitcake\/laravel-cors[39m\nDiscovered Package: [32mlaravel\/tinker[39m\nDiscovered Package: [32mnesbot\/carbon[39m\nDiscovered Package: [32mnunomaduro\/collision[39m\n[32mPackage manifest generated successfully.[39m\n73 packages you are using are looking for funding.\nUse the `composer fund` command to find out more!\nReloading PHP FPM...\n\n   Illuminate\\Database\\QueryException \n\n  SQLSTATE[HY000] [1049] Unknown database 'laravel' (SQL: select * from information_schema.tables where table_schema = laravel and table_name = migrations and table_type = 'BASE TABLE')\n\n  at vendor\/laravel\/framework\/src\/Illuminate\/Database\/Connection.php:671\n    667▕         \/\/ If an exception occurs when attempting to run a query, we'll format the error\n    668▕         \/\/ message to include the bindings with SQL, which will make this exception a\n    669▕         \/\/ lot more helpful to the developer instead of just the database's errors.\n    670▕         catch (Exception $e) {\n  ➜ 671▕             throw new QueryException(\n    672▕                 $query, $this->prepareBindings($bindings), $e\n    673▕             );\n    674▕         }\n    675▕\n\n      [2m+33 vendor frames [22m\n  34  artisan:37\n      Illuminate\\Foundation\\Console\\Kernel::handle()\n"
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/sites/{siteId}/deployment-history/{deploymentId}/output`

# Configuration Files

## Get Nginx Configuration

The response is a string for this request.

### HTTP Request

`GET /api/v1/servers/{serverId}/sites/{siteId}/nginx`

## Update Nginx Configuration

> Payload

```json
{
    "content": "CONTENT"
}
```

### HTTP Request

`PUT /api/v1/servers/{serverId}/sites/{siteId}/nginx`

## Get .env File

The response is a string for this request.

### HTTP Request

`GET /api/v1/servers/{serverId}/sites/{siteId}/env`

## Update .env File

> Payload

```json
{
    "content": "CONTENT"
}
```

### HTTP Request

`PUT /api/v1/servers/{serverId}/sites/{siteId}/env`

# Git Projects

## Install New

> Payload

```json
{
    "provider": "github",
    "repository": "username/repository",
    "branch": "master",
    "composer": true
}
```

### HTTP Request

`POST /api/v1/servers/{serverId}/sites/{siteId}/git`

### Parameters

| Key | Description |
| --- | --- |
| provider | The repository provider. Valid values are `github`, `gitlab`, `gitlab-custom`, `bitbucket`, and `custom`. |
| composer | Whether to install Composer dependencies. Valid values are `true` or `false`. |

## Update Repository

> Payload

```json
{
    "provider": "github",
    "repository": "username/repository",
    "branch": "master"
}
```

### HTTP Request

`PUT /api/v1/servers/{serverId}/sites/{siteId}/git`

## Remove Project

### HTTP Request

`DELETE /api/v1/servers/{serverId}/sites/{siteId}/git`

## Create Deploy Key

> Response

```json
{
  "key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDM9Uq4P4rrJCwFHqfvA5vJ6CfvlVZNpGeopmxXRKmN7yjyaMDXEBHSoOwftTsaqNE+Y1M12yctCUyyFpHxVHOhvqiT6XRsYVDSASMYm7rZWQjt\/zXJSdl80NvY\/2m5dSbLWfr9CoHcRROx3Ja213b8Qc9BOQfINbsT4OsGPrOlvpHyCFWDgu4wQvmVcmaGa2soJr92TaGTkJv6T73BjUXD8ZdfYmCkX5y3L2cXgUiNhcDTDm9G+tQebfdr77CRhGcOUi473MsSDuEPCV7RtDHSVA5\/SSSReZyOW3MocObl3LPyq18gmiX9kUO5bVCAev7Yf5QCB2SJFUl5StZ9Wn1yLtY+P02fFZNr+GrmqbAlhv2rTf8UqOBzal46j8oGYbBaRC4BvUKzmxjM7VbUVGO3+8DJIiJYSZoEr+9ptQbs+0YVo1lVah8O1TGm1uoh1LEV36d3GzHbeUjfN71Oqrq5929gt3Ppt\/phxSli7VAQgBIKYvhtWlVxeAz\/EekxMCYSmWT8ZtGOLWlFLQMNEFn5wP\/+CW5VxzQaQbvzsW4EEkBFH5BW0BRt99FQlbhJ7RZmRYD+v1r8Du8Er9I8WGj8f\/cP8PmxlfvVVPokrxvWr4E7GU5mNHIFQTz7hAq4DIxbaR96IEcd6INxa1wfoWgib+9YW77edyX1C1iF2bGGIw== worker@forge.laravel.com"
}
```

### HTTP Request

`POST /api/v1/servers/{serverId}/sites/{siteId}/deploy-key`

## Delete Deploy Key

### HTTP Request

`DELETE /api/v1/servers/{serverId}/sites/{siteId}/deploy-key`

# Site Commands

## Execute Command

> Payload

```json
{
  "command": "ls -la"
}
```

### HTTP Request

`POST /api/v1/servers/{serverId}/sites/{siteId}/commands`

## List Command History

> Response

```json
{
  "commands": [
    {
      "id": 68,
      "server_id": 34,
      "site_id": 48,
      "user_id": 1,
      "event_id": 730,
      "command": "ls -lah",
      "status": "finished",
      "created_at": "2021-04-16 14:46:55",
      "updated_at": "2021-04-16 14:47:00",
      "profile_photo_url": "https:\/\/unavatar.vercel.app\/james%40brooks.page?fallback=https%3A%2F%2Fui-avatars.com%2Fapi%3Fname%3DJames%2BBrooks%26color%3D7F9CF4%26background%3DEBF4FF",
      "user_name": "James Brooks"
    },
    {
      "id": 69,
      "server_id": 34,
      "site_id": 48,
      "user_id": 1,
      "event_id": 731,
      "command": "echo 'Hello!'",
      "status": "finished",
      "created_at": "2021-04-16 14:48:01",
      "updated_at": "2021-04-16 14:48:07",
      "profile_photo_url": "https:\/\/unavatar.vercel.app\/james%40brooks.page?fallback=https%3A%2F%2Fui-avatars.com%2Fapi%3Fname%3DJames%2BBrooks%26color%3D7F9CF4%26background%3DEBF4FF",
      "user_name": "James Brooks"
    }
  ]
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/sites/{siteId}/commands`

## Get Command

> Response

```json
{
  "command": {
      "id": 69,
      "server_id": 34,
      "site_id": 48,
      "user_id": 1,
      "event_id": 731,
      "command": "echo 'Hello!'",
      "status": "finished",
      "created_at": "2021-04-16 14:48:01",
      "updated_at": "2021-04-16 14:48:07",
      "profile_photo_url": "https:\/\/unavatar.vercel.app\/james%40brooks.page?fallback=https%3A%2F%2Fui-avatars.com%2Fapi%3Fname%3DJames%2BBrooks%26color%3D7F9CF4%26background%3DEBF4FF",
      "user_name": "James Brooks"
  },
  "output": "Hello!"
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/sites/{siteId}/commands/{commandId}`

# WordPress

## Install

> Payload

```json
{
    "database": "forge",
    "user": 1
}
```

### HTTP Request

`POST /api/v1/servers/{serverId}/sites/{siteId}/wordpress`

## Uninstall WordPress

### HTTP Request

`DELETE /api/v1/servers/{serverId}/sites/{siteId}/wordpress`

This endpoint will uninstall WordPress and revert the site back to a default state.

# phpMyAdmin

## Install

> Payload

```json
{
    "database": "forge",
    "user": 1
}
```

### HTTP Request

`POST /api/v1/servers/{serverId}/sites/{siteId}/phpmyadmin`

## Uninstall phpMyAdmin

### HTTP Request

`DELETE /api/v1/servers/{serverId}/sites/{siteId}/phpmyadmin`

This endpoint will uninstall phpMyAdmin and revert the site back to a default state.

# Webhooks

## List

> Response

```json
{
    "webhooks": [
        {
            "id": 10,
            "url": "http://domain.com",
            "created_at": "2018-10-10 17:01:18"
        }
    ]
}
```

### HTTP Request

`GET /api/v1/servers/{server_id}/sites/{site_id}/webhooks`

## Show

> Response

```json
{
    "webhook": {
        "id": 10,
        "url": "http://domain.com",
        "created_at": "2018-10-10 17:01:18"
    }
}
```

### HTTP Request

`GET /api/v1/servers/{server_id}/sites/{site_id}/webhooks/{id}`

## Create

> Response

```json
{
    "url": "http://domain.com"
}
```

### HTTP Request

`POST /api/v1/servers/{server_id}/sites/{site_id}/webhooks`

## Delete

> Response

```json
{
    "url": "http://domain.com"
}
```

### HTTP Request

`DELETE /api/v1/servers/{server_id}/sites/{site_id}/webhooks/{id}`

# Deployment Failure Emails

## Set Deployment Failure Email

Set the deployment failure emails for a site.

It will replace the existing deployment failure emails

> Payload

```json
{
    "emails": ["failure@domain.org"]
}
```

### HTTP Request

`POST /api/v1/servers/{server_id}/sites/{site_id}/deployment-failure-emails`

# Recipes

## Create Recipe

> Payload

```json
{
    "name": "Recipe Name",
    "user": "root",
    "script": "SCRIPT_CONTENT"
}
```

> Response

```json
{
    "recipe": {
        "id": 1,
        "name": "Recipe Name",
        "user": "root",
        "script": "SCRIPT_CONTENT",
        "created_at": "2016-12-16 16:24:05"
    }
}
```

### HTTP Request

`POST /api/v1/recipes`

## List Recipes

> Response

```json
{
    "recipes": [
        {
            "id": 1,
            "name": "Recipe Name",
            "user": "root",
            "script": "SCRIPT_CONTENT",
            "created_at": "2016-12-16 16:24:05"
        }
    ]
}
```

### HTTP Request

`GET /api/v1/recipes`

## Get Recipe

> Response

```json
{
    "recipe": {
        "id": 1,
        "name": "Recipe Name",
        "user": "root",
        "script": "SCRIPT_CONTENT",
        "created_at": "2016-12-16 16:24:05"
    }
}
```

### HTTP Request

`GET /api/v1/recipes/{recipeId}`

## Update Recipe

> Payload

```json
{
    "name": "Recipe Name",
    "user": "root",
    "script": "SCRIPT_CONTENT"
}
```

> Response

```json
{
    "recipe": {
        "id": 1,
        "name": "Recipe Name",
        "user": "root",
        "script": "SCRIPT_CONTENT",
        "created_at": "2016-12-16 16:24:05"
    }
}
```

### HTTP Request

`PUT /api/v1/recipes/{recipeId}`

## Delete Recipe

### HTTP Request

`DELETE /api/v1/recipes/{recipeId}`

## Run Recipe

> Payload

```json
{
    "servers": [1,2],
    "notify": true
}
```

### HTTP Request

`POST /api/v1/recipes/{recipeId}/run`

# Regions

## List

> Response

```json
{
  "regions": {
    "ocean2": [
      {
        "id": "ams2",
        "name": "Amsterdam 2",
        "sizes": [
          {
            "id": "01",
            "size": "s-1vcpu-1gb",
            "name": "1GB RAM - 1 CPU Core - 25GB SSD"
          }
        ]
      }
    ],
    "linode": [],
    "vultr": [],
    "aws": []
  }
}
```

### HTTP Request

`GET /api/v1/regions`

# Credentials

## List

> Response

```json
{
  "credentials": [
    {
      "id": 1,
      "type": "ocean2",
      "name": "Personal"
    }
  ]
}
```

### HTTP Request

`GET /api/v1/credentials`

# Backups

## List Backup Configurations

### HTTP Request

`GET /api/v1/servers/{serverId}/backup-configs`

> Response

```json
{
    "backups": [
        {
            "id": 10,
            "day_of_week": null,
            "time": null,
            "provider" : "spaces",
            "provider_name": "DigitalOcean Spaces",
            "status": "installed",
            "databases": [{
                "id": 100,
                "name": "forge",
                "status": "installed",
                "created_at": "2020-01-01 10:00:00"
            }],
            "backups": [
                {
                    "id": 144,
                    "backup_id": 10,
                    "status": "success",
                    "restore_status": null,
                    "archive_path": "s3://backup-configs/server/db/backup-10-20200101123601.tar.gz",
                    "duration": 4,
                    "date": "1st Jan 12:36 PM"
                }
            ],
            "last_backup_time": "3 days ago"
        }
    ]
}
```

## Create Backup Configuration

> Payload

```json
{
    "provider": "spaces",
    "credentials": {
        "endpoint": "https://my-endpoint.com",
        "region": "region-key",
        "bucket": "bucket-name",
        "access_key": "",
        "secret_key": ""
    },
    "frequency": {
        "type": "weekly",
        "time": "12:30",
        "day": 1
    },
    "directory": "backups/server/db",
    "email": "forge@laravel.com",
    "retention": 7,
    "databases": [
        24
    ]
}
```

> Response

```json
{
    "backup": {
        "id": 10,
        "day_of_week": null,
        "time": null,
        "provider" : "spaces",
        "provider_name": "DigitalOcean Spaces",
        "status": "installing",
        "databases": [
            {
                "id": 24,
                "name": "forge",
                "status": "installed",
                "created_at": "2020-01-13 15:47:33"
            }
        ],
        "backups": [],
        "last_backup_time": null
    }
}
```

### HTTP Request

`POST /api/v1/servers/{serverId}/backup-configs`

### Available providers

| Key | Description |
| --- | --- |
| `s3` | Amazon S3 |
| `spaces` | DigitalOcean Spaces |
| `custom` | Custom (S3 Compatible, e.g. MinIO) |

When supplying a `custom` provider, you **must** also provide an `endpoint`.

### Frequency options

*   `hourly`
*   `daily` - you must supply a `time` in 24 hour format
*   `weekly` - you must supply a `time` in 24 hour format and a `day` option 0 (Sunday) - 6 (Saturday)
*   `custom` - you must supply a `custom` value, as a valid cron expression

## Update Backup Configuration

> Payload

```json
{
    "provider": "spaces",
    "credentials": {
        "endpoint": "https://my-endpoint.com",
        "region": "region-key",
        "bucket": "bucket-name",
        "access_key": "",
        "secret_key": ""
    },
    "frequency": {
        "type": "weekly",
        "time": "12:30",
        "day": 1
    },
    "directory": "backups/server/db",
    "email": "forge@laravel.com",
    "retention": 7,
    "databases": [
      24,
      25
    ]
}
```

> Response

```json
{
    "backup": {
        "id": 10,
        "day_of_week": null,
        "time": null,
        "provider" : "spaces",
        "provider_name": "DigitalOcean Spaces",
        "status": "updating",
        "databases": [
            {
                "id": 24,
                "name": "forge",
                "status": "installed",
                "created_at": "2020-01-13 15:47:33"
            }
        ],
        "backups": [],
        "last_backup_time": null
    }
}
```

The payload and options are the same for updating as they are for creating.

### HTTP Request

`PUT /api/v1/servers/{serverId}/backup-configs/{backupConfigurationId}`

## Get Backup Configuration

> Response

```json
{
    "backup": {
        "id": 10,
        "day_of_week": null,
        "time": null,
        "provider" : "spaces",
        "provider_name": "DigitalOcean Spaces",
        "status": "installed",
        "databases": [
            {
                "id": 24,
                "name": "forge",
                "status": "installed",
                "created_at": "2020-01-13 15:47:33"
            }
        ],
        "backups": [],
        "last_backup_time": null
    }
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/backup-configs/{backupConfigurationId}`

## Run Backup Configuration

Manually run a backup configuration.

### HTTP Request

`POST /api/v1/servers/{serverId}/backup-configs/{backupConfigurationId}`

## Delete Backup Configuration

### HTTP Request

`DELETE /api/v1/servers/{serverId}/backup-configs/{backupConfigurationId}`

## Restore Backup

> Payload

```json
{
  "database": 7
}
```

If no `database` value is provided, Forge will restore the first database available.

### HTTP Request

`POST /api/v1/servers/{serverId}/backup-configs/{backupConfigurationId}/backups/{backupId}`

## Delete Backup

### HTTP Request

`DELETE /api/v1/servers/{serverId}/backup-configs/{backupConfigurationId}/backups/{backupId}`

# Monitoring

## List Monitors

### HTTP Request

`GET /api/v1/servers/{serverId}/monitors`

> Response

```json
{
    "monitors": [
        {
            "id": 3,
            "status": "installed",
            "type": "free_memory",
            "operator": "lte",
            "threshold": 70,
            "minutes": 5,
            "state": "ALERT",
            "state_changed_at": "2020-03-01 12:45:00"
        },
        {
            "id": 7,
            "status": "installed",
            "type": "disk",
            "operator": "lte",
            "threshold": 25,
            "minutes": 0,
            "state": "OK",
            "state_changed_at": "2020-03-01 12:45:00"
        }
    ]
}
```

## Create Monitor

> Payload

```json
{
    "type": "cpu_load",
    "operator": "gte",
    "threshold": "1.3",
    "minutes": "5",
    "notify": "forge@laravel.com"
}
```

> Response

```json
{
    "monitor": {
        "id": 8,
        "status": "installed",
        "type": "disk",
        "operator": "lte",
        "threshold": 25,
        "minutes": 0,
        "state": "OK",
        "state_changed_at": "2020-03-01 12:45:00"
    }
}
```

### HTTP Request

`POST /api/v1/servers/{serverId}/monitors`

### Available Monitors

Monitors work on a % value.

| Key | Description |
| --- | --- |
| `disk` | The used disk space |
| `used_memory` | The amount of used memory |
| `cpu_load` | The CPU load |

### Operators

*   `gte` - greater than or equal to
*   `lte` - less than or equal to

## Get Monitor

### HTTP Request

`GET /api/v1/servers/{serverId}/monitors/{monitorId}`

> Response

```json
{
    "monitor": {
        "id": 3,
        "status": "installed",
        "type": "free_memory",
        "operator": "lte",
        "threshold": 70,
        "minutes": 5,
        "state": "ALERT",
        "state_changed_at": "2020-03-01 12:45:00"
    }
}
```

## Delete Monitor

### HTTP Request

`DELETE /api/v1/servers/{serverId}/monitors/{monitorId}`

# Server Logs

## Get Log

### HTTP Request

`GET /api/v1/servers/{serverId}/logs`

> Response

```json
{
      "path": "\/var\/log\/mysql\/error.log",
      "content": "2020-08-18T10:22:01.238990Z 0 [System] [MY-010931] [Server] \/usr\/sbin\/mysqld: ready for connections. Version: '8.0.21'  socket: '\/var\/run\/mysqld\/mysqld.sock'  port: 3306  MySQL Community Server - GPL.\n2020-08-18T10:22:01.359309Z 0 [System] [MY-013172] [Server] Received SHUTDOWN from user <via user signal>. Shutting down mysqld (Version: 8.0.21).\n2020-08-18T10:22:03.644177Z 0 [System] [MY-010910] [Server] \/usr\/sbin\/mysqld: Shutdown complete (mysqld 8.0.21)  MySQL Community Server - GPL.\n2020-08-18T10:22:04.183385Z 0 [System] [MY-010116] [Server] \/usr\/sbin\/mysqld (mysqld 8.0.21) starting as process 42752\n2020-08-18T10:22:04.193962Z 1 [System] [MY-013576] [InnoDB] InnoDB initialization has started.\n2020-08-18T10:22:04.530305Z 1 [System] [MY-013577] [InnoDB] InnoDB initialization has ended.\n2020-08-18T10:22:04.662877Z 0 [System] [MY-011323] [Server] X Plugin ready for connections. Bind-address: '::' port: 33060, socket: \/var\/run\/mysqld\/mysqlx.sock\n2020-08-18T10:22:04.749019Z 0 [Warning] [MY-010068] [Server] CA certificate ca.pem is self signed.\n2020-08-18T10:22:04.749453Z 0 [System] [MY-013602] [Server] Channel mysql_main configured to support TLS. Encrypted connections are now supported for this channel.\n2020-08-18T10:22:04.775994Z 0 [System] [MY-010931] [Server] \/usr\/sbin\/mysqld: ready for connections. Version: '8.0.21'  socket: '\/var\/run\/mysqld\/mysqld.sock'  port: 3306  MySQL Community Server - GPL.\n"
}
```

### File Types

You must specify one of the below `file` types:

*   `nginx_access`
*   `nginx_error`
*   `database`
*   `php7x` (where `x` is a valid version number, e.g. `php71`) or `php56`

# Integrations

## Check Horizon Daemon Status

Checks if a Horizon daemon is enabled for the site.

> Response (enabled)

```json
{
    "enabled": true,
    "daemon": {
        "id": 2,
        "command": "php8.3 artisan horizon",
        "user": "forge",
        "directory": "/home/forge/foo.com",
        "processes": 1,
        "startsecs": 1,
        "stopwaitsecs": 5,
        "stopsignal": "SIGTERM",
        "status": "installed",
        "created_at": "2024-07-04 08:34:26"
    },
    "horizon_installed": true
}
```

> Response (disabled)

```json
{
    "enabled": false,
    "daemon": null,
    "horizon_installed": true
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/sites/{siteId}/integrations/horizon`

## Enable/Create Horizon Daemon

Creates a new Horizon daemon for the site.

If the site had a previously configured daemon, it will be converted to a site managed daemon.

Requires `server:create-daemons` permission.

> Payload

> Response

```json
{
    "daemon": {
        "command": "php8.3 artisan horizon",
        "user": "forge",
        "directory": "/home/forge/foo.com",
        "processes": 1,
        "startsecs": 1,
        "stopwaitsecs": 5,
        "stopsignal": "SIGTERM",
        "status": "installing",
        "created_at": "2024-07-04 08:34:26",
        "id": 2
    }
}
```

### HTTP Request

`POST /api/v1/servers/{serverId}/sites/{siteId}/integrations/horizon`

## Disable/remove Horizon Daemon

Removes the Horizon daemon for the site.

Requires `server:delete-daemons` permission.

### HTTP Request

`DELETE /api/v1/servers/{serverId}/sites/{siteId}/integrations/horizon`

## Check Octane Daemon Status

Checks if a Octane daemon is enabled for the site.

> Response (enabled)

```json
{
    "enabled": true,
    "daemon": {
        "id": 3,
        "command": "php8.3 artisan octane:start --no-interaction",
        "user": "forge",
        "directory": "/home/forge/foo.bar",
        "processes": 1,
        "startsecs": 1,
        "stopwaitsecs": 5,
        "stopsignal": "SIGTERM",
        "status": "installed",
        "created_at": "2024-07-04 08:49:27"
    },
    "octane_installed": true
}
```

> Response (disabled)

```json
{
    "enabled": false,
    "daemon": null,
    "octane_installed": true
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/sites/{siteId}/integrations/octane`

## Enable/Create Octane Daemon

Creates a new Octane daemon for the site.

If the site had a previously configured daemon, it will be converted to a site managed daemon.

Requires `server:create-daemons` permission.

> Payload

```json
{
    "port": 8000,
    "server": "swoole"
}
```

> Response

```json
{
    "daemon": {
        "command": "php8.3 artisan octane:start --no-interaction",
        "user": "forge",
        "directory": "/home/forge/foo.bar",
        "processes": 1,
        "startsecs": 1,
        "stopwaitsecs": 5,
        "stopsignal": "SIGTERM",
        "status": "installing",
        "created_at": "2024-07-04 08:49:27",
        "id": 3
    }
}
```

### Parameters

| Key | Description |
| --- | --- |
| `port` | This field is optional (default: `8000`). The port for which Octane should run on. |
| `server` | This field is optional (default: `swoole`). The server type to use with Octane (`swoole`/`roadrunner`/`frankenphp`) |

### HTTP Request

`POST /api/v1/servers/{serverId}/sites/{siteId}/integrations/octane`

## Disable/remove Octane Daemon

Removes the Octane daemon for the site.

Requires `server:delete-daemons` permission.

### HTTP Request

`DELETE /api/v1/servers/{serverId}/sites/{siteId}/integrations/octane`

## Check Reverb Daemon Status

Checks if a Reverb daemon is enabled for the site.

> Response (enabled)

```json
{
    "enabled": true,
    "daemon": {
        "id": 6,
        "command": "php8.3 artisan reverb:start --no-interaction --port=8080",
        "user": "forge",
        "directory": "/home/forge/bar.test",
        "processes": 1,
        "startsecs": 1,
        "stopwaitsecs": 5,
        "stopsignal": "SIGTERM",
        "status": "installed",
        "created_at": "2024-07-05 12:36:07"
    },
    "reverb_host": "ws.anotther.test",
    "reverb_port": 8080,
    "reverb_connections": 1000,
    "reverb_installed": true
}
```

> Response (disabled)

```json
{
    "enabled": false,
    "daemon": null,
    "reverb_host": null,
    "reverb_port": null,
    "reverb_connections": null,
    "reverb_installed": true
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/sites/{siteId}/integrations/reverb`

## Enable/Create Reverb Daemon

Creates a new Reverb daemon for the site.

If the site had a previously configured daemon, it will be converted to a site managed daemon.

Requires `server:create-daemons` permission.

> Payload

```json
{
    "port": 8080,
    "host": "ws.bar.test",
    "connections": 1000
}
```

> Response

```json
{
    "daemon": {
        "command": "php8.3 artisan reverb:start --no-interaction --port=8080",
        "user": "forge",
        "directory": "/home/forge/foo.test",
        "processes": 1,
        "startsecs": 1,
        "stopwaitsecs": 5,
        "stopsignal": "SIGTERM",
        "status": "installing",
        "created_at": "2024-07-05 12:36:07",
        "id": 6
    },
    "reverb_host": "ws.bar.test",
    "reverb_port": 8080,
    "reverb_connections": 1000
}
```

### Parameters

| Key | Description |
| --- | --- |
| `port` | This field is optional (default: `8080`). The port for which Reverb should run on. |
| `host` | This field is optional (default: `ws.bar.test`). The host for Reverb. |
| `connections` | The field is optional (default: `1000`). The amount of connections allowed. |

### HTTP Request

`POST /api/v1/servers/{serverId}/sites/{siteId}/integrations/reverb`

## Disable/remove Reverb Daemon

Removes the Reverb daemon for the site.

Requires `server:delete-daemons` permission.

### HTTP Request

`DELETE /api/v1/servers/{serverId}/sites/{siteId}/integrations/reverb`

## Check Pulse Daemon Status

Checks if a Pulse daemon is enabled for the site.

> Response (enabled)

```json
{
    "enabled": true,
    "daemon": {
        "id": 6,
        "command": "php8.3 artisan pulse:check",
        "user": "forge",
        "directory": "/home/forge/bar.test",
        "processes": 1,
        "startsecs": 1,
        "stopwaitsecs": 5,
        "stopsignal": "SIGTERM",
        "status": "installed",
        "created_at": "2024-07-05 12:36:07"
    },
    "pulse_installed": true
}
```

> Response (disabled)

```json
{
    "enabled": false,
    "daemon": null,
    "pulse_installed": true
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/sites/{siteId}/integrations/pulse`

## Enable/Create Pulse Daemon

Creates a new Pulse daemon for the site.

If the site had a previously configured daemon, it will be converted to a site managed daemon.

Requires `server:create-daemons` permission.

> Payload

> Response

```json
{
    "daemon": {
        "command": "php8.3 artisan pulse:check",
        "user": "forge",
        "directory": "/home/forge/foo.com",
        "processes": 1,
        "startsecs": 1,
        "stopwaitsecs": 5,
        "stopsignal": "SIGTERM",
        "status": "installing",
        "created_at": "2024-07-04 08:34:26",
        "id": 2
    }
}
```

### HTTP Request

`POST /api/v1/servers/{serverId}/sites/{siteId}/integrations/pulse`

## Disable/remove Pulse Daemon

Removes the Pulse daemon for the site.

Requires `server:delete-daemons` permission.

### HTTP Request

`DELETE /api/v1/servers/{serverId}/sites/{siteId}/integrations/pulse`

## Check Inertia Daemon Status

Checks if a Inertia daemon is enabled for the site.

> Response (enabled)

```json
{
    "enabled": true,
    "daemon": {
        "id": 6,
        "command": "php8.3 artisan inertia:start-ssr",
        "user": "forge",
        "directory": "/home/forge/bar.test",
        "processes": 1,
        "startsecs": 1,
        "stopwaitsecs": 5,
        "stopsignal": "SIGQUIT",
        "status": "installed",
        "created_at": "2024-07-05 12:36:07"
    },
    "inertia_installed": true
}
```

> Response (disabled)

```json
{
    "enabled": false,
    "daemon": null,
    "inertia_installed": true
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/sites/{siteId}/integrations/inertia`

## Enable/Create Inertia Daemon

Creates a new Inertia daemon for the site.

If the site had a previously configured daemon, it will be converted to a site managed daemon.

Requires `server:create-daemons` permission.

> Payload

```json
{
    "deploys_restart_inertia_daemon": true
}
```

> Response

```json
{
    "daemon": {
        "command": "php8.3 artisan inertia:start-ssr",
        "user": "forge",
        "directory": "/home/forge/foo.com",
        "processes": 1,
        "startsecs": 1,
        "stopwaitsecs": 5,
        "stopsignal": "SIGQUIT",
        "status": "installing",
        "created_at": "2024-07-04 08:34:26",
        "id": 21
    }
}
```

### Parameters

| Key | Description |
| --- | --- |
| `deploys_restart_inertia_daemon` | This field is optional (default: `false`). Updates the project's deployment script to restart the Inertia SSR daemon during deployments. |

### HTTP Request

`POST /api/v1/servers/{serverId}/sites/{siteId}/integrations/inertia`

## Check Laravel Maintenance Status

Checks if Laravel maintenance mode is enabled or disabled for the site.

> Response (enabled)

```json
{
    "enabled": true,
    "status": null,
    "laravel_installed": true
}
```

> Response (disabled)

```json
{
    "enabled": false,
    "status": null,
    "laravel_installed": true
}
```

> Response (enabling)

```json
{
    "enabled": false,
    "status": "enabling",
    "laravel_installed": true
}
```

> Response (disabling)

```json
{
    "enabled": true,
    "status": "disabling",
    "laravel_installed": true
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/sites/{siteId}/integrations/laravel-maintenance`

## Enable Laravel Maintenance Mode

Enables Laravel maintenance mode for the site.

Requires `site:manage-deploys` permission.

> Payload

```json
{
    "secret": "my-super-secret",
    "status": 503
}
```

> Response

```json
{
    "enabled": false,
    "status": "enabling"
}
```

### Parameters

| Key | Description |
| --- | --- |
| `secret` | This field is optional. The secret key to bypass maintenance mode. |
| `status` | This field is optional (default: `503`). The status code for the maintenance page. |

### HTTP Request

`POST /api/v1/servers/{serverId}/sites/{siteId}/integrations/laravel-maintenance`

## Disable Laravel Maintenance Mode

Disables Laravel maintenance mode for the site.

Requires `site:manage-deploys` permission.

> Response

```json
{
    "enabled": true,
    "status": "disabling"
}
```

### HTTP Request

`DELETE /api/v1/servers/{serverId}/sites/{siteId}/integrations/laravel-maintenance`

## Check Laravel Scheduler Status

Checks if Laravel Scheduler is enabled or disabled for the site.

> Response (enabled)

```json
{
    "enabled": true,
    "job": {
        "id": 1,
        "command": "/home/forge/cronin.com/artisan schedule:run",
        "user": "root",
        "frequency": "",
        "cron": "0 * * * *",
        "status": "",
        "created_at": "2024-07-05 14:44:46",
        "next_run_time": "2024-07-05T15:00:00+00:00"
    },
    "laravel_installed": true
}
```

> Response (disabled)

```json
{
    "enabled": false,
    "job": null,
    "laravel_installed": true
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/sites/{siteId}/integrations/laravel-scheduler`

## Enable/Create Laravel Scheduler

Creates a Laravel Scheduler for the site

Requires `server:create-schedulers` permission.

> Payload

> Response

```json
{
    "job": {
        "command": "php8.3 /home/forge/forging.site/artisan schedule:run",
        "user": "forge",
        "frequency": "Minutely",
        "cron": "* * * * *",
        "status": "installing",
        "created_at": "2024-07-08 08:55:39",
        "id": 7,
        "next_run_time": "2024-07-08T08:56:00+00:00"
    }
}
```

### HTTP Request

`POST /api/v1/servers/{serverId}/sites/{siteId}/integrations/laravel-scheduler`

## Disable/remove Laravel Scheduler

Removes the Laravel scheduler for the site.

Requires `server:delete-schedulers` permission.

### HTTP Request

`DELETE /api/v1/servers/{serverId}/sites/{siteId}/integrations/laravel-scheduler`

# Composer Packages Authentication

## Get Composer Packages Authentication

Get the Composer packages authentication for the site.

> Response

```json
{
    "credentials": {
        "http-basic": {
            "packages.laravel.com": {
                "username": "oliver",
                "password": "secret"
            }
        }
    }
}
```

### HTTP Request

`GET /api/v1/servers/{serverId}/sites/{siteId}/packages`

## Update Composer Packages Authentication

Update the Composer packages authentication for the site.

Requires `server:manage-packages` permission.

> Payload

```json
{
    "credentials": [
        {
            "repository_url": "packages.laravel.com",
            "username": "oliver",
            "password": "secret"
        }
    ]
}
```

> Response

### Parameters

| Key | Description |
| --- | --- |
| `credentials.*.repostiory_url` | The repository URL for composer to look at. |
| `credentials.*.username` | The username for authenticating the repository. |
| `credentials.*.password` | The password for authenticating the repository. |

### HTTP Request

`PUT /api/v1/servers/{serverId}/sites/{siteId}/packages`
