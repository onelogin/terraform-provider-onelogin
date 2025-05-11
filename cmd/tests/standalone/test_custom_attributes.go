package standalone

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/onelogin/onelogin-go-sdk/pkg/client"
	"github.com/onelogin/onelogin-go-sdk/pkg/oltypes"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/users"
)

func TestCustomAttributes() {
	clientID := os.Getenv("ONELOGIN_CLIENT_ID")
	clientSecret := os.Getenv("ONELOGIN_CLIENT_SECRET")
	region := os.Getenv("ONELOGIN_REGION")
	url := os.Getenv("ONELOGIN_OAPI_URL")

	fmt.Println("Testing OneLogin Custom Attributes")
	fmt.Printf("URL: %s, Region: %s\n", url, region)

	// Create client for user operations
	oneloginClient, err := client.NewClient(&client.APIClientConfig{
		Timeout:      60,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Region:       region,
		Url:          url,
	})
	if err != nil {
		log.Fatalf("Failed to create OneLogin client: %v", err)
	}

	// Check existing users to see if our test user already exists
	testUserName := "test.custom.attr.user"
	
	query := users.UserQuery{
		Username: oltypes.String(testUserName),
	}

	existingUsers, err := oneloginClient.Services.UsersV2.Query(&query)
	if err != nil {
		log.Fatalf("Failed to query users: %v", err)
	}

	var testUser users.User
	var isNewUser bool

	if len(existingUsers) > 0 {
		// Use the existing user
		testUser = existingUsers[0]
		fmt.Printf("Found existing test user with ID: %d\n", *testUser.ID)
	} else {
		// Create a new user with a timestamp to make it unique
		uniqueUsername := fmt.Sprintf("%s_%d", testUserName, time.Now().Unix())
		uniqueEmail := fmt.Sprintf("test_%d@example.com", time.Now().Unix())
		
		testUser = users.User{
			Username:  oltypes.String(uniqueUsername),
			Email:     oltypes.String(uniqueEmail),
			Firstname: oltypes.String("Test"),
			Lastname:  oltypes.String("User"),
		}

		fmt.Printf("Creating new test user %s...\n", uniqueUsername)
		err = oneloginClient.Services.UsersV2.Create(&testUser)
		if err != nil {
			log.Fatalf("Failed to create user: %v", err)
		}
		fmt.Printf("Created user with ID: %d\n", *testUser.ID)
		isNewUser = true
	}

	// Test custom attribute - we'll try to use an existing one
	// Get all users
	allUsers, err := oneloginClient.Services.UsersV2.Query(nil)
	if err != nil {
		log.Fatalf("Failed to query all users: %v", err)
	}

	// Look for any custom attributes already in use
	var customAttrName string
	for _, user := range allUsers {
		if user.CustomAttributes != nil && len(user.CustomAttributes) > 0 {
			fmt.Printf("User %s (ID: %d) has custom attributes:\n", *user.Username, *user.ID)
			for key, value := range user.CustomAttributes {
				fmt.Printf("  %s: %v\n", key, value)
				customAttrName = key
				break
			}
			if customAttrName != "" {
				break
			}
		}
	}

	if customAttrName == "" {
		fmt.Println("No existing custom attributes found in the account.")
		fmt.Println("You need to manually create a custom attribute in the OneLogin admin portal.")
		
		// We'll try a common attribute name in case it exists
		customAttrName = "department"
		fmt.Printf("Trying with a common attribute name: %s\n", customAttrName)
	} else {
		fmt.Printf("Found existing custom attribute: %s\n", customAttrName)
	}

	// Set custom attributes for the user
	if testUser.CustomAttributes == nil {
		testUser.CustomAttributes = make(map[string]interface{})
	}
	testUser.CustomAttributes[customAttrName] = "Engineering"

	fmt.Printf("Setting custom attribute '%s' for the user...\n", customAttrName)
	err = oneloginClient.Services.UsersV2.Update(&testUser)
	if err != nil {
		log.Printf("Failed to update user with custom attributes: %v", err)
		log.Println("This may happen if the custom attribute doesn't exist.")
	} else {
		fmt.Println("Custom attributes set successfully!")

		// Verify the custom attributes were set
		fmt.Println("Retrieving user to verify custom attributes...")
		retrievedUser, err := oneloginClient.Services.UsersV2.GetOne(*testUser.ID)
		if err != nil {
			log.Fatalf("Failed to retrieve user: %v", err)
		}

		if retrievedUser.CustomAttributes != nil && len(retrievedUser.CustomAttributes) > 0 {
			fmt.Println("Custom attributes for the user:")
			for key, value := range retrievedUser.CustomAttributes {
				fmt.Printf("  %s: %v\n", key, value)
			}
		} else {
			fmt.Println("No custom attributes found for the user.")
		}
	}

	// Cleanup - delete the test user if we created it
	if isNewUser {
		fmt.Println("Cleaning up - deleting test user...")
		err = oneloginClient.Services.UsersV2.Destroy(*testUser.ID)
		if err != nil {
			log.Printf("WARNING: Failed to delete test user: %v", err)
		} else {
			fmt.Println("Test user deleted successfully.")
		}
	} else {
		fmt.Println("Test complete - keeping existing user.")
	}
}