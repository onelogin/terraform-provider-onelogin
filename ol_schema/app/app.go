package appschema

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	appconfigurationschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/app/configuration"
	appparametersschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/app/parameters"
	appprovisioningschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/app/provisioning"
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
	}
}

// Inflate takes a map of interfaces and constructs a OneLogin App.
func Inflate(s map[string]interface{}) (models.App, error) {
	var appID, connectorID int32
	var name, description, notes string
	var visible, allowAssumedSignin bool

	// Set required/common fields
	name = s["name"].(string)

	if s["description"] != nil {
		description = s["description"].(string)
	}

	if s["notes"] != nil {
		notes = s["notes"].(string)
	}

	if s["connector_id"] != nil {
		connectorID = int32(s["connector_id"].(int))
	}

	if s["visible"] != nil {
		visible = s["visible"].(bool)
	}

	if s["allow_assumed_signin"] != nil {
		allowAssumedSignin = s["allow_assumed_signin"].(bool)
	}

	app := models.App{
		Name:               &name,
		Description:        &description,
		Notes:              &notes,
		ConnectorID:        &connectorID,
		Visible:            &visible,
		AllowAssumedSignin: &allowAssumedSignin,
	}

	// Set optional fields
	if s["id"] != nil {
		if id, err := strconv.Atoi(s["id"].(string)); err == nil {
			appID = int32(id)
			app.ID = &appID
		}
	}

	if s["brand_id"] != nil {
		brandID := s["brand_id"].(int)
		app.BrandID = &brandID
	}

	// Handle parameters
	if s["parameters"] != nil {
		p := s["parameters"].(*schema.Set).List()
		params := make(map[string]models.Parameter, len(p))
		for _, val := range p {
			valMap := val.(map[string]interface{})
			params[valMap["param_key_name"].(string)] = appparametersschema.Inflate(valMap)
		}
		app.Parameters = &params
	}

	// Handle provisioning
	if s["provisioning"] != nil {
		prov := appprovisioningschema.Inflate(s["provisioning"].(map[string]interface{}))
		app.Provisioning = &prov
	}

	// Handle configuration
	if s["configuration"] != nil {
		conf, err := appconfigurationschema.Inflate(s["configuration"].(map[string]interface{}))
		if err != nil {
			return app, err
		}
		app.Configuration = conf
	}

	return app, nil
}
