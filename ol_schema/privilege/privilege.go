package privilegeschema

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
)

// Schema returns a key/value map of the various fields that make up a Privilege at OneLogin.
// DefaultVersion is the privilege document version used when a configuration
// does not name one, and when a response omits it. Defined here so the schema
// default and privilegeRead's fallback cannot drift apart.
const DefaultVersion = "2018-05-18"

// OrderStatementsLikeState puts the statements a read returned into the order
// the configuration already has them in.
//
// statement is a TypeList, so the order of a document's statements is part of
// the value -- but the API does not preserve it. The same privilege comes back
// with its statements in a different order from one read to the next, so a
// configuration whose statements happen to come back swapped plans a
// replacement of the whole block, and the read after that may swap them again.
// The acceptance test failed roughly every other run for exactly this reason,
// which read as a test isolation problem and was not one.
//
// Matching the order already in state is not pretending the documents agree:
// the statements are compared on their contents, and only ones that are
// genuinely the same statement are moved. Anything the configuration does not
// have keeps the order the API gave it. An empty prior state is an import,
// where there is no order to respect yet.
//
// A set would express "unordered" properly, and is the right fix whenever the
// schema can take a breaking change.
func OrderStatementsLikeState(prior interface{}, statements []map[string]interface{}) []map[string]interface{} {
	priorStatements := statementsFromState(prior)
	if len(priorStatements) == 0 {
		return statements
	}

	used := make([]bool, len(statements))
	out := make([]map[string]interface{}, 0, len(statements))

	for _, priorStatement := range priorStatements {
		key := statementKey(priorStatement)
		for i, statement := range statements {
			if !used[i] && statementKey(statement) == key {
				used[i] = true
				out = append(out, statement)
				break
			}
		}
	}

	for i, statement := range statements {
		if !used[i] {
			out = append(out, statement)
		}
	}
	return out
}

// statementsFromState digs the statement blocks out of a privilege as Terraform
// holds it: one document, which holds the list of statements.
//
// privilege is declared as a TypeSet -- only ever with one element, to give the
// block a name -- so a read gets a *schema.Set. The list form is accepted too,
// for callers that build the value by hand.
func statementsFromState(prior interface{}) []map[string]interface{} {
	var documents []interface{}
	switch typed := prior.(type) {
	case *schema.Set:
		documents = typed.List()
	case []interface{}:
		documents = typed
	}
	if len(documents) == 0 {
		return nil
	}
	document, ok := documents[0].(map[string]interface{})
	if !ok {
		return nil
	}

	raw, ok := document["statement"].([]interface{})
	if !ok {
		return nil
	}

	out := make([]map[string]interface{}, 0, len(raw))
	for _, item := range raw {
		if statement, ok := item.(map[string]interface{}); ok {
			out = append(out, statement)
		}
	}
	return out
}

// statementKey identifies a statement by its contents. The action and scope
// lists are sorted for the comparison only -- the values written to state are
// the ones the API returned, untouched.
//
// Untouched is safe: unlike the statements themselves, the API preserves the
// order of the values inside one. A statement created with
// ["users:List","users:Get","users:Update","users:Delete"] comes back in that
// order on every read. Sorting here just means a statement is still recognised
// if that ever stops being true.
func statementKey(statement map[string]interface{}) string {
	parts := []string{fmt.Sprint(statement["effect"])}

	for _, field := range []string{"action", "scope"} {
		values := toStrings(statement[field])
		sort.Strings(values)
		parts = append(parts, strings.Join(values, ","))
	}
	return strings.Join(parts, "|")
}

func toStrings(value interface{}) []string {
	items, ok := value.([]interface{})
	if !ok {
		return nil
	}

	out := make([]string, 0, len(items))
	for _, item := range items {
		out = append(out, fmt.Sprint(item))
	}
	return out
}

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
						Default:  DefaultVersion,
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
