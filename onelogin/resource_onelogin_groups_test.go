package onelogin

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOneLoginGroup_crud(t *testing.T) {
	// Skip this test for now since the OneLogin API doesn't support group creation/update/deletion
	t.Skip("Skipping test as OneLogin API doesn't support group CRUD operations")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOneLoginGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_groups.test", "name", "Test Group"),
					resource.TestCheckResourceAttr("onelogin_groups.test", "reference", "test-group"),
				),
			},
			{
				Config: testAccCheckOneLoginGroupConfigUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_groups.test", "name", "Updated Test Group"),
					resource.TestCheckResourceAttr("onelogin_groups.test", "reference", "updated-test-group"),
				),
			},
		},
	})
}

const testAccCheckOneLoginGroupConfig = `
resource "onelogin_groups" "test" {
  name      = "Test Group"
  reference = "test-group"
}
`

const testAccCheckOneLoginGroupConfigUpdated = `
resource "onelogin_groups" "test" {
  name      = "Updated Test Group"
  reference = "updated-test-group"
}
`
