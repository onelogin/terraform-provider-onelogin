---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "onelogin Provider"
subcategory: ""
description: |-
  
---

# onelogin Provider





<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `authorization` (String)
- `bearer_auth` (String)
- `content_type` (String)
- `endpoints` (Block Set) (see [below for nested schema](#nestedblock--endpoints))

<a id="nestedblock--endpoints"></a>
### Nested Schema for `endpoints`

Optional:

- `apps` (String) Use this to override the resource endpoint URL (the default one or the one constructed from the `region`).
- `rules` (String) Use this to override the resource endpoint URL (the default one or the one constructed from the `region`).
- `users` (String) Use this to override the resource endpoint URL (the default one or the one constructed from the `region`).