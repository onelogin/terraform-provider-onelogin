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
		// place rather than clearing it. Write policy_id = 0 to unassign: an
		// explicit zero is still a value in the configuration, so it diffs
		// against state, and Inflate turns it into the null the API wants.
		"policy_id": &schema.Schema{
			Type:         schema.TypeInt,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validAppAssignmentID,
		},
		// The brand whose login page this app uses. Optional and Computed for
		// the same reason as policy_id, and previously neither: brand_id was
		// Optional alone while being populated by every read, so an app branded
		// in the OneLogin UI proposed brand_id -> 0 on every plan. On the basic
		// app resource, the only one that sent the field, that apply then failed
		// with 422 "The associated AccountBrand with ID 0 could not be found";
		// on the OIDC and SAML resources the field was never sent at all, so the
		// diff simply came back for ever.
		//
		// Write brand_id = 0 to unassign, as with policy_id.
		"brand_id": &schema.Schema{
			Type:         schema.TypeInt,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validAppAssignmentID,
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

// validAppAssignmentID rejects an ID that cannot name a policy or a brand.
//
// 0 is allowed and means "none": Inflate turns it into the JSON null the app
// endpoint accepts as an unassignment. A negative is simply wrong, and saying
// so during planning beats an apply that comes back with 422 "The associated
// resource with the given id could not be found".
func validAppAssignmentID(val interface{}, key string) ([]string, []error) {
	id, ok := val.(int)
	if !ok || id >= 0 {
		return nil, nil
	}

	return nil, []error{fmt.Errorf("%s must be a positive ID, or 0 to unassign, got %d", key, id)}
}

// assignmentID reads policy_id or brand_id out of an inflate map.
//
// present is false when the key is absent, which is how a caller says "leave
// this assignment alone". A value that cannot be read is an error rather than
// something to step over: the key being present is the whole instruction, and
// the worst case is a dropped 0, where the plan promises an unassignment and
// the apply quietly does nothing, leaving a diff that never settles.
//
// The type cannot actually be wrong today -- both fields are TypeInt, so d.Get
// hands back an int -- but a map built by hand can get it wrong, and the
// neighbouring fields answer that with an unchecked assertion and a panic.
func assignmentID(s map[string]interface{}, key string) (id int, present bool, err error) {
	raw, ok := s[key]
	if !ok || raw == nil {
		return 0, false, nil
	}

	id, ok = raw.(int)
	if !ok {
		return 0, false, fmt.Errorf("%s must be an int, got %T", key, raw)
	}

	return id, true, nil
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

	// Both assignments follow the same rule; see the policy_id note below.
	brandID, brandGiven, err := assignmentID(s, "brand_id")
	if err != nil {
		return app, err
	}
	if brandGiven {
		if brandID == 0 {
			app.ClearBrandID = true
		} else {
			app.BrandID = &brandID
		}
	}

	// Presence of the key, not truth of the value, decides whether the policy
	// is sent, so that the caller controls it and this only translates.
	//
	// 0 is not sent as 0. The app endpoint refuses it -- 422 "The associated
	// Policy with ID 0 could not be found" -- and takes a JSON null as the
	// unassignment instead, which ClearPolicyID is how models.App expresses
	// (onelogin-go-sdk v4.16.0). A group, whose API does accept 0, is the
	// reason this is worth spelling out.
	policyID, policyGiven, err := assignmentID(s, "policy_id")
	if err != nil {
		return app, err
	}
	if policyGiven {
		if policyID == 0 {
			app.ClearPolicyID = true
		} else {
			app.PolicyID = &policyID
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
