# Onelogin Terraform Provider
[![Go Report Card](https://goreportcard.com/badge/github.com/onelogin/terraform-provider-onelogin)](https://goreportcard.com/report/github.com/onelogin/terraform-provider-onelogin)
<a href='https://github.com/dcaponi/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-100%25-brightgreen.svg?longCache=true&style=flat)</a>

## Latest Updates

### v0.6.0 - Version Alignment

This version aligns the version number with the GitHub releases:

- Version number synchronized with GitHub releases (current v0.5.4)

This version fixes the custom attribute support in the OneLogin v4 API:

- Fixed custom attribute creation by wrapping payload with `user_field` object
- **Now supports** creating, reading, updating, and deleting custom attribute definitions
- Improved user management with custom attributes
- Updated provider to use subdomain instead of region for API connections
- See examples in `examples/onelogin_user_custom_attributes_example.tf`

### v0.1.10 - Custom User Attributes Support

This version includes support for Custom User Attributes using the OneLogin v4 API:

- Added new resource `onelogin_user_custom_attributes` for setting values of existing custom user attributes
- Updated to OneLogin Go SDK v4.1.0
- Improved user management with custom attributes

## Prerequisites
1. Install Go 1.18 or later
2. Install Terraform v0.13.x or later
3. Install gosec (for security scanning):
   ```bash
   curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.18.2
   ```

## Development Setup

1. Clone this repository
2. Set up your OneLogin API credentials:
   ```bash
   export ONELOGIN_CLIENT_ID=<your client id>
   export ONELOGIN_CLIENT_SECRET=<your client secret>
   export ONELOGIN_SUBDOMAIN=<your OneLogin subdomain>
   ```
3. Build and install the provider locally:
   ```bash
   make sideload
   ```

### Example Provider Configuration
```hcl
terraform {
  required_providers {
    onelogin = {
      source  = "onelogin.com/onelogin/onelogin"
      version = "0.6.0"
    }
  }
}

provider "onelogin" {
  # Configuration options
}
```

## Development Workflow

### Adding a New Resource
1. Add the service to the [OneLogin SDK](https://github.com/onelogin/onelogin-go-sdk) (see `AppsService` for example)
2. Define the resource in `onelogin/provider.go`
3. Create resource files:
   - `onelogin/resource_onelogin_<resource>.go`
   - `onelogin/resource_onelogin<resource>_test.go`
4. Add schema definitions in `ol_schemas/<resource>/<sub-resource>`
5. Add examples in `examples/`
6. Add documentation in `docs/resources/`

### Testing
- Run unit tests: `make test`
- Run security checks: `make secure`
- Run acceptance tests: `make testacc` (requires API credentials)
- Debug with: `export TF_LOG=trace`

### Helpful Commands
```bash
# Build and install locally
make sideload

# Clean terraform state
make clean-terraform

# Run tests (skips acceptance tests)
make test

# Run security checks
make secure
```

## Release Process
1. Create feature branch from `main`
2. Create PR against `main`
3. After approval and CI passing, merge to `main`
4. Create a new release either:
   - Through GitHub UI:
     1. Go to "Releases" on GitHub
     2. Click "Draft a new release"
     3. Create a new tag (e.g., v0.5.2)
     4. Fill in release details
     5. Click "Publish release"
   - Or via command line:
     ```bash
     git tag vX.X.X
     git push origin vX.X.X
     ```
5. The GitHub Action will automatically:
   - Build the provider
   - Create a GitHub release
   - Publish to the Terraform Registry

Note: Tags should follow semantic versioning (e.g., v0.5.2)

## Terraform Overview
Terraform enables declarative infrastructure management using HashiCorp Configuration Language (HCL). It tracks the desired state in `.tf` files and the current state in `.tfstate` files.

Basic commands:
```bash
# Initialize working directory
terraform init

# Preview changes
terraform plan

# Apply changes
terraform apply
```

# Dependency Management
We use go mod for dependency management.

To add a package:

```
go get -u "package-name"
```

To re-install dependencies for this project:
```
rm go.sum
go mod download
```

To update dependencies for this project:
```
go mod -u ./...
```

# Helpful Makefile Commands

**testacc** runs acceptance tests (actually creates resources in OL then cleans them up)
```
make testacc
```

**sideload** builds and sideloads the provider for local dev/testing
```
make sideload
```

**clean-terraform** reset terraform state in the local folder
```
make clean-terraform
```

**test** runs unit tests (non-acceptance and no real requests made) and applies coverage badge
```
make test
```

**secure** runs gosec code analysis to warn about possible exploits specific to go
```
make secure
```
