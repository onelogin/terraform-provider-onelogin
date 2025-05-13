package privilegeschema

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
)

// Schema returns a key/value map of the various fields that make up a Privilege at OneLogin.
func Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"user_ids": &schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeInt},
		},
		"role_ids": &schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeInt},
		},
		"privilege": &schema.Schema{
			Type:     schema.TypeSet, // lets us define a sub-model and dictate the key name is privilege
			Required: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"version": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Default:  "2018-05-18",
					},
					"statement": &schema.Schema{
						Type:     schema.TypeList,
						Required: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"effect": &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},
								"action": &schema.Schema{
									Type:     schema.TypeList,
									Required: true,
									Elem:     &schema.Schema{Type: schema.TypeString},
								},
								"scope": &schema.Schema{
									Type:     schema.TypeList,
									Required: true,
									Elem:     &schema.Schema{Type: schema.TypeString},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Inflate takes a map of interfaces and constructs a Privilege object for the OneLogin API.
func Inflate(d map[string]interface{}) (models.Privilege, error) {
	pd, ok := d["privilege"].(*schema.Set).List()[0].(map[string]interface{})
	if !ok {
		return models.Privilege{}, errors.New("unable to parse terraform data for privilege")
	}

	// Process role IDs and user IDs
	var roleIDs []int
	var userIDs []int

	if d["role_ids"] != nil {
		rIDs := d["role_ids"].(*schema.Set).List()
		roleIDs = make([]int, len(rIDs))
		for i, r := range rIDs {
			roleIDs[i] = r.(int)
		}
	}

	if d["user_ids"] != nil {
		uIDs := d["user_ids"].(*schema.Set).List()
		userIDs = make([]int, len(uIDs))
		for i, u := range uIDs {
			userIDs[i] = u.(int)
		}
	}

	// Create the basic privilege object
	privilege := models.Privilege{
		RoleIDs:   roleIDs,
		UserIDs:   userIDs,
		Privilege: &models.PrivilegeData{},
	}

	// Handle basic fields
	if name, ok := d["name"].(string); ok {
		privilege.Name = &name
	}

	if desc, ok := d["description"].(string); ok {
		privilege.Description = &desc
	}

	if id, ok := d["id"].(string); ok {
		privilege.ID = &id
	}

	// Handle version
	if version, ok := pd["version"].(string); ok {
		privilege.Privilege.Version = &version
	}

	// Process statements
	if pd["statement"] != nil {
		ps := pd["statement"].([]interface{})
		privilege.Privilege.Statement = make([]models.StatementData, len(ps))

		for i, s := range ps {
			st := s.(map[string]interface{})

			stAct := st["action"].([]interface{})
			stSco := st["scope"].([]interface{})

			statementActions := make([]string, len(stAct))
			statementScopes := make([]string, len(stSco))

			for j, ac := range stAct {
				statementActions[j] = ac.(string)
			}
			for j, sc := range stSco {
				statementScopes[j] = sc.(string)
			}

			effect := st["effect"].(string)

			privilege.Privilege.Statement[i] = models.StatementData{
				Effect: &effect,
				Action: statementActions,
				Scope:  statementScopes,
			}
		}
	}

	return privilege, nil
}

// FlattenPrivilegeData converts a PrivilegeData struct to a format suitable for terraform state
func FlattenPrivilegeData(p models.PrivilegeData) []map[string]interface{} {
	statements := make([]map[string]interface{}, len(p.Statement))
	for i, s := range p.Statement {
		statements[i] = map[string]interface{}{
			"effect": *s.Effect,
			"action": s.Action,
			"scope":  s.Scope,
		}
	}
	return []map[string]interface{}{
		map[string]interface{}{
			"version":   *p.Version,
			"statement": statements,
		},
	}
}
