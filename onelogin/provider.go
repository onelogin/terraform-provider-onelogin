package onelogin

import (
	"errors"

	"github.com/onelogin/onelogin-go-sdk/pkg/client"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var (
	errClientCredentials = errors.New("client_id or client_sercret or region missing")
)

// Provider creates a new provider with all the neccessary configurations.
// It returns a pointer to the created provider.
func Provider() terraform.ResourceProvider {
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
				Default:  client.USRegion,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"onelogin_app_role_attachments": AppRoleAttachment(),
			"onelogin_apps":                 Apps(),
			"onelogin_oidc_apps":            OIDCApps(),
			"onelogin_saml_apps":            SAMLApps(),
			"onelogin_app_rules":            AppRules(),
			"onelogin_user_mappings":        UserMappings(),
			"onelogin_users":                Users(),
			"onelogin_auth_servers":         AuthServers(),
			"onelogin_roles":                Roles(),
			"onelogin_smarthooks":           SmartHooks(),
		},
		ConfigureFunc: configProvider,
	}
}

// configProvider configures the provider, and if successful, it returns
// an interface containing the api client.
func configProvider(d *schema.ResourceData) (interface{}, error) {
	clientID := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)
	region := d.Get("region").(string)
	url := d.Get("url").(string)

	timeout := client.DefaultTimeout

	oneloginClient, err := client.NewClient(&client.APIClientConfig{
		Timeout:      timeout,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Region:       region,
		Url:          url,
	})
	if err != nil {
		return nil, err
	}
	return oneloginClient, nil
}
