package error

import (
	"errors"
	"fmt"
	"strings"
)

// OneloginError used for any errors.
type OneloginError struct {
	context string
	err     error
}

// OneloginErrorWrapper creates a new OneloginError and returns, if an error is passed in,
// the pointer to the error struct.
func OneloginErrorWrapper(context string, err error) error {
	if err == nil {
		return nil
	}

	return &OneloginError{
		context,
		err,
	}
}

func (olError *OneloginError) Error() string {
	errMsg := ""
	if olError.err != nil {
		errMsg = olError.err.Error()
	}
	return fmt.Sprintf("error: context: [%s], error_message: [%s]", olError.context, errMsg)
}

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
