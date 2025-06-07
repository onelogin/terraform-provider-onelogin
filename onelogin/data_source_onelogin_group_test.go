package onelogin

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOneLoginGroup(t *testing.T) {
	// This test requires a valid group ID to exist in the OneLogin account
	// For now, we'll skip it in CI and only run it locally with proper setup
	if testGroupID == "" {
		t.Skip("Skipping test as no test group ID is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOneLoginGroupDataSourceConfig(testGroupID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.onelogin_group.test", "id", testGroupID),
					resource.TestCheckResourceAttrSet("data.onelogin_group.test", "name"),
				),
			},
		},
	})
}

// testGroupID should be set to a valid group ID for testing
// In a real environment, this would be set via environment variables
var testGroupID = ""

func testAccCheckOneLoginGroupDataSourceConfig(groupID string) string {
	return fmt.Sprintf(`
data "onelogin_group" "test" {
  id = %s
}
`, groupID)
}
