package smarthooksschema

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/smarthooks"
	"github.com/onelogin/terraform-provider-onelogin/utils"
	"log"
)

// Schema returns a key/value map of the various fields that make up the Rules of a OneLogin App.
func Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": &schema.Schema{
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validTypes,
		},
		"status": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"function": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"disabled": &schema.Schema{
			Type:     schema.TypeBool,
			Required: true,
		},
		"risk_enabled": &schema.Schema{
			Type:     schema.TypeBool,
			Required: true,
		},
		"location_enabled": &schema.Schema{
			Type:     schema.TypeBool,
			Required: true,
		},
		"retries": &schema.Schema{
			Type:     schema.TypeInt,
			Required: true,
		},
		"timeout": &schema.Schema{
			Type:     schema.TypeInt,
			Required: true,
		},
		"env_vars": &schema.Schema{
			Type:     schema.TypeList,
			Required: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"created_at": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"updated_at": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func validTypes(val interface{}, key string) (warns []string, errs []error) {
	return utils.OneOf(key, val.(string), []string{"pre-authentication"})
}

// Inflate takes a key/value map of interfaces and uses the fields to construct
// a AppProvisioning struct, a sub-field of a OneLogin App.
func Inflate(s map[string]interface{}) smarthooks.SmartHook {
	out := smarthooks.SmartHook{}
	if id, notNil := s["id"].(string); notNil {
		out.ID = oltypes.String(id)
	}

	if hookType, notNil := s["type"].(string); notNil {
		out.Type = oltypes.String(hookType)
	}

	if function, notNil := s["function"].(string); notNil {
		out.Function = oltypes.String(function)
	}

	if disabled, notNil := s["disabled"].(bool); notNil {
		out.Disabled = oltypes.Bool(disabled)
	}

	if riskEnabled, notNil := s["risk_enabled"].(bool); notNil {
		out.RiskEnabled = oltypes.Bool(riskEnabled)
	}

	if locationEnabled, notNil := s["location_enabled"].(bool); notNil {
		out.LocationEnabled = oltypes.Bool(locationEnabled)
	}

	if retries, notNil := s["retries"].(int); notNil {
		out.Retries = oltypes.Int32(int32(retries))
	}

	if timeout, notNil := s["timeout"].(int); notNil {
		out.Timeout = oltypes.Int32(int32(timeout))
	}

	if s["env_vars"] != nil {
		out.EnvVars = make([]string, len(s["env_vars"].([]interface{})))
		for i, envVar := range s["env_vars"].([]interface{}) {
			out.EnvVars[i] = envVar.(string)
		}
		log.Println("ENVVARS", out.EnvVars, s["env_vars"])
	}

	return out
}
