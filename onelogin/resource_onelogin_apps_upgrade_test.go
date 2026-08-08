package onelogin

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// upgradeFromVersion is the last release where a parameter's values and
// default_values were strings. The first step below runs that provider from
// the registry so the state under test is the real thing rather than something
// hand-written to look like it.
const upgradeFromVersion = "0.16.0"

// TestAccSAMLApp_upgradeParameterValues proves the state upgrader is actually
// wired up, which the unit tests on UpgradeParameterValuesV0 cannot: they call
// the function directly and say nothing about whether SchemaVersion,
// StateUpgraders or the v0 type are right.
//
// Step one applies with the published provider, writing default_values as a
// string. Step two hands that state to the provider under test and asks for a
// plan. An upgrader that is missing, or whose v0 type does not match what was
// written, fails here -- and it fails the way a practitioner would experience
// it, on the first plan after upgrading.
//
// This is the one test that needs the registry. If it fails with
//
//	no available releases match the given constraints 0.16.0
//
// while https://registry.terraform.io/v1/providers/onelogin/onelogin/versions
// plainly lists the version, the cause is almost certainly a stale
// ~/.terraform.d/plugins directory: Terraform treats it as an implied
// filesystem mirror, and a provider found there is never looked up in the
// registry at all. Check for onelogin/onelogin under it, and run with
//
//	TF_CLI_CONFIG_FILE=<file containing: provider_installation { direct {} }>
//
// or clear the stale copies out.
func TestAccSAMLApp_upgradeParameterValues(t *testing.T) {
	suffix := acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)

	// One value, so the string that v0.16.0 stores and the single-element list
	// this version stores are the same value expressed two ways. That is the
	// case the upgrade has to get right for everyone who is already using the
	// provider.
	v0Config := fmt.Sprintf(`
resource "onelogin_saml_apps" "upgrade" {
  connector_id = 50534
  name         = "Upgrade Probe %s"

  parameters {
    param_key_name = "upgraded"
    label          = "Upgraded"
    default_values = "one"
  }

  configuration = {
    signature_algorithm = "SHA-1"
  }
}
`, suffix)

	v1Config := fmt.Sprintf(`
resource "onelogin_saml_apps" "upgrade" {
  connector_id = 50534
  name         = "Upgrade Probe %s"

  parameters {
    param_key_name = "upgraded"
    label          = "Upgraded"
    default_values = ["one"]
  }

  configuration = {
    signature_algorithm = "SHA-1"
  }
}
`, suffix)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"onelogin": {
						Source:            "onelogin/onelogin",
						VersionConstraint: upgradeFromVersion,
					},
				},
				Config: v0Config,
				Check:  checkParameterExists("onelogin_saml_apps.upgrade", "upgraded"),
			},
			{
				ProviderFactories: testAccProviderFactories,
				Config:            v1Config,
				// The whole point: the same value, now written as a list,
				// must plan to no change at all. A non-empty plan here is an
				// upgrade that rewrote somebody's parameter.
				PlanOnly: true,
			},
			{
				ProviderFactories: testAccProviderFactories,
				Config:            v1Config,
				Check:             checkParameterValues("onelogin_saml_apps.upgrade", "upgraded", "one"),
			},
		},
	})
}
