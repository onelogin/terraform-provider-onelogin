package standalone

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	ol "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/api"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/authentication"
)

func main() {
	clientID := os.Getenv("ONELOGIN_CLIENT_ID")
	clientSecret := os.Getenv("ONELOGIN_CLIENT_SECRET")
	url := os.Getenv("ONELOGIN_OAPI_URL")

	fmt.Println("Testing custom attribute creation with SDK v4.3.0")
	fmt.Printf("URL: %s\n", url)

	// Initialize v4 client
	authenticator := authentication.NewAuthenticator("api")
	timeoutDuration := time.Second * 60
	apiClient := &api.Client{
		HttpClient: &http.Client{Timeout: timeoutDuration},
		Auth:       authenticator,
		OLdomain:   url,
		Timeout:    timeoutDuration,
	}

	client := &ol.OneloginSDK{
		Client: apiClient,
	}

	// Make the credentials available to the authenticator
	os.Setenv("ONELOGIN_CLIENT_ID", clientID)
	os.Setenv("ONELOGIN_CLIENT_SECRET", clientSecret)
	os.Setenv("ONELOGIN_OAPI_URL", url)

	// Test creating a custom attribute
	timestamp := time.Now().Unix()
	attributeName := fmt.Sprintf("Test-%d", timestamp)
	attributeShortname := fmt.Sprintf("test_%d", timestamp)
	
	fmt.Printf("Creating custom attribute: %s (%s)\n", attributeName, attributeShortname)
	
	payload := map[string]interface{}{
		"name":      attributeName,
		"shortname": attributeShortname,
	}
	
	result, err := client.CreateCustomAttributes(payload)
	if err != nil {
		log.Fatalf("Error creating custom attribute: %v", err)
	}
	
	fmt.Printf("Success! Created attribute: %+v\n", result)
	
	// Try to get all custom attributes to verify
	fmt.Println("\nListing all custom attributes:")
	attributes, err := client.GetCustomAttributes()
	if err != nil {
		log.Fatalf("Error getting custom attributes: %v", err)
	}
	
	fmt.Printf("All custom attributes: %+v\n", attributes)
}