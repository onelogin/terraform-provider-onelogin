package onelogin

import (
	"testing"

	mod "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	"github.com/stretchr/testify/assert"
)

func TestBuildUserMappingUpdatePayload(t *testing.T) {
	name := "mapping"
	match := "all"
	position := int32(9)
	enabled := true
	disabled := false
	id := int32(100)

	t.Run("sends position null when disabling", func(t *testing.T) {
		payload := buildUserMappingUpdatePayload(mod.UserMapping{
			ID:       &id,
			Name:     &name,
			Match:    &match,
			Enabled:  &disabled,
			Position: &position,
		})

		val, hasPosition := payload["position"]
		_, hasID := payload["id"]

		assert.True(t, hasPosition)
		assert.Nil(t, val)
		assert.False(t, hasID)
	})

	t.Run("keeps explicit position when enabled", func(t *testing.T) {
		payload := buildUserMappingUpdatePayload(mod.UserMapping{
			Name:     &name,
			Match:    &match,
			Enabled:  &enabled,
			Position: &position,
		})

		val, hasPosition := payload["position"]
		assert.True(t, hasPosition)
		assert.Equal(t, &position, val)
	})

	t.Run("omits position when not provided", func(t *testing.T) {
		payload := buildUserMappingUpdatePayload(mod.UserMapping{
			Name:    &name,
			Match:   &match,
			Enabled: &enabled,
		})

		_, hasPosition := payload["position"]
		assert.False(t, hasPosition)
	})
}
