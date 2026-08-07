package utils

import (
	"strings"
	"testing"
)

// TestOneOfValue covers the guard that keeps a validator from panicking.
//
// Terraform should only ever hand a string to a TypeString attribute, so the
// assertion "cannot" fail. But a panic inside a validator takes the whole
// provider process down and reports a stack trace rather than the offending
// value, which is a poor trade against one type check.
func TestOneOfValue(t *testing.T) {
	opts := []string{"all", "any"}

	t.Run("accepts an allowed value", func(t *testing.T) {
		if _, errs := OneOfValue("match", "all", opts); len(errs) > 0 {
			t.Fatalf("expected %q to be accepted, got %v", "all", errs)
		}
	})

	t.Run("rejects a disallowed value", func(t *testing.T) {
		_, errs := OneOfValue("match", "sometimes", opts)

		if len(errs) == 0 {
			t.Fatal("expected a disallowed value to be rejected")
		}
		if !strings.Contains(errs[0].Error(), "sometimes") {
			t.Fatalf("expected the offending value in the error, got %v", errs[0])
		}
	})

	t.Run("reports a non-string rather than panicking", func(t *testing.T) {
		// The whole point. Each of these would have panicked on a bare
		// assertion, taking the provider down with it.
		for _, val := range []interface{}{42, true, nil, []string{"all"}, map[string]string{}} {
			func() {
				defer func() {
					if r := recover(); r != nil {
						t.Fatalf("panicked on %#v: %v", val, r)
					}
				}()

				_, errs := OneOfValue("match", val, opts)
				if len(errs) == 0 {
					t.Fatalf("expected %#v to be reported as an error", val)
				}
				if !strings.Contains(errs[0].Error(), "expected a string") {
					t.Fatalf("expected a type error for %#v, got %v", val, errs[0])
				}
			}()
		}
	})

	t.Run("names the key and the type it got", func(t *testing.T) {
		// A validator error is all the practitioner sees, so it has to say
		// which attribute and what arrived.
		_, errs := OneOfValue("match", 42, opts)

		if len(errs) == 0 {
			t.Fatal("expected an error")
		}
		msg := errs[0].Error()
		if !strings.Contains(msg, "match") || !strings.Contains(msg, "int") {
			t.Fatalf("expected the key and the type in %q", msg)
		}
	})
}
