package onelogin

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	brandschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/brand"
)

// brandData builds a ResourceData holding a brand, so the delete guard can be
// exercised without a client.
func brandData(t *testing.T, id string, master bool) *schema.ResourceData {
	t.Helper()

	d := schema.TestResourceDataRaw(t, brandschema.Schema(), map[string]interface{}{
		"name": "Engineering",
	})
	d.SetId(id)
	if err := d.Set("master", master); err != nil {
		t.Fatalf("set master: %v", err)
	}
	return d
}

// TestBrandDeleteRefusesTheMasterBrand covers the one destroy that must not
// reach the API.
//
// The master brand is the account's own branding and the fallback for every app
// naming no other brand. It is reachable here only by importing it, which
// somebody does to change a colour, not to delete it. Whether the endpoint
// would allow the delete was deliberately not tested -- that would have meant
// deleting a shared tenant's master brand -- and it does not change the answer:
// if the API refuses, this is the clearer error; if it allows, this is the
// difference between a destroy and an outage.
//
// The nil client is the assertion. brandDelete takes the SDK from m, so
// reaching the API at all would panic rather than quietly pass.
func TestBrandDeleteRefusesTheMasterBrand(t *testing.T) {
	diags := brandDelete(context.Background(), brandData(t, "1324", true), nil)

	if !diags.HasError() {
		t.Fatal("expected deleting the master brand to be refused")
	}

	message := diags[0].Summary
	if !strings.Contains(message, "master brand") {
		t.Errorf("the error should say why it was refused, got %q", message)
	}
	if !strings.Contains(message, "terraform state rm") {
		t.Errorf("the error should say what to do instead, got %q", message)
	}
	if !strings.Contains(message, "1324") {
		t.Errorf("the error should name the brand, got %q", message)
	}
}

// TestBrandDeleteAllowsAnOrdinaryBrand is the control: the guard has to be
// specific to the master brand, not a refusal to delete anything.
//
// A nil client would panic on the first SDK call, so recovering from that panic
// is how this asserts the guard let it through. Returning without a diagnostic
// would mean the guard fired for the wrong brand.
func TestBrandDeleteAllowsAnOrdinaryBrand(t *testing.T) {
	reached := false

	func() {
		defer func() {
			if recover() != nil {
				reached = true
			}
		}()
		_ = brandDelete(context.Background(), brandData(t, "7844", false), nil)
	}()

	if !reached {
		t.Fatal("expected an ordinary brand to reach the delete call rather than be refused by the master-brand guard")
	}
}
