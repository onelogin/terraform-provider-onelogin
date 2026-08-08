package onelogin

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
)

// TestApps tests the CRUD operations of the app resource
func TestApps(t *testing.T) {
	r := Apps()
	assert.NotNil(t, r)
	assert.NotNil(t, r.Schema)
	assert.NotNil(t, r.CreateContext)
	assert.NotNil(t, r.ReadContext)
	assert.NotNil(t, r.UpdateContext)
	assert.NotNil(t, r.DeleteContext)
}

// TestAppsSchema verifies the schema has required fields
func TestAppsSchema(t *testing.T) {
	schema := Apps().Schema

	// Verify required fields exist
	requiredFields := []string{"name", "connector_id"}
	for _, field := range requiredFields {
		assert.Contains(t, schema, field, "Schema is missing required field: %s", field)
	}
}

// testResourceData creates a ResourceData with the given attributes for testing
func testResourceData(t *testing.T, resourceType string, attrs map[string]interface{}) *schema.ResourceData {
	var r *schema.Resource
	switch resourceType {
	case "onelogin_apps":
		r = Apps()
	default:
		t.Fatalf("Unknown resource type: %s", resourceType)
	}

	return schema.TestResourceDataRaw(t, r.Schema, attrs)
}

// TestLogicalImplementation tests the resource implementation.
//
// The read is given the provider's own meta rather than nil. Handing it nil
// guaranteed the panic the test was written to catch -- every read opens with
// m.(*onelogin.OneloginSDK), and a type assertion on a nil interface panics --
// and it went unnoticed because Meta() is only non-nil once something has
// configured the provider. Run alone the body was skipped and the test passed;
// run after any acceptance test in the same process it failed. Terraform never
// calls a read with a nil meta, so what it caught was not a real state.
//
// With a real client and no ID, the read reaches the API and comes back with
// diagnostics, which is the thing worth pinning: a read for an app that is not
// there reports rather than crashes.
func TestLogicalImplementation(t *testing.T) {
	// Skip if this is not an acceptance test
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	ctx := context.Background()
	d := testResourceData(t, "onelogin_apps", map[string]interface{}{
		"name":         "Test App",
		"description":  "Test App Description",
		"connector_id": 123456,
	})

	appRes := Apps()

	// Verify that the function signatures are compatible with the schema interfaces
	var _ schema.CreateContextFunc = appRes.CreateContext
	var _ schema.ReadContextFunc = appRes.ReadContext
	var _ schema.UpdateContextFunc = appRes.UpdateContext
	var _ schema.DeleteContextFunc = appRes.DeleteContext

	// Configure the provider if nothing else in this process has. Relying on
	// another acceptance test to do it means running this one on its own
	// skips the whole point of it, silently -- which is how the nil meta it
	// used to pass went unnoticed for so long.
	m := testAccProvider.Meta()
	if m == nil {
		if diags := testAccProvider.Configure(ctx, terraform.NewResourceConfigRaw(nil)); diags.HasError() {
			t.Fatalf("could not configure the provider from the environment: %v", diags)
		}
		m = testAccProvider.Meta()
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Panic in implementation: %v", r)
		}
	}()

	var diags diag.Diagnostics
	diags = appRes.ReadContext(ctx, d, m)
	assert.NotNil(t, diags)
}
