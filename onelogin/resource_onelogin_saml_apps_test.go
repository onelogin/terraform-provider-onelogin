package onelogin

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
					// Check that the parameters exist and have the correct values
					checkParameterExists("onelogin_saml_apps.saml", "email"),
					checkParameterExists("onelogin_saml_apps.saml", "firstname"),
					checkParameterExists("onelogin_saml_apps.saml", "lastname"),
					checkParameterExists("onelogin_saml_apps.saml", "department"),
				),
			},
			{
				Config: update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_saml_apps.saml", "name", "Updated SAML App"),
					resource.TestCheckResourceAttr("onelogin_saml_apps.saml", "description", "Updated SAML"),
					resource.TestCheckResourceAttr("onelogin_saml_apps.saml", "configuration.signature_algorithm", "SHA-256"),
					// Check that the parameters exist and have the correct values
					checkParameterExists("onelogin_saml_apps.saml", "email"),
					checkParameterExists("onelogin_saml_apps.saml", "firstname"),
					checkParameterExists("onelogin_saml_apps.saml", "lastname"),
					checkParameterExists("onelogin_saml_apps.saml", "department"),
					checkParameterExists("onelogin_saml_apps.saml", "title"),
				),
			},
		},
	})
}

// paramKeyName matches the state address of a parameter's key name. parameters
// is a TypeSet, so the middle element is the set hash rather than an index.
//
// The pattern used to end in `\s*=\s*<name>` and was matched against the
// attribute key, which is only ever an address -- so it matched nothing and no
// parameter was ever found. The name belongs in the value comparison below,
// which was already there and never reached.
var paramKeyName = regexp.MustCompile(`^parameters\.\d+\.param_key_name$`)

// checkParameterExists verifies that a parameter with the given key name exists in the SAML app
func checkParameterExists(resourceName, keyName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource %s not found", resourceName)
		}

		for k, v := range rs.Primary.Attributes {
			if paramKeyName.MatchString(k) && v == keyName {
				return nil
			}
		}

		return fmt.Errorf("Parameter with key_name %s not found in %s", keyName, resourceName)
	}
}
