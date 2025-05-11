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

func TestCreateAttribute_V3() {
	clientID := os.Getenv("ONELOGIN_CLIENT_ID")
	clientSecret := os.Getenv("ONELOGIN_CLIENT_SECRET")
	url := os.Getenv("ONELOGIN_OAPI_URL")

	fmt.Println("Testing direct API call for creating custom attribute")
	fmt.Printf("URL: %s\n", url)

	// Create client
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
	attributeName := fmt.Sprintf("Age-%d", timestamp)
	attributeShortname := fmt.Sprintf("age_years_%d", timestamp)
	
	fmt.Printf("Creating custom attribute: %s (%s)\n", attributeName, attributeShortname)
	
	payload := map[string]interface{}{
		"name":      attributeName,
		"shortname": attributeShortname,
		"position":  3,
	}
	
	result, err := client.CreateCustomAttributes(payload)
	if err != nil {
		log.Printf("Error with map[string]interface{} payload: %v", err)
		
		// Try with string map
		simplePayload := map[string]string{
			"name":      attributeName,
			"shortname": attributeShortname,
		}
		
		result, err = client.CreateCustomAttributes(simplePayload)
		if err != nil {
			log.Fatalf("Error with map[string]string payload: %v", err)
		}
	}
	
	fmt.Printf("Success! Created attribute: %+v\n", result)
}