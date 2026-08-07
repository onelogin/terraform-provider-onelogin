package onelogin

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestUserCustomAttributesNameIsOptional guards the schema shape. A value set on
// one user needs no name, so requiring it made the documented usage -- and the
// user_attr_value acceptance test below -- impossible to plan.
func TestUserCustomAttributesNameIsOptional(t *testing.T) {
	name := UserCustomAttributes().Schema["name"]

	if name.Required {
		t.Fatal("expected name to be optional, a user-value resource has no name to give")
	}
	if !name.Optional {
		t.Fatal("expected name to be optional")
	}
}

// TestCheckCustomAttributeDefinitionName covers the apply-time check that stands
// in for the Required the schema can no longer carry.
func TestCheckCustomAttributeDefinitionName(t *testing.T) {
	t.Run("accepts a definition with a name", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, UserCustomAttributes().Schema, map[string]interface{}{
			"name":      "Employee ID",
			"shortname": "employee_id",
		})

		if err := checkCustomAttributeDefinitionName(d); err != nil {
			t.Fatalf("expected the definition to be accepted, got %v", err)
		}
	})

	t.Run("rejects a definition without a name", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, UserCustomAttributes().Schema, map[string]interface{}{
			"shortname": "employee_id",
		})

		err := checkCustomAttributeDefinitionName(d)
		if err == nil {
			t.Fatal("expected a definition with no name to be rejected")
		}
		// The shortname is the only thing identifying which resource is at
		// fault, so it has to appear in the message.
		if !strings.Contains(err.Error(), "employee_id") {
			t.Fatalf("expected the shortname in the error, got %v", err)
		}
	})
}

// TestValidCustomAttributePosition covers the plan-time guard on position. A
// negative one reads differently in each payload builder -- sent as-is on
// create, treated as a clear on update -- so it has to be rejected up front.
func TestValidCustomAttributePosition(t *testing.T) {
	for _, position := range []int{0, 1, 27} {
		if _, errs := validCustomAttributePosition(position, "position"); len(errs) > 0 {
			t.Fatalf("expected position %d to be accepted, got %v", position, errs)
		}
	}

	for _, position := range []int{-1, -27} {
		if _, errs := validCustomAttributePosition(position, "position"); len(errs) == 0 {
			t.Fatalf("expected position %d to be rejected", position)
		}
	}
}

// TestUserCustomAttributeDefinitionCreateInput covers the "user_field" body sent
// when a definition is created. OneLogin defaults position to null, so an unset
// position has to stay out of the body rather than go in as a zero.
func TestUserCustomAttributeDefinitionCreateInput(t *testing.T) {
	definition := func(extra map[string]interface{}) map[string]interface{} {
		state := map[string]interface{}{
			"name":      "Employee ID",
			"shortname": "employee_id",
		}
		for k, v := range extra {
			state[k] = v
		}
		return state
	}

	t.Run("always carries name and shortname", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, UserCustomAttributes().Schema, definition(nil))

		input := userCustomAttributeDefinitionCreateInput(d)
		if input["name"] != "Employee ID" {
			t.Fatalf("expected name %q, got %v", "Employee ID", input["name"])
		}
		if input["shortname"] != "employee_id" {
			t.Fatalf("expected shortname %q, got %v", "employee_id", input["shortname"])
		}
	})

	t.Run("includes an explicit position", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, UserCustomAttributes().Schema, definition(map[string]interface{}{
			"position": 27,
		}))

		if input := userCustomAttributeDefinitionCreateInput(d); input["position"] != 27 {
			t.Fatalf("expected position 27, got %v", input["position"])
		}
	})

	t.Run("omits an unset position", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, UserCustomAttributes().Schema, definition(nil))

		if value, present := userCustomAttributeDefinitionCreateInput(d)["position"]; present {
			t.Fatalf("expected an unset position to be omitted, got %v", value)
		}
	})
}

// TestUserCustomAttributeDefinitionUpdateInput covers the body sent when a
// definition is updated. The API requires name and shortname on every update,
// and only a literal null resets a position -- an omitted key leaves it alone.
func TestUserCustomAttributeDefinitionUpdateInput(t *testing.T) {
	t.Run("always carries the required fields", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, UserCustomAttributes().Schema, map[string]interface{}{
			"name":      "Employee ID",
			"shortname": "employee_id",
		})

		input := userCustomAttributeDefinitionUpdateInput(d)
		if input["name"] != "Employee ID" {
			t.Fatalf("expected name %q, got %v", "Employee ID", input["name"])
		}
		if input["shortname"] != "employee_id" {
			t.Fatalf("expected shortname %q, got %v", "employee_id", input["shortname"])
		}
	})

	t.Run("sends a configured position", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, UserCustomAttributes().Schema, map[string]interface{}{
			"name":      "Employee ID",
			"shortname": "employee_id",
			"position":  27,
		})

		if input := userCustomAttributeDefinitionUpdateInput(d); input["position"] != 27 {
			t.Fatalf("expected position 27, got %v", input["position"])
		}
	})

	t.Run("sends an explicit null once the position is cleared", func(t *testing.T) {
		// A position dropped from the configuration plans as 0, which is the
		// state this update has to turn back into a null.
		d := schema.TestResourceDataRaw(t, UserCustomAttributes().Schema, map[string]interface{}{
			"name":      "Employee ID",
			"shortname": "employee_id",
			"position":  0,
		})

		input := userCustomAttributeDefinitionUpdateInput(d)
		value, present := input["position"]
		if !present {
			t.Fatal("expected a cleared position to be sent, an omitted key leaves the old value in place")
		}
		if value != nil {
			t.Fatalf("expected a cleared position to be sent as null, got %v", value)
		}
	})
}

func TestAccUserCustomAttributes_basic(t *testing.T) {
	// Inline rather than a fixture, so the token is substituted here. Custom
	// attribute shortnames are unique per tenant, and a definition left behind
	// by an earlier run would otherwise collide.
	suffix := acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckUserCustomAttributesConfig(suffix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_user_custom_attributes.test_attr", "name", "Test Attribute"),
					resource.TestCheckResourceAttr("onelogin_user_custom_attributes.test_attr", "shortname", "test_attr_"+suffix),
				),
			},
		},
	})
}

func testAccCheckUserCustomAttributesConfig(suffix string) string {
	return strings.ReplaceAll(`
resource onelogin_user_custom_attributes test_attr {
  name      = "Test Attribute"
  shortname = "test_attr_acctest"
}
`, fixtureUniqueToken, suffix)
}

// TestAccUserCustomAttributes_position walks a definition through the lifecycle
// reported in issue #224: created without a position, given one, then put back
// to having none.
func TestAccUserCustomAttributes_position(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckUserCustomAttributesPositionConfig(""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_user_custom_attributes.positioned_attr", "position", "0"),
				),
			},
			{
				Config: testAccCheckUserCustomAttributesPositionConfig("position = 27"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_user_custom_attributes.positioned_attr", "position", "27"),
				),
			},
			{
				ResourceName:      "onelogin_user_custom_attributes.positioned_attr",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCheckUserCustomAttributesPositionConfig(""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_user_custom_attributes.positioned_attr", "position", "0"),
				),
			},
		},
	})
}

func testAccCheckUserCustomAttributesPositionConfig(position string) string {
	return `
resource onelogin_user_custom_attributes positioned_attr {
  name      = "Positioned Attribute"
  shortname = "positioned_attr"
  ` + position + `
}
`
}

func TestAccUserCustomAttributesWithUser_basic(t *testing.T) {
	// This config is inline rather than a fixture, so it needs the token
	// substituted here: OneLogin enforces unique usernames per tenant.
	suffix := acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckUserCustomAttributesWithUserConfig(suffix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("onelogin_users.test_user", "username"),
					resource.TestCheckResourceAttr("onelogin_user_custom_attributes.test_attr", "name", "Test Attribute"),
					resource.TestCheckResourceAttr("onelogin_user_custom_attributes.test_attr", "shortname", "test_attr_"+suffix),
					resource.TestCheckResourceAttrSet("onelogin_user_custom_attributes.user_attr_value", "user_id"),
					resource.TestCheckResourceAttr("onelogin_user_custom_attributes.user_attr_value", "value", "test_value"),
				),
			},
		},
	})
}

func testAccCheckUserCustomAttributesWithUserConfig(suffix string) string {
	return strings.ReplaceAll(`
resource onelogin_users test_user {
  username = "test.user.for.attrs.acctest"
  email    = "test.user.attrs.acctest@example.com"
}

resource onelogin_user_custom_attributes test_attr {
  name      = "Test Attribute"
  shortname = "test_attr_acctest"
}

resource onelogin_user_custom_attributes user_attr_value {
  user_id   = onelogin_users.test_user.id
  shortname = onelogin_user_custom_attributes.test_attr.shortname
  value     = "test_value"
}
`, fixtureUniqueToken, suffix)
}
