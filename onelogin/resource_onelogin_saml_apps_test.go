package onelogin

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccSAMLApp_crud(t *testing.T) {
	base := GetFixture("onelogin_saml_app_example.tf", t)
	update := GetFixture("onelogin_saml_app_updated_example.tf", t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: base,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_saml_apps.saml", "name", "SAML App"),
					resource.TestCheckResourceAttr("onelogin_saml_apps.saml", "description", "SAML"),
					resource.TestCheckResourceAttr("onelogin_saml_apps.saml", "configuration.signature_algorithm", "SHA-1"),
					resource.TestCheckResourceAttr("onelogin_saml_apps.saml", "rules.0.name", "first rule"),
					resource.TestCheckResourceAttr("onelogin_saml_apps.saml", "rules.1.name", "second rule"),
				),
			},
			{
				Config: update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_saml_apps.saml", "name", "Updated SAML App"),
					resource.TestCheckResourceAttr("onelogin_saml_apps.saml", "description", "Updated SAML"),
					resource.TestCheckResourceAttr("onelogin_saml_apps.saml", "configuration.signature_algorithm", "SHA-256"),
					resource.TestCheckResourceAttr("onelogin_saml_apps.saml", "rules.0.name", "second rule"),
					resource.TestCheckResourceAttr("onelogin_saml_apps.saml", "rules.1.name", "first rule"),
				),
			},
		},
	})
}
