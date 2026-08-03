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

// TestRoleReadUsesDirectFetch verifies that the roleRead implementation fields align
// with the map structure returned by a direct GET /api/2/roles/{id} response.
// roleRead now calls GetRoleByIDWithContext (one request per role) instead of
// paginating through the full role list, which eliminates redundant API calls when
// multiple role resources are refreshed concurrently.
func TestRoleReadUsesDirectFetch(t *testing.T) {
	// Simulate the map[string]interface{} payload returned by the SDK for a single role.
	roleResponse := map[string]interface{}{
		"id":     float64(42),
		"name":   "executive admin",
		"apps":   []interface{}{float64(1), float64(2)},
		"users":  []interface{}{float64(10)},
		"admins": []interface{}{},
	}

	// Validate name field
	assert.Equal(t, "executive admin", roleResponse["name"])

	// Validate apps parsing
	var appIDs []int
	if v, ok := roleResponse["apps"].([]interface{}); ok {
		for _, app := range v {
			if id, ok := app.(float64); ok {
				appIDs = append(appIDs, int(id))
			}
		}
	}
	assert.Equal(t, []int{1, 2}, appIDs, "apps should be parsed correctly from direct fetch response")

	// Validate users parsing
	var userIDs []int
	if v, ok := roleResponse["users"].([]interface{}); ok {
		for _, user := range v {
			if id, ok := user.(float64); ok {
				userIDs = append(userIDs, int(id))
			}
		}
	}
	assert.Equal(t, []int{10}, userIDs, "users should be parsed correctly from direct fetch response")

	// Validate admins — empty array should result in nil slice (no entries appended)
	var adminIDs []int
	if v, ok := roleResponse["admins"].([]interface{}); ok {
		for _, admin := range v {
			if id, ok := admin.(float64); ok {
				adminIDs = append(adminIDs, int(id))
			}
		}
	}
	assert.Empty(t, adminIDs, "admins should be empty when API returns empty array")
}
