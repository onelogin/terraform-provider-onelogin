package onelogin

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// The fixtures deliberately give one statement SEVERAL actions, listed out of
// alphabetical order. The API returns the values inside a statement in an order
// of its own, and a single-action statement cannot express that -- which is why
// this test passed throughout gh-254. Do not simplify them back to one action.
func TestAccPrivilege_crud(t *testing.T) {
	// One suffix across both steps: loaded separately, step 2 would get a
	// different one and replace every resource instead of updating them.
	fixtures := GetFixtures([]string{
		"onelogin_privilege_example.tf",
		"onelogin_privilege_updated_example.tf",
	}, t)
	base, update := fixtures[0], fixtures[1]

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: base,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_privileges.super_admin", "name", "super admin"),
					resource.TestCheckResourceAttr("onelogin_privileges.super_admin", "description", "description"),
					// privilege is a TypeSet, so its members are keyed by hash
					// rather than index: "privilege.statement.0..." can never
					// match. TestCheckTypeSetElemNestedAttrs searches the set.
					resource.TestCheckResourceAttr("onelogin_privileges.super_admin", "privilege.#", "1"),
					// The full action list, in the order the configuration
					// declares it. Asserting only element 0 cannot see the
					// reordering the API applies inside a statement.
					resource.TestCheckTypeSetElemNestedAttrs("onelogin_privileges.super_admin", "privilege.*", map[string]string{
						"statement.0.action.0": "apps:List",
						"statement.1.action.0": "users:List",
						"statement.1.action.1": "users:Get",
						"statement.1.action.2": "users:Update",
						"statement.1.action.3": "users:Delete",
					}),
				),
			},
			{
				Config: update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_privileges.super_admin", "name", "super duper admin"),
					resource.TestCheckResourceAttr("onelogin_privileges.super_admin", "description", "description"),
					resource.TestCheckResourceAttr("onelogin_privileges.super_admin", "privilege.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("onelogin_privileges.super_admin", "privilege.*", map[string]string{
						"statement.0.action.0": "apps:List",
						"statement.1.action.0": "users:List",
						"statement.1.action.1": "users:Get",
						"statement.1.action.2": "users:Update",
						"statement.1.action.3": "users:Delete",
					}),
				),
			},
		},
	})
}
