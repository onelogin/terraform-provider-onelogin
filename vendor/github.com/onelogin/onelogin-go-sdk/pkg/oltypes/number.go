package oltypes

// Int64 takes in a int64 and returns
// the pointer to it.
func Int64(val int64) *int64 {
	return &val
}

// GetInt64Val takes in a pointer, and if not nil,
// it returns the pointers underlying value.
func GetInt64Val(val *int64) (int64, bool) {
	if val != nil {
		return *val, true
	}

	return 0, false
}

// Int32 takes in a int32 and returns
// the pointer to it.
func Int32(val int32) *int32 {
	return &val
}

// GetInt32Val takes in a pointer, and if not nil,
// it returns the pointers underlying value.
func GetInt32Val(val *int32) (int32, bool) {
	if val != nil {
		return *val, true
	}

	return 0, false
}
