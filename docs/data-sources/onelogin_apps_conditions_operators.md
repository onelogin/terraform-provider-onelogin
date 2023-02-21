---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "onelogin_apps_conditions_operators Data Source - terraform-provider-onelogin-1"
subcategory: ""
description: |-
  
---

# onelogin_apps_conditions_operators (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `apps_id` (String)
- `conditions_id` (String)

### Optional

- `filter` (Block Set) (see [below for nested schema](#nestedblock--filter))
- `name` (String) Name of the operator
- `value` (String) The condition operator value to use when creating or updating rules.

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--filter"></a>
### Nested Schema for `filter`

Required:

- `name` (String)
- `values` (List of String)

