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

func TestCreateAttribute_V430_WithUserfield() {
	clientID := os.Getenv("ONELOGIN_CLIENT_ID")
	clientSecret := os.Getenv("ONELOGIN_CLIENT_SECRET")
	url := os.Getenv("ONELOGIN_OAPI_URL")

	fmt.Println("Testing custom attribute creation with user_field")
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

	// Now create a custom attribute with ONLY name, shortname, and user_field
	timestamp := fmt.Sprintf("%d", os.Getpid())
	attributeName := fmt.Sprintf("UserFieldAttr-%s", timestamp)
	attributeShortname := fmt.Sprintf("user_field_attr_%s", timestamp)

	fmt.Printf("Creating custom attribute: %s (%s)\n", attributeName, attributeShortname)

	// Try different values for user_field
	userFieldValues := []string{"", "custom_attribute", "custom_fields", "user_fields", "true"}
	
	for i, userFieldValue := range userFieldValues {
		fmt.Printf("\nTrying with user_field = '%s'\n", userFieldValue)
		
		payload := map[string]string{
			"name":       attributeName,
			"shortname":  attributeShortname,
		}
		
		// Add user_field if it's not empty
		if userFieldValue != "" {
			payload["user_field"] = userFieldValue
		}
		
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
			fmt.Printf("Successfully created custom attribute with user_field = '%s'!\n", userFieldValue)
			break
		} else {
			fmt.Printf("Failed to create custom attribute with user_field = '%s' (status %d)\n", userFieldValue, createResp.StatusCode)
			
			// Try a different approach if we're still getting 400s
			if i == len(userFieldValues)-1 && createResp.StatusCode == 400 {
				fmt.Println("\n\nTrying one last approach with JSON structure:")
				
				// Try with a boolean value
				jsonPayload := fmt.Sprintf(`{
					"name": "%s",
					"shortname": "%s",
					"user_field": true
				}`, attributeName, attributeShortname)
				
				createReq, err := http.NewRequest("POST", createURL, bytes.NewBuffer([]byte(jsonPayload)))
				if err != nil {
					log.Fatalf("Error creating final custom attribute request: %v", err)
				}

				createReq.Header.Set("Authorization", "Bearer "+accessToken)
				createReq.Header.Set("Content-Type", "application/json")

				createResp, err := client.Do(createReq)
				if err != nil {
					log.Fatalf("Error in final create attempt: %v", err)
				}
				
				createBody, _ := ioutil.ReadAll(createResp.Body)
				createResp.Body.Close()

				fmt.Printf("Final response status: %s\n", createResp.Status)
				fmt.Printf("Final response body: %s\n", string(createBody))
			}
		}
	}
}