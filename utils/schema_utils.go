package utils

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// CommonFields returns schema fields that are common across multiple resources
func CommonFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"updated_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

// BooleanField returns a schema definition for a boolean field
func BooleanField(required bool, computed bool, description string, defaultValue bool) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeBool,
		Required:    required,
		Optional:    !required && !computed,
		Computed:    computed,
		Description: description,
		Default:     defaultValue,
	}
}

// StringField returns a schema definition for a string field
func StringField(required bool, computed bool, description string, defaultValue string) *schema.Schema {
	schema := &schema.Schema{
		Type:        schema.TypeString,
		Required:    required,
		Optional:    !required && !computed,
		Computed:    computed,
		Description: description,
	}

	if defaultValue != "" {
		schema.Default = defaultValue
	}

	return schema
}

// IntField returns a schema definition for an integer field
func IntField(required bool, computed bool, description string, defaultValue int) *schema.Schema {
	schema := &schema.Schema{
		Type:        schema.TypeInt,
		Required:    required,
		Optional:    !required && !computed,
		Computed:    computed,
		Description: description,
	}

	if defaultValue != 0 {
		schema.Default = defaultValue
	}

	return schema
}

// StringListField returns a schema definition for a list of strings
func StringListField(required bool, computed bool, description string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Required:    required,
		Optional:    !required && !computed,
		Computed:    computed,
		Description: description,
		Elem:        &schema.Schema{Type: schema.TypeString},
	}
}

// SetField returns a schema definition for a set field
func SetField(required bool, computed bool, description string, elem interface{}) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Required:    required,
		Optional:    !required && !computed,
		Computed:    computed,
		Description: description,
		Elem:        elem,
	}
}

// MapField returns a schema definition for a map field
func MapField(required bool, computed bool, description string, elem interface{}) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeMap,
		Required:    required,
		Optional:    !required && !computed,
		Computed:    computed,
		Description: description,
		Elem:        elem,
	}
}

// MergeSchemas combines multiple schema maps into one
func MergeSchemas(schemas ...map[string]*schema.Schema) map[string]*schema.Schema {
	result := make(map[string]*schema.Schema)

	for _, s := range schemas {
		for k, v := range s {
			result[k] = v
		}
	}

	return result
}
