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
  shortname = "dept_code"       # Must match an existing attribute in OneLogin
  value     = "IT-DEPT"
}