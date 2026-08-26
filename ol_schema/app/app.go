package appschema

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	appconfigurationschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/app/configuration"
	appparametersschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/app/parameters"
	appprovisioningschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/app/provisioning"
)

// SchemaV0 returns the app schema as it was before a parameter's values and
// default_values became lists. The state upgrader needs it to decode state
// written by an earlier provider version.
//
// Only parameters differ; the rest comes from Schema so the two cannot drift
// apart as fields are added.
func SchemaV0() map[string]*schema.Schema {
	v0 := Schema()
	v0["parameters"] = &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Computed: true,
		Set:      appparametersschema.HashByKeyName,
		Elem: &schema.Resource{
			Schema: appparametersschema.SchemaV0(),
		},
	}
	return v0
}

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
		// The app policy enforced when users sign in to this app. Only a
		// policy whose kind is "app" resolves here; the endpoint answers a
		// user policy with 422 "The associated Policy with ID <n> could not be
		// found", the same thing it says about an ID that does not exist.
		//
		// Optional *and* Computed, which is the opposite of the group
		// resource's policy_id. There the attribute was new, so Computed would
		// only have taken away the ability to say "no policy". Here it has
		// been Computed since the resource shipped, so state already holds
		// whatever policy the app was given in the OneLogin UI; dropping
		// Computed would read every configuration that never mentioned
		// policy_id as asking for 0 and unassign the policy on the first apply
		// after the provider was upgraded. #260.
		//
		// The cost is that removing the attribute leaves the last value in
		// place rather than clearing it, and validAppPolicyID explains why
		// writing 0 is not the way out.
		"policy_id": &schema.Schema{
			Type:         schema.TypeInt,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validAppPolicyID,
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
			Set:      appparametersschema.HashByKeyName,
			Elem: &schema.Resource{
				Schema: appparametersschema.Schema(),
			},
		},
	}
}

// validAppPolicyID rejects the two values that cannot reach the API as a policy
// assignment, so that a plan fails with the reason rather than an apply failing
// with a 422.
//
// A group clears its policy by being sent a 0, and an app cannot: OneLogin
// answers 0 with 422 "The associated Policy with ID 0 could not be found". What
// the app endpoint accepts as "no policy" is a JSON null, and models.App
// declares PolicyID as *int tagged omitempty, where a nil pointer is omitted
// from the request entirely and so leaves the assignment alone. There is no
// third state to send, which is why unassigning is not expressible here yet and
// needs a change to the SDK model rather than to this provider.
func validAppPolicyID(val interface{}, key string) ([]string, []error) {
	id, ok := val.(int)
	if !ok || id > 0 {
		return nil, nil
	}

	if id < 0 {
		return nil, []error{fmt.Errorf("%s must be the ID of an app policy, got %d", key, id)}
	}

	return nil, []error{fmt.Errorf(
		"%s cannot be set to 0: OneLogin rejects it with \"The associated Policy with ID 0 could not be found\". "+
			"Unassigning an app's policy is not supported yet -- it requires sending policy_id as null, which the "+
			"OneLogin Go SDK's App model cannot express. Assign a different app policy, or clear the assignment in "+
			"the OneLogin admin UI. Note that removing policy_id from the configuration does not clear it either, "+
			"because the attribute is computed as well as optional", key)}
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
		ConnectorID:        &connectorID,
		Visible:            &visible,
		AllowAssumedSignin: &allowAssumedSignin,
	}

	// Left nil when empty rather than pointed at "". Both fields are tagged
	// omitempty, but a pointer to an empty string is not empty -- only a nil
	// pointer is -- so taking the address unconditionally sent `"notes": ""`
	// for a field the API had returned as null, and the API stores what it is
	// sent. An update touching only the description rewrote them.
	//
	// Omitting a field leaves it alone: the app endpoint takes a PUT but
	// merges it, which is what makes leaving them out the right repair rather
	// than reading the app first and echoing every field back.
	if description != "" {
		app.Description = &description
	}
	if notes != "" {
		app.Notes = &notes
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

	// Presence of the key, not truth of the value, decides whether the policy
	// is sent, so that the caller controls it and this only translates. Unlike
	// a group, an app has no value meaning "no policy": 0 is refused and null
	// is unreachable through a *int tagged omitempty, so the callers only ever
	// hand over a real policy ID. See validAppPolicyID.
	if policyID, ok := s["policy_id"]; ok && policyID != nil {
		if policyIDInt, ok := policyID.(int); ok {
			id := policyIDInt
			app.PolicyID = &id
		}
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
