package privilegeschema

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
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
	// Create test variables as pointers
	name := "name"
	description := "description"
	version := "version"
	effect := "allow"

	tests := map[string]struct {
		ResourceData   map[string]interface{}
		ExpectedOutput models.Privilege
	}{
		"creates and returns the address of a privilege struct": {
			ResourceData: map[string]interface{}{
				"name":        name,
				"description": description,
				"role_ids":    schema.NewSet(mockRoleSetFn, []interface{}{1, 2, 3}),
				"user_ids":    schema.NewSet(mockUserSetFn, []interface{}{4, 5, 6}),
				"privilege": schema.NewSet(mockPrivilegeSetFn,
					[]interface{}{
						map[string]interface{}{
							"version": version,
							"statement": []interface{}{
								map[string]interface{}{
									"effect": effect,
									"action": []interface{}{"Apps:Create"},
									"scope":  []interface{}{"*"},
								},
							},
						},
					},
				),
			},
			ExpectedOutput: models.Privilege{
				Name:        &name,
				Description: &description,
				UserIDs:     []int{4, 5, 6},
				RoleIDs:     []int{1, 2, 3},
				Privilege: &models.PrivilegeData{
					Version: &version,
					Statement: []models.StatementData{
						{
							Effect: &effect,
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

			// Compare the pointers correctly
			if subj.Name != nil && test.ExpectedOutput.Name != nil {
				assert.Equal(t, *subj.Name, *test.ExpectedOutput.Name)
			}
			if subj.Description != nil && test.ExpectedOutput.Description != nil {
				assert.Equal(t, *subj.Description, *test.ExpectedOutput.Description)
			}

			assert.Equal(t, subj.UserIDs, test.ExpectedOutput.UserIDs)
			assert.Equal(t, subj.RoleIDs, test.ExpectedOutput.RoleIDs)

			if subj.Privilege != nil && test.ExpectedOutput.Privilege != nil {
				if subj.Privilege.Version != nil && test.ExpectedOutput.Privilege.Version != nil {
					assert.Equal(t, *subj.Privilege.Version, *test.ExpectedOutput.Privilege.Version)
				}

				if len(subj.Privilege.Statement) > 0 && len(test.ExpectedOutput.Privilege.Statement) > 0 {
					if subj.Privilege.Statement[0].Effect != nil && test.ExpectedOutput.Privilege.Statement[0].Effect != nil {
						assert.Equal(t, *subj.Privilege.Statement[0].Effect, *test.ExpectedOutput.Privilege.Statement[0].Effect)
					}
					assert.Equal(t, subj.Privilege.Statement[0].Action, test.ExpectedOutput.Privilege.Statement[0].Action)
					assert.Equal(t, subj.Privilege.Statement[0].Scope, test.ExpectedOutput.Privilege.Statement[0].Scope)
				}
			}
		})
	}
}

func TestFlatten(t *testing.T) {
	// Create test variables as pointers
	version := "version"
	effect := "allow"

	tests := map[string]struct {
		InputData      models.PrivilegeData
		ExpectedOutput []map[string]interface{}
	}{
		"creates and returns the address of a privilege struct": {
			InputData: models.PrivilegeData{
				Version: &version,
				Statement: []models.StatementData{
					{
						Effect: &effect,
						Action: []string{"Apps:Create"},
						Scope:  []string{"*"},
					},
				},
			},
			ExpectedOutput: []map[string]interface{}{
				{
					"version": version,
					"statement": []map[string]interface{}{
						{
							"effect": effect,
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
			subj := FlattenPrivilegeData(test.InputData)
			assert.Equal(t, test.ExpectedOutput, subj)
		})
	}
}
