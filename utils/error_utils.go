package utils

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// ErrorCategory represents the category of error for consistent formatting
type ErrorCategory string

const (
	ErrorCategoryCreate  ErrorCategory = "CREATE"
	ErrorCategoryRead    ErrorCategory = "READ"
	ErrorCategoryUpdate  ErrorCategory = "UPDATE"
	ErrorCategoryDelete  ErrorCategory = "DELETE"
	ErrorCategoryImport  ErrorCategory = "IMPORT"
	ErrorCategoryUnknown ErrorCategory = "UNKNOWN"
)

// ErrorSeverity represents how severe an error is
type ErrorSeverity string

const (
	ErrorSeverityDebug ErrorSeverity = "DEBUG"
	ErrorSeverityInfo  ErrorSeverity = "INFO"
	ErrorSeverityWarn  ErrorSeverity = "WARN"
	ErrorSeverityError ErrorSeverity = "ERROR"
	ErrorSeverityFatal ErrorSeverity = "FATAL"
)

// FormatError creates a standardized error message
func FormatError(category ErrorCategory, resourceType string, operation string, id string, err error) string {
	return fmt.Sprintf("[%s] Error %s %s (ID: %s): %v", category, operation, resourceType, id, err)
}

// LogAndReturnError logs an error and returns it as a diag.Diagnostics
func LogAndReturnError(
	ctx context.Context,
	severity ErrorSeverity,
	category ErrorCategory,
	resourceType string,
	operation string,
	id string,
	err error,
) diag.Diagnostics {
	message := FormatError(category, resourceType, operation, id, err)

	data := map[string]interface{}{
		"resource_type": resourceType,
		"operation":     operation,
		"id":            id,
		"error":         err.Error(),
	}

	switch severity {
	case ErrorSeverityDebug:
		tflog.Debug(ctx, message, data)
	case ErrorSeverityInfo:
		tflog.Info(ctx, message, data)
	case ErrorSeverityWarn:
		tflog.Warn(ctx, message, data)
	case ErrorSeverityError, ErrorSeverityFatal:
		tflog.Error(ctx, message, data)
	}

	return diag.FromErr(err)
}

// HandleAPIError provides a simplified way to handle API errors
func HandleAPIError(ctx context.Context, err error, category ErrorCategory, resourceType string, id string) diag.Diagnostics {
	return LogAndReturnError(
		ctx,
		ErrorSeverityError,
		category,
		resourceType,
		"calling API",
		id,
		err,
	)
}

// HandleSchemaError provides a simplified way to handle schema errors
func HandleSchemaError(ctx context.Context, err error, category ErrorCategory, resourceType string, id string) diag.Diagnostics {
	return LogAndReturnError(
		ctx,
		ErrorSeverityError,
		category,
		resourceType,
		"processing schema",
		id,
		err,
	)
}

// IsNotFoundError checks if an error represents a 404/not found condition
// This is used to distinguish between "resource doesn't exist" (which should
// remove the resource from state) and actual errors (which should fail the operation)
func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "status: 404") ||
		strings.Contains(errMsg, "not found") ||
		strings.Contains(errMsg, "does not exist")
}
