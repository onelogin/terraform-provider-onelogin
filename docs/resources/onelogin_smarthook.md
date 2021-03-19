---
layout: "onelogin"
page_title: "OneLogin: onelogin_smarthooks"
sidebar_current: "docs-onelogin-resource-smarthook"
description: |-
  Manage SmartHook resources.
---

# onelogin_smarthooks

Manage SmartHook resources.

This resource allows you to create and configure SmartHooks.

## Example Usage

```hcl
resource onelogin_smarthooks basic_test {
  type = "pre-authentication"
  packages = {
    mysql = "^2.18.1"
  }
  env_vars = [ "API_KEY" ]
  retries = 0
  timeout = 2
  disabled = false
  options = {
    risk_enabled = false
    location_enabled = false
  }
  function = <<EOF
    exports.handler = async context => {
      console.log("Pre-auth executing for " + context.user.user_identifier);
      return { user: context.user };
    };
	EOF
}

resource onelogin_smarthooks basic_test {
  type = "pre-authentication"
  packages = {
    mysql = "^2.18.1"
  }
  env_vars = [ "API_KEY" ]
  retries = 0
  timeout = 2
  disabled = false
  options = {
    risk_enabled = false
    location_enabled = false
  }
  function = "CQlmdW5jdGlvbiBteUZ1bmMoKSB7CgkJCWxldCBhID0gMTsKCQkJbGV0IGIgPSAxOwoJCQlsZXQgYyA9IGEgKyBiOwoJCSAgY29uc29sZS5sb2coIkRpbmcgRG9uZyIsIGEsIGIsIGMpOwoJCX0K"
}

```

## Argument Reference

The following arguments are supported:
* `type` - (Required) The name of the hook. Must be one of: `user-migration` `pre-authentication` `pre-user-create` `post-user-create` `pre-user-update` `post-user-update`

* `status` - (Computed) The smarthook's status.

* `packages` - (Required) A list of public npm packages than will be installed as part of the function build process. These packages names must be on our allowlist. See Node Modules section of this doc. Packages can be any version and support the semantic versioning syntax used by NPM.

* `function` - (Required) A base64 encoded blob, or Heredoc string containing the javascript function code.

* `disabled` - (Required) Indicates if function is available for execution or not. Default true

* `options` - (Required if type = pre-authentication) A list of options for the hook
  * `risk_enabled` - (Required) When true a risk score and risk reasons will be passed in the context. Only applies authentication time hooks. E.g. pre-authentication, user-migration. Default false

  * `location_enabled` - (Required) When true an ip to location lookup is done and the location info is passed in the context. Only applies authentication time hooks. E.g. pre-authentication, user-migration. Default false

* `retries` - (Required) Number of retries if execution fails. Default 0, Max 4

* `timeout` - (Required) The number of milliseconds to allow before timeout. Default 1000, Max 10000

* `env_vars` - (Required) An array of predefined environment variables to be supplied to the function at runtime.

* `created_at` - (Computed) Timestamp for smarthook's last update

* `updated_at` - (Computed) Timestamp for smarthook's last update

## Attributes Reference

No further attributes are exported

## Import

A SmartHook can be imported via the OneLogin SmartHook.

```
$ terraform import onelogin_smarthooks.example <smarthook_id>
```
