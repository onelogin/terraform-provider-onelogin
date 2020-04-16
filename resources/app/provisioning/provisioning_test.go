package provisioning

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
	tests := map[string]struct {
		ResourceData   map[string]interface{}
		ExpectedOutput *models.AppProvisioning
	}{
		"creates and returns the address of an AppProvisioning struct": {
			ResourceData:   map[string]interface{}{"enabled": true},
			ExpectedOutput: &models.AppProvisioning{Enabled: oltypes.Bool(true)},
		},
		"ignores unprovided field": {
			ResourceData:   map[string]interface{}{},
			ExpectedOutput: &models.AppProvisioning{},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			prov := InflateProvisioning(&test.ResourceData)
			assert.Equal(t, prov, test.ExpectedOutput)
		})
	}
}

func TestFlatten(t *testing.T){
	t.Run("It flattens the AppProvisioning Struct", func(t *testing.T){
		appProvisioning := models.AppProvisioning{
			Enabled: oltypes.Bool(true),
		}
		subj := Flatten(&appProvisioning)
		expected := []map[string]interface{}{ {"enabled": true} }
		assert.Equal(t, subj, expected)
	})
}
