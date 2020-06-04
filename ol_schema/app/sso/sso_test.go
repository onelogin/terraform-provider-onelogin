package sso

import (
	"testing"

	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/apps"
	"github.com/stretchr/testify/assert"
)

func TestOIDCSSOSchema(t *testing.T) {
	t.Run("creates and returns a map of a OIDC SSO Schema", func(t *testing.T) {
		schema := OIDCSchema()
		assert.NotNil(t, schema["client_id"])
		assert.NotNil(t, schema["client_secret"])
	})
}

func TestSAMLSSOSchema(t *testing.T) {
	t.Run("creates and returns a map of a SAML SSO Schema", func(t *testing.T) {
		schema := SAMLSchema()
		assert.NotNil(t, schema["acs_url"])
		assert.NotNil(t, schema["metadata_url"])
		assert.NotNil(t, schema["issuer"])
		assert.NotNil(t, schema["certificate"])
	})
}

func TestFlattenOIDCSSO(t *testing.T) {
	tests := map[string]struct {
		InputData      apps.AppSso
		ExpectedOutput []map[string]interface{}
	}{
		"creates and returns the address of an AppSso struct for a OIDC app": {
			InputData: apps.AppSso{
				ClientID:     oltypes.String("test"),
				ClientSecret: oltypes.String("test"),
			},
			ExpectedOutput: []map[string]interface{}{
				map[string]interface{}{
					"client_id":     oltypes.String("test"),
					"client_secret": oltypes.String("test"),
				},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj := FlattenOIDC(test.InputData)
			assert.Equal(t, test.ExpectedOutput, subj)
		})
	}
}

func TestFlattenSAML(t *testing.T) {
	tests := map[string]struct {
		InputData      apps.AppSso
		ExpectedOutput []map[string]interface{}
	}{
		"creates and returns the address of an AppSso struct for a OIDC app": {
			InputData: apps.AppSso{
				MetadataURL: oltypes.String("test"),
				AcsURL:      oltypes.String("test"),
				SlsURL:      oltypes.String("test"),
				Issuer:      oltypes.String("test"),
				Certificate: &apps.AppSsoCertificate{
					Name:  oltypes.String("test"),
					ID:    oltypes.Int32(123),
					Value: oltypes.String("test"),
				},
			},
			ExpectedOutput: []map[string]interface{}{
				map[string]interface{}{
					"metadata_url": oltypes.String("test"),
					"acs_url":      oltypes.String("test"),
					"sls_url":      oltypes.String("test"),
					"issuer":       oltypes.String("test"),
					"certificate": []map[string]interface{}{
						map[string]interface{}{
							"name":  oltypes.String("test"),
							"id":    oltypes.Int32(123),
							"value": oltypes.String("test"),
						},
					},
				},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj := FlattenSAML(test.InputData)
			assert.Equal(t, test.ExpectedOutput, subj)
		})
	}
}
