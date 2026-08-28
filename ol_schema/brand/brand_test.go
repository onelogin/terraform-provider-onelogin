package brand

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestSchema(t *testing.T) {
	s := Schema()

	t.Run("name is required", func(t *testing.T) {
		// The API answers a create with no name with
		// 422 {"name":"Value is required."}.
		assert.True(t, s["name"].Required)
	})

	t.Run("master is read-only", func(t *testing.T) {
		assert.True(t, s["master"].Computed)
		assert.False(t, s["master"].Optional, "a brand cannot be made the master through this API")
	})

	t.Run("the images are optional, not computed, and not printed", func(t *testing.T) {
		for _, name := range []string{"logo", "background"} {
			assert.True(t, s[name].Optional, "%s should be settable", name)
			assert.False(t, s[name].Computed,
				"%s must not be Computed: the API returns an object of URLs, not the base64 that was sent, so there is nothing to write back", name)
			// Not because an image is a secret, but because Terraform prints a
			// string attribute in full: a 200KB logo measured at a single
			// 270,909-character line and a 273KB plan, against 2.2KB with this
			// set. The API allows five times that for a background.
			assert.True(t, s[name].Sensitive,
				"%s holds base64 image data and would otherwise be printed in full on every plan", name)
		}
	})

	t.Run("every plain field is optional and computed", func(t *testing.T) {
		for _, f := range fields {
			attribute, ok := s[f.Name]
			if assert.True(t, ok, "%s is missing from the schema", f.Name) {
				assert.True(t, attribute.Optional, "%s should be optional", f.Name)
				assert.True(t, attribute.Computed, "%s should be computed, or the API's own default shows as a diff", f.Name)
				assert.Equal(t, f.Type, attribute.Type, "%s has the wrong type", f.Name)
			}
		}
	})

	t.Run("only the localised fields suppress JSON differences", func(t *testing.T) {
		for _, f := range fields {
			if f.Localised {
				assert.NotNil(t, s[f.Name].DiffSuppressFunc, "%s is localised and should suppress equivalent JSON", f.Name)
			} else {
				assert.Nil(t, s[f.Name].DiffSuppressFunc, "%s is not localised and should compare literally", f.Name)
			}
		}
	})

	t.Run("every field has a description", func(t *testing.T) {
		for name, attribute := range s {
			assert.NotEmpty(t, attribute.Description, "%s has no description", name)
		}
	})
}

func TestSuppressEquivalentJSON(t *testing.T) {
	cases := []struct {
		name       string
		old, new   string
		suppressed bool
	}{
		{
			name:       "identical strings",
			old:        `{"en":"Username"}`,
			new:        `{"en":"Username"}`,
			suppressed: true,
		},
		{
			// The case this exists for: the API re-encodes what it stores, so
			// the string coming back is not the one that went in.
			name:       "same object, different whitespace",
			old:        `{"en":"Username"}`,
			new:        "{\n  \"en\": \"Username\"\n}",
			suppressed: true,
		},
		{
			name:       "same object, different key order",
			old:        `{"en":"Username","de":"Benutzername"}`,
			new:        `{"de":"Benutzername","en":"Username"}`,
			suppressed: true,
		},
		{
			name:       "a changed value is a real diff",
			old:        `{"en":"Username"}`,
			new:        `{"en":"Email"}`,
			suppressed: false,
		},
		{
			name:       "an added locale is a real diff",
			old:        `{"en":"Username"}`,
			new:        `{"en":"Username","de":"Benutzername"}`,
			suppressed: false,
		},
		{
			// Malformed input must not be quietly treated as equal to
			// anything: the practitioner needs to see the diff.
			name:       "unparseable values are compared literally",
			old:        `not json`,
			new:        `also not json`,
			suppressed: false,
		},
		{
			name:       "one side unparseable",
			old:        `{"en":"Username"}`,
			new:        `{"en":`,
			suppressed: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.suppressed, suppressEquivalentJSON("", c.old, c.new, nil))
		})
	}
}

func TestConfiguredKeys(t *testing.T) {
	t.Run("returns the attributes that are not null", func(t *testing.T) {
		keys := ConfiguredKeys(cty.ObjectVal(map[string]cty.Value{
			"name":         cty.StringVal("Engineering"),
			"enabled":      cty.False,
			"custom_color": cty.NullVal(cty.String),
		}))

		assert.True(t, keys["name"])
		// A configured false has to survive: it is the difference between
		// "leave it alone" and "turn it off".
		assert.True(t, keys["enabled"], "a configured false is still configured")
		assert.False(t, keys["custom_color"], "an unwritten attribute is not configured")
	})

	t.Run("returns nil when there is no configuration to read", func(t *testing.T) {
		assert.Nil(t, ConfiguredKeys(cty.NullVal(cty.EmptyObject)))
	})
}

// getter is a Getter built from a plain map, so RequestBody can be exercised
// without a live ResourceData.
type getter map[string]interface{}

func (g getter) Get(key string) interface{} { return g[key] }

func (g getter) GetOk(key string) (interface{}, bool) {
	value, ok := g[key]
	if !ok {
		return nil, false
	}
	switch v := value.(type) {
	case string:
		return value, v != ""
	case bool:
		return value, v
	case int:
		return value, v != 0
	}
	return value, true
}

func TestRequestBody(t *testing.T) {
	t.Run("always sends the name", func(t *testing.T) {
		// models.Brand.Name is the one field without omitempty, so a nil Name
		// marshals to "name": null -- which the API answers with 422 "Value
		// must be a string." even though omitting the key is a 200.
		body := RequestBody(getter{"name": "Engineering"}, map[string]bool{"name": true})

		if assert.NotNil(t, body.Name, "a nil Name would be sent as an explicit null") {
			assert.Equal(t, "Engineering", *body.Name)
		}
	})

	t.Run("sends only what was configured", func(t *testing.T) {
		body := RequestBody(
			getter{"name": "Engineering", "custom_color": "#123456", "enabled": true},
			map[string]bool{"name": true, "custom_color": true},
		)

		if assert.NotNil(t, body.CustomColor) {
			assert.Equal(t, "#123456", *body.CustomColor)
		}
		assert.Nil(t, body.Enabled, "enabled was not configured, so it must not be sent")
	})

	t.Run("sends a configured false rather than dropping it", func(t *testing.T) {
		body := RequestBody(
			getter{"name": "Engineering", "enabled": false},
			map[string]bool{"name": true, "enabled": true},
		)

		if assert.NotNil(t, body.Enabled, "a configured false is an instruction to disable the brand") {
			assert.False(t, *body.Enabled)
		}
	})

	t.Run("sends a configured zero opacity", func(t *testing.T) {
		body := RequestBody(
			getter{"name": "Engineering", "custom_masking_opacity": 0},
			map[string]bool{"name": true, "custom_masking_opacity": true},
		)

		if assert.NotNil(t, body.CustomMaskingOpacity) {
			assert.Equal(t, int32(0), *body.CustomMaskingOpacity)
		}
	})

	t.Run("sends the images when configured", func(t *testing.T) {
		body := RequestBody(
			getter{"name": "Engineering", "logo": "aGVsbG8=", "background": "d29ybGQ="},
			map[string]bool{"name": true, "logo": true, "background": true},
		)

		if assert.NotNil(t, body.Logo) {
			assert.Equal(t, "aGVsbG8=", *body.Logo)
		}
		if assert.NotNil(t, body.Background) {
			assert.Equal(t, "d29ybGQ=", *body.Background)
		}
	})

	t.Run("falls back to GetOk when there is no raw configuration", func(t *testing.T) {
		body := RequestBody(getter{"name": "Engineering", "custom_color": "#123456"}, nil)

		if assert.NotNil(t, body.CustomColor) {
			assert.Equal(t, "#123456", *body.CustomColor)
		}
	})
}

// TestRequestBodyCoversEveryField guards the assign switch. A field added to
// the table but not to the switch would be in the schema, accepted in a
// configuration, and then silently never sent -- the resource would report
// success having changed nothing.
func TestRequestBodyCoversEveryField(t *testing.T) {
	values := getter{"name": "Engineering"}

	for _, f := range fields {
		switch f.Type {
		case schema.TypeBool:
			values[f.Name] = true
		case schema.TypeInt:
			values[f.Name] = 42
		default:
			values[f.Name] = "set"
		}
	}

	// Every attribute in the table must reach some field of the model. A body
	// carrying nothing but the name is the baseline; configuring one more
	// attribute has to change it.
	nameOnly := RequestBody(getter{"name": "Engineering"}, map[string]bool{"name": true})

	for _, f := range fields {
		withField := RequestBody(values, map[string]bool{"name": true, f.Name: true})

		assert.NotEqual(t, nameOnly, withField,
			"%s is in the schema but assign() never writes it to models.Brand, so configuring it would silently do nothing", f.Name)
	}
}
