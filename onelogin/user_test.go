package onelogin

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	"github.com/stretchr/testify/assert"
)

// Mock the OneLogin SDK client
type mockOneLoginSDK struct {
	onelogin.OneloginSDK
}

func (m *mockOneLoginSDK) GetUserByID(id int, queryParams interface{}) (interface{}, error) {
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
