package onelogin

import (
	"context"
	"errors"
	"fmt"
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

// subdomainFromURL pulls the tenant subdomain out of the configured API URL.
//
// The SDK still asks for ONELOGIN_SUBDOMAIN even though every call goes to
// ONELOGIN_API_URL, so this exists to satisfy it. A regional API host has no
// tenant subdomain in it at all, hence the placeholder.
func subdomainFromURL(url string) (string, error) {
	parts := strings.Split(strings.TrimPrefix(strings.TrimPrefix(url, "https://"), "http://"), ".")

	switch {
	case len(parts) > 0 && parts[0] != "" && parts[0] != "api":
		// Direct subdomain URL (e.g., company.onelogin.com)
		return parts[0], nil
	case len(parts) > 1 && parts[0] == "api":
		// API URL format (e.g., api.us.onelogin.com or api.eu.onelogin.com)
		if region := parts[1]; region != "us" && region != "eu" {
			return "", fmt.Errorf("invalid API URL format, expected api.us.onelogin.com or api.eu.onelogin.com, got %q", url)
		}
		return "dummy", nil
	default:
		return "", fmt.Errorf("could not extract a subdomain from %q, please provide a valid OneLogin URL", url)
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
	subdomain, err := subdomainFromURL(url)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	os.Setenv("ONELOGIN_SUBDOMAIN", subdomain)

	// Initialize the SDK
	client, err := ol.NewOneloginSDK()
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return client, nil
}
