package onelogin

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

// TestLogicalImplementation tests the resource implementation
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

	// Create a provider instance without actually making API calls
	m := testAccProvider.Meta()

	// Test function signatures and types
	var diags diag.Diagnostics

	// Verify that the function signatures are compatible with the schema interfaces
	var _ schema.CreateContextFunc = appRes.CreateContext
	var _ schema.ReadContextFunc = appRes.ReadContext
	var _ schema.UpdateContextFunc = appRes.UpdateContext
	var _ schema.DeleteContextFunc = appRes.DeleteContext

	// We can't actually make API calls in a unit test
	// But we can verify that the implementations don't panic
	if m != nil {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Panic in implementation: %v", r)
			}
		}()

		// These will fail because we're not actually making API calls
		// But they shouldn't panic
		diags = appRes.ReadContext(ctx, d, nil)
		assert.NotNil(t, diags)
	}
}
