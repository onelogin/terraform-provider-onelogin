terraform {
  required_providers {
    onelogin = {
      source  = "onelogin/onelogin"
      version = ">= 0.8.0"
    }
  }
}

provider "onelogin" {
  # Set these variables with ONELOGIN_CLIENT_ID, ONELOGIN_CLIENT_SECRET environment variables
  # Or provide them directly here (not recommended for sensitive values)
}

# Basic user mapping that applies a role to users with email domain @example.com
resource "onelogin_user_mappings" "example_mapping" {
  name     = "Example Domain Mapping"
  match    = "all"                  # Match all conditions
  enabled  = true

  # Condition to check if user's email contains @example.com
  conditions {
    source   = "email"              # User attribute to check
    operator = "~"           # Operator for comparison
    value    = "@example.com"       # Value to compare against
  }

  # Action to assign a role to matching users
  # Note: Replace 12345 with actual role ID from your OneLogin account
  actions {
    action = "add_role"
    value  = ["380586"]
  }
}

# More complex user mapping with multiple conditions and actions
resource "onelogin_user_mappings" "department_mapping" {
  name     = "Department Based Mapping"
  match    = "all"                  # Match all conditions
  enabled  = true

  # Check if user belongs to IT department
  conditions {
    source   = "department"
    operator = "="
    value    = "IT"
  }

  # Check if user's title contains "Engineer"
  conditions {
    source   = "title"
    operator = "~"
    value    = "Engineer"
  }

  # Assign multiple roles to matching users
  # Note: Replace IDs with actual role IDs from your OneLogin account
  actions {
    action = "add_role"
    value  = ["380586"]
  }

  # Set the user's group memberships
}