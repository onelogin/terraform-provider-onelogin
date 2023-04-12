package main

import (
	"log"

	"github.com/dikhan/terraform-provider-openapi/v3/openapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

var (
	// Version specifies the version of the provider (will be set statically at compile time)
	Version = "dev"
	// Commit specifies the commit hash of the provider at the time of building the binary (will be set statically at compile time)
	Commit = "none"
	// Date specifies the data which the binary was build (will be set statically at compile time)
	Date = "unknown"

// Generate the Terraform provider documentation using `tfplugindocs`:
)

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
func main() {

	//log.Printf("[INFO] Running Terraform Provider %s v%s-%s; Released on: %s", ProviderName, Version, Commit, Date)

	//log.Printf("[INFO] Initializing OpenAPI Terraform provider '%s' with service provider's OpenAPI document: %s", ProviderName, ProviderOpenAPIURL)

	var providerName = "onelogin"
	var providerOpenAPIURL = "https://raw.githubusercontent.com/onelogin/terraform-provider-onelogin/develop/swag-api.yml"

	p := openapi.ProviderOpenAPI{ProviderName: providerName}
	serviceProviderConfig := &openapi.ServiceConfigV1{
		SwaggerURL: providerOpenAPIURL,
	}

	provider, err := p.CreateSchemaProviderFromServiceConfiguration(serviceProviderConfig)
	if err != nil {
		log.Fatalf("[ERROR] Failed to initialize the terraform provider: %s", err)
	}

	plugin.Serve(
		&plugin.ServeOpts{
			ProviderFunc: func() *schema.Provider {
				return provider
			},
		},
	)
}
