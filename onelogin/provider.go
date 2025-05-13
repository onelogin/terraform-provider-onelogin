package onelogin

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

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
				DefaultFunc: schema.EnvDefaultFunc("ONELOGIN_API_URL", nil),
				Optional:    true,
				Description: "OneLogin API URL. This is an alternative to using subdomain. If both are provided, subdomain takes precedence.",
			},
			// Both region and subdomain are deprecated and will be removed in a future version
			"region": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Use url instead",
			},
			"subdomain": {
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("ONELOGIN_SUBDOMAIN", nil),
				Optional:    true,
				Deprecated:  "Use url instead",
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
			"onelogin_app_role_attachments": AppRoleAttachment(),
			"onelogin_apps":                 Apps(),
			"onelogin_oidc_apps":            OIDCApps(),
			"onelogin_saml_apps":            SAMLApps(),
			"onelogin_app_rules":            AppRules(),
			// "onelogin_user_mappings":                   UserMappings(), // Disabled until SDK support is added
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

	// Set client credentials environment variables for the SDK
	os.Setenv("ONELOGIN_CLIENT_ID", clientID)
	os.Setenv("ONELOGIN_CLIENT_SECRET", clientSecret)

	// Prioritize URL over subdomain
	url := d.Get("url").(string)
	subdomain := d.Get("subdomain").(string)

	if url != "" {
		// Set API URL for SDK - now the preferred way
		os.Setenv("ONELOGIN_API_URL", url)

		// Extract subdomain from URL for backward compatibility with SDK
		// URL format is typically https://company.onelogin.com
		urlParts := strings.Split(strings.TrimPrefix(strings.TrimPrefix(url, "https://"), "http://"), ".")
		if len(urlParts) > 0 {
			extractedSubdomain := urlParts[0]
			os.Setenv("ONELOGIN_SUBDOMAIN", extractedSubdomain)
		} else {
			return nil, diag.Errorf("Could not extract subdomain from URL. Please provide a valid OneLogin URL.")
		}
	} else if subdomain != "" {
		// For backward compatibility
		os.Setenv("ONELOGIN_SUBDOMAIN", subdomain)
		// Also set the API URL for consistency
		os.Setenv("ONELOGIN_API_URL", fmt.Sprintf("https://%s.onelogin.com", subdomain))
	} else {
		return nil, diag.Errorf("Either OneLogin API URL or subdomain is required. Please set the ONELOGIN_API_URL environment variable.")
	}

	// Initialize the SDK
	client, err := ol.NewOneloginSDK()
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return client, nil
}
