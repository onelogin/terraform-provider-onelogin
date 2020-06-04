package utils

import "fmt"

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
