# Terraform Provider Onelogin

This guide lists the configuration for 'onelogin' Terraform provider resources that can be managed using [Terraform v0.12](https://www.hashicorp.com/blog/announcing-terraform-0-12/).

- [Terraform Provider Onelogin](#terraform-provider-onelogin)
  - [Provider Installation](#provider-installation)
  - [Provider Configuration](#provider-configuration)
    - [Example Usage](#example-usage)
  - [Provider Resources](#provider-resources)
    - [onelogin\_apps](#onelogin_apps)
      - [Example usage](#example-usage-1)
      - [Arguments Reference](#arguments-reference)
      - [Attributes Reference](#attributes-reference)
      - [Import](#import)
    - [onelogin\_rules](#onelogin_rules)
      - [Example usage](#example-usage-2)
      - [Arguments Reference](#arguments-reference-1)
      - [Attributes Reference](#attributes-reference-1)
      - [Import](#import-1)
    - [onelogin\_users](#onelogin_users)
      - [Example usage](#example-usage-3)
      - [Arguments Reference](#arguments-reference-2)
      - [Attributes Reference](#attributes-reference-2)
      - [Import](#import-2)
  - [Data Sources (using resource id)](#data-sources-using-resource-id)
    - [onelogin\_apps\_instance](#onelogin_apps_instance)
      - [Example usage](#example-usage-4)
      - [Arguments Reference](#arguments-reference-3)
      - [Attributes Reference](#attributes-reference-3)
    - [onelogin\_rules\_instance](#onelogin_rules_instance)
      - [Example usage](#example-usage-5)
      - [Arguments Reference](#arguments-reference-4)
      - [Attributes Reference](#attributes-reference-4)
    - [onelogin\_users\_instance](#onelogin_users_instance)
      - [Example usage](#example-usage-6)
      - [Arguments Reference](#arguments-reference-5)
      - [Attributes Reference](#attributes-reference-5)
  - [Data Sources (using filters)](#data-sources-using-filters)
    - [onelogin\_api\_authorizations (filters)](#onelogin_api_authorizations-filters)
      - [Example usage](#example-usage-7)
      - [Arguments Reference](#arguments-reference-6)
      - [Attributes Reference](#attributes-reference-6)
    - [onelogin\_api\_authorizations\_claims (filters)](#onelogin_api_authorizations_claims-filters)
      - [Example usage](#example-usage-8)
      - [Arguments Reference](#arguments-reference-7)
      - [Attributes Reference](#attributes-reference-7)
    - [onelogin\_api\_authorizations\_scopes (filters)](#onelogin_api_authorizations_scopes-filters)
      - [Example usage](#example-usage-9)
      - [Arguments Reference](#arguments-reference-8)
      - [Attributes Reference](#attributes-reference-8)
    - [onelogin\_apps\_actions (filters)](#onelogin_apps_actions-filters)
      - [Example usage](#example-usage-10)
      - [Arguments Reference](#arguments-reference-9)
      - [Attributes Reference](#attributes-reference-9)
    - [onelogin\_apps\_actions\_values (filters)](#onelogin_apps_actions_values-filters)
      - [Example usage](#example-usage-11)
      - [Arguments Reference](#arguments-reference-10)
      - [Attributes Reference](#attributes-reference-10)
    - [onelogin\_apps\_conditions (filters)](#onelogin_apps_conditions-filters)
      - [Example usage](#example-usage-12)
      - [Arguments Reference](#arguments-reference-11)
      - [Attributes Reference](#attributes-reference-11)
    - [onelogin\_apps\_conditions\_operators (filters)](#onelogin_apps_conditions_operators-filters)
      - [Example usage](#example-usage-13)
      - [Arguments Reference](#arguments-reference-12)
      - [Attributes Reference](#attributes-reference-12)
    - [onelogin\_apps\_conditions\_values (filters)](#onelogin_apps_conditions_values-filters)
      - [Example usage](#example-usage-14)
      - [Arguments Reference](#arguments-reference-13)
      - [Attributes Reference](#attributes-reference-13)
    - [onelogin\_apps\_rules (filters)](#onelogin_apps_rules-filters)
      - [Example usage](#example-usage-15)
      - [Arguments Reference](#arguments-reference-14)
      - [Attributes Reference](#attributes-reference-14)
    - [onelogin\_apps\_users (filters)](#onelogin_apps_users-filters)
      - [Example usage](#example-usage-16)
      - [Arguments Reference](#arguments-reference-15)
      - [Attributes Reference](#attributes-reference-15)
    - [onelogin\_mappings (filters)](#onelogin_mappings-filters)
      - [Example usage](#example-usage-17)
      - [Arguments Reference](#arguments-reference-16)
      - [Attributes Reference](#attributes-reference-16)
    - [onelogin\_privileges (filters)](#onelogin_privileges-filters)
      - [Example usage](#example-usage-18)
      - [Arguments Reference](#arguments-reference-17)
      - [Attributes Reference](#attributes-reference-17)
    - [onelogin\_rules (filters)](#onelogin_rules-filters)
      - [Example usage](#example-usage-19)
      - [Arguments Reference](#arguments-reference-18)
      - [Attributes Reference](#attributes-reference-18)
    - [onelogin\_users (filters)](#onelogin_users-filters)
      - [Example usage](#example-usage-20)
      - [Arguments Reference](#arguments-reference-19)
      - [Attributes Reference](#attributes-reference-19)
    - [onelogin\_users\_apps (filters)](#onelogin_users_apps-filters)
      - [Example usage](#example-usage-21)
      - [Arguments Reference](#arguments-reference-20)
      - [Attributes Reference](#attributes-reference-20)
    - [onelogin\_users\_devices (filters)](#onelogin_users_devices-filters)
      - [Example usage](#example-usage-22)
      - [Arguments Reference](#arguments-reference-21)
      - [Attributes Reference](#attributes-reference-21)
    - [onelogin\_roles (filters)](#onelogin_roles-filters)
      - [Example usage](#example-usage-23)
      - [Arguments Reference](#arguments-reference-22)
      - [Attributes Reference](#attributes-reference-22)
    - [onelogin\_roles\_admins (filters)](#onelogin_roles_admins-filters)
      - [Example usage](#example-usage-24)
      - [Arguments Reference](#arguments-reference-23)
      - [Attributes Reference](#attributes-reference-23)
    - [onelogin\_roles\_apps (filters)](#onelogin_roles_apps-filters)
      - [Example usage](#example-usage-25)
      - [Arguments Reference](#arguments-reference-24)
      - [Attributes Reference](#attributes-reference-24)
    - [onelogin\_roles\_users (filters)](#onelogin_roles_users-filters)
      - [Example usage](#example-usage-26)
      - [Arguments Reference](#arguments-reference-25)
      - [Attributes Reference](#attributes-reference-25)
    - [onelogin\_rules (filters)](#onelogin_rules-filters-1)
      - [Example usage](#example-usage-27)
      - [Arguments Reference](#arguments-reference-26)
      - [Attributes Reference](#attributes-reference-26)

## Provider Installation

In order to provision 'onelogin' Terraform resources, you need to first install the 'onelogin' Terraform plugin by running the following command (you must be running Terraform >= 0.12):

```shell
 export PROVIDER_NAME=onelogin && curl -fsSL https://raw.githubusercontent.com/dikhan/terraform-provider-openapi/master/scripts/install.sh | bash -s -- --provider-name $PROVIDER_NAME
```

```shell
[INFO] Downloading https://github.com/dikhan/terraform-provider-openapi/v3/releases/download/v3.0.0/terraform-provider-openapi_3.0.0_darwin_amd64.tar.gz in temporally folder /var/folders/n_/1lrwb99s7f50xmn9jpmfnddh0000gp/T/tmp.Xv1AkIZh...  
[INFO] Extracting terraform-provider-openapi from terraform-provider-openapi_0.29.4_darwin_amd64.tar.gz...  
[INFO] Cleaning up tmp dir created for installation purposes: /var/folders/n_/1lrwb99s7f50xmn9jpmfnddh0000gp/T/tmp.Xv1AkIZh  
[INFO] Terraform provider 'terraform-provider-onelogin' successfully installed at: '~/.terraform.d/plugins'!
```

You can then start running the Terraform provider:

```shell
 export OTF_VAR_onelogin_PLUGIN_CONFIGURATION_FILE="https://raw.githubusercontent.com/onelogin/terraform-provider-onelogin/openapi/swag-api.yml" && terraform init && terraform plan
```

**Note:** As of Terraform >= 0.13 each Terraform module must declare which providers it requires, so that Terraform can install and use them. If you are using Terraform >= 0.13, copy into your .tf file the following snippet already populated with the provider configuration:

```hcl
terraform {
  required_providers {
    onelogin = {
      source  = "onelogin.com/onelogin/onelogin"
      version = ">= 2.0.1" 
    }
  }
}
```

## Provider Configuration

### Example Usage

```hcl
provider "onelogin" {
  bearer_auth  = "..."
  content_type  = "..."
  authorization  = "..."
}
```

## Provider Resources

---

### onelogin\_apps

#### Example usage

```hcl
resource "onelogin_apps" "my_apps"{
}
```

#### Arguments Reference

The following arguments are supported:

- updated_at [string] - (Optional) The date the app was last updated.
- \* provisioning [object] - (Optional) . The following properties compose the object schema :
  - enabled [boolean] - (Optional) Indicates if provisioning is enabled for this app.
- name [string] - (Optional) App name.
- tab_id [integer] - (Optional) ID of the OneLogin portal tab that the app is assigned to.
- connector_id [integer] - (Optional) ID of the apps underlying connector.
- brand_id [integer] - (Optional) The custom login page branding to use for this app. Applies to app initiated logins via OIDC or SAML.
- allow_assumed_signin [boolean] - (Optional) Indicates whether or not administrators can access the app as a user that they have assumed control over.
- auth_method [integer] - (Optional) An ID indicating the type of app.
- policy_id [integer] - (Optional) The security policy assigned to the app.
- notes [string] - (Optional) Freeform notes about the app.
- id [integer] - (Optional) Apps unique ID in OneLogin.
- visible [boolean] - (Optional) Indicates if the app is visible in the OneLogin portal.
- role_ids [list of integers] - (Optional) A list of OneLogin Role IDs of the user
- created_at [string] - (Optional) The date the app was created.
- description [string] - (Optional) Freeform description of the app.
- icon_url [string] - (Optional) A link to the apps icon url.

\* Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin\_apps.my\_apps.**provisioning[0]**.object\_property)

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

\* Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin\_apps.my\_apps.**provisioning[0]**.object\_property)

#### Import

apps resources can be imported using the `id` , e.g:

```shell
terraform import onelogin_apps.my_apps id
```

**Note**: In order for the import to work, the 'onelogin' terraform provider must be [properly installed](#provider_installation). Read more about Terraform import usage [here](https://www.terraform.io/docs/import/usage.html).

### onelogin\_rules

#### Example usage

```hcl
resource "onelogin_rules" "my_rules"{
}
```

#### Arguments Reference

The following arguments are supported:

- id [string] - (Optional)
- type [string] - (Optional) The type parameter specifies the type of rule that will be created.
- target [string] - (Optional) The target parameter that will be used when evaluating the rule against an incoming event.
- filters [list of strings] - (Optional) A list of IP addresses or country codes or names to evaluate against each event.
- \* source [object] - (Optional) Used for targeting custom rules based on a group of people, customers, accounts, or even a single user.. The following properties compose the object schema :
  - name [string] - (Optional) The name of the source
  - id [string] - (Optional) A unique id that represents the source of the event.
- name [string] - (Optional) The name of this rule
- description [string] - (Optional)

\* Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin\_rules.my\_rules.**source[0]**.object\_property)

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

\* Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin\_rules.my\_rules.**source[0]**.object\_property)

#### Import

rules resources can be imported using the `id` , e.g:

```hcl
terraform import onelogin_rules.my_rules id
```

**Note**: In order for the import to work, the 'onelogin' terraform provider must be [properly installed](#provider_installation). Read more about Terraform import usage [here](https://www.terraform.io/docs/import/usage.html).

### onelogin\_users

#### Example usage

```hcl
resource "onelogin_users" "my_users"{
}
```

#### Arguments Reference

The following arguments are supported:

- manager\_ad\_id [string] - (Optional) The ID of the user's manager in Active Directory.
- salt [string] - (Optional) The salt value used with the password\_algorithm.
- password\_changed\_at [string] - (Optional)
- firstname [string] - (Optional) The user's first name.
- invitation\_sent\_at [string] - (Optional)
- password [string] - (Optional) The password to set for a user.
- username [string] - (Optional) A username for the user.
- status [integer] - (Optional)
- password\_confirmation [string] - (Optional) Required if the password is being set.
- password\_algorithm [string] - (Optional) Use this when importing a password that's already hashed. Prepend the salt value to the cleartext password value before SHA-256-encoding it
- phone [string] - (Optional) The E.164 format phone number for a user.
- invalid\_login\_attempts [integer] - (Optional)
- email [string] - (Optional) A valid email for the user.
- lastname [string] - (Optional) The user's last name.
- locked\_until [string] - (Optional)
- id [integer] - (Optional)
- title [string] - (Optional) The user's job title.
- userprincipalname [string] - (Optional) The principle name of the user.
- member\_of [string] - (Optional) The user's directory membership.
- role\_ids [list of integers] - (Optional) A list of OneLogin Role IDs of the user
- state [integer] - (Optional)
- updated\_at [string] - (Optional)
- trusted\_idp\_id [integer] - (Optional) The ID of the OneLogin Trusted IDP of the user.
- created\_at [string] - (Optional)
- preferred\_locale\_code [string] - (Optional)
- group\_id [integer] - (Optional) The ID of the Group in OneLogin that the user is assigned to.
- directory\_id [integer] - (Optional) The ID of the OneLogin Directory of the user.
- distinguished\_name [string] - (Optional) The distinguished name of the user.
- company [string] - (Optional) The user's company.
- manager\_user\_id [string] - (Optional) The OneLogin User ID for the user's manager.
- comment [string] - (Optional) Free text related to the user.
- samaccount\_name [string] - (Optional) The user's Active Directory username.
- department [string] - (Optional) The user's department.
- external\_id [string] - (Optional) The ID of the user in an external directory.
- activated\_at [string] - (Optional)
- last\_login [string] - (Optional)

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

#### Import

users resources can be imported using the `id` , e.g:

```shell
 terraform import onelogin_users.my_users id
```

**Note**: In order for the import to work, the 'onelogin' terraform provider must be [properly installed](#provider_installation). Read more about Terraform import usage [here](https://www.terraform.io/docs/import/usage.html).

## Data Sources (using resource id)

---

### onelogin\_apps\_instance

Retrieve an existing resource using it's ID

#### Example usage

```hcl
data "onelogin_apps_instance" "my_apps_instance"{
    id = "existing_resource_id"
}
```

#### Arguments Reference

The following arguments are supported:

- id - (Required) ID of the existing resource to retrieve

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- name [string] - App name.
- tab\_id [integer] - ID of the OneLogin portal tab that the app is assigned to.
- updated\_at [string] - The date the app was last updated.
- connector\_id [integer] - ID of the apps underlying connector.
- auth\_method [integer] - An ID indicating the type of app.
- allow\_assumed\_signin [boolean] - Indicates whether or not administrators can access the app as a user that they have assumed control over.
- brand\_id [integer] - The custom login page branding to use for this app. Applies to app initiated logins via OIDC or SAML.
- id [integer] - Apps unique ID in OneLogin.
- notes [string] - Freeform notes about the app.
- visible [boolean] - Indicates if the app is visible in the OneLogin portal.
- \* provisioning [object] The following properties compose the object schema:
  - enabled [boolean] - Indicates if provisioning is enabled for this app.
- policy\_id [integer] - The security policy assigned to the app.
- role\_ids [list of integers] - A list of OneLogin Role IDs of the user
- icon\_url [string] - A link to the apps icon url.
- created\_at [string] - The date the app was created.
- description [string] - Freeform description of the app.

\* Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin\_apps\_instance.my_apps_instance.**provisioning[0]**.object\_property)

### onelogin\_rules\_instance

Retrieve an existing resource using it's ID

#### Example usage

```hcl
data "onelogin_rules_instance" "my_rules_instance"{
    id = "existing_resource_id"
}
```

#### Arguments Reference

The following arguments are supported:

- id - (Required) ID of the existing resource to retrieve

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- id [string]
- type [string] - The type parameter specifies the type of rule that will be created.
- filters [list of strings] - A list of IP addresses or country codes or names to evaluate against each event.
- target [string] - The target parameter that will be used when evaluating the rule against an incoming event.
- \* source [object] - Used for targeting custom rules based on a group of people, customers, accounts, or even a single user. The following properties compose the object schema:
  - name [string] - The name of the source
  - id [string] - A unique id that represents the source of the event.
- name [string] - The name of this rule
- description [string]

\* Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin\_rules\_instance.my\_rules\_instance.**source[0]**.object\_property)

### onelogin\_users\_instance

Retrieve an existing resource using it's ID

#### Example usage

```hcl
data "onelogin_users_instance" "my_users_instance"{
    id = "existing_resource_id"
}
```

#### Arguments Reference

The following arguments are supported:

- id - (Required) ID of the existing resource to retrieve

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- invitation\_sent\_at [string]
- firstname [string] - The user's first name.
- salt [string] - The salt value used with the password\_algorithm.
- password\_changed\_at [string]
- manager\_ad\_id [string] - The ID of the user's manager in Active Directory.
- phone [string] - The E.164 format phone number for a user.
- password\_algorithm [string] - Use this when importing a password that's already hashed. Prepend the salt value to the cleartext password value before SHA-256-encoding it
- password\_confirmation [string] - Required if the password is being set.
- password [string] - The password to set for a user.
- status [integer]
- username [string] - A username for the user.
- locked\_until [string]
- lastname [string] - The user's last name.
- email [string] - A valid email for the user.
- invalid\_login\_attempts [integer]
- userprincipalname [string] - The principle name of the user.
- member\_of [string] - The user's directory membership.
- title [string] - The user's job title.
- id [integer]
- updated\_at [string]
- state [integer]
- role\_ids [list of integers] - A list of OneLogin Role IDs of the user
- group\_id [integer] - The ID of the Group in OneLogin that the user is assigned to.
- preferred\_locale\_code [string]
- directory\_id [integer] - The ID of the OneLogin Directory of the user.
- created\_at [string]
- trusted\_idp\_id [integer] - The ID of the OneLogin Trusted IDP of the user.
- company [string] - The user's company.
- distinguished\_name [string] - The distinguished name of the user.
- activated\_at [string]
- external\_id [string] - The ID of the user in an external directory.
- last\_login [string]
- comment [string] - Free text related to the user.
- samaccount\_name [string] - The user's Active Directory username.
- department [string] - The user's department.
- manager\_user\_id [string] - The OneLogin User ID for the user's manager.

## Data Sources (using filters)

---

### onelogin\_api\_authorizations (filters)

The api\_authorizations data source allows you to retrieve an already existing api\_authorizations resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_api_authorizations" "my_api_authorizations"{
    filter  {
        name  = "property name to filter by, see docs below for more info about available filter name options"
        values  = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: description, name, id,
- values \[array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- description [string] - Description of what the API does.
- name [string] - Name of the API.
- id [integer] - Auth server unique ID in Onelogin
- \* configuration [object] - Authorization server configuration The following properties compose the object schema:
  - access\_token\_expiration\_minutes [integer] - The number of minutes until access token expires. There is no maximum expiry limit.
  - refresh\_token\_expiration\_minutes [integer] - The number of minutes until refresh token expires. There is no maximum expiry limit.
  - audiences [list of strings] - List of API endpoints that will be returned in Access Tokens.
  - resource\_identifier [string] - Unique identifier for the API that the Authorization Server will issue Access Tokens for.

\* Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin\_api\_authorizations.my\_api\_authorizations.**configuration[0]**.object\_property)

### onelogin\_api\_authorizations\_claims (filters)

The api\_authorizations\_claims data source allows you to retrieve an already existing api\_authorizations\_claims resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_api_authoriztions_claims" "my_api_authorizations_claims"{
    filter  {
        name  = "property name to filter by, see docs below for more info about available filter name options"
        values  = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: provisioned\_entitlements, skip\_if\_blank, id, user\_attribute\_mappings, attribute\_transformations, api\_authorizations\_id, user\_attribute\_macros, default\_values, label,
- values \[array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- provisioned\_entitlements [boolean] - Relates to Rules/Entitlements. Not supported yet.
- skip\_if\_blank [boolean] - not used
- id [integer] - The unique ID of the claim.
- user\_attribute\_mappings [string] - A user attribute to map values from.
- attribute\_transformations [string] - The type of transformation to perform on multi valued attributes.
- user\_attribute\_macros [string] - When \`user\_attribute\_mappings\` is set to \`\_macro\_\` this macro will be used to assign the claims value.
- values [list of strings] - Relates to Rules/Entitlements. Not supported yet.
- default\_values [string] - Relates to Rules/Entitlements. Not supported yet.
- label [string] - The UI label for the claims.

### onelogin\_api\_authorizations\_scopes (filters)

The api\_authorizations\_scopes data source allows you to retrieve an already existing api\_authorizations\_scopes resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_api_authoriztions_scopes" "my_api_authorizations_scopes"{
    filter  {
        name  = "property name to filter by, see docs below for more info about available filter name options"
        values  = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: description, name, id, api\_authorizations\_id,
- values \[array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- description [string] - Description of what the API does.
- name [string] - Name of the API.
- id [integer] - Auth server unique ID in Onelogin
- \* configuration [object] - Authorization server configuration The following properties compose the object schema:
  - access\_token\_expiration\_minutes [integer] - The number of minutes until access token expires. There is no maximum expiry limit.
  - refresh\_token\_expiration\_minutes [integer] - The number of minutes until refresh token expires. There is no maximum expiry limit.
  - audiences [list of strings] - List of API endpoints that will be returned in Access Tokens.
  - resource\_identifier [string] - Unique identifier for the API that the Authorization Server will issue Access Tokens for.

\* Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin\_api\_authorizations\_scopes.my\_api\_authorizations\_scopes.**configuration[0]**.object\_property)

### onelogin\_apps\_actions (filters)

The apps\_actions data source allows you to retrieve an already existing apps\_actions resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_apps_actions "my_apps_actions"{
    filter  {
        name  = "property name to filter by, see docs below for more info about available filter name options"
        values  = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: apps\_id, name, value,
- values \[array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- name [string] - Name of the Action
- value [string] - The unique identifier of the action. This should be used when defining actions for a User Mapping.

### onelogin\_apps\_actions\_values (filters)

The apps\_actions\_values data source allows you to retrieve an already existing apps\_actions\_values resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_apps_actions_values" "my_apps_actions_values"{
    filter  {
        name  = "property name to filter by, see docs below for more info about available filter name options"
        values  = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: actions\_id, apps\_id, name, value,
- values \[array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- name [string] - Name of the Action
- value [string] - The unique identifier of the action. This should be used when defining actions for a User Mapping.

### onelogin\_apps\_conditions (filters)

The apps\_conditions data source allows you to retrieve an already existing apps\_conditions resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_apps_conditins" "my_apps_conditions"{
    filter  {
        name  = "property name to filter by, see docs below for more info about available filter name options"
        values  = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: value, apps\_id, name,
- values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- value [string] - The unique identifier of the condition. This should be used when defining conditions for a rule.
- name [string] - Name of the rule condition

### onelogin\_apps\_conditions\_operators (filters)

The apps\_conditions\_operators data source allows you to retrieve an already existing apps\_conditions\_operators resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_apps_conditins_operators" "my_apps_conditions_operators"{
    filter  {
        name  = "property name to filter by, see docs below for more info about available filter name options"
        values  = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: name, apps\_id, value, conditions\_id,
- values \[array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- name [string] - Name of the operator
- value [string] - The condition operator value to use when creating or updating rules.

### onelogin\_apps\_conditions\_values (filters)

The apps\_conditions\_values data source allows you to retrieve an already existing apps\_conditions\_values resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_apps_conditins_values" "my_apps_conditions_values"{
    filter  {
        name  = "property name to filter by, see docs below for more info about available filter name options"
        values  = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: value, apps\_id, name, conditions\_id,
- values \[array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- value [string] - The unique identifier of the condition. This should be used when defining conditions for a rule.
- name [string] - Name of the rule condition

### onelogin\_apps\_rules (filters)

The apps\_rules data source allows you to retrieve an already existing apps\_rules resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_apps_rules" my_apps_rules"{
    filter  {
        name  = "property name to filter by, see docs below for more info about available filter name options"
        values  = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: position, match, name, apps\_id, enabled, id,
- values \[array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- position [integer] - Indicates the order of the rule. When \`""\` this will default to last position.
- match [string] - Indicates how conditions should be matched.
- conditions [list of objects] - An array of conditions that the user must meet in order for the rule to be applied. The following properties compose the object schema:
  - source [string] - source field to check.
  - operator [string] - A valid operator for the selected condition source
  - value [string] - A plain text string or valid value for the selected condition source
- name [string] - Rule Name
- enabled [boolean] - Indicates if the rule is enabled or not.
- id [integer] - App Rule ID
- actions [list of objects] The following properties compose the object schema:
  - value [list of strings] - Only applicable to provisioned and set\_\* actions. Items in the array will be a plain text string or valid value for the selected action.
  - action [string] - The action to apply

### onelogin\_apps\_users (filters)

The apps\_users data source allows you to retrieve an already existing apps\_users resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_apps_users" my_apps_users"{
    filter  {
        name  = "property name to filter by, see docs below for more info about available filter name options"
        values  = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: invitation\_sent\_at, firstname, salt, password\_changed\_at, manager\_ad\_id, phone, password\_algorithm, password\_confirmation, password, status, username, locked\_until, lastname, email, apps\_id, invalid\_login\_attempts, userprincipalname, member\_of, title, id, updated\_at, state, group\_id, preferred\_locale\_code, directory\_id, created\_at, trusted\_idp\_id, company, distinguished\_name, activated\_at, external\_id, last\_login, comment, samaccount\_name, department, manager\_user\_id,
- values \[array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- invitation\_sent\_at [string]
- firstname [string] - The user's first name.
- salt [string] - The salt value used with the password\_algorithm.
- password\_changed\_at [string]
- manager\_ad\_id [string] - The ID of the user's manager in Active Directory.
- phone [string] - The E.164 format phone number for a user.
- password\_algorithm [string] - Use this when importing a password that's already hashed. Prepend the salt value to the cleartext password value before SHA-256-encoding it
- password\_confirmation [string] - Required if the password is being set.
- password [string] - The password to set for a user.
- status [integer]
- username [string] - A username for the user.
- locked\_until [string]
- lastname [string] - The user's last name.
- email [string] - A valid email for the user.
- invalid\_login\_attempts [integer]
- userprincipalname [string] - The principle name of the user.
- member\_of [string] - The user's directory membership.
- title [string] - The user's job title.
- id [integer]
- updated\_at [string]
- state [integer]
- role\_ids [list of integers] - A list of OneLogin Role IDs of the user
- group\_id [integer] - The ID of the Group in OneLogin that the user is assigned to.
- preferred\_locale\_code [string]
- directory\_id [integer] - The ID of the OneLogin Directory of the user.
- created\_at [string]
- trusted\_idp\_id [integer] - The ID of the OneLogin Trusted IDP of the user.
- company [string] - The user's company.
- distinguished\_name [string] - The distinguished name of the user.
- activated\_at [string]
- external\_id [string] - The ID of the user in an external directory.
- last\_login [string]
- comment [string] - Free text related to the user.
- samaccount\_name [string] - The user's Active Directory username.
- department [string] - The user's department.
- manager\_user\_id [string] - The OneLogin User ID for the user's manager.

### onelogin\_mappings (filters)

The mappings data source allows you to retrieve an already existing mappings resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_mappings" "my_mappings"{
    filter  {
        name  = "property name to filter by, see docs below for more info about available filter name options"
        values  = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: match, name, id, enabled, position,
- values \[array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- match [string] - Indicates how conditions should be matched.
- name [string] - The name of the mapping.
- actions [list of objects] - An array of actions that will be applied to the users that are matched by the conditions. The following properties compose the object schema:
  - value [list of strings] - Only applicable to provisioned and set\_\* actions. Items in the array will be a plain text string or valid value for the selected action.
  - action [string] - The action to apply
- id [integer]
- enabled [boolean] - Indicates if the mapping is enabled or not.
- position [integer] - Indicates the order of the mapping. When \`""\` this will default to last position.
- conditions [list of objects] - An array of conditions that the user must meet in order for the mapping to be applied. The following properties compose the object schema:
  - source [string] - source field to check.
  - operator [string] - A valid operator for the selected condition source
  - value [string] - A plain text string or valid value for the selected condition source

### onelogin\_privileges (filters)

The privileges data source allows you to retrieve an already existing privileges resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_privileges" "my_privileges"{
    filter  {
        name  = "property name to filter by, see docs below for more info about available filter name options"
        values  = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: id, description, name,
- values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- id [string]
- \* privilege [object] The following properties compose the object schema:
  - statement [list of objects] The following properties compose the object schema:
    - scope [list of strings] - Target the privileged action against specific resources with the scope. The scope pattern is the class of object used by the Action, followed by an ID that represents a resource in OneLogin. e.g. apps/1234, where apps is the class and 1234 is the ID of an app. The wildcard \* is supported and indicates that all resources of the class type declared, in the Action, are in scope. The Action and Scope classes must match. However, there is an exception, a scope of roles/{role\_id} can be combined with Actions on the user or app class. The exception allows you to target groups of users or apps with specific actions.
    - action [list of strings] - An array of strings that represent actions within OneLogin. Actions are prefixed with the class of object they are related to and followed by a specific action for the given class. e.g. users:List, where the class is users and the specific action is List. Don’t mix classes within an Action array. To create a privilege that includes multiple different classes, create multiple statements. A wildcard \* that includes all actions is supported. Use wildcards to create a Super User privilege.
    - effect [string] - Set to “Allow.” By default, all actions are denied, this Statement allows the listed actions to be executed.
  - version [string]
- description [string]
- name [string]

\* Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin\_privileges.my\_privileges.**privilege[0]**.object\_property)

### onelogin\_rules (filters)

The rules data source allows you to retrieve an already existing rules resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_rules" "my_rules"{
    filter  {
        name  = "property name to filter by, see docs below for more info about available filter name options"
        values  = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: id, type, target, name, description,
- values \[array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- id [string]
- type [string] - The type parameter specifies the type of rule that will be created.
- filters [list of strings] - A list of IP addresses or country codes or names to evaluate against each event.
- target [string] - The target parameter that will be used when evaluating the rule against an incoming event.
- \* source [object] - Used for targeting custom rules based on a group of people, customers, accounts, or even a single user. The following properties compose the object schema:
  - name [string] - The name of the source
  - id [string] - A unique id that represents the source of the event.
- name [string] - The name of this rule
- description [string]

\* Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin\_rules.my\_rules.**source[0]**.object\_property)

### onelogin\_users (filters)

The users data source allows you to retrieve an already existing users resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_users" "my_users"{
    filter  {
        name  = "property name to filter by, see docs below for more info about available filter name options"
        values  = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: invitation\_sent\_at, firstname, salt, password\_changed\_at, manager\_ad\_id, phone, password\_algorithm, password\_confirmation, password, status, username, locked\_until, lastname, email, invalid\_login\_attempts, userprincipalname, member\_of, title, id, updated\_at, state, group\_id, preferred\_locale\_code, directory\_id, created\_at, trusted\_idp\_id, company, distinguished\_name, activated\_at, external\_id, last\_login, comment, samaccount\_name, department, manager\_user\_id,
- values \[array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- invitation\_sent\_at [string]
- firstname [string] - The user's first name.
- salt [string] - The salt value used with the password\_algorithm.
- password\_changed\_at [string]
- manager\_ad\_id [string] - The ID of the user's manager in Active Directory.
- phone [string] - The E.164 format phone number for a user.
- password\_algorithm [string] - Use this when importing a password that's already hashed. Prepend the salt value to the cleartext password value before SHA-256-encoding it
- password\_confirmation [string] - Required if the password is being set.
- password [string] - The password to set for a user.
- status [integer]
- username [string] - A username for the user.
- locked\_until [string]
- lastname [string] - The user's last name.
- email [string] - A valid email for the user.
- invalid\_login\_attempts [integer]
- userprincipalname [string] - The principle name of the user.
- member\_of [string] - The user's directory membership.
- title [string] - The user's job title.
- id [integer]
- updated\_at [string]
- state [integer]
- role\_ids [list of integers] - A list of OneLogin Role IDs of the user
- group\_id [integer] - The ID of the Group in OneLogin that the user is assigned to.
- preferred\_locale\_code [string]
- directory\_id [integer] - The ID of the OneLogin Directory of the user.
- created\_at [string]
- trusted\_idp\_id [integer] - The ID of the OneLogin Trusted IDP of the user.
- company [string] - The user's company.
- distinguished\_name [string] - The distinguished name of the user.
- activated\_at [string]
- external\_id [string] - The ID of the user in an external directory.
- last\_login [string]
- comment [string] - Free text related to the user.
- samaccount\_name [string] - The user's Active Directory username.
- department [string] - The user's department.
- manager\_user\_id [string] - The OneLogin User ID for the user's manager.

### onelogin\_users\_apps (filters)

The users\_apps data source allows you to retrieve an already existing users\_apps resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_users_apps" "my_users_apps"{
    filter  {
        name  = "property name to filter by, see docs below for more info about available filter name options"
        values  = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: login\_id, extension, name, icon\_url, id, provisioning\_state, provisioning\_enabled, users\_id, provisioning\_status,
- values \[array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- login\_id [integer] - Unqiue identifier for this user and app combination.
- extension [boolean] - Boolean that indicates if the OneLogin browser extension is required to launch this app.
- name [string] - The name of the app.
- icon\_url [string] - A url for the icon that represents the app in the OneLogin portal
- id [integer] - The App ID
- provisioning\_state [string] - If provisioning is enabled this indicates the state of provisioning for the given user.
- provisioning\_enabled [boolean] - Indicates if provisioning is enabled for this app.
- provisioning\_status [string]

### onelogin\_users\_devices (filters)

The users\_devices data source allows you to retrieve an already existing users\_devices resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_users_devices" "my_users_devices"{
    filter  {
        name  = "property name to filter by, see docs below for more info about available filter name options"
        values  = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: device\_id, default, user\_display\_name, auth\_factor\_name, users\_id, type\_display\_name,
- values \[array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- device\_id [string] - MFA device identifier.
- default [boolean] - true = is user’s default MFA device for OneLogin.
- user\_display\_name [string] - Authentication factor display name assigned by users when they register the device.
- auth\_factor\_name [string] - Authentication factor name, as it appears to administrators in OneLogin.
- type\_display\_name [string] - Authentication factor display name as it appears to users upon initial registration, as defined by admins at Settings > Authentication Factors.

### onelogin\_roles (filters)

The roles data source allows you to retrieve an already existing roles resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_roles" "my_roles"{
    filter  {
        name  = "property name to filter by, see docs below for more info about available filter name options"
        values  = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: id, name,
- values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- id [integer]
- apps [list of integers]
- admins [list of integers]
- users [list of integers]
- name [string]

### onelogin\_roles\_admins (filters)

The roles\_admins data source allows you to retrieve an already existing roles\_admins resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_roles_admins" "my_roles_admins"{
    filter  {
        name  = "property name to filter by, see docs below for more info about available filter name options"
        values  = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: invitation\_sent\_at, firstname, salt, password\_changed\_at, manager\_ad\_id, phone, password\_algorithm, password\_confirmation, password, status, username, locked\_until, lastname, email, invalid\_login\_attempts, userprincipalname, member\_of, title, id, updated\_at, state, group\_id, preferred\_locale\_code, directory\_id, created\_at, trusted\_idp\_id, company, distinguished\_name, activated\_at, external\_id, last\_login, comment, samaccount\_name, department, roles\_id, manager\_user\_id,
- values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- invitation\_sent\_at [string]
- firstname [string] - The user's first name.
- salt [string] - The salt value used with the password\_algorithm.
- password\_changed\_at [string]
- manager\_ad\_id [string] - The ID of the user's manager in Active Directory.
- phone [string] - The E.164 format phone number for a user.
- password\_algorithm [string] - Use this when importing a password that's already hashed. Prepend the salt value to the cleartext password value before SHA-256-encoding it
- password\_confirmation [string] - Required if the password is being set.
- password [string] - The password to set for a user.
- status [integer]
- username [string] - A username for the user.
- locked\_until [string]
- lastname [string] - The user's last name.
- email [string] - A valid email for the user.
- invalid\_login\_attempts [integer]
- userprincipalname [string] - The principle name of the user.
- member\_of [string] - The user's directory membership.
- title [string] - The user's job title.
- id [integer]
- updated\_at [string]
- state [integer]
- role\_ids [list of integers] - A list of OneLogin Role IDs of the user
- group\_id [integer] - The ID of the Group in OneLogin that the user is assigned to.
- preferred\_locale\_code [string]
- directory\_id [integer] - The ID of the OneLogin Directory of the user.
- created\_at [string]
- trusted\_idp\_id [integer] - The ID of the OneLogin Trusted IDP of the user.
- company [string] - The user's company.
- distinguished\_name [string] - The distinguished name of the user.
- activated\_at [string]
- external\_id [string] - The ID of the user in an external directory.
- last\_login [string]
- comment [string] - Free text related to the user.
- samaccount\_name [string] - The user's Active Directory username.
- department [string] - The user's department.
- manager\_user\_id [string] - The OneLogin User ID for the user's manager.

### onelogin\_roles\_apps (filters)

The roles\_apps data source allows you to retrieve an already existing roles\_apps resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_roles_apps" "my_roles_apps"{
    filter  {
        name  = "property name to filter by, see docs below for more info about available filter name options"
        values  = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: icon\_url, id, name, roles\_id,
- values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- icon\_url [string] - url of Icon
- id [integer] - app id
- name [string] - app name

### onelogin\_roles\_users (filters)

The roles\_users data source allows you to retrieve an already existing roles\_users resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_roles_users" "my_roles_users"{
    filter  {
        name  = "property name to filter by, see docs below for more info about available filter name options"
        values  = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: invitation\_sent\_at, firstname, salt, password\_changed\_at, manager\_ad\_id, phone, password\_algorithm, password\_confirmation, password, status, username, locked\_until, lastname, email, invalid\_login\_attempts, userprincipalname, member\_of, title, id, updated\_at, state, group\_id, preferred\_locale\_code, directory\_id, created\_at, trusted\_idp\_id, company, distinguished\_name, activated\_at, external\_id, last\_login, comment, samaccount\_name, department, roles\_id, manager\_user\_id,
- values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- invitation\_sent\_at [string]
- firstname [string] - The user's first name.
- salt [string] - The salt value used with the password\_algorithm.
- password\_changed\_at [string]
- manager\_ad\_id [string] - The ID of the user's manager in Active Directory.
- phone [string] - The E.164 format phone number for a user.
- password\_algorithm [string] - Use this when importing a password that's already hashed. Prepend the salt value to the cleartext password value before SHA-256-encoding it
- password\_confirmation [string] - Required if the password is being set.
- password [string] - The password to set for a user.
- status [integer]
- username [string] - A username for the user.
- locked\_until [string]
- lastname [string] - The user's last name.
- email [string] - A valid email for the user.
- invalid\_login\_attempts [integer]
- userprincipalname [string] - The principle name of the user.
- member\_of [string] - The user's directory membership.
- title [string] - The user's job title.
- id [integer]
- updated\_at [string]
- state [integer]
- role\_ids [list of integers] - A list of OneLogin Role IDs of the user
- group\_id [integer] - The ID of the Group in OneLogin that the user is assigned to.
- preferred\_locale\_code [string]
- directory\_id [integer] - The ID of the OneLogin Directory of the user.
- created\_at [string]
- trusted\_idp\_id [integer] - The ID of the OneLogin Trusted IDP of the user.
- company [string] - The user's company.
- distinguished\_name [string] - The distinguished name of the user.
- activated\_at [string]
- external\_id [string] - The ID of the user in an external directory.
- last\_login [string]
- comment [string] - Free text related to the user.
- samaccount\_name [string] - The user's Active Directory username.
- department [string] - The user's department.
- manager\_user\_id [string] - The OneLogin User ID for the user's manager.

### onelogin\_rules (filters)

The rules data source allows you to retrieve an already existing rules resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_rules" "my_rules"{
    filter  {
        name  = "property name to filter by, see docs below for more info about available filter name options"
        values  = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: id, type, target, name, description,
- values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- id [string]
- type [string] - The type parameter specifies the type of rule that will be created.
- filters [list of strings] - A list of IP addresses or country codes or names to evaluate against each event.
- target [string] - The target parameter that will be used when evaluating the rule against an incoming event.
- \* source [object] - Used for targeting custom rules based on a group of people, customers, accounts, or even a single user. The following properties compose the object schema:
  - name [string] - The name of the source
  - id [string] - A unique id that represents the source of the event.
- name [string] - The name of this rule
- description [string]

\* Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin\_rules.my\_rules.**source[0]**.object\_property)
