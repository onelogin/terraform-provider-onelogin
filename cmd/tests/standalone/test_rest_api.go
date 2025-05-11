package standalone

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func TestRestAPI() {
	clientID := os.Getenv("ONELOGIN_CLIENT_ID")
	clientSecret := os.Getenv("ONELOGIN_CLIENT_SECRET")
	baseURL := os.Getenv("ONELOGIN_OAPI_URL")

	fmt.Println("Testing OneLogin API directly")
	fmt.Printf("URL: %s\n", baseURL)

	// Get an access token
	token, err := getToken(baseURL, clientID, clientSecret)
	if err != nil {
		log.Fatalf("Failed to get token: %v", err)
	}
	fmt.Println("Got access token successfully")

	// Create a custom attribute
	attributeName := fmt.Sprintf("Test Attr %d", time.Now().Unix())
	attributeShortname := fmt.Sprintf("test_attr_%d", time.Now().Unix())
	
	fmt.Printf("Creating custom attribute %s (%s)...\n", attributeName, attributeShortname)
	attrID, err := createCustomAttribute(baseURL, token, attributeName, attributeShortname)
	if err != nil {
		log.Printf("Failed to create custom attribute: %v", err)
	} else {
		fmt.Printf("Created custom attribute with ID: %d\n", attrID)
		
		// Delete the custom attribute
		fmt.Printf("Deleting custom attribute %d...\n", attrID)
		err = deleteCustomAttribute(baseURL, token, attrID)
		if err != nil {
			log.Printf("Failed to delete custom attribute: %v", err)
		} else {
			fmt.Println("Custom attribute deleted successfully")
		}
	}
}

// getToken fetches an access token from the OneLogin API
func getToken(baseURL, clientID, clientSecret string) (string, error) {
	url := baseURL + "/auth/oauth2/v2/token"
	payload := []byte(`{"grant_type":"client_credentials"}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "client_id:"+clientID+", client_secret:"+clientSecret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: %s", string(body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	token, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("invalid token response: %s", string(body))
	}

	return token, nil
}

// createCustomAttribute creates a custom attribute
func createCustomAttribute(baseURL, token, name, shortname string) (int, error) {
	url := baseURL + "/api/2/users/custom_attributes"
	payload, _ := json.Marshal(map[string]string{
		"name":      name,
		"shortname": shortname,
	})

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "bearer:"+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode >= 400 {
		return 0, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, err
	}

	id, ok := result["id"].(float64)
	if !ok {
		return 0, fmt.Errorf("failed to parse ID from response: %s", string(body))
	}

	return int(id), nil
}

// deleteCustomAttribute deletes a custom attribute
func deleteCustomAttribute(baseURL, token string, id int) error {
	url := fmt.Sprintf("%s/api/2/users/custom_attributes/%d", baseURL, id)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "bearer:"+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	return nil
}