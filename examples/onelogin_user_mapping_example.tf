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

# Mappings act on roles, so this example creates the ones it assigns rather than
# naming IDs from somebody else's account. Point these at your own roles if you
# already have them -- the API rejects an ID it does not recognise with
# "Invalid action value(s)".
resource "onelogin_roles" "mapping_primary" {
  name = "Mapping Example Primary acctest"
}

resource "onelogin_roles" "mapping_secondary" {
  name = "Mapping Example Secondary acctest"
}

# Basic user mapping that applies a role to users with email domain @example.com
resource "onelogin_user_mappings" "example_mapping" {
  name    = "Example Domain Mapping"
  match   = "all" # Match all conditions
  enabled = true

  # Condition to check if user's email contains @example.com
  conditions {
    source   = "email"        # User attribute to check
    operator = "~"            # Operator for comparison
    value    = "@example.com" # Value to compare against
  }

  # Action to assign a role to matching users
  actions {
    action = "add_role"
    value  = [onelogin_roles.mapping_primary.id]
  }
}

# More complex user mapping with multiple conditions and actions
resource "onelogin_user_mappings" "department_mapping" {
  name    = "Department Based Mapping"
  match   = "all" # Match all conditions
  enabled = true

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

  actions {
    action = "add_role"
    value  = [onelogin_roles.mapping_primary.id]
  }
}
