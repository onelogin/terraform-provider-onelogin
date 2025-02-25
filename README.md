# Onelogin Terraform Provider
[![Go Report Card](https://goreportcard.com/badge/github.com/onelogin/terraform-provider-onelogin)](https://goreportcard.com/report/github.com/onelogin/terraform-provider-onelogin)
<a href='https://github.com/dcaponi/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-100%25-brightgreen.svg?longCache=true&style=flat)</a>

# Prerequisites
    1) Install Golang
    2) Install Terraform v0.12.24 or later
    3) Install gosec (for security scanning):
       ```
       curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b $(go env GOPATH)/bin latest
       ```
    4) Install code dependencies

# Getting Started W/ Local Testing & Development
If you are sideloading this provider (i.e. not getting this via the Terraform store) You must clone this repository to run the following commands.

1) In the ./terraform-provider-onelogin directory run:
    ```
    make sideload
    ```

    If you are using Terraform v0.13.x or later you can use following Terraform configuration for sideloaded version of this provider:
    ```
    terraform {
      required_providers {
        onelogin = {
          source  = "onelogin.com/onelogin/onelogin"
          version = "0.1.10"
        }
      }
    }

    provider "onelogin" {
      # Configuration options
    }
    ```

2) You'll need admin access to a OneLogin account where you can create API credentials. Create a set of API credentials with _manage all_ permission. For applying the credentials, there are 2 ways

    * Export these credentials to your environment and the provider will read them in from there
    ```
    export ONELOGIN_CLIENT_ID=<your client id>
    export ONELOGIN_CLIENT_SECRET=<your client secret>
    export ONELOGIN_OAPI_URL=<the api url for your region>
    ```

3) You are ready to use the provider, just follow the terraform commands!

# Shipping Code
### Development
1. Adding a new resource generally requires that service to be defined in the [OneLogin SDK](https://github.com/onelogin/onelogin-go-sdk) see `AppsService` for an example.

2. Define the new resource in `onelogin/provider.go` in a similar fashion to how the other resources are identified there.

3. `onelogin/resource_onelogin_<resource>.go` and `onelogin/resource_onelogin<resource>_test.go` are required. See existing code for examples. This layer is the interface to the Terraform and OneLogin SDKs for making the actual requests to OneLogin. The accompanying test file is used by Terraform for running the acceptance tests.

4. `ol_schemas/<resource>/<sub-resource>` contains the logic for packing & unpacking resources between json and their golang struct definitions.

5. Add examples in `examples/onelogin_<resource>_example.tf` `examples/onelogin_<resource>_updated_example.tf` to be used by the acceptance tests to ensure applications happen correctly.

6. Add a doc page to `docs/resources/onelogin_<resource>.md`

To debug / troubleshoot, set an environment variable `export TF_LOG=trace` to see the output of any loggers in the Terraform workflow. 

### Deployment
1. Open a PR against `develop` branch. Once approved and CI/CD pass merge it to `develop` via github.

2. Once ready to ship cut a `release` off of `develop`. Release branches should follow the naming convention `vX.X.XX` so if you use `git flow` cut one with `git flow release start v0.0.00`

3. Do a smoke test and any last minute updates then merge the release into both `master` and `develop` and tag the merge commit. if you use `git flow ` its `git flow release finish v0.0.00` (This also tags your release commit).

4. Push everything to github. From `develop` you can run `git push && git checkout master && git push && git push --tags`

5. The new tag will trigger the release action which makes builds for the OSes described in the `release.yaml`.

6. Once the release action completes, go to the [Releases](https://github.com/onelogin/terraform-provider-onelogin/releases) section of the repository in github and look for the draft release with your version number.

7. Click on that and ensure the build artifacts were uploaded to the release. Once you've verified this, click "Edit Draft" and "Publish Release".

# Terraform for non-users

Wildly simplified explanation - it's a thing that lets you describe the final state of all your OneLogin things (apps, users, associations via roles etc) via a `.tf` file using HashiCorp Language (HCL) and users that description to fire off a bunch of API requests via OneLogin APIs to make that desired state a reality. Also tracks the known state of OneLogin in `.tfstate` and users that as the source of truth.


**Install**:
```
brew install terraform
```

**Initialize** sets up the tfstate and prepares to track:
 ```
terraform init
```

**Plan** Shows the diff between current and desired state:
```
terraform plan
```

**Apply** does all the actual work of updating OneLogin:
```
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