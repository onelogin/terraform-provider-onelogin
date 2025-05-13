package oltypes

// Bool takes in a boolean and returns
// the pointer to it.
func Bool(val bool) *bool {
	return &val
}

// GetBoolVal takes in a pointer, and if not nil,
// it returns the pointers underlying value.
func GetBoolVal(val *bool) (bool, bool) {
	if val != nil {
		return *val, true
	}

	return false, false
}
