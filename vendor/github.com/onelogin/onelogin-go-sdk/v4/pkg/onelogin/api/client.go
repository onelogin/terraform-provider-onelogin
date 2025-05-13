package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/authentication"
	olerror "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/error"
	mod "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	utl "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/utilities"
)

// Client represents the API client.
type Client struct {
	HttpClient HTTPClient                    // HTTPClient interface for making HTTP requests
	Auth       *authentication.Authenticator // Authenticator for managing authentication
	OLdomain   string                        // OneLogin domain
	Timeout    time.Duration
}

// HTTPClient is an interface that defines the Do method for making HTTP requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Authenticator is an interface that defines the GetToken method for retrieving authentication tokens.
type Authenticator interface {
	GetToken() (string, error)
	NewAuthenticator() *authentication.Authenticator
}

// NewClient creates a new instance of the API client.
func NewClient() (*Client, error) {
	subdomain := os.Getenv("ONELOGIN_SUBDOMAIN")
	if subdomain == "" {
		return nil, olerror.NewSDKError("ONELOGIN_SUBDOMAIN environment variable is required")
	}

	old := fmt.Sprintf("https://%s.onelogin.com", subdomain)
	authenticator := authentication.NewAuthenticator(subdomain)
	timeoutStr := os.Getenv("ONELOGIN_TIMEOUT")
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil || timeout <= 0 {
		timeout = 10
	}
	timeoutDuration := time.Second * time.Duration(timeout)

	err = authenticator.GenerateToken()
	if err != nil {
		return nil, err
	}
	return &Client{
		HttpClient: &http.Client{
			Timeout: timeoutDuration,
		},
		Auth:     authenticator,
		OLdomain: old,
	}, nil
}

// newRequest creates a new HTTP request with the specified method, path, query parameters, and request body.
func (c *Client) newRequest(method string, path *string, queryParams mod.Queryable, body io.Reader) (*http.Request, error) {
	return c.newRequestWithContext(context.Background(), method, path, queryParams, body)
}

// newRequestWithContext creates a new HTTP request with context and the specified method, path, query parameters, and request body.
func (c *Client) newRequestWithContext(ctx context.Context, method string, path *string, queryParams mod.Queryable, body io.Reader) (*http.Request, error) {
	p, err := utl.AddQueryToPath(*path, queryParams)
	if err != nil {
		return nil, err
	}
	log.Println("Path:", p)
	// Parse the OneLogin domain and path
	u, err := url.Parse(c.OLdomain + p)
	if err != nil {
		return nil, err
	}

	// Create a new HTTP request with context
	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, err
	}

	// Get authentication token
	log.Println("Getting authentication token...")
	tk, err := c.Auth.GetToken()
	if err != nil {
		log.Println("Error getting authentication token:", err)
		return nil, olerror.NewAuthenticationError("Access Token Retrieval Error")
	}
	log.Println("Authentication token retrieved successfully.")

	// Set request headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tk))
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

// Get sends a GET request to the specified path with the given query parameters.
func (c *Client) Get(path *string, queryParams mod.Queryable) (*http.Response, error) {
	return c.GetWithContext(context.Background(), path, queryParams)
}

// GetWithContext sends a GET request with the given context, path, and query parameters.
func (c *Client) GetWithContext(ctx context.Context, path *string, queryParams mod.Queryable) (*http.Response, error) {
	req, err := c.newRequestWithContext(ctx, http.MethodGet, path, queryParams, http.NoBody)
	if err != nil {
		return nil, err
	}

	return c.sendRequest(req)
}

// Delete sends a DELETE request to the specified path with the given query parameters.
func (c *Client) Delete(path *string) (*http.Response, error) {
	return c.DeleteWithContext(context.Background(), path)
}

// DeleteWithContext sends a DELETE request with the given context and path.
func (c *Client) DeleteWithContext(ctx context.Context, path *string) (*http.Response, error) {
	req, err := c.newRequestWithContext(ctx, http.MethodDelete, path, nil, http.NoBody)
	if err != nil {
		return nil, err
	}

	return c.sendRequest(req)
}

// DeleteWithBody sends a DELETE request to the specified path with the given request body.
func (c *Client) DeleteWithBody(path *string, body interface{}) (*http.Response, error) {
	return c.DeleteWithBodyAndContext(context.Background(), path, body)
}

// DeleteWithBodyAndContext sends a DELETE request with the given context, path, and request body.
func (c *Client) DeleteWithBodyAndContext(ctx context.Context, path *string, body interface{}) (*http.Response, error) {
	// Convert request body to JSON
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := c.newRequestWithContext(ctx, http.MethodDelete, path, nil, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	return c.sendRequest(req)
}

// Post sends a POST request to the specified path with the given request body.
func (c *Client) Post(path *string, body interface{}) (*http.Response, error) {
	return c.PostWithContext(context.Background(), path, body)
}

// PostWithContext sends a POST request with the given context, path, and request body.
func (c *Client) PostWithContext(ctx context.Context, path *string, body interface{}) (*http.Response, error) {
	// Convert request body to JSON
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := c.newRequestWithContext(ctx, http.MethodPost, path, nil, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	return c.sendRequest(req)
}

// Put sends a PUT request to the specified path with the given request body.
func (c *Client) Put(path *string, body interface{}) (*http.Response, error) {
	return c.PutWithContext(context.Background(), path, body)
}

// PutWithContext sends a PUT request with the given context, path, and request body.
func (c *Client) PutWithContext(ctx context.Context, path *string, body interface{}) (*http.Response, error) {
	// Convert request body to JSON
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := c.newRequestWithContext(ctx, http.MethodPut, path, nil, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	return c.sendRequest(req)
}

// sendRequest sends the specified HTTP request and returns the HTTP response.
func (c *Client) sendRequest(req *http.Request) (*http.Response, error) {
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Check for API errors
	if resp.StatusCode == http.StatusUnauthorized {
		// Regenerate the token and reattempt the request
		err := c.Auth.GenerateToken()
		if err != nil {
			return nil, olerror.NewAuthenticationError("Failed to refresh access token")
		}

		// Retry the request
		resp, err = c.HttpClient.Do(req)
		if err != nil {
			return nil, err
		}
	}

	return resp, nil
}
