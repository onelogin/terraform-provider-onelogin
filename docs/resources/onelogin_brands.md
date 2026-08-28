---
layout: "onelogin"
page_title: "OneLogin: onelogin_brands"
sidebar_current: "docs-onelogin-resource-brands"
description: |-
  Create and configure OneLogin account brands.
---

# onelogin_brands

Manages a OneLogin account brand: the look of the login page and the portal for the users
and apps it applies to — colours, logo, background, the login instruction screen, and the
label on the Username/Email field.

Attach a brand to an application through that application's `brand_id`.

Every account has one **master** brand, created with the account. It cannot be created or
destroyed through this resource, and no brand can be promoted to master. To manage the
master brand's colours and images, `terraform import` it.

## Example Usage

```hcl
resource "onelogin_brands" "engineering" {
  name    = "Engineering"
  enabled = true

  custom_color         = "#123456"
  custom_accent_color  = "#016B91"
  custom_masking_color = "#000000"

  hide_onelogin_footer = true

  enable_custom_label_for_login_screen = true
  custom_label_text_for_login_screen = jsonencode({
    en = "Username"
    de = "Benutzername"
  })

  login_instruction_title = jsonencode({ en = "Need help signing in?" })
  login_instruction       = jsonencode({ en = "Use your **corporate** account." })

  logo = filebase64("${path.module}/logo.png")
}

resource "onelogin_saml_apps" "internal" {
  name         = "Internal"
  connector_id = 108232
  brand_id     = onelogin_brands.engineering.id
}
```

## Argument Reference

* `name` - (Required) Name of the brand. The API rejects a brand without one.
* `enabled` - (Optional) Whether the brand is enabled. Defaults to `false` on create.
* `custom_support_enabled` - (Optional) Whether the login page offers a support request link.
* `custom_color` - (Optional) Primary brand colour, as a hex value such as `#000000`.
* `custom_accent_color` - (Optional) Secondary brand colour, as a hex value such as `#016B91`.
* `custom_masking_color` - (Optional) Colour of the masking layer drawn over the background image.
* `custom_masking_opacity` - (Optional) Opacity of `custom_masking_color`. Defaults to `0`.
* `enable_custom_label_for_login_screen` - (Optional) Whether the Username/Email field uses a custom label.
* `custom_label_text_for_login_screen` - (Optional) Label for the Username/Email field. Localised, see below.
* `login_instruction_title` - (Optional) Link text that opens the login instruction screen. Localised, see below.
* `login_instruction` - (Optional) Text of the login instruction screen, styled in Markdown. Localised, see below.
* `hide_onelogin_footer` - (Optional) Whether to hide the OneLogin footer at the bottom of the login page.
* `mfa_enrollment_message` - (Optional) Text replacing the default message on the first screen of MFA registration.
* `logo` - (Optional) Base64-encoded PNG for the login page logo, under 1MB. Write-only, see below.
* `background` - (Optional) Base64-encoded JPG or PNG for the login page background, under 5MB. Write-only, see below.

## Attribute Reference

* `id` - The brand's ID.
* `master` - Whether this is the account's master brand. Read-only.

## Localised text

`custom_label_text_for_login_screen`, `login_instruction_title` and `login_instruction`
hold one value per locale. OneLogin stores them as a JSON object encoded into a string, so
build them with `jsonencode`:

```hcl
login_instruction = jsonencode({
  en = "Use your corporate account."
  de = "Verwenden Sie Ihr Firmenkonto."
})
```

The API re-encodes what it stores, so the string it returns is rarely byte-identical to the
one sent. These three attributes therefore compare their **decoded** values: a difference in
whitespace or key order is not a change, while a different locale or a different text is.

## Images are write-only

`logo` and `background` are sent as base64 but read back as an object of image URLs and
metadata, not as the data that was uploaded. There is nothing for Terraform to compare, so
neither attribute is refreshed from the API.

The consequence: **an image replaced outside Terraform is not reported as drift.** State
keeps the base64 that was last applied. Changing the configured value still replaces the
image as usual, and `filebase64` on a changed file produces a changed value.

Both are also marked sensitive, so a plan shows `(sensitive value)` rather than the data.
That is not a claim that an image is a secret — it stops Terraform printing the base64 in
full. A 200KB logo otherwise renders as a single 270,909-character line, turning a 2.2KB
plan into a 273KB one, and the API accepts a background five times that size.

## Attributes this resource does not manage

The API returns some brand fields that the Go SDK's `models.Brand` cannot send, so they are
deliberately absent rather than offered and silently ignored: `custom_links`,
`use_custom_smtp_setting`, and — on the master brand — `domain_name`,
`allowed_redirect_urls`, `hide_login_forgot_password` and the `password_self_service_*`
fields. Manage those in the OneLogin admin UI until the SDK model carries them.

## Removing an argument

Every optional argument is also computed, because the API fills in defaults on create.
Deleting an argument from the configuration therefore does not reset it: state keeps the
last value the API returned, and nothing is sent. Set it back to the value you want
explicitly instead.

## Import

```sh
terraform import onelogin_brands.engineering 7844
```
