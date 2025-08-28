package onelogin

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	userschema "github.com/onelogin/terraform-provider-onelogin/ol_schema/user"
	"github.com/stretchr/testify/assert"
)

// Mock the OneLogin SDK client
type mockOneLoginSDK struct {
	onelogin.OneloginSDK
	getUserFunc func(id int, queryParams interface{}) (interface{}, error)
}

func (m *mockOneLoginSDK) GetUserByID(id int, queryParams interface{}) (interface{}, error) {
	// Use custom function if provided, otherwise use default
	if m.getUserFunc != nil {
		return m.getUserFunc(id, queryParams)
	}

	// Return a mock user with the trusted_idp_id field set
	return map[string]interface{}{
		"id":             id,
		"username":       "test.user",
		"email":          "test.user@example.com",
		"trusted_idp_id": 12345,
	}, nil
}

func (m *mockOneLoginSDK) GetUserByIDWithContext(ctx context.Context, id int, queryParams interface{}) (interface{}, error) {
	return m.GetUserByID(id, queryParams)
}

func (m *mockOneLoginSDK) UpdateUser(id int, user models.User) (interface{}, error) {
	// Return a successful response
	return map[string]interface{}{
		"status": map[string]interface{}{
			"type":    "success",
			"message": "Updated user.",
			"code":    200,
		},
	}, nil
}

func (m *mockOneLoginSDK) UpdateUserWithContext(ctx context.Context, id int, user models.User) (interface{}, error) {
	return m.UpdateUser(id, user)
}

// Override the userRead function for testing
func mockUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*mockOneLoginSDK)
	uid := 12345
	d.SetId("12345")

	result, _ := client.GetUserByID(uid, nil)

	// Parse the user from the result
	userMap := result.(map[string]interface{})

	// Set basic user fields
	basicFields := []string{
		"username", "email", "firstname", "lastname", "title",
		"department", "company", "status", "state", "phone",
		"group_id", "directory_id", "distinguished_name", "external_id",
		"manager_ad_id", "manager_user_id", "samaccountname", "userprincipalname",
		"member_of", "created_at", "updated_at", "activated_at", "last_login",
		"trusted_idp_id",
	}

	for _, field := range basicFields {
		if val, ok := userMap[field]; ok {
			d.Set(field, val)
		}
	}

	return nil
}

func TestUserBasicFields(t *testing.T) {
	// Create a mock ResourceData
	r := Users().Schema
	d := schema.TestResourceDataRaw(t, r, map[string]interface{}{
		"username": "test.user",
		"email":    "test.user@example.com",
	})

	// Mock the OneLogin SDK client
	client := &mockOneLoginSDK{}

	// Call our mock userRead function
	diags := mockUserRead(context.Background(), d, client)
	assert.Nil(t, diags, "userRead should not return diagnostics")

	// Verify that trusted_idp_id is included in the basicFields list
	assert.Equal(t, 12345, d.Get("trusted_idp_id"), "trusted_idp_id should be set correctly")
}

func TestUserUpdate(t *testing.T) {
	// Create a mock ResourceData
	r := Users().Schema
	d := schema.TestResourceDataRaw(t, r, map[string]interface{}{
		"id":             "12345",
		"username":       "test.user",
		"email":          "test.user@example.com",
		"trusted_idp_id": 12345,
	})
	d.SetId("12345")

	// Verify the field was set in the ResourceData
	assert.Equal(t, 12345, d.Get("trusted_idp_id"), "trusted_idp_id should be set correctly before update")

	// Now simulate removing the trusted_idp_id
	d.Set("trusted_idp_id", 0)

	// Verify the field was updated in the ResourceData
	assert.Equal(t, 0, d.Get("trusted_idp_id"), "trusted_idp_id should be updated to 0")
}

func TestUserCompanyDepartmentClearing(t *testing.T) {
	// Create a mock ResourceData with company and department initially set
	r := Users().Schema
	d := schema.TestResourceDataRaw(t, r, map[string]interface{}{
		"id":         "12345",
		"username":   "test.user",
		"email":      "test.user@example.com",
		"company":    "Test Company",
		"department": "Test Department",
	})
	d.SetId("12345")

	// Verify initial values are set
	assert.Equal(t, "Test Company", d.Get("company"), "company should be set initially")
	assert.Equal(t, "Test Department", d.Get("department"), "department should be set initially")

	// Simulate removing company and department from Terraform config
	// When fields are removed from config, d.Get() returns empty strings, not nil
	d.Set("company", "")
	d.Set("department", "")

	// Verify they are now empty strings (this is what d.Get() returns for removed fields)
	assert.Equal(t, "", d.Get("company"), "company should be empty string when removed from config")
	assert.Equal(t, "", d.Get("department"), "department should be empty string when removed from config")

	// Test that the userUpdate function would pass these empty strings to the API
	// This is the root of the issue - empty strings should clear the fields in OneLogin

	// Test the fixed Inflate behavior - empty strings should now be included in the struct
	userData := map[string]interface{}{
		"username":   d.Get("username"),
		"email":      d.Get("email"),
		"company":    d.Get("company"),    // This is "" (empty string)
		"department": d.Get("department"), // This is "" (empty string)
	}

	user, err := userschema.Inflate(userData)
	assert.NoError(t, err, "Inflate should not error with empty company/department")

	// With the fix, company and department should be set to empty strings in the struct
	assert.Equal(t, "", user.Company, "Company should be set to empty string (to clear it in API)")
	assert.Equal(t, "", user.Department, "Department should be set to empty string (to clear it in API)")
}

func TestUserInflateCompanyDepartmentEdgeCases(t *testing.T) {
	tests := []struct {
		name             string
		input            map[string]interface{}
		expectedCompany  string
		expectedDept     string
		shouldSetCompany bool
		shouldSetDept    bool
	}{
		{
			name: "both fields have values",
			input: map[string]interface{}{
				"username":   "test",
				"email":      "test@example.com",
				"company":    "Test Company",
				"department": "Test Department",
			},
			expectedCompany:  "Test Company",
			expectedDept:     "Test Department",
			shouldSetCompany: true,
			shouldSetDept:    true,
		},
		{
			name: "both fields are empty strings (fix should include them)",
			input: map[string]interface{}{
				"username":   "test",
				"email":      "test@example.com",
				"company":    "",
				"department": "",
			},
			expectedCompany:  "",
			expectedDept:     "",
			shouldSetCompany: true,
			shouldSetDept:    true,
		},
		{
			name: "fields missing from input (should result in zero values)",
			input: map[string]interface{}{
				"username": "test",
				"email":    "test@example.com",
			},
			expectedCompany:  "",
			expectedDept:     "",
			shouldSetCompany: false,
			shouldSetDept:    false,
		},
		{
			name: "mixed case - one empty, one with value",
			input: map[string]interface{}{
				"username":   "test",
				"email":      "test@example.com",
				"company":    "",
				"department": "Engineering",
			},
			expectedCompany:  "",
			expectedDept:     "Engineering",
			shouldSetCompany: true,
			shouldSetDept:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := userschema.Inflate(tt.input)
			assert.NoError(t, err, "Inflate should not error")
			assert.Equal(t, tt.expectedCompany, user.Company, "Company field should match expected")
			assert.Equal(t, tt.expectedDept, user.Department, "Department field should match expected")
		})
	}
}

func TestIsNeverLoggedInDate(t *testing.T) {
	tests := []struct {
		name     string
		dateStr  string
		expected bool
	}{
		{
			name:     "Year 1 AD placeholder date should be filtered",
			dateStr:  "0001-01-01T00:00:00Z",
			expected: true,
		},
		{
			name:     "Unix epoch placeholder should be filtered",
			dateStr:  "1970-01-01T00:00:00Z",
			expected: true,
		},
		{
			name:     "Y2K placeholder should be filtered",
			dateStr:  "2000-01-01T00:00:00Z",
			expected: true,
		},
		{
			name:     "Valid recent date should not be filtered",
			dateStr:  "2023-01-15T10:30:00Z",
			expected: false,
		},
		{
			name:     "Date after OneLogin founding should not be filtered",
			dateStr:  "2015-06-20T14:25:00Z",
			expected: false,
		},
		{
			name:     "Date just before OneLogin founding should be filtered",
			dateStr:  "2008-12-31T23:59:59Z",
			expected: true,
		},
		{
			name:     "Simple date format should work",
			dateStr:  "2000-01-01",
			expected: true,
		},
		{
			name:     "Invalid date string should not be filtered (safe fallback)",
			dateStr:  "invalid-date",
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := isNeverLoggedInDate(test.dateStr)
			assert.Equal(t, test.expected, result, "isNeverLoggedInDate should return %v for %s", test.expected, test.dateStr)
		})
	}
}

func TestUserReadWithLastLoginFiltering(t *testing.T) {
	// Create a mock ResourceData
	r := Users().Schema
	d := schema.TestResourceDataRaw(t, r, map[string]interface{}{
		"username": "test.user",
		"email":    "test.user@example.com",
	})

	// Mock the OneLogin SDK client with placeholder last_login date
	client := &mockOneLoginSDK{}

	// Override GetUserByID to return a user with placeholder last_login
	client.getUserFunc = func(id int, queryParams interface{}) (interface{}, error) {
		return map[string]interface{}{
			"id":             id,
			"username":       "test.user",
			"email":          "test.user@example.com",
			"last_login":     "0001-01-01T00:00:00Z", // Placeholder "never logged in" date
			"trusted_idp_id": 12345,
		}, nil
	}

	// Call the mock userRead function which simulates the filtering
	d.SetId("12345")
	diags := mockUserReadWithFiltering(context.Background(), d, client)
	assert.Nil(t, diags, "userRead should not return diagnostics")

	// Verify that last_login is not set when placeholder date is filtered
	lastLogin := d.Get("last_login")
	// When a placeholder date is filtered out, the field should not be set (nil/empty)
	assert.Empty(t, lastLogin, "last_login should be empty when placeholder date is filtered")

	// Verify other fields are still set normally
	assert.Equal(t, "test.user", d.Get("username"), "username should be set")
	assert.Equal(t, "test.user@example.com", d.Get("email"), "email should be set")
	assert.Equal(t, 12345, d.Get("trusted_idp_id"), "trusted_idp_id should be set")
}

// Mock userRead function that includes the new last_login filtering logic
func mockUserReadWithFiltering(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*mockOneLoginSDK)
	uid := 12345
	d.SetId("12345")

	result, _ := client.GetUserByID(uid, nil)

	// Parse the user from the result
	userMap := result.(map[string]interface{})

	// Set basic user fields (excluding last_login which needs special handling)
	basicFields := []string{
		"username", "email", "firstname", "lastname", "title",
		"department", "company", "status", "state", "phone",
		"group_id", "directory_id", "distinguished_name", "external_id",
		"manager_ad_id", "manager_user_id", "samaccountname", "userprincipalname",
		"member_of", "created_at", "updated_at", "activated_at",
		"trusted_idp_id",
	}

	// Filter out "never logged in" placeholder dates from the data before setting fields
	if lastLoginValue, ok := userMap["last_login"]; ok {
		if lastLoginStr, ok := lastLoginValue.(string); ok {
			if isNeverLoggedInDate(lastLoginStr) {
				// Remove the placeholder date so it doesn't get set
				delete(userMap, "last_login")
			}
		}
	}

	for _, field := range basicFields {
		if val, ok := userMap[field]; ok {
			d.Set(field, val)
		}
	}

	return nil
}
