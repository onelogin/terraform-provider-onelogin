package utils

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceIDFromData extracts and converts an ID from ResourceData
func ResourceIDFromData(d *schema.ResourceData) (string, error) {
	id := d.Id()
	if id == "" {
		return "", fmt.Errorf("ID is empty")
	}
	return id, nil
}

// ResourceIDToInt converts a resource ID to an integer
func ResourceIDToInt(id string) (int32, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return 0, fmt.Errorf("failed to parse ID %s: %w", id, err)
	}
	return int32(intID), nil
}

// StandardDeleteFunc provides a standard delete function pattern
func StandardDeleteFunc(
	ctx context.Context,
	d *schema.ResourceData,
	deleteFunc func(id string) (interface{}, error),
	resourceType string,
) diag.Diagnostics {
	id, err := ResourceIDFromData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("[DELETE] Deleting %s", resourceType), map[string]interface{}{
		"id": id,
	})

	_, err = deleteFunc(id)
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("[ERROR] Error deleting %s", resourceType), map[string]interface{}{
			"id":    id,
			"error": err,
		})
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("[DELETED] Successfully deleted %s", resourceType), map[string]interface{}{
		"id": id,
	})

	d.SetId("")
	return nil
}

// SetResourceFields sets multiple fields from a map to ResourceData
func SetResourceFields(d *schema.ResourceData, data map[string]interface{}, fields []string) {
	for _, field := range fields {
		if value, exists := data[field]; exists {
			d.Set(field, value)
		}
	}
}

// SetNestedFields sets multiple nested fields from a map
func SetNestedFields(ctx context.Context, d *schema.ResourceData, data map[string]interface{}, nestedField string, subFields []string) error {
	if nested, ok := data[nestedField].(map[string]interface{}); ok {
		nestedMap := make(map[string]interface{})

		for _, field := range subFields {
			if value, exists := nested[field]; exists {
				nestedMap[field] = value
			}
		}

		if err := d.Set(nestedField, []map[string]interface{}{nestedMap}); err != nil {
			tflog.Error(ctx, fmt.Sprintf("[ERROR] Error setting %s", nestedField), map[string]interface{}{
				"error": err,
			})
			return err
		}
	}
	return nil
}

// HandleAPIResponse handles standard API response checking pattern
func HandleAPIResponse(ctx context.Context, result interface{}, err error, operation string, resourceType string, id string) (interface{}, diag.Diagnostics) {
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("[ERROR] Error %s %s", operation, resourceType), map[string]interface{}{
			"id":    id,
			"error": err,
		})
		return nil, diag.FromErr(err)
	}

	if result == nil {
		tflog.Info(ctx, fmt.Sprintf("[NOT FOUND] %s with ID %s was not found", resourceType, id), nil)
		return nil, nil
	}

	return result, nil
}
