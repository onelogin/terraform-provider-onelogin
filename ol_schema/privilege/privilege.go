package privilegeschema

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/privileges"
)

// Schema returns a key/value map of the various fields that make up an App at OneLogin.
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

func Inflate(d map[string]interface{}) (privileges.Privilege, error) {
	pd, ok := d["privilege"].(*schema.Set).List()[0].(map[string]interface{})
	if !ok {
		return privileges.Privilege{}, errors.New("unable to parse terraform data for privilege")
	}

	rIDs := d["role_ids"].(*schema.Set).List()
	uIDs := d["user_ids"].(*schema.Set).List()

	roleIDs := make([]int, len(rIDs))
	userIDs := make([]int, len(uIDs))

	for i, r := range rIDs {
		roleIDs[i] = r.(int)
	}

	for i, u := range uIDs {
		userIDs[i] = u.(int)
	}
	privilege := privileges.Privilege{
		Name:        oltypes.String(d["name"].(string)),
		Description: oltypes.String(d["description"].(string)),
		RoleIDs:     roleIDs,
		UserIDs:     userIDs,
		Privilege: &privileges.PrivilegeData{
			Version: oltypes.String(pd["version"].(string)),
		},
	}
	if d["id"] != nil {
		privilege.ID = oltypes.String(d["id"].(string))
	}
	ps := pd["statement"].([]interface{})
	privilege.Privilege.Statement = make([]privileges.StatementData, len(ps))
	for i, s := range ps {
		st := s.(map[string]interface{})

		stAct := st["action"].([]interface{})
		stSco := st["scope"].([]interface{})

		statementActions := make([]string, len(stAct))
		statementScopes := make([]string, len(stSco))

		for i, ac := range stAct {
			statementActions[i] = ac.(string)
		}
		for i, sc := range stSco {
			statementScopes[i] = sc.(string)
		}
		privilege.Privilege.Statement[i] = privileges.StatementData{
			Effect: oltypes.String(st["effect"].(string)),
			Action: statementActions,
			Scope:  statementScopes,
		}
	}
	return privilege, nil
}
func FlattenPrivilegeData(p privileges.PrivilegeData) []map[string]interface{} {
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
