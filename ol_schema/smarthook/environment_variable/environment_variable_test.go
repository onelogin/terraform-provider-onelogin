package smarthookenvironmentvariablesschema

import (
	"testing"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
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
	// Create test variables
	id := "32f9dfee-a02c-4932-98ec-37838ce62ba0"
	name := "API_KEY"
	value := "123-456-789"

	tests := map[string]struct {
		ResourceData   map[string]interface{}
		ExpectedOutput models.EnvVar
	}{
		"creates and returns the address of a SmartHook": {
			ResourceData: map[string]interface{}{
				"id":    id,
				"name":  name,
				"value": value,
			},
			ExpectedOutput: models.EnvVar{
				ID:    &id,
				Name:  &name,
				Value: &value,
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj := Inflate(test.ResourceData)
			if subj.ID != nil && test.ExpectedOutput.ID != nil {
				assert.Equal(t, *subj.ID, *test.ExpectedOutput.ID)
			}
			if subj.Name != nil && test.ExpectedOutput.Name != nil {
				assert.Equal(t, *subj.Name, *test.ExpectedOutput.Name)
			}
			if subj.Value != nil && test.ExpectedOutput.Value != nil {
				assert.Equal(t, *subj.Value, *test.ExpectedOutput.Value)
			}
		})
	}
}
