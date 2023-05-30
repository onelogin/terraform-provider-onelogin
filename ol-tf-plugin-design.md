# Design Document for Terraform-Provider-OneLogin

## Overview

This document outlines the design for the new Terraform provider for OneLogin, which aims to leverage the Terraform Plugin Framework. The objective is to provide complete functionality for OneLogin's API and address the limitations of the current provider.

## Architecture

The Terraform provider will be developed as a standalone application that communicates with the Terraform core and OneLogin API. It will be implemented in Go, following HashiCorp's recommended practices for creating Terraform providers.

The major components of the provider will include:

1. **Provider Configuration**: This component will handle the authentication mechanism for interacting with the OneLogin API. It will securely fetch the API token using environment variables, specifically `ONELOGIN_CLIENT_ID` and `ONELOGIN_CLIENT_SECRET`.

2. **Data Sources**: Data sources will be designed to fetch and read data from OneLogin. Each data source will be implemented as a Go struct that adheres to the `datasource.DataSource` interface. These data sources will ensure error-free retrieval of data, with proper handling of OneLogin subdomains.

3. **Resources**: Resources will enable create, read, update, and delete operations on OneLogin entities. Each resource will be implemented as a Go struct following the `resource.Resource` interface. All fields exposed by the OneLogin API will be included in the resource implementations.

4. **Error Handling**: The provider will prioritize meaningful error messages and graceful handling of API failures. Error handling capabilities of Go will be utilized to achieve this.

## Provider Configuration

The provider will be initialized with a configuration that sets up the authentication mechanism. This will be done using the `Configure` function provided by the Terraform Plugin Framework. The `Configure` function will fetch the API token securely using environment variables.

```go
type providerConfig struct {
	APIToken string
}

func (c *providerConfig) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Fetch API token using environment variables here
}
```

## Data Sources

Data sources will be designed to fetch data from OneLogin. Each data source will be a Go struct that implements the `datasource.DataSource` interface. The `Read` function of each data source will be responsible for retrieving and processing the data from OneLogin.

```go
type dataSourceOneLoginUser struct{}

func (ds dataSourceOneLoginUser) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Implement data reading logic here
}
```

## Resources

Resources will be designed to manage OneLogin entities. Each resource will be a Go struct that implements the `resource.Resource` interface. The resource struct will have functions to handle create, read, update, and delete operations.

```go
type resourceOneLoginUser struct{}

func (r resourceOneLoginUser) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Implement create operation here
}

func (r resourceOneLoginUser) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Implement read operation here
}

func (r resourceOneLoginUser) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Implement update operation here
}

func (r resourceOneLoginUser) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Implement delete operation here
}
```

## Error Handling

The provider will leverage Go's error handling capabilities to handle errors effectively.

 This includes handling API errors, configuration errors, and other potential issues. Error messages will be meaningful and user-friendly, enabling users to identify and resolve issues easily.

```go
func (r resourceOneLoginUser) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Attempt to read the resource
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource",
			fmt.Sprintf("An error occurred while reading the resource: %s", err),
		)
		return
	}
}
```

## Step-by-step Process of Creating a New Terraform Provider

Here is a detailed step-by-step process for creating a new Terraform provider using the Terraform Plugin Framework:

1. **Create the Provider Server**: Set up the provider server, which encapsulates Terraform plugin details and handles provider, resource, and data source operations. The provider server is implemented as a binary that the Terraform CLI downloads, starts, and stops.

2. **Define the Provider**: Define the provider, which is the top-level abstraction that exposes available resources and data sources for users. The provider can have its own configuration, such as authentication information. It is recommended to focus the provider on a single API or SDK to simplify connectivity and authentication requirements.

3. **Define the Schema**: Define schemas for provider, resource, and provisioner configuration blocks. Schemas describe the available fields and provide metadata to Terraform. The resource and attribute schemas should align closely with the underlying API to enable easy conversion and interoperability with other tools.

4. **Create Resources**: Create resources that allow Terraform to manage infrastructure objects, such as compute instances or access policies. Resources act as a translation layer between Terraform and the API, providing a way to define infrastructure using Terraform's declarative language. Each resource should represent a single API object and support import operations.

5. **Create Data Sources**: Implement data sources that allow Terraform to reference external data. Data sources define how to request external data and convert the response into an interpolatable format. Data sources enable users to fetch and use data from external systems in their Terraform configurations.

6. **Handle Sensitive Values in State**: Implement handling of sensitive values in state by using the sensitive flag in the schema of fields that contain sensitive information. This ensures that sensitive values are not displayed in CLI output or Terraform Cloud.

7. **Consider State and Versioning**: Ensure continuity of state and configurations when releasing a provider. Providers should follow Semantic Versioning 2.0.0 for user state and configurations to maintain compatibility. Care should be taken to avoid breaking changes that could impact existing deployments.

8. **Test Your Provider**: Thoroughly test the provider before releasing it. Unit tests and acceptance tests should be conducted to ensure the functionality and stability of the provider. Unit tests focus on individual functions, while acceptance tests cover end-to-end scenarios.

9. **Release Your Provider**: Once confident in the provider's functionality and stability, release it to users. This involves tagging a release in the version control system, building a release binary, and publishing it to a registry or distribution platform. 