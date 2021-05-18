package privilegeschema

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/privileges"
	"github.com/stretchr/testify/assert"
)

func mockRoleSetFn(i interface{}) int {
	return i.(int)
}

func mockUserSetFn(i interface{}) int {
	return i.(int)
}

func mockPrivilegeSetFn(i interface{}) int {
	return 0
}

func TestSchema(t *testing.T) {
	t.Run("creates and returns a map of a privilege Schema", func(t *testing.T) {
		provSchema := Schema()
		assert.NotNil(t, provSchema["name"])
		assert.NotNil(t, provSchema["description"])
		assert.NotNil(t, provSchema["user_ids"])
		assert.NotNil(t, provSchema["role_ids"])
		assert.NotNil(t, provSchema["privilege"])
	})
}

func TestInflate(t *testing.T) {
	tests := map[string]struct {
		ResourceData   map[string]interface{}
		ExpectedOutput privileges.Privilege
	}{
		"creates and returns the address of a privilege struct": {
			ResourceData: map[string]interface{}{
				"name":        "name",
				"description": "description",
				"role_ids":    schema.NewSet(mockRoleSetFn, []interface{}{1, 2, 3}),
				"user_ids":    schema.NewSet(mockUserSetFn, []interface{}{4, 5, 6}),
				"privilege": schema.NewSet(mockPrivilegeSetFn,
					[]interface{}{
						map[string]interface{}{
							"version": "version",
							"statement": []interface{}{
								map[string]interface{}{
									"effect": "allow",
									"action": []interface{}{"Apps:Create"},
									"scope":  []interface{}{"*"},
								},
							},
						},
					},
				),
			},
			ExpectedOutput: privileges.Privilege{
				Name:        oltypes.String("name"),
				Description: oltypes.String("description"),
				UserIDs:     []int{4, 5, 6},
				RoleIDs:     []int{1, 2, 3},
				Privilege: &privileges.PrivilegeData{
					Version: oltypes.String("version"),
					Statement: []privileges.StatementData{
						privileges.StatementData{
							Effect: oltypes.String("allow"),
							Action: []string{"Apps:Create"},
							Scope:  []string{"*"},
						},
					},
				},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj, _ := Inflate(test.ResourceData)
			assert.Equal(t, test.ExpectedOutput, subj)
		})
	}
}

func TestFlatten(t *testing.T) {
	tests := map[string]struct {
		InputData      privileges.Privilege
		ExpectedOutput []map[string]interface{}
	}{
		"creates and returns the address of a privilege struct": {
			InputData: privileges.Privilege{
				Name:        oltypes.String("name"),
				Description: oltypes.String("description"),
				UserIDs:     []int{4, 5, 6},
				RoleIDs:     []int{1, 2, 3},
				Privilege: &privileges.PrivilegeData{
					Version: oltypes.String("version"),
					Statement: []privileges.StatementData{
						privileges.StatementData{
							Effect: oltypes.String("allow"),
							Action: []string{"Apps:Create"},
							Scope:  []string{"*"},
						},
					},
				},
			},
			ExpectedOutput: []map[string]interface{}{
				map[string]interface{}{
					"version": "version",
					"statement": []map[string]interface{}{
						map[string]interface{}{
							"effect": "allow",
							"action": []string{"Apps:Create"},
							"scope":  []string{"*"},
						},
					},
				},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj := FlattenPrivilegeData(*test.InputData.Privilege)
			assert.Equal(t, test.ExpectedOutput, subj)
		})
	}
}
