package onelogin

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
)

func TestAccOneLoginGroup_crud(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { TestAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOneLoginGroupDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOneLoginGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_groups.test", "name", "Test Group"),
					// The API drops reference rather than storing it, so this
					// holds only because the read leaves a configured value
					// alone instead of writing the null back over it.
					resource.TestCheckResourceAttr("onelogin_groups.test", "reference", "test-group-ref"),
				),
			},
			{
				Config: testAccCheckOneLoginGroupConfigUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_groups.test", "name", "Updated Test Group"),
					resource.TestCheckResourceAttr("onelogin_groups.test", "reference", "updated-group-ref"),
				),
			},
			{
				ResourceName:      "onelogin_groups.test",
				ImportState:       true,
				ImportStateVerify: true,
				// reference cannot survive an import: the API never returns
				// it, and an import has no prior state to preserve it from.
				// An imported group needs the attribute written back into the
				// configuration by hand, which is the best the API allows.
				ImportStateVerifyIgnore: []string{"reference"},
			},
		},
	})
}

func testAccCheckOneLoginGroupDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*onelogin.OneloginSDK)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "onelogin_groups" {
			continue
		}
		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}
		resp, err := client.GetGroupByIDV2(id)
		if err == nil && resp != nil {
			return fmt.Errorf("group %d still exists", id)
		}
	}
	return nil
}

const testAccCheckOneLoginGroupConfig = `
resource "onelogin_groups" "test" {
  name      = "Test Group"
  reference = "test-group-ref"
}
`

const testAccCheckOneLoginGroupConfigUpdated = `
resource "onelogin_groups" "test" {
  name      = "Updated Test Group"
  reference = "updated-group-ref"
}
`
