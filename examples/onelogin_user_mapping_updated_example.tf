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
resource "onelogin_user_mappings" "example_mapping" {
  name     = "Updated Domain Mapping"  # Changed name
  match    = "any"                     # Changed from "all" to "any"
  enabled  = true

  # Original condition
  conditions {
    source   = "email"
    operator = "~"
    value    = "@example.com"
  }

  # Added condition
  conditions {
    source   = "email"
    operator = "~" 
    value    = "@partner.com"
  }

  # Original action with updated role IDs
  actions {
    action = "add_role"
    value  = ["380586"]                # Updated role ID
  }
  
  # Added action to set custom attributes
  actions {
    action = "set_userprincipalname"
    value  = ["$${user.email}"]        # Dynamic value using user's email
  }
}

# Updated version of the department mapping with different approach
resource "onelogin_user_mappings" "department_mapping" {
  name     = "Engineering Team Mapping"  # Updated name
  match    = "all"
  enabled  = true

  # Simplified condition - now just checking for Engineering department
  conditions {
    source   = "department"
    operator = "="
    value    = "Engineering"             # Changed from IT to Engineering
  }

  # Updated role assignments
  actions {
    action = "add_role"
    value  = ["34567", "56789"]          # Updated role IDs
  }
}