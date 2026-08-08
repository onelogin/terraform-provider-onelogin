package onelogin

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

// fixtureUniqueToken is replaced in every fixture with a value unique to the
// run. Fixtures in examples/ are documentation as well as test input, so they
// carry readable names — but OneLogin enforces uniqueness on usernames and
// emails within a tenant, and a fixture that hard-codes them passes once and
// fails on every run afterwards with:
//
//	Validation failed: Email must be unique within <tenant>,
//	Username must be unique within <tenant>
//
// A token keeps the examples readable while letting the suite be re-run.
const fixtureUniqueToken = "acctest"

// GetFixture returns the HCL example to be used in an acceptance test, with
// fixtureUniqueToken replaced by a value unique to this run.
func GetFixture(name string, t *testing.T) string {
	t.Helper()
	return getFixtureWithSuffix(name, acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum), t)
}

// GetFixturesWithSuffix returns several fixtures sharing one suffix, and the
// suffix itself, for tests that assert on a value the token appears in.
func GetFixturesWithSuffix(names []string, t *testing.T) ([]string, string) {
	t.Helper()
	suffix := acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)

	out := make([]string, 0, len(names))
	for _, name := range names {
		out = append(out, getFixtureWithSuffix(name, suffix, t))
	}
	return out, suffix
}

// GetFixtures returns several fixtures sharing one suffix, for a test whose
// steps must refer to the same resources. Loading them separately would give
// each step a different suffix, and the second step would replace every
// resource rather than updating it.
func GetFixtures(names []string, t *testing.T) []string {
	t.Helper()
	suffix := acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)

	out := make([]string, 0, len(names))
	for _, name := range names {
		out = append(out, getFixtureWithSuffix(name, suffix, t))
	}
	return out
}

func getFixtureWithSuffix(name, suffix string, t *testing.T) string {
	t.Helper()

	_, filename, _, _ := runtime.Caller(0)
	p := filepath.Join(filepath.Dir(filename), "..", "examples", name)

	rawFile, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("failed to load fixture %s for acceptance test: %v", name, err)
	}

	return stripProviderConfig(strings.ReplaceAll(string(rawFile), fixtureUniqueToken, suffix))
}

// stripProviderConfig removes terraform{} and provider{} blocks from a fixture.
//
// They belong in examples/ — a reader copying one needs to know which provider
// it wants — but an acceptance test is given the provider under test, and a
// required_providers block asking the registry for a version makes Terraform
// try to resolve it instead:
//
//	Error: Inconsistent dependency lock file
//	  provider registry.terraform.io/onelogin/onelogin: required by this
//	  configuration but no version is selected
func stripProviderConfig(config string) string {
	var out strings.Builder
	depth, skipping := 0, false

	for _, line := range strings.Split(config, "\n") {
		trimmed := strings.TrimSpace(line)

		// `provider "` rather than `provider `: the latter also matches the
		// provider meta-argument inside a resource, `provider = onelogin.eu`,
		// and would silently drop that line and the resource's provider with
		// it. A provider block always carries a quoted name.
		if !skipping && (strings.HasPrefix(trimmed, "terraform {") || strings.HasPrefix(trimmed, "provider \"")) {
			skipping = true
			depth = 0
		}
		if !skipping {
			out.WriteString(line + "\n")
			continue
		}

		// Counting braces rather than looking for a closing line on its own:
		// these blocks nest, and required_providers is two deep.
		depth += strings.Count(line, "{") - strings.Count(line, "}")
		if depth <= 0 {
			skipping = false
		}
	}
	return out.String()
}
