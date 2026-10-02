---
layout: "onelogin"
page_title: "OneLogin: onelogin_app_role_attachments"
sidebar_current: "docs-onelogin-resource-app_role_attachments"
description: |-
  Manage App Role Attachment resources.
---

# onelogin_app_role_attachments

Manage App Role Attachment resources.

This resource allows you to create and configure App Roles Attachments. The App Role Attachment is not an API-managed resource and is only a means to attach roles to apps in Terraform to avoid complications with circular dependencies.

Each attachment reads the app's roles, adds or removes its own, and writes the whole list
back. Attachments to the same app are applied one at a time within a run, so `for_each` or
`count` over many roles is safe. To assign several roles to an app in a single write, set
`role_ids` on the app resource instead
([`onelogin_apps`](onelogin_apps.md), [`onelogin_saml_apps`](onelogin_saml_apps.md),
[`onelogin_oidc_apps`](onelogin_oidc_apps.md)). Do not use both on the same app: once
`role_ids` is set it is authoritative, so the next apply removes the roles attached here,
and every plan after that fails with `App <id> does not have role <id>`.

## Example Usage

```hcl
resource onelogin_app_role_attachments example {
	app_id = onelogin_saml_apps.saml.id
	role_id = 12345
}
```

## Argument Reference

The following arguments are supported:

* `app_id` - (Required) The id of the App resource to which the role should belong.

* `role_id` - (Required) The id of the Role being attached to the App.

## Attributes Reference

No further attributes are exported.

## Import

An App Role Attachment cannot be imported at this time.
