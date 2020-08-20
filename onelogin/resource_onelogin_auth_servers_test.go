package onelogin

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAuthServer_crud(t *testing.T) {
	base := GetFixture("onelogin_auth_server_example.tf", t)
	update := GetFixture("onelogin_auth_server_updated_example.tf", t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: base,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_auth_servers.test", "name", "test"),
					resource.TestCheckResourceAttr("onelogin_auth_servers.test", "description", "test"),
					resource.TestCheckResourceAttr("onelogin_auth_servers.test", "configuration.0.resource_identifier", "https://example.com/contacts"),
					resource.TestCheckResourceAttr("onelogin_auth_servers.test", "configuration.0.audiences.0", "https://example.com/contacts"),
					resource.TestCheckResourceAttr("onelogin_auth_servers.test", "configuration.0.refresh_token_expiration_minutes", "30"),
					resource.TestCheckResourceAttr("onelogin_auth_servers.test", "configuration.0.access_token_expiration_minutes", "10"),
				),
			},
			{
				Config: update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_auth_servers.test", "name", "updated"),
					resource.TestCheckResourceAttr("onelogin_auth_servers.test", "description", "updated test"),
					resource.TestCheckResourceAttr("onelogin_auth_servers.test", "configuration.0.resource_identifier", "https://example.com/users/contacts"),
					resource.TestCheckResourceAttr("onelogin_auth_servers.test", "configuration.0.audiences.0", "https://example.com/contacts"),
					resource.TestCheckResourceAttr("onelogin_auth_servers.test", "configuration.0.audiences.1", "https://example.com/users/contacts"),
					resource.TestCheckResourceAttr("onelogin_auth_servers.test", "configuration.0.refresh_token_expiration_minutes", "30"),
					resource.TestCheckResourceAttr("onelogin_auth_servers.test", "configuration.0.access_token_expiration_minutes", "10"),
				),
			},
		},
	})
}
