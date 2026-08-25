package onelogin

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
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
					// The API drops reference rather than storing it, so this
					// holds only because the read leaves a configured value
					// alone instead of writing the null back over it.
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
				ResourceName:      "onelogin_groups.test",
				ImportState:       true,
				ImportStateVerify: true,
				// reference cannot survive an import: the API never returns
				// it, and an import has no prior state to preserve it from.
				// An imported group needs the attribute written back into the
				// configuration by hand, which is the best the API allows.
				ImportStateVerifyIgnore: []string{"reference"},
			},
		},
	})
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
