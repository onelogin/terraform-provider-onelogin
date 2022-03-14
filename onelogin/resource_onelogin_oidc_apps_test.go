package onelogin

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOIDCApp_crud(t *testing.T) {
	base := GetFixture("onelogin_oidc_app_example.tf", t)
	update := GetFixture("onelogin_oidc_app_updated_example.tf", t)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			TestAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: base,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_oidc_apps.oidc", "name", "OIDC APP"),
					resource.TestCheckResourceAttr("onelogin_oidc_apps.oidc", "description", "OIDC"),
					resource.TestCheckResourceAttr("onelogin_oidc_apps.oidc", "connector_id", "108419"),
					resource.TestCheckResourceAttr("onelogin_oidc_apps.oidc", "configuration.redirect_uri", "https://localhost:3000/callback"),
					resource.TestCheckResourceAttr("onelogin_oidc_apps.oidc", "configuration.refresh_token_expiration_minutes", "1"),
					resource.TestCheckResourceAttr("onelogin_oidc_apps.oidc", "configuration.login_url", "https://www.test.com"),
					resource.TestCheckResourceAttr("onelogin_oidc_apps.oidc", "configuration.oidc_application_type", "0"),
					resource.TestCheckResourceAttr("onelogin_oidc_apps.oidc", "configuration.token_endpoint_auth_method", "1"),
					resource.TestCheckResourceAttr("onelogin_oidc_apps.oidc", "configuration.access_token_expiration_minutes", "1"),
				),
			},
			{
				Config: update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_oidc_apps.oidc", "name", "Updated OIDC APP"),
					resource.TestCheckResourceAttr("onelogin_oidc_apps.oidc", "description", "OIDC"),
					resource.TestCheckResourceAttr("onelogin_oidc_apps.oidc", "connector_id", "108419"),
					resource.TestCheckResourceAttr("onelogin_oidc_apps.oidc", "configuration.redirect_uri", "https://localhost:3000/callback"),
					resource.TestCheckResourceAttr("onelogin_oidc_apps.oidc", "configuration.refresh_token_expiration_minutes", "1"),
					resource.TestCheckResourceAttr("onelogin_oidc_apps.oidc", "configuration.login_url", "https://www.updated.com"),
					resource.TestCheckResourceAttr("onelogin_oidc_apps.oidc", "configuration.oidc_application_type", "0"),
					resource.TestCheckResourceAttr("onelogin_oidc_apps.oidc", "configuration.token_endpoint_auth_method", "1"),
					resource.TestCheckResourceAttr("onelogin_oidc_apps.oidc", "configuration.access_token_expiration_minutes", "1"),
				),
			},
		},
	})
}
