package groupschema

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
)

// Schema returns a key/value map of the various fields that make up a OneLogin Group.
func Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": &schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
		},
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"reference": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		// The user policy applied to the group's members. Optional and not
		// Computed: Computed would read an empty configuration as "keep
		// whatever is in state", which leaves no way to say "no policy" once
		// one has been set. Optional lets the attribute be removed, and the
		// API accepts a 0 to clear the assignment.
		//
		// Only user policies are assignable. An app policy is refused with
		// 422 "Policy must reference a user policy".
		"policy_id": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
	}
}

// Inflate takes a map of interfaces and uses the fields to construct a Group
func Inflate(s map[string]interface{}) (models.Group, error) {
	var group models.Group

	if id, ok := s["id"]; ok && id != nil {
		if idInt, ok := id.(int); ok {
			group.ID = idInt
		}
	}

	if name, ok := s["name"]; ok && name != nil {
		if nameStr, ok := name.(string); ok {
			group.Name = nameStr
		}
	}

	if ref, ok := s["reference"]; ok && ref != nil {
		if refStr, ok := ref.(string); ok {
			refStrPtr := refStr
			group.Reference = &refStrPtr
		}
	}

	// Presence of the key, not truth of the value, decides whether the policy
	// is sent. A 0 is meaningful here -- it is how the assignment is cleared --
	// so the caller controls the key and this only translates it.
	if policyID, ok := s["policy_id"]; ok && policyID != nil {
		if policyIDInt, ok := policyID.(int); ok {
			id := policyIDInt
			group.PolicyID = &id
		}
	}

	return group, nil
}

// FlattenMany takes a slice of Group instances and converts them to a slice of maps
func FlattenMany(groups []models.Group) []map[string]interface{} {
	out := make([]map[string]interface{}, len(groups))
	for i, group := range groups {
		out[i] = map[string]interface{}{
			"id":   group.ID,
			"name": group.Name,
		}
		if group.Reference != nil {
			out[i]["reference"] = *group.Reference
		}
		if group.PolicyID != nil {
			out[i]["policy_id"] = *group.PolicyID
		}
	}
	return out
}

// Flatten takes a Group instance and converts it to a map
func Flatten(group models.Group) map[string]interface{} {
	out := map[string]interface{}{
		"id":   group.ID,
		"name": group.Name,
	}
	if group.Reference != nil {
		out["reference"] = *group.Reference
	}
	if group.PolicyID != nil {
		out["policy_id"] = *group.PolicyID
	}
	return out
}
