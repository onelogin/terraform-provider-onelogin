package onelogin

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
	"github.com/stretchr/testify/assert"

	appconfigurationschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/app/configuration"
)

// Mock SDK client for testing
type mockSAMLAppSDK struct {
	onelogin.OneloginSDK
}

func (m *mockSAMLAppSDK) GetAppByID(id int, queryParams interface{}) (interface{}, error) {
	// Return a mock app with configuration
	return map[string]interface{}{
		"id":          id,
		"name":        "Test SAML App",
		"description": "Test Description",
		"configuration": map[string]interface{}{
			"signature_algorithm": "SHA-1",
		},
	}, nil
}

func TestConfigurationConsistency(t *testing.T) {
	// Test empty configuration map
	emptyConfig := map[string]interface{}{}
	flattened := appconfigurationschema.Flatten(emptyConfig)
	assert.NotNil(t, flattened, "Flattened empty config should not be nil")
	assert.Equal(t, map[string]interface{}{}, flattened, "Flattened empty config should be empty map")

	// Test SAML configuration
	samlConfig := map[string]interface{}{
		"signature_algorithm": "SHA-1",
	}
	flattened = appconfigurationschema.Flatten(samlConfig)
	assert.NotNil(t, flattened, "Flattened SAML config should not be nil")
	assert.Equal(t, "SHA-1", flattened["signature_algorithm"], "Signature algorithm should be preserved")

	// Test nil configuration
	var nilConfig map[string]interface{}
	flattened = appconfigurationschema.Flatten(nilConfig)
	assert.NotNil(t, flattened, "Flattened nil config should not be nil")
	assert.Equal(t, map[string]interface{}{}, flattened, "Flattened nil config should be empty map")
}

func TestSAMLAppReadConfigurationHandling(t *testing.T) {
	// Create a mock ResourceData
	r := SAMLApps().Schema
	d := schema.TestResourceDataRaw(t, r, map[string]interface{}{
		"name":        "Test SAML App",
		"description": "Test Description",
		"configuration": map[string]interface{}{
			"signature_algorithm": "SHA-1",
		},
	})
	d.SetId("12345")

	// Mock the OneLogin SDK client
	client := &mockSAMLAppSDK{}

	// Call samlAppRead with our mock data
	diags := samlAppRead(context.Background(), d, client)
	assert.Nil(t, diags, "samlAppRead should not return diagnostics")

	// Verify that configuration is set correctly
	config := d.Get("configuration").(map[string]interface{})
	assert.NotNil(t, config, "Configuration should not be nil")
	assert.Equal(t, "SHA-1", config["signature_algorithm"], "Signature algorithm should be preserved")
}
