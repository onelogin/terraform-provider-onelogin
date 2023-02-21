---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "onelogin_brands_apps Data Source - terraform-provider-onelogin-1"
subcategory: ""
description: |-
  
---

# onelogin_brands_apps (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `brands_id` (String)

### Optional

- `auth_method` (Number)
- `auth_method_description` (String)
- `connector_id` (Number)
- `created_at` (String)
- `description` (String)
- `filter` (Block Set) (see [below for nested schema](#nestedblock--filter))
- `name` (String)
- `updated_at` (String)
- `visible` (Boolean)

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--filter"></a>
### Nested Schema for `filter`

Required:

- `name` (String)
- `values` (List of String)

