package models

import (
	"time"
)

type Queryable interface {
	GetKeyValidators() map[string]func(interface{}) bool
}

// validateString checks if the provided value is a string.
func validateString(val interface{}) bool {
	switch v := val.(type) {
	case string:
		return true
	case *string:
		return v != nil
	default:
		return false
	}
}

// validateTime checks if the provided value is a time.Time.
func validateTime(val interface{}) bool {
	switch v := val.(type) {
	case time.Time:
		return true
	case *time.Time:
		return v != nil
	default:
		return false
	}
}

// validateInt checks if the provided value is an int.
func validateInt(val interface{}) bool {
	switch v := val.(type) {
	case int:
		return true
	case *int:
		return v != nil
	default:
		return false
	}
}

// validateBool checks if the provided value is a bool.
func validateBool(val interface{}) bool {
	switch v := val.(type) {
	case bool:
		return true
	case *bool:
		return v != nil
	default:
		return false
	}
}
