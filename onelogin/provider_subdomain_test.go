package onelogin

import "testing"

// TestSubdomainFromURL pins the shapes the provider accepts for ONELOGIN_API_URL.
// The SDK still asks for a subdomain even though every call goes to the URL, and
// this is the only thing that supplies it.
func TestSubdomainFromURL(t *testing.T) {
	for _, tc := range []struct {
		url  string
		want string
	}{
		{"https://chicken.onelogin.com", "chicken"},
		{"http://chicken.onelogin.com", "chicken"},
		{"chicken.onelogin.com", "chicken"},
		// A regional API host carries no tenant subdomain, so the SDK gets a
		// placeholder to satisfy its own requirement.
		{"https://api.us.onelogin.com", "dummy"},
		{"https://api.eu.onelogin.com", "dummy"},
	} {
		got, err := subdomainFromURL(tc.url)
		if err != nil {
			t.Fatalf("%s: unexpected error %v", tc.url, err)
		}
		if got != tc.want {
			t.Fatalf("%s: expected %q, got %q", tc.url, tc.want, got)
		}
	}

	for _, url := range []string{"https://api.ap.onelogin.com", ""} {
		if got, err := subdomainFromURL(url); err == nil {
			t.Fatalf("%q: expected an error, got %q", url, got)
		}
	}
}
