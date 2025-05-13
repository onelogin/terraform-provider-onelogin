package appssoschema

import (
	"testing"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	"github.com/stretchr/testify/assert"
)

func TestFlattenOIDCSSO(t *testing.T) {
	tests := map[string]struct {
		InputData      models.SSOOpenId
		ExpectedOutput map[string]interface{}
	}{
		"creates and returns a map of SSO fields from an OIDC app": {
			InputData: models.SSOOpenId{
				ClientID: "test",
			},
			ExpectedOutput: map[string]interface{}{
				"client_id": "test",
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

func TestFlattenSAMLCert(t *testing.T) {
	tests := map[string]struct {
		InputData      models.SSOSAML
		ExpectedOutput map[string]interface{}
	}{
		"creates and returns a map of SAML SSO Certificate fields for the given SAML app": {
			InputData: models.SSOSAML{
				MetadataURL: "test",
				AcsURL:      "test",
				SlsURL:      "test",
				Issuer:      "test",
				Certificate: models.Certificate{
					Name:  "test",
					ID:    123,
					Value: "test",
				},
			},
			ExpectedOutput: map[string]interface{}{
				"name":  "test",
				"value": "test",
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj := FlattenSAMLCert(test.InputData)
			assert.Equal(t, test.ExpectedOutput, subj)
		})
	}
}

func TestFlattenSAML(t *testing.T) {
	tests := map[string]struct {
		InputData      models.SSOSAML
		ExpectedOutput map[string]interface{}
	}{
		"creates and returns a map of SSO fields for a SAML app": {
			InputData: models.SSOSAML{
				MetadataURL: "test",
				AcsURL:      "test",
				SlsURL:      "test",
				Issuer:      "test",
				Certificate: models.Certificate{
					Name:  "test",
					ID:    123,
					Value: "test",
				},
			},
			ExpectedOutput: map[string]interface{}{
				"metadata_url": "test",
				"acs_url":      "test",
				"sls_url":      "test",
				"issuer":       "test",
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
