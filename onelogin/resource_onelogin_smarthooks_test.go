package onelogin

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	ol "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
)

// skipIfHookTypeExists skips when the tenant already has a hook of this type.
//
// A synchronous hook is a singleton: creating a second pre-authentication hook
// is refused with
//
//	409 The 'pre-authentication' hook is synchronous and can only have one
//	defined function
//
// so on any tenant with one configured this test cannot run. The alternative is
// deleting whatever is there, and on a shared tenant that is somebody's live
// authentication path -- not something to do for a green tick.
func skipIfHookTypeExists(t *testing.T, hookType string) {
	t.Helper()

	// The SDK still insists on a subdomain even though every call goes to
	// ONELOGIN_API_URL, and nothing has configured the provider this early.
	subdomain, err := subdomainFromURL(os.Getenv("ONELOGIN_API_URL"))
	if err != nil {
		t.Fatalf("could not work out the subdomain to check for an existing %s hook: %v", hookType, err)
	}
	os.Setenv("ONELOGIN_SUBDOMAIN", subdomain)

	client, err := ol.NewOneloginSDK()
	if err != nil {
		t.Fatalf("could not build a client to check for an existing %s hook: %v", hookType, err)
	}

	result, err := client.ListHooks(nil)
	if err != nil {
		t.Fatalf("could not list hooks: %v", err)
	}

	hooks, ok := result.([]interface{})
	if !ok {
		// Not a shape this check understands. Say so and let the test run
		// rather than skipping on a guess.
		t.Logf("could not read the hook list (%T); running anyway", result)
		return
	}

	for _, raw := range hooks {
		hook, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		if hook["type"] == hookType {
			t.Skipf("this tenant already has a %s hook, and it is a singleton: the API refuses a second one with a 409", hookType)
		}
	}
}

func TestAccSmartHook_crud(t *testing.T) {
	base := GetFixture("onelogin_smarthooks_example.tf", t)
	update := GetFixture("onelogin_smarthooks_updated_example.tf", t)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			TestAccPreCheck(t)
			skipIfHookTypeExists(t, "pre-authentication")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: base,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_smarthooks.basic_test", "type", "pre-authentication"),
					resource.TestCheckResourceAttr("onelogin_smarthooks.basic_test", "retries", "0"),
					resource.TestCheckResourceAttr("onelogin_smarthooks.basic_test", "context_version", "1.0.0"),
					resource.TestCheckResourceAttr("onelogin_smarthooks.basic_test", "timeout", "1"),
					resource.TestCheckResourceAttr("onelogin_smarthooks.basic_test", "disabled", "false"),
					resource.TestCheckResourceAttr("onelogin_smarthooks.basic_test", "options.risk_enabled", "false"),
					resource.TestCheckResourceAttr("onelogin_smarthooks.basic_test", "options.location_enabled", "false"),
					resource.TestCheckResourceAttr("onelogin_smarthooks.basic_test", "packages.mysql", "2.18.1"),
					resource.TestCheckResourceAttr("onelogin_smarthooks.basic_test", "function", `ICAgIGV4cG9ydHMuaGFuZGxlciA9IGFzeW5jIGNvbnRleHQgPT4gewogICAgICBjb25zb2xlLmxvZygiUHJlLWF1dGggZXhlY3V0aW5nIGZvciAiICsgY29udGV4dC51c2VyLnVzZXJfaWRlbnRpZmllcik7CiAgICAgIHJldHVybiB7IHVzZXI6IGNvbnRleHQudXNlciB9OwogICAgfTsK`),
				),
			},
			{
				Config: update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_smarthooks.basic_test", "type", "pre-authentication"),
					resource.TestCheckResourceAttr("onelogin_smarthooks.basic_test", "retries", "0"),
					resource.TestCheckResourceAttr("onelogin_smarthooks.basic_test", "context_version", "1.0.0"),
					resource.TestCheckResourceAttr("onelogin_smarthooks.basic_test", "timeout", "1"),
					resource.TestCheckResourceAttr("onelogin_smarthooks.basic_test", "disabled", "false"),
					resource.TestCheckResourceAttr("onelogin_smarthooks.basic_test", "options.risk_enabled", "false"),
					resource.TestCheckResourceAttr("onelogin_smarthooks.basic_test", "options.location_enabled", "false"),
					resource.TestCheckResourceAttr("onelogin_smarthooks.basic_test", "packages.mysql", "2.18.1"),
					resource.TestCheckResourceAttr("onelogin_smarthooks.basic_test", "function", `ICAgIGV4cG9ydHMuaGFuZGxlciA9IGFzeW5jIGNvbnRleHQgPT4gewogICAgICBjb25zb2xlLmxvZygiUHJlLWF1dGggZXhlY3V0aW5nIGZvciAiICsgY29udGV4dC51c2VyLnVzZXJfaWRlbnRpZmllcik7CiAgICAgIHJldHVybiB7IHVzZXI6IGNvbnRleHQudXNlciB9OwogICAgfTsK`),
				),
			},
		},
	})
}
