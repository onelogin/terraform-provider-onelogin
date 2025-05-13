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
		ResourceData   map[string]interface{}
		ExpectedOutput models.Role
	}{
		"creates and returns the address of a role struct": {
			ResourceData: map[string]interface{}{
				"id":     "1",
				"name":   "name",
				"apps":   schema.NewSet(mockSetFn, []interface{}{1, 2, 3}),
				"users":  schema.NewSet(mockSetFn, []interface{}{4, 5, 6}),
				"admins": schema.NewSet(mockSetFn, []interface{}{4, 7}),
			},
			ExpectedOutput: models.Role{
				ID:     &id,
				Name:   &name,
				Apps:   []int32{1, 2, 3},
				Users:  []int32{4, 5, 6},
				Admins: []int32{4, 7},
			},
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

			// Compare slices directly
			assert.Equal(t, test.ExpectedOutput.Apps, subj.Apps)
			assert.Equal(t, test.ExpectedOutput.Users, subj.Users)
			assert.Equal(t, test.ExpectedOutput.Admins, subj.Admins)
		})
	}
}
