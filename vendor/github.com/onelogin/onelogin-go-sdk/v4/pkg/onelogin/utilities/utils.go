package utilities

import (
	"encoding/base64"
	"fmt"
	"log"
	"regexp"
	"strings"
)

// ReplaceSpecialChar replaces any non-alphanumeric character in a string with a character of choice
func ReplaceSpecialChar(str string, rep string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")
	return reg.ReplaceAllString(str, rep)
}

// ToSnakeCase turns camel/pascal case strings into snake_case strings
func ToSnakeCase(str string) string {
	matchFirstCap := regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap := regexp.MustCompile("([a-z0-9])([A-Z])")
	matchAllSpaces := regexp.MustCompile("(\\s)")
	cleanUpHack := regexp.MustCompile("i_ds")

	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	snake = matchAllSpaces.ReplaceAllString(snake, "_")
	snake = strings.ToLower(snake)
	snake = cleanUpHack.ReplaceAllString(snake, "ids")

	return snake
}

// OneOf takes the name of a field to analyze, the field's value, and a list of allowed values.
// Returns an error if a field's value is not in the list of legal values.
func OneOf(key string, v string, opts []string) error {
	isValid := false
	for _, o := range opts {
		isValid = v == o
		if isValid {
			break
		}
	}
	if !isValid {
		return fmt.Errorf("%s must be one of %v, got: %s", key, opts, v)
	}
	return nil
}

/*

	String Encoding Convenience Methods

*/

// IsEncoded takes a string and returns whether or not the string is base64 encoded
func IsEncoded(s string) bool {
	B64EncodedRegex := *regexp.MustCompile(`^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$`)
	return B64EncodedRegex.MatchString(s)
}

// EncodeString takes a multi-line string representation of a function body and
// encodes it to a base64 encoded string for use with OneLogin API.
// Returns the original encoded string if it happens to be base64 encoded
func EncodeString(s string) string {
	if IsEncoded(s) {
		return s
	}
	encoded := base64.StdEncoding.EncodeToString([]byte(s))
	return encoded
}

// DecodeString takes a base64 encoded string and returns the decoeded result.
// returns the original string if it is not base64 encoded
func DecodeString(s string) string {
	if !IsEncoded(s) {
		return s
	}
	decodedBytes, e := base64.StdEncoding.DecodeString(s)
	if e != nil {
		log.Println("Unexpected decoding error:", e)
		return s
	}
	return string(decodedBytes)
}
