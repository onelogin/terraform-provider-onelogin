package onelogin

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"onelogin": testAccProvider,
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

	if v := os.Getenv("ONELOGIN_OAPI_URL"); v == "" {
		t.Fatal("ONELOGIN_OAPI_URL must be set for acceptance tests")
	}
	if v := os.Getenv("ONELOGIN_CLIENT_ID"); v == "" {
		t.Fatal("ONELOGIN_CLIENT_ID must be set for acceptance tests")
	}
	if v := os.Getenv("ONELOGIN_CLIENT_SECRET"); v == "" {
		t.Fatal("ONELOGIN_CLIENT_SECRET must be set for acceptance tests")
	}
}
