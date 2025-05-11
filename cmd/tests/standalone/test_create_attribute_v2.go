package standalone

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	ol "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/api"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/authentication"
)

func TestCreateAttribute_V2() {
	clientID := os.Getenv("ONELOGIN_CLIENT_ID")
	clientSecret := os.Getenv("ONELOGIN_CLIENT_SECRET")
	url := os.Getenv("ONELOGIN_OAPI_URL")

	fmt.Println("Testing creating Age custom attribute with v4 SDK - Improved Version")
	fmt.Printf("URL: %s\n", url)

	// Set up environment variables for the v4 SDK
	os.Setenv("ONELOGIN_CLIENT_ID", clientID)
	os.Setenv("ONELOGIN_CLIENT_SECRET", clientSecret)

	// Extract subdomain from URL
	var subdomain string
	if url != "" {
		subdomain = "api" // Use api as subdomain when custom URL is provided
		os.Setenv("ONELOGIN_OAPI_URL", url)
	} else {
		subdomain = "api"
	}
	os.Setenv("ONELOGIN_SUBDOMAIN", subdomain)
	os.Setenv("ONELOGIN_TIMEOUT", "60")

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

	// First, print the raw client credentials to check
	fmt.Println("API Client ID (length):", len(os.Getenv("ONELOGIN_CLIENT_ID")))
	fmt.Println("API Client Secret (length):", len(os.Getenv("ONELOGIN_CLIENT_SECRET")))

	// For debugging, let's check our credentials by getting a list of users
	fmt.Println("Testing credentials by getting a list of users...")
	_, err := client.GetUsers(nil)
	if err != nil {
		log.Printf("Error getting users: %v", err)
	} else {
		fmt.Println("Successfully retrieved users - credentials are working!")
	}

	// Get current attributes to see what exists
	fmt.Println("Getting current custom attributes...")
	currentAttrs, err := client.GetCustomAttributes()
	if err != nil {
		log.Printf("Error getting current custom attributes: %v", err)
	} else {
		// Pretty print the attributes
		jsonData, _ := json.MarshalIndent(currentAttrs, "", "  ")
		fmt.Println(string(jsonData))
	}

	// Test creating a custom attribute - use a timestamp to make it unique
	timestamp := time.Now().Unix()
	fmt.Println("Creating Age custom attribute...")
	attributeName := fmt.Sprintf("Age-%d", timestamp)
	attributeShortname := fmt.Sprintf("age_years_%d", timestamp)
	
	payload := map[string]interface{}{
		"name":      attributeName,
		"shortname": attributeShortname,
		"position":  3,
	}
	
	// Print the payload for verification
	jsonPayload, _ := json.MarshalIndent(payload, "", "  ")
	fmt.Printf("Request payload: %s\n", string(jsonPayload))
	
	result, err := client.CreateCustomAttributes(payload)
	if err != nil {
		log.Printf("Error creating custom attribute: %v", err)
		
		// Try with a simpler payload
		fmt.Println("Trying with a simpler payload...")
		simplePayload := map[string]string{
			"name":      attributeName,
			"shortname": attributeShortname,
		}
		
		jsonPayload, _ := json.MarshalIndent(simplePayload, "", "  ")
		fmt.Printf("Simple payload: %s\n", string(jsonPayload))
		
		result, err = client.CreateCustomAttributes(simplePayload)
		if err != nil {
			log.Fatalf("Still failed to create custom attribute: %v", err)
		}
	}
	
	// If we get here, creation succeeded
	fmt.Printf("Created custom attribute: %+v\n", result)
	
	// Pretty print the result
	jsonResult, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsonResult))
	
	// Extract the custom attribute ID if possible
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		log.Printf("Warning: Couldn't parse result as map[string]interface{}")
	} else {
		if id, ok := resultMap["id"].(float64); ok {
			attributeID := int(id)
			fmt.Printf("Custom attribute ID: %d\n", attributeID)
		} else {
			fmt.Println("Could not extract ID from result")
		}
	}
}