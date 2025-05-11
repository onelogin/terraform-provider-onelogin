package v1

import (
	"fmt"
	"log"
	"os"

	"github.com/onelogin/onelogin-go-sdk/pkg/client"
)

func TestV1Users() {
	clientID := os.Getenv("ONELOGIN_CLIENT_ID")
	clientSecret := os.Getenv("ONELOGIN_CLIENT_SECRET")
	region := os.Getenv("ONELOGIN_REGION")
	url := os.Getenv("ONELOGIN_OAPI_URL")

	fmt.Println("Testing OneLogin Users API with v1 SDK")
	fmt.Printf("URL: %s, Region: %s\n", url, region)

	// Create client
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

	// Get users
	users, err := oneloginClient.Services.UsersV2.Query(nil)
	if err != nil {
		log.Fatalf("Failed to query users: %v", err)
	}

	fmt.Printf("Successfully retrieved %d users\n", len(users))
	if len(users) > 0 {
		fmt.Printf("First user: ID=%d, Username=%s\n", *users[0].ID, *users[0].Username)
	}
}