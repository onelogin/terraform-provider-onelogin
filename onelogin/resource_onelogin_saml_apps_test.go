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

// checkParameterExists verifies that a parameter with the given key name exists in the SAML app
func checkParameterExists(resourceName, paramKeyName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Create a regex pattern to match the parameter set item
		pattern := regexp.MustCompile(fmt.Sprintf(`parameters\.\d+\.param_key_name\s*=\s*%s`, paramKeyName))

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource %s not found", resourceName)
		}

		// Iterate through all attributes to find matching parameters
		for k, v := range rs.Primary.Attributes {
			if pattern.MatchString(k) && v == paramKeyName {
				return nil
			}
		}

		return fmt.Errorf("Parameter with key_name %s not found in %s", paramKeyName, resourceName)
	}
}
