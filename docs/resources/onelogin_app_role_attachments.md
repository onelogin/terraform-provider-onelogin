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
