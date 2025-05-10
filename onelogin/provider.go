package onelogin

import (
	"context"
	"errors"

	"github.com/onelogin/onelogin-go-sdk/pkg/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	// USRegion is the US region identifier
	USRegion string = "us"
	// EURegion is the EU region identifier
	EURegion string = "eu"
)

var (
	errClientCredentials = errors.New("client_id or client_sercret or region missing")
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
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  USRegion,
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
	region := d.Get("region").(string)
	url := d.Get("url").(string)

	// Get timeout from configuration
	timeoutSeconds := d.Get("timeout").(int)

	// Use the original client for now for backward compatibility
	oneloginClient, err := client.NewClient(&client.APIClientConfig{
		Timeout:      timeoutSeconds,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Region:       region,
		Url:          url,
	})
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return oneloginClient, nil
}