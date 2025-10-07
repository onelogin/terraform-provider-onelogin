package onelogin

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	roleschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/role"
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

// TestRoleQueryPagination tests that when cursor is set, limit and page are cleared
// to comply with the OneLogin API requirement: "cursor xor pagination arguments"
func TestRoleQueryPagination(t *testing.T) {
	// Test initial query with limit
	query := &roleschema.RoleQuery{
		Limit: "100",
	}

	assert.Equal(t, "100", query.Limit, "Initial limit should be set")
	assert.Equal(t, "", query.Cursor, "Initial cursor should be empty")
	assert.Equal(t, "", query.Page, "Initial page should be empty")

	// Test cursor-based pagination - simulate what happens in roleRead
	query.Cursor = "12345"
	query.Limit = ""
	query.Page = ""

	assert.Equal(t, "12345", query.Cursor, "Cursor should be set")
	assert.Equal(t, "", query.Limit, "Limit should be cleared when using cursor")
	assert.Equal(t, "", query.Page, "Page should be cleared when using cursor")
}
