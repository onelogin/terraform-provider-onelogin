package onelogin

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccSmartHookEnvVar_crud(t *testing.T) {
	base := GetFixture("onelogin_smarthook_environment_variables_example.tf", t)
	update := GetFixture("onelogin_smarthook_environment_variables_updated_example.tf", t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: base,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_smarthook_environment_variables.api_key", "name", "SOME_KEY"),
					resource.TestCheckResourceAttr("onelogin_smarthook_environment_variables.api_key", "value", "123-456-789"),
				),
			},
			{
				Config: update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_smarthook_environment_variables.api_key", "name", "SOME_KEY"),
					resource.TestCheckResourceAttr("onelogin_smarthook_environment_variables.api_key", "value", "987-654-321"),
				),
			},
		},
	})
}
