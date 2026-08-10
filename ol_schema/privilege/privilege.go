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
				out = append(out, orderValuesLikeState(priorStatement, statement))
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

// orderValuesLikeState puts the action and scope values of a matched statement
// back into the order the configuration has them in.
//
// The API does not preserve the order of the values inside a statement any more
// than it preserves the order of the statements themselves: a statement written
// as ["users:List","users:Get"] can be read back as ["users:Get","users:List"].
// action and scope are TypeLists, so that is a changed value, and because
// privilege is a TypeSet the changed value changes the element's hash -- which
// is why the plan proposes removing the whole privilege block and adding it
// back rather than editing it in place.
//
// Only values the prior statement already has are moved, and every value the
// API returned is kept, so this reorders without ever inventing or dropping
// one. A matched statement is copied rather than edited in place, so the
// caller's map is left as it was. Statements that match nothing are passed
// through as they arrived, and are the caller's own maps.
func orderValuesLikeState(prior, statement map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(statement))
	for key, value := range statement {
		out[key] = value
	}

	for _, field := range []string{"action", "scope"} {
		// A field the response did not give us as a list is left exactly as it
		// arrived, rather than becoming an empty one. The test is the same one
		// statementKey uses, so a statement can never match on its contents and
		// then quietly skip the reordering -- that would be a perpetual diff
		// with nothing to show for it.
		values, ok := toStrings(statement[field])
		if !ok {
			continue
		}
		priorValues, _ := toStrings(prior[field])
		out[field] = orderLikePrior(priorValues, values)
	}
	return out
}

// orderLikePrior returns returned, with the values prior names first and in
// prior's order, then anything left over in the order the API gave it.
func orderLikePrior(prior, returned []string) []interface{} {
	used := make([]bool, len(returned))
	out := make([]interface{}, 0, len(returned))

	for _, want := range prior {
		for i, got := range returned {
			if !used[i] && got == want {
				used[i] = true
				out = append(out, got)
				break
			}
		}
	}

	for i, got := range returned {
		if !used[i] {
			out = append(out, got)
		}
	}
	return out
}

// statementKey identifies a statement by its contents. The action and scope
// lists are sorted so that a statement is recognised as the same statement
// however the API happens to have ordered its values; orderValuesLikeState then
// restores the configuration's order in the value that reaches state.
//
// The separators are NUL so that a value containing one cannot forge a key:
// joined on ",", ["a,b"] and ["a","b"] both render as "a,b" and compare equal.
// No OneLogin action or scope contains a NUL.
func statementKey(statement map[string]interface{}) string {
	parts := []string{fmt.Sprint(statement["effect"])}

	for _, field := range []string{"action", "scope"} {
		values, _ := toStrings(statement[field])
		sort.Strings(values)
		parts = append(parts, strings.Join(values, "\x00"))
	}
	return strings.Join(parts, "\x00|\x00")
}

// toStrings reads the two shapes a statement's values arrive in: []interface{}
// from a decoded API response, and []string from FlattenPrivilegeData. It
// reports whether the value was one of them, so that callers can tell an empty
// list from something that is not a list at all -- keying those the same would
// make any two statements sharing an effect look identical.
func toStrings(value interface{}) ([]string, bool) {
	switch items := value.(type) {
	case []interface{}:
		out := make([]string, 0, len(items))
		for _, item := range items {
			out = append(out, fmt.Sprint(item))
		}
		return out, true
	case []string:
		return append([]string(nil), items...), true
	}
	return nil, false
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
