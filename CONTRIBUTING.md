# Contributing to OneLogin Terraform Provider

Thank you for your interest in contributing! This guide will help you set up your development environment and understand our development workflow.

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
   export ONELOGIN_API_URL=<your OneLogin API URL, e.g., https://company.onelogin.com>
   ```
3. Build and install the provider locally:
   ```bash
   make sideload
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

### Helpful Makefile Commands

```bash
# Build and install locally
make sideload

# Clean terraform state
make clean-terraform

# Run tests (skips acceptance tests)
make test

# Run security checks
make secure

# Run acceptance tests (creates real resources)
make testacc
```

## Dependency Management

We use go mod for dependency management.

To add a package:
```bash
go get -u "package-name"
```

To re-install dependencies:
```bash
rm go.sum
go mod download
```

To update dependencies:
```bash
go mod -u ./...
```

## Release Process

To create a new release, simply publish a GitHub Release:

1. Go to the [Releases page](../../releases) in GitHub
2. Click **"Draft a new release"**
3. Click **"Choose a tag"** and create a new tag following semantic versioning (e.g., `v0.11.1`)
4. Set the release title and description (you can use "Generate release notes" for automatic changelog)
5. Click **"Publish release"**

The Release workflow will automatically:
- Build the provider binaries with GoReleaser (using the tag version)
- Generate checksums and sign them with GPG
- Attach binaries and checksums to the GitHub release
- Publish to the Terraform Registry

**That's it!** The entire release process is automated from a single GitHub Release creation.

## Code Quality

- All code must pass `gosec` security scanning
- Unit tests are required for new functionality
- Aim for 100% test coverage where practical
- Follow Go best practices and idioms
