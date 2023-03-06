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
    - [onelogin\_users](#onelogin_users-1)
      - [Example usage](#example-usage-4)
      - [Arguments Reference](#arguments-reference-3)
      - [Attributes Reference](#attributes-reference-3)
      - [Import](#import-3)
  - [Data Sources (using resource id)](#data-sources-using-resource-id)
    - [onelogin\_apps\_instance](#onelogin_apps_instance)
      - [Example usage](#example-usage-5)
      - [Arguments Reference](#arguments-reference-4)
      - [Attributes Reference](#attributes-reference-4)
    - [onelogin\_rules\_instance](#onelogin_rules_instance)
      - [Example usage](#example-usage-6)
      - [Arguments Reference](#arguments-reference-5)
      - [Attributes Reference](#attributes-reference-5)
    - [onelogin\_users\_instance](#onelogin_users_instance)
      - [Example usage](#example-usage-7)
      - [Arguments Reference](#arguments-reference-6)
      - [Attributes Reference](#attributes-reference-6)
    - [onelogin\_users\_instance](#onelogin_users_instance-1)
      - [Example usage](#example-usage-8)
      - [Arguments Reference](#arguments-reference-7)
      - [Attributes Reference](#attributes-reference-7)
  - [Data Sources (using filters)](#data-sources-using-filters)
    - [onelogin\_api\_authorizations (filters)](#onelogin_api_authorizations-filters)
      - [Example usage](#example-usage-9)
      - [Arguments Reference](#arguments-reference-8)
      - [Attributes Reference](#attributes-reference-8)
    - [onelogin\_api\_authorizations\_claims (filters)](#onelogin_api_authorizations_claims-filters)
      - [Example usage](#example-usage-10)
      - [Arguments Reference](#arguments-reference-9)
      - [Attributes Reference](#attributes-reference-9)
    - [onelogin\_api\_authorizations\_scopes (filters)](#onelogin_api_authorizations_scopes-filters)
      - [Example usage](#example-usage-11)
      - [Arguments Reference](#arguments-reference-10)
      - [Attributes Reference](#attributes-reference-10)
    - [onelogin\_apps (filters)](#onelogin_apps-filters)
      - [Example usage](#example-usage-12)
      - [Arguments Reference](#arguments-reference-11)
      - [Attributes Reference](#attributes-reference-11)
    - [onelogin\_apps\_actions (filters)](#onelogin_apps_actions-filters)
      - [Example usage](#example-usage-13)
      - [Arguments Reference](#arguments-reference-12)
      - [Attributes Reference](#attributes-reference-12)
    - [onelogin\_apps\_actions\_values (filters)](#onelogin_apps_actions_values-filters)
      - [Example usage](#example-usage-14)
      - [Arguments Reference](#arguments-reference-13)
      - [Attributes Reference](#attributes-reference-13)
    - [onelogin\_apps\_conditions (filters)](#onelogin_apps_conditions-filters)
      - [Example usage](#example-usage-15)
      - [Arguments Reference](#arguments-reference-14)
      - [Attributes Reference](#attributes-reference-14)
    - [onelogin\_apps\_conditions\_operators (filters)](#onelogin_apps_conditions_operators-filters)
      - [Example usage](#example-usage-16)
      - [Arguments Reference](#arguments-reference-15)
      - [Attributes Reference](#attributes-reference-15)
    - [onelogin\_apps\_conditions\_values (filters)](#onelogin_apps_conditions_values-filters)
      - [Example usage](#example-usage-17)
      - [Arguments Reference](#arguments-reference-16)
      - [Attributes Reference](#attributes-reference-16)
    - [onelogin\_apps\_rules (filters)](#onelogin_apps_rules-filters)
      - [Example usage](#example-usage-18)
      - [Arguments Reference](#arguments-reference-17)
      - [Attributes Reference](#attributes-reference-17)
    - [onelogin\_apps\_users (filters)](#onelogin_apps_users-filters)
      - [Example usage](#example-usage-19)
      - [Arguments Reference](#arguments-reference-18)
      - [Attributes Reference](#attributes-reference-18)
    - [onelogin\_brands (filters)](#onelogin_brands-filters)
      - [Example usage](#example-usage-20)
      - [Arguments Reference](#arguments-reference-19)
      - [Attributes Reference](#attributes-reference-19)
    - [onelogin\_brands\_apps (filters)](#onelogin_brands_apps-filters)
      - [Example usage](#example-usage-21)
      - [Arguments Reference](#arguments-reference-20)
      - [Attributes Reference](#attributes-reference-20)
    - [onelogin\_brands\_templates (filters)](#onelogin_brands_templates-filters)
      - [Example usage](#example-usage-22)
      - [Arguments Reference](#arguments-reference-21)
      - [Attributes Reference](#attributes-reference-21)
    - [onelogin\_mappings (filters)](#onelogin_mappings-filters)
      - [Example usage](#example-usage-23)
      - [Arguments Reference](#arguments-reference-22)
      - [Attributes Reference](#attributes-reference-22)
    - [onelogin\_privileges (filters)](#onelogin_privileges-filters)
      - [Example usage](#example-usage-24)
      - [Arguments Reference](#arguments-reference-23)
      - [Attributes Reference](#attributes-reference-23)
    - [onelogin\_roles (filters)](#onelogin_roles-filters)
      - [Example usage](#example-usage-25)
      - [Arguments Reference](#arguments-reference-24)
      - [Attributes Reference](#attributes-reference-24)
    - [onelogin\_roles-admins (filters)](#onelogin_roles-admins-filters)
      - [Example usage](#example-usage-26)
      - [Arguments Reference](#arguments-reference-25)
      - [Attributes Reference](#attributes-reference-25)
    - [onelogin\_roles-apps (filters)](#onelogin_roles-apps-filters)
      - [Example usage](#example-usage-27)
      - [Arguments Reference](#arguments-reference-26)
      - [Attributes Reference](#attributes-reference-26)
    - [onelogin\_roles-users (filters)](#onelogin_roles-users-filters)
      - [Example usage](#example-usage-28)
      - [Arguments Reference](#arguments-reference-27)
      - [Attributes Reference](#attributes-reference-27)
    - [onelogin\_rules (filters)](#onelogin_rules-filters)
      - [Example usage](#example-usage-29)
      - [Arguments Reference](#arguments-reference-28)
      - [Attributes Reference](#attributes-reference-28)
    - [onelogin\_users (filters)](#onelogin_users-filters)
      - [Example usage](#example-usage-30)
      - [Arguments Reference](#arguments-reference-29)
      - [Attributes Reference](#attributes-reference-29)
    - [onelogin\_users (filters)](#onelogin_users-filters-1)
      - [Example usage](#example-usage-31)
      - [Arguments Reference](#arguments-reference-30)
      - [Attributes Reference](#attributes-reference-30)
    - [onelogin\_users\_apps (filters)](#onelogin_users_apps-filters)
      - [Example usage](#example-usage-32)
      - [Arguments Reference](#arguments-reference-31)
      - [Attributes Reference](#attributes-reference-31)
    - [onelogin\_users\_apps (filters)](#onelogin_users_apps-filters-1)
      - [Example usage](#example-usage-33)
      - [Arguments Reference](#arguments-reference-32)
      - [Attributes Reference](#attributes-reference-32)
    - [onelogin\_users\_devices (filters)](#onelogin_users_devices-filters)
      - [Example usage](#example-usage-34)
      - [Arguments Reference](#arguments-reference-33)
      - [Attributes Reference](#attributes-reference-33)

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
export OTF_VAR_onelogin_PLUGIN_CONFIGURATION_FILE="https://api.service.com/openapi.yaml"  
➜ ~ terraform init && terraform plan
```

**Note:** As of Terraform >= 0.13 each Terraform module must declare which providers it requires, so that Terraform can install and use them. If you are using Terraform >= 0.13, copy into your .tf file the following snippet already populated with the provider configuration:

```hcl
terraform {
  required_providers {
    onelogin = {
      source  = "onelogin/onelogin"
      version = ">= 2.0.1" 
    }
  }
}
```

## Provider Configuration

#### Example Usage

```hcl
provider "onelogin" {
 content_type = "..."
}
```

## Provider Resources

### onelogin_apps

#### Example usage

```hcl
resource "onelogin_apps" "my_apps"{
}
```

#### Arguments Reference

The following arguments are supported:

- * provisioning [object] - (Optional) Indicates if provisioning is enabled for this app.. The following properties compose the object schema :
  - enabled [boolean] - (Optional)
- tab_id [integer] - (Optional) ID of the OneLogin portal tab that the app is assigned to.
- name [string] - (Optional) The name of the app.
- role_ids [list of integers] - (Optional) List of Role IDs that are assigned to the app. On App Create or Update the entire array is replaced with the values provided.
- allow_assumed_signin [boolean] - (Optional) Indicates whether or not administrators can access the app as a user that they have assumed control over.
- auth_method [integer] - (Optional) An ID indicating the type of app: - 0: Password - 1: OpenId - 2: SAML - 3: API - 4: Google - 6: Forms Based App - 7: WSFED - 8: OpenId Connect
- policy_id [integer] - (Optional) The security policy assigned to the app.
- * parameters [object] - (Optional) The parameters section contains parameterized attributes that have defined at the connector level as well as custom attributes that have been defined specifically for this app. Regardless of how they are defined, all parameters have the following attributes. Each parameter is an object with the key for the object being set as the parameters short name.. The following properties compose the object schema :
  - include_in_saml_assertion [boolean] - (Optional) When true, this parameter will be included in a SAML assertion payload.
  - label [string] - (Optional) The can only be set when creating a new parameter. It can not be updated.
  - user_attribute_mappings [string] - (Optional) A user attribute to map values from For custom attributes prefix the name of the attribute with `custom_attribute_`. e.g. To get the value for custom attribute `employee_id` use `custom_attribute_employee_id`.
  - user_attribute_macros [string] - (Optional) When `user_attribute_mappings` is set to `_macro_` this macro will be used to assign the parameter value.
- notes [string] - (Optional) Freeform notes about the app.
- * enforcement_point [object] - (Optional) For apps that connect to a OneLogin Access Enforcement Point the following enforcement_point object will be included with the app payload.. The following properties compose the object schema :
  - * session_expiry_inactivity [object] - (Optional) unit: - 0 = Seconds - 1 = Minutes - 2 = Hours value: - When Unit = 0 or 1 value must be 0-60 - When Unit = 2 value must be 0-24\. The following properties compose the object schema :
    - unit [integer] - (Optional)
    - value [integer] - (Optional)
  - case_sensitive [boolean] - (Optional) The URL path evaluation is case insensitive by default. Resources hosted on web servers such as Apache, NGINX and Java EE are case sensitive paths. Web servers such as Microsoft IIS are not case-sensitive.
  - permissions [string] - (Optional) Specify to always `allow`, `deny` access to resources, of if access is `conditional`.
  - target [string] - (Optional) A fully-qualified URL to the internal application including scheme, authority and path. The target host authority must be an IP address, not a hostname.
  - vhost [string] - (Optional) A comma-delimited list of one or more virtual hosts that map to applications assigned to the enforcement point. A VHOST may be a host name or an IP address. VHOST distinguish between applications that are at the same context root.
  - require_sitewide_authentication [boolean] - (Optional) Require user authentication to access any resource protected by this enforcement point.
  - * session_expiry_fixed [object] - (Optional) unit: - 0 = Seconds - 1 = Minutes - 2 = Hours value: - When Unit = 0 or 1 value must be 0-60 - When Unit = 2 value must be 0-24\. The following properties compose the object schema :
    - unit [integer] - (Optional)
    - value [integer] - (Optional)
  - use_target_host_header [boolean] - (Optional) Use the target host header as opposed to the original gateway or upstream host header.
  - landing_page [string] - (Optional) The location within the context root to which the browser will be redirected for IdP-initiated single sign-on. For example, the landing page might be an index page in the context root such as index.html or default.aspx. The landing page cannot begin with a slash and must use valid URL characters.
  - resources [list of objects] - (Optional) Array of resource objects. The following properties compose the object schema :
    - path [string] - (Optional)
    - is_path_regex [boolean] - (Optional)
    - permission [string] - (Optional)
    - require_auth [boolean] - (Optional)
    - conditions [string] - (Optional) required if permission == "conditions"
  - context_root [string] - (Optional) The root path to the application, often the name of the application. Can be any name, path or just a slash (“/”). The context root uniquely identifies the application within the enforcement point.
  - conditions [string] - (Optional) If access is conditional, the conditions that must evaluate to true to allow access to a resource. For example, to require the user must be authenticated and have either the role Admin or User
- visible [boolean] - (Optional) Indicates if the app is visible in the OneLogin portal.
- icon_url [string] - (Optional) A link to the apps icon url
- description [string] - (Optional) Freeform description of the app.

* Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin_apps.my_apps.**enforcement_point[0]**.object_property)

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- connector_id [integer] - ID of the connector to base the app from.
- created_at [string] - the date the app was created
- * enforcement_point [object] - For apps that connect to a OneLogin Access Enforcement Point the following enforcement_point object will be included with the app payload. The following properties compose the object schema:
  - token [string] - Can only be set on create. Access Gateway Token.
- id [integer] - Apps unique ID in OneLogin.
- updated_at [string] - the date the app was last updated

* Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin_apps.my_apps.**enforcement_point[0]**.object_property)

#### Import

apps resources can be imported using the `id` , e.g:

```shell
terraform import onelogin_apps.my_apps id
```

**Note**: In order for the import to work, the 'onelogin' terraform provider must be [properly installed](#provider-installation). Read more about Terraform import usage [here](https://www.terraform.io/docs/import/usage.html).

### onelogin_rules

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
- * source [object] - (Optional) Used for targeting custom rules based on a group of people, customers, accounts, or even a single user.. The following properties compose the object schema :
  - name [string] - (Optional) The name of the source
  - id [string] - (Optional) A unique id that represents the source of the event.
- name [string] - (Optional) The name of this rule
- description [string] - (Optional)

* Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin_rules.my_rules.**source[0]**.object_property)

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

* Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin_rules.my_rules.**source[0]**.object_property)

#### Import

rules resources can be imported using the `id` , e.g:

```shell
terraform import onelogin_rules.my_rules id
```

**Note**: In order for the import to work, the 'onelogin' terraform provider must be [properly installed](#provider-installation). Read more about Terraform import usage [here](https://www.terraform.io/docs/import/usage.html).

### onelogin_users

#### Example usage

```hcl
resource "onelogin_users" "my_users"{
}

```

#### Arguments Reference

The following arguments are supported:

- manager_ad_id [string] - (Optional) The ID of the user's manager in Active Directory.
- salt [string] - (Optional) The salt value used with the password_algorithm.
- password_changed_at [string] - (Optional)
- firstname [string] - (Optional) The user's first name.
- invitation_sent_at [string] - (Optional)
- password [string] - (Optional) The password to set for a user.
- username [string] - (Optional) A username for the user.
- status [integer] - (Optional)
- password_confirmation [string] - (Optional) Required if the password is being set.
- password_algorithm [string] - (Optional) Use this when importing a password that's already hashed. Prepend the salt value to the cleartext password value before SHA-256-encoding it
- phone [string] - (Optional) The E.164 format phone number for a user.
- invalid_login_attempts [integer] - (Optional)
- email [string] - (Optional) A valid email for the user.
- lastname [string] - (Optional) The user's last name.
- locked_until [string] - (Optional)
- id [integer] - (Optional)
- title [string] - (Optional) The user's job title.
- userprincipalname [string] - (Optional) The principle name of the user.
- member_of [string] - (Optional) The user's directory membership.
- role_ids [list of integers] - (Optional) A list of OneLogin Role IDs of the user
- state [integer] - (Optional)
- updated_at [string] - (Optional)
- trusted_idp_id [integer] - (Optional) The ID of the OneLogin Trusted IDP of the user.
- created_at [string] - (Optional)
- preferred_locale_code [string] - (Optional)
- group_id [integer] - (Optional) The ID of the Group in OneLogin that the user is assigned to.
- directory_id [integer] - (Optional) The ID of the OneLogin Directory of the user.
- distinguished_name [string] - (Optional) The distinguished name of the user.
- company [string] - (Optional) The user's company.
- manager_user_id [string] - (Optional) The OneLogin User ID for the user's manager.
- comment [string] - (Optional) Free text related to the user.
- samaccount_name [string] - (Optional) The user's Active Directory username.
- department [string] - (Optional) The user's department.
- external_id [string] - (Optional) The ID of the user in an external directory.
- activated_at [string] - (Optional)
- last_login [string] - (Optional)

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

#### Import

users resources can be imported using the `id` , e.g:

```shell
terraform import onelogin_users.my_users id
```

**Note**: In order for the import to work, the 'onelogin' terraform provider must be [properly installed](#provider-installation). Read more about Terraform import usage [here](https://www.terraform.io/docs/import/usage.html).

### onelogin_users

#### Example usage

```hcl
resource "onelogin_users" "my_users"{
}

```

#### Arguments Reference

The following arguments are supported:

- manager_ad_id [string] - (Optional) The ID of the user's manager in Active Directory.
- salt [string] - (Optional) The salt value used with the password_algorithm.
- password_changed_at [string] - (Optional)
- firstname [string] - (Optional) The user's first name.
- invitation_sent_at [string] - (Optional)
- password [string] - (Optional) The password to set for a user.
- username [string] - (Optional) A username for the user.
- status [integer] - (Optional)
- password_confirmation [string] - (Optional) Required if the password is being set.
- password_algorithm [string] - (Optional) Use this when importing a password that's already hashed. Prepend the salt value to the cleartext password value before SHA-256-encoding it
- phone [string] - (Optional) The E.164 format phone number for a user.
- invalid_login_attempts [integer] - (Optional)
- email [string] - (Optional) A valid email for the user.
- lastname [string] - (Optional) The user's last name.
- locked_until [string] - (Optional)
- id [integer] - (Optional)
- title [string] - (Optional) The user's job title.
- userprincipalname [string] - (Optional) The principle name of the user.
- member_of [string] - (Optional) The user's directory membership.
- role_ids [list of integers] - (Optional) A list of OneLogin Role IDs of the user
- state [integer] - (Optional)
- updated_at [string] - (Optional)
- trusted_idp_id [integer] - (Optional) The ID of the OneLogin Trusted IDP of the user.
- created_at [string] - (Optional)
- preferred_locale_code [string] - (Optional)
- group_id [integer] - (Optional) The ID of the Group in OneLogin that the user is assigned to.
- directory_id [integer] - (Optional) The ID of the OneLogin Directory of the user.
- distinguished_name [string] - (Optional) The distinguished name of the user.
- company [string] - (Optional) The user's company.
- manager_user_id [string] - (Optional) The OneLogin User ID for the user's manager.
- comment [string] - (Optional) Free text related to the user.
- samaccount_name [string] - (Optional) The user's Active Directory username.
- department [string] - (Optional) The user's department.
- external_id [string] - (Optional) The ID of the user in an external directory.
- activated_at [string] - (Optional)
- last_login [string] - (Optional)

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

#### Import

users resources can be imported using the `id` , e.g:

```hcl
terraform import onelogin_users.my_users id
```

**Note**: In order for the import to work, the 'onelogin' terraform provider must be [properly installed](#provider-installation). Read more about Terraform import usage [here](https://www.terraform.io/docs/import/usage.html).

## Data Sources (using resource id)

### onelogin_apps_instance

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

- tab_id [integer] - ID of the OneLogin portal tab that the app is assigned to.
- * provisioning [object] - Indicates if provisioning is enabled for this app. The following properties compose the object schema:
  - enabled [boolean]
- role_ids [list of integers] - List of Role IDs that are assigned to the app. On App Create or Update the entire array is replaced with the values provided.
- name [string] - The name of the app.
- auth_method [integer] - An ID indicating the type of app: - 0: Password - 1: OpenId - 2: SAML - 3: API - 4: Google - 6: Forms Based App - 7: WSFED - 8: OpenId Connect
- connector_id [integer] - ID of the connector to base the app from.
- allow_assumed_signin [boolean] - Indicates whether or not administrators can access the app as a user that they have assumed control over.
- created_at [string] - the date the app was created
- notes [string] - Freeform notes about the app.
- visible [boolean] - Indicates if the app is visible in the OneLogin portal.
- policy_id [integer] - The security policy assigned to the app.
- id [integer] - Apps unique ID in OneLogin.
- updated_at [string] - the date the app was last updated
- * parameters [object] - The parameters section contains parameterized attributes that have defined at the connector level as well as custom attributes that have been defined specifically for this app. Regardless of how they are defined, all parameters have the following attributes. Each parameter is an object with the key for the object being set as the parameters short name. The following properties compose the object schema:
  - include_in_saml_assertion [boolean] - When true, this parameter will be included in a SAML assertion payload.
  - label [string] - The can only be set when creating a new parameter. It can not be updated.
  - user_attribute_mappings [string] - A user attribute to map values from For custom attributes prefix the name of the attribute with `custom_attribute_`. e.g. To get the value for custom attribute `employee_id` use `custom_attribute_employee_id`.
  - user_attribute_macros [string] - When `user_attribute_mappings` is set to `_macro_` this macro will be used to assign the parameter value.
- description [string] - Freeform description of the app.
- * enforcement_point [object] - For apps that connect to a OneLogin Access Enforcement Point the following enforcement_point object will be included with the app payload. The following properties compose the object schema:
  - permissions [string] - Specify to always `allow`, `deny` access to resources, of if access is `conditional`.
  - case_sensitive [boolean] - The URL path evaluation is case insensitive by default. Resources hosted on web servers such as Apache, NGINX and Java EE are case sensitive paths. Web servers such as Microsoft IIS are not case-sensitive.
  - * session_expiry_fixed [object] - unit: - 0 = Seconds - 1 = Minutes - 2 = Hours value: - When Unit = 0 or 1 value must be 0-60 - When Unit = 2 value must be 0-24 The following properties compose the object schema:
    - unit [integer]
    - value [integer]
  - resources [list of objects] - Array of resource objects The following properties compose the object schema:
    - path [string]
    - permission [string]
    - is_path_regex [boolean]
    - require_auth [boolean]
    - conditions [string] - required if permission == "conditions"
  - vhost [string] - A comma-delimited list of one or more virtual hosts that map to applications assigned to the enforcement point. A VHOST may be a host name or an IP address. VHOST distinguish between applications that are at the same context root.
  - target [string] - A fully-qualified URL to the internal application including scheme, authority and path. The target host authority must be an IP address, not a hostname.
  - token [string] - Can only be set on create. Access Gateway Token.
  - landing_page [string] - The location within the context root to which the browser will be redirected for IdP-initiated single sign-on. For example, the landing page might be an index page in the context root such as index.html or default.aspx. The landing page cannot begin with a slash and must use valid URL characters.
  - use_target_host_header [boolean] - Use the target host header as opposed to the original gateway or upstream host header.
  - require_sitewide_authentication [boolean] - Require user authentication to access any resource protected by this enforcement point.
  - context_root [string] - The root path to the application, often the name of the application. Can be any name, path or just a slash (“/”). The context root uniquely identifies the application within the enforcement point.
  - * session_expiry_inactivity [object] - unit: - 0 = Seconds - 1 = Minutes - 2 = Hours value: - When Unit = 0 or 1 value must be 0-60 - When Unit = 2 value must be 0-24 The following properties compose the object schema:
    - unit [integer]
    - value [integer]
  - conditions [string] - If access is conditional, the conditions that must evaluate to true to allow access to a resource. For example, to require the user must be authenticated and have either the role Admin or User
- icon_url [string] - A link to the apps icon url

* Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin_apps_instance.my_apps_instance.**enforcement_point[0]**.object_property)

### onelogin_rules_instance

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
- * source [object] - Used for targeting custom rules based on a group of people, customers, accounts, or even a single user. The following properties compose the object schema:
  - name [string] - The name of the source
  - id [string] - A unique id that represents the source of the event.
- name [string] - The name of this rule
- description [string]

* Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin_rules_instance.my_rules_instance.**source[0]**.object_property)

### onelogin_users_instance

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

- invitation_sent_at [string]
- firstname [string] - The user's first name.
- salt [string] - The salt value used with the password_algorithm.
- password_changed_at [string]
- manager_ad_id [string] - The ID of the user's manager in Active Directory.
- phone [string] - The E.164 format phone number for a user.
- password_algorithm [string] - Use this when importing a password that's already hashed. Prepend the salt value to the cleartext password value before SHA-256-encoding it
- password_confirmation [string] - Required if the password is being set.
- password [string] - The password to set for a user.
- status [integer]
- username [string] - A username for the user.
- locked_until [string]
- lastname [string] - The user's last name.
- email [string] - A valid email for the user.
- invalid_login_attempts [integer]
- userprincipalname [string] - The principle name of the user.
- member_of [string] - The user's directory membership.
- title [string] - The user's job title.
- id [integer]
- updated_at [string]
- state [integer]
- role_ids [list of integers] - A list of OneLogin Role IDs of the user
- group_id [integer] - The ID of the Group in OneLogin that the user is assigned to.
- preferred_locale_code [string]
- directory_id [integer] - The ID of the OneLogin Directory of the user.
- created_at [string]
- trusted_idp_id [integer] - The ID of the OneLogin Trusted IDP of the user.
- company [string] - The user's company.
- distinguished_name [string] - The distinguished name of the user.
- activated_at [string]
- external_id [string] - The ID of the user in an external directory.
- last_login [string]
- comment [string] - Free text related to the user.
- samaccount_name [string] - The user's Active Directory username.
- department [string] - The user's department.
- manager_user_id [string] - The OneLogin User ID for the user's manager.

### onelogin_users_instance

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

- invitation_sent_at [string]
- firstname [string] - The user's first name.
- salt [string] - The salt value used with the password_algorithm.
- password_changed_at [string]
- manager_ad_id [string] - The ID of the user's manager in Active Directory.
- phone [string] - The E.164 format phone number for a user.
- password_algorithm [string] - Use this when importing a password that's already hashed. Prepend the salt value to the cleartext password value before SHA-256-encoding it
- password_confirmation [string] - Required if the password is being set.
- password [string] - The password to set for a user.
- status [integer]
- username [string] - A username for the user.
- locked_until [string]
- lastname [string] - The user's last name.
- email [string] - A valid email for the user.
- invalid_login_attempts [integer]
- userprincipalname [string] - The principle name of the user.
- member_of [string] - The user's directory membership.
- title [string] - The user's job title.
- id [integer]
- updated_at [string]
- state [integer]
- role_ids [list of integers] - A list of OneLogin Role IDs of the user
- group_id [integer] - The ID of the Group in OneLogin that the user is assigned to.
- preferred_locale_code [string]
- directory_id [integer] - The ID of the OneLogin Directory of the user.
- created_at [string]
- trusted_idp_id [integer] - The ID of the OneLogin Trusted IDP of the user.
- company [string] - The user's company.
- distinguished_name [string] - The distinguished name of the user.
- activated_at [string]
- external_id [string] - The ID of the user in an external directory.
- last_login [string]
- comment [string] - Free text related to the user.
- samaccount_name [string] - The user's Active Directory username.
- department [string] - The user's department.
- manager_user_id [string] - The OneLogin User ID for the user's manager.

## Data Sources (using filters)

### onelogin_api_authorizations (filters)

The api_authorizations data source allows you to retrieve an already existing api_authorizations resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_api_authorizations" "my_api_authorizations"{
    filter {
        name = "property name to filter by, see docs below for more info about available filter name options"
        values = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: description, name, id,
- values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- description [string] - Description of what the API does.
- name [string] - Name of the API.
- id [integer] - Auth server unique ID in Onelogin
- * configuration [object] - Authorization server configuration The following properties compose the object schema:
  - access_token_expiration_minutes [integer] - The number of minutes until access token expires. There is no maximum expiry limit.
  - refresh_token_expiration_minutes [integer] - The number of minutes until refresh token expires. There is no maximum expiry limit.
  - audiences [list of strings] - List of API endpoints that will be returned in Access Tokens.
  - resource_identifier [string] - Unique identifier for the API that the Authorization Server will issue Access Tokens for.

* Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin_api_authorizations.my_api_authorizations.**configuration[0]**.object_property)

### onelogin_api_authorizations_claims (filters)

The api_authorizations_claims data source allows you to retrieve an already existing api_authorizations_claims resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_api_authorizations_claims" "my_api_authorizations_claims"{
    filter {
        name = "property name to filter by, see docs below for more info about available filter name options"
        values = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: provisioned_entitlements, skip_if_blank, id, user_attribute_mappings, attribute_transformations, api_authorizations_id, user_attribute_macros, default_values, label,
- values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- provisioned_entitlements [boolean] - Relates to Rules/Entitlements. Not supported yet.
- skip_if_blank [boolean] - not used
- id [integer] - The unique ID of the claim.
- user_attribute_mappings [string] - A user attribute to map values from.
- attribute_transformations [string] - The type of transformation to perform on multi valued attributes.
- user_attribute_macros [string] - When `user_attribute_mappings` is set to `_macro_` this macro will be used to assign the claims value.
- values [list of strings] - Relates to Rules/Entitlements. Not supported yet.
- default_values [string] - Relates to Rules/Entitlements. Not supported yet.
- label [string] - The UI label for the claims.

### onelogin_api_authorizations_scopes (filters)

The api_authorizations_scopes data source allows you to retrieve an already existing api_authorizations_scopes resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_api_authorizations_scopes" "my_api_authorizations_scopes"{
    filter {
        name = "property name to filter by, see docs below for more info about available filter name options"
        values = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: description, name, id, api_authorizations_id,
- values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- description [string] - Description of what the API does.
- name [string] - Name of the API.
- id [integer] - Auth server unique ID in Onelogin
- * configuration [object] - Authorization server configuration The following properties compose the object schema:
  - access_token_expiration_minutes [integer] - The number of minutes until access token expires. There is no maximum expiry limit.
  - refresh_token_expiration_minutes [integer] - The number of minutes until refresh token expires. There is no maximum expiry limit.
  - audiences [list of strings] - List of API endpoints that will be returned in Access Tokens.
  - resource_identifier [string] - Unique identifier for the API that the Authorization Server will issue Access Tokens for.

* Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin_api_authorizations_scopes.my_api_authorizations_scopes.**configuration[0]**.object_property)

### onelogin_apps (filters)

The apps data source allows you to retrieve an already existing apps resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_apps" "my_apps"{
    filter {
        name = "property name to filter by, see docs below for more info about available filter name options"
        values = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: tab_id, name, auth_method, connector_id, allow_assumed_signin, created_at, notes, visible, policy_id, id, updated_at, description, icon_url,
- values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- tab_id [integer] - ID of the OneLogin portal tab that the app is assigned to.
- * provisioning [object] - Indicates if provisioning is enabled for this app. The following properties compose the object schema:
  - enabled [boolean]
- role_ids [list of integers] - List of Role IDs that are assigned to the app. On App Create or Update the entire array is replaced with the values provided.
- name [string] - The name of the app.
- auth_method [integer] - An ID indicating the type of app: - 0: Password - 1: OpenId - 2: SAML - 3: API - 4: Google - 6: Forms Based App - 7: WSFED - 8: OpenId Connect
- connector_id [integer] - ID of the connector to base the app from.
- allow_assumed_signin [boolean] - Indicates whether or not administrators can access the app as a user that they have assumed control over.
- created_at [string] - the date the app was created
- notes [string] - Freeform notes about the app.
- visible [boolean] - Indicates if the app is visible in the OneLogin portal.
- policy_id [integer] - The security policy assigned to the app.
- id [integer] - Apps unique ID in OneLogin.
- updated_at [string] - the date the app was last updated
- * parameters [object] - The parameters section contains parameterized attributes that have defined at the connector level as well as custom attributes that have been defined specifically for this app. Regardless of how they are defined, all parameters have the following attributes. Each parameter is an object with the key for the object being set as the parameters short name. The following properties compose the object schema:
  - include_in_saml_assertion [boolean] - When true, this parameter will be included in a SAML assertion payload.
  - label [string] - The can only be set when creating a new parameter. It can not be updated.
  - user_attribute_mappings [string] - A user attribute to map values from For custom attributes prefix the name of the attribute with `custom_attribute_`. e.g. To get the value for custom attribute `employee_id` use `custom_attribute_employee_id`.
  - user_attribute_macros [string] - When `user_attribute_mappings` is set to `_macro_` this macro will be used to assign the parameter value.
- description [string] - Freeform description of the app.
- * enforcement_point [object] - For apps that connect to a OneLogin Access Enforcement Point the following enforcement_point object will be included with the app payload. The following properties compose the object schema:
  - permissions [string] - Specify to always `allow`, `deny` access to resources, of if access is `conditional`.
  - case_sensitive [boolean] - The URL path evaluation is case insensitive by default. Resources hosted on web servers such as Apache, NGINX and Java EE are case sensitive paths. Web servers such as Microsoft IIS are not case-sensitive.
  - * session_expiry_fixed [object] - unit: - 0 = Seconds - 1 = Minutes - 2 = Hours value: - When Unit = 0 or 1 value must be 0-60 - When Unit = 2 value must be 0-24 The following properties compose the object schema:
    - unit [integer]
    - value [integer]
  - resources [list of objects] - Array of resource objects The following properties compose the object schema:
    - path [string]
    - permission [string]
    - is_path_regex [boolean]
    - require_auth [boolean]
    - conditions [string] - required if permission == "conditions"
  - vhost [string] - A comma-delimited list of one or more virtual hosts that map to applications assigned to the enforcement point. A VHOST may be a host name or an IP address. VHOST distinguish between applications that are at the same context root.
  - target [string] - A fully-qualified URL to the internal application including scheme, authority and path. The target host authority must be an IP address, not a hostname.
  - token [string] - Can only be set on create. Access Gateway Token.
  - landing_page [string] - The location within the context root to which the browser will be redirected for IdP-initiated single sign-on. For example, the landing page might be an index page in the context root such as index.html or default.aspx. The landing page cannot begin with a slash and must use valid URL characters.
  - use_target_host_header [boolean] - Use the target host header as opposed to the original gateway or upstream host header.
  - require_sitewide_authentication [boolean] - Require user authentication to access any resource protected by this enforcement point.
  - context_root [string] - The root path to the application, often the name of the application. Can be any name, path or just a slash (“/”). The context root uniquely identifies the application within the enforcement point.
  - * session_expiry_inactivity [object] - unit: - 0 = Seconds - 1 = Minutes - 2 = Hours value: - When Unit = 0 or 1 value must be 0-60 - When Unit = 2 value must be 0-24 The following properties compose the object schema:
    - unit [integer]
    - value [integer]
  - conditions [string] - If access is conditional, the conditions that must evaluate to true to allow access to a resource. For example, to require the user must be authenticated and have either the role Admin or User
- icon_url [string] - A link to the apps icon url

* Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin_apps.my_apps.**enforcement_point[0]**.object_property)

### onelogin_apps_actions (filters)

The apps_actions data source allows you to retrieve an already existing apps_actions resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_apps_actions" "my_apps_actions"{
    filter {
        name = "property name to filter by, see docs below for more info about available filter name options"
        values = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: apps_id, name, value,
- values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- name [string] - Name of the Action
- value [string] - The unique identifier of the action. This should be used when defining actions for a User Mapping.

### onelogin_apps_actions_values (filters)

The apps_actions_values data source allows you to retrieve an already existing apps_actions_values resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_apps_actions_values" "my_apps_actions_values"{
    filter {
        name = "property name to filter by, see docs below for more info about available filter name options"
        values = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: actions_id, apps_id, name, value,
- values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- name [string] - Name of the Action
- value [string] - The unique identifier of the action. This should be used when defining actions for a User Mapping.

### onelogin_apps_conditions (filters)

The apps_conditions data source allows you to retrieve an already existing apps_conditions resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_apps_conditions" "my_apps_conditions"{
    filter {
        name = "property name to filter by, see docs below for more info about available filter name options"
        values = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: value, apps_id, name,
- values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- value [string] - The unique identifier of the condition. This should be used when defining conditions for a rule.
- name [string] - Name of the rule condition

### onelogin_apps_conditions_operators (filters)

The apps_conditions_operators data source allows you to retrieve an already existing apps_conditions_operators resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_apps_conditions_operators" "my_apps_conditions_operators"{
    filter {
        name = "property name to filter by, see docs below for more info about available filter name options"
        values = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: name, apps_id, value, conditions_id,
- values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- name [string] - Name of the operator
- value [string] - The condition operator value to use when creating or updating rules.

### onelogin_apps_conditions_values (filters)

The apps_conditions_values data source allows you to retrieve an already existing apps_conditions_values resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_apps_conditions_values" "my_apps_conditions_values"{
    filter {
        name = "property name to filter by, see docs below for more info about available filter name options"
        values = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: value, apps_id, name, conditions_id,
- values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- value [string] - The unique identifier of the condition. This should be used when defining conditions for a rule.
- name [string] - Name of the rule condition

### onelogin_apps_rules (filters)

The apps_rules data source allows you to retrieve an already existing apps_rules resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_apps_rules" "my_apps_rules"{
    filter {
        name = "property name to filter by, see docs below for more info about available filter name options"
        values = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: match, name, apps_id, enabled, id, position,
- values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- match [string] - Indicates how conditions should be matched.
- conditions [list of objects] - An array of conditions that the user must meet in order for the rule to be applied. The following properties compose the object schema:
  - source [string] - source field to check.
  - operator [string] - A valid operator for the selected condition source
  - value [string] - A plain text string or valid value for the selected condition source
- name [string] - Rule Name
- enabled [boolean] - Indicates if the rule is enabled or not.
- id [integer] - App Rule ID
- actions [list of objects] The following properties compose the object schema:
  - value [list of strings] - Only applicable to provisioned and set_* actions. Items in the array will be a plain text string or valid value for the selected action.
  - action [string] - The action to apply
- position [integer] - Indicates the order of the rule. When `null` this will default to last position.

### onelogin_apps_users (filters)

The apps_users data source allows you to retrieve an already existing apps_users resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_apps_users" "my_apps_users"{
    filter {
        name = "property name to filter by, see docs below for more info about available filter name options"
        values = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: invitation_sent_at, firstname, salt, password_changed_at, manager_ad_id, phone, password_algorithm, password_confirmation, password, status, username, locked_until, lastname, email, apps_id, invalid_login_attempts, userprincipalname, member_of, title, id, updated_at, state, group_id, preferred_locale_code, directory_id, created_at, trusted_idp_id, company, distinguished_name, activated_at, external_id, last_login, comment, samaccount_name, department, manager_user_id,
- values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- invitation_sent_at [string]
- firstname [string] - The user's first name.
- salt [string] - The salt value used with the password_algorithm.
- password_changed_at [string]
- manager_ad_id [string] - The ID of the user's manager in Active Directory.
- phone [string] - The E.164 format phone number for a user.
- password_algorithm [string] - Use this when importing a password that's already hashed. Prepend the salt value to the cleartext password value before SHA-256-encoding it
- password_confirmation [string] - Required if the password is being set.
- password [string] - The password to set for a user.
- status [integer]
- username [string] - A username for the user.
- locked_until [string]
- lastname [string] - The user's last name.
- email [string] - A valid email for the user.
- invalid_login_attempts [integer]
- userprincipalname [string] - The principle name of the user.
- member_of [string] - The user's directory membership.
- title [string] - The user's job title.
- id [integer]
- updated_at [string]
- state [integer]
- role_ids [list of integers] - A list of OneLogin Role IDs of the user
- group_id [integer] - The ID of the Group in OneLogin that the user is assigned to.
- preferred_locale_code [string]
- directory_id [integer] - The ID of the OneLogin Directory of the user.
- created_at [string]
- trusted_idp_id [integer] - The ID of the OneLogin Trusted IDP of the user.
- company [string] - The user's company.
- distinguished_name [string] - The distinguished name of the user.
- activated_at [string]
- external_id [string] - The ID of the user in an external directory.
- last_login [string]
- comment [string] - Free text related to the user.
- samaccount_name [string] - The user's Active Directory username.
- department [string] - The user's department.
- manager_user_id [string] - The OneLogin User ID for the user's manager.

### onelogin_brands (filters)

The brands data source allows you to retrieve an already existing brands resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_brands" "my_brands"{
    filter {
        name = "property name to filter by, see docs below for more info about available filter name options"
        values = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: enabled, id, name,
- values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- enabled [boolean] - Indicates if the brand is enabled or not.
- id [integer] - Brand’s unique ID in OneLogin.
- name [string] - Brand name for humans. This isn’t related to subdomains.

### onelogin_brands_apps (filters)

The brands_apps data source allows you to retrieve an already existing brands_apps resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_brands_apps" "my_brands_apps"{
    filter {
        name = "property name to filter by, see docs below for more info about available filter name options"
        values = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: visible, connector_id, id, updated_at, auth_method, created_at, description, brands_id, auth_method_description, name,
- values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- visible [boolean]
- connector_id [integer]
- id [integer]
- updated_at [string]
- auth_method [integer]
- created_at [string]
- description [string]
- auth_method_description [string]
- name [string]

### onelogin_brands_templates (filters)

The brands_templates data source allows you to retrieve an already existing brands_templates resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_brands_templates" "my_brands_templates"{
    filter {
        name = "property name to filter by, see docs below for more info about available filter name options"
        values = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: name, enabled, id, brands_id,
- values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- name [string] - name of message template
- enabled [boolean] - indicator if template is enabled
- id [integer] - template ID

### onelogin_mappings (filters)

The mappings data source allows you to retrieve an already existing mappings resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_mappings" "my_mappings"{
    filter {
        name = "property name to filter by, see docs below for more info about available filter name options"
        values = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: match, name, id, enabled, position,
- values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- match [string] - Indicates how conditions should be matched.
- name [string] - The name of the mapping.
- actions [list of objects] - An array of actions that will be applied to the users that are matched by the conditions. The following properties compose the object schema:
  - value [list of strings] - Only applicable to provisioned and set_* actions. Items in the array will be a plain text string or valid value for the selected action.
  - action [string] - The action to apply
- id [integer]
- enabled [boolean] - Indicates if the mapping is enabled or not.
- conditions [list of objects] - An array of conditions that the user must meet in order for the mapping to be applied. The following properties compose the object schema:
  - source [string] - source field to check.
  - operator [string] - A valid operator for the selected condition source
  - value [string] - A plain text string or valid value for the selected condition source
- position [integer] - Indicates the order of the mapping. When `null` this will default to last position.

### onelogin_privileges (filters)

The privileges data source allows you to retrieve an already existing privileges resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_privileges" "my_privileges"{
    filter {
        name = "property name to filter by, see docs below for more info about available filter name options"
        values = ["filter value"]
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
- * privilege [object] The following properties compose the object schema:
  - statement [list of objects] The following properties compose the object schema:
    - scope [list of strings] - Target the privileged action against specific resources with the scope. The scope pattern is the class of object used by the Action, followed by an ID that represents a resource in OneLogin. e.g. apps/1234, where apps is the class and 1234 is the ID of an app. The wildcard * is supported and indicates that all resources of the class type declared, in the Action, are in scope. The Action and Scope classes must match. However, there is an exception, a scope of roles/{role_id} can be combined with Actions on the user or app class. The exception allows you to target groups of users or apps with specific actions.
    - action [list of strings] - An array of strings that represent actions within OneLogin. Actions are prefixed with the class of object they are related to and followed by a specific action for the given class. e.g. users:List, where the class is users and the specific action is List. Don’t mix classes within an Action array. To create a privilege that includes multiple different classes, create multiple statements. A wildcard * that includes all actions is supported. Use wildcards to create a Super User privilege.
    - effect [string] - Set to “Allow.” By default, all actions are denied, this Statement allows the listed actions to be executed.
  - version [string]
- description [string]
- name [string]

* Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin_privileges.my_privileges.**privilege[0]**.object_property)

### onelogin_roles (filters)

The roles data source allows you to retrieve an already existing roles resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_roles" "my_roles"{
    filter {
        name = "property name to filter by, see docs below for more info about available filter name options"
        values = ["filter value"]
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

- users [list of integers] - array of user IDs
- id [integer] - Role ID
- name [string] - Role Name
- admins [list of integers]
- apps [list of integers] - array of app IDs

### onelogin_roles-admins (filters)

The roles-admins data source allows you to retrieve an already existing roles-admins resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_roles-admins" "my_roles-admins"{
    filter {
        name = "property name to filter by, see docs below for more info about available filter name options"
        values = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: invitation_sent_at, firstname, salt, password_changed_at, manager_ad_id, phone, password_algorithm, password_confirmation, password, status, username, locked_until, lastname, email, invalid_login_attempts, userprincipalname, member_of, title, id, updated_at, state, group_id, preferred_locale_code, directory_id, created_at, trusted_idp_id, company, distinguished_name, activated_at, external_id, last_login, comment, samaccount_name, department, roles-id, manager_user_id,
- values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- invitation_sent_at [string]
- firstname [string] - The user's first name.
- salt [string] - The salt value used with the password_algorithm.
- password_changed_at [string]
- manager_ad_id [string] - The ID of the user's manager in Active Directory.
- phone [string] - The E.164 format phone number for a user.
- password_algorithm [string] - Use this when importing a password that's already hashed. Prepend the salt value to the cleartext password value before SHA-256-encoding it
- password_confirmation [string] - Required if the password is being set.
- password [string] - The password to set for a user.
- status [integer]
- username [string] - A username for the user.
- locked_until [string]
- lastname [string] - The user's last name.
- email [string] - A valid email for the user.
- invalid_login_attempts [integer]
- userprincipalname [string] - The principle name of the user.
- member_of [string] - The user's directory membership.
- title [string] - The user's job title.
- id [integer]
- updated_at [string]
- state [integer]
- role_ids [list of integers] - A list of OneLogin Role IDs of the user
- group_id [integer] - The ID of the Group in OneLogin that the user is assigned to.
- preferred_locale_code [string]
- directory_id [integer] - The ID of the OneLogin Directory of the user.
- created_at [string]
- trusted_idp_id [integer] - The ID of the OneLogin Trusted IDP of the user.
- company [string] - The user's company.
- distinguished_name [string] - The distinguished name of the user.
- activated_at [string]
- external_id [string] - The ID of the user in an external directory.
- last_login [string]
- comment [string] - Free text related to the user.
- samaccount_name [string] - The user's Active Directory username.
- department [string] - The user's department.
- manager_user_id [string] - The OneLogin User ID for the user's manager.

### onelogin_roles-apps (filters)

The roles-apps data source allows you to retrieve an already existing roles-apps resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_roles-apps" "my_roles-apps"{
    filter {
        name = "property name to filter by, see docs below for more info about available filter name options"
        values = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: icon_url, id, name, roles-id,
- values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- icon_url [string] - url of Icon
- id [integer] - app id
- name [string] - app name

### onelogin_roles-users (filters)

The roles-users data source allows you to retrieve an already existing roles-users resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_roles-users" "my_roles-users"{
    filter {
        name = "property name to filter by, see docs below for more info about available filter name options"
        values = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: invitation_sent_at, firstname, salt, password_changed_at, manager_ad_id, phone, password_algorithm, password_confirmation, password, status, username, locked_until, lastname, email, invalid_login_attempts, userprincipalname, member_of, title, id, updated_at, state, group_id, preferred_locale_code, directory_id, created_at, trusted_idp_id, company, distinguished_name, activated_at, external_id, last_login, comment, samaccount_name, department, roles-id, manager_user_id,
- values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- invitation_sent_at [string]
- firstname [string] - The user's first name.
- salt [string] - The salt value used with the password_algorithm.
- password_changed_at [string]
- manager_ad_id [string] - The ID of the user's manager in Active Directory.
- phone [string] - The E.164 format phone number for a user.
- password_algorithm [string] - Use this when importing a password that's already hashed. Prepend the salt value to the cleartext password value before SHA-256-encoding it
- password_confirmation [string] - Required if the password is being set.
- password [string] - The password to set for a user.
- status [integer]
- username [string] - A username for the user.
- locked_until [string]
- lastname [string] - The user's last name.
- email [string] - A valid email for the user.
- invalid_login_attempts [integer]
- userprincipalname [string] - The principle name of the user.
- member_of [string] - The user's directory membership.
- title [string] - The user's job title.
- id [integer]
- updated_at [string]
- state [integer]
- role_ids [list of integers] - A list of OneLogin Role IDs of the user
- group_id [integer] - The ID of the Group in OneLogin that the user is assigned to.
- preferred_locale_code [string]
- directory_id [integer] - The ID of the OneLogin Directory of the user.
- created_at [string]
- trusted_idp_id [integer] - The ID of the OneLogin Trusted IDP of the user.
- company [string] - The user's company.
- distinguished_name [string] - The distinguished name of the user.
- activated_at [string]
- external_id [string] - The ID of the user in an external directory.
- last_login [string]
- comment [string] - Free text related to the user.
- samaccount_name [string] - The user's Active Directory username.
- department [string] - The user's department.
- manager_user_id [string] - The OneLogin User ID for the user's manager.

### onelogin_rules (filters)

The rules data source allows you to retrieve an already existing rules resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_rules" "my_rules"{
    filter {
        name = "property name to filter by, see docs below for more info about available filter name options"
        values = ["filter value"]
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
- * source [object] - Used for targeting custom rules based on a group of people, customers, accounts, or even a single user. The following properties compose the object schema:
  - name [string] - The name of the source
  - id [string] - A unique id that represents the source of the event.
- name [string] - The name of this rule
- description [string]

* Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin_rules.my_rules.**source[0]**.object_property)

### onelogin_users (filters)

The users data source allows you to retrieve an already existing users resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_users" "my_users"{
    filter {
        name = "property name to filter by, see docs below for more info about available filter name options"
        values = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: invitation_sent_at, firstname, salt, password_changed_at, manager_ad_id, phone, password_algorithm, password_confirmation, password, status, username, locked_until, lastname, email, invalid_login_attempts, userprincipalname, member_of, title, id, updated_at, state, group_id, preferred_locale_code, directory_id, created_at, trusted_idp_id, company, distinguished_name, activated_at, external_id, last_login, comment, samaccount_name, department, manager_user_id,
- values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- invitation_sent_at [string]
- firstname [string] - The user's first name.
- salt [string] - The salt value used with the password_algorithm.
- password_changed_at [string]
- manager_ad_id [string] - The ID of the user's manager in Active Directory.
- phone [string] - The E.164 format phone number for a user.
- password_algorithm [string] - Use this when importing a password that's already hashed. Prepend the salt value to the cleartext password value before SHA-256-encoding it
- password_confirmation [string] - Required if the password is being set.
- password [string] - The password to set for a user.
- status [integer]
- username [string] - A username for the user.
- locked_until [string]
- lastname [string] - The user's last name.
- email [string] - A valid email for the user.
- invalid_login_attempts [integer]
- userprincipalname [string] - The principle name of the user.
- member_of [string] - The user's directory membership.
- title [string] - The user's job title.
- id [integer]
- updated_at [string]
- state [integer]
- role_ids [list of integers] - A list of OneLogin Role IDs of the user
- group_id [integer] - The ID of the Group in OneLogin that the user is assigned to.
- preferred_locale_code [string]
- directory_id [integer] - The ID of the OneLogin Directory of the user.
- created_at [string]
- trusted_idp_id [integer] - The ID of the OneLogin Trusted IDP of the user.
- company [string] - The user's company.
- distinguished_name [string] - The distinguished name of the user.
- activated_at [string]
- external_id [string] - The ID of the user in an external directory.
- last_login [string]
- comment [string] - Free text related to the user.
- samaccount_name [string] - The user's Active Directory username.
- department [string] - The user's department.
- manager_user_id [string] - The OneLogin User ID for the user's manager.

### onelogin_users (filters)

The users data source allows you to retrieve an already existing users resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_users" "my_users"{
    filter {
        name = "property name to filter by, see docs below for more info about available filter name options"
        values = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: invitation_sent_at, firstname, salt, password_changed_at, manager_ad_id, phone, password_algorithm, password_confirmation, password, status, username, locked_until, lastname, email, invalid_login_attempts, userprincipalname, member_of, title, id, updated_at, state, group_id, preferred_locale_code, directory_id, created_at, trusted_idp_id, company, distinguished_name, activated_at, external_id, last_login, comment, samaccount_name, department, manager_user_id,
- values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- invitation_sent_at [string]
- firstname [string] - The user's first name.
- salt [string] - The salt value used with the password_algorithm.
- password_changed_at [string]
- manager_ad_id [string] - The ID of the user's manager in Active Directory.
- phone [string] - The E.164 format phone number for a user.
- password_algorithm [string] - Use this when importing a password that's already hashed. Prepend the salt value to the cleartext password value before SHA-256-encoding it
- password_confirmation [string] - Required if the password is being set.
- password [string] - The password to set for a user.
- status [integer]
- username [string] - A username for the user.
- locked_until [string]
- lastname [string] - The user's last name.
- email [string] - A valid email for the user.
- invalid_login_attempts [integer]
- userprincipalname [string] - The principle name of the user.
- member_of [string] - The user's directory membership.
- title [string] - The user's job title.
- id [integer]
- updated_at [string]
- state [integer]
- role_ids [list of integers] - A list of OneLogin Role IDs of the user
- group_id [integer] - The ID of the Group in OneLogin that the user is assigned to.
- preferred_locale_code [string]
- directory_id [integer] - The ID of the OneLogin Directory of the user.
- created_at [string]
- trusted_idp_id [integer] - The ID of the OneLogin Trusted IDP of the user.
- company [string] - The user's company.
- distinguished_name [string] - The distinguished name of the user.
- activated_at [string]
- external_id [string] - The ID of the user in an external directory.
- last_login [string]
- comment [string] - Free text related to the user.
- samaccount_name [string] - The user's Active Directory username.
- department [string] - The user's department.
- manager_user_id [string] - The OneLogin User ID for the user's manager.

### onelogin_users_apps (filters)

The users_apps data source allows you to retrieve an already existing users_apps resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_users_apps" "my_users_apps"{
    filter {
        name = "property name to filter by, see docs below for more info about available filter name options"
        values = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: login_id, extension, name, icon_url, id, provisioning_state, provisioning_enabled, users_id, provisioning_status,
- values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- login_id [integer] - Unqiue identifier for this user and app combination.
- extension [boolean] - Boolean that indicates if the OneLogin browser extension is required to launch this app.
- name [string] - The name of the app.
- icon_url [string] - A url for the icon that represents the app in the OneLogin portal
- id [integer] - The App ID
- provisioning_state [string] - If provisioning is enabled this indicates the state of provisioning for the given user.
- provisioning_enabled [boolean] - Indicates if provisioning is enabled for this app.
- provisioning_status [string]

### onelogin_users_apps (filters)

The users_apps data source allows you to retrieve an already existing users_apps resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_users_apps" "my_users_apps"{
    filter {
        name = "property name to filter by, see docs below for more info about available filter name options"
        values = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: login_id, extension, name, icon_url, id, provisioning_state, provisioning_enabled, users_id, provisioning_status,
- values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- login_id [integer] - Unqiue identifier for this user and app combination.
- extension [boolean] - Boolean that indicates if the OneLogin browser extension is required to launch this app.
- name [string] - The name of the app.
- icon_url [string] - A url for the icon that represents the app in the OneLogin portal
- id [integer] - The App ID
- provisioning_state [string] - If provisioning is enabled this indicates the state of provisioning for the given user.
- provisioning_enabled [boolean] - Indicates if provisioning is enabled for this app.
- provisioning_status [string]

### onelogin_users_devices (filters)

The users_devices data source allows you to retrieve an already existing users_devices resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

```hcl
data "onelogin_users_devices" "my_users_devices"{
    filter {
        name = "property name to filter by, see docs below for more info about available filter name options"
        values = ["filter value"]
    }
}
```

#### Arguments Reference

The following arguments are supported:

- filter - (Required) Object containing two properties.

- name [string]: the name should match one of the properties to filter by. The following property names are supported: device_id, default, user_display_name, auth_factor_name, users_id, type_display_name,
- values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

- device_id [string] - MFA device identifier.
- default [boolean] - true = is user’s default MFA device for OneLogin.
- user_display_name [string] - Authentication factor display name assigned by users when they register the device.
- auth_factor_name [string] - Authentication factor name, as it appears to administrators in OneLogin.
- type_display_name [string] - Authentication factor display name as it appears to users upon initial registration, as defined by admins at Settings > Authentication Factors.
