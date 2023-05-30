# Terraform Provider Onelogin Engineering Requirements

## Overview

The objective of this project is to develop a Terraform provider for OneLogin that fully utilizes the Terraform Plugin Framework. The new provider aims to offer complete functionality for OneLogin's API, addressing the limitations of the current provider which uses the Terraform Provider OpenAPI tool.

## Service Area

The service area of the new Terraform provider will cover the full functionality of OneLogin's API, as detailed in the [OneLogin API documentation](https://developers.onelogin.com/api-docs/2/getting-started/dev-overview). Key features to be implemented include support for all API endpoints, secure handling of OneLogin API credentials, correct processing of OneLogin subdomains, inclusion of all fields exposed by the OneLogin API, robust error handling, and comprehensive documentation.

## Known Issues with the Existing Implementation and Solutions

Based on user feedback and reported issues, the following problems in the current Terraform provider for OneLogin have been identified along with their corresponding solutions in the new implementation:

1. **Subdomain setting issue**: Users have experienced failures when attempting to set the subdomain in the provider, resulting in data source usage issues. The new provider will ensure correct processing of OneLogin subdomains to allow error-free usage of data sources.

2. **Provider credential source**: The existing provider requires a separate manual REST operation to retrieve an API token using OneLogin API credentials. The new provider will securely manage OneLogin API credentials and fetch the API token using environment variables `ONELOGIN_CLIENT_ID` and `ONELOGIN_CLIENT_SECRET`.

3. **Missing fields in onelogin_saml_apps**: Certain fields available in the OneLogin API are not accessible in the current Terraform provider, including `policy_id`, `tab_id`, `configuration.external_role`, `configuration.external_id`, `configuration.certificate_id`, and `provisioning.enabled`. The new provider will include all fields exposed by the OneLogin API, particularly those currently missing from the `onelogin_saml_apps` resource.

## Engineering Requirements

1. **Full API Coverage**: The new provider must support all endpoints and features exposed by the OneLogin API.

2. **Secure Credential Handling**: OneLogin API credentials should be securely managed by the provider, preferably by fetching the API token using environment variables.

3. **Subdomain Support**: The provider must correctly handle OneLogin subdomains to ensure error-free usage of data sources.

4. **Field Inclusion**: All fields exposed by the OneLogin API, including those currently absent from the `onelogin_saml_apps` resource, should be included in the new provider.

5. **Error Handling**: The provider must have robust error handling capabilities, providing users with meaningful error messages to aid in issue resolution.

6. **Documentation**: The provider's documentation should be comprehensive, including usage examples for every resource and data source.

## Milestones

While the precise timeline is difficult to determine, the following provides a general guide:

1. **API Exploration and Design**: Thorough exploration of the OneLogin API and design of the new provider's architecture.

2. **Initial Development**: Begin coding the new provider, prioritizing high-priority features and addressing known issues in the current provider.

3. **Testing and Iteration**: Extensive testing of the new provider, iterating its design and implementation based on test results.

4. **Documentation and Cleanup**: Creation of comprehensive documentation and resolution of any remaining codebase issues.

5. **Release and Maintenance**: Launch the new provider, maintain it, and incorporate improvements based on user feedback.

## Success Criteria

The success of this project will be evaluated based on the following criteria:

1. **Technical Benchmarks**: The provider supports all OneLogin API endpoints and features, handles credentials securely, correctly processes subdomains, includes all fields exposed by the API, and provides robust error handling.

2. **User Feedback**: Positive feedback from users regarding the functionality and usability of the new provider.

3. **Adoption Rate**: The number of projects adopting the new Terraform provider within six months of its release.

4. **Issue Resolution**: Known issues with the existing provider are effectively addressed in the new implementation.

5. **Documentation Quality**: The provider's documentation is comprehensive, includes usage examples for every resource and data source, and facilitates easy understanding and usage of the provider.

Please note that the milestones provided may be subject to change based on availability and project constraints.
