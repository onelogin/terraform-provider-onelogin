package authserverconfigurationschema

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	authservers "github.com/onelogin/onelogin-go-sdk/pkg/services/auth_servers"
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
// a AppProvisioning struct, a sub-field of a OneLogin App.
func Inflate(in []interface{}) authservers.AuthServerConfiguration {
	s := in[0].(map[string]interface{})
	out := authservers.AuthServerConfiguration{}
	if val, notNil := s["audiences"].([]string); notNil {
		out.Audiences = make([]string, len(val))
		for i, str := range val {
			out.Audiences[i] = str
		}
	}
	if ri, notNil := s["resource_identifier"].(string); notNil {
		out.ResourceIdentifier = oltypes.String(ri)
	}
	if at, notNil := s["access_token_expiration_minutes"].(int); notNil {
		out.AccessTokenExpirationMinutes = oltypes.Int32(int32(at))
	}
	if rt, notNil := s["refresh_token_expiration_minutes"].(int); notNil {
		out.RefreshTokenExpirationMinutes = oltypes.Int32(int32(rt))
	}
	return out
}

// Flatten takes an AuthServer configuration and converts it to a map of varied types
func Flatten(asc authservers.AuthServerConfiguration) []map[string]interface{} {
	out := make([]map[string]interface{}, 1)
	out[0] = map[string]interface{}{}
	if asc.ResourceIdentifier != nil {
		out[0]["resource_identifier"] = *asc.ResourceIdentifier
	}
	if asc.Audiences != nil {
		out[0]["audiences"] = asc.Audiences
	}
	if asc.AccessTokenExpirationMinutes != nil {
		out[0]["access_token_expiration_minutes"] = *asc.AccessTokenExpirationMinutes
	}
	if asc.RefreshTokenExpirationMinutes != nil {
		out[0]["refresh_token_expiration_minutes"] = *asc.RefreshTokenExpirationMinutes
	}
	return out
}
