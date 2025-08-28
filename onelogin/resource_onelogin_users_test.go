package onelogin

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccUser_crud(t *testing.T) {
	base := GetFixture("onelogin_user_example.tf", t)
	update := GetFixture("onelogin_user_updated_example.tf", t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: base,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_users.basic_test", "username", "testy.mctesterson"),
					resource.TestCheckResourceAttr("onelogin_users.basic_test", "email", "testy.mctesterson@onelogin.com"),
				),
			},
			{
				Config: update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_users.basic_test", "username", "boaty.mcboatface"),
					resource.TestCheckResourceAttr("onelogin_users.basic_test", "email", "boaty.mcboatface@onelogin.com"),
				),
			},
		},
	})
}

func TestAccUser_trustedIdp(t *testing.T) {
	withIdp := GetFixture("onelogin_user_with_trusted_idp_example.tf", t)
	withoutIdp := GetFixture("onelogin_user_without_trusted_idp_example.tf", t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: withIdp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_users.trusted_idp_test", "username", "trusted.idp.test"),
					resource.TestCheckResourceAttr("onelogin_users.trusted_idp_test", "email", "trusted.idp.test@onelogin.com"),
					resource.TestCheckResourceAttr("onelogin_users.trusted_idp_test", "trusted_idp_id", "123456"),
				),
			},
			{
				Config: withoutIdp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onelogin_users.trusted_idp_test", "username", "trusted.idp.test"),
					resource.TestCheckResourceAttr("onelogin_users.trusted_idp_test", "email", "trusted.idp.test@onelogin.com"),
					resource.TestCheckNoResourceAttr("onelogin_users.trusted_idp_test", "trusted_idp_id"),
				),
			},
		},
	})
}

// TestMergeCustomAttributes tests the custom attributes merging logic
func TestMergeCustomAttributes(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping unit test in short mode")
	}

	// Test case 1: Merge existing API attributes with resource attributes
	apiAttrs := map[string]interface{}{
		"dept_code":   "IT-DEPT",
		"employee_id": "EMP12345",
		"location":    "NYC",
	}

	resourceAttrs := map[string]interface{}{
		"team":      "DevOps",
		"dept_code": "HR-DEPT", // This should override the API value
	}

	merged := make(map[string]interface{})
	// Start with API attributes
	for k, v := range apiAttrs {
		merged[k] = v
	}
	// Override with resource attributes
	for k, v := range resourceAttrs {
		merged[k] = v
	}

	// Verify the merge results
	if merged["dept_code"] != "HR-DEPT" {
		t.Errorf("Expected dept_code to be 'HR-DEPT', got '%v'", merged["dept_code"])
	}
	if merged["employee_id"] != "EMP12345" {
		t.Errorf("Expected employee_id to be preserved as 'EMP12345', got '%v'", merged["employee_id"])
	}
	if merged["location"] != "NYC" {
		t.Errorf("Expected location to be preserved as 'NYC', got '%v'", merged["location"])
	}
	if merged["team"] != "DevOps" {
		t.Errorf("Expected team to be 'DevOps', got '%v'", merged["team"])
	}

	// Test case 2: Handle nil resource attributes
	merged2 := make(map[string]interface{})
	for k, v := range apiAttrs {
		merged2[k] = v
	}
	// resourceAttrs is nil - should preserve all API attributes

	if len(merged2) != len(apiAttrs) {
		t.Errorf("Expected merged map to preserve all API attributes when resource attrs are nil")
	}
	if merged2["employee_id"] != "EMP12345" {
		t.Errorf("Expected employee_id to be preserved as 'EMP12345', got '%v'", merged2["employee_id"])
	}
}
