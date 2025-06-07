package onelogin

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOneLoginGroups(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOneLoginGroupsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.onelogin_groups.groups", "groups.#"),
				),
			},
		},
	})
}

const testAccCheckOneLoginGroupsDataSourceConfig = `
data "onelogin_groups" "groups" {}
`
