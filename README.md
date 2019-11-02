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
