package appssoschema

import (
	"testing"

	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/apps"
	"github.com/stretchr/testify/assert"
)

func TestFlattenOIDCSSO(t *testing.T) {
	tests := map[string]struct {
		InputData      apps.AppSso
		ExpectedOutput map[string]interface{}
	}{
		"creates and returns the address of an AppSso struct for a OIDC app": {
			InputData: apps.AppSso{
				ClientID:     oltypes.String("test"),
				ClientSecret: oltypes.String("test"),
			},
			ExpectedOutput: map[string]interface{}{
				"client_id":     oltypes.String("test"),
				"client_secret": oltypes.String("test"),
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
		ExpectedOutput map[string]interface{}
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
			ExpectedOutput: map[string]interface{}{
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
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj := FlattenSAML(test.InputData)
			assert.Equal(t, test.ExpectedOutput, subj)
		})
	}
}
