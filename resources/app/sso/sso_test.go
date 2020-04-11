package sso

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOIDCSSOSchema(t *testing.T) {
	t.Run("creates and returns a map of a OIDC SSO Schema", func(t *testing.T) {
		schema := OIDCSSOSchema()
		assert.NotNil(t, schema["client_id"])
		assert.NotNil(t, schema["client_secret"])
	})
}

func TestSAMLSSOSchema(t *testing.T) {
	t.Run("creates and returns a map of a SAML SSO Schema", func(t *testing.T) {
		schema := SAMLSSOSchema()
		assert.NotNil(t, schema["acs_url"])
		assert.NotNil(t, schema["metadata_url"])
		assert.NotNil(t, schema["issuer"])
		assert.NotNil(t, schema["certificate"])
	})
}
