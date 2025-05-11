package v4

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	ol "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/api"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/authentication"
	"github.com/onelogin/terraform-provider-onelogin/cmd/tests/utils"
)

func TestV4CustomAttributes() {
	clientID, clientSecret, url, subdomain := utils.GetAPICredentials()

	fmt.Println("Testing OneLogin Custom Attributes with v4 SDK")
	fmt.Printf("URL: %s\n", url)

	// Set up environment variables for the v4 SDK
	utils.SetupSDKEnvironment(clientID, clientSecret, subdomain, url)

	// Initialize v4 client
	authenticator := authentication.NewAuthenticator(subdomain)
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

	// Test getting all custom attributes
	fmt.Println("Getting all custom attributes...")
	attributes, err := client.GetCustomAttributes()
	if err != nil {
		log.Printf("Error getting custom attributes: %v", err)
	} else {
		fmt.Printf("Custom attributes: %+v\n", attributes)
	}

	// Test creating a custom attribute
	fmt.Println("Creating a custom attribute...")
	attributeName := fmt.Sprintf("Test Attr %d", time.Now().Unix())
	attributeShortname := fmt.Sprintf("test_attr_%d", time.Now().Unix())

	// Print request details for debugging
	fmt.Printf("Creating attribute with name: %s, shortname: %s\n", attributeName, attributeShortname)

	// Follow the exact structure used in the resource
	userFieldPayload := map[string]interface{}{
		"name":      attributeName,
		"shortname": attributeShortname,
	}

	// Wrap in user_field object as required by API
	payload := map[string]interface{}{
		"user_field": userFieldPayload,
	}

	// Convert to JSON for logging
	payloadJSON, _ := json.Marshal(payload)
	fmt.Printf("Request payload: %s\n", string(payloadJSON))
	
	result, err := client.CreateCustomAttributes(payload)
	if err != nil {
		log.Fatalf("Error creating custom attribute: %v", err)
	}
	
	fmt.Printf("Created custom attribute: %+v\n", result)
	
	// Extract the custom attribute ID
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		log.Fatalf("Failed to parse custom attribute creation response")
	}
	
	id, ok := resultMap["id"].(float64)
	if !ok {
		log.Fatalf("Failed to extract custom attribute ID from response")
	}
	
	attributeID := int(id)
	fmt.Printf("Custom attribute ID: %d\n", attributeID)
	
	// Cleanup
	fmt.Printf("Deleting custom attribute %d...\n", attributeID)
	_, err = client.DeleteCustomAttributes(attributeID)
	if err != nil {
		log.Printf("Error deleting custom attribute: %v", err)
	} else {
		fmt.Println("Custom attribute deleted successfully")
	}
}