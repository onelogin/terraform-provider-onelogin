package app

import (
	"testing"

	"github.com/onelogin/onelogin-go-sdk/pkg/models"
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/stretchr/testify/assert"
)

func TestProvisioningSchema(t *testing.T) {
	t.Run("creates and returns a map of an AppProvisioning Schema", func(t *testing.T) {
		provSchema := ProvisioningSchema()
		assert.NotNil(t, provSchema["enabled"])
	})
}

func TestInflateProvisioning(t *testing.T) {
	test := struct {
		ResourceData   map[string]interface{}
		ExpectedOutput *models.AppProvisioning
	}{

		ResourceData:   map[string]interface{}{"enabled": true},
		ExpectedOutput: &models.AppProvisioning{Enabled: oltypes.Bool(true)},
	}
	t.Run("creates and returns the address of an AppProvisioning struct", func(t *testing.T) {
		prov := InflateProvisioning(&test.ResourceData)
		assert.Equal(t, prov, test.ExpectedOutput)
		assert.Equal(t, *prov.Enabled, true)
	})

}
