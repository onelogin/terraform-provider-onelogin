---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "onelogin_auth_servers Data Source - terraform-provider-onelogin"
subcategory: ""
description: |-
  
---

# onelogin_auth_servers (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `configuration` (Block List, Max: 1) Authorization server configuration (see [below for nested schema](#nestedblock--configuration))
- `description` (String) Description of what the API does.
- `filter` (Block Set) (see [below for nested schema](#nestedblock--filter))
- `name` (String) Name of the API.

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--configuration"></a>
### Nested Schema for `configuration`

Optional:

- `access_token_expiration_minutes` (Number) The number of minutes until access token expires. There is no maximum expiry limit.
- `audiences` (List of String) List of API endpoints that will be returned in Access Tokens.
- `refresh_token_expiration_minutes` (Number) The number of minutes until refresh token expires. There is no maximum expiry limit.
- `resource_identifier` (String) Unique identifier for the API that the Authorization Server will issue Access Tokens for.


<a id="nestedblock--filter"></a>
### Nested Schema for `filter`

Required:

- `name` (String)
- `values` (List of String)

