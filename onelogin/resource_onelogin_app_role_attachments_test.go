package onelogin

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
)

// TestLockAppRolesSerialisesOneApp is the race from #272 without the API: each
// worker copies the shared list, waits, and writes back its copy plus one
// entry, which is the attachment's GET, add, PUT. Unserialised, workers that
// read before another writes drop that write. The acceptance test below shows
// the same against the real API, but only this one runs in CI.
func TestLockAppRolesSerialisesOneApp(t *testing.T) {
	const appID, workers = -272, 20

	var roles []int
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(role int) {
			defer wg.Done()
			defer lockAppRoles(appID)()

			read := append([]int(nil), roles...)
			time.Sleep(time.Millisecond)
			roles = append(read, role)
		}(i)
	}
	wg.Wait()

	if len(roles) != workers {
		t.Fatalf("expected all %d roles to survive, got %d: %v", workers, len(roles), roles)
	}
}

// TestLockAppRolesIndependentPerApp checks the lock is per app, so attachments
// to different apps still run in parallel.
func TestLockAppRolesIndependentPerApp(t *testing.T) {
	defer lockAppRoles(-2721)()

	done := make(chan struct{})
	go func() {
		lockAppRoles(-2722)()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("locking one app blocked another")
	}
}

// testAccAttachmentRoleCount is how many attachments race for one app. Each
// attachment rewrites the app's whole role_ids array, so any two running at
// once can lose a role; five gives the default parallelism of ten plenty of
// chances to interleave.
const testAccAttachmentRoleCount = 5

// TestAccAppRoleAttachment_parallel is gh issue 272: several attachments on one
// app, applied concurrently.
//
// Each attachment reads the app's role_ids, adds its own role and writes the
// array back. Two doing that at once both read the same array, the later write
// drops the role the earlier one added, and that attachment's read fails with
// "App <id> does not have role <id>". The race does not lose every time, which
// is why this uses several roles rather than two.
//
// Removal has the same race in reverse, and the second step covers it: three
// attachments deleted at once, with the app and every role kept, so a removal
// that another one overwrote would leave its role behind.
func TestAccAppRoleAttachment_parallel(t *testing.T) {
	const kept = 2

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAppRoleAttachmentConfig(testAccAttachmentRoleCount, testAccAttachmentRoleCount),
				Check:  testAccCheckAppRoles("onelogin_saml_apps.attach", testAccAttachmentRoleNames(testAccAttachmentRoleCount)...),
			},
			{
				Config: testAccAppRoleAttachmentConfig(testAccAttachmentRoleCount, kept),
				Check:  testAccCheckAppRoles("onelogin_saml_apps.attach", testAccAttachmentRoleNames(kept)...),
			},
		},
	})
}

func testAccAttachmentRoleNames(count int) []string {
	names := make([]string, count)
	for i := range names {
		names[i] = fmt.Sprintf("onelogin_roles.attach_%d", i)
	}
	return names
}

// testAccCheckAppRoles asks the API which roles the app has and compares them
// with the IDs of the given role resources, exactly.
//
// The API rather than state, because app state is written by the app's own
// read, which runs before any attachment and so cannot see what they did.
func testAccCheckAppRoles(app string, roles ...string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[app]
		if !ok {
			return fmt.Errorf("%s not found in state", app)
		}
		appID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		want := make([]int, 0, len(roles))
		for _, role := range roles {
			rs, ok := s.RootModule().Resources[role]
			if !ok {
				return fmt.Errorf("%s not found in state", role)
			}
			id, err := strconv.Atoi(rs.Primary.ID)
			if err != nil {
				return err
			}
			want = append(want, id)
		}

		got, err := testAccAppRoleIDs(testAccProvider.Meta().(*onelogin.OneloginSDK), appID)
		if err != nil {
			return err
		}

		sort.Ints(want)
		sort.Ints(got)
		if fmt.Sprint(got) != fmt.Sprint(want) {
			return fmt.Errorf("app %d has roles %v, want %v", appID, got, want)
		}
		return nil
	}
}

// testAccAppRoleIDs reads the app's roles from the API.
func testAccAppRoleIDs(client *onelogin.OneloginSDK, appID int) ([]int, error) {
	result, err := client.GetAppByID(appID, nil)
	if err != nil {
		return nil, err
	}
	appMap, ok := result.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected app response %T", result)
	}
	raw, _ := appMap["role_ids"].([]interface{})
	ids := make([]int, 0, len(raw))
	for _, v := range raw {
		ids = append(ids, int(v.(float64)))
	}
	return ids, nil
}

// testAccAppRoleAttachmentConfig builds an app, roleCount roles, and
// attachments for the first attachmentCount of them.
//
// The roles skip their apps refresh. onelogin_roles.apps is Optional and not
// Computed, so a role whose configuration leaves it out reads an attached app
// back as drift and plans to remove it on every run.
//
// The app names one configuration key because a SAML app created without any
// records every key the API sends back, which the next plan then proposes to
// remove. That is a separate bug, and naming a key keeps it out of this test.
func testAccAppRoleAttachmentConfig(roleCount, attachmentCount int) string {
	roles := make([]string, roleCount)
	for i := range roles {
		roles[i] = fmt.Sprintf(`
resource "onelogin_roles" "attach_%d" {
  name                    = "TF Acc Attachment Role %d"
  skip_membership_refresh = ["apps"]
}`, i, i)
	}

	// count rather than the issue's for_each, which applies just as
	// concurrently: this SDK's test harness cannot read a string-keyed
	// instance back out of state.
	refs := testAccAttachmentRoleNames(roleCount)
	for i := range refs {
		refs[i] += ".id"
	}

	return strings.Join(roles, "\n") + fmt.Sprintf(`

resource "onelogin_saml_apps" "attach" {
  connector_id = 110016
  name         = "TF Acc Attachment App"

  configuration = {
    signature_algorithm = "SHA-256"
  }
}

resource "onelogin_app_role_attachments" "attach" {
  count = %d

  app_id  = onelogin_saml_apps.attach.id
  role_id = [%s][count.index]
}
`, attachmentCount, strings.Join(refs, ", "))
}
