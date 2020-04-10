# Onelogin Terraform Provider SDK
Custom terraform provider for onelogin

# Prerequisites
    1) Install Golang
    2) Install Terraform v0.12.2
    3) Install code dependencies

# Getting Started
1) In the ./onelogin-terraform-provider folder run:
    ```
    make compile
    ```
2) You are ready to use the provider, just follow the terraform commands!

# Managing App Resources
Refer to [creating an App](https://developers.onelogin.com/api-docs/2/apps/create-app)
On create, omitted fields are ignored and set to their empty or default values.

On update, omitted fields are treated as if the intent is to clear the field and
an empty or zero value is sent to the API to clear the field. E.G. creating an app with
a description, then removing the description field in your HCL file, will result in
setting the description to `""`

### Required Fields
1) name [string]
2) connector_id [int]

### Optional Fields:
1) description [string]
2) notes [string]
3) visible [bool] - Defaults to `true`
4) allow_assumed_signin [bool] - Defaults to `false`
5) parameters [set] - see below
6) configuration [set] - see below
7) provisioning [set] - see below

### Computed Fields
These fields cannot be set via Terraform and are returned by the API
1) auth_method
2) icon_url
3) policy_id
4) tab_id
5) updated_at
6) created_at

# Terraform
Install:
```
brew install terraform
```

Initialize:
 ```
terraform init
```

Plan:
```
terraform plan
```

Apply:
```
terraform apply
```

# Dependency Management
We use go mod for dependency management.

To add a package:

```
go get -u "package-name"
```

# Folder Structure

    /cmd
        Main applications for project (main file for the app)
