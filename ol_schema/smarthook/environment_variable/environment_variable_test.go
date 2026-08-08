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

// TestInflateUpdateBody pins the shape the update endpoint accepts. It takes
// the variable's ID from the URL and rejects a body carrying anything but the
// value, with `instance is not allowed to have the additional property "id"`,
// so an update that inflates from the resource ID fails every time.
func TestInflateUpdateBody(t *testing.T) {
	out := Inflate(map[string]interface{}{"value": "987-654-321"})

	if out.ID != nil {
		t.Fatalf("expected no id in the update body, got %q", *out.ID)
	}
	if out.Name != nil {
		t.Fatalf("expected no name in the update body, got %q", *out.Name)
	}
	if out.Value == nil || *out.Value != "987-654-321" {
		t.Fatalf("expected the value to be carried through, got %v", out.Value)
	}
}

// TestNameForcesNew records that the API cannot rename a variable, so a changed
// name has to replace it rather than update it.
func TestNameForcesNew(t *testing.T) {
	if name := Schema()["name"]; !name.ForceNew {
		t.Fatal("expected name to be ForceNew: the update endpoint refuses a name outright")
	}
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
