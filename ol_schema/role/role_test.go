package roleschema

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/roles"
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
	tests := map[string]struct {
		ResourceData   map[string]interface{}
		ExpectedOutput roles.Role
	}{
		"creates and returns the address of a role struct": {
			ResourceData: map[string]interface{}{
				"id":     "1",
				"name":   "name",
				"apps":   schema.NewSet(mockSetFn, []interface{}{1, 2, 3}),
				"users":  schema.NewSet(mockSetFn, []interface{}{4, 5, 6}),
				"admins": schema.NewSet(mockSetFn, []interface{}{4, 7}),
			},
			ExpectedOutput: roles.Role{
				ID:     oltypes.Int32(int32(1)),
				Name:   oltypes.String("name"),
				Apps:   []int32{1, 2, 3},
				Users:  []int32{4, 5, 6},
				Admins: []int32{4, 7},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj := Inflate(test.ResourceData)
			assert.Equal(t, test.ExpectedOutput, subj)
		})
	}
}
