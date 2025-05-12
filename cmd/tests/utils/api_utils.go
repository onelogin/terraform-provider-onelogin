package utils

import (
	"fmt"
	"os"
	"strings"
)

// GetAPICredentials returns the OneLogin API credentials and URL
// It supports both new subdomain-based and legacy URL-based configurations
func GetAPICredentials() (clientID, clientSecret, baseURL, subdomain string) {
	clientID = os.Getenv("ONELOGIN_CLIENT_ID")
	clientSecret = os.Getenv("ONELOGIN_CLIENT_SECRET")

	// Check for subdomain first (preferred)
	subdomain = os.Getenv("ONELOGIN_SUBDOMAIN")
	if subdomain != "" {
		// Generate URL from subdomain
		if !strings.HasSuffix(subdomain, ".onelogin.com") {
			baseURL = fmt.Sprintf("https://%s.onelogin.com", subdomain)
		} else {
			baseURL = fmt.Sprintf("https://%s", subdomain)
		}
		return
	}

	// Fallback to ONELOGIN_OAPI_URL (legacy)
	baseURL = os.Getenv("ONELOGIN_OAPI_URL")
	if baseURL != "" {
		fmt.Println("WARNING: Using legacy ONELOGIN_OAPI_URL environment variable.")
		fmt.Println("Please consider using ONELOGIN_SUBDOMAIN instead.")

		// Try to extract subdomain from URL
		if strings.Contains(baseURL, "onelogin.com") {
			parts := strings.Split(baseURL, "//")
			if len(parts) > 1 {
				domainParts := strings.Split(parts[1], ".")
				if len(domainParts) > 0 {
					subdomain = domainParts[0]
				}
			}
		}
	} else if subdomain == "" {
		// If neither subdomain nor URL is set, default to the subdomain.onelogin.com format
		fmt.Println("WARNING: Neither ONELOGIN_SUBDOMAIN nor ONELOGIN_OAPI_URL is set.")
		fmt.Println("Using subdomain to construct API URL.")
		baseURL = fmt.Sprintf("https://%s.onelogin.com", subdomain)
	}

	// Ensure URL has a scheme
	if baseURL != "" && !strings.HasPrefix(baseURL, "http") {
		baseURL = "https://" + baseURL
	}

	return
}

// SetupSDKEnvironment sets up the environment variables required by the OneLogin SDK
func SetupSDKEnvironment(clientID, clientSecret, subdomain, url string) {
	os.Setenv("ONELOGIN_CLIENT_ID", clientID)
	os.Setenv("ONELOGIN_CLIENT_SECRET", clientSecret)
	
	if subdomain != "" {
		os.Setenv("ONELOGIN_SUBDOMAIN", subdomain)
	}
	
	if url != "" {
		os.Setenv("ONELOGIN_OAPI_URL", url)
	}
	
	// Default timeout
	os.Setenv("ONELOGIN_TIMEOUT", "60")
}