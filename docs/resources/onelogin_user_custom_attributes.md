# onelogin_user_custom_attributes Resource

This resource allows you to manage custom user attributes in OneLogin.

> **Note:** Due to a bug in the OneLogin API (returning "Missing param: user_field"), creating, updating, and deleting custom attribute definitions must currently be done through the OneLogin UI. Once created, you can use this resource to set the values for those attributes on specific users.

## Example Usage

```hcl
# Create a user
resource onelogin_users test_user {
  username = "test.user"
  email    = "test.user@example.com"
}

# Reference a custom attribute that was created in the OneLogin UI
resource onelogin_user_custom_attributes employee_id_reference {
  name      = "Employee ID"    # For reference only - must match the UI
  shortname = "employee_id"    # For reference only - must match the UI
}

# Set a custom attribute value for a specific user
resource onelogin_user_custom_attributes user_employee_id {
  user_id   = onelogin_users.test_user.id
  shortname = "employee_id"     # Must match an existing attribute in OneLogin
  value     = "EMP12345"
}

# Set another custom attribute value for the same user
resource onelogin_user_custom_attributes user_department_code {
  user_id   = onelogin_users.test_user.id
  shortname = "dept_code"       # Must match an existing attribute in OneLogin
  value     = "IT-DEPT"
}
```

> **Important:** Due to a bug in the OneLogin API, the custom attributes must first be created in the OneLogin UI before they can be used in Terraform.

## Argument Reference

The following arguments are supported:

### For Custom Attribute References

* `name` - (Required) The human-readable name of the custom attribute (for reference only, must be created in the OneLogin UI).
* `shortname` - (Required) The short name (identifier) of the custom attribute (for reference only, must be created in the OneLogin UI).
* `position` - (Optional) The position of the custom attribute in the UI (for reference only).

### For User-Specific Custom Attribute Values

* `user_id` - (Required) The ID of the user to set this custom attribute for.
* `shortname` - (Required) The short name (identifier) of the custom attribute. Must match an existing attribute in OneLogin.
* `value` - (Required) The value to set for this custom attribute.

## Import

> Note: Due to the OneLogin API bug, importing custom attribute definitions is currently limited. You can only import user-specific values.

User-specific custom attribute values can be imported using the format `{user_id}_{shortname}`:

```bash
terraform import onelogin_user_custom_attributes.user_employee_id 789012_employee_id
```

## Attribute Reference

* `id` - The composite ID for user-specific attribute values in the format `{user_id}_{shortname}`.

## Future Enhancements

Once OneLogin fixes the API bug with the "Missing param: user_field" error, this resource will be updated to fully support creating, updating, and deleting custom attribute definitions via Terraform.