package appprovisioningschema

import (
	"testing"

	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/apps"
	"github.com/stretchr/testify/assert"
)

func TestInflateProvisioning(t *testing.T) {
	tests := map[string]struct {
		ResourceData   map[string]interface{}
		ExpectedOutput apps.AppProvisioning
	}{
		"creates and returns the address of an AppProvisioning struct": {
			ResourceData:   map[string]interface{}{"enabled": true},
			ExpectedOutput: apps.AppProvisioning{Enabled: oltypes.Bool(true)},
		},
		"ignores unprovided field": {
			ResourceData:   map[string]interface{}{},
			ExpectedOutput: apps.AppProvisioning{},
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
		appProvisioning := apps.AppProvisioning{
			Enabled: oltypes.Bool(true),
		}
		subj := Flatten(appProvisioning)
		expected := map[string]interface{}{"enabled": true}
		assert.Equal(t, subj, expected)
	})
}
