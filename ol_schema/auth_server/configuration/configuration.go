package authserverconfigurationschema

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
)

// Schema returns a key/value map of the various fields that make up the Rules of a OneLogin App.
func Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"resource_identifier": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"audiences": &schema.Schema{
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"access_token_expiration_minutes": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"refresh_token_expiration_minutes": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
	}
}

// Inflate takes a key/value map of interfaces and uses the fields to construct
// an AuthServerConfiguration struct for a OneLogin AuthServer.
func Inflate(in []interface{}) models.AuthServerConfiguration {
	s := in[0].(map[string]interface{})
	out := models.AuthServerConfiguration{}

	// Handle audiences
	if s["audiences"] != nil {
		switch audiences := s["audiences"].(type) {
		case []interface{}:
			// Handle the case where it's a []interface{}
			out.Audiences = make([]string, len(audiences))
			for i, val := range audiences {
				out.Audiences[i] = val.(string)
			}
		case []string:
			// Handle the case where it's already a []string
			out.Audiences = audiences
		}
	}

	// Handle resource identifier
	if ri, notNil := s["resource_identifier"].(string); notNil {
		out.ResourceIdentifier = &ri
	}

	// Handle token expirations
	if at, notNil := s["access_token_expiration_minutes"].(int); notNil {
		at32 := int32(at)
		out.AccessTokenExpirationMinutes = &at32
	}

	if rt, notNil := s["refresh_token_expiration_minutes"].(int); notNil {
		rt32 := int32(rt)
		out.RefreshTokenExpirationMinutes = &rt32
	}

	return out
}

// Flatten takes an AuthServer configuration and converts it to a map of varied types
func Flatten(asc models.AuthServerConfiguration) map[string]interface{} {
	out := map[string]interface{}{}

	if asc.ResourceIdentifier != nil {
		out["resource_identifier"] = *asc.ResourceIdentifier
	}

	if asc.Audiences != nil {
		out["audiences"] = asc.Audiences
	}

	if asc.AccessTokenExpirationMinutes != nil {
		out["access_token_expiration_minutes"] = *asc.AccessTokenExpirationMinutes
	}

	if asc.RefreshTokenExpirationMinutes != nil {
		out["refresh_token_expiration_minutes"] = *asc.RefreshTokenExpirationMinutes
	}

	return out
}
