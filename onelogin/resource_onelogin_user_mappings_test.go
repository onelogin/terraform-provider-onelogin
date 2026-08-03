package onelogin

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccUserMapping_crud(t *testing.T) {
	base := GetFixture("onelogin_user_mapping_example.tf", t)
	update := GetFixture("onelogin_user_mapping_updated_example.tf", t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: base,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_user_mappings.basic_test", "name", "Select Login"),
					resource.TestCheckResourceAttr("onelogin_user_mappings.basic_test", "enabled", "true"),
					resource.TestCheckResourceAttr("onelogin_user_mappings.basic_test", "match", "all"),
					resource.TestCheckResourceAttr("onelogin_user_mappings.basic_test", "actions.0.action", "set_status"),
					resource.TestCheckResourceAttr("onelogin_user_mappings.basic_test", "conditions.0.value", "90"),
					resource.TestCheckResourceAttr("onelogin_user_mappings.basic_test", "conditions.0.source", "last_login"),
					resource.TestCheckResourceAttr("onelogin_user_mappings.basic_test", "conditions.0.operator", ">"),
				),
			},
			{
				Config: update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_user_mappings.basic_test", "name", "Updated Login"),
					resource.TestCheckResourceAttr("onelogin_user_mappings.basic_test", "enabled", "true"),
					resource.TestCheckResourceAttr("onelogin_user_mappings.basic_test", "match", "all"),
					resource.TestCheckResourceAttr("onelogin_user_mappings.basic_test", "actions.0.action", "set_status"),
					resource.TestCheckResourceAttr("onelogin_user_mappings.basic_test", "conditions.0.value", "120"),
					resource.TestCheckResourceAttr("onelogin_user_mappings.basic_test", "conditions.0.source", "last_login"),
					resource.TestCheckResourceAttr("onelogin_user_mappings.basic_test", "conditions.0.operator", ">"),
				),
			},
		},
	})
}

// TestUserMappingUpdateInput covers the update body assembled by the provider.
// The API rejects an update whose body carries the mapping ID, and re-enabling a
// mapping with the position it held before it was disabled fails the same way.
func TestUserMappingUpdateInput(t *testing.T) {
	baseState := func(extra map[string]interface{}) map[string]interface{} {
		state := map[string]interface{}{
			"name":    "Select Login",
			"match":   "all",
			"enabled": true,
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
		}
		for k, v := range extra {
			state[k] = v
		}
		return state
	}

	t.Run("never sends the id in the body", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, UserMappings().Schema, baseState(nil))

		if _, present := userMappingUpdateInput(d)["id"]; present {
			t.Fatal("expected the update input to omit id, the API takes it from the URL")
		}
	})

	t.Run("includes an explicit position", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, UserMappings().Schema, baseState(map[string]interface{}{
			"position": 5,
		}))

		input := userMappingUpdateInput(d)
		if input["position"] != 5 {
			t.Fatalf("expected position 5, got %v", input["position"])
		}
	})

	t.Run("omits a position cleared by a previous disable", func(t *testing.T) {
		// userMappingRead clears position when the API reports none, which is what
		// a disabled mapping looks like. Re-enabling must not resurrect it.
		d := schema.TestResourceDataRaw(t, UserMappings().Schema, baseState(map[string]interface{}{
			"position": 5,
		}))
		if err := d.Set("position", nil); err != nil {
			t.Fatalf("expected to clear position, got %v", err)
		}

		if value, present := userMappingUpdateInput(d)["position"]; present {
			t.Fatalf("expected a cleared position to be omitted, got %v", value)
		}
	})

	t.Run("omits the position while disabled", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, UserMappings().Schema, baseState(map[string]interface{}{
			"enabled": false,
		}))

		if value, present := userMappingUpdateInput(d)["position"]; present {
			t.Fatalf("expected no position for a disabled mapping, got %v", value)
		}
	})

	t.Run("carries the mapping fields through", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, UserMappings().Schema, baseState(nil))
		input := userMappingUpdateInput(d)

		for _, field := range []string{"name", "match", "enabled", "conditions", "actions"} {
			if _, present := input[field]; !present {
				t.Fatalf("expected %q in the update input", field)
			}
		}
		if input["name"] != "Select Login" {
			t.Fatalf("expected name to be carried through, got %v", input["name"])
		}
	})
}
