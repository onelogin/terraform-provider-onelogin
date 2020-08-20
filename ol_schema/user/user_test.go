package userschema

import (
	"testing"

	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/users"
	"github.com/stretchr/testify/assert"
)

func TestSchema(t *testing.T) {
	t.Run("creates and returns a map of a user Schema", func(t *testing.T) {
		provSchema := Schema()
		assert.NotNil(t, provSchema["username"])
		assert.NotNil(t, provSchema["email"])
		assert.NotNil(t, provSchema["firstname"])
		assert.NotNil(t, provSchema["lastname"])
		assert.NotNil(t, provSchema["distinguished_name"])
		assert.NotNil(t, provSchema["samaccountname"])
		assert.NotNil(t, provSchema["userprincipalname"])
		assert.NotNil(t, provSchema["member_of"])
		assert.NotNil(t, provSchema["phone"])
		assert.NotNil(t, provSchema["title"])
		assert.NotNil(t, provSchema["company"])
		assert.NotNil(t, provSchema["department"])
		assert.NotNil(t, provSchema["comment"])
		assert.NotNil(t, provSchema["state"])
		assert.NotNil(t, provSchema["status"])
		assert.NotNil(t, provSchema["group_id"])
		assert.NotNil(t, provSchema["directory_id"])
		assert.NotNil(t, provSchema["trusted_idp_id"])
		assert.NotNil(t, provSchema["manager_ad_id"])
		assert.NotNil(t, provSchema["manager_user_id"])
		assert.NotNil(t, provSchema["external_id"])
	})
}

func TestInflate(t *testing.T) {
	tests := map[string]struct {
		ResourceData   map[string]interface{}
		ExpectedOutput users.User
	}{
		"creates and returns the address of a user struct": {
			ResourceData: map[string]interface{}{
				"id":                 1,
				"username":           "username",
				"email":              "email",
				"firstname":          "firstname",
				"lastname":           "lastname",
				"distinguished_name": "distinguished_name",
				"samaccountname":     "samaccountname",
				"userprincipalname":  "userprincipalname",
				"member_of":          "member_of",
				"phone":              "phone",
				"title":              "title",
				"company":            "company",
				"department":         "department",
				"comment":            "comment",
				"state":              1,
				"status":             1,
				"group_id":           1,
				"directory_id":       1,
				"trusted_idp_id":     1,
				"manager_ad_id":      1,
				"manager_user_id":    1,
				"external_id":        1,
			},
			ExpectedOutput: users.User{
				ID:                oltypes.Int32(int32(1)),
				Username:          oltypes.String("username"),
				Email:             oltypes.String("email"),
				Firstname:         oltypes.String("firstname"),
				Lastname:          oltypes.String("lastname"),
				DistinguishedName: oltypes.String("distinguished_name"),
				Samaccountname:    oltypes.String("samaccountname"),
				UserPrincipalName: oltypes.String("userprincipalname"),
				MemberOf:          oltypes.String("member_of"),
				Phone:             oltypes.String("phone"),
				Title:             oltypes.String("title"),
				Company:           oltypes.String("company"),
				Department:        oltypes.String("department"),
				Comment:           oltypes.String("comment"),
				State:             oltypes.Int32(int32(1)),
				Status:            oltypes.Int32(int32(1)),
				GroupID:           oltypes.Int32(int32(1)),
				DirectoryID:       oltypes.Int32(int32(1)),
				TrustedIDPID:      oltypes.Int32(int32(1)),
				ManagerADID:       oltypes.Int32(int32(1)),
				ManagerUserID:     oltypes.Int32(int32(1)),
				ExternalID:        oltypes.Int32(int32(1)),
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			subj, _ := Inflate(test.ResourceData)
			assert.Equal(t, test.ExpectedOutput, subj)
		})
	}
}
