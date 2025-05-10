# onelogin_user_custom_attributes Resource

This resource allows you to manage custom user attributes in OneLogin. You can:

1. Create, read, update, and delete custom attribute definitions
2. Set, update, and remove custom attribute values for specific users

## Example Usage

```hcl
# Create a user
resource onelogin_users test_user {
  username = "test.user"
  email    = "test.user@example.com"
}

# Create a custom attribute definition (schema)
resource onelogin_user_custom_attributes employee_id_definition {
  name      = "Employee ID"    # Display name shown in the UI
  shortname = "employee_id"    # Technical name/identifier for the attribute
}

# Create another custom attribute definition
resource onelogin_user_custom_attributes department_definition {
  name      = "Department Code"
  shortname = "dept_code"
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
  shortname = "dept_code"       # Must match an existing attribute
  value     = "IT-DEPT"
}
```

## Argument Reference

The following arguments are supported:

### For Custom Attribute Definitions

* `name` - (Required) The human-readable display name of the custom attribute.
* `shortname` - (Required) The short name (identifier) of the custom attribute.

### For User-Specific Custom Attribute Values

* `user_id` - (Required) The ID of the user to set this custom attribute for.
* `shortname` - (Required) The short name (identifier) of the custom attribute. Must match an existing attribute in OneLogin.
* `value` - (Required) The value to set for this custom attribute.

## Import

### Custom Attribute Definitions

Custom attribute definitions can be imported using the attribute ID:

```bash
terraform import onelogin_user_custom_attributes.employee_id_definition attr_12345
```

### User-Specific Custom Attribute Values

User-specific custom attribute values can be imported using the format `{user_id}_{shortname}`:

```bash
terraform import onelogin_user_custom_attributes.user_employee_id 789012_employee_id
```

## Attribute Reference

* `id` - For attribute definitions, the ID with the format `attr_{id}`.
* `id` - For user-specific attribute values, the composite ID with the format `{user_id}_{shortname}`.