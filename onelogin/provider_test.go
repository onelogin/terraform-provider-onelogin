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

	// Check for API URL or subdomain as fallback
	apiURL := os.Getenv("ONELOGIN_API_URL")
	subdomain := os.Getenv("ONELOGIN_SUBDOMAIN")

	if apiURL == "" && subdomain == "" {
		t.Fatal("ONELOGIN_API_URL must be set for acceptance tests")
	}

	// Warn if using subdomain instead of API URL
	if apiURL == "" && subdomain != "" {
		t.Logf("WARNING: Using ONELOGIN_SUBDOMAIN which is deprecated. Please switch to ONELOGIN_API_URL.")
	}

	// Set a longer timeout for tests (5 minutes) if not already set
	if os.Getenv("ONELOGIN_CLIENT_TIMEOUT") == "" {
		t.Logf("Setting ONELOGIN_CLIENT_TIMEOUT to 300 seconds for tests")
		os.Setenv("ONELOGIN_CLIENT_TIMEOUT", "300")
	} else {
		t.Logf("Using existing ONELOGIN_CLIENT_TIMEOUT: %s", os.Getenv("ONELOGIN_CLIENT_TIMEOUT"))
	}
}
