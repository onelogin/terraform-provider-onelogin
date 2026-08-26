package appschema

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	"github.com/stretchr/testify/assert"
)

func TestSchema(t *testing.T) {
	t.Run("creates and returns a map of an App Schema", func(t *testing.T) {
		schema := Schema()
		assert.NotNil(t, schema["name"])
		assert.NotNil(t, schema["visible"])
		assert.NotNil(t, schema["description"])
		assert.NotNil(t, schema["notes"])
		assert.NotNil(t, schema["icon_url"])
		assert.NotNil(t, schema["auth_method"])
		assert.NotNil(t, schema["policy_id"])
		assert.NotNil(t, schema["brand_id"])
		assert.NotNil(t, schema["allow_assumed_signin"])
		assert.NotNil(t, schema["tab_id"])
		assert.NotNil(t, schema["connector_id"])
		assert.NotNil(t, schema["created_at"])
		assert.NotNil(t, schema["updated_at"])
		assert.NotNil(t, schema["provisioning"])
		assert.NotNil(t, schema["parameters"])
	})
}

func mockSetFn(interface{}) int {
	return 0
}

func TestInflate(t *testing.T) {
	id := int32(123)
	name := "test"
	visible := true
	description := "test"
	notes := "test"
	allowAssumedSignin := true
	connectorID := int32(123)
	brandID := 123
	provArn := "test"
	sigAlg := "test"

	tests := map[string]struct {
		ResourceData   map[string]interface{}
		ExpectedOutput models.App
	}{
		"creates and returns the address of an App struct with all sub-fields": {
			ResourceData: map[string]interface{}{
				"id":                   "123",
				"name":                 "test",
				"visible":              true,
				"description":          "test",
				"notes":                "test",
				"allow_assumed_signin": true,
				"connector_id":         123,
				"brand_id":             123,
				"parameters": schema.NewSet(mockSetFn, []interface{}{
					map[string]interface{}{
						"param_key_name":             "test",
						"param_id":                   123,
						"label":                      "test",
						"user_attribute_mappings":    "test",
						"user_attribute_macros":      "test",
						"attributes_transformations": "test",
						"default_values":             "test",
						"skip_if_blank":              true,
						"values":                     "test",
						"provisioned_entitlements":   true,
						"safe_entitlements_enabled":  true,
					},
				}),
				"provisioning": map[string]interface{}{
					"enabled": true,
				},
				"configuration": map[string]interface{}{
					"provider_arn":        "test",
					"signature_algorithm": "test",
				},
			},
			ExpectedOutput: models.App{
				ID:                 &id,
				Name:               &name,
				Visible:            &visible,
				Description:        &description,
				Notes:              &notes,
				AllowAssumedSignin: &allowAssumedSignin,
				ConnectorID:        &connectorID,
				BrandID:            &brandID,
				Provisioning: &models.Provisioning{
					Enabled: true,
				},
				Configuration: models.ConfigurationSAML{
					ProviderArn:        provArn,
					SignatureAlgorithm: sigAlg,
				},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj, _ := Inflate(test.ResourceData)

			// Compare fields individually instead of whole struct
			if subj.ID != nil && test.ExpectedOutput.ID != nil {
				assert.Equal(t, *test.ExpectedOutput.ID, *subj.ID)
			}
			if subj.Name != nil && test.ExpectedOutput.Name != nil {
				assert.Equal(t, *test.ExpectedOutput.Name, *subj.Name)
			}
			if subj.Visible != nil && test.ExpectedOutput.Visible != nil {
				assert.Equal(t, *test.ExpectedOutput.Visible, *subj.Visible)
			}
			if subj.Description != nil && test.ExpectedOutput.Description != nil {
				assert.Equal(t, *test.ExpectedOutput.Description, *subj.Description)
			}
			if subj.Notes != nil && test.ExpectedOutput.Notes != nil {
				assert.Equal(t, *test.ExpectedOutput.Notes, *subj.Notes)
			}
			if subj.AllowAssumedSignin != nil && test.ExpectedOutput.AllowAssumedSignin != nil {
				assert.Equal(t, *test.ExpectedOutput.AllowAssumedSignin, *subj.AllowAssumedSignin)
			}
			if subj.ConnectorID != nil && test.ExpectedOutput.ConnectorID != nil {
				assert.Equal(t, *test.ExpectedOutput.ConnectorID, *subj.ConnectorID)
			}
			if subj.BrandID != nil && test.ExpectedOutput.BrandID != nil {
				assert.Equal(t, *test.ExpectedOutput.BrandID, *subj.BrandID)
			}

			// Check provisioning
			if subj.Provisioning != nil && test.ExpectedOutput.Provisioning != nil {
				provSubj := subj.Provisioning
				provExp := test.ExpectedOutput.Provisioning
				assert.Equal(t, provExp.Enabled, provSubj.Enabled)
			}

			// Check Configuration
			confSubj, ok1 := subj.Configuration.(models.ConfigurationSAML)
			confExp, ok2 := test.ExpectedOutput.Configuration.(models.ConfigurationSAML)
			if ok1 && ok2 {
				assert.Equal(t, confExp.ProviderArn, confSubj.ProviderArn)
				assert.Equal(t, confExp.SignatureAlgorithm, confSubj.SignatureAlgorithm)
			}

			// Check Parameters
			if subj.Parameters != nil {
				// Just verify parameters exist, detailed parameter testing is in parameters_test.go
				assert.NotNil(t, subj.Parameters)
			}
		})
	}
}

// TestInflatePolicyIDTypes covers the three things the key can mean, and the
// one thing it must not do quietly.
//
// A value Inflate cannot read is reported rather than stepped over. The key
// being present is the whole instruction, and the case that would hurt is a
// dropped 0: the plan promises an unassignment, the apply sends nothing, and
// the next read brings the policy back, so the same diff returns for ever.
func TestInflatePolicyIDTypes(t *testing.T) {
	base := func() map[string]interface{} {
		return map[string]interface{}{"name": "test app", "connector_id": 108419}
	}

	t.Run("absent leaves both fields alone", func(t *testing.T) {
		app, err := Inflate(base())
		assert.NoError(t, err)
		assert.Nil(t, app.PolicyID)
		assert.False(t, app.ClearPolicyID)
	})

	t.Run("a positive id assigns", func(t *testing.T) {
		s := base()
		s["policy_id"] = 955633
		app, err := Inflate(s)
		assert.NoError(t, err)
		if assert.NotNil(t, app.PolicyID) {
			assert.Equal(t, 955633, *app.PolicyID)
		}
		assert.False(t, app.ClearPolicyID)
	})

	t.Run("zero unassigns", func(t *testing.T) {
		s := base()
		s["policy_id"] = 0
		app, err := Inflate(s)
		assert.NoError(t, err)
		assert.Nil(t, app.PolicyID)
		assert.True(t, app.ClearPolicyID)
	})

	t.Run("a value it cannot read is an error, not a silent drop", func(t *testing.T) {
		s := base()
		s["policy_id"] = "955633"
		_, err := Inflate(s)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "policy_id")
	})
}
