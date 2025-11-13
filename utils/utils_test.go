package utils

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOneOf(t *testing.T) {
	validOpts := []string{"SHA-1", "SHA-256", "SHA-348", "SHA-512"}
	tests := map[string]struct {
		InputKey       string
		InputValue     string
		ExpectedOutput []error
	}{
		"no errors on valid input": {
			InputKey:       "signature_algorithm",
			InputValue:     "SHA-1",
			ExpectedOutput: nil,
		},
		"errors on invalid input": {
			InputKey:       "signature_algorithm",
			InputValue:     "asdf",
			ExpectedOutput: []error{fmt.Errorf("signature_algorithm must be one of %v, got: %s", validOpts, "asdf")},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, errs := OneOf(test.InputKey, test.InputValue, validOpts)
			assert.Equal(t, test.ExpectedOutput, errs)
		})
	}
}

func TestIsNotFoundError(t *testing.T) {
	tests := map[string]struct {
		InputError     error
		ExpectedResult bool
	}{
		"returns false for nil error": {
			InputError:     nil,
			ExpectedResult: false,
		},
		"returns true for error containing 404": {
			InputError:     errors.New("request failed with status: 404"),
			ExpectedResult: true,
		},
		"returns true for error containing 'not found'": {
			InputError:     errors.New("resource not found"),
			ExpectedResult: true,
		},
		"returns true for error containing 'Not Found' (capitalized)": {
			InputError:     errors.New("Resource Not Found"),
			ExpectedResult: true,
		},
		"returns true for error containing 'does not exist'": {
			InputError:     errors.New("app does not exist"),
			ExpectedResult: true,
		},
		"returns true for error containing 'Does Not Exist' (capitalized)": {
			InputError:     errors.New("App Does Not Exist"),
			ExpectedResult: true,
		},
		"returns false for 500 error": {
			InputError:     errors.New("request failed with status: 500"),
			ExpectedResult: false,
		},
		"returns false for network timeout": {
			InputError:     errors.New("network timeout"),
			ExpectedResult: false,
		},
		"returns false for invalid credentials": {
			InputError:     errors.New("invalid credentials"),
			ExpectedResult: false,
		},
		"returns false for generic error": {
			InputError:     errors.New("something went wrong"),
			ExpectedResult: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := IsNotFoundError(test.InputError)
			assert.Equal(t, test.ExpectedResult, result, "IsNotFoundError should return %v for error: %v", test.ExpectedResult, test.InputError)
		})
	}
}
