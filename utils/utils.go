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

func ParseNestedResourceImportId(id string) (string, string, error) {
	parts := strings.SplitN(id, ":", 2)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected attribute1:attribute2", id)
	}

	return parts[0], parts[1], nil
}
