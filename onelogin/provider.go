package onelogin

import (
	"context"
	"errors"
	"os"
	"strconv"
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
				Required:    true,
				Description: "OneLogin API URL (e.g. https://api.us.onelogin.com or https://api.eu.onelogin.com)",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ONELOGIN_TIMEOUT", 180),
				Description: "Timeout in seconds for API operations. Defaults to 180 seconds if not specified.",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"onelogin_user":   dataSourceUser(),
			"onelogin_users":  dataSourceUsers(),
			"onelogin_group":  dataSourceOneLoginGroup(),
			"onelogin_groups": dataSourceOneLoginGroups(),
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
			"onelogin_groups":                          resourceOneLoginGroups(),
			"onelogin_self_registration_profiles":      SelfRegistrationProfiles(),
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

	// Set a longer timeout for API operations
	timeout := d.Get("timeout").(int)
	os.Setenv("ONELOGIN_TIMEOUT", strconv.Itoa(timeout))
	// Keep setting the old env var for backward compatibility
	os.Setenv("ONELOGIN_CLIENT_TIMEOUT", strconv.Itoa(timeout))

	// Set the API URL
	url := d.Get("url").(string)
	if url == "" {
		return nil, diag.Errorf("OneLogin API URL is required. Please set the ONELOGIN_API_URL environment variable.")
	}

	// Set API URL for SDK
	os.Setenv("ONELOGIN_API_URL", url)

	// Extract subdomain from URL for backward compatibility with SDK's internals
	// Most SDK functions still use the subdomain internally
	urlParts := strings.Split(strings.TrimPrefix(strings.TrimPrefix(url, "https://"), "http://"), ".")
	if len(urlParts) > 0 && urlParts[0] != "api" {
		// Direct subdomain URL (e.g., company.onelogin.com)
		extractedSubdomain := urlParts[0]
		os.Setenv("ONELOGIN_SUBDOMAIN", extractedSubdomain)
	} else if len(urlParts) > 1 && urlParts[0] == "api" {
		// API URL format (e.g., api.us.onelogin.com or api.eu.onelogin.com)
		region := urlParts[1]
		if region == "us" || region == "eu" {
			// This is a valid API URL, but we need to set a dummy subdomain
			// as the SDK still requires ONELOGIN_SUBDOMAIN to be set
			os.Setenv("ONELOGIN_SUBDOMAIN", "dummy")
		} else {
			return nil, diag.Errorf("Invalid API URL format. Expected api.us.onelogin.com or api.eu.onelogin.com")
		}
	} else {
		return nil, diag.Errorf("Could not extract subdomain from URL. Please provide a valid OneLogin URL.")
	}

	// Initialize the SDK
	client, err := ol.NewOneloginSDK()
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return client, nil
}
