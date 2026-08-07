package onelogin

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

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
					resource.TestCheckTypeSetElemNestedAttrs("onelogin_privileges.super_admin", "privilege.*", map[string]string{
						"statement.0.action.0": "apps:List",
						"statement.1.action.0": "users:List",
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
					}),
				),
			},
		},
	})
}
