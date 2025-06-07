package groupschema

import (
	"testing"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	"github.com/stretchr/testify/assert"
)

func TestSchema(t *testing.T) {
	t.Run("creates and returns a map of fields", func(t *testing.T) {
		schema := Schema()
		assert.NotNil(t, schema["id"])
		assert.NotNil(t, schema["name"])
		assert.NotNil(t, schema["reference"])
	})
}

func TestInflate(t *testing.T) {
	tests := map[string]struct {
		ResourceData map[string]interface{}
		Expected     models.Group
	}{
		"creates and returns the group struct": {
			ResourceData: map[string]interface{}{
				"id":        123,
				"name":      "test group",
				"reference": "test-ref",
			},
			Expected: models.Group{
				ID:        123,
				Name:      "test group",
				Reference: stringPtr("test-ref"),
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			group, _ := Inflate(test.ResourceData)
			assert.Equal(t, test.Expected.ID, group.ID)
			assert.Equal(t, test.Expected.Name, group.Name)
			assert.Equal(t, *test.Expected.Reference, *group.Reference)
		})
	}
}

func TestFlattenMany(t *testing.T) {
	t.Run("flattens group struct to map", func(t *testing.T) {
		groups := []models.Group{
			{
				ID:        123,
				Name:      "test group",
				Reference: stringPtr("test-ref"),
			},
			{
				ID:        456,
				Name:      "another group",
				Reference: nil,
			},
		}
		flattened := FlattenMany(groups)
		assert.Equal(t, 2, len(flattened))
		assert.Equal(t, 123, flattened[0]["id"])
		assert.Equal(t, "test group", flattened[0]["name"])
		assert.Equal(t, "test-ref", flattened[0]["reference"])
		assert.Equal(t, 456, flattened[1]["id"])
		assert.Equal(t, "another group", flattened[1]["name"])
		_, hasRef := flattened[1]["reference"]
		assert.False(t, hasRef)
	})
}

func TestFlatten(t *testing.T) {
	t.Run("flattens group struct to map", func(t *testing.T) {
		group := models.Group{
			ID:        123,
			Name:      "test group",
			Reference: stringPtr("test-ref"),
		}
		flattened := Flatten(group)
		assert.Equal(t, 123, flattened["id"])
		assert.Equal(t, "test group", flattened["name"])
		assert.Equal(t, "test-ref", flattened["reference"])
	})
}

func stringPtr(s string) *string {
	return &s
}
