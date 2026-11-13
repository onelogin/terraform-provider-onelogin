package onelogin

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
)

func TestAccOneLoginGroup_crud(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { TestAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOneLoginGroupDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOneLoginGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_groups.test", "name", "Test Group"),
					resource.TestCheckResourceAttr("onelogin_groups.test", "reference", "test-group-ref"),
				),
			},
			{
				Config: testAccCheckOneLoginGroupConfigUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_groups.test", "name", "Updated Test Group"),
					resource.TestCheckResourceAttr("onelogin_groups.test", "reference", "updated-group-ref"),
				),
			},
			{
				// reference included: the read takes it from api/1.
				ResourceName:      "onelogin_groups.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Changes made outside Terraform. A reference changed in OneLogin used to be
// invisible (api/2 never returns it), and a group deleted in OneLogin failed
// every plan with a 404 instead of being recreated.
func TestAccOneLoginGroup_outOfBand(t *testing.T) {
	suffix := acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	name := "tf-acc-group-oob-" + suffix
	config := fmt.Sprintf(`
resource "onelogin_groups" "test" {
  name      = %q
  reference = "managed"
}
`, name)

	var groupID int
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { TestAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOneLoginGroupDestroyed,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccGroupID("onelogin_groups.test", &groupID),
					testAccCheckGroupReference(&groupID, "managed"),
				),
			},
			{
				// Changed in OneLogin. The refresh sees it and the apply
				// puts the configured value back.
				PreConfig: func() {
					client := testAccProvider.Meta().(*onelogin.OneloginSDK)
					changed := "changed-outside"
					if _, err := client.UpdateGroup(groupID, &models.Group{Name: name, Reference: &changed}); err != nil {
						t.Fatalf("changing reference outside Terraform: %v", err)
					}
				},
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_groups.test", "reference", "managed"),
					testAccCheckGroupReference(&groupID, "managed"),
				),
			},
			{
				// Deleted in OneLogin. The refresh drops it and the apply
				// creates it again.
				PreConfig: func() {
					client := testAccProvider.Meta().(*onelogin.OneloginSDK)
					if _, err := client.DeleteGroup(groupID); err != nil {
						t.Fatalf("deleting group outside Terraform: %v", err)
					}
				},
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						previous := groupID
						if err := testAccGroupID("onelogin_groups.test", &groupID)(s); err != nil {
							return err
						}
						if groupID == previous {
							return fmt.Errorf("expected a new group, still %d", groupID)
						}
						return nil
					},
					resource.TestCheckResourceAttr("onelogin_groups.test", "name", name),
					testAccCheckGroupReference(&groupID, "managed"),
				),
			},
		},
	})
}

// testAccGroupID records the ID of a group in state.
func testAccGroupID(name string, id *int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("%s not in state", name)
		}
		parsed, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}
		*id = parsed
		return nil
	}
}

// testAccCheckGroupReference checks the reference OneLogin holds, read from
// api/1 rather than from state.
func testAccCheckGroupReference(id *int, want string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*onelogin.OneloginSDK)
		resp, err := client.GetGroupByID(*id)
		if err != nil {
			return err
		}
		group, err := groupFromResponse(resp, *id)
		if err != nil {
			return err
		}
		if got, _ := group["reference"].(string); got != want {
			return fmt.Errorf("group %d has reference %q in OneLogin, want %q", *id, got, want)
		}
		return nil
	}
}

type stubResponse struct {
	status int
	body   string
}

// groupResourceRead runs the resource's read for group 1 against a stub API,
// through the real SDK, starting from a state whose reference is "old".
func groupResourceRead(t *testing.T, routes map[string]stubResponse) (*schema.ResourceData, error) {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/auth/oauth2/v2/token" {
			io.WriteString(w, `{"access_token":"x","token_type":"bearer","expires_in":36000}`)
			return
		}
		route, ok := routes[r.URL.Path]
		if !ok {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(route.status)
		io.WriteString(w, route.body)
	}))
	t.Cleanup(srv.Close)

	t.Setenv("ONELOGIN_API_URL", srv.URL)
	t.Setenv("ONELOGIN_SUBDOMAIN", "test")
	t.Setenv("ONELOGIN_CLIENT_ID", "id")
	t.Setenv("ONELOGIN_CLIENT_SECRET", "secret")
	client, err := onelogin.NewOneloginSDK()
	if err != nil {
		t.Fatalf("client: %v", err)
	}

	d := resourceOneLoginGroups().Data(&terraform.InstanceState{ID: "1"})
	if err := d.Set("reference", "old"); err != nil {
		t.Fatal(err)
	}
	if diags := groupRead(context.Background(), d, client); diags.HasError() {
		return d, fmt.Errorf("%s", diags[0].Summary)
	}
	return d, nil
}

const groupV2Body = `{"id":1,"name":"Engineering","policy_id":null,"admins":[],"users":[],"policy":{}}`

// The read used to keep whatever reference state held, because api/2 never
// returns one. It now comes from api/1.
func TestGroupReadReferenceFromV1(t *testing.T) {
	d, err := groupResourceRead(t, map[string]stubResponse{
		"/api/2/groups/1": {http.StatusOK, groupV2Body},
		"/api/1/groups/1": {http.StatusOK, `{"status":{"error":false,"code":200},"data":[{"id":1,"name":"Engineering","reference":"changed"}]}`},
	})
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if got := d.Get("reference"); got != "changed" {
		t.Errorf("reference = %q, want \"changed\"", got)
	}
	if got := d.Get("name"); got != "Engineering" {
		t.Errorf("name = %q, want \"Engineering\"", got)
	}
}

// A reference cleared in OneLogin reads back as "", so a configuration that
// sets one sees the change.
func TestGroupReadNullReference(t *testing.T) {
	d, err := groupResourceRead(t, map[string]stubResponse{
		"/api/2/groups/1": {http.StatusOK, groupV2Body},
		"/api/1/groups/1": {http.StatusOK, `{"status":{"error":false,"code":200},"data":[{"id":1,"name":"Engineering","reference":null}]}`},
	})
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if got := d.Get("reference"); got != "" {
		t.Errorf("reference = %q, want \"\"", got)
	}
}

// A group deleted outside Terraform is dropped from state, so the plan
// recreates it, whichever of the two reads finds it gone.
func TestGroupReadNotFoundRemovesFromState(t *testing.T) {
	notFound := stubResponse{http.StatusNotFound, `{"message":"The resource with the given id could not be found","statusCode":404}`}
	cases := map[string]map[string]stubResponse{
		"api/2": {"/api/2/groups/1": notFound},
		"api/1": {
			"/api/2/groups/1": {http.StatusOK, groupV2Body},
			"/api/1/groups/1": notFound,
		},
	}
	for name, routes := range cases {
		t.Run(name, func(t *testing.T) {
			d, err := groupResourceRead(t, routes)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if d.Id() != "" {
				t.Errorf("expected the group removed from state, id is %q", d.Id())
			}
		})
	}
}

// Anything but a 404 is still an error.
func TestGroupReadServerErrorFails(t *testing.T) {
	d, err := groupResourceRead(t, map[string]stubResponse{
		"/api/2/groups/1": {http.StatusOK, groupV2Body},
		"/api/1/groups/1": {http.StatusInternalServerError, `{}`},
	})
	if err == nil {
		t.Fatal("expected an error")
	}
	if d.Id() != "1" {
		t.Errorf("expected the group kept in state, id is %q", d.Id())
	}
}

func testAccCheckOneLoginGroupDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*onelogin.OneloginSDK)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "onelogin_groups" {
			continue
		}
		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}
		resp, err := client.GetGroupByIDV2(id)
		if err == nil && resp != nil {
			return fmt.Errorf("group %d still exists", id)
		}
	}
	return nil
}

const testAccCheckOneLoginGroupConfig = `
resource "onelogin_groups" "test" {
  name      = "Test Group"
  reference = "test-group-ref"
}
`

const testAccCheckOneLoginGroupConfigUpdated = `
resource "onelogin_groups" "test" {
  name      = "Updated Test Group"
  reference = "updated-group-ref"
}
`

// groupState builds the state a read leaves behind for a group.
func groupState(t *testing.T, policyID int) *schema.ResourceData {
	t.Helper()
	d := resourceOneLoginGroups().Data(nil)
	d.SetId("590390")
	for key, value := range map[string]interface{}{
		"name":      "Engineering",
		"reference": "eng",
		"policy_id": policyID,
	} {
		if err := d.Set(key, value); err != nil {
			t.Fatal(err)
		}
	}
	return d
}

func TestGroupPolicyIDNoPerpetualDiff(t *testing.T) {
	t.Run("an assigned policy settles", func(t *testing.T) {
		d := groupState(t, 955633)
		config := map[string]interface{}{
			"name":      "Engineering",
			"reference": "eng",
			"policy_id": 955633,
		}

		diff, err := resourceOneLoginGroups().Diff(context.Background(), d.State(), terraform.NewResourceConfigRaw(config), nil)
		if err != nil {
			t.Fatalf("diff: %v", err)
		}
		if diff != nil && len(diff.Attributes) > 0 {
			for attribute, change := range diff.Attributes {
				t.Errorf("unexpected diff on %s: %q -> %q", attribute, change.Old, change.New)
			}
		}
	})

	// A group that never had a policy reads back 0, which is also what an
	// absent attribute holds. If those did not agree every plan would propose
	// the same change.
	t.Run("an unassigned group settles", func(t *testing.T) {
		d := groupState(t, 0)
		config := map[string]interface{}{
			"name":      "Engineering",
			"reference": "eng",
		}

		diff, err := resourceOneLoginGroups().Diff(context.Background(), d.State(), terraform.NewResourceConfigRaw(config), nil)
		if err != nil {
			t.Fatalf("diff: %v", err)
		}
		if diff != nil && len(diff.Attributes) > 0 {
			for attribute, change := range diff.Attributes {
				t.Errorf("unexpected diff on %s: %q -> %q", attribute, change.Old, change.New)
			}
		}
	})
}

// Removing the attribute has to be expressible. Were policy_id Computed, an
// empty configuration would read as "keep what is in state" and a policy could
// never be taken off a group once set.
func TestGroupPolicyIDRemovable(t *testing.T) {
	d := groupState(t, 955633)
	config := map[string]interface{}{
		"name":      "Engineering",
		"reference": "eng",
	}

	diff, err := resourceOneLoginGroups().Diff(context.Background(), d.State(), terraform.NewResourceConfigRaw(config), nil)
	if err != nil {
		t.Fatalf("diff: %v", err)
	}
	if diff == nil {
		t.Fatal("expected removing policy_id to propose a change")
	}
	change, ok := diff.Attributes["policy_id"]
	if !ok {
		t.Fatalf("expected a diff on policy_id, got %v", diff.Attributes)
	}
	if change.Old != "955633" || change.New != "0" {
		t.Fatalf("expected 955633 -> 0, got %q -> %q", change.Old, change.New)
	}
}

// The other half of the feature request this resource pair exists for: create
// a policy, attach it to a group, and take it off again.
func TestAccOneLoginGroup_policy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { TestAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOneLoginGroupDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOneLoginGroupPolicyConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"onelogin_groups.policy_test", "policy_id",
						"onelogin_policies.group_policy", "id",
					),
				),
			},
			{
				// policy_id dropped from the configuration. The assignment is
				// cleared rather than quietly preserved, and the plan that
				// follows is empty.
				Config: testAccCheckOneLoginGroupPolicyConfigCleared,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_groups.policy_test", "policy_id", "0"),
				),
			},
		},
	})
}

const testAccCheckOneLoginGroupPolicyConfig = `
resource "onelogin_policies" "group_policy" {
  name = "TF Acc Group Policy"
  kind = "user"
}

resource "onelogin_groups" "policy_test" {
  name      = "TF Acc Policy Group"
  policy_id = onelogin_policies.group_policy.id
}
`

const testAccCheckOneLoginGroupPolicyConfigCleared = `
resource "onelogin_policies" "group_policy" {
  name = "TF Acc Group Policy"
  kind = "user"
}

resource "onelogin_groups" "policy_test" {
  name = "TF Acc Policy Group"
}
`
