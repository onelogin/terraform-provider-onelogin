package standalone

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func TestCreateAttribute_V430_Direct() {
	clientID := os.Getenv("ONELOGIN_CLIENT_ID")
	clientSecret := os.Getenv("ONELOGIN_CLIENT_SECRET")
	url := os.Getenv("ONELOGIN_OAPI_URL")

	fmt.Println("Testing direct HTTP API call for custom attributes with SDK v4.3.0")
	fmt.Printf("URL: %s\n", url)

	// First get an access token
	tokenURL := fmt.Sprintf("%s/auth/oauth2/v2/token", url)
	fmt.Printf("Token URL: %s\n", tokenURL)

	tokenReq, err := http.NewRequest("POST", tokenURL, bytes.NewBuffer([]byte("grant_type=client_credentials")))
	if err != nil {
		log.Fatalf("Error creating token request: %v", err)
	}

	tokenReq.SetBasicAuth(clientID, clientSecret)
	tokenReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	tokenResp, err := client.Do(tokenReq)
	if err != nil {
		log.Fatalf("Error getting token: %v", err)
	}
	defer tokenResp.Body.Close()

	tokenBody, err := ioutil.ReadAll(tokenResp.Body)
	if err != nil {
		log.Fatalf("Error reading token response: %v", err)
	}

	var tokenData map[string]interface{}
	err = json.Unmarshal(tokenBody, &tokenData)
	if err != nil {
		log.Fatalf("Error parsing token response: %v", err)
	}

	accessToken, ok := tokenData["access_token"].(string)
	if !ok {
		log.Fatalf("No access token in response: %s", string(tokenBody))
	}

	fmt.Println("Got access token successfully")

	// Get all custom attributes to see what exists
	fmt.Println("\nGetting all custom attributes:")
	getURL := fmt.Sprintf("%s/api/2/users/custom_attributes", url)
	getReq, err := http.NewRequest("GET", getURL, nil)
	if err != nil {
		log.Fatalf("Error creating get request: %v", err)
	}
	getReq.Header.Set("Authorization", "Bearer "+accessToken)
	getResp, err := client.Do(getReq)
	if err != nil {
		log.Fatalf("Error getting custom attributes: %v", err)
	}
	defer getResp.Body.Close()
	
	getBody, err := ioutil.ReadAll(getResp.Body)
	if err != nil {
		log.Fatalf("Error reading get response: %v", err)
	}
	
	fmt.Printf("Existing custom attributes: %s\n\n", string(getBody))

	// Now create a custom attribute
	timestamp := fmt.Sprintf("%d", os.Getpid())
	attributeName := fmt.Sprintf("TestAttribute-%s", timestamp)
	attributeShortname := fmt.Sprintf("test_attr_%s", timestamp)

	fmt.Printf("Creating custom attribute: %s (%s)\n", attributeName, attributeShortname)

	// Try with a variety of payloads
	payloads := []map[string]interface{}{
		{
			"name":      attributeName,
			"shortname": attributeShortname,
		},
		{
			"name":       attributeName,
			"shortname":  attributeShortname,
			"user_field": "custom_attribute",
		},
		{
			"name":          attributeName,
			"shortname":     attributeShortname,
			"sort_order":    1,
		},
		{
			"name":          attributeName,
			"shortname":     attributeShortname,
			"sort_order":    1,
			"user_field":    "custom_attribute",
		},
	}

	for i, payload := range payloads {
		fmt.Printf("\nTrying payload %d: %+v\n", i+1, payload)
		
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			log.Fatalf("Error marshalling payload: %v", err)
		}

		createURL := fmt.Sprintf("%s/api/2/users/custom_attributes", url)
		fmt.Printf("Create URL: %s\n", createURL)

		createReq, err := http.NewRequest("POST", createURL, bytes.NewBuffer(payloadBytes))
		if err != nil {
			log.Fatalf("Error creating custom attribute request: %v", err)
		}

		createReq.Header.Set("Authorization", "Bearer "+accessToken)
		createReq.Header.Set("Content-Type", "application/json")

		createResp, err := client.Do(createReq)
		if err != nil {
			log.Printf("Error creating custom attribute: %v", err)
			continue
		}
		
		createBody, _ := ioutil.ReadAll(createResp.Body)
		createResp.Body.Close()

		fmt.Printf("Response status: %s\n", createResp.Status)
		fmt.Printf("Response body: %s\n", string(createBody))

		if createResp.StatusCode >= 200 && createResp.StatusCode < 300 {
			fmt.Println("Successfully created custom attribute!")
			break
		} else {
			fmt.Printf("Failed to create custom attribute with payload %d (status %d)\n", i+1, createResp.StatusCode)
		}
	}
}