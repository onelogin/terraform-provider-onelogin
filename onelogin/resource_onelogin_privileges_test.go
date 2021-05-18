package onelogin

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccPrivilege_crud(t *testing.T) {
	base := GetFixture("onelogin_privilege_example.tf", t)
	update := GetFixture("onelogin_privilege_updated_example.tf", t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: base,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_privileges.super_admin", "name", "super admin"),
					resource.TestCheckResourceAttr("onelogin_privileges.super_admin", "description", "description"),
					resource.TestCheckResourceAttr("onelogin_privileges.super_admin", "privilege.statement.0.action.0", "apps:List"),
					resource.TestCheckResourceAttr("onelogin_privileges.super_admin", "privilege.statement.0.action.1", "users:List"),
				),
			},
			{
				Config: update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_privileges.super_admin", "name", "super duper admin"),
					resource.TestCheckResourceAttr("onelogin_privileges.super_admin", "description", "description"),
					resource.TestCheckResourceAttr("onelogin_privileges.super_admin", "privilege.statement.0.action.0", "apps:List"),
					resource.TestCheckResourceAttr("onelogin_privileges.super_admin", "privilege.statement.0.action.1", "users:List"),
				),
			},
		},
	})
}
