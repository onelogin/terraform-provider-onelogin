package customerrors

import (
	"errors"
	"strings"
)

// Errors in services
var (
	ErrValueMissing = errors.New("A required parameter was not given")
)

// StackErrors amalgamates a list of error messages into a single comma-separated error message.
func StackErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	stackedErrors := make([]string, len(errs))
	for i, e := range errs {
		if e != nil {
			stackedErrors[i] = e.Error()
		}
	}
	return errors.New(strings.Join(stackedErrors, ", "))
}
