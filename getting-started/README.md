# Getting Started with OneLogin Terraform Provider

This directory contains a minimal setup to get started with the OneLogin Terraform provider.

## Prerequisites

- [Terraform](https://www.terraform.io/downloads.html) installed (v0.13+)
- OneLogin account with API credentials
- API credentials with required permissions

## Setup

1. **Obtain OneLogin API Credentials**
   - Log in to your OneLogin admin portal
   - Navigate to Developers > API Credentials
   - Create credentials with appropriate permissions
   - Note the Client ID, Client Secret, and your subdomain

2. **Configure the Provider**
   - Open `main.tf` 
   - Replace `YOUR_CLIENT_ID`, `YOUR_CLIENT_SECRET`, and `YOUR_SUBDOMAIN` with your actual credentials
   - The subdomain is the part before `.onelogin.com` in your OneLogin URL

3. **Initialize Terraform**
   ```
   terraform init
   ```

4. **Apply Configuration**
   ```
   terraform apply
   ```

## Example Resources

The `main.tf` file includes a commented example of creating a OneLogin user. Uncomment it to use, or consult the [documentation](https://registry.terraform.io/providers/onelogin/onelogin/latest/docs) for more resource types.

## Security Note

Do not commit your API credentials to version control. Consider using environment variables or another secure method to manage secrets.

Example with environment variables:
```hcl
provider "onelogin" {
  client_id     = "${env("ONELOGIN_CLIENT_ID")}"
  client_secret = "${env("ONELOGIN_CLIENT_SECRET")}"
  subdomain     = "${env("ONELOGIN_SUBDOMAIN")}"
}
```