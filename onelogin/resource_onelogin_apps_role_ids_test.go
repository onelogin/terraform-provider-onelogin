package onelogin

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	"github.com/onelogin/terraform-provider-onelogin/utils"
)

// roleIDsState builds the state a read leaves behind for an app with roles.
func roleIDsState(t *testing.T, r *schema.Resource, roleIDs ...int) *terraform.InstanceState {
	t.Helper()

	d := r.Data(nil)
	d.SetId("1543259")
	for key, value := range map[string]interface{}{
		"name":         "my SAML APP",
		"connector_id": 110016,
		"role_ids":     roleIDs,
	} {
		if err := d.Set(key, value); err != nil {
			t.Fatal(err)
		}
	}
	return d.State()
}

// TestAppRoleIDsConfigurable is gh issue 272: role_ids has to be settable for
// an app's roles to be declared in one write, and stay computed so that roles
// assigned some other way survive a configuration that does not mention them.
func TestAppRoleIDsConfigurable(t *testing.T) {
	for name, newResource := range appResourcesUnderTest() {
		t.Run(name, func(t *testing.T) {
			attribute, ok := newResource().CoreConfigSchema().Attributes["role_ids"]
			if !ok {
				t.Fatal("role_ids is missing from the core schema")
			}
			if !attribute.Optional {
				t.Error("role_ids is not Optional, so a configuration cannot set it")
			}
			if !attribute.Computed {
				t.Error("role_ids is not Computed, so roles assigned outside the app resource would be cleared")
			}
		})
	}
}

// TestAppRoleIDsLeftAloneWhenUnset covers everyone already assigning roles
// with onelogin_app_role_attachments, onelogin_roles.apps or the OneLogin UI.
// Their state picks role_ids up on the first read after upgrading, and their
// configuration does not mention it; that must not plan to take the roles off.
func TestAppRoleIDsLeftAloneWhenUnset(t *testing.T) {
	for name, newResource := range appResourcesUnderTest() {
		t.Run(name, func(t *testing.T) {
			r := newResource()
			diff := appDiff(t, r, roleIDsState(t, r, 380586, 406973), map[string]interface{}{
				"name":         "my SAML APP",
				"connector_id": 110016,
			})

			if diff == nil {
				return
			}
			for key, change := range diff.Attributes {
				if strings.HasPrefix(key, "role_ids") {
					t.Fatalf("expected role_ids to be left alone, got %s: %q -> %q", key, change.Old, change.New)
				}
			}
		})
	}
}

// TestAppRoleIDsSettles guards the other direction: roles a configuration
// asked for and got are not proposed again, whatever order they were written
// in. The API does not promise an order, which is why this is a set.
func TestAppRoleIDsSettles(t *testing.T) {
	for name, newResource := range appResourcesUnderTest() {
		t.Run(name, func(t *testing.T) {
			r := newResource()
			diff := appDiff(t, r, roleIDsState(t, r, 380586, 406973), map[string]interface{}{
				"name":         "my SAML APP",
				"connector_id": 110016,
				"role_ids":     []interface{}{406973, 380586},
			})

			if diff == nil {
				return
			}
			for key, change := range diff.Attributes {
				if strings.HasPrefix(key, "role_ids") {
					t.Fatalf("expected role_ids to settle, got %s: %q -> %q", key, change.Old, change.New)
				}
			}
		})
	}
}

// TestAppRoleIDsBody covers what reaches the API, on every resource.
func TestAppRoleIDsBody(t *testing.T) {
	for name, newResource := range appResourcesUnderTest() {
		t.Run(name, func(t *testing.T) {
			r := newResource()

			t.Run("create sends every configured role", func(t *testing.T) {
				got := appBody(t, r, nil, map[string]interface{}{
					"name":         "my SAML APP",
					"connector_id": 110016,
					"role_ids":     []interface{}{380586, 406973},
				}, "role_ids", addAppAssignmentForCreate)

				if ids := bodyRoleIDs(got); ids != "380586,406973" {
					t.Fatalf("expected both roles to be sent, got %s", got)
				}
			})

			t.Run("create omits an absent role_ids", func(t *testing.T) {
				got := appBody(t, r, nil, map[string]interface{}{
					"name":         "my SAML APP",
					"connector_id": 110016,
				}, "role_ids", addAppAssignmentForCreate)

				if strings.Contains(got, `"role_ids"`) {
					t.Fatalf("expected role_ids to be omitted, got %s", got)
				}
			})

			// Same as an explicit 0 for policy_id: a create has no roles to
			// take off, and an app created without role_ids comes back with
			// [] anyway.
			t.Run("create omits an empty role_ids", func(t *testing.T) {
				got := appBody(t, r, nil, map[string]interface{}{
					"name":         "my SAML APP",
					"connector_id": 110016,
					"role_ids":     []interface{}{},
				}, "role_ids", addAppAssignmentForCreate)

				if strings.Contains(got, `"role_ids"`) {
					t.Fatalf("expected role_ids to be omitted on create, got %s", got)
				}
			})

			t.Run("update replaces the roles with the configured set", func(t *testing.T) {
				got := appBody(t, r, roleIDsState(t, r, 380586, 406973), map[string]interface{}{
					"name":         "my SAML APP",
					"connector_id": 110016,
					"role_ids":     []interface{}{391265},
				}, "role_ids", addAppAssignmentForUpdate)

				if ids := bodyRoleIDs(got); ids != "391265" {
					t.Fatalf("expected role_ids to be replaced with [391265], got %s", got)
				}
			})

			// The case policy_id and brand_id need a 0 sentinel for. An empty
			// set is a value, not a null, so it diffs on its own and has to
			// reach the API as [] rather than being dropped.
			t.Run("update sends [] to take every role off", func(t *testing.T) {
				got := appBody(t, r, roleIDsState(t, r, 380586, 406973), map[string]interface{}{
					"name":         "my SAML APP",
					"connector_id": 110016,
					"role_ids":     []interface{}{},
				}, "role_ids", addAppAssignmentForUpdate)

				if !strings.Contains(got, `"role_ids":[]`) {
					t.Fatalf("expected role_ids to be sent as [], got %s", got)
				}
			})

			t.Run("update omits unchanged roles", func(t *testing.T) {
				got := appBody(t, r, roleIDsState(t, r, 380586, 406973), map[string]interface{}{
					"name":         "a renamed app",
					"connector_id": 110016,
					"role_ids":     []interface{}{380586, 406973},
				}, "role_ids", addAppAssignmentForUpdate)

				if strings.Contains(got, `"role_ids"`) {
					t.Fatalf("expected an unrelated update to leave role_ids out, got %s", got)
				}
			})

			t.Run("update omits role_ids the configuration does not mention", func(t *testing.T) {
				got := appBody(t, r, roleIDsState(t, r, 380586, 406973), map[string]interface{}{
					"name":         "a renamed app",
					"connector_id": 110016,
				}, "role_ids", addAppAssignmentForUpdate)

				if strings.Contains(got, `"role_ids"`) {
					t.Fatalf("expected role_ids to be left out, got %s", got)
				}
			})
		})
	}
}

// bodyRoleIDs pulls the role_ids array out of a request body as a sorted,
// comma-separated string, since a set carries no order.
func bodyRoleIDs(body string) string {
	_, rest, found := strings.Cut(body, `"role_ids":[`)
	if !found {
		return ""
	}
	list, _, _ := strings.Cut(rest, "]")
	ids := strings.Split(list, ",")
	sort.Strings(ids)
	return strings.Join(ids, ",")
}

// TestAppReadFieldsCoverEveryAssignment holds the shared read list to the
// assignments the create and update maps send. A field that is sent but never
// read back cannot show drift, which for role_ids would quietly break the
// promise that a configured set is authoritative.
func TestAppReadFieldsCoverEveryAssignment(t *testing.T) {
	read := make(map[string]bool, len(appBasicReadFields))
	for _, field := range appBasicReadFields {
		read[field] = true
	}
	for _, field := range []string{"policy_id", "brand_id", "role_ids"} {
		if !read[field] {
			t.Errorf("%s is sent on create and update but missing from appBasicReadFields", field)
		}
	}
}

// TestAppRoleIDsReadFromAPI puts the shape the API actually returns -- JSON
// numbers, so float64 -- through the shared read fields, and checks it lands in
// state as the IDs on every resource.
func TestAppRoleIDsReadFromAPI(t *testing.T) {
	for name, newResource := range appResourcesUnderTest() {
		t.Run(name, func(t *testing.T) {
			d := newResource().Data(nil)
			utils.SetResourceFields(d, map[string]interface{}{
				"role_ids": []interface{}{float64(380586), float64(406973)},
			}, appBasicReadFields)

			got := d.Get("role_ids").(*schema.Set)
			if got.Len() != 2 || !got.Contains(380586) || !got.Contains(406973) {
				t.Fatalf("expected role_ids [380586 406973] in state, got %v", got.List())
			}
		})
	}
}

// TestAppRoleIDsNonPositiveRejected catches a role ID that cannot exist at plan
// time, rather than leaving it to the 422 the API answers with.
func TestAppRoleIDsNonPositiveRejected(t *testing.T) {
	for name, newResource := range appResourcesUnderTest() {
		for _, bad := range []int{0, -1} {
			t.Run(fmt.Sprintf("%s/%d", name, bad), func(t *testing.T) {
				diags := schema.InternalMap(newResource().Schema).Validate(terraform.NewResourceConfigRaw(map[string]interface{}{
					"name":         "my SAML APP",
					"connector_id": 110016,
					"role_ids":     []interface{}{380586, bad},
				}))

				if !diags.HasError() {
					t.Fatalf("expected a role ID of %d to be rejected", bad)
				}
			})
		}
	}
}

// TestAccSAMLApp_roleIDs walks role_ids through everything it can mean against
// the real API: two roles in one write, replaced by a third, all taken off with
// [], put back, a role added outside Terraform taken off again, and finally
// left alone by a configuration that stops mentioning them.
//
// The unit tests above stop at the request body. This is what covers the API
// accepting role_ids on create as well as update, and the read bringing them
// back so that every step settles.
func TestAccSAMLApp_roleIDs(t *testing.T) {
	const app = "onelogin_saml_apps.role_ids"

	// Captured by one step for the next one's PreConfig, which cannot see state.
	var appID, strayRoleID int

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSAMLAppRoleIDsConfig("role_ids = [onelogin_roles.role_ids_0.id, onelogin_roles.role_ids_1.id]"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(app, "role_ids.#", "2"),
					testAccCheckAppRoles(app, "onelogin_roles.role_ids_0", "onelogin_roles.role_ids_1"),
				),
			},
			{
				Config: testAccSAMLAppRoleIDsConfig("role_ids = [onelogin_roles.role_ids_2.id]"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(app, "role_ids.#", "1"),
					testAccCheckAppRoles(app, "onelogin_roles.role_ids_2"),
				),
			},
			{
				Config: testAccSAMLAppRoleIDsConfig("role_ids = []"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(app, "role_ids.#", "0"),
					testAccCheckAppRoles(app),
				),
			},
			{
				Config: testAccSAMLAppRoleIDsConfig("role_ids = [onelogin_roles.role_ids_0.id]"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppRoles(app, "onelogin_roles.role_ids_0"),
					testAccCaptureIDs(map[string]*int{app: &appID, "onelogin_roles.role_ids_1": &strayRoleID}),
				),
			},
			{
				// Once set, role_ids is authoritative: a role given the app
				// anywhere else is drift, and the apply takes it off again.
				PreConfig: func() {
					roleIDs := []int{strayRoleID}
					if err := testAccAddAppRoles(appID, roleIDs); err != nil {
						t.Fatal(err)
					}
				},
				Config: testAccSAMLAppRoleIDsConfig("role_ids = [onelogin_roles.role_ids_0.id]"),
				Check:  testAccCheckAppRoles(app, "onelogin_roles.role_ids_0"),
			},
			{
				// Dropping the argument is not a request for no roles; that
				// is what [] is for.
				Config: testAccSAMLAppRoleIDsConfig(""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(app, "role_ids.#", "1"),
					testAccCheckAppRoles(app, "onelogin_roles.role_ids_0"),
				),
			},
		},
	})
}

// testAccSAMLAppRoleIDsConfig is three roles and a SAML app, with roleIDs
// written into the app as-is.
//
// The roles skip their apps refresh for the reason given on
// testAccAppRoleAttachmentConfig, and the app names a configuration key for
// the same reason too.
func testAccSAMLAppRoleIDsConfig(roleIDs string) string {
	roles := make([]string, 3)
	for i := range roles {
		roles[i] = fmt.Sprintf(`
resource "onelogin_roles" "role_ids_%d" {
  name                    = "TF Acc role_ids Role %d"
  skip_membership_refresh = ["apps"]
}`, i, i)
	}

	return strings.Join(roles, "\n") + fmt.Sprintf(`

resource "onelogin_saml_apps" "role_ids" {
  connector_id = 110016
  name         = "TF Acc role_ids App"
  %s

  configuration = {
    signature_algorithm = "SHA-256"
  }
}
`, roleIDs)
}

// testAccCaptureIDs copies the IDs of the named resources out of state.
func testAccCaptureIDs(into map[string]*int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for name, dst := range into {
			rs, ok := s.RootModule().Resources[name]
			if !ok {
				return fmt.Errorf("%s not found in state", name)
			}
			id, err := strconv.Atoi(rs.Primary.ID)
			if err != nil {
				return err
			}
			*dst = id
		}
		return nil
	}
}

// testAccAddAppRoles gives the app more roles directly through the API, the
// way an admin in the OneLogin UI would, and confirms they took. Without that
// the step after it would pass whether or not there was any drift to correct.
func testAccAddAppRoles(appID int, roleIDs []int) error {
	client := testAccProvider.Meta().(*onelogin.OneloginSDK)
	current, err := testAccAppRoleIDs(client, appID)
	if err != nil {
		return err
	}
	want := append(append([]int{}, current...), roleIDs...)
	if _, err := client.UpdateApp(appID, models.App{RoleIDs: &want}); err != nil {
		return err
	}

	got, err := testAccAppRoleIDs(client, appID)
	if err != nil {
		return err
	}
	sort.Ints(want)
	sort.Ints(got)
	if fmt.Sprint(got) != fmt.Sprint(want) {
		return fmt.Errorf("app %d has roles %v after adding %v, want %v", appID, got, roleIDs, want)
	}
	return nil
}
