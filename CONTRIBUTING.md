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

A release is created by publishing a GitHub Release:

1. Go to the [Releases page](../../releases) in GitHub
2. Click **"Draft a new release"**
3. Click **"Choose a tag"** and create a new tag following semantic versioning (e.g., `v0.11.1`)
4. Set the release title and description (you can use "Generate release notes" for automatic changelog)
5. **Save it as a draft while you write the notes.** A draft creates no tag and starts
   nothing, so there is no rush and no half-finished release visible to anyone.
6. Click **"Publish release"** when the notes are ready

Publishing starts the Release workflow, which builds the binaries with GoReleaser,
signs the checksums with GPG, and attaches everything to the release.

### Always verify a release before considering it done

Publishing is what starts the build, so a release necessarily exists for a minute or
two before it has any binaries. That is normal. What is not normal is it staying that
way, and a release with no binaries looks perfectly healthy on GitHub while being
uninstallable.

After publishing, confirm both of these:

```bash
# 1. the release has its artifacts: 5 platform zips, SHA256SUMS, and the .sig
gh release view v0.14.0 --json assets --jq '.assets[].name'

# 2. the Terraform Registry has actually picked the version up
curl -s https://registry.terraform.io/v1/providers/onelogin/onelogin/versions \
  | jq -r '.versions[].version' | sort -V | tail -3
```

The Registry ingests from the release-published event and does not keep retrying
indefinitely. If the build is slow or fails, the Registry can conclude there is
nothing to publish and never look again — so step 2 is not redundant with step 1.

### If a release ends up with no binaries

Rebuild the tag without touching the release:

**Actions → Release → Run workflow**, and give it the tag (e.g. `v0.14.0`). The tag is
a required input because a dispatch otherwise builds the branch it was started from.

Then re-check the Registry. If the release now has its artifacts but the Registry still
does not list the version, the ingest event was missed and needs to be sent again.
Delete the release and re-create it against the same tag:

```bash
gh release delete v0.14.0 --yes          # keeps the tag; do NOT pass --cleanup-tag
gh release create v0.14.0 --verify-tag --title v0.14.0 --notes-file notes.md
```

Keep the notes in a file first — deleting the release deletes its body along with its
assets. Re-creating fires a fresh event, the workflow rebuilds, and the Registry picks
it up a few minutes later.

*(Both situations occurred on 2026-08-06 during a GitHub Actions incident: v0.13.0
published, its build was cancelled, and the release sat with no binaries — invisible to
`terraform init` — for nine hours.)*

## Code Quality

- All code must pass `gosec` security scanning
- Unit tests are required for new functionality
- Aim for 100% test coverage where practical
- Follow Go best practices and idioms
