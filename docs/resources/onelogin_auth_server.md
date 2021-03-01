---
layout: "onelogin"
page_title: "OneLogin: onelogin_auth_server"
sidebar_current: "docs-onelogin-resource-auth_server"
description: |-
  Creates an Authentication Server Resource.
---

# onelogin_auth_server

Creates an Authentication Server Resource.

This resource allows you to create and configure an Authentication Server.

## Example Usage

```hcl
resource onelogin_auth_servers example {
  name = "Contacts API"
  description = "This is an api"
  configuration = {
    resource_identifier = "https://example.com/contacts"
    audiences = ["https://example.com/contacts"]
    refresh_token_expiration_minutes = 30
    access_token_expiration_minutes = 10
  }
}
```

## Argument Reference

The following arguments are supported:
* `name` - (Required) The resource's name.

* `description` - (Required) A brief description about the resource.

* `configuration` - (Required) Configuration parameters
  * `resource_identifier` - (Required) Unique identifier for the API that the Authorization Server will issue Access Tokens for.

  * `audiences` - (Required) List of API endpoints that will be returned in Access Tokens.

  * `access_token_expiration_minutes` (Optional) The number of minutes until the token expires

  * `refresh_token_expiration_minutes` (Optional) The number of minutes until the token expires


## Attributes Reference

No further attributes are exported

## Import

An Auth Server can be imported via the OneLogin Auth Server ID.

```
$ terraform import onelogin_auth_servers.example <auth_server_id>
```
