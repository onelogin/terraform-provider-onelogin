package appconfigurationschema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestInflateSAML covers what a SAML app sends, signature_algorithm or not.
func TestInflateSAML(t *testing.T) {
	t.Run("sends a configuration without signature_algorithm", func(t *testing.T) {
		got, err := InflateSAML(map[string]interface{}{
			"account_id": "my-domain",
			"audience":   "https://example.com/metadata",
		})
		assert.NoError(t, err)
		assert.Equal(t, map[string]interface{}{
			"account_id": "my-domain",
			"audience":   "https://example.com/metadata",
		}, got)
	})

	// Same as the signature_algorithm path: a value the map holds as ""
	// is not one the configuration set.
	t.Run("drops empty values", func(t *testing.T) {
		got, err := InflateSAML(map[string]interface{}{
			"audience":     "https://example.com/metadata",
			"consumer_url": "",
		})
		assert.NoError(t, err)
		assert.Equal(t, map[string]interface{}{"audience": "https://example.com/metadata"}, got)
	})

	t.Run("sends certificate_id as a number", func(t *testing.T) {
		got, err := InflateSAML(map[string]interface{}{"certificate_id": "123"})
		assert.NoError(t, err)
		assert.Equal(t, map[string]interface{}{"certificate_id": 123}, got)
	})

	t.Run("rejects a certificate_id that is not a number", func(t *testing.T) {
		_, err := InflateSAML(map[string]interface{}{"certificate_id": "abc"})
		assert.Error(t, err)
	})

	t.Run("an empty configuration sends an empty one", func(t *testing.T) {
		got, err := InflateSAML(map[string]interface{}{})
		assert.NoError(t, err)
		assert.Equal(t, map[string]interface{}{}, got)
	})

	// The guessing Inflate is unchanged, so this is still what any caller
	// that does not say its configuration is SAML gets.
	t.Run("Inflate still cannot place it", func(t *testing.T) {
		got, err := Inflate(map[string]interface{}{"audience": "https://example.com/metadata"})
		assert.NoError(t, err)
		assert.Equal(t, map[string]interface{}{}, got)
	})
}
