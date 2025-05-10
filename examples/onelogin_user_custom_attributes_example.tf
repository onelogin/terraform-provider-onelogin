# Create a user
resource onelogin_users test_user {
  username = "test.user"
  email    = "test.user@example.com"
}

# Reference a custom attribute that was created in the OneLogin UI
# Note: Due to an API bug, custom attributes must first be created in the OneLogin UI
resource onelogin_user_custom_attributes employee_id_reference {
  name      = "Employee ID"    # For reference only - must match the UI
  shortname = "employee_id"    # For reference only - must match the UI
}

# Reference another custom attribute from the OneLogin UI
resource onelogin_user_custom_attributes department_code_reference {
  name      = "Department Code"    # For reference only - must match the UI
  shortname = "dept_code"          # For reference only - must match the UI
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