package onelogin

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccUserCustomAttributes_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckUserCustomAttributesConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_user_custom_attributes.test_attr", "name", "Test Attribute"),
					resource.TestCheckResourceAttr("onelogin_user_custom_attributes.test_attr", "shortname", "test_attr"),
				),
			},
		},
	})
}

func testAccCheckUserCustomAttributesConfig() string {
	return `
resource onelogin_user_custom_attributes test_attr {
  name      = "Test Attribute"
  shortname = "test_attr"
  position  = 1
}
`
}

func TestAccUserCustomAttributesWithUser_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckUserCustomAttributesWithUserConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_users.test_user", "username", "test.user.for.attrs"),
					resource.TestCheckResourceAttr("onelogin_user_custom_attributes.test_attr", "name", "Test Attribute"),
					resource.TestCheckResourceAttr("onelogin_user_custom_attributes.test_attr", "shortname", "test_attr"),
					resource.TestCheckResourceAttrSet("onelogin_user_custom_attributes.user_attr_value", "user_id"),
					resource.TestCheckResourceAttr("onelogin_user_custom_attributes.user_attr_value", "value", "test_value"),
				),
			},
		},
	})
}

func testAccCheckUserCustomAttributesWithUserConfig() string {
	return `
resource onelogin_users test_user {
  username = "test.user.for.attrs"
  email    = "test.user.attrs@example.com"
}

resource onelogin_user_custom_attributes test_attr {
  name      = "Test Attribute"
  shortname = "test_attr"
  position  = 1
}

resource onelogin_user_custom_attributes user_attr_value {
  user_id   = onelogin_users.test_user.id
  shortname = onelogin_user_custom_attributes.test_attr.shortname
  value     = "test_value"
}
`
}