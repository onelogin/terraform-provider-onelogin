package customerrors

import "fmt"

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
