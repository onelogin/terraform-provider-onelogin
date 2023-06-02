# Terraform Provider OneLogin

This repository contains the code for the new Terraform provider for OneLogin, a work-in-progress project that aims to provide complete functionality for the OneLogin API and address limitations in the existing provider. The provider leverages the Terraform Plugin Framework and is written in Go following HashiCorp's best practices for creating Terraform providers.

## Overview

The Terraform provider for OneLogin serves as a standalone application, communicating with the Terraform core and the OneLogin API. The provider primarily consists of the following components:

- **Provider Configuration**: Handles the authentication mechanism with the OneLogin API. Securely fetches the API token using the `ONELOGIN_CLIENT_ID` and `ONELOGIN_CLIENT_SECRET` environment variables.
- **Data Sources**: Fetches and reads data from OneLogin, each data source is implemented as a Go struct that adheres to the `datasource.DataSource` interface.
- **Resources**: Enables create, read, update, and delete operations for OneLogin entities, each resource is implemented as a Go struct following the `resource.Resource` interface.
- **Error Handling**: Prioritizes meaningful error messages and graceful handling of API failures.

## Engineering Requirements

The main objectives of this project include:

- **Full API Coverage**: The provider must support all endpoints and features exposed by the OneLogin API.
- **Secure Credential Handling**: OneLogin API credentials should be securely managed by the provider by fetching the API token using environment variables.
- **Subdomain Support**: The provider must correctly handle OneLogin subdomains to ensure error-free usage of data sources.
- **Field Inclusion**: All fields exposed by the OneLogin API should be included in the new provider.
- **Error Handling**: The provider must have robust error handling capabilities, providing users with meaningful error messages to aid in issue resolution.
- **Documentation**: The provider's documentation should be comprehensive, including usage examples for every resource and data source.

## Milestones

The general progression of this project is as follows:

1. **API Exploration and Design**: Thorough exploration of the OneLogin API and design of the new provider's architecture.
2. **Initial Development**: Begin coding the new provider, prioritizing high-priority features and addressing known issues in the current provider.
3. **Testing and Iteration**: Extensive testing of the new provider, iterating its design and implementation based on test results.
4. **Documentation and Cleanup**: Creation of comprehensive documentation and resolution of any remaining codebase issues.
5. **Release and Maintenance**: Launch the new provider, maintain it, and incorporate improvements based on user feedback.

## Success Criteria

The success of this project will be evaluated based on the following criteria:

- **Technical Benchmarks**: The provider supports all OneLogin API endpoints and features, handles credentials securely, correctly processes subdomains, includes all fields exposed by the API, and provides robust error handling.
- **User Feedback**: Positive feedback from users regarding the functionality and usability of the new provider.
- **Adoption Rate**: The number of projects adopting the new Terraform provider within six months of its release.
- **Issue Resolution**: Known issues with the existing provider are effectively addressed in the new implementation.
- **Documentation Quality**: The provider's documentation is comprehensive, includes usage examples for every resource and data source, and facilitates easy understanding and usage of the provider.

Please note that the milestones and success criteria may be subject to change based on availability and project constraints.

## Contributing

As this is a work-in-progress project, contributions and feedback are very much welcome. Please open an issue or submit a pull request if you'd like to contribute.

## Updates
6-1-23: In order to provide full functionality on Terraform, I need to rewrite the Onelogin Go SDK. Here is the [work in progress branch](https://github.com/onelogin/onelogin-go-sdk/tree/non-generated-updates)