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
provider "onelogin" {
  client_id = <your client id>
  client_secret = <your client secret>
  url = <the api url for your region>
}

# Add an App to your account
resource "onelogin_saml_app" "my_saml_app" {
  # ...
}
```

You're also welcome to leave the provider field blank and export your
credentials to your environment

## Argument Reference

The following arguments are supported in the `provider` block:

* `client_id` - (Required) This is the client_id for your OneLogin account that is used to authenticate requests to the OneLogin APIs on your behalf. You can create this by visiting your OneLogin account and selecting from the top ribbon Developers > API Credentials

* `client_secret` - (Required) This is the client_secret for your OneLogin account that is used to authenticate requests to the OneLogin APIs on your behalf. You can create this by visiting your OneLogin account and selecting from the top ribbon Developers > API Credentials

* `url` - (Optional, if no region given) This is the url for your API endpoint depending on your location. It can be api.<us or eu>.onelogin.com

* `region` - (Optional, if no url given) This is the region for your API endpoint. It will be interpolated into the url as shown above.
