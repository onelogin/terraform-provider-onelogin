package utils

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin"
)

// AddRequestResponseLogging adds request/response logging to the SDK client
func AddRequestResponseLogging(ctx context.Context, client *onelogin.OneloginSDK) {
	// Just log that debug mode is enabled
	tflog.Info(ctx, "[DEBUG] Debug logging enabled for OneLogin SDK")
}

// LogAPIRequest logs an API request for debugging purposes
func LogAPIRequest(ctx context.Context, method string, url string, body interface{}) {
	// Format the body if it's provided
	var bodyJSON string
	if body != nil {
		bodyBytes, err := json.MarshalIndent(body, "", "  ")
		if err == nil {
			bodyJSON = string(bodyBytes)
		} else {
			bodyJSON = fmt.Sprintf("%v", body)
		}
	}

	tflog.Info(ctx, "[DEBUG] API Request", map[string]interface{}{
		"method": method,
		"url":    url,
		"body":   bodyJSON,
	})
}

// LogAPIResponse logs an API response for debugging purposes
func LogAPIResponse(ctx context.Context, statusCode int, response interface{}, err error) {
	if err != nil {
		tflog.Error(ctx, "[DEBUG] API Error", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Format the response
	var responseJSON string
	if response != nil {
		respBytes, jsonErr := json.MarshalIndent(response, "", "  ")
		if jsonErr == nil {
			responseJSON = string(respBytes)
		} else {
			responseJSON = fmt.Sprintf("%v", response)
		}
	}

	logFunc := tflog.Info
	if statusCode >= 400 {
		logFunc = tflog.Error
	}

	logFunc(ctx, "[DEBUG] API Response", map[string]interface{}{
		"status_code": statusCode,
		"body":        responseJSON,
	})
}
