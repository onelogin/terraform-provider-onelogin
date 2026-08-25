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

func TestInflatePolicyID(t *testing.T) {
	t.Run("absent leaves the policy alone", func(t *testing.T) {
		group, _ := Inflate(map[string]interface{}{"name": "Engineering"})
		assert.Nil(t, group.PolicyID, "an unmentioned policy must stay out of the request body")
	})

	t.Run("a zero is carried through", func(t *testing.T) {
		// The clear. A nil here would drop the key and silently preserve the
		// existing assignment, so the zero has to survive Inflate.
		group, _ := Inflate(map[string]interface{}{"name": "Engineering", "policy_id": 0})
		assert.NotNil(t, group.PolicyID)
		assert.Equal(t, 0, *group.PolicyID)
	})

	t.Run("an id is carried through", func(t *testing.T) {
		group, _ := Inflate(map[string]interface{}{"name": "Engineering", "policy_id": 955633})
		assert.NotNil(t, group.PolicyID)
		assert.Equal(t, 955633, *group.PolicyID)
	})
}

func TestPolicyIDSchema(t *testing.T) {
	// Computed here would read an empty configuration as "keep what is in
	// state", leaving no way to remove a policy once set.
	policy := Schema()["policy_id"]
	assert.NotNil(t, policy)
	assert.True(t, policy.Optional)
	assert.False(t, policy.Computed, "Computed would make the assignment impossible to clear")
}
