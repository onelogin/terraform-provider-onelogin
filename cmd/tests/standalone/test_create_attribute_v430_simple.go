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

func TestCreateAttribute_V430_Simple() {
	clientID := os.Getenv("ONELOGIN_CLIENT_ID")
	clientSecret := os.Getenv("ONELOGIN_CLIENT_SECRET")
	url := os.Getenv("ONELOGIN_OAPI_URL")

	fmt.Println("Testing simple custom attribute creation (name & shortname only)")
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

	// Now create a custom attribute - ONLY using name and shortname
	timestamp := fmt.Sprintf("%d", os.Getpid())
	attributeName := fmt.Sprintf("SimpleAttr-%s", timestamp)
	attributeShortname := fmt.Sprintf("simple_attr_%s", timestamp)

	fmt.Printf("Creating custom attribute: %s (%s)\n", attributeName, attributeShortname)

	payload := map[string]string{
		"name":      attributeName,
		"shortname": attributeShortname,
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
		log.Fatalf("Error creating custom attribute: %v", err)
	}
	defer createResp.Body.Close()

	createBody, err := ioutil.ReadAll(createResp.Body)
	if err != nil {
		log.Fatalf("Error reading create response: %v", err)
	}

	fmt.Printf("Response status: %s\n", createResp.Status)
	fmt.Printf("Response body: %s\n", string(createBody))

	if createResp.StatusCode >= 200 && createResp.StatusCode < 300 {
		fmt.Println("Successfully created custom attribute!")
	} else {
		fmt.Printf("Failed to create custom attribute (status %d)\n", createResp.StatusCode)
	}
}