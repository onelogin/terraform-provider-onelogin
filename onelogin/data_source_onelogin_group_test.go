package onelogin

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
)

// TestAccDataSourceOneLoginGroup reads back a group the test creates, so it
// runs on any tenant. reference is checked too: only api/1 returns it.
func TestAccDataSourceOneLoginGroup(t *testing.T) {
	suffix := acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	name := "tf-acc-group-ds-" + suffix
	reference := "tf-acc-ref-" + suffix

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { TestAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOneLoginGroupDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOneLoginGroupDataSourceConfig(name, reference),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.onelogin_group.test", "id", "onelogin_groups.test", "id"),
					resource.TestCheckResourceAttr("data.onelogin_group.test", "name", name),
					resource.TestCheckResourceAttr("data.onelogin_group.test", "reference", reference),
				),
			},
		},
	})
}

func testAccCheckOneLoginGroupDataSourceConfig(name, reference string) string {
	return fmt.Sprintf(`
resource "onelogin_groups" "test" {
  name      = %q
  reference = %q
}

data "onelogin_group" "test" {
  id = onelogin_groups.test.id
}
`, name, reference)
}

// groupDataSourceRead runs the data source's read against a stub API that
// answers api/1/groups/{id} with status and body, through the real SDK.
func groupDataSourceRead(t *testing.T, id, status int, body string) (map[string]string, error) {
	t.Helper()

	path := fmt.Sprintf("/api/1/groups/%d", id)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/auth/oauth2/v2/token":
			io.WriteString(w, `{"access_token":"x","token_type":"bearer","expires_in":36000}`)
		case path:
			w.WriteHeader(status)
			io.WriteString(w, body)
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
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

	d := dataSourceOneLoginGroup().TestResourceData()
	if err := d.Set("id", id); err != nil {
		t.Fatal(err)
	}
	if diags := dataSourceOneLoginGroupRead(context.Background(), d, client); diags.HasError() {
		return nil, fmt.Errorf("%s", diags[0].Summary)
	}
	return map[string]string{
		"id":        d.Id(),
		"name":      d.Get("name").(string),
		"reference": d.Get("reference").(string),
	}, nil
}

// The api/1 response as a test tenant returns it. Before #282 was fixed, name
// and reference were looked for beside "status" and both came back empty.
func TestDataSourceOneLoginGroupReadUnwrapsV1Response(t *testing.T) {
	got, err := groupDataSourceRead(t, 731, http.StatusOK,
		`{"status":{"error":false,"code":200,"type":"success","message":"Success"},`+
			`"data":[{"id":731,"name":"Example Group","reference":"example"}]}`)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	want := map[string]string{"id": "731", "name": "Example Group", "reference": "example"}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("%s = %q, want %q", k, got[k], v)
		}
	}
}

func TestDataSourceOneLoginGroupReadNullReference(t *testing.T) {
	got, err := groupDataSourceRead(t, 731, http.StatusOK,
		`{"status":{"error":false,"code":200},"data":[{"id":731,"name":"Example Group","reference":null}]}`)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if got["name"] != "Example Group" || got["reference"] != "" {
		t.Errorf("got name %q reference %q, want \"Example Group\" and \"\"", got["name"], got["reference"])
	}
}

func TestDataSourceOneLoginGroupReadNotFound(t *testing.T) {
	_, err := groupDataSourceRead(t, 731, http.StatusNotFound,
		`{"status":{"error":true,"code":404,"type":"not found","message":"Group for id 731 was not found"}}`)
	if err == nil || !strings.Contains(err.Error(), "status: 404") {
		t.Fatalf("expected the API's 404 for a group that does not exist, got %v", err)
	}
}

// A response the reader can't find the group in is an error, never an empty
// name. Each case here returned "" for name before the fix.
func TestGroupFromResponseRejects(t *testing.T) {
	cases := map[string]interface{}{
		"not an object": []interface{}{map[string]interface{}{"id": 731.0, "name": "Example Group"}},
		"empty data":    map[string]interface{}{"data": []interface{}{}},
		"data not a list": map[string]interface{}{
			"data": map[string]interface{}{"id": 731.0, "name": "Example Group"},
		},
		"two groups": map[string]interface{}{"data": []interface{}{
			map[string]interface{}{"id": 731.0, "name": "Example Group"},
			map[string]interface{}{"id": 732.0, "name": "Other Group"},
		}},
		"group not an object": map[string]interface{}{"data": []interface{}{"731"}},
		"other group": map[string]interface{}{"data": []interface{}{
			map[string]interface{}{"id": 732.0, "name": "Other Group"},
		}},
		"no id": map[string]interface{}{"data": []interface{}{
			map[string]interface{}{"name": "Example Group"},
		}},
		"no name": map[string]interface{}{"data": []interface{}{
			map[string]interface{}{"id": 731.0, "reference": "example"},
		}},
		"empty name": map[string]interface{}{"data": []interface{}{
			map[string]interface{}{"id": 731.0, "name": ""},
		}},
		"status only": map[string]interface{}{"status": map[string]interface{}{"error": false, "code": 200.0}},
		// api/2's shape. It never includes reference, so taking it would set
		// reference to "" for a group that has one.
		"bare object": map[string]interface{}{"id": 731.0, "name": "Example Group"},
	}
	for name, resp := range cases {
		t.Run(name, func(t *testing.T) {
			if group, err := groupFromResponse(resp, 731); err == nil {
				t.Fatalf("expected an error, got %v", group)
			}
		})
	}
}

func TestGroupFromResponseErrorNamesTheGroup(t *testing.T) {
	_, err := groupFromResponse(map[string]interface{}{"data": []interface{}{
		map[string]interface{}{"id": 732.0, "name": "Other Group"},
	}}, 731)
	if err == nil || !strings.Contains(err.Error(), "732") || !strings.Contains(err.Error(), "731") {
		t.Fatalf("expected an error naming both groups, got %v", err)
	}
}
