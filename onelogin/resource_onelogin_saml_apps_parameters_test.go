package onelogin

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
)

// TestAccSAMLApp_parameterChanges covers gh issue 274: changing an attribute of
// a parameter that is already on the app, and nothing else, planned no changes
// and never reached OneLogin.
//
// Parameters are a set, and the SDK compares a set by its hash codes before it
// looks at a single attribute. While the hash covered the key name alone, an
// element whose key had not changed always matched, so no attribute inside it
// could produce a diff. TestAccSAMLApp_crud changed
// default_values on an existing parameter and passed regardless, because the
// same step added another parameter and the new code forced the comparison.
// Every step here leaves the set of key names alone.
func TestAccSAMLApp_parameterChanges(t *testing.T) {
	const app = "onelogin_saml_apps.parameters"

	var appID int

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSAMLAppParametersConfig("email", `["one"]`),
				Check: resource.ComposeTestCheckFunc(
					checkParameterAttr(app, "email", "user_attribute_mappings", "email"),
					testAccCheckAppParameter(app, "email", "user_attribute_mappings", "email"),
					testAccCaptureIDs(map[string]*int{app: &appID}),
				),
			},
			{
				// The change the issue reported.
				Config: testAccSAMLAppParametersConfig("username", `["one"]`),
				Check: resource.ComposeTestCheckFunc(
					checkParameterAttr(app, "email", "user_attribute_mappings", "username"),
					testAccCheckAppParameter(app, "email", "user_attribute_mappings", "username"),
				),
			},
			{
				// A list attribute on an existing parameter, the change
				// TestAccSAMLApp_crud only appeared to cover.
				Config: testAccSAMLAppParametersConfig("username", `["one", "two"]`),
				Check: resource.ComposeTestCheckFunc(
					checkParameterValues(app, "groups", "one", "two"),
					testAccCheckAppParameter(app, "groups", "default_values", "[one two]"),
				),
			},
			{
				// An empty default_values goes out as null, which clears it.
				Config: testAccSAMLAppParametersConfig("username", `[]`),
				Check: resource.ComposeTestCheckFunc(
					checkParameterValues(app, "groups"),
					testAccCheckAppParameter(app, "groups", "default_values", "<nil>"),
				),
			},
			{
				// The same blind spot from the other side: a mapping changed
				// in OneLogin is drift, and the apply puts it back. The
				// configuration is the previous step's, so the drift is the
				// only thing to apply.
				PreConfig: func() {
					if err := testAccSetAppParameterMapping(appID, "email", "firstname"); err != nil {
						t.Fatal(err)
					}
				},
				Config: testAccSAMLAppParametersConfig("username", `[]`),
				Check:  testAccCheckAppParameter(app, "email", "user_attribute_mappings", "username"),
			},
		},
	})
}

// testAccSAMLAppParametersConfig is a SAML app with two parameters, the email
// mapping and the groups default_values written in as-is.
//
// groups spells out false for its flags, the way the documentation's examples
// do. The API is never sent a false -- the model drops it -- so this is what
// would show a perpetual diff if a plan started proposing values an apply
// cannot write. The test framework's empty-plan check after each step is what
// asserts that.
func testAccSAMLAppParametersConfig(mapping, groups string) string {
	return fmt.Sprintf(`
resource "onelogin_saml_apps" "parameters" {
  connector_id = 110016
  name         = "TF Acc parameters App"

  parameters {
    param_key_name            = "email"
    label                     = "Email"
    user_attribute_mappings   = %q
    include_in_saml_assertion = true
  }

  parameters {
    param_key_name            = "groups"
    label                     = "Groups"
    default_values            = %s
    include_in_saml_assertion = false
    skip_if_blank             = false
  }

  configuration = {
    signature_algorithm = "SHA-256"
  }
}
`, mapping, groups)
}

// checkParameterAttr asserts one attribute of the parameter with the given
// key name, found by key since parameters is a set.
func checkParameterAttr(resourceName, keyName, attr, want string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource %s not found", resourceName)
		}

		for k, v := range rs.Primary.Attributes {
			if !paramKeyName.MatchString(k) || v != keyName {
				continue
			}
			key := strings.TrimSuffix(k, "param_key_name") + attr
			if got := rs.Primary.Attributes[key]; got != want {
				return fmt.Errorf("expected %s of %s to be %q, got %q", attr, keyName, want, got)
			}
			return nil
		}

		return fmt.Errorf("Parameter with key_name %s not found in %s", keyName, resourceName)
	}
}

// testAccCheckAppParameter asserts what OneLogin itself holds for one attribute
// of a parameter, rather than what state says it holds. The bug this covers
// left the two agreeing with each other and both wrong, so state alone cannot
// tell a write that happened from one that was never attempted.
//
// want is compared against fmt.Sprint of the decoded JSON, so a list reads as
// [a b].
func testAccCheckAppParameter(resourceName, keyName, attr, want string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource %s not found", resourceName)
		}
		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		param, err := testAccAppParameter(testAccProvider.Meta().(*onelogin.OneloginSDK), id, keyName)
		if err != nil {
			return err
		}
		// Present, so that "<nil>" means the API said null rather than
		// that it left the attribute out.
		value, ok := param[attr]
		if !ok {
			return fmt.Errorf("expected OneLogin to report %s for %s on app %d, got %v", attr, keyName, id, param)
		}
		if got := fmt.Sprint(value); got != want {
			return fmt.Errorf("expected OneLogin to hold %s = %q for %s on app %d, got %q", attr, want, keyName, id, got)
		}
		return nil
	}
}

// testAccSetAppParameterMapping changes a parameter's mapping directly through
// the API, the way an admin in the OneLogin UI would, and confirms it took.
// Without that the step after it would pass whether or not there was any drift
// to correct.
func testAccSetAppParameterMapping(appID int, keyName, mapping string) error {
	client := testAccProvider.Meta().(*onelogin.OneloginSDK)
	current, err := testAccAppParameter(client, appID, keyName)
	if err != nil {
		return err
	}

	id, _ := current["id"].(float64)
	label, _ := current["label"].(string)
	params := map[string]models.Parameter{
		keyName: {ID: int(id), Label: label, UserAttributeMappings: mapping},
	}
	if _, err := client.UpdateApp(appID, models.App{Parameters: &params}); err != nil {
		return err
	}

	got, err := testAccAppParameter(client, appID, keyName)
	if err != nil {
		return err
	}
	if got["user_attribute_mappings"] != mapping {
		return fmt.Errorf("app %d parameter %s maps %v after setting %q", appID, keyName, got["user_attribute_mappings"], mapping)
	}
	return nil
}

func testAccAppParameter(client *onelogin.OneloginSDK, appID int, keyName string) (map[string]interface{}, error) {
	result, err := client.GetAppByID(appID, nil)
	if err != nil {
		return nil, err
	}
	appMap, ok := result.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected app response %T", result)
	}
	params, _ := appMap["parameters"].(map[string]interface{})
	param, ok := params[keyName].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("app %d has no parameter %s", appID, keyName)
	}
	return param, nil
}
