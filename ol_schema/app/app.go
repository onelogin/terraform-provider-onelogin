package appschema

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models" // Replace with the new SDK package path
)

// Schema returns a key/value map of the various fields that make up an App at OneLogin.
func Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"visible": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"notes": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		"icon_url": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"auth_method": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
		},
		"policy_id": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
		},
		"brand_id": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"allow_assumed_signin": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"tab_id": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
		},
		"connector_id": &schema.Schema{
			Type:     schema.TypeInt,
			Required: true,
		},
		"created_at": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"updated_at": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"provisioning": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeBool},
		},
		"parameters": &schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: appparametersschema.Schema(),
			},
		},
		"app_type": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
	}
}

// Inflate takes a map of interfaces and constructs a OneLogin App.
func Inflate(s map[string]interface{}) (models.App, error) {
	var err error
	app := models.App{
		Name:               oltypes.String(s["name"].(string)),
		Description:        oltypes.String(s["description"].(string)),
		Notes:              oltypes.String(s["notes"].(string)),
		ConnectorID:        oltypes.Int32(int32(s["connector_id"].(int))),
		Visible:            oltypes.Bool(s["visible"].(bool)),
		AllowAssumedSignin: oltypes.Bool(s["allow_assumed_signin"].(bool)),
	}
	if s["id"] != nil {
		if id, err := strconv.Atoi(s["id"].(string)); err == nil {
			app.ID = oltypes.Int32(int32(id))
		}
	}
	if s["brand_id"] != nil {
		brandID := s["brand_id"].(int)
		app.BrandID = oltypes.Int32(int32(brandID))
	}
	if s["parameters"] != nil {
		p := s["parameters"].(*schema.Set).List()
		app.Parameters = make(map[string]models.Parameter, len(p))
		for _, val := range p {
			valMap := val.(map[string]interface{})
			app.Parameters[valMap["param_key_name"].(string)] = appparametersschema.Inflate(valMap)
		}
	}
	if s["provisioning"] != nil {
		prov := appprovisioningschema.Inflate(s["provisioning"].(map[string]interface{}))
		app.Provisioning = &prov
	}
	if s["app_type"] != nil {
		app.AppType = s["app_type"].(string)
	}
	if s["sso"] != nil {
		switch ssoType := s["sso"].(type) {
		case map[string]interface{}:
			switch ssoType["type"].(string) {
			case "openid":
				app.SSO = &SSOOpenID{
					ClientID: ssoType["client_id"].(string),
				}
			case "saml":
				app.SSO = &SSOSAML{
					MetadataURL: ssoType["metadata_url"].(string),
					AcsURL:      ssoType["acs_url"].(string),
					SlsURL:      ssoType["sls_url"].(string),
					Issuer:      ssoType["issuer"].(string),
					Certificate: Certificate{
						ID:    ssoType["certificate"].(map[string]interface{})["id"].(int),
						Name:  ssoType["certificate"].(map[string]interface{})["name"].(string),
						Value: ssoType["certificate"].(map[string]interface{})["value"].(string),
					},
				}
			}
		}
	}

	if s["configuration"] != nil {
		switch configType := s["configuration"].(type) {
		case map[string]interface{}:
			switch configType["type"].(string) {
			case "openid":
				app.Configuration = &ConfigurationOpenID{
					RedirectURI:                   configType["redirect_uri"].(string),
					LoginURL:                      configType["login_url"].(string),
					OidcApplicationType:           configType["oidc_application_type"].(int),
					TokenEndpointAuthMethod:       configType["token_endpoint_auth_method"].(int),
					AccessTokenExpirationMinutes:  configType["access_token_expiration_minutes"].(int),
					RefreshTokenExpirationMinutes: configType["refresh_token_expiration_minutes"].(int),
				}
			case "saml":
				app.Configuration = &ConfigurationSAML{
					ProviderArn:        configType["provider_arn"],
					SignatureAlgorithm: configType["signature_algorithm"].(string),
					CertificateID:      configType["certificate_id"].(int),
				}
			}
		}
	}
	return app, err
}
