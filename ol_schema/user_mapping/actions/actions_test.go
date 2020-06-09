package usermappingactionsschema

import (
	"testing"

	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/user_mappings"
	"github.com/stretchr/testify/assert"
)

func TestRulesSchema(t *testing.T) {
	t.Run("creates and returns a map of a user mapping Schema", func(t *testing.T) {
		provSchema := Schema()
		assert.NotNil(t, provSchema["action"])
		assert.NotNil(t, provSchema["value"])
	})
}

func TestInflate(t *testing.T) {
	tests := map[string]struct {
		ResourceData   map[string]interface{}
		ExpectedOutput usermappings.UserMappingActions
	}{
		"creates and returns the address of an user mappings action struct": {
			ResourceData: map[string]interface{}{
				"action": "test",
				"value":  []interface{}{"test"},
			},
			ExpectedOutput: usermappings.UserMappingActions{
				Action: oltypes.String("test"),
				Value:  []string{"test"},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj := Inflate(test.ResourceData)
			assert.Equal(t, subj, test.ExpectedOutput)
		})
	}
}

func TestFlatten(t *testing.T) {
	t.Run("It flattens the user mapping actions Struct", func(t *testing.T) {
		appConditionStruct := []usermappings.UserMappingActions{
			usermappings.UserMappingActions{
				Action: oltypes.String("test"),
				Value:  []string{"test"},
			},
			usermappings.UserMappingActions{
				Action: oltypes.String("test2"),
				Value:  []string{"test2"},
			},
		}
		subj := Flatten(appConditionStruct)
		expected := []map[string]interface{}{
			map[string]interface{}{
				"action": oltypes.String("test"),
				"value":  []string{"test"},
			},
			map[string]interface{}{
				"action": oltypes.String("test2"),
				"value":  []string{"test2"},
			},
		}
		assert.Equal(t, expected, subj)
	})
}
