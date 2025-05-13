package error

import (
	"fmt"
)

type SerializationError struct {
	Message string
}

func (e SerializationError) Error() string {
	return fmt.Sprintf("Serialization error: %s", e.Message)
}

func NewSerializationError(message string) error {
	return SerializationError{
		Message: message,
	}
}
