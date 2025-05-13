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
}

# Updated version of the example mapping
resource "onelogin_user_mapping" "example_mapping" {
  name     = "Updated Domain Mapping"  # Changed name
  match    = "any"                     # Changed from "all" to "any"
  enabled  = true
  position = 3                         # Changed position

  # Original condition
  conditions {
    source   = "email"
    operator = "contains"
    value    = "@example.com"
  }

  # Added condition
  conditions {
    source   = "email"
    operator = "contains" 
    value    = "@partner.com"
  }

  # Original action with updated role IDs
  actions {
    action = "set_role"
    value  = ["45678"]                # Updated role ID
  }
  
  # Added action to set custom attributes
  actions {
    action = "set_userprincipalname"
    value  = ["${user.email}"]        # Dynamic value using user's email
  }
}

# Updated version of the department mapping with different approach
resource "onelogin_user_mapping" "department_mapping" {
  name     = "Engineering Team Mapping"  # Updated name
  match    = "all"
  enabled  = true
  position = 4                           # Updated position

  # Simplified condition - now just checking for Engineering department
  conditions {
    source   = "department"
    operator = "equals"
    value    = "Engineering"             # Changed from IT to Engineering
  }

  # Set specific custom attribute for engineers
  actions {
    action = "set_custom_attribute"
    value  = ["employee_type", "technical"]
  }

  # Updated role assignments
  actions {
    action = "set_role"
    value  = ["34567", "56789"]          # Updated role IDs
  }
}