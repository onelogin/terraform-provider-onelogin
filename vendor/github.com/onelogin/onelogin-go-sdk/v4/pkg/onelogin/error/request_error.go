package error

import (
	"errors"
	"fmt"
	"net/http"
)

// type RequestError struct {
// 	Message string
// }

type RequestError struct {
	context    string
	err        error
	statusCode int
	Message    string
}

func (e RequestError) Error() string {
	return fmt.Sprintf("Request error: %s", e.Message)
}

func NewRequestError(message string) *RequestError {
	return &RequestError{
		Message: message,
	}
}

// ReqErrorWrapper creates a new Request error and returns,
// the pointer to the request error.
// func ReqErrorWrapper(resp *http.Response, context string, err error) error {
// 	code := 0
// 	errToUse := err

// 	if resp != nil {
// 		code = resp.StatusCode
// 	}

// 	if errToUse == nil && code >= http.StatusBadRequest {
// 		errToUse = errors.New(http.StatusText(code))
// 	}

// 	if errToUse == nil {
// 		return nil
// 	}

// 	return &RequestError{
// 		context: context,
// 		err: errToUse,
// 		statusCode: code,
// 	}
// }

func ReqErrorWrapper(resp *http.Response, context string, err error) error {
	if err == nil && resp != nil && resp.StatusCode >= http.StatusBadRequest {
		err = errors.New(http.StatusText(resp.StatusCode))
	}
	if err == nil {
		return nil
	}
	return &RequestError{
		context:    context,
		err:        err,
		statusCode: resp.StatusCode,
	}
}
