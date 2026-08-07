package utils

import (
	"fmt"
	"strings"
)

// OneOf returns errors or warnings for the given key/value pair if the value is not
// included in the given list of allowed options
func OneOf(key string, v string, opts []string) (warns []string, errs []error) {
	isValid := false
	for _, o := range opts {
		isValid = v == o
		if isValid {
			break
		}
	}
	if !isValid {
		errs = append(errs, fmt.Errorf("%s must be one of %v, got: %s", key, opts, v))
	}
	return
}

// OneOfValue is OneOf for a value arriving from a schema.SchemaValidateFunc,
// where it is an interface{} and asserting the string bare would panic.
//
// Terraform should only ever hand a string to a TypeString attribute, so the
// assertion "cannot" fail. But a panic inside a validator takes the whole
// provider process down and reports a stack trace rather than the bad value,
// which is a poor trade against one type check. Corrupted state and schema
// changes are the realistic routes in.
func OneOfValue(key string, val interface{}, opts []string) (warns []string, errs []error) {
	v, ok := val.(string)
	if !ok {
		return nil, []error{fmt.Errorf("%s: expected a string, got %T", key, val)}
	}
	return OneOf(key, v, opts)
}

func ParseNestedResourceImportId(id string) (string, string, error) {
	parts := strings.SplitN(id, ":", 2)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected attribute1:attribute2", id)
	}

	return parts[0], parts[1], nil
}
