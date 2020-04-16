package onelogin

import (
	"errors"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"onelogin": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
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

func TestAccPreCheck(t *testing.T) {
	err := accPreCheck()
	if err != nil {
		t.Fatalf("%v", err)
	}
}
