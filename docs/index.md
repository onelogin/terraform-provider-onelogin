---
layout: "onelogin"
page_title: "Provider: Onelogin"
description: |-
  The OneLogin provider is used to interact with OneLogin resources.
---

# OneLogin Provider

The OneLogin provider is used to interact with OneLogin resources.

The provider allows you to manage your OneLogin organization's resources easily.
It needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the OneLogin Provider
provider "onelogin" {}

# Add an App to your account
resource "onelogin_saml_app" "my_saml_app" {
  # ...
}
```

You're also welcome to leave the provider field blank and export your
credentials to your environment

## Argument Reference

The following arguments are supported in the `provider` block:

None: This provider reads API credentials from your environment. You need to export
your OneLogin credentials like so:

```
export ONELOGIN_CLIENT_ID=<your client id>
export ONELOGIN_CLIENT_SECRET=<your client secret>
export ONELOGIN_API_URL=<the complete api url, e.g., https://company.onelogin.com>
```

For backward compatibility, you can also use the subdomain approach, but this is deprecated:

```
export ONELOGIN_CLIENT_ID=<your client id>
export ONELOGIN_CLIENT_SECRET=<your client secret>
export ONELOGIN_SUBDOMAIN=<your onelogin subdomain, e.g., 'company' for company.onelogin.com>
```

Using the `ONELOGIN_API_URL` is now the recommended approach.
