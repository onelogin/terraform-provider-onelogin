package onelogin

import (
	"context"
	"errors"
	"os"

	ol "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	errClientCredentials = errors.New("client_id or client_secret missing")
)

// Provider creates a new provider with all the neccessary configurations.
// It returns a pointer to the created provider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("ONELOGIN_CLIENT_ID", nil),
				Required:    true,
			},
			"client_secret": {
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("ONELOGIN_CLIENT_SECRET", nil),
				Required:    true,
			},
			"url": {
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("ONELOGIN_OAPI_URL", nil),
				Optional:    true,
			},
			// Region is deprecated and will be removed in a future version
			"region": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Use subdomain instead",
			},
			"subdomain": {
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("ONELOGIN_SUBDOMAIN", nil),
				Required:    true,
				Description: "OneLogin subdomain (e.g. 'company' for company.onelogin.com)",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ONELOGIN_CLIENT_TIMEOUT", 60),
				Description: "Timeout in seconds for API operations. Defaults to 60 seconds if not specified.",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"onelogin_user":  dataSourceUser(),
			"onelogin_users": dataSourceUsers(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"onelogin_app_role_attachments":            AppRoleAttachment(),
			"onelogin_apps":                            Apps(),
			"onelogin_oidc_apps":                       OIDCApps(),
			"onelogin_saml_apps":                       SAMLApps(),
			"onelogin_app_rules":                       AppRules(),
			"onelogin_user_mappings":                   UserMappings(),
			"onelogin_users":                           Users(),
			"onelogin_auth_servers":                    AuthServers(),
			"onelogin_roles":                           Roles(),
			"onelogin_smarthooks":                      SmartHooks(),
			"onelogin_smarthook_environment_variables": SmarthookEnvironmentVariables(),
			"onelogin_privileges":                      Privileges(),
			"onelogin_user_custom_attributes":          UserCustomAttributes(),
		},
		ConfigureContextFunc: configProvider,
	}
}

// configProvider configures the provider, and if successful, it returns
// an interface containing the api client.
func configProvider(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	clientID := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)

	// These are no longer used but kept for reference
	// url := d.Get("url").(string)
	// timeoutSeconds := d.Get("timeout").(int)

	// We're not using the v1 client anymore, but keeping this commented code for reference
	// in case we need to initialize it in the future
	// _, err := client.NewClient(&client.APIClientConfig{
	//	Timeout:      60,
	//	ClientID:     clientID,
	//	ClientSecret: clientSecret,
	//	Region:       "us",
	// })
	// if err != nil {
	//	return nil, diag.FromErr(err)
	// }

	// For resource types that use v4 client (like custom attributes)
	subdomain := d.Get("subdomain").(string)
	if subdomain == "" {
		return nil, diag.Errorf("OneLogin subdomain is required. Please set the ONELOGIN_SUBDOMAIN environment variable.")
	}

	// Set environment variables for the SDK
	os.Setenv("ONELOGIN_CLIENT_ID", clientID)
	os.Setenv("ONELOGIN_CLIENT_SECRET", clientSecret)
	os.Setenv("ONELOGIN_SUBDOMAIN", subdomain)

	// Initialize the SDK
	clientV4, err := ol.NewOneloginSDK()
	if err != nil {
		return nil, diag.FromErr(err)
	}

	// For now, return the v4 client as we're updating custom attributes
	return clientV4, nil
}