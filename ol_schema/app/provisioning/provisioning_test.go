package appprovisioningschema

import (
	"testing"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	"github.com/stretchr/testify/assert"
)

func TestInflateProvisioning(t *testing.T) {
	tests := map[string]struct {
		ResourceData   map[string]interface{}
		ExpectedOutput models.Provisioning
	}{
		"creates and returns the address of an AppProvisioning struct": {
			ResourceData:   map[string]interface{}{"enabled": true},
			ExpectedOutput: models.Provisioning{Enabled: true},
		},
		"ignores unprovided field": {
			ResourceData:   map[string]interface{}{},
			ExpectedOutput: models.Provisioning{},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			prov := Inflate(test.ResourceData)
			assert.Equal(t, prov, test.ExpectedOutput)
		})
	}
}

func TestFlatten(t *testing.T) {
	t.Run("It flattens the AppProvisioning Struct", func(t *testing.T) {
		appProvisioning := models.Provisioning{
			Enabled: true,
		}
		subj := Flatten(appProvisioning)
		expected := map[string]interface{}{"enabled": true}
		assert.Equal(t, subj, expected)
	})
}
