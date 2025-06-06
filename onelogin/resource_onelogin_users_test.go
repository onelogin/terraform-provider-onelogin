package onelogin

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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

func TestAccUser_trustedIdp(t *testing.T) {
	withIdp := GetFixture("onelogin_user_with_trusted_idp_example.tf", t)
	withoutIdp := GetFixture("onelogin_user_without_trusted_idp_example.tf", t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: withIdp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_users.trusted_idp_test", "username", "trusted.idp.test"),
					resource.TestCheckResourceAttr("onelogin_users.trusted_idp_test", "email", "trusted.idp.test@onelogin.com"),
					resource.TestCheckResourceAttr("onelogin_users.trusted_idp_test", "trusted_idp_id", "123456"),
				),
			},
			{
				Config: withoutIdp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_users.trusted_idp_test", "username", "trusted.idp.test"),
					resource.TestCheckResourceAttr("onelogin_users.trusted_idp_test", "email", "trusted.idp.test@onelogin.com"),
					resource.TestCheckNoResourceAttr("onelogin_users.trusted_idp_test", "trusted_idp_id"),
				),
			},
		},
	})
}
