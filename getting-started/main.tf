terraform {
  required_providers {
    onelogin = {
      source  = "onelogin/onelogin"
      version = "0.8.0"
    }
  }
}

provider "onelogin" {
  client_id     = "YOUR_CLIENT_ID"     # Replace with your OneLogin Client ID
  client_secret = "YOUR_CLIENT_SECRET" # Replace with your OneLogin Client Secret
  url           = "YOUR_API_URL"       # Replace with your OneLogin API URL (e.g., "https://company.onelogin.com")
}

# Example resource (uncomment if needed)
# resource "onelogin_user" "example" {
#   firstname   = "Example"
#   lastname    = "User"
#   username    = "example.user@example.com"
#   email       = "example.user@example.com"
#   status      = 1
#   department  = "IT"
#   phone       = "123-456-7890"
# }
