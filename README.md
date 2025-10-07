# OneLogin Terraform Provider

[![Go Report Card](https://goreportcard.com/badge/github.com/onelogin/terraform-provider-onelogin)](https://goreportcard.com/report/github.com/onelogin/terraform-provider-onelogin)
<a href='https://github.com/dcaponi/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-100%25-brightgreen.svg?longCache=true&style=flat)</a>

Manage your OneLogin resources with Terraform! This official provider allows you to configure users, groups, roles, applications, and more using infrastructure as code.

## Installation

The OneLogin provider is available on the [Terraform Registry](https://registry.terraform.io/providers/onelogin/onelogin/latest). Terraform will automatically download it when you run `terraform init`.

Add the provider to your Terraform configuration:

```hcl
terraform {
  required_providers {
    onelogin = {
      source  = "onelogin/onelogin"
      version = "~> 1.0"  # Use the latest version from the Terraform Registry
    }
  }
}

provider "onelogin" {
  # Configuration options
}
```

## Authentication

The provider requires OneLogin API credentials. You can configure these via environment variables or provider configuration.

### Option 1: Environment Variables (Recommended)

```bash
export ONELOGIN_CLIENT_ID="your_client_id"
export ONELOGIN_CLIENT_SECRET="your_client_secret"
export ONELOGIN_API_URL="https://your-subdomain.onelogin.com"
```

### Option 2: Provider Configuration

```hcl
provider "onelogin" {
  client_id     = "your_client_id"
  client_secret = "your_client_secret"
  url           = "https://your-subdomain.onelogin.com"
}
```

### Getting API Credentials

1. Log in to your OneLogin admin portal
2. Go to **Developers** → **API Credentials**
3. Create a new API credential with appropriate permissions
4. Save the Client ID and Client Secret

## Quick Start Example

Here's a simple example that creates a user and assigns them to a role:

```hcl
# Create a role
resource "onelogin_roles" "developers" {
  name = "Developers"
}

# Create a user
resource "onelogin_users" "john_doe" {
  username   = "john.doe@example.com"
  email      = "john.doe@example.com"
  firstname  = "John"
  lastname   = "Doe"
}

# Create a group
resource "onelogin_groups" "engineering" {
  name = "Engineering"
}
```

## Available Resources

The provider supports the following OneLogin resources:

- `onelogin_users` - Manage users
- `onelogin_groups` - Manage groups
- `onelogin_roles` - Manage roles
- `onelogin_apps` - Manage applications
- `onelogin_saml_apps` - Manage SAML applications
- `onelogin_oidc_apps` - Manage OIDC applications
- `onelogin_app_rules` - Manage application provisioning rules
- `onelogin_app_role_attachments` - Attach roles to applications
- `onelogin_auth_servers` - Manage OAuth authorization servers
- `onelogin_privileges` - Manage custom privileges
- `onelogin_user_mappings` - Manage user attribute mappings
- `onelogin_user_custom_attributes` - Manage custom user attributes
- `onelogin_smarthooks` - Manage SmartHooks
- `onelogin_smarthook_environment_variables` - Manage SmartHook environment variables
- `onelogin_self_registration_profiles` - Manage self-registration profiles

## Available Data Sources

Use data sources to reference existing OneLogin resources:

- `onelogin_user` - Look up a single user
- `onelogin_users` - Query multiple users
- `onelogin_group` - Look up a single group
- `onelogin_groups` - Query multiple groups

## Documentation

For detailed documentation on each resource and data source, see:

- [Terraform Registry Documentation](https://registry.terraform.io/providers/onelogin/onelogin/latest/docs)
- [Examples](./examples/) - Example configurations for common use cases

## Support

- **Issues**: Report bugs or request features via [GitHub Issues](https://github.com/onelogin/terraform-provider-onelogin/issues)
- **Questions**: For questions about using the provider, please use GitHub Discussions or OneLogin support channels

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for details on:

- Setting up your development environment
- Running tests
- Submitting pull requests
- Release process

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
