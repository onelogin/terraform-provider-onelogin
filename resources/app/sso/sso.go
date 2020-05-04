package sso

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/models"
)

// OIDCSSOSchema returns a key/value map of the various fields that make up
// the SSO field for a OneLogin App.
func OIDCSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"client_id": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"client_secret": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

// SAMLSSOSchema returns a key/value map of the various fields that make up
// the SSO field for a OneLogin App.
func SAMLSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"acs_url": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"metadata_url": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"issuer": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"sls_url": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"certificate": &schema.Schema{
			Type:     schema.TypeList,
			MaxItems: 1,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"id": &schema.Schema{
						Type:     schema.TypeInt,
						Computed: true,
					},
					"value": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}

func FlattenOIDC(sso models.AppSso) []map[string]interface{} {
	return []map[string]interface{}{
		map[string]interface{}{
			"client_id":     sso.ClientID,
			"client_secret": sso.ClientSecret,
		},
	}
}

func FlattenSAML(sso models.AppSso) []map[string]interface{} {
	return []map[string]interface{}{
		map[string]interface{}{
			"metadata_url": sso.MetadataURL,
			"acs_url":      sso.AcsURL,
			"sls_url":      sso.SlsURL,
			"issuer":       sso.Issuer,
			"certificate": []map[string]interface{}{
				map[string]interface{}{
					"name":  sso.Certificate.Name,
					"id":    sso.Certificate.ID,
					"value": sso.Certificate.Value,
				},
			},
		},
	}
}
