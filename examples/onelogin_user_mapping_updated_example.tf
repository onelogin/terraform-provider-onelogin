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

# Unchanged from the example this updates -- the roles the mappings assign.
resource "onelogin_roles" "mapping_primary" {
  name = "Mapping Example Primary acctest"
}

resource "onelogin_roles" "mapping_secondary" {
  name = "Mapping Example Secondary acctest"
}

# Updated version of the example mapping
resource "onelogin_user_mappings" "example_mapping" {
  name    = "Updated Domain Mapping" # Changed name
  match   = "any"                    # Changed from "all" to "any"
  enabled = true

  # Original condition
  # A domain of its own, not plain @example.com. The user and privilege
  # fixtures create accounts at @example.com, and a mapping that outlives a
  # failed run then rewrites them -- which is exactly what happened here: a
  # leftover set_userprincipalname mapping made TestAccUser_crud fail with a
  # userprincipalname nobody had configured.
  conditions {
    source   = "email"
    operator = "~"
    value    = "@sales.example.com"
  }

  # Added condition
  conditions {
    source   = "email"
    operator = "~"
    value    = "@partner.example.com"
  }

  actions {
    action = "add_role"
    value  = [onelogin_roles.mapping_primary.id]
  }

  # Added action, using a OneLogin macro rather than a Terraform interpolation.
  # The doubled $ is HCL's escape: what reaches OneLogin is ${user.email}.
  actions {
    action = "set_userprincipalname"
    value  = ["$${user.email}"]
  }
}

# Updated version of the department mapping with different approach
resource "onelogin_user_mappings" "department_mapping" {
  name    = "Engineering Team Mapping" # Updated name
  match   = "all"
  enabled = true

  # Simplified condition - now just checking for Engineering department
  conditions {
    source   = "department"
    operator = "="
    value    = "Engineering" # Changed from IT to Engineering
  }

  # Now assigning both roles -- one block each, because add_role takes a single
  # value. Two IDs in one block are rejected with
  # "Invalid action value(s): <both ids>", which names the values when what is
  # actually wrong is the count.
  actions {
    action = "add_role"
    value  = [onelogin_roles.mapping_primary.id]
  }

  actions {
    action = "add_role"
    value  = [onelogin_roles.mapping_secondary.id]
  }
}
