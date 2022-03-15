package onelogin

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAppRule_crud(t *testing.T) {
	base := GetFixture("onelogin_app_rules_example.tf", t)
	update := GetFixture("onelogin_app_rules_updated_example.tf", t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: base,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_app_rules.test", "name", "first rule"),
					resource.TestCheckResourceAttr("onelogin_app_rules.test", "conditions.0.source", "last_login"),
					resource.TestCheckResourceAttr("onelogin_app_rules.test", "conditions.0.operator", ">"),
					resource.TestCheckResourceAttr("onelogin_app_rules.test", "actions.0.action", "set_amazonusername"),

					resource.TestCheckResourceAttr("onelogin_app_rules.check", "name", "second rule"),
					resource.TestCheckResourceAttr("onelogin_app_rules.check", "conditions.0.source", "has_role"),
					resource.TestCheckResourceAttr("onelogin_app_rules.check", "conditions.0.operator", "ri"),
					resource.TestCheckResourceAttr("onelogin_app_rules.check", "actions.0.action", "set_amazonusername"),
				),
			},
			{
				Config: update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_app_rules.test", "name", "updated first rule"),
					resource.TestCheckResourceAttr("onelogin_app_rules.test", "conditions.0.source", "last_login"),
					resource.TestCheckResourceAttr("onelogin_app_rules.test", "conditions.0.operator", "<"),
					resource.TestCheckResourceAttr("onelogin_app_rules.test", "actions.0.action", "set_amazonusername"),

					resource.TestCheckResourceAttr("onelogin_app_rules.check", "name", "updated second rule"),
					resource.TestCheckResourceAttr("onelogin_app_rules.check", "conditions.0.source", "has_role"),
					resource.TestCheckResourceAttr("onelogin_app_rules.check", "conditions.0.operator", "ri"),
					resource.TestCheckResourceAttr("onelogin_app_rules.check", "actions.0.action", "set_amazonusername"),
				),
			},
		},
	})
}
