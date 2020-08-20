package onelogin

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccUser_crud(t *testing.T) {
	base := GetFixture("onelogin_user_example.tf", t)
	update := GetFixture("onelogin_user_updated_example.tf", t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: base,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_users.basic_test", "username", "testy.mctesterson"),
					resource.TestCheckResourceAttr("onelogin_users.basic_test", "email", "testy.mctesterson@onelogin.com"),
				),
			},
			{
				Config: update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_users.basic_test", "username", "boaty.mcboatface"),
					resource.TestCheckResourceAttr("onelogin_users.basic_test", "email", "boaty.mcboatface@onelogin.com"),
				),
			},
		},
	})
}
