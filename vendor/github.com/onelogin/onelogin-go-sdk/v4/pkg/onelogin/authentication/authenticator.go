package authentication

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	olError "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/error"
)

const (
	TkPath     string = "/auth/oauth2/v2/token"
	RevokePath string = "/auth/oauth2/revoke"
)

type Authenticator struct {
	accessToken string
	subdomain   string
}

func NewAuthenticator(subdomain string) *Authenticator {
	return &Authenticator{subdomain: subdomain}
}

// BaseURL is the root every request is sent to, without a trailing slash.
//
// ONELOGIN_API_URL wins when it is set. Deriving the host from the subdomain
// alone assumes every tenant lives at <subdomain>.onelogin.com, which is true
// of production and of nothing else: a development or staging deployment, or a
// tenant on a custom domain, cannot be reached that way. Worse, the derived
// host is a real tenant belonging to someone, so a caller pointed elsewhere had
// its credentials posted to production rather than being told the host was
// unsupported.
//
// With no ONELOGIN_API_URL set the subdomain is used exactly as before, so
// existing callers are unaffected.
func BaseURL(subdomain string) string {
	apiURL := strings.TrimSpace(os.Getenv("ONELOGIN_API_URL"))
	if apiURL == "" {
		return fmt.Sprintf("https://%s.onelogin.com", subdomain)
	}

	// A bare host is the likely way to get this wrong, and http.NewRequest
	// rejects a URL with no scheme rather than assuming one.
	if !strings.Contains(apiURL, "://") {
		apiURL = "https://" + apiURL
	}

	// Reduced to scheme and authority. Every path this SDK builds is absolute
	// -- /auth/oauth2/v2/token, /api/2/... -- and is appended to whatever this
	// returns, so a value carrying a path would produce
	// https://host/auth/auth/oauth2/v2/token rather than replacing anything.
	// Such a value cannot work, so trimming it is the difference between one
	// wrong URL and every wrong URL.
	if parsed, err := url.Parse(apiURL); err == nil && parsed.Host != "" {
		return parsed.Scheme + "://" + parsed.Host
	}

	// Unparseable: hand it back tidied rather than silently substituting a
	// different host, which is the failure this function exists to prevent.
	return strings.TrimRight(apiURL, "/")
}

func (a *Authenticator) GenerateToken() error {
	// Read & Check environment variables
	clientID := os.Getenv("ONELOGIN_CLIENT_ID")
	if len(clientID) == 0 {
		return olError.NewAuthenticationError("Missing ONELOGIN_CLIENT_ID Env Variable")
	}
	//fmt.Println("clientID", clientID)
	clientSecret := os.Getenv("ONELOGIN_CLIENT_SECRET")
	if len(clientSecret) == 0 {
		return olError.NewAuthenticationError("Missing ONELOGIN_CLIENT_SECRET Env Variable")
	}

	// Construct the authentication URL
	authURL := BaseURL(a.subdomain) + TkPath

	// Create authentication request payload
	data := map[string]string{
		"grant_type": "client_credentials",
	}

	// Convert payload to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return olError.NewSerializationError("Unable to convert payload to JSON")
	}

	// Create HTTP request
	req, err := http.NewRequest(http.MethodPost, authURL, strings.NewReader(string(jsonData)))
	if err != nil {
		return olError.NewRequestError("Failed to create authentication request")
	}

	// Add authorization header with base64-encoded credentials
	encodedCredentials := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", clientID, clientSecret)))
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", encodedCredentials))
	req.Header.Add("Content-Type", "application/json")

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return olError.NewRequestError("Failed to send authentication request")
	}

	// Parse the authentication response
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return olError.NewSerializationError("Failed to read authentication response")
	}

	// Check if authentication failed
	if resp.StatusCode != http.StatusOK {
		return olError.NewAuthenticationError("Authentication failed")
	}

	// Extract access token from the response
	accessToken, ok := result["access_token"].(string)
	if !ok {
		return olError.NewAuthenticationError("Authentication Failed at Endpoint")
	}
	// Store access token
	a.accessToken = accessToken

	return nil
}

func (a *Authenticator) RevokeToken(token *string) error {
	// Read environment variables
	clientID := os.Getenv("ONELOGIN_CLIENT_ID")
	clientSecret := os.Getenv("ONELOGIN_CLIENT_SECRET")

	// Check if required environment variables are missing
	if clientID == "" || clientSecret == "" {
		return errors.New("missing client ID, client secret, or subdomain")
	}

	// Construct the revoke URL. This previously carried no scheme at all, which
	// http.NewRequest rejects, so revocation could not have worked.
	revokeURL := BaseURL(a.subdomain) + RevokePath

	// Create revoke request payload
	data := struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: *token,
	}

	// Convert payload to JSON
	jsonData, err := json.Marshal(data) // #nosec G117 -- access_token field name is required by the OAuth2 revocation spec
	if err != nil {
		return fmt.Errorf("failed to create revocation request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", revokeURL, strings.NewReader(string(jsonData)))
	if err != nil {
		return fmt.Errorf("failed to create revocation request: %w", err)
	}

	// Add authorization header with base64-encoded credentials
	encodedCredentials := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", clientID, clientSecret)))
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", encodedCredentials))
	req.Header.Add("Content-Type", "application/json")

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to revoke: %w", err)
	}

	// Check if revocation failed
	if resp.StatusCode != http.StatusOK {
		return olError.NewAuthenticationError("Revocation failed")
	}

	// Success condition feedback
	fmt.Println("Revocation successful")

	return nil
}

func (a *Authenticator) GetToken() (string, error) {
	return a.accessToken, nil
}
