# terraform provider onelogin

This guide lists the configuration for 'onelogin' Terraform provider resources that can be managed using [Terraform v0.12](https://www.hashicorp.com/blog/announcing-terraform-0-12/). 


## Provider Installation

In order to provision 'onelogin' Terraform resources, you need to first install the 'onelogin' Terraform plugin by running the following command (you must be running Terraform >= 0.12):

<pre>$ export PROVIDER_NAME=onelogin && curl -fsSL https://raw.githubusercontent.com/dikhan/terraform-provider-openapi/master/scripts/install.sh | bash -s -- --provider-name $PROVIDER_NAME  
[INFO] Downloading https://github.com/dikhan/terraform-provider-openapi/v3/releases/download/v3.0.0/terraform-provider-openapi_3.0.0_darwin_amd64.tar.gz in temporally folder /var/folders/n_/1lrwb99s7f50xmn9jpmfnddh0000gp/T/tmp.Xv1AkIZh...  
[INFO] Extracting terraform-provider-openapi from terraform-provider-openapi_0.29.4_darwin_amd64.tar.gz...  
[INFO] Cleaning up tmp dir created for installation purposes: /var/folders/n_/1lrwb99s7f50xmn9jpmfnddh0000gp/T/tmp.Xv1AkIZh  
[INFO] Terraform provider 'terraform-provider-onelogin' successfully installed at: '~/.terraform.d/plugins'!</pre>

<span>You can then start running the Terraform provider:</span>

<pre dir="ltr">$ export OTF_VAR_onelogin_PLUGIN_CONFIGURATION_FILE="https://api.service.com/openapi.yaml"  


```shell
export OTF_VAR_onelogin_PLUGIN_CONFIGURATION_FILE="[https://api.service.com/openapi.yaml](https://raw.githubusercontent.com/onelogin/terraform-provider-onelogin/develop/swag-api.yml)"

➜ ~ terraform init && terraform plan
```

**Note:** As of Terraform >= 0.13 each Terraform module must declare which providers it requires, so that Terraform can install and use them. If you are using Terraform >= 0.13, copy into your .tf file the following snippet already populated with the provider configuration:

<pre dir="ltr">terraform {
  required_providers {
    onelogin = {
      source  = "onelogin/onelogin"
      version = ">= 0.3.8"

    }
  }
}
</pre>

## Provider Configuration

#### Example Usage

<pre><span>provider</span> <span>"onelogin"</span> {
 <span>apikey_auth</span> = <span>"..."</span>
 <span>content_type</span> = <span>"..."</span>
<span>}</span>
</pre>

## Provider Resources

### onelogin_apps

#### Example usage

<pre><span>resource</span> <span>"onelogin_apps" "my_apps"</span>{
<span>}</span>
</pre>

#### Arguments Reference

The following arguments are supported:

*   <span class="wysiwyg-color-red">*</span> provisioning [object] - (Optional) Indicates if provisioning is enabled for this app.. The following properties compose the object schema :
    *   enabled [boolean] - (Optional)
*   tab_id [integer] - (Optional) ID of the OneLogin portal tab that the app is assigned to.
*   name [string] - (Optional) The name of the app.
*   role_ids [list of integers] - (Optional) List of Role IDs that are assigned to the app. On App Create or Update the entire array is replaced with the values provided.
*   allow_assumed_signin [boolean] - (Optional) Indicates whether or not administrators can access the app as a user that they have assumed control over.
*   auth_method [integer] - (Optional) An ID indicating the type of app: - 0: Password - 1: OpenId - 2: SAML - 3: API - 4: Google - 6: Forms Based App - 7: WSFED - 8: OpenId Connect
*   policy_id [integer] - (Optional) The security policy assigned to the app.
*   <span class="wysiwyg-color-red">*</span> parameters [object] - (Optional) The parameters section contains parameterized attributes that have defined at the connector level as well as custom attributes that have been defined specifically for this app. Regardless of how they are defined, all parameters have the following attributes. Each parameter is an object with the key for the object being set as the parameters short name.. The following properties compose the object schema :
    *   include_in_saml_assertion [boolean] - (Optional) When true, this parameter will be included in a SAML assertion payload.
    *   label [string] - (Optional) The can only be set when creating a new parameter. It can not be updated.
    *   user_attribute_mappings [string] - (Optional) A user attribute to map values from For custom attributes prefix the name of the attribute with `custom_attribute_`. e.g. To get the value for custom attribute `employee_id` use `custom_attribute_employee_id`.
    *   user_attribute_macros [string] - (Optional) When `user_attribute_mappings` is set to `_macro_` this macro will be used to assign the parameter value.
*   notes [string] - (Optional) Freeform notes about the app.
*   updated_at [string] - (Optional) the date the app was last updated
*   <span class="wysiwyg-color-red">*</span> enforcement_point [object] - (Optional) For apps that connect to a OneLogin Access Enforcement Point the following enforcement_point object will be included with the app payload.. The following properties compose the object schema :
    *   <span class="wysiwyg-color-red">*</span> session_expiry_inactivity [object] - (Optional) unit: - 0 = Seconds - 1 = Minutes - 2 = Hours value: - When Unit = 0 or 1 value must be 0-60 - When Unit = 2 value must be 0-24\. The following properties compose the object schema :
        *   unit [integer] - (Optional)
        *   value [integer] - (Optional)
    *   case_sensitive [boolean] - (Optional) The URL path evaluation is case insensitive by default. Resources hosted on web servers such as Apache, NGINX and Java EE are case sensitive paths. Web servers such as Microsoft IIS are not case-sensitive.
    *   permissions [string] - (Optional) Specify to always `allow`, `deny` access to resources, of if access is `conditional`.
    *   target [string] - (Optional) A fully-qualified URL to the internal application including scheme, authority and path. The target host authority must be an IP address, not a hostname.
    *   vhost [string] - (Optional) A comma-delimited list of one or more virtual hosts that map to applications assigned to the enforcement point. A VHOST may be a host name or an IP address. VHOST distinguish between applications that are at the same context root.
    *   require_sitewide_authentication [boolean] - (Optional) Require user authentication to access any resource protected by this enforcement point.
    *   <span class="wysiwyg-color-red">*</span> session_expiry_fixed [object] - (Optional) unit: - 0 = Seconds - 1 = Minutes - 2 = Hours value: - When Unit = 0 or 1 value must be 0-60 - When Unit = 2 value must be 0-24\. The following properties compose the object schema :
        *   unit [integer] - (Optional)
        *   value [integer] - (Optional)
    *   use_target_host_header [boolean] - (Optional) Use the target host header as opposed to the original gateway or upstream host header.
    *   landing_page [string] - (Optional) The location within the context root to which the browser will be redirected for IdP-initiated single sign-on. For example, the landing page might be an index page in the context root such as index.html or default.aspx. The landing page cannot begin with a slash and must use valid URL characters.
    *   resources [list of objects] - (Optional) Array of resource objects. The following properties compose the object schema :
        *   path [string] - (Optional)
        *   is_path_regex [boolean] - (Optional)
        *   permission [string] - (Optional)
        *   require_auth [boolean] - (Optional)
        *   conditions [string] - (Optional) required if permission == "conditions"
    *   context_root [string] - (Optional) The root path to the application, often the name of the application. Can be any name, path or just a slash (“/”). The context root uniquely identifies the application within the enforcement point.
    *   conditions [string] - (Optional) If access is conditional, the conditions that must evaluate to true to allow access to a resource. For example, to require the user must be authenticated and have either the role Admin or User
*   visible [boolean] - (Optional) Indicates if the app is visible in the OneLogin portal.
*   created_at [string] - (Optional) the date the app was created
*   connector_id [integer] - (Optional) ID of the connector to base the app from.
*   icon_url [string] - (Optional) A link to the apps icon url
*   description [string] - (Optional) Freeform description of the app.

<span class="wysiwyg-color-red">*</span> Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin_apps.my_apps.**enforcement_point[0]**.object_property)

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   <span class="wysiwyg-color-red">*</span> enforcement_point [object] - For apps that connect to a OneLogin Access Enforcement Point the following enforcement_point object will be included with the app payload. The following properties compose the object schema:
    *   token [string] - Can only be set on create. Access Gateway Token.
*   id [integer] - Apps unique ID in OneLogin.

<span class="wysiwyg-color-red">*</span> Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin_apps.my_apps.**enforcement_point[0]**.object_property)

#### Import

apps resources can be imported using the `id` , e.g:

<pre dir="ltr">$ terraform import onelogin_apps.my_apps id</pre>

**Note**: In order for the import to work, the 'onelogin' terraform provider must be [properly installed](#provider_installation). Read more about Terraform import usage [here](https://www.terraform.io/docs/import/usage.html).

### onelogin_apps_rules

#### Example usage

<pre><span>resource</span> <span>"onelogin_apps_rules" "my_apps_rules"</span>{
    <span>apps_id</span> = <span>"apps_id"</span>
<span>}</span>
</pre>

#### Arguments Reference

The following arguments are supported:

*   apps_id [string] - (Required) The apps_id that this resource belongs to
*   match [string] - (Optional) Indicates how conditions should be matched.
*   actions [list of objects] - (Optional) . The following properties compose the object schema :
    *   value [list of strings] - (Optional) Only applicable to provisioned and set_* actions. Items in the array will be a plain text string or valid value for the selected action.
    *   action [string] - (Optional) The action to apply
*   name [string] - (Optional) Rule Name
*   conditions [list of objects] - (Optional) An array of conditions that the user must meet in order for the rule to be applied.. The following properties compose the object schema :
    *   operator [string] - (Optional) A valid operator for the selected condition source
    *   source [string] - (Optional) source field to check.
    *   value [string] - (Optional) A plain text string or valid value for the selected condition source
*   enabled [boolean] - (Optional) Indicates if the rule is enabled or not.
*   id [integer] - (Optional) App Rule ID
*   position [integer] - (Optional) Indicates the order of the rule. When `null` this will default to last position.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

#### Import

apps_rules resources can be imported using the `id` . This is a sub-resource so the parent resource IDs (`[apps_id]`) are required to be able to retrieve an instance of this resource, e.g:

<pre dir="ltr">$ terraform import onelogin_apps_rules.my_apps_rules apps_id/apps_rules_id</pre>

**Note**: In order for the import to work, the 'onelogin' terraform provider must be [properly installed](#provider_installation). Read more about Terraform import usage [here](https://www.terraform.io/docs/import/usage.html).

### onelogin_auth_servers

#### Example usage

<pre><span>resource</span> <span>"onelogin_auth_servers" "my_auth_servers"</span>{
    <span>configuration</span> <span>{</span>

    <span>audiences</span> = <span>["audiences1", "audiences2"]</span>

    <span>resource_identifier</span> = <span>"resource_identifier"</span>
            <span>}</span>
    <span>description</span> = <span>"description"</span>
    <span>name</span> = <span>"name"</span>
<span>}</span>
</pre>

#### Arguments Reference

The following arguments are supported:

*   <span class="wysiwyg-color-red">*</span> configuration [object] - (Required) Authorization server configuration. The following properties compose the object schema :
    *   access_token_expiration_minutes [integer] - (Optional) The number of minutes until access token expires. There is no maximum expiry limit.
    *   audiences [list of strings] - (Required) List of API endpoints that will be returned in Access Tokens.
    *   refresh_token_expiration_minutes [integer] - (Optional) The number of minutes until refresh token expires. There is no maximum expiry limit.
    *   resource_identifier [string] - (Required) Unique identifier for the API that the Authorization Server will issue Access Tokens for.
*   description [string] - (Required) Description of what the API does.
*   name [string] - (Required) Name of the API.

<span class="wysiwyg-color-red">*</span> Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin_auth_servers.my_auth_servers.**configuration[0]**.object_property)

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   id [integer] - Auth server unique ID in Onelogin

<span class="wysiwyg-color-red">*</span> Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin_auth_servers.my_auth_servers.**configuration[0]**.object_property)

#### Import

auth_servers resources can be imported using the `id` , e.g:

<pre dir="ltr">$ terraform import onelogin_auth_servers.my_auth_servers id</pre>

**Note**: In order for the import to work, the 'onelogin' terraform provider must be [properly installed](#provider_installation). Read more about Terraform import usage [here](https://www.terraform.io/docs/import/usage.html).

### onelogin_brands

#### Example usage

<pre><span>resource</span> <span>"onelogin_brands" "my_brands"</span>{
    <span>custom_accent_color</span> = <span>"custom_accent_color"</span>
    <span>custom_label_text_for_login_screen</span> = <span>"custom_label_text_for_login_screen"</span>
    <span>hide_onelogin_footer</span> = <span>true</span>
    <span>login_instruction</span> = <span>"login_instruction"</span>
    <span>custom_color</span> = <span>"custom_color"</span>
    <span>background</span> <span>{</span>

    <span>urls</span> <span>{</span>

    <span>branding</span> = <span>"branding"</span>

    <span>login</span> = <span>"login"</span>

    <span>original</span> = <span>"original"</span>
            <span>}</span>

    <span>updated_at</span> = <span>"updated_at"</span>

    <span>content_type</span> = <span>"content_type"</span>

    <span>file_size</span> = <span>1234</span>
            <span>}</span>
    <span>enable_custom_label_for_login_screen</span> = <span>true</span>
    <span>custom_masking_color</span> = <span>"custom_masking_color"</span>
    <span>mfa_enrollment_message</span> = <span>"mfa_enrollment_message"</span>
    <span>custom_masking_opacity</span> = <span>1234</span>
    <span>logo</span> <span>{</span>

    <span>urls</span> <span>{</span>

    <span>navigation</span> = <span>"navigation"</span>

    <span>login</span> = <span>"login"</span>

    <span>original</span> = <span>"original"</span>
            <span>}</span>

    <span>updated_at</span> = <span>"updated_at"</span>

    <span>content_type</span> = <span>"content_type"</span>

    <span>file_size</span> = <span>1234</span>
            <span>}</span>
    <span>id</span> = <span>1234</span>
    <span>login_instruction_title</span> = <span>"login_instruction_title"</span>
    <span>enabled</span> = <span>true</span>
    <span>custom_support_enabled</span> = <span>true</span>
<span>}</span>
</pre>

#### Arguments Reference

The following arguments are supported:

*   custom_accent_color [string] - (Required) Secondary brand color
*   custom_label_text_for_login_screen [string] - (Required) Custom label for the Username/Email field on the login screen. See example here.
*   hide_onelogin_footer [boolean] - (Required) Indicates if the OneLogin footer will appear at the bottom of the login page.
*   login_instruction [string] - (Required) Text for the login instruction screen, styled in Markdown.
*   custom_color [string] - (Required) Primary brand color
*   <span class="wysiwyg-color-red">*</span> background [object] - (Required) . The following properties compose the object schema :
    *   <span class="wysiwyg-color-red">*</span> urls [object] - (Required) . The following properties compose the object schema :
        *   branding [string] - (Required)
        *   login [string] - (Required)
        *   original [string] - (Required)
    *   updated_at [string] - (Required)
    *   content_type [string] - (Required)
    *   file_size [integer] - (Required)
*   enable_custom_label_for_login_screen [boolean] - (Required) Indicates if the custom Username/Email field label is enabled or not
*   custom_masking_color [string] - (Required) Color for the masking layer above the background image of the branded login screen.
*   mfa_enrollment_message [string] - (Required) Text that replaces the default text displayed on the initial screen of the MFA Registration.
*   custom_masking_opacity [integer] - (Required) Opacity for the custom_masking_color.
*   <span class="wysiwyg-color-red">*</span> logo [object] - (Required) . The following properties compose the object schema :
    *   <span class="wysiwyg-color-red">*</span> urls [object] - (Required) . The following properties compose the object schema :
        *   navigation [string] - (Required)
        *   login [string] - (Required)
        *   original [string] - (Required)
    *   updated_at [string] - (Required)
    *   content_type [string] - (Required)
    *   file_size [integer] - (Required)
*   id [integer] - (Required)
*   login_instruction_title [string] - (Required) Link text to show login instruction screen.
*   enabled [boolean] - (Required) Indicates if the brand is enabled or not. Default value is: false
*   custom_support_enabled [boolean] - (Required) Indicates if the custom support is enabled. If enabled, the login page includes the ability to submit a support request.

<span class="wysiwyg-color-red">*</span> Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin_brands.my_brands.**logo[0]**.object_property)

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

<span class="wysiwyg-color-red">*</span> Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin_brands.my_brands.**logo[0]**.object_property)

#### Import

brands resources can be imported using the `id` , e.g:

<pre dir="ltr">$ terraform import onelogin_brands.my_brands id</pre>

**Note**: In order for the import to work, the 'onelogin' terraform provider must be [properly installed](#provider_installation). Read more about Terraform import usage [here](https://www.terraform.io/docs/import/usage.html).

### onelogin_privileges

#### Example usage

<pre><span>resource</span> <span>"onelogin_privileges" "my_privileges"</span>{
    <span>name</span> = <span>"name"</span>
    <span>privilege</span> <span>{</span>

            <span>}</span>
<span>}</span>
</pre>

#### Arguments Reference

The following arguments are supported:

*   name [string] - (Required)
*   <span class="wysiwyg-color-red">*</span> privilege [object] - (Required) . The following properties compose the object schema :
    *   statement [list of objects] - (Optional) . The following properties compose the object schema :
        *   effect [string] - (Required) Set to “Allow.” By default, all actions are denied, this Statement allows the listed actions to be executed.
        *   action [list of strings] - (Required) An array of strings that represent actions within OneLogin. Actions are prefixed with the class of object they are related to and followed by a specific action for the given class. e.g. users:List, where the class is users and the specific action is List. Don’t mix classes within an Action array. To create a privilege that includes multiple different classes, create multiple statements. A wildcard * that includes all actions is supported. Use wildcards to create a Super User privilege.
        *   scope [list of strings] - (Required) Target the privileged action against specific resources with the scope. The scope pattern is the class of object used by the Action, followed by an ID that represents a resource in OneLogin. e.g. apps/1234, where apps is the class and 1234 is the ID of an app. The wildcard * is supported and indicates that all resources of the class type declared, in the Action, are in scope. The Action and Scope classes must match. However, there is an exception, a scope of roles/{role_id} can be combined with Actions on the user or app class. The exception allows you to target groups of users or apps with specific actions.
    *   version [string] - (Optional)
*   id [string] - (Optional)
*   description [string] - (Optional)

<span class="wysiwyg-color-red">*</span> Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin_privileges.my_privileges.**privilege[0]**.object_property)

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

<span class="wysiwyg-color-red">*</span> Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin_privileges.my_privileges.**privilege[0]**.object_property)

#### Import

privileges resources can be imported using the `id` , e.g:

<pre dir="ltr">$ terraform import onelogin_privileges.my_privileges id</pre>

**Note**: In order for the import to work, the 'onelogin' terraform provider must be [properly installed](#provider_installation). Read more about Terraform import usage [here](https://www.terraform.io/docs/import/usage.html).

### onelogin_risk_rules

#### Example usage

<pre><span>resource</span> <span>"onelogin_risk_rules" "my_risk_rules"</span>{
<span>}</span>
</pre>

#### Arguments Reference

The following arguments are supported:

*   id [string] - (Optional)
*   type [string] - (Optional) The type parameter specifies the type of rule that will be created.
*   target [string] - (Optional) The target parameter that will be used when evaluating the rule against an incoming event.
*   filters [list of strings] - (Optional) A list of IP addresses or country codes or names to evaluate against each event.
*   <span class="wysiwyg-color-red">*</span> source [object] - (Optional) Used for targeting custom rules based on a group of people, customers, accounts, or even a single user.. The following properties compose the object schema :
    *   name [string] - (Optional) The name of the source
    *   id [string] - (Optional) A unique id that represents the source of the event.
*   name [string] - (Optional) The name of this rule
*   description [string] - (Optional)

<span class="wysiwyg-color-red">*</span> Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin_risk_rules.my_risk_rules.**source[0]**.object_property)

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

<span class="wysiwyg-color-red">*</span> Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin_risk_rules.my_risk_rules.**source[0]**.object_property)

#### Import

risk_rules resources can be imported using the `id` , e.g:

<pre dir="ltr">$ terraform import onelogin_risk_rules.my_risk_rules id</pre>

**Note**: In order for the import to work, the 'onelogin' terraform provider must be [properly installed](#provider_installation). Read more about Terraform import usage [here](https://www.terraform.io/docs/import/usage.html).

### onelogin_roles

#### Example usage

<pre><span>resource</span> <span>"onelogin_roles" "my_roles"</span>{
    <span>name</span> = <span>"name"</span>
<span>}</span>
</pre>

#### Arguments Reference

The following arguments are supported:

*   name [string] - (Required) The name of the role.
*   admins [list of integers] - (Optional) A list of user IDs to assign as role administrators.
*   apps [list of integers] - (Optional) A list of app IDs that will be assigned to the role.
*   users [list of integers] - (Optional) A list of user IDs to assign to the role.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   id [integer] - Role ID

#### Import

roles resources can be imported using the `id` , e.g:

<pre dir="ltr">$ terraform import onelogin_roles.my_roles id</pre>

**Note**: In order for the import to work, the 'onelogin' terraform provider must be [properly installed](#provider_installation). Read more about Terraform import usage [here](https://www.terraform.io/docs/import/usage.html).

### onelogin_users

#### Example usage

<pre><span>resource</span> <span>"onelogin_users" "my_users"</span>{
<span>}</span>
</pre>

#### Arguments Reference

The following arguments are supported:

*   manager_ad_id [string] - (Optional) The ID of the user's manager in Active Directory.
*   salt [string] - (Optional) The salt value used with the password_algorithm.
*   password_changed_at [string] - (Optional)
*   firstname [string] - (Optional) The user's first name.
*   invitation_sent_at [string] - (Optional)
*   password [string] - (Optional) The password to set for a user.
*   username [string] - (Optional) A username for the user.
*   status [integer] - (Optional)
*   password_confirmation [string] - (Optional) Required if the password is being set.
*   password_algorithm [string] - (Optional) Use this when importing a password that's already hashed. Prepend the salt value to the cleartext password value before SHA-256-encoding it
*   phone [string] - (Optional) The E.164 format phone number for a user.
*   samaccountname [string] - (Optional) The user's Active Directory username.
*   invalid_login_attempts [integer] - (Optional)
*   email [string] - (Optional) A valid email for the user.
*   lastname [string] - (Optional) The user's last name.
*   locked_until [string] - (Optional)
*   id [integer] - (Optional)
*   title [string] - (Optional) The user's job title.
*   userprincipalname [string] - (Optional) The principle name of the user.
*   member_of [string] - (Optional) The user's directory membership.
*   role_ids [list of integers] - (Optional) A list of OneLogin Role IDs of the user
*   state [integer] - (Optional)
*   updated_at [string] - (Optional)
*   trusted_idp_id [integer] - (Optional) The ID of the OneLogin Trusted IDP of the user.
*   created_at [string] - (Optional)
*   preferred_locale_code [string] - (Optional)
*   group_id [integer] - (Optional) The ID of the Group in OneLogin that the user is assigned to.
*   directory_id [integer] - (Optional) The ID of the OneLogin Directory of the user.
*   distinguished_name [string] - (Optional) The distinguished name of the user.
*   company [string] - (Optional) The user's company.
*   manager_user_id [string] - (Optional) The OneLogin User ID for the user's manager.
*   comment [string] - (Optional) Free text related to the user.
*   department [string] - (Optional) The user's department.
*   external_id [string] - (Optional) The ID of the user in an external directory.
*   activated_at [string] - (Optional)
*   last_login [string] - (Optional)

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

#### Import

users resources can be imported using the `id` , e.g:

<pre dir="ltr">$ terraform import onelogin_users.my_users id</pre>

**Note**: In order for the import to work, the 'onelogin' terraform provider must be [properly installed](#provider_installation). Read more about Terraform import usage [here](https://www.terraform.io/docs/import/usage.html).

### onelogin_users_v1

#### Example usage

<pre><span>resource</span> <span>"onelogin_users_v1" "my_users_v1"</span>{
<span>}</span>
</pre>

#### Arguments Reference

The following arguments are supported:

*   manager_ad_id [string] - (Optional) The ID of the user's manager in Active Directory.
*   salt [string] - (Optional) The salt value used with the password_algorithm.
*   password_changed_at [string] - (Optional)
*   firstname [string] - (Optional) The user's first name.
*   invitation_sent_at [string] - (Optional)
*   password [string] - (Optional) The password to set for a user.
*   username [string] - (Optional) A username for the user.
*   status [integer] - (Optional)
*   password_confirmation [string] - (Optional) Required if the password is being set.
*   password_algorithm [string] - (Optional) Use this when importing a password that's already hashed. Prepend the salt value to the cleartext password value before SHA-256-encoding it
*   phone [string] - (Optional) The E.164 format phone number for a user.
*   samaccountname [string] - (Optional) The user's Active Directory username.
*   invalid_login_attempts [integer] - (Optional)
*   email [string] - (Optional) A valid email for the user.
*   lastname [string] - (Optional) The user's last name.
*   locked_until [string] - (Optional)
*   id [integer] - (Optional)
*   title [string] - (Optional) The user's job title.
*   userprincipalname [string] - (Optional) The principle name of the user.
*   member_of [string] - (Optional) The user's directory membership.
*   role_ids [list of integers] - (Optional) A list of OneLogin Role IDs of the user
*   state [integer] - (Optional)
*   updated_at [string] - (Optional)
*   trusted_idp_id [integer] - (Optional) The ID of the OneLogin Trusted IDP of the user.
*   created_at [string] - (Optional)
*   preferred_locale_code [string] - (Optional)
*   group_id [integer] - (Optional) The ID of the Group in OneLogin that the user is assigned to.
*   directory_id [integer] - (Optional) The ID of the OneLogin Directory of the user.
*   distinguished_name [string] - (Optional) The distinguished name of the user.
*   company [string] - (Optional) The user's company.
*   manager_user_id [string] - (Optional) The OneLogin User ID for the user's manager.
*   comment [string] - (Optional) Free text related to the user.
*   department [string] - (Optional) The user's department.
*   external_id [string] - (Optional) The ID of the user in an external directory.
*   activated_at [string] - (Optional)
*   last_login [string] - (Optional)

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

#### Import

users_v1 resources can be imported using the `id` , e.g:

<pre dir="ltr">$ terraform import onelogin_users_v1.my_users_v1 id</pre>

**Note**: In order for the import to work, the 'onelogin' terraform provider must be [properly installed](#provider_installation). Read more about Terraform import usage [here](https://www.terraform.io/docs/import/usage.html).

## Data Sources (using resource id)

### onelogin_apps_instance

Retrieve an existing resource using it's ID

#### Example usage

<pre><span>data</span> <span>"onelogin_apps_instance" "my_apps_instance"</span>{
    id = "existing_resource_id"
<span>}</span></pre>

#### Arguments Reference

The following arguments are supported:

*   id - (Required) ID of the existing resource to retrieve

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   tab_id [integer] - ID of the OneLogin portal tab that the app is assigned to.
*   <span class="wysiwyg-color-red">*</span> provisioning [object] - Indicates if provisioning is enabled for this app. The following properties compose the object schema:
    *   enabled [boolean]
*   role_ids [list of integers] - List of Role IDs that are assigned to the app. On App Create or Update the entire array is replaced with the values provided.
*   name [string] - The name of the app.
*   auth_method [integer] - An ID indicating the type of app: - 0: Password - 1: OpenId - 2: SAML - 3: API - 4: Google - 6: Forms Based App - 7: WSFED - 8: OpenId Connect
*   allow_assumed_signin [boolean] - Indicates whether or not administrators can access the app as a user that they have assumed control over.
*   notes [string] - Freeform notes about the app.
*   visible [boolean] - Indicates if the app is visible in the OneLogin portal.
*   updated_at [string] - the date the app was last updated
*   policy_id [integer] - The security policy assigned to the app.
*   id [integer] - Apps unique ID in OneLogin.
*   connector_id [integer] - ID of the connector to base the app from.
*   <span class="wysiwyg-color-red">*</span> parameters [object] - The parameters section contains parameterized attributes that have defined at the connector level as well as custom attributes that have been defined specifically for this app. Regardless of how they are defined, all parameters have the following attributes. Each parameter is an object with the key for the object being set as the parameters short name. The following properties compose the object schema:
    *   include_in_saml_assertion [boolean] - When true, this parameter will be included in a SAML assertion payload.
    *   label [string] - The can only be set when creating a new parameter. It can not be updated.
    *   user_attribute_mappings [string] - A user attribute to map values from For custom attributes prefix the name of the attribute with `custom_attribute_`. e.g. To get the value for custom attribute `employee_id` use `custom_attribute_employee_id`.
    *   user_attribute_macros [string] - When `user_attribute_mappings` is set to `_macro_` this macro will be used to assign the parameter value.
*   created_at [string] - the date the app was created
*   description [string] - Freeform description of the app.
*   <span class="wysiwyg-color-red">*</span> enforcement_point [object] - For apps that connect to a OneLogin Access Enforcement Point the following enforcement_point object will be included with the app payload. The following properties compose the object schema:
    *   permissions [string] - Specify to always `allow`, `deny` access to resources, of if access is `conditional`.
    *   case_sensitive [boolean] - The URL path evaluation is case insensitive by default. Resources hosted on web servers such as Apache, NGINX and Java EE are case sensitive paths. Web servers such as Microsoft IIS are not case-sensitive.
    *   <span class="wysiwyg-color-red">*</span> session_expiry_fixed [object] - unit: - 0 = Seconds - 1 = Minutes - 2 = Hours value: - When Unit = 0 or 1 value must be 0-60 - When Unit = 2 value must be 0-24 The following properties compose the object schema:
        *   unit [integer]
        *   value [integer]
    *   resources [list of objects] - Array of resource objects The following properties compose the object schema:
        *   path [string]
        *   permission [string]
        *   is_path_regex [boolean]
        *   require_auth [boolean]
        *   conditions [string] - required if permission == "conditions"
    *   vhost [string] - A comma-delimited list of one or more virtual hosts that map to applications assigned to the enforcement point. A VHOST may be a host name or an IP address. VHOST distinguish between applications that are at the same context root.
    *   target [string] - A fully-qualified URL to the internal application including scheme, authority and path. The target host authority must be an IP address, not a hostname.
    *   token [string] - Can only be set on create. Access Gateway Token.
    *   landing_page [string] - The location within the context root to which the browser will be redirected for IdP-initiated single sign-on. For example, the landing page might be an index page in the context root such as index.html or default.aspx. The landing page cannot begin with a slash and must use valid URL characters.
    *   use_target_host_header [boolean] - Use the target host header as opposed to the original gateway or upstream host header.
    *   require_sitewide_authentication [boolean] - Require user authentication to access any resource protected by this enforcement point.
    *   context_root [string] - The root path to the application, often the name of the application. Can be any name, path or just a slash (“/”). The context root uniquely identifies the application within the enforcement point.
    *   <span class="wysiwyg-color-red">*</span> session_expiry_inactivity [object] - unit: - 0 = Seconds - 1 = Minutes - 2 = Hours value: - When Unit = 0 or 1 value must be 0-60 - When Unit = 2 value must be 0-24 The following properties compose the object schema:
        *   unit [integer]
        *   value [integer]
    *   conditions [string] - If access is conditional, the conditions that must evaluate to true to allow access to a resource. For example, to require the user must be authenticated and have either the role Admin or User
*   icon_url [string] - A link to the apps icon url

<span class="wysiwyg-color-red">*</span> Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin_apps_instance.my_apps_instance.**enforcement_point[0]**.object_property)

### onelogin_apps_rules_instance

Retrieve an existing resource using it's ID

#### Example usage

<pre><span>data</span> <span>"onelogin_apps_rules_instance" "my_apps_rules_instance"</span>{
    id = "existing_resource_id"
<span>}</span></pre>

#### Arguments Reference

The following arguments are supported:

*   id - (Required) ID of the existing resource to retrieve

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   match [string] - Indicates how conditions should be matched.
*   conditions [list of objects] - An array of conditions that the user must meet in order for the rule to be applied. The following properties compose the object schema:
    *   source [string] - source field to check.
    *   operator [string] - A valid operator for the selected condition source
    *   value [string] - A plain text string or valid value for the selected condition source
*   name [string] - Rule Name
*   enabled [boolean] - Indicates if the rule is enabled or not.
*   id [integer] - App Rule ID
*   actions [list of objects] The following properties compose the object schema:
    *   value [list of strings] - Only applicable to provisioned and set_* actions. Items in the array will be a plain text string or valid value for the selected action.
    *   action [string] - The action to apply
*   position [integer] - Indicates the order of the rule. When `null` this will default to last position.

### onelogin_auth_servers_instance

Retrieve an existing resource using it's ID

#### Example usage

<pre><span>data</span> <span>"onelogin_auth_servers_instance" "my_auth_servers_instance"</span>{
    id = "existing_resource_id"
<span>}</span></pre>

#### Arguments Reference

The following arguments are supported:

*   id - (Required) ID of the existing resource to retrieve

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   description [string] - Description of what the API does.
*   name [string] - Name of the API.
*   id [integer] - Auth server unique ID in Onelogin
*   <span class="wysiwyg-color-red">*</span> configuration [object] - Authorization server configuration The following properties compose the object schema:
    *   access_token_expiration_minutes [integer] - The number of minutes until access token expires. There is no maximum expiry limit.
    *   refresh_token_expiration_minutes [integer] - The number of minutes until refresh token expires. There is no maximum expiry limit.
    *   audiences [list of strings] - List of API endpoints that will be returned in Access Tokens.
    *   resource_identifier [string] - Unique identifier for the API that the Authorization Server will issue Access Tokens for.

<span class="wysiwyg-color-red">*</span> Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin_auth_servers_instance.my_auth_servers_instance.**configuration[0]**.object_property)

### onelogin_brands_instance

Retrieve an existing resource using it's ID

#### Example usage

<pre><span>data</span> <span>"onelogin_brands_instance" "my_brands_instance"</span>{
    id = "existing_resource_id"
<span>}</span></pre>

#### Arguments Reference

The following arguments are supported:

*   id - (Required) ID of the existing resource to retrieve

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   login_instruction_title [string] - Link text to show login instruction screen.
*   <span class="wysiwyg-color-red">*</span> logo [object] The following properties compose the object schema:
    *   file_size [integer]
    *   <span class="wysiwyg-color-red">*</span> urls [object] The following properties compose the object schema:
        *   login [string]
        *   original [string]
        *   navigation [string]
    *   updated_at [string]
    *   content_type [string]
*   custom_support_enabled [boolean] - Indicates if the custom support is enabled. If enabled, the login page includes the ability to submit a support request.
*   mfa_enrollment_message [string] - Text that replaces the default text displayed on the initial screen of the MFA Registration.
*   custom_masking_color [string] - Color for the masking layer above the background image of the branded login screen.
*   custom_masking_opacity [integer] - Opacity for the custom_masking_color.
*   id [integer]
*   enabled [boolean] - Indicates if the brand is enabled or not
*   login_instruction [string] - Text for the login instruction screen, styled in Markdown.
*   custom_color [string] - Primary brand color
*   enable_custom_label_for_login_screen [boolean] - Indicates if the custom Username/Email field label is enabled or not
*   custom_accent_color [string] - Secondary brand color
*   custom_label_text_for_login_screen [string] - Custom label for the Username/Email field on the login screen. See example here.
*   hide_onelogin_footer [boolean] - Indicates if the OneLogin footer will appear at the bottom of the login page.
*   <span class="wysiwyg-color-red">*</span> background [object] The following properties compose the object schema:
    *   file_size [integer]
    *   updated_at [string]
    *   content_type [string]
    *   <span class="wysiwyg-color-red">*</span> urls [object] The following properties compose the object schema:
        *   branding [string]
        *   login [string]
        *   original [string]

<span class="wysiwyg-color-red">*</span> Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin_brands_instance.my_brands_instance.**background[0]**.object_property)

### onelogin_privileges_instance

Retrieve an existing resource using it's ID

#### Example usage

<pre><span>data</span> <span>"onelogin_privileges_instance" "my_privileges_instance"</span>{
    id = "existing_resource_id"
<span>}</span></pre>

#### Arguments Reference

The following arguments are supported:

*   id - (Required) ID of the existing resource to retrieve

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   id [string]
*   <span class="wysiwyg-color-red">*</span> privilege [object] The following properties compose the object schema:
    *   statement [list of objects] The following properties compose the object schema:
        *   scope [list of strings] - Target the privileged action against specific resources with the scope. The scope pattern is the class of object used by the Action, followed by an ID that represents a resource in OneLogin. e.g. apps/1234, where apps is the class and 1234 is the ID of an app. The wildcard * is supported and indicates that all resources of the class type declared, in the Action, are in scope. The Action and Scope classes must match. However, there is an exception, a scope of roles/{role_id} can be combined with Actions on the user or app class. The exception allows you to target groups of users or apps with specific actions.
        *   action [list of strings] - An array of strings that represent actions within OneLogin. Actions are prefixed with the class of object they are related to and followed by a specific action for the given class. e.g. users:List, where the class is users and the specific action is List. Don’t mix classes within an Action array. To create a privilege that includes multiple different classes, create multiple statements. A wildcard * that includes all actions is supported. Use wildcards to create a Super User privilege.
        *   effect [string] - Set to “Allow.” By default, all actions are denied, this Statement allows the listed actions to be executed.
    *   version [string]
*   description [string]
*   name [string]

<span class="wysiwyg-color-red">*</span> Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin_privileges_instance.my_privileges_instance.**privilege[0]**.object_property)

### onelogin_risk_rules_instance

Retrieve an existing resource using it's ID

#### Example usage

<pre><span>data</span> <span>"onelogin_risk_rules_instance" "my_risk_rules_instance"</span>{
    id = "existing_resource_id"
<span>}</span></pre>

#### Arguments Reference

The following arguments are supported:

*   id - (Required) ID of the existing resource to retrieve

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   id [string]
*   type [string] - The type parameter specifies the type of rule that will be created.
*   filters [list of strings] - A list of IP addresses or country codes or names to evaluate against each event.
*   target [string] - The target parameter that will be used when evaluating the rule against an incoming event.
*   <span class="wysiwyg-color-red">*</span> source [object] - Used for targeting custom rules based on a group of people, customers, accounts, or even a single user. The following properties compose the object schema:
    *   name [string] - The name of the source
    *   id [string] - A unique id that represents the source of the event.
*   name [string] - The name of this rule
*   description [string]

<span class="wysiwyg-color-red">*</span> Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin_risk_rules_instance.my_risk_rules_instance.**source[0]**.object_property)

### onelogin_roles_instance

Retrieve an existing resource using it's ID

#### Example usage

<pre><span>data</span> <span>"onelogin_roles_instance" "my_roles_instance"</span>{
    id = "existing_resource_id"
<span>}</span></pre>

#### Arguments Reference

The following arguments are supported:

*   id - (Required) ID of the existing resource to retrieve

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   admins [list of integers] - A list of user IDs to assign as role administrators.
*   name [string] - The name of the role.
*   apps [list of integers] - A list of app IDs that will be assigned to the role.
*   id [integer] - Role ID
*   users [list of integers] - A list of user IDs to assign to the role.

### onelogin_users_instance

Retrieve an existing resource using it's ID

#### Example usage

<pre><span>data</span> <span>"onelogin_users_instance" "my_users_instance"</span>{
    id = "existing_resource_id"
<span>}</span></pre>

#### Arguments Reference

The following arguments are supported:

*   id - (Required) ID of the existing resource to retrieve

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   invitation_sent_at [string]
*   firstname [string] - The user's first name.
*   salt [string] - The salt value used with the password_algorithm.
*   password_changed_at [string]
*   manager_ad_id [string] - The ID of the user's manager in Active Directory.
*   phone [string] - The E.164 format phone number for a user.
*   samaccountname [string] - The user's Active Directory username.
*   password_algorithm [string] - Use this when importing a password that's already hashed. Prepend the salt value to the cleartext password value before SHA-256-encoding it
*   password_confirmation [string] - Required if the password is being set.
*   password [string] - The password to set for a user.
*   status [integer]
*   username [string] - A username for the user.
*   locked_until [string]
*   lastname [string] - The user's last name.
*   email [string] - A valid email for the user.
*   invalid_login_attempts [integer]
*   userprincipalname [string] - The principle name of the user.
*   member_of [string] - The user's directory membership.
*   title [string] - The user's job title.
*   id [integer]
*   updated_at [string]
*   state [integer]
*   role_ids [list of integers] - A list of OneLogin Role IDs of the user
*   group_id [integer] - The ID of the Group in OneLogin that the user is assigned to.
*   preferred_locale_code [string]
*   directory_id [integer] - The ID of the OneLogin Directory of the user.
*   created_at [string]
*   trusted_idp_id [integer] - The ID of the OneLogin Trusted IDP of the user.
*   company [string] - The user's company.
*   distinguished_name [string] - The distinguished name of the user.
*   activated_at [string]
*   external_id [string] - The ID of the user in an external directory.
*   last_login [string]
*   comment [string] - Free text related to the user.
*   department [string] - The user's department.
*   manager_user_id [string] - The OneLogin User ID for the user's manager.

### onelogin_users_v1_instance

Retrieve an existing resource using it's ID

#### Example usage

<pre><span>data</span> <span>"onelogin_users_v1_instance" "my_users_v1_instance"</span>{
    id = "existing_resource_id"
<span>}</span></pre>

#### Arguments Reference

The following arguments are supported:

*   id - (Required) ID of the existing resource to retrieve

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   invitation_sent_at [string]
*   firstname [string] - The user's first name.
*   salt [string] - The salt value used with the password_algorithm.
*   password_changed_at [string]
*   manager_ad_id [string] - The ID of the user's manager in Active Directory.
*   phone [string] - The E.164 format phone number for a user.
*   samaccountname [string] - The user's Active Directory username.
*   password_algorithm [string] - Use this when importing a password that's already hashed. Prepend the salt value to the cleartext password value before SHA-256-encoding it
*   password_confirmation [string] - Required if the password is being set.
*   password [string] - The password to set for a user.
*   status [integer]
*   username [string] - A username for the user.
*   locked_until [string]
*   lastname [string] - The user's last name.
*   email [string] - A valid email for the user.
*   invalid_login_attempts [integer]
*   userprincipalname [string] - The principle name of the user.
*   member_of [string] - The user's directory membership.
*   title [string] - The user's job title.
*   id [integer]
*   updated_at [string]
*   state [integer]
*   role_ids [list of integers] - A list of OneLogin Role IDs of the user
*   group_id [integer] - The ID of the Group in OneLogin that the user is assigned to.
*   preferred_locale_code [string]
*   directory_id [integer] - The ID of the OneLogin Directory of the user.
*   created_at [string]
*   trusted_idp_id [integer] - The ID of the OneLogin Trusted IDP of the user.
*   company [string] - The user's company.
*   distinguished_name [string] - The distinguished name of the user.
*   activated_at [string]
*   external_id [string] - The ID of the user in an external directory.
*   last_login [string]
*   comment [string] - Free text related to the user.
*   department [string] - The user's department.
*   manager_user_id [string] - The OneLogin User ID for the user's manager.

## Data Sources (using filters)

### onelogin_apps (filters)

The apps data source allows you to retrieve an already existing apps resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

<pre><span>data</span> <span>"onelogin_apps" "my_apps"</span>{
    <span>filter</span> <span>{</span>
        <span>name</span> = <span>"property name to filter by, see docs below for more info about available filter name options"</span>
        <span>values</span> = <span>["filter value"]</span>
    <span>}</span>
<span>}</span></pre>

#### Arguments Reference

The following arguments are supported:

*   filter - (Required) Object containing two properties.

*   name [string]: the name should match one of the properties to filter by. The following property names are supported: <span>tab_id,</span> <span>name,</span> <span>auth_method,</span> <span>allow_assumed_signin,</span> <span>notes,</span> <span>visible,</span> <span>updated_at,</span> <span>policy_id,</span> <span>id,</span> <span>connector_id,</span> <span>created_at,</span> <span>description,</span> <span>icon_url,</span>
*   values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   tab_id [integer] - ID of the OneLogin portal tab that the app is assigned to.
*   <span class="wysiwyg-color-red">*</span> provisioning [object] - Indicates if provisioning is enabled for this app. The following properties compose the object schema:
    *   enabled [boolean]
*   role_ids [list of integers] - List of Role IDs that are assigned to the app. On App Create or Update the entire array is replaced with the values provided.
*   name [string] - The name of the app.
*   auth_method [integer] - An ID indicating the type of app: - 0: Password - 1: OpenId - 2: SAML - 3: API - 4: Google - 6: Forms Based App - 7: WSFED - 8: OpenId Connect
*   allow_assumed_signin [boolean] - Indicates whether or not administrators can access the app as a user that they have assumed control over.
*   notes [string] - Freeform notes about the app.
*   visible [boolean] - Indicates if the app is visible in the OneLogin portal.
*   updated_at [string] - the date the app was last updated
*   policy_id [integer] - The security policy assigned to the app.
*   id [integer] - Apps unique ID in OneLogin.
*   connector_id [integer] - ID of the connector to base the app from.
*   <span class="wysiwyg-color-red">*</span> parameters [object] - The parameters section contains parameterized attributes that have defined at the connector level as well as custom attributes that have been defined specifically for this app. Regardless of how they are defined, all parameters have the following attributes. Each parameter is an object with the key for the object being set as the parameters short name. The following properties compose the object schema:
    *   include_in_saml_assertion [boolean] - When true, this parameter will be included in a SAML assertion payload.
    *   label [string] - The can only be set when creating a new parameter. It can not be updated.
    *   user_attribute_mappings [string] - A user attribute to map values from For custom attributes prefix the name of the attribute with `custom_attribute_`. e.g. To get the value for custom attribute `employee_id` use `custom_attribute_employee_id`.
    *   user_attribute_macros [string] - When `user_attribute_mappings` is set to `_macro_` this macro will be used to assign the parameter value.
*   created_at [string] - the date the app was created
*   description [string] - Freeform description of the app.
*   <span class="wysiwyg-color-red">*</span> enforcement_point [object] - For apps that connect to a OneLogin Access Enforcement Point the following enforcement_point object will be included with the app payload. The following properties compose the object schema:
    *   permissions [string] - Specify to always `allow`, `deny` access to resources, of if access is `conditional`.
    *   case_sensitive [boolean] - The URL path evaluation is case insensitive by default. Resources hosted on web servers such as Apache, NGINX and Java EE are case sensitive paths. Web servers such as Microsoft IIS are not case-sensitive.
    *   <span class="wysiwyg-color-red">*</span> session_expiry_fixed [object] - unit: - 0 = Seconds - 1 = Minutes - 2 = Hours value: - When Unit = 0 or 1 value must be 0-60 - When Unit = 2 value must be 0-24 The following properties compose the object schema:
        *   unit [integer]
        *   value [integer]
    *   resources [list of objects] - Array of resource objects The following properties compose the object schema:
        *   path [string]
        *   permission [string]
        *   is_path_regex [boolean]
        *   require_auth [boolean]
        *   conditions [string] - required if permission == "conditions"
    *   vhost [string] - A comma-delimited list of one or more virtual hosts that map to applications assigned to the enforcement point. A VHOST may be a host name or an IP address. VHOST distinguish between applications that are at the same context root.
    *   target [string] - A fully-qualified URL to the internal application including scheme, authority and path. The target host authority must be an IP address, not a hostname.
    *   token [string] - Can only be set on create. Access Gateway Token.
    *   landing_page [string] - The location within the context root to which the browser will be redirected for IdP-initiated single sign-on. For example, the landing page might be an index page in the context root such as index.html or default.aspx. The landing page cannot begin with a slash and must use valid URL characters.
    *   use_target_host_header [boolean] - Use the target host header as opposed to the original gateway or upstream host header.
    *   require_sitewide_authentication [boolean] - Require user authentication to access any resource protected by this enforcement point.
    *   context_root [string] - The root path to the application, often the name of the application. Can be any name, path or just a slash (“/”). The context root uniquely identifies the application within the enforcement point.
    *   <span class="wysiwyg-color-red">*</span> session_expiry_inactivity [object] - unit: - 0 = Seconds - 1 = Minutes - 2 = Hours value: - When Unit = 0 or 1 value must be 0-60 - When Unit = 2 value must be 0-24 The following properties compose the object schema:
        *   unit [integer]
        *   value [integer]
    *   conditions [string] - If access is conditional, the conditions that must evaluate to true to allow access to a resource. For example, to require the user must be authenticated and have either the role Admin or User
*   icon_url [string] - A link to the apps icon url

<span class="wysiwyg-color-red">*</span> Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin_apps.my_apps.**enforcement_point[0]**.object_property)

### onelogin_apps_actions (filters)

The apps_actions data source allows you to retrieve an already existing apps_actions resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

<pre><span>data</span> <span>"onelogin_apps_actions" "my_apps_actions"</span>{
    <span>filter</span> <span>{</span>
        <span>name</span> = <span>"property name to filter by, see docs below for more info about available filter name options"</span>
        <span>values</span> = <span>["filter value"]</span>
    <span>}</span>
<span>}</span></pre>

#### Arguments Reference

The following arguments are supported:

*   filter - (Required) Object containing two properties.

*   name [string]: the name should match one of the properties to filter by. The following property names are supported: <span>apps_id,</span> <span>name,</span> <span>value,</span>
*   values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   name [string] - Name of the Action
*   value [string] - The unique identifier of the action. This should be used when defining actions for a User Mapping.

### onelogin_apps_actions_values (filters)

The apps_actions_values data source allows you to retrieve an already existing apps_actions_values resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

<pre><span>data</span> <span>"onelogin_apps_actions_values" "my_apps_actions_values"</span>{
    <span>filter</span> <span>{</span>
        <span>name</span> = <span>"property name to filter by, see docs below for more info about available filter name options"</span>
        <span>values</span> = <span>["filter value"]</span>
    <span>}</span>
<span>}</span></pre>

#### Arguments Reference

The following arguments are supported:

*   filter - (Required) Object containing two properties.

*   name [string]: the name should match one of the properties to filter by. The following property names are supported: <span>actions_id,</span> <span>apps_id,</span> <span>name,</span> <span>value,</span>
*   values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   name [string] - Name of the Action
*   value [string] - The unique identifier of the action. This should be used when defining actions for a User Mapping.

### onelogin_apps_conditions (filters)

The apps_conditions data source allows you to retrieve an already existing apps_conditions resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

<pre><span>data</span> <span>"onelogin_apps_conditions" "my_apps_conditions"</span>{
    <span>filter</span> <span>{</span>
        <span>name</span> = <span>"property name to filter by, see docs below for more info about available filter name options"</span>
        <span>values</span> = <span>["filter value"]</span>
    <span>}</span>
<span>}</span></pre>

#### Arguments Reference

The following arguments are supported:

*   filter - (Required) Object containing two properties.

*   name [string]: the name should match one of the properties to filter by. The following property names are supported: <span>value,</span> <span>apps_id,</span> <span>name,</span>
*   values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   value [string] - The unique identifier of the condition. This should be used when defining conditions for a rule.
*   name [string] - Name of the rule condition

### onelogin_apps_conditions_operators (filters)

The apps_conditions_operators data source allows you to retrieve an already existing apps_conditions_operators resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

<pre><span>data</span> <span>"onelogin_apps_conditions_operators" "my_apps_conditions_operators"</span>{
    <span>filter</span> <span>{</span>
        <span>name</span> = <span>"property name to filter by, see docs below for more info about available filter name options"</span>
        <span>values</span> = <span>["filter value"]</span>
    <span>}</span>
<span>}</span></pre>

#### Arguments Reference

The following arguments are supported:

*   filter - (Required) Object containing two properties.

*   name [string]: the name should match one of the properties to filter by. The following property names are supported: <span>name,</span> <span>apps_id,</span> <span>value,</span> <span>conditions_id,</span>
*   values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   name [string] - Name of the operator
*   value [string] - The condition operator value to use when creating or updating rules.

### onelogin_apps_rules (filters)

The apps_rules data source allows you to retrieve an already existing apps_rules resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

<pre><span>data</span> <span>"onelogin_apps_rules" "my_apps_rules"</span>{
    <span>filter</span> <span>{</span>
        <span>name</span> = <span>"property name to filter by, see docs below for more info about available filter name options"</span>
        <span>values</span> = <span>["filter value"]</span>
    <span>}</span>
<span>}</span></pre>

#### Arguments Reference

The following arguments are supported:

*   filter - (Required) Object containing two properties.

*   name [string]: the name should match one of the properties to filter by. The following property names are supported: <span>match,</span> <span>name,</span> <span>apps_id,</span> <span>enabled,</span> <span>id,</span> <span>position,</span>
*   values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   match [string] - Indicates how conditions should be matched.
*   conditions [list of objects] - An array of conditions that the user must meet in order for the rule to be applied. The following properties compose the object schema:
    *   source [string] - source field to check.
    *   operator [string] - A valid operator for the selected condition source
    *   value [string] - A plain text string or valid value for the selected condition source
*   name [string] - Rule Name
*   enabled [boolean] - Indicates if the rule is enabled or not.
*   id [integer] - App Rule ID
*   actions [list of objects] The following properties compose the object schema:
    *   value [list of strings] - Only applicable to provisioned and set_* actions. Items in the array will be a plain text string or valid value for the selected action.
    *   action [string] - The action to apply
*   position [integer] - Indicates the order of the rule. When `null` this will default to last position.

### onelogin_apps_users (filters)

The apps_users data source allows you to retrieve an already existing apps_users resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

<pre><span>data</span> <span>"onelogin_apps_users" "my_apps_users"</span>{
    <span>filter</span> <span>{</span>
        <span>name</span> = <span>"property name to filter by, see docs below for more info about available filter name options"</span>
        <span>values</span> = <span>["filter value"]</span>
    <span>}</span>
<span>}</span></pre>

#### Arguments Reference

The following arguments are supported:

*   filter - (Required) Object containing two properties.

*   name [string]: the name should match one of the properties to filter by. The following property names are supported: <span>invitation_sent_at,</span> <span>firstname,</span> <span>salt,</span> <span>password_changed_at,</span> <span>manager_ad_id,</span> <span>phone,</span> <span>samaccountname,</span> <span>password_algorithm,</span> <span>password_confirmation,</span> <span>password,</span> <span>status,</span> <span>username,</span> <span>locked_until,</span> <span>lastname,</span> <span>email,</span> <span>apps_id,</span> <span>invalid_login_attempts,</span> <span>userprincipalname,</span> <span>member_of,</span> <span>title,</span> <span>id,</span> <span>updated_at,</span> <span>state,</span> <span>group_id,</span> <span>preferred_locale_code,</span> <span>directory_id,</span> <span>created_at,</span> <span>trusted_idp_id,</span> <span>company,</span> <span>distinguished_name,</span> <span>activated_at,</span> <span>external_id,</span> <span>last_login,</span> <span>comment,</span> <span>department,</span> <span>manager_user_id,</span>
*   values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   invitation_sent_at [string]
*   firstname [string] - The user's first name.
*   salt [string] - The salt value used with the password_algorithm.
*   password_changed_at [string]
*   manager_ad_id [string] - The ID of the user's manager in Active Directory.
*   phone [string] - The E.164 format phone number for a user.
*   samaccountname [string] - The user's Active Directory username.
*   password_algorithm [string] - Use this when importing a password that's already hashed. Prepend the salt value to the cleartext password value before SHA-256-encoding it
*   password_confirmation [string] - Required if the password is being set.
*   password [string] - The password to set for a user.
*   status [integer]
*   username [string] - A username for the user.
*   locked_until [string]
*   lastname [string] - The user's last name.
*   email [string] - A valid email for the user.
*   invalid_login_attempts [integer]
*   userprincipalname [string] - The principle name of the user.
*   member_of [string] - The user's directory membership.
*   title [string] - The user's job title.
*   id [integer]
*   updated_at [string]
*   state [integer]
*   role_ids [list of integers] - A list of OneLogin Role IDs of the user
*   group_id [integer] - The ID of the Group in OneLogin that the user is assigned to.
*   preferred_locale_code [string]
*   directory_id [integer] - The ID of the OneLogin Directory of the user.
*   created_at [string]
*   trusted_idp_id [integer] - The ID of the OneLogin Trusted IDP of the user.
*   company [string] - The user's company.
*   distinguished_name [string] - The distinguished name of the user.
*   activated_at [string]
*   external_id [string] - The ID of the user in an external directory.
*   last_login [string]
*   comment [string] - Free text related to the user.
*   department [string] - The user's department.
*   manager_user_id [string] - The OneLogin User ID for the user's manager.

### onelogin_auth_servers (filters)

The auth_servers data source allows you to retrieve an already existing auth_servers resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

<pre><span>data</span> <span>"onelogin_auth_servers" "my_auth_servers"</span>{
    <span>filter</span> <span>{</span>
        <span>name</span> = <span>"property name to filter by, see docs below for more info about available filter name options"</span>
        <span>values</span> = <span>["filter value"]</span>
    <span>}</span>
<span>}</span></pre>

#### Arguments Reference

The following arguments are supported:

*   filter - (Required) Object containing two properties.

*   name [string]: the name should match one of the properties to filter by. The following property names are supported: <span>description,</span> <span>name,</span> <span>id,</span>
*   values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   description [string] - Description of what the API does.
*   name [string] - Name of the API.
*   id [integer] - Auth server unique ID in Onelogin
*   <span class="wysiwyg-color-red">*</span> configuration [object] - Authorization server configuration The following properties compose the object schema:
    *   access_token_expiration_minutes [integer] - The number of minutes until access token expires. There is no maximum expiry limit.
    *   refresh_token_expiration_minutes [integer] - The number of minutes until refresh token expires. There is no maximum expiry limit.
    *   audiences [list of strings] - List of API endpoints that will be returned in Access Tokens.
    *   resource_identifier [string] - Unique identifier for the API that the Authorization Server will issue Access Tokens for.

<span class="wysiwyg-color-red">*</span> Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin_auth_servers.my_auth_servers.**configuration[0]**.object_property)

### onelogin_auth_servers_claims (filters)

The auth_servers_claims data source allows you to retrieve an already existing auth_servers_claims resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

<pre><span>data</span> <span>"onelogin_auth_servers_claims" "my_auth_servers_claims"</span>{
    <span>filter</span> <span>{</span>
        <span>name</span> = <span>"property name to filter by, see docs below for more info about available filter name options"</span>
        <span>values</span> = <span>["filter value"]</span>
    <span>}</span>
<span>}</span></pre>

#### Arguments Reference

The following arguments are supported:

*   filter - (Required) Object containing two properties.

*   name [string]: the name should match one of the properties to filter by. The following property names are supported: <span>provisioned_entitlements,</span> <span>skip_if_blank,</span> <span>auth_servers_id,</span> <span>id,</span> <span>user_attribute_mappings,</span> <span>attribute_transformations,</span> <span>user_attribute_macros,</span> <span>default_values,</span> <span>label,</span>
*   values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   provisioned_entitlements [boolean] - Relates to Rules/Entitlements. Not supported yet.
*   skip_if_blank [boolean] - not used
*   id [integer] - The unique ID of the claim.
*   user_attribute_mappings [string] - A user attribute to map values from.
*   attribute_transformations [string] - The type of transformation to perform on multi valued attributes.
*   user_attribute_macros [string] - When `user_attribute_mappings` is set to `_macro_` this macro will be used to assign the claims value.
*   values [list of strings] - Relates to Rules/Entitlements. Not supported yet.
*   default_values [string] - Relates to Rules/Entitlements. Not supported yet.
*   label [string] - The UI label for the claims.

### onelogin_auth_servers_scopes (filters)

The auth_servers_scopes data source allows you to retrieve an already existing auth_servers_scopes resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

<pre><span>data</span> <span>"onelogin_auth_servers_scopes" "my_auth_servers_scopes"</span>{
    <span>filter</span> <span>{</span>
        <span>name</span> = <span>"property name to filter by, see docs below for more info about available filter name options"</span>
        <span>values</span> = <span>["filter value"]</span>
    <span>}</span>
<span>}</span></pre>

#### Arguments Reference

The following arguments are supported:

*   filter - (Required) Object containing two properties.

*   name [string]: the name should match one of the properties to filter by. The following property names are supported: <span>auth_servers_id,</span> <span>id,</span> <span>description,</span> <span>value,</span>
*   values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   id [integer] - Unique ID for the Scope
*   description [string] - A description of what access the scope enables
*   value [string] - A value representing the api scope that with be authorized

### onelogin_brands (filters)

The brands data source allows you to retrieve an already existing brands resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

<pre><span>data</span> <span>"onelogin_brands" "my_brands"</span>{
    <span>filter</span> <span>{</span>
        <span>name</span> = <span>"property name to filter by, see docs below for more info about available filter name options"</span>
        <span>values</span> = <span>["filter value"]</span>
    <span>}</span>
<span>}</span></pre>

#### Arguments Reference

The following arguments are supported:

*   filter - (Required) Object containing two properties.

*   name [string]: the name should match one of the properties to filter by. The following property names are supported: <span>enabled,</span> <span>id,</span> <span>name,</span>
*   values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   enabled [boolean] - Indicates if the brand is enabled or not.
*   id [integer] - Brand’s unique ID in OneLogin.
*   name [string] - Brand name for humans. This isn’t related to subdomains.

### onelogin_brands_apps (filters)

The brands_apps data source allows you to retrieve an already existing brands_apps resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

<pre><span>data</span> <span>"onelogin_brands_apps" "my_brands_apps"</span>{
    <span>filter</span> <span>{</span>
        <span>name</span> = <span>"property name to filter by, see docs below for more info about available filter name options"</span>
        <span>values</span> = <span>["filter value"]</span>
    <span>}</span>
<span>}</span></pre>

#### Arguments Reference

The following arguments are supported:

*   filter - (Required) Object containing two properties.

*   name [string]: the name should match one of the properties to filter by. The following property names are supported: <span>visible,</span> <span>connector_id,</span> <span>id,</span> <span>updated_at,</span> <span>auth_method,</span> <span>created_at,</span> <span>description,</span> <span>brands_id,</span> <span>auth_method_description,</span> <span>name,</span>
*   values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   visible [boolean]
*   connector_id [integer]
*   id [integer]
*   updated_at [string]
*   auth_method [integer]
*   created_at [string]
*   description [string]
*   auth_method_description [string]
*   name [string]

### onelogin_brands_templates (filters)

The brands_templates data source allows you to retrieve an already existing brands_templates resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

<pre><span>data</span> <span>"onelogin_brands_templates" "my_brands_templates"</span>{
    <span>filter</span> <span>{</span>
        <span>name</span> = <span>"property name to filter by, see docs below for more info about available filter name options"</span>
        <span>values</span> = <span>["filter value"]</span>
    <span>}</span>
<span>}</span></pre>

#### Arguments Reference

The following arguments are supported:

*   filter - (Required) Object containing two properties.

*   name [string]: the name should match one of the properties to filter by. The following property names are supported: <span>name,</span> <span>enabled,</span> <span>id,</span> <span>brands_id,</span>
*   values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   name [string] - name of message template
*   enabled [boolean] - indicator if template is enabled
*   id [integer] - template ID

### onelogin_mappings (filters)

The mappings data source allows you to retrieve an already existing mappings resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

<pre><span>data</span> <span>"onelogin_mappings" "my_mappings"</span>{
    <span>filter</span> <span>{</span>
        <span>name</span> = <span>"property name to filter by, see docs below for more info about available filter name options"</span>
        <span>values</span> = <span>["filter value"]</span>
    <span>}</span>
<span>}</span></pre>

#### Arguments Reference

The following arguments are supported:

*   filter - (Required) Object containing two properties.

*   name [string]: the name should match one of the properties to filter by. The following property names are supported: <span>match,</span> <span>name,</span> <span>id,</span> <span>enabled,</span> <span>position,</span>
*   values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   match [string] - Indicates how conditions should be matched.
*   name [string] - The name of the mapping.
*   actions [list of objects] - An array of actions that will be applied to the users that are matched by the conditions. The following properties compose the object schema:
    *   value [list of strings] - Only applicable to provisioned and set_* actions. Items in the array will be a plain text string or valid value for the selected action.
    *   action [string] - The action to apply
*   id [integer]
*   enabled [boolean] - Indicates if the mapping is enabled or not.
*   conditions [list of objects] - An array of conditions that the user must meet in order for the mapping to be applied. The following properties compose the object schema:
    *   source [string] - source field to check.
    *   operator [string] - A valid operator for the selected condition source
    *   value [string] - A plain text string or valid value for the selected condition source
*   position [integer] - Indicates the order of the mapping. When `null` this will default to last position.

### onelogin_privileges (filters)

The privileges data source allows you to retrieve an already existing privileges resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

<pre><span>data</span> <span>"onelogin_privileges" "my_privileges"</span>{
    <span>filter</span> <span>{</span>
        <span>name</span> = <span>"property name to filter by, see docs below for more info about available filter name options"</span>
        <span>values</span> = <span>["filter value"]</span>
    <span>}</span>
<span>}</span></pre>

#### Arguments Reference

The following arguments are supported:

*   filter - (Required) Object containing two properties.

*   name [string]: the name should match one of the properties to filter by. The following property names are supported: <span>id,</span> <span>description,</span> <span>name,</span>
*   values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   id [string]
*   <span class="wysiwyg-color-red">*</span> privilege [object] The following properties compose the object schema:
    *   statement [list of objects] The following properties compose the object schema:
        *   scope [list of strings] - Target the privileged action against specific resources with the scope. The scope pattern is the class of object used by the Action, followed by an ID that represents a resource in OneLogin. e.g. apps/1234, where apps is the class and 1234 is the ID of an app. The wildcard * is supported and indicates that all resources of the class type declared, in the Action, are in scope. The Action and Scope classes must match. However, there is an exception, a scope of roles/{role_id} can be combined with Actions on the user or app class. The exception allows you to target groups of users or apps with specific actions.
        *   action [list of strings] - An array of strings that represent actions within OneLogin. Actions are prefixed with the class of object they are related to and followed by a specific action for the given class. e.g. users:List, where the class is users and the specific action is List. Don’t mix classes within an Action array. To create a privilege that includes multiple different classes, create multiple statements. A wildcard * that includes all actions is supported. Use wildcards to create a Super User privilege.
        *   effect [string] - Set to “Allow.” By default, all actions are denied, this Statement allows the listed actions to be executed.
    *   version [string]
*   description [string]
*   name [string]

<span class="wysiwyg-color-red">*</span> Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin_privileges.my_privileges.**privilege[0]**.object_property)

### onelogin_risk_rules (filters)

The risk_rules data source allows you to retrieve an already existing risk_rules resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

<pre><span>data</span> <span>"onelogin_risk_rules" "my_risk_rules"</span>{
    <span>filter</span> <span>{</span>
        <span>name</span> = <span>"property name to filter by, see docs below for more info about available filter name options"</span>
        <span>values</span> = <span>["filter value"]</span>
    <span>}</span>
<span>}</span></pre>

#### Arguments Reference

The following arguments are supported:

*   filter - (Required) Object containing two properties.

*   name [string]: the name should match one of the properties to filter by. The following property names are supported: <span>id,</span> <span>type,</span> <span>target,</span> <span>name,</span> <span>description,</span>
*   values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   id [string]
*   type [string] - The type parameter specifies the type of rule that will be created.
*   filters [list of strings] - A list of IP addresses or country codes or names to evaluate against each event.
*   target [string] - The target parameter that will be used when evaluating the rule against an incoming event.
*   <span class="wysiwyg-color-red">*</span> source [object] - Used for targeting custom rules based on a group of people, customers, accounts, or even a single user. The following properties compose the object schema:
    *   name [string] - The name of the source
    *   id [string] - A unique id that represents the source of the event.
*   name [string] - The name of this rule
*   description [string]

<span class="wysiwyg-color-red">*</span> Note: Object type properties are internally represented (in the state file) as a list of one elem due to [Terraform SDK's limitation for supporting complex object types](https://github.com/hashicorp/terraform-plugin-sdk/issues/155#issuecomment-489699737). Please index on the first elem of the array to reference the object values (eg: onelogin_risk_rules.my_risk_rules.**source[0]**.object_property)

### onelogin_roles (filters)

The roles data source allows you to retrieve an already existing roles resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

<pre><span>data</span> <span>"onelogin_roles" "my_roles"</span>{
    <span>filter</span> <span>{</span>
        <span>name</span> = <span>"property name to filter by, see docs below for more info about available filter name options"</span>
        <span>values</span> = <span>["filter value"]</span>
    <span>}</span>
<span>}</span></pre>

#### Arguments Reference

The following arguments are supported:

*   filter - (Required) Object containing two properties.

*   name [string]: the name should match one of the properties to filter by. The following property names are supported: <span>name,</span> <span>id,</span>
*   values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   admins [list of integers] - A list of user IDs to assign as role administrators.
*   name [string] - The name of the role.
*   apps [list of integers] - A list of app IDs that will be assigned to the role.
*   id [integer] - Role ID
*   users [list of integers] - A list of user IDs to assign to the role.

### onelogin_roles_admins (filters)

The roles_admins data source allows you to retrieve an already existing roles_admins resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

<pre><span>data</span> <span>"onelogin_roles_admins" "my_roles_admins"</span>{
    <span>filter</span> <span>{</span>
        <span>name</span> = <span>"property name to filter by, see docs below for more info about available filter name options"</span>
        <span>values</span> = <span>["filter value"]</span>
    <span>}</span>
<span>}</span></pre>

#### Arguments Reference

The following arguments are supported:

*   filter - (Required) Object containing two properties.

*   name [string]: the name should match one of the properties to filter by. The following property names are supported: <span>invitation_sent_at,</span> <span>firstname,</span> <span>salt,</span> <span>password_changed_at,</span> <span>manager_ad_id,</span> <span>phone,</span> <span>samaccountname,</span> <span>password_algorithm,</span> <span>password_confirmation,</span> <span>password,</span> <span>status,</span> <span>username,</span> <span>locked_until,</span> <span>lastname,</span> <span>email,</span> <span>invalid_login_attempts,</span> <span>userprincipalname,</span> <span>member_of,</span> <span>title,</span> <span>id,</span> <span>updated_at,</span> <span>state,</span> <span>group_id,</span> <span>preferred_locale_code,</span> <span>directory_id,</span> <span>created_at,</span> <span>trusted_idp_id,</span> <span>company,</span> <span>distinguished_name,</span> <span>activated_at,</span> <span>external_id,</span> <span>last_login,</span> <span>comment,</span> <span>department,</span> <span>roles_id,</span> <span>manager_user_id,</span>
*   values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   invitation_sent_at [string]
*   firstname [string] - The user's first name.
*   salt [string] - The salt value used with the password_algorithm.
*   password_changed_at [string]
*   manager_ad_id [string] - The ID of the user's manager in Active Directory.
*   phone [string] - The E.164 format phone number for a user.
*   samaccountname [string] - The user's Active Directory username.
*   password_algorithm [string] - Use this when importing a password that's already hashed. Prepend the salt value to the cleartext password value before SHA-256-encoding it
*   password_confirmation [string] - Required if the password is being set.
*   password [string] - The password to set for a user.
*   status [integer]
*   username [string] - A username for the user.
*   locked_until [string]
*   lastname [string] - The user's last name.
*   email [string] - A valid email for the user.
*   invalid_login_attempts [integer]
*   userprincipalname [string] - The principle name of the user.
*   member_of [string] - The user's directory membership.
*   title [string] - The user's job title.
*   id [integer]
*   updated_at [string]
*   state [integer]
*   role_ids [list of integers] - A list of OneLogin Role IDs of the user
*   group_id [integer] - The ID of the Group in OneLogin that the user is assigned to.
*   preferred_locale_code [string]
*   directory_id [integer] - The ID of the OneLogin Directory of the user.
*   created_at [string]
*   trusted_idp_id [integer] - The ID of the OneLogin Trusted IDP of the user.
*   company [string] - The user's company.
*   distinguished_name [string] - The distinguished name of the user.
*   activated_at [string]
*   external_id [string] - The ID of the user in an external directory.
*   last_login [string]
*   comment [string] - Free text related to the user.
*   department [string] - The user's department.
*   manager_user_id [string] - The OneLogin User ID for the user's manager.

### onelogin_roles_apps (filters)

The roles_apps data source allows you to retrieve an already existing roles_apps resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

<pre><span>data</span> <span>"onelogin_roles_apps" "my_roles_apps"</span>{
    <span>filter</span> <span>{</span>
        <span>name</span> = <span>"property name to filter by, see docs below for more info about available filter name options"</span>
        <span>values</span> = <span>["filter value"]</span>
    <span>}</span>
<span>}</span></pre>

#### Arguments Reference

The following arguments are supported:

*   filter - (Required) Object containing two properties.

*   name [string]: the name should match one of the properties to filter by. The following property names are supported: <span>icon_url,</span> <span>id,</span> <span>name,</span> <span>roles_id,</span>
*   values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   icon_url [string] - url of Icon
*   id [integer] - app id
*   name [string] - app name

### onelogin_roles_users (filters)

The roles_users data source allows you to retrieve an already existing roles_users resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

<pre><span>data</span> <span>"onelogin_roles_users" "my_roles_users"</span>{
    <span>filter</span> <span>{</span>
        <span>name</span> = <span>"property name to filter by, see docs below for more info about available filter name options"</span>
        <span>values</span> = <span>["filter value"]</span>
    <span>}</span>
<span>}</span></pre>

#### Arguments Reference

The following arguments are supported:

*   filter - (Required) Object containing two properties.

*   name [string]: the name should match one of the properties to filter by. The following property names are supported: <span>invitation_sent_at,</span> <span>firstname,</span> <span>salt,</span> <span>password_changed_at,</span> <span>manager_ad_id,</span> <span>phone,</span> <span>samaccountname,</span> <span>password_algorithm,</span> <span>password_confirmation,</span> <span>password,</span> <span>status,</span> <span>username,</span> <span>locked_until,</span> <span>lastname,</span> <span>email,</span> <span>invalid_login_attempts,</span> <span>userprincipalname,</span> <span>member_of,</span> <span>title,</span> <span>id,</span> <span>updated_at,</span> <span>state,</span> <span>group_id,</span> <span>preferred_locale_code,</span> <span>directory_id,</span> <span>created_at,</span> <span>trusted_idp_id,</span> <span>company,</span> <span>distinguished_name,</span> <span>activated_at,</span> <span>external_id,</span> <span>last_login,</span> <span>comment,</span> <span>department,</span> <span>roles_id,</span> <span>manager_user_id,</span>
*   values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   invitation_sent_at [string]
*   firstname [string] - The user's first name.
*   salt [string] - The salt value used with the password_algorithm.
*   password_changed_at [string]
*   manager_ad_id [string] - The ID of the user's manager in Active Directory.
*   phone [string] - The E.164 format phone number for a user.
*   samaccountname [string] - The user's Active Directory username.
*   password_algorithm [string] - Use this when importing a password that's already hashed. Prepend the salt value to the cleartext password value before SHA-256-encoding it
*   password_confirmation [string] - Required if the password is being set.
*   password [string] - The password to set for a user.
*   status [integer]
*   username [string] - A username for the user.
*   locked_until [string]
*   lastname [string] - The user's last name.
*   email [string] - A valid email for the user.
*   invalid_login_attempts [integer]
*   userprincipalname [string] - The principle name of the user.
*   member_of [string] - The user's directory membership.
*   title [string] - The user's job title.
*   id [integer]
*   updated_at [string]
*   state [integer]
*   role_ids [list of integers] - A list of OneLogin Role IDs of the user
*   group_id [integer] - The ID of the Group in OneLogin that the user is assigned to.
*   preferred_locale_code [string]
*   directory_id [integer] - The ID of the OneLogin Directory of the user.
*   created_at [string]
*   trusted_idp_id [integer] - The ID of the OneLogin Trusted IDP of the user.
*   company [string] - The user's company.
*   distinguished_name [string] - The distinguished name of the user.
*   activated_at [string]
*   external_id [string] - The ID of the user in an external directory.
*   last_login [string]
*   comment [string] - Free text related to the user.
*   department [string] - The user's department.
*   manager_user_id [string] - The OneLogin User ID for the user's manager.

### onelogin_users (filters)

The users data source allows you to retrieve an already existing users resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

<pre><span>data</span> <span>"onelogin_users" "my_users"</span>{
    <span>filter</span> <span>{</span>
        <span>name</span> = <span>"property name to filter by, see docs below for more info about available filter name options"</span>
        <span>values</span> = <span>["filter value"]</span>
    <span>}</span>
<span>}</span></pre>

#### Arguments Reference

The following arguments are supported:

*   filter - (Required) Object containing two properties.

*   name [string]: the name should match one of the properties to filter by. The following property names are supported: <span>invitation_sent_at,</span> <span>firstname,</span> <span>salt,</span> <span>password_changed_at,</span> <span>manager_ad_id,</span> <span>phone,</span> <span>samaccountname,</span> <span>password_algorithm,</span> <span>password_confirmation,</span> <span>password,</span> <span>status,</span> <span>username,</span> <span>locked_until,</span> <span>lastname,</span> <span>email,</span> <span>invalid_login_attempts,</span> <span>userprincipalname,</span> <span>member_of,</span> <span>title,</span> <span>id,</span> <span>updated_at,</span> <span>state,</span> <span>group_id,</span> <span>preferred_locale_code,</span> <span>directory_id,</span> <span>created_at,</span> <span>trusted_idp_id,</span> <span>company,</span> <span>distinguished_name,</span> <span>activated_at,</span> <span>external_id,</span> <span>last_login,</span> <span>comment,</span> <span>department,</span> <span>manager_user_id,</span>
*   values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   invitation_sent_at [string]
*   firstname [string] - The user's first name.
*   salt [string] - The salt value used with the password_algorithm.
*   password_changed_at [string]
*   manager_ad_id [string] - The ID of the user's manager in Active Directory.
*   phone [string] - The E.164 format phone number for a user.
*   samaccountname [string] - The user's Active Directory username.
*   password_algorithm [string] - Use this when importing a password that's already hashed. Prepend the salt value to the cleartext password value before SHA-256-encoding it
*   password_confirmation [string] - Required if the password is being set.
*   password [string] - The password to set for a user.
*   status [integer]
*   username [string] - A username for the user.
*   locked_until [string]
*   lastname [string] - The user's last name.
*   email [string] - A valid email for the user.
*   invalid_login_attempts [integer]
*   userprincipalname [string] - The principle name of the user.
*   member_of [string] - The user's directory membership.
*   title [string] - The user's job title.
*   id [integer]
*   updated_at [string]
*   state [integer]
*   role_ids [list of integers] - A list of OneLogin Role IDs of the user
*   group_id [integer] - The ID of the Group in OneLogin that the user is assigned to.
*   preferred_locale_code [string]
*   directory_id [integer] - The ID of the OneLogin Directory of the user.
*   created_at [string]
*   trusted_idp_id [integer] - The ID of the OneLogin Trusted IDP of the user.
*   company [string] - The user's company.
*   distinguished_name [string] - The distinguished name of the user.
*   activated_at [string]
*   external_id [string] - The ID of the user in an external directory.
*   last_login [string]
*   comment [string] - Free text related to the user.
*   department [string] - The user's department.
*   manager_user_id [string] - The OneLogin User ID for the user's manager.

### onelogin_users_apps (filters)

The users_apps data source allows you to retrieve an already existing users_apps resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

<pre><span>data</span> <span>"onelogin_users_apps" "my_users_apps"</span>{
    <span>filter</span> <span>{</span>
        <span>name</span> = <span>"property name to filter by, see docs below for more info about available filter name options"</span>
        <span>values</span> = <span>["filter value"]</span>
    <span>}</span>
<span>}</span></pre>

#### Arguments Reference

The following arguments are supported:

*   filter - (Required) Object containing two properties.

*   name [string]: the name should match one of the properties to filter by. The following property names are supported: <span>login_id,</span> <span>extension,</span> <span>name,</span> <span>icon_url,</span> <span>id,</span> <span>provisioning_state,</span> <span>provisioning_enabled,</span> <span>users_id,</span> <span>provisioning_status,</span>
*   values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   login_id [integer] - Unqiue identifier for this user and app combination.
*   extension [boolean] - Boolean that indicates if the OneLogin browser extension is required to launch this app.
*   name [string] - The name of the app.
*   icon_url [string] - A url for the icon that represents the app in the OneLogin portal
*   id [integer] - The App ID
*   provisioning_state [string] - If provisioning is enabled this indicates the state of provisioning for the given user.
*   provisioning_enabled [boolean] - Indicates if provisioning is enabled for this app.
*   provisioning_status [string]

### onelogin_users_devices (filters)

The users_devices data source allows you to retrieve an already existing users_devices resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

<pre><span>data</span> <span>"onelogin_users_devices" "my_users_devices"</span>{
    <span>filter</span> <span>{</span>
        <span>name</span> = <span>"property name to filter by, see docs below for more info about available filter name options"</span>
        <span>values</span> = <span>["filter value"]</span>
    <span>}</span>
<span>}</span></pre>

#### Arguments Reference

The following arguments are supported:

*   filter - (Required) Object containing two properties.

*   name [string]: the name should match one of the properties to filter by. The following property names are supported: <span>device_id,</span> <span>default,</span> <span>user_display_name,</span> <span>auth_factor_name,</span> <span>users_id,</span> <span>type_display_name,</span>
*   values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   device_id [string] - MFA device identifier.
*   default [boolean] - true = is user’s default MFA device for OneLogin.
*   user_display_name [string] - Authentication factor display name assigned by users when they register the device.
*   auth_factor_name [string] - Authentication factor name, as it appears to administrators in OneLogin.
*   type_display_name [string] - Authentication factor display name as it appears to users upon initial registration, as defined by admins at Settings > Authentication Factors.

### onelogin_users_v1 (filters)

The users_v1 data source allows you to retrieve an already existing users_v1 resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

<pre><span>data</span> <span>"onelogin_users_v1" "my_users_v1"</span>{
    <span>filter</span> <span>{</span>
        <span>name</span> = <span>"property name to filter by, see docs below for more info about available filter name options"</span>
        <span>values</span> = <span>["filter value"]</span>
    <span>}</span>
<span>}</span></pre>

#### Arguments Reference

The following arguments are supported:

*   filter - (Required) Object containing two properties.

*   name [string]: the name should match one of the properties to filter by. The following property names are supported: <span>invitation_sent_at,</span> <span>firstname,</span> <span>salt,</span> <span>password_changed_at,</span> <span>manager_ad_id,</span> <span>phone,</span> <span>samaccountname,</span> <span>password_algorithm,</span> <span>password_confirmation,</span> <span>password,</span> <span>status,</span> <span>username,</span> <span>locked_until,</span> <span>lastname,</span> <span>email,</span> <span>invalid_login_attempts,</span> <span>userprincipalname,</span> <span>member_of,</span> <span>title,</span> <span>id,</span> <span>updated_at,</span> <span>state,</span> <span>group_id,</span> <span>preferred_locale_code,</span> <span>directory_id,</span> <span>created_at,</span> <span>trusted_idp_id,</span> <span>company,</span> <span>distinguished_name,</span> <span>activated_at,</span> <span>external_id,</span> <span>last_login,</span> <span>comment,</span> <span>department,</span> <span>manager_user_id,</span>
*   values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   invitation_sent_at [string]
*   firstname [string] - The user's first name.
*   salt [string] - The salt value used with the password_algorithm.
*   password_changed_at [string]
*   manager_ad_id [string] - The ID of the user's manager in Active Directory.
*   phone [string] - The E.164 format phone number for a user.
*   samaccountname [string] - The user's Active Directory username.
*   password_algorithm [string] - Use this when importing a password that's already hashed. Prepend the salt value to the cleartext password value before SHA-256-encoding it
*   password_confirmation [string] - Required if the password is being set.
*   password [string] - The password to set for a user.
*   status [integer]
*   username [string] - A username for the user.
*   locked_until [string]
*   lastname [string] - The user's last name.
*   email [string] - A valid email for the user.
*   invalid_login_attempts [integer]
*   userprincipalname [string] - The principle name of the user.
*   member_of [string] - The user's directory membership.
*   title [string] - The user's job title.
*   id [integer]
*   updated_at [string]
*   state [integer]
*   role_ids [list of integers] - A list of OneLogin Role IDs of the user
*   group_id [integer] - The ID of the Group in OneLogin that the user is assigned to.
*   preferred_locale_code [string]
*   directory_id [integer] - The ID of the OneLogin Directory of the user.
*   created_at [string]
*   trusted_idp_id [integer] - The ID of the OneLogin Trusted IDP of the user.
*   company [string] - The user's company.
*   distinguished_name [string] - The distinguished name of the user.
*   activated_at [string]
*   external_id [string] - The ID of the user in an external directory.
*   last_login [string]
*   comment [string] - Free text related to the user.
*   department [string] - The user's department.
*   manager_user_id [string] - The OneLogin User ID for the user's manager.

### onelogin_users_v1_apps (filters)

The users_v1_apps data source allows you to retrieve an already existing users_v1_apps resource using filters. Refer to the arguments section to learn more about how to configure the filters.

#### Example usage

<pre><span>data</span> <span>"onelogin_users_v1_apps" "my_users_v1_apps"</span>{
    <span>filter</span> <span>{</span>
        <span>name</span> = <span>"property name to filter by, see docs below for more info about available filter name options"</span>
        <span>values</span> = <span>["filter value"]</span>
    <span>}</span>
<span>}</span></pre>

#### Arguments Reference

The following arguments are supported:

*   filter - (Required) Object containing two properties.

*   name [string]: the name should match one of the properties to filter by. The following property names are supported: <span>login_id,</span> <span>extension,</span> <span>name,</span> <span>icon_url,</span> <span>id,</span> <span>provisioning_state,</span> <span>provisioning_enabled,</span> <span>users_v1_id,</span> <span>provisioning_status,</span>
*   values [array of string]: Values to filter by (only one value is supported at the moment).

**Note:** If more or less than a single match is returned by the search, Terraform will fail. Ensure that your search is specific enough to return a single result only.

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   login_id [integer] - Unqiue identifier for this user and app combination.
*   extension [boolean] - Boolean that indicates if the OneLogin browser extension is required to launch this app.
*   name [string] - The name of the app.
*   icon_url [string] - A url for the icon that represents the app in the OneLogin portal
*   id [integer] - The App ID
*   provisioning_state [string] - If provisioning is enabled this indicates the state of provisioning for the given user.
*   provisioning_enabled [boolean] - Indicates if provisioning is enabled for this app.
*   provisioning_status [string]