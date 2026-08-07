package onelogin

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider
var testAccProviderFactories map[string]func() (*schema.Provider, error)

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"onelogin": testAccProvider,
	}
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"onelogin": func() (*schema.Provider, error) {
			return Provider(), nil
		},
	}
}

// TestProvider checks the validity of a provider and stops further testing
// if a problem is found
func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

// TestAccPreCheck performs a check to ensure requisite credentials are in
// the environment and stops further testing if a problem is found
func TestAccPreCheck(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping acceptance test in short mode")
	}

	// Check for client credentials
	if v := os.Getenv("ONELOGIN_CLIENT_ID"); v == "" {
		t.Fatal("ONELOGIN_CLIENT_ID must be set for acceptance tests")
	}
	if v := os.Getenv("ONELOGIN_CLIENT_SECRET"); v == "" {
		t.Fatal("ONELOGIN_CLIENT_SECRET must be set for acceptance tests")
	}

	// ONELOGIN_API_URL, and only that. This used to accept a bare
	// ONELOGIN_SUBDOMAIN and log that it was deprecated, which read as though
	// the older setting still worked -- it does not. The provider's url is
	// Required, so a run with only a subdomain set gets past this check and
	// then fails every step with
	//
	//	Error: Missing required argument
	//	The argument "url" is required, but no definition was found.
	//
	// and the SDK, given only a subdomain, would aim at
	// <subdomain>.onelogin.com regardless of which tenant was meant.
	if v := os.Getenv("ONELOGIN_API_URL"); v == "" {
		t.Fatal("ONELOGIN_API_URL must be set for acceptance tests: the provider requires a url, and ONELOGIN_SUBDOMAIN alone cannot configure it")
	}

	// Set a longer timeout for tests (5 minutes) if not already set
	if os.Getenv("ONELOGIN_CLIENT_TIMEOUT") == "" {
		t.Logf("Setting ONELOGIN_CLIENT_TIMEOUT to 300 seconds for tests")
		os.Setenv("ONELOGIN_CLIENT_TIMEOUT", "300")
	} else {
		t.Logf("Using existing ONELOGIN_CLIENT_TIMEOUT: %s", os.Getenv("ONELOGIN_CLIENT_TIMEOUT"))
	}
}
