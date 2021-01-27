package onelogin

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccSmartHook_crud(t *testing.T) {
	base := GetFixture("onelogin_smarthooks_example.tf", t)
	update := GetFixture("onelogin_smarthooks_updated_example.tf", t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: base,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_smarthook.basic_test", "type", "pre-authentication"),
					resource.TestCheckResourceAttr("onelogin_smarthook.basic_test", "env_vars.0", "API_KEY"),
					resource.TestCheckResourceAttr("onelogin_smarthook.basic_test", "retries", "0"),
					resource.TestCheckResourceAttr("onelogin_smarthook.basic_test", "timeout", "2"),
					resource.TestCheckResourceAttr("onelogin_smarthook.basic_test", "disabled", "false"),
					resource.TestCheckResourceAttr("onelogin_smarthook.basic_test", "status", "ready"),
					resource.TestCheckResourceAttr("onelogin_smarthook.basic_test", "risk_enabled", "false"),
					resource.TestCheckResourceAttr("onelogin_smarthook.basic_test", "location_enabled", "false"),
					resource.TestCheckResourceAttr("onelogin_smarthook.basic_test", "packages.mysql", "2.18.1"),
					resource.TestCheckResourceAttr("onelogin_smarthook.basic_test", "function", `function myFunc() {
            console.log('DING DONG')
          }`),
				),
			},
			{
				Config: update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_smarthook.basic_test", "type", "pre-authentication"),
					resource.TestCheckResourceAttr("onelogin_smarthook.basic_test", "env_vars.0", "API_KEY"),
					resource.TestCheckResourceAttr("onelogin_smarthook.basic_test", "retries", "0"),
					resource.TestCheckResourceAttr("onelogin_smarthook.basic_test", "timeout", "2"),
					resource.TestCheckResourceAttr("onelogin_smarthook.basic_test", "disabled", "false"),
					resource.TestCheckResourceAttr("onelogin_smarthook.basic_test", "status", "ready"),
					resource.TestCheckResourceAttr("onelogin_smarthook.basic_test", "risk_enabled", "false"),
					resource.TestCheckResourceAttr("onelogin_smarthook.basic_test", "location_enabled", "false"),
					resource.TestCheckResourceAttr("onelogin_smarthook.basic_test", "packages.mysql", "2.18.1"),
					resource.TestCheckResourceAttr("onelogin_smarthook.basic_test", "function", `function myFunc() {
            console.log('WOO WOO')
          }`),
				),
			},
		},
	})
}
