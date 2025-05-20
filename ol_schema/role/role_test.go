package roleschema

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	"github.com/stretchr/testify/assert"
)

func mockSetFn(i interface{}) int {
	return i.(int)
}

func TestSchema(t *testing.T) {
	t.Run("creates and returns a map of a role Schema", func(t *testing.T) {
		provSchema := Schema()
		assert.NotNil(t, provSchema["name"])
		assert.NotNil(t, provSchema["apps"])
		assert.NotNil(t, provSchema["users"])
		assert.NotNil(t, provSchema["admins"])
	})
}

func TestInflate(t *testing.T) {
	id := int32(1)
	name := "name"

	tests := map[string]struct {
		ResourceData    map[string]interface{}
		ExpectedOutput  models.Role
		ExpectNilUsers  bool
		ExpectNilApps   bool
		ExpectNilAdmins bool
	}{
		"creates and returns the address of a role struct": {
			ResourceData: map[string]interface{}{
				"id":                   "1",
				"include_id_in_output": true,
				"name":                 "name",
				"apps":                 schema.NewSet(mockSetFn, []interface{}{1, 2, 3}),
				"users":                schema.NewSet(mockSetFn, []interface{}{4, 5, 6}),
				"admins":               schema.NewSet(mockSetFn, []interface{}{4, 7}),
			},
			ExpectedOutput: models.Role{
				ID:     &id,
				Name:   &name,
				Apps:   []int32{1, 2, 3},
				Users:  []int32{4, 5, 6},
				Admins: []int32{4, 7},
			},
			ExpectNilUsers:  false,
			ExpectNilApps:   false,
			ExpectNilAdmins: false,
		},
		"only name field provided": {
			ResourceData: map[string]interface{}{
				"id":                   "1",
				"include_id_in_output": true,
				"name":                 "name",
			},
			ExpectedOutput: models.Role{
				ID:   &id,
				Name: &name,
			},
			ExpectNilUsers:  true,
			ExpectNilApps:   true,
			ExpectNilAdmins: true,
		},
		"only name and users provided": {
			ResourceData: map[string]interface{}{
				"id":                   "1",
				"include_id_in_output": true,
				"name":                 "name",
				"users":                schema.NewSet(mockSetFn, []interface{}{4, 5, 6}),
			},
			ExpectedOutput: models.Role{
				ID:    &id,
				Name:  &name,
				Users: []int32{4, 5, 6},
			},
			ExpectNilUsers:  false,
			ExpectNilApps:   true,
			ExpectNilAdmins: true,
		},
		"only name and apps provided": {
			ResourceData: map[string]interface{}{
				"id":                   "1",
				"include_id_in_output": true,
				"name":                 "name",
				"apps":                 schema.NewSet(mockSetFn, []interface{}{1, 2, 3}),
			},
			ExpectedOutput: models.Role{
				ID:   &id,
				Name: &name,
				Apps: []int32{1, 2, 3},
			},
			ExpectNilUsers:  true,
			ExpectNilApps:   false,
			ExpectNilAdmins: true,
		},
		"only name and admins provided": {
			ResourceData: map[string]interface{}{
				"id":                   "1",
				"include_id_in_output": true,
				"name":                 "name",
				"admins":               schema.NewSet(mockSetFn, []interface{}{4, 7}),
			},
			ExpectedOutput: models.Role{
				ID:     &id,
				Name:   &name,
				Admins: []int32{4, 7},
			},
			ExpectNilUsers:  true,
			ExpectNilApps:   true,
			ExpectNilAdmins: false,
		},
		"name and without id in output": {
			ResourceData: map[string]interface{}{
				"id":                   "1",
				"include_id_in_output": false, // ID should not be in output
				"name":                 "name",
				"users":                schema.NewSet(mockSetFn, []interface{}{4, 5, 6}),
			},
			ExpectedOutput: models.Role{
				// ID field should be nil
				Name:  &name,
				Users: []int32{4, 5, 6},
			},
			ExpectNilUsers:  false,
			ExpectNilApps:   true,
			ExpectNilAdmins: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj := Inflate(test.ResourceData)

			// Compare pointer values
			if subj.ID != nil && test.ExpectedOutput.ID != nil {
				assert.Equal(t, *test.ExpectedOutput.ID, *subj.ID)
			}

			if subj.Name != nil && test.ExpectedOutput.Name != nil {
				assert.Equal(t, *test.ExpectedOutput.Name, *subj.Name)
			}

			// Check nil states
			if test.ExpectNilUsers {
				assert.Nil(t, subj.Users, "Users should be nil")
			} else {
				assert.NotNil(t, subj.Users, "Users should not be nil")
				assert.Equal(t, test.ExpectedOutput.Users, subj.Users)
			}

			if test.ExpectNilApps {
				assert.Nil(t, subj.Apps, "Apps should be nil")
			} else {
				assert.NotNil(t, subj.Apps, "Apps should not be nil")
				assert.Equal(t, test.ExpectedOutput.Apps, subj.Apps)
			}

			if test.ExpectNilAdmins {
				assert.Nil(t, subj.Admins, "Admins should be nil")
			} else {
				assert.NotNil(t, subj.Admins, "Admins should not be nil")
				assert.Equal(t, test.ExpectedOutput.Admins, subj.Admins)
			}
		})
	}
}
