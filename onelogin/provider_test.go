package onelogin

import (
	"errors"
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
	err := accPreCheck()
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func accPreCheck() error {
	if v := os.Getenv("ONELOGIN_OAPI_URL"); v == "" {
		return errors.New("ONELOGIN_OAPI_URL must be set for acceptance tests")
	}
	if v := os.Getenv("ONELOGIN_CLIENT_ID"); v == "" {
		return errors.New("ONELOGIN_CLIENT_ID must be set for acceptance tests")
	}
	if v := os.Getenv("ONELOGIN_CLIENT_SECRET"); v == "" {
		return errors.New("ONELOGIN_CLIENT_SECRET must be set for acceptance tests")
	}

	return nil
}
