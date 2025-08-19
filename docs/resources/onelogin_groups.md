---
layout: "onelogin"
page_title: "OneLogin: onelogin_groups"
sidebar_current: "docs-onelogin-resource-groups"
description: |-
  Provides a OneLogin group resource.
---

# onelogin_groups

Provides a OneLogin group resource.

> **Note:** The OneLogin API currently only supports reading groups. Creating, updating, and deleting groups through the API is not supported at this time. This resource is provided for future compatibility when these operations become available.

## Example Usage

```hcl
# This resource is read-only for now
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
