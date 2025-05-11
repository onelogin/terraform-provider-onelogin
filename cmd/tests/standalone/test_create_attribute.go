package standalone

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

func TestCreateAttribute() {
	clientID, clientSecret, url, subdomain := utils.GetAPICredentials()

	fmt.Println("Testing creating Age custom attribute with v4 SDK")
	fmt.Printf("URL: %s, Subdomain: %s\n", url, subdomain)

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

	// Test creating a custom attribute
	fmt.Println("Creating Age custom attribute...")
	timestamp := time.Now().Unix()
	attributeName := fmt.Sprintf("Age_%d", timestamp)
	attributeShortname := fmt.Sprintf("age_years_%d", timestamp)

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

	// Now verify that we can get the custom attribute
	fmt.Println("Getting all custom attributes to verify...")
	attributes, err := client.GetCustomAttributes()
	if err != nil {
		log.Printf("Error getting custom attributes: %v", err)
	} else {
		fmt.Printf("Custom attributes: %+v\n", attributes)
		
		// Try to find our newly created attribute
		attrList, ok := attributes.([]interface{})
		if !ok {
			log.Fatalf("Invalid custom attributes response format")
		}
		
		found := false
		for _, attr := range attrList {
			attrMap, ok := attr.(map[string]interface{})
			if !ok {
				continue
			}
			
			if attrMap["shortname"] == attributeShortname {
				fmt.Printf("Found our attribute: %+v\n", attrMap)
				found = true
				break
			}
		}
		
		if !found {
			fmt.Println("Could not find our newly created attribute in the list!")
		}
	}
}