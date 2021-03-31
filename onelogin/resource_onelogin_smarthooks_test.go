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
					resource.TestCheckResourceAttr("onelogin_smarthooks.basic_test", "type", "pre-authentication"),
					resource.TestCheckResourceAttr("onelogin_smarthooks.basic_test", "retries", "0"),
					resource.TestCheckResourceAttr("onelogin_smarthooks.basic_test", "context_version", "1.0.0"),
					resource.TestCheckResourceAttr("onelogin_smarthooks.basic_test", "timeout", "1"),
					resource.TestCheckResourceAttr("onelogin_smarthooks.basic_test", "disabled", "false"),
					resource.TestCheckResourceAttr("onelogin_smarthooks.basic_test", "options.risk_enabled", "false"),
					resource.TestCheckResourceAttr("onelogin_smarthooks.basic_test", "options.location_enabled", "false"),
					resource.TestCheckResourceAttr("onelogin_smarthooks.basic_test", "packages.mysql", "2.18.1"),
					resource.TestCheckResourceAttr("onelogin_smarthooks.basic_test", "function", `ICAgIGV4cG9ydHMuaGFuZGxlciA9IGFzeW5jIGNvbnRleHQgPT4gewogICAgICBjb25zb2xlLmxvZygiUHJlLWF1dGggZXhlY3V0aW5nIGZvciAiICsgY29udGV4dC51c2VyLnVzZXJfaWRlbnRpZmllcik7CiAgICAgIHJldHVybiB7IHVzZXI6IGNvbnRleHQudXNlciB9OwogICAgfTsK`),
				),
			},
			{
				Config: update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_smarthooks.basic_test", "type", "pre-authentication"),
					resource.TestCheckResourceAttr("onelogin_smarthooks.basic_test", "retries", "0"),
					resource.TestCheckResourceAttr("onelogin_smarthooks.basic_test", "context_version", "1.0.0"),
					resource.TestCheckResourceAttr("onelogin_smarthooks.basic_test", "timeout", "1"),
					resource.TestCheckResourceAttr("onelogin_smarthooks.basic_test", "disabled", "false"),
					resource.TestCheckResourceAttr("onelogin_smarthooks.basic_test", "options.risk_enabled", "false"),
					resource.TestCheckResourceAttr("onelogin_smarthooks.basic_test", "options.location_enabled", "false"),
					resource.TestCheckResourceAttr("onelogin_smarthooks.basic_test", "packages.mysql", "2.18.1"),
					resource.TestCheckResourceAttr("onelogin_smarthooks.basic_test", "function", `ICAgIGV4cG9ydHMuaGFuZGxlciA9IGFzeW5jIGNvbnRleHQgPT4gewogICAgICBjb25zb2xlLmxvZygiUHJlLWF1dGggZXhlY3V0aW5nIGZvciAiICsgY29udGV4dC51c2VyLnVzZXJfaWRlbnRpZmllcik7CiAgICAgIHJldHVybiB7IHVzZXI6IGNvbnRleHQudXNlciB9OwogICAgfTsK`),
				),
			},
		},
	})
}
