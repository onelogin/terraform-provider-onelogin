package oltypes

// EmptyString is the default value of a string value.
const (
	EmptyString = ""
)

// String takes in a string and returns
// the pointer to it.
func String(val string) *string {
	return &val
}

// GetStringVal takes in a pointer, and if not nil,
// it returns the pointers underlying value.
func GetStringVal(val *string) (string, bool) {
	if val != nil {
		return *val, true
	}

	return EmptyString, false
}
