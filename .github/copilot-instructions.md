# OneLogin Terraform Provider

Always reference these instructions first and fallback to search or bash commands only when you encounter unexpected information that does not match the info here.

## Working Effectively

### Bootstrap, Build, and Test the Repository:
- `go mod download` -- downloads dependencies, takes ~40 seconds. Set timeout to 120+ seconds.
- `make build` -- compiles the provider binary, takes ~25 seconds. Set timeout to 60+ seconds.
- `make test` -- runs unit tests (non-acceptance), takes ~18 seconds. Set timeout to 60+ seconds.
- `make sideload` -- builds and installs locally for Terraform testing, takes ~1 second after build.

### Prerequisites:
- Go 1.23.0+ (current version: go1.24.6)
- Terraform v0.13.x or later (install with: `curl -O https://releases.hashicorp.com/terraform/1.5.7/terraform_1.5.7_linux_amd64.zip && unzip terraform_1.5.7_linux_amd64.zip && sudo mv terraform /usr/local/bin/`)
- Optional: gosec for security scanning (may have connectivity issues)

### Development Workflow:
1. **Always run bootstrapping steps first:** `go mod download && make build`
2. **Before committing:** Run `go fmt ./...` (~0.3s) and `go vet ./...` (~5s) 
3. **For local testing:** `make sideload` then `terraform init` in your test directory
4. **Run tests frequently:** `make test` for quick validation

### Long-Running Commands (NEVER CANCEL):
- **Acceptance tests:** `make testacc` -- takes up to 120 minutes. NEVER CANCEL. Set timeout to 150+ minutes.
  - Requires OneLogin API credentials: `ONELOGIN_CLIENT_ID`, `ONELOGIN_CLIENT_SECRET`, `ONELOGIN_API_URL`
  - Actually creates/modifies real OneLogin resources, then cleans them up
  - Only run when you have valid API credentials and want full integration testing

## Validation

### Always validate changes with these steps:
1. **Build validation:** `make build` - must succeed without errors
2. **Unit test validation:** `make test` - all tests must pass
3. **Code formatting:** `go fmt ./...` and `go vet ./...` - must pass for CI
4. **Provider installation:** `make sideload` - must install without errors
5. **Terraform integration:** Create a simple `.tf` file and run `terraform init` - must find provider

### Manual Testing Scenarios:
- **Provider loading test:** Create a test `.tf` file with OneLogin provider block and run `terraform init`
- **Resource validation:** Use examples from `examples/` directory to test specific resources
- **Schema validation:** Run unit tests in `ol_schema/` to verify schema definitions

### CI Requirements:
The GitHub Actions workflow (`.github/workflows/go.yml`) runs:
- `go fmt ./...` 
- `go build ./...`
- `go vet ./...`
- `make test`
- `make secure` (gosec security scan - may fail due to connectivity)

Always run these locally before pushing to ensure CI passes.

## Common Tasks

### Repository Structure:
```
/
├── main.go                    # Provider entry point
├── GNUmakefile               # Build commands
├── go.mod                    # Go module dependencies
├── onelogin/                 # Core provider implementation
│   ├── provider.go           # Provider configuration
│   ├── resource_*.go         # Individual resource implementations
│   └── *_test.go            # Acceptance tests (require API credentials)
├── ol_schema/                # Terraform schema definitions
│   └── */                   # Schema packages for each resource type
├── examples/                 # Example Terraform configurations
├── docs/                     # Provider documentation
├── getting-started/          # Quick start guide and example
└── .github/workflows/        # CI/CD pipelines
```

### Key Files and Their Purposes:
- **onelogin/provider.go:** Defines all available resources and data sources
- **onelogin/resource_onelogin_*.go:** Individual resource CRUD operations
- **ol_schema/*:** Terraform schema definitions for resources
- **examples/*.tf:** Working examples for each resource type
- **getting-started/main.tf:** Basic provider configuration example

### Adding New Resources:
1. Add the resource to `onelogin/provider.go` in the `ResourcesMap`
2. Create `onelogin/resource_onelogin_<resource>.go` with CRUD operations
3. Create `onelogin/resource_onelogin_<resource>_test.go` with acceptance tests
4. Add schema in `ol_schema/<resource>/` directory
5. Add example in `examples/onelogin_<resource>_example.tf`
6. Add documentation in `docs/resources/onelogin_<resource>.md`

### Testing Patterns:
- **Unit tests:** Use testify assertions, focus on schema validation and data transformation
- **Acceptance tests:** Use `resource.Test()` framework, require `TF_ACC=1` environment variable
- **All tests use:** Standard Go testing with table-driven test patterns

### Build Output:
- Binary: `./dist/terraform-provider-onelogin` (~23MB)
- Installed to: `~/.terraform.d/plugins/onelogin.com/onelogin/onelogin/0.8.5/linux_amd64/`

### Common Commands:
```bash
# Quick development cycle
make build && make test

# Full local validation
go fmt ./... && go vet ./... && make build && make test

# Install for local Terraform testing  
make sideload

# Initialize Terraform with local provider
terraform init

# Clean up Terraform state (use with caution - files may not exist)
make clean-terraform

# Clean build artifacts
make clean
```

### Environment Variables for Testing:
```bash
export ONELOGIN_CLIENT_ID=<your_client_id>
export ONELOGIN_CLIENT_SECRET=<your_client_secret>
export ONELOGIN_API_URL=<your_api_url>  # e.g., https://company.onelogin.com
export TF_ACC=1  # Enables acceptance tests
export TF_LOG=trace  # Enables debug logging
```

### Current Version: 0.8.5
Update the `VERSION` in `GNUmakefile` when releasing new versions.

## Troubleshooting

### Common Issues:
- **"make secure" fails:** gosec may have connectivity issues, this is known and can be ignored during development
- **"make clean-terraform" fails:** This is expected when terraform files don't exist, ignore the error
- **Acceptance tests fail:** Check that OneLogin API credentials are properly set and valid
- **Provider not found:** Run `make sideload` to reinstall the provider locally
- **Build fails:** Run `go mod download` to refresh dependencies

### Performance Notes:
- Initial `go mod download`: ~40 seconds
- Build from scratch: ~25 seconds  
- Incremental builds: ~1 second
- Unit tests: ~18 seconds
- Acceptance tests: Up to 120 minutes (NEVER CANCEL)