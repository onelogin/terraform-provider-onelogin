package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/onelogin/onelogin-go-sdk/pkg/client"
)

// Ensure oneloginProvider satisfies various provider interfaces.
var _ provider.Provider = &oneloginProvider{}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &oneloginProvider{
			version: version,
		}
	}
}

// oneloginProvider defines the provider implementation.
type oneloginProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// oneloginProviderModel describes the provider data model.
type oneloginProviderModel struct {
	ClientId     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	SubDomain    types.String `tfsdk:"subdomain"`
}

func (p *oneloginProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "onelogin"
	resp.Version = p.version
}

func (p *oneloginProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with Onelogin",
		Attributes: map[string]schema.Attribute{
			"client_id": schema.StringAttribute{
				Description: "onelogin Client Id",
				Optional:    true,
			},
			"client_secret": schema.StringAttribute{
				Description: "onelogin Client Secret",
				Optional:    true,
				Sensitive:   true,
			},
			"subdomain": schema.StringAttribute{
				MarkdownDescription: "onelogin subdomain",
				Optional:            true,
			},
		},
	}
}

func (p *oneloginProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Onelogin client")

	var config oneloginProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	if config.ClientId.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_id"),
			"Unknown Onelogin Client ID",
			"The provider cannot create the Onelogin API client as there is an unknown configuration value for the Onelogin API client_id. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ONELOGIN_CLIENT_ID environment variable.",
		)
	}

	if config.ClientSecret.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_secret"),
			"Unknown Onelogin Client Secret",
			"The provider cannot create the Onelogin API client as there is an unknown configuration value for the Onelogin API client_secret. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ONELOGIN_CLIENT_SECRET environment variable.",
		)
	}

	if config.SubDomain.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("subomain"),
			"Unknown Onelogin subdomain",
			"The provider cannot create the Onelogin API client as there is an unknown configuration value for the Onelogin subdomain. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ONELOGIN_SUBDOMAIN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}
	// Example client configuration for data sources and resources
	client_id := os.Getenv("ONELOGIN_CLIENT_ID")
	client_secret := os.Getenv("ONELOGIN_CLIENT_SECRET")
	subdomain := os.Getenv("ONELOGIN_SUBDOMAIN")

	if !config.ClientId.IsNull() {
		client_id = config.ClientId.ValueString()
	}

	if !config.ClientSecret.IsNull() {
		client_secret = config.ClientSecret.ValueString()
	}

	if !config.SubDomain.IsNull() {
		subdomain = config.SubDomain.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if client_id == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_id"),
			"Missing onelogin API client_id",
			"The provider cannot create the onelogin API client as there is a missing or empty value for the onelogin API client_id. "+
				"Set the client_id value in the configuration or use the ONELOGIN_CLIENT_ID environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if client_secret == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_secret"),
			"Missing onelogin API client_secret",
			"The provider cannot create the onelogin API client as there is a missing or empty value for the onelogin API client_secret. "+
				"Set the client_secret value in the configuration or use the ONELOGIN_CLIENT_SECRET environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if subdomain == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("subdomain"),
			"Missing onelogin API subdomain",
			"The provider cannot create the onelogin API client as there is a missing or empty value for the onelogin subdomain. "+
				"Set the url value in the configuration or use the ONELOGIN_SUBDOMAIN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "ONELOGIN_CLIENT_ID", client_id)
	ctx = tflog.SetField(ctx, "ONELOGIN_CLIENT_SECRET", client_secret)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "ONELOGIN_CLIENT_SECRET")
	ctx = tflog.SetField(ctx, "ONELOGIN_SUBDOMAIN", subdomain)

	tflog.Debug(ctx, "Creating onelogin client")

	// Create a new onelogin client using the configuration values
	client, err := client.NewClient(&client.APIClientConfig{
		ClientID:     client_id,
		ClientSecret: client_secret,
		SubDomain:    subdomain,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create onelogin API Client",
			"An unexpected error occurred when creating the onelogin API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"onelogin Client Error: "+err.Error(),
		)
		return
	}

	// Make the onelogin client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured onelogin client", map[string]any{"success": true})
}

func (p *oneloginProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewExampleResource,
	}
}

func (p *oneloginProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewExampleDataSource,
	}
}
