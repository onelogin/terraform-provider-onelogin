package appparametersschema

import (
	"reflect"
	"testing"
)

// TestParameterValues covers the shape sent to the API.
//
// The API stores exactly what it is given and returns it unchanged, so this is
// what decides whether an upgrade rewrites a parameter that has been a string
// since it was created.
func TestParameterValues(t *testing.T) {
	for _, tc := range []struct {
		name string
		in   interface{}
		want interface{}
	}{
		{"one value goes as a bare string", []interface{}{"memberOf"}, "memberOf"},
		{"several go as an array", []interface{}{"a", "b"}, []string{"a", "b"}},
		{"an empty list is nothing at all", []interface{}{}, nil},
		{"an empty string is a value, not nothing", []interface{}{""}, ""},
		{"a value containing a comma stays one value", []interface{}{"a,b"}, "a,b"},
		{"nil is nothing", nil, nil},
		{"a stray non-list is nothing", "not a list", nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := parameterValues(tc.in); !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("expected %#v, got %#v", tc.want, got)
			}
		})
	}
}

// TestParameterValueList covers reading the field back. The same field is a
// string on one parameter and an array on another, decided by nothing but what
// was last written to it.
func TestParameterValueList(t *testing.T) {
	for _, tc := range []struct {
		name string
		in   interface{}
		want []interface{}
	}{
		{"a string becomes a list of one", "memberOf", []interface{}{"memberOf"}},
		{"an empty string is kept", "", []interface{}{""}},
		{"an array comes through in order", []interface{}{"a", "b"}, []interface{}{"a", "b"}},
		{"an empty array stays empty", []interface{}{}, []interface{}{}},
		// Not a shape the API produces -- JSON gives []interface{} -- but the
		// model field is interface{}, so a caller building a Parameter in Go
		// can put one here. Dropping it would lose the value silently.
		{"a []string from a Go caller is read", []string{"a", "b"}, []interface{}{"a", "b"}},
		{"a null is nothing, so the key can be left out", nil, nil},
		{"an unexpected shape is nothing rather than a panic", 42, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := parameterValueList(tc.in); !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("expected %#v, got %#v", tc.want, got)
			}
		})
	}
}
