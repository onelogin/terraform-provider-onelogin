package onelogin

import (
	"strings"
	"testing"
)

// TestStripProviderConfig covers what the fixture loader removes, and more
// importantly what it must not.
func TestStripProviderConfig(t *testing.T) {
	t.Run("removes the blocks an acceptance test supplies itself", func(t *testing.T) {
		out := stripProviderConfig(`terraform {
  required_providers {
    onelogin = {
      source  = "onelogin/onelogin"
      version = ">= 0.8.0"
    }
  }
}

provider "onelogin" {
  # comment
}

resource onelogin_roles role {
  name = "kept"
}
`)

		for _, gone := range []string{"required_providers", "onelogin/onelogin", `provider "onelogin"`} {
			if strings.Contains(out, gone) {
				t.Fatalf("expected %q to be stripped, got:\n%s", gone, out)
			}
		}
		if !strings.Contains(out, `name = "kept"`) {
			t.Fatalf("expected the resource to survive, got:\n%s", out)
		}
	})

	t.Run("keeps the provider meta-argument inside a resource", func(t *testing.T) {
		// `provider = onelogin.eu` selects an aliased provider for one
		// resource. Matching on "provider " rather than `provider "` dropped
		// this line, quietly moving the resource to the default provider.
		out := stripProviderConfig(`resource onelogin_roles role {
  provider = onelogin.eu
  name     = "kept"
}
`)

		if !strings.Contains(out, "provider = onelogin.eu") {
			t.Fatalf("expected the provider meta-argument to survive, got:\n%s", out)
		}
		if !strings.Contains(out, `name     = "kept"`) {
			t.Fatalf("expected the rest of the resource to survive, got:\n%s", out)
		}
	})
}
