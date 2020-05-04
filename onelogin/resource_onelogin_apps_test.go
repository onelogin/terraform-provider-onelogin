package onelogin

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccApp_crud(t *testing.T) {
	base := GetFixture("onelogin_app_example.tf", t)
	update := GetFixture("onelogin_app_updated_example.tf", t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: base,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_apps.basic_test", "name", "Form-Based Fitbit App"),
				),
			},
			{
				Config: update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_apps.basic_test", "name", "Updated Form-Based Fitbit App"),
				),
			},
		},
	})
}
