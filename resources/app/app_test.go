package app

import (
	"testing"

	"github.com/onelogin/onelogin-terraform-provider/resources/app/configuration"
	"github.com/stretchr/testify/assert"
)

func TestAddSubSchema(t *testing.T) {
	t.Run("adds sub schema to given resrouce schema", func(t *testing.T) {
		appSchema := AppSchema()
		AddSubSchema("sub", &appSchema, configuration.SAMLConfigurationSchema)
		assert.NotNil(t, appSchema["sub"])
	})
}
