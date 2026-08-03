package utils

import (
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
		err      error
		expected bool
	}{
		"nil error is not a 404": {
			err:      nil,
			expected: false,
		},
		"404 status error is recognised": {
			err:      fmt.Errorf("request failed with status: 404"),
			expected: true,
		},
		"502 status error is not a 404": {
			err:      fmt.Errorf("request failed with status: 502"),
			expected: false,
		},
		"unrelated error is not a 404": {
			err:      fmt.Errorf("connection refused"),
			expected: false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, IsNotFoundError(tc.err))
		})
	}
}
