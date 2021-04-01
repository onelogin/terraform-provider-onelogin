package onelogin

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
