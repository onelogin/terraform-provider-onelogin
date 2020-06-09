package onelogin

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccUserMapping_crud(t *testing.T) {
	base := GetFixture("onelogin_user_mapping_example.tf", t)
	update := GetFixture("onelogin_user_mapping_updated_example.tf", t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: base,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_user_mappings.basic_test", "name", "Select Login"),
					resource.TestCheckResourceAttr("onelogin_user_mappings.basic_test", "enabled", "true"),
					resource.TestCheckResourceAttr("onelogin_user_mappings.basic_test", "match", "all"),
					resource.TestCheckResourceAttr("onelogin_user_mappings.basic_test", "actions.0.action", "set_status"),
					resource.TestCheckResourceAttr("onelogin_user_mappings.basic_test", "conditions.0.value", "90"),
					resource.TestCheckResourceAttr("onelogin_user_mappings.basic_test", "conditions.0.source", "last_login"),
					resource.TestCheckResourceAttr("onelogin_user_mappings.basic_test", "conditions.0.operator", ">"),
				),
			},
			{
				Config: update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_user_mappings.basic_test", "name", "Updated Login"),
					resource.TestCheckResourceAttr("onelogin_user_mappings.basic_test", "enabled", "true"),
					resource.TestCheckResourceAttr("onelogin_user_mappings.basic_test", "match", "all"),
					resource.TestCheckResourceAttr("onelogin_user_mappings.basic_test", "actions.0.action", "set_status"),
					resource.TestCheckResourceAttr("onelogin_user_mappings.basic_test", "conditions.0.value", "120"),
					resource.TestCheckResourceAttr("onelogin_user_mappings.basic_test", "conditions.0.source", "last_login"),
					resource.TestCheckResourceAttr("onelogin_user_mappings.basic_test", "conditions.0.operator", ">"),
				),
			},
		},
	})
}
