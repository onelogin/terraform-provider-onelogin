#Onelogin Terraform Provider SDK

#Getting Started
To run app: "go run cmd/main.go"

#Dependency Management
We use dep for dependency management.

To install: 

```
brew install dep
```

To install dependency into vendor folder:
```
dep ensure
```

To add a package:
```
dep ensure -add "package"
```

To update package:
```
dep ensure -update
```

To check status:
```
dep status
```

#Folder Structure

    /cmd
        Main applications for project (main file for the app)
    /internal
        Internal packages for current app
    /pkg
        Packages available for external apps
    /vendor
        Application dependencies


#Tests