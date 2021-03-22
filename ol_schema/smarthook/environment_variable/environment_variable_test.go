package smarthookenvironmentvariablesschema

import (
	"testing"

	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/smarthooks/envs"
	"github.com/stretchr/testify/assert"
)

func TestSmartHookSchema(t *testing.T) {
	t.Run("creates and returns a map of a Smarthooks Environment Variable Schema", func(t *testing.T) {
		provSchema := Schema()
		assert.NotNil(t, provSchema["name"])
		assert.NotNil(t, provSchema["value"])
		assert.NotNil(t, provSchema["created_at"])
		assert.NotNil(t, provSchema["updated_at"])
	})
}

func TestInflate(t *testing.T) {
	tests := map[string]struct {
		ResourceData   map[string]interface{}
		ExpectedOutput smarthookenvs.EnvVar
	}{
		"creates and returns the address of a SmartHook": {
			ResourceData: map[string]interface{}{
				"id":    "32f9dfee-a02c-4932-98ec-37838ce62ba0",
				"name":  "API_KEY",
				"value": "123-456-789",
			},
			ExpectedOutput: smarthookenvs.EnvVar{
				ID:    oltypes.String("32f9dfee-a02c-4932-98ec-37838ce62ba0"),
				Name:  oltypes.String("API_KEY"),
				Value: oltypes.String("123-456-789"),
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj := Inflate(test.ResourceData)
			assert.Equal(t, test.ExpectedOutput, subj)
		})
	}
}
