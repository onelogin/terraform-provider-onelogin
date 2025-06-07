package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/utilities"
)

func TestCheckHTTPResponseWithForbiddenStatus(t *testing.T) {
	// Create a test server that returns a 403 Forbidden status
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"error": "Forbidden"}`))
	}))
	defer ts.Close()

	// Create a test HTTP response
	resp, err := http.Get(ts.URL)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Call the CheckHTTPResponse function
	_, err = utilities.CheckHTTPResponse(resp)
	
	// Verify that the error message contains information about the status code
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "403")
}

func TestCheckHTTPResponseAndUnmarshalWithForbiddenStatus(t *testing.T) {
	// Create a test server that returns a 403 Forbidden status
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"error": "Forbidden"}`))
	}))
	defer ts.Close()

	// Create a test HTTP response
	resp, err := http.Get(ts.URL)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Call the CheckHTTPResponseAndUnmarshal function
	var result map[string]interface{}
	err = utilities.CheckHTTPResponseAndUnmarshal(resp, &result)
	
	// Verify that the error message contains information about the status code
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "403")
}
