package smarthooksschema

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/smarthooks"
	smarthookenvs "github.com/onelogin/onelogin-go-sdk/pkg/services/smarthooks/envs"
	smarthookconditionsschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/smarthook/conditions"
	smarthookoptions "github.com/onelogin/terraform-provider-onelogin/ol_schema/smarthook/options"
	"github.com/onelogin/terraform-provider-onelogin/utils"
)

// Schema returns a key/value map of the various fields that make up the Rules of a OneLogin App.
func Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validTypes,
		},
		"disabled": {
			Type:     schema.TypeBool,
			Required: true,
		},
		"timeout": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"env_vars": {
			Type:     schema.TypeList,
			Required: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"runtime": {
			Type:     schema.TypeString,
			Required: true,
		},
		"context_version": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"retries": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"options": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: smarthookoptions.Schema(),
			},
		},
		"packages": {
			Type:     schema.TypeMap,
			Required: true,

			Elem: &schema.Schema{Type: schema.TypeString},
		},
		"function": {
			Type:     schema.TypeString,
			Required: true,
		},
		"conditions": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: smarthookconditionsschema.Schema(),
			},
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"updated_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func validTypes(val interface{}, key string) (warns []string, errs []error) {
	return utils.OneOf(key, val.(string), []string{"pre-authentication", "user-migration"})
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

	if runtime, notNil := s["runtime"].(string); notNil {
		out.Runtime = oltypes.String(runtime)
	}

	if function, notNil := s["function"].(string); notNil {
		out.Function = oltypes.String(function)
	}

	if disabled, notNil := s["disabled"].(bool); notNil {
		out.Disabled = oltypes.Bool(disabled)
	}

	if retries, notNil := s["retries"].(int); notNil {
		out.Retries = oltypes.Int32(int32(retries))
	}

	if timeout, notNil := s["timeout"].(int); notNil {
		out.Timeout = oltypes.Int32(int32(timeout))
	}

	if s["env_vars"] != nil {
		out.EnvVars = make([]smarthookenvs.EnvVar, len(s["env_vars"].([]interface{})))
		for i, envVar := range s["env_vars"].([]interface{}) {
			out.EnvVars[i] = smarthookenvs.EnvVar{Name: oltypes.String(envVar.(string))}
		}
	}

	if s["conditions"] != nil {
		out.Conditions = []smarthooks.Condition{}
		for _, val := range s["conditions"].([]interface{}) {
			cond := smarthookconditionsschema.Inflate(val.(map[string]interface{}))
			out.Conditions = append(out.Conditions, cond)
		}
	}

	if s["options"] != nil {
		opts := smarthookoptions.Inflate(s["options"].(map[string]interface{}))
		out.Options = &opts
	}

	if s["packages"] != nil {
		out.Packages = make(map[string]string, len(s["packages"].(map[string]interface{})))
		for pkg, ver := range s["packages"].(map[string]interface{}) {
			out.Packages[pkg] = ver.(string)
		}
	}
	return out
}

// FlattenEnvVars takes a SmartHook and gets a list of env_var names
func FlattenEnvVars(vars []smarthookenvs.EnvVar) []string {
	out := make([]string, len(vars))
	for i, v := range vars {
		out[i] = *v.Name
	}
	return out
}
