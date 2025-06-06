package utilities

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	olerror "github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/error"
)

// receive http response, check error code status, if good return json of resp.Body
// else return error
func CheckHTTPResponse(resp *http.Response) (any, error) {
	// Handle 204 No Content responses - this is a success but with no content
	if resp.StatusCode == http.StatusNoContent {
		return map[string]any{"status": "success"}, nil
	}

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusAccepted {
		// If the status code is 403 (Forbidden), it likely means the API credentials are read-only
		if resp.StatusCode == http.StatusForbidden {
			return nil, fmt.Errorf("request failed with status: %d - API credentials may have read-only access", resp.StatusCode)
		}
		return nil, fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Close the response body
	err = resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close response body: %w", err)
	}

	// Try to unmarshal the response body into a map[string]any or []any
	var data any
	bodyStr := string(body)
	//log.Printf("Response body: %s\n", bodyStr)
	if strings.HasPrefix(bodyStr, "[") {
		var slice []any
		err = json.Unmarshal(body, &slice)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal response body into []any: %w", err)
		}
		data = slice
	} else if strings.HasPrefix(bodyStr, "{") {
		var dict map[string]any
		err = json.Unmarshal(body, &dict)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal response body into map[string]any: %w", err)
		}
		data = dict
	} else {
		data = bodyStr
	}

	//log.Printf("Response body unmarshaled successfully: %v\n", data)
	return data, nil
}

// CheckHTTPResponseAndUnmarshal checks the HTTP response and unmarshals the response body into the target struct
func CheckHTTPResponseAndUnmarshal(resp *http.Response, target any) error {
	// Handle 204 No Content responses - this is a success but with no content
	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusAccepted {
		// If the status code is 403 (Forbidden), it likely means the API credentials are read-only
		if resp.StatusCode == http.StatusForbidden {
			return fmt.Errorf("request failed with status: %d - API credentials may have read-only access", resp.StatusCode)
		}
		return fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Close the response body
	err = resp.Body.Close()
	if err != nil {
		return fmt.Errorf("failed to close response body: %w", err)
	}

	// Unmarshal the response body into the target struct
	err = json.Unmarshal(body, target)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return nil
}

func BuildAPIPath(parts ...any) (string, error) {
	var path string
	for _, part := range parts {
		switch p := part.(type) {
		case string:
			path += "/" + p
		case int:
			path += fmt.Sprintf("/%d", p)
		case int32:
			path += fmt.Sprintf("/%d", p)
		case int64:
			path += fmt.Sprintf("/%d", p)
		default:
			// Handle other types if needed
			return path, olerror.NewSDKError("Unsupported path type")
		}
	}

	// Check if the path is valid
	if !IsPathValid(path) {
		return path, olerror.NewSDKError("Invalid path")
	}

	return path, nil
}

// AddQueryToPath adds the model as a JSON-encoded query parameter to the path and returns the new path.
func AddQueryToPath(path string, query any) (string, error) {
	if query == nil {
		return path, nil
	}

	// Convert query parameters to URL-encoded string
	values, err := queryToValues(query)
	if err != nil {
		return "", err
	}

	// Append query parameters to path
	if values.Encode() != "" {
		path += "?" + values.Encode()
	}

	return path, nil
}

func queryToValues(query any) (url.Values, error) {
	values := url.Values{}

	// Convert query parameters to URL-encoded string using reflection
	if query != nil {
		// First, get the json field names from struct tags
		queryBytes, err := json.Marshal(query)
		if err != nil {
			return nil, err
		}

		// Unmarshal to map[string]interface{} to handle all types of values
		var data map[string]any
		if err := json.Unmarshal(queryBytes, &data); err != nil {
			return nil, err
		}

		// Add each field to query parameters
		for key, value := range data {
			if value != nil {
				// Handle different value types
				switch v := value.(type) {
				case string:
					values.Set(key, v)
				case float64:
					values.Set(key, fmt.Sprintf("%v", v))
				case []any:
					// For arrays, convert to comma-separated string
					if len(v) > 0 {
						// Convert array to comma-separated string
						strItems := make([]string, len(v))
						for i, item := range v {
							strItems[i] = fmt.Sprintf("%v", item)
						}
						values.Set(key, strings.Join(strItems, ","))
					}
				default:
					// Convert other types to string
					values.Set(key, fmt.Sprintf("%v", v))
				}
			}
		}
	}

	return values, nil
}
