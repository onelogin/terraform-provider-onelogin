---
layout: "onelogin"
page_title: "OneLogin: onelogin_groups"
sidebar_current: "docs-onelogin-resource-groups"
description: |-
  Manages a OneLogin group.
---

# onelogin_groups

Manages a OneLogin group. Supports full CRUD operations via the V2 API.

## Example Usage

```hcl
resource "onelogin_groups" "engineering" {
  name      = "Engineering"
  reference = "eng-group"
}
```

## Argument Reference

* `name` - (Required) The name of the group.
* `reference` - (Optional) A reference identifier for the group.

## Attribute Reference

* `id` - The ID of the group.

## Import

Groups can be imported using the group ID:

```
$ terraform import onelogin_groups.engineering 123456
```
