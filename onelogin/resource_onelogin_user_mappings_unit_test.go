package onelogin

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestUserMappingUpdateInput(t *testing.T) {
	t.Run("includes position when present in state", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, UserMappings().Schema, map[string]interface{}{
			"name":     "test",
			"match":    "all",
			"enabled":  true,
			"position": 5,
			"conditions": []interface{}{
				map[string]interface{}{
					"source":   "last_login",
					"operator": ">",
					"value":    "30",
				},
			},
			"actions": []interface{}{
				map[string]interface{}{
					"action": "set_status",
					"value":  []interface{}{"1"},
				},
			},
		})

		input := userMappingUpdateInput(d)
		_, hasPosition := input["position"]

		assert.True(t, hasPosition)
		assert.Equal(t, 5, input["position"])
		assert.NotContains(t, input, "id")
	})

	t.Run("omits position when not present", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, UserMappings().Schema, map[string]interface{}{
			"name":    "test",
			"match":   "all",
			"enabled": false,
			"conditions": []interface{}{
				map[string]interface{}{
					"source":   "last_login",
					"operator": ">",
					"value":    "30",
				},
			},
			"actions": []interface{}{
				map[string]interface{}{
					"action": "set_status",
					"value":  []interface{}{"1"},
				},
			},
		})

		input := userMappingUpdateInput(d)
		_, hasPosition := input["position"]

		assert.False(t, hasPosition)
		assert.NotContains(t, input, "id")
	})
}
