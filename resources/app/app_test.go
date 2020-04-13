package app

import (
	"testing"

	"github.com/onelogin/onelogin-go-sdk/pkg/models"
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-terraform-provider/resources/app/configuration"
	"github.com/stretchr/testify/assert"
)

func TestAppSchema(t *testing.T) {
	t.Run("creates and returns a map of an AppConfiguration Schema", func(t *testing.T) {
		schema := AppSchema()
		assert.NotNil(t, schema["name"])
		assert.NotNil(t, schema["visible"])
		assert.NotNil(t, schema["description"])
		assert.NotNil(t, schema["notes"])
		assert.NotNil(t, schema["icon_url"])
		assert.NotNil(t, schema["auth_method"])
		assert.NotNil(t, schema["policy_id"])
		assert.NotNil(t, schema["allow_assumed_signin"])
		assert.NotNil(t, schema["tab_id"])
		assert.NotNil(t, schema["connector_id"])
		assert.NotNil(t, schema["created_at"])
		assert.NotNil(t, schema["updated_at"])
		assert.NotNil(t, schema["provisioning"])
		assert.NotNil(t, schema["parameters"])
	})
}

func TestAddSubSchema(t *testing.T) {
	t.Run("adds sub schema to given resrouce schema", func(t *testing.T) {
		appSchema := AppSchema()
		AddSubSchema("sub", &appSchema, configuration.SAMLConfigurationSchema)
		assert.NotNil(t, appSchema["sub"])
	})
}

func TestInflateApp(t *testing.T) {
	tests := map[string]struct {
		ResourceData   map[string]interface{}
		ExpectedOutput *models.AppParameters
	}{
		"creates and returns the address of an AppParameters struct": {
			ResourceData: map[string]interface{}{
				"name":                 "test",
				"visible":              true,
				"description":          "test",
				"notes":                "test",
				"allow_assumed_signin": true,
				"connector_id":         123,
				"provisioning": map[string]interface{}{
					"enabled": true,
				},
				"parameters": map[string]interface{}{
					"test": map[string]interface{}{
						"user_attribute_mappings": "test",
					},
				},
			},
			ExpectedOutput: &models.App{
				Name:               oltypes.String("test"),
				Visible:            oltypes.Bool(true),
				Description:        oltypes.String("test"),
				Notes:              oltypes.String("test"),
				AllowAssumedSignin: oltypes.Bool(true),
				ConnectorID:        oltypes.Int32(123),
				Provisioning: &models.AppProvisioning{
					Enabled: oltypes.Bool(true),
				},
				Parameters: map[string]models.AppParameters{
					"test": models.AppParameters{
						UserAttributeMappings: oltypes.String("test"),
					},
				},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj := InflateApp(&test.ResourceData)
			assert.Equal(t, subj, test.ExpectedOutput)
		})
	}
}
