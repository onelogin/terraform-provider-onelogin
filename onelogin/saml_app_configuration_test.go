package onelogin

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
	"github.com/stretchr/testify/assert"

	appconfigurationschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/app/configuration"
)

// Mock SDK client for testing
type mockSAMLAppSDK struct {
	// Embed the interface to satisfy the contract but only implement what we need
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

// TestSAMLAppReadConfigurationHandlingMock tests the configuration handling in samlAppRead
// using a custom implementation that doesn't rely on the actual samlAppRead function
func TestSAMLAppReadConfigurationHandlingMock(t *testing.T) {
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

	// Get mock app data
	mockSDK := &mockSAMLAppSDK{}
	aid, _ := strconv.Atoi(d.Id())
	appData, err := mockSDK.GetAppByID(aid, nil)
	assert.NoError(t, err, "GetAppByID should not return an error")
	
	// Extract configuration from app data
	appMap, ok := appData.(map[string]interface{})
	assert.True(t, ok, "App data should be a map")
	
	if v, ok := appMap["configuration"]; ok {
		if configData, ok := v.(map[string]interface{}); ok {
			flattenedConfig := appconfigurationschema.Flatten(configData)
			// Set the configuration field
			d.Set("configuration", flattenedConfig)
		} else {
			d.Set("configuration", map[string]interface{}{})
		}
	} else {
		d.Set("configuration", map[string]interface{}{})
	}

	// Verify that configuration is set correctly
	config := d.Get("configuration").(map[string]interface{})
	assert.NotNil(t, config, "Configuration should not be nil")
	assert.Equal(t, "SHA-1", config["signature_algorithm"], "Signature algorithm should be preserved")
}
