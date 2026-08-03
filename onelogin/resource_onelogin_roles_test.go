package onelogin

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestAccRole_crud(t *testing.T) {
	base := GetFixture("onelogin_role_example.tf", t)
	update := GetFixture("onelogin_role_updated_example.tf", t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: base,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_roles.executive_admin", "name", "executive admin"),
				),
			},
			{
				Config: update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_roles.executive_admin", "name", "updated executive admin"),
				),
			},
		},
	})
}

// TestExtractMemberIDs tests the extractMemberIDs helper that converts sub-endpoint
// object arrays (GET /api/2/roles/{id}/apps|users|admins) into flat integer ID slices.
// Sub-endpoints return objects like {"id": 1, "name": "..."} rather than bare IDs.
func TestExtractMemberIDs(t *testing.T) {
	tests := map[string]struct {
		input    interface{}
		expected []int
	}{
		"nil input returns empty slice": {
			input:    nil,
			expected: []int{},
		},
		"empty array returns empty slice": {
			input:    []interface{}{},
			expected: []int{},
		},
		"objects with id fields": {
			input: []interface{}{
				map[string]interface{}{"id": float64(1), "name": "app-one"},
				map[string]interface{}{"id": float64(2), "name": "app-two"},
			},
			expected: []int{1, 2},
		},
		"skips objects missing id field": {
			input: []interface{}{
				map[string]interface{}{"name": "no-id"},
				map[string]interface{}{"id": float64(5), "name": "has-id"},
			},
			expected: []int{5},
		},
		"non-slice input returns empty slice": {
			input:    "unexpected-string",
			expected: []int{},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := extractMemberIDs(tc.input)
			assert.Equal(t, tc.expected, got)
		})
	}
}
